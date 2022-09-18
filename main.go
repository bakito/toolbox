package main

import (
	"bytes"
	"fmt"
	"github.com/bakito/toolbox/version"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/bakito/toolbox/pkg/extract"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/cavaliergopher/grab/v3"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

var (
	aliases = map[string][]string{
		"amd64":   {"x86_64", "64"},
		"windows": {"win", "win64"},
		"linux":   {"linux64"},
	}
)

func main() {

	log.Printf("toolbox v%s", version.Version)

	tb, err := readToolbox()
	if err != nil {
		panic(err)
	}

	if tb.Target == "" {
		tb.Target = "./tools"
	}
	if tb.Aliases != nil {
		aliases = *tb.Aliases
	}

	if _, err := os.Stat(tb.Target); err != nil {
		if os.IsNotExist(err) {
			if tb.CreateTarget == nil || *tb.CreateTarget {
				log.Printf("Creating target dir %q\n", tb.Target)
				_ = os.MkdirAll(tb.Target, 0700)
			} else {
				log.Fatalf("Target dir %q does not exist and may not be created.\n", tb.Target)
			}
		} else {
			panic(err)
		}
	}

	tmp, err := os.MkdirTemp("", "toolbox")
	if err != nil {
		panic(err)
	}

	defer func() { _ = os.RemoveAll(tmp) }()

	client := resty.New()

	for _, tool := range tb.Tools {
		log.Printf("Download %s\n", tool.Name)
		var ghr *types.GithubRelease
		if tool.Github != "" {
			ghr = &types.GithubRelease{}

			ghc := client.R().
				SetResult(ghr).
				SetHeader("Accept", "application/json")
			if t, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
				ghc = ghc.SetAuthToken(t)
			}
			_, err := ghc.
				Get(tool.LatestURL())

			if err != nil {
				panic(err)
			}

			if tool.Version == "" {
				tool.Version = ghr.TagName
				log.Printf("Latest Version: %s", tool.Version)
			}
		}

		if tool.DownloadURL != "" {
			if strings.HasPrefix(tool.Version, "http") {
				resp, err := client.R().
					EnableTrace().
					Get(tool.Version)
				if err != nil {
					panic(err)
				}
				tool.Version = string(resp.Body())
				log.Printf("Latest Version: %s", tool.Version)
			}

			if err := fetchTool(tmp, tool.Name, tool.Name, parseTemplate(tool.DownloadURL, tool.Version), tb.Target); err != nil {
				panic(err)
			}
		} else if ghr != nil {
			matching := findMatching(tool.Name, ghr.Assets)
			if matching != nil {
				if err := fetchTool(tmp, tool.Name, tool.Name, matching.BrowserDownloadUrl, tb.Target); err != nil {
					panic(err)
				}
			}
			for _, add := range tool.Additional {
				matching := findMatching(add, ghr.Assets)
				if matching != nil {
					if err := fetchTool(tmp, add, add, matching.BrowserDownloadUrl, tb.Target); err != nil {
						panic(err)
					}
				}
			}
		}
		println()
	}
}

func findMatching(toolName string, assets []types.Asset) *types.Asset {
	var matching []*types.Asset
	for i := range assets {
		a := assets[i]
		if strings.Contains(a.Name, toolName) && matches(runtime.GOOS, a.Name) {
			matching = append(matching, &a)
		}
	}
	sort.Slice(matching, func(i, j int) bool {
		mi := matches(runtime.GOARCH, matching[i].Name)
		mj := matches(runtime.GOARCH, matching[j].Name)

		if mi == mj {
			return true
		}

		return mi
	})
	if len(matching) == 0 {
		return nil
	}
	return matching[0]
}

func parseTemplate(templ string, version string) string {
	ut, err := template.New("url").Parse(templ)
	if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	if err := ut.Execute(&b, templateData(version)); err != nil {
		panic(err)
	}
	return b.String()
}

func templateData(version string) map[string]string {
	return map[string]string{
		"Version": version,
		"OS":      runtime.GOOS,
		"Arch":    runtime.GOARCH,
		"ArchBIT": fmt.Sprintf("%d", strconv.IntSize),
		"FileExt": defaultFileExtension(),
	}
}

func fetchTool(tmp string, remoteToolName string, trueToolName string, url string, targetDir string) error {
	dir := fmt.Sprintf("%s/%s", tmp, remoteToolName)
	paths := strings.Split(url, "/")
	fileName := paths[len(paths)-1]
	path := fmt.Sprintf("%s/%s", dir, fileName)
	log.Printf("Downloading %s", url)
	if err := downloadFile(path, url); err != nil {
		return err
	}
	extracted, err := extract.File(path, dir)
	if err != nil {
		return err
	}
	if !extracted {
		remoteToolName = fileName
	}
	ok, err := copyTool(dir, remoteToolName, targetDir, trueToolName)
	if err != nil {
		return err
	}
	if !ok {
		log.Printf("WARN: Could not find: %s", remoteToolName)
	}
	return nil
}

func copyTool(dir string, fileName string, targetDir string, targetName string) (bool, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	var dirs []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
		if file.Name() == binaryName(fileName) || file.Name() == binaryName(fmt.Sprintf("%s_%s_%s", fileName, runtime.GOOS, runtime.GOARCH)) {

			if err := copyFile(dir, file, targetDir, targetName); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	for _, d := range dirs {
		ok, err := copyTool(filepath.Join(dir, d.Name()), fileName, targetDir, targetName)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func copyFile(dir string, file os.DirEntry, targetDir string, targetName string) error {
	from, err := os.Open(filepath.Join(dir, file.Name()))
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(filepath.Join(targetDir, binaryName(targetName)), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer to.Close()
	log.Printf("Copy %s to %s", from.Name(), to.Name())
	_, err = to.ReadFrom(from)
	return err
}

func readToolbox() (*types.Toolbox, error) {
	b, err := os.ReadFile(".toolbox.yaml")
	if err != nil {
		return nil, err
	}
	tb := &types.Toolbox{}
	err = yaml.Unmarshal(b, tb)
	if err != nil {
		return nil, err
	}
	for name, tool := range tb.Tools {
		if tool.Name == "" {
			tool.Name = name
		}
	}
	return tb, nil
}

func matches(info string, name string) bool {
	if strings.Contains(name, info) {
		return true
	}

	for _, a := range aliases[info] {
		if strings.Contains(name, a) {
			return true
		}
	}

	return false
}

func downloadFile(path string, url string) (err error) {
	resp, err := grab.Get(path, url)
	if err != nil {
		return err
	}

	log.Printf("Download saved to %s", resp.Filename)
	return nil
}
