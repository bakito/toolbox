package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/bakito/toolbox/pkg/extract"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/cavaliergopher/grab/v3"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

func main() {

	tb, err := readToolbox()
	if err != nil {
		panic(err)
	}

	if tb.Target == "" {
		tb.Target = "tools"
	}

	tmp, err := os.MkdirTemp("", "toolbox")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(tmp)

	client := resty.New()

	for _, tool := range tb.Tools {
		log.Printf("Download %s\n", tool.Name)
		var ghr *types.GithubRelease
		if tool.Github != "" {
			ghr = &types.GithubRelease{}
			_, err := client.R().
				EnableTrace().
				SetResult(ghr).
				Get(tool.LatestURL())
			//Get(os.Args[0])

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
			ut, err := template.New("url").Parse(tool.DownloadURL)
			if err != nil {
				panic(err)
			}

			var b bytes.Buffer
			if err := ut.Execute(&b, map[string]string{"Version": tool.Version, "OS": runtime.GOOS, "Arch": runtime.GOARCH}); err != nil {
				panic(err)
			}

			if err := fetchTool(tmp, tool.Name, b.String(), tb.Target); err != nil {
				panic(err)
			}
		} else if ghr != nil {

			for _, a := range ghr.Assets {
				if strings.Contains(a.Name, tool.Name) && matches(runtime.GOOS, a.Name) && matches(runtime.GOARCH, a.Name) {
					if err := fetchTool(tmp, tool.Name, a.BrowserDownloadUrl, tb.Target); err != nil {
						panic(err)
					}
				}
				for _, add := range tool.Additional {

					if strings.Contains(a.Name, add) && matches(runtime.GOOS, a.Name) && matches(runtime.GOARCH, a.Name) {
						if err := fetchTool(tmp, add, a.BrowserDownloadUrl, tb.Target); err != nil {
							panic(err)
						}
					}
				}
			}
		}
	}
}

func fetchTool(tmp string, toolName string, url string, targetDir string) error {
	dir := fmt.Sprintf("%s/%s", tmp, toolName)
	paths := strings.Split(url, "/")
	fileName := paths[len(paths)-1]
	path := fmt.Sprintf("%s/%s", dir, fileName)
	log.Printf("Downloading %s", url)
	if err := downloadFile(path, url); err != nil {
		return err
	}
	if err := extract.File(path, dir); err != nil {
		return err
	}
	if _, err := copyTool(dir, toolName, targetDir); err != nil {
		return err
	}
	return nil
}

func copyTool(dir string, tool string, targetDir string) (bool, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	var dirs []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
		if file.Name() == tool || (runtime.GOOS == "windows" && file.Name() == tool+".exe") {

			if err := copyFile(dir, file, targetDir); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	for _, d := range dirs {
		ok, err := copyTool(filepath.Join(dir, d.Name()), tool, targetDir)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func copyFile(dir string, file os.DirEntry, targetDir string) error {
	from, err := os.Open(filepath.Join(dir, file.Name()))
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.OpenFile(filepath.Join(targetDir, file.Name()), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
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

// https://get.helm.sh/helm-v3.9.4-darwin-amd64.tar.gz

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

var (
	aliases = map[string][]string{"amd64": {"x86_64"}}
)

func downloadFile(path string, url string) (err error) {
	resp, err := grab.Get(path, url)
	if err != nil {
		return err
	}

	log.Printf("Download saved to %s", resp.Filename)
	return nil
}
