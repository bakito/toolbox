package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/bakito/toolbox/pkg/extract"
	"github.com/bakito/toolbox/pkg/github"
	"github.com/bakito/toolbox/pkg/http"
	"github.com/bakito/toolbox/pkg/quietly"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/bakito/toolbox/version"
	"github.com/cavaliergopher/grab/v3"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	flagConfig          = "config"
	toolboxConfFile     = ".toolbox.yaml"
	toolboxVersionsFile = ".toolbox-versions.yaml"
)

var (
	aliases = map[string][]string{
		"amd64":   {"x86_64", "64", "64bit"},
		"windows": {"win", "win64"},
		"linux":   {"linux64"},
	}
	stopAliases = map[string][]string{
		"amd64":   {"arm"},
		"windows": {"darwin"},
	}

	// fetchCmd represents the fetch command
	fetchCmd = &cobra.Command{
		Use:   "fetch",
		Short: "Fetch all tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := cmd.Flags().GetString(flagConfig)
			if err != nil {
				return err
			}
			return fetch(cfg)
		},
	}
)

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().StringP(flagConfig, "c", "",
		"The config file to be used. (default 1. '.toolbox.yaml' current dir, 2. '~/.toolbox.yaml')")
}

func fetch(cfgFile string) error {
	log.Printf("🧰 toolbox %s", version.Version)

	client := resty.New()

	tbRel, err := github.LatestRelease(client, "bakito/toolbox", true)
	if err != nil {
		return err
	}
	if tbRel.TagName != version.Version {
		log.Printf("🌟 A new toolbox version is available %s (current: %s)\n", tbRel.TagName, version.Version)
	}

	tb, err := readToolbox(cfgFile)
	if err != nil {
		return err
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
				_ = os.MkdirAll(tb.Target, 0o700)
			} else {
				return fmt.Errorf("target dir %q does not exist and may not be created", tb.Target)
			}
		} else {
			return err
		}
	}

	ver, err := readVersions(tb.Target)
	if err != nil {
		return err
	}

	tmp, err := os.MkdirTemp("", "toolbox")
	if err != nil {
		return err
	}

	defer func() { _ = os.RemoveAll(tmp) }()

	tools := tb.GetTools()
	for i := range tools {
		tool := tools[i]
		if err := handleTool(client, ver, tmp, tb, tool); err != nil {
			return err
		}
	}

	// save versions
	tv := tb.Versions()
	var b bytes.Buffer
	env := yaml.NewEncoder(&b)
	env.SetIndent(2)

	if err := env.Encode(&tv); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(tb.Target, toolboxVersionsFile), b.Bytes(), 0o600)
}

func handleTool(client *resty.Client, ver map[string]string, tmp string, tb *types.Toolbox, tool *types.Tool) error {
	log.Printf("⚙️ Processing %s\n", tool.Name)
	defer func() { println() }()
	var ghr *types.GithubRelease
	var err error
	configVersion := tool.Version
	currentVersion := ver[tool.Name]
	if tool.Github != "" {
		if configVersion == "" {
			ghr, err = github.LatestRelease(client, tool.Github, false)
		} else {
			ghr, err = github.Release(client, tool.Github, configVersion, false)
		}
		if err != nil {
			return err
		}

		if tool.Version == "" {
			tool.Version = ghr.TagName
			if currentVersion != "" && tool.Version != currentVersion {
				log.Printf("Latest Version: %s (current: %s)", tool.Version, currentVersion)
			} else {
				log.Printf("Latest Version: %s", tool.Version)
			}
		}
	}
	if tool.Version == currentVersion {
		if configVersion != "" {
			log.Printf("✅ Skipping since already configured version %s\n", configVersion)
		} else {
			log.Printf("✅ Skipping since already latest version\n")
		}
		return nil
	}

	if tool.DownloadURL != "" {
		return downloadFromURL(client, ver, tmp, tb, tool)
	} else if ghr != nil {
		return downloadViaGithub(tool, ghr, tmp, tb)
	}
	return nil
}

func downloadViaGithub(tool *types.Tool, ghr *types.GithubRelease, tmp string, tb *types.Toolbox) error {
	matching := findMatching(tool.Name, ghr.Assets)
	tool.CouldNotBeFound = true
	if matching != nil {
		tool.CouldNotBeFound = false
		if err := fetchTool(tmp, tool.Name, tool.Name, matching.BrowserDownloadURL, tb.Target); err != nil {
			return err
		}
	}
	for _, add := range tool.Additional {
		matching := findMatching(add, ghr.Assets)
		if matching != nil {
			tool.CouldNotBeFound = false
			if err := fetchTool(tmp, add, add, matching.BrowserDownloadURL, tb.Target); err != nil {
				return err
			}
		}
	}
	if tool.CouldNotBeFound {
		log.Printf("❌ Couldn't find a file here!\n")
	}
	return nil
}

func downloadFromURL(client *resty.Client, ver map[string]string, tmp string, tb *types.Toolbox, tool *types.Tool) error {
	currentVersion := ver[tool.Name]
	if strings.HasPrefix(tool.Version, "http") {
		resp, err := client.R().
			EnableTrace().
			Get(tool.Version)
		if err != nil {
			return nil
		}
		tool.Version = string(resp.Body())

		if currentVersion != "" && tool.Version != currentVersion {
			log.Printf("Latest Version: %s (current: %s)", tool.Version, currentVersion)
		} else {
			log.Printf("Latest Version: %s", tool.Version)
		}
	}

	if tool.Version == currentVersion {
		log.Printf("✅ Skipping since already latest version\n")
		return nil
	}
	return fetchTool(tmp, tool.Name, tool.Name, parseTemplate(tool.DownloadURL, tool.Version), tb.Target)
}

func findMatching(toolName string, assets []types.Asset) *types.Asset {
	var matching []*types.Asset
	for i := range assets {
		a := assets[i]
		if strings.Contains(a.Name, toolName) && matches(runtime.GOOS, a.Name) && !strings.HasSuffix(a.Name, "sum") {
			matching = append(matching, &a)
		}
	}
	sort.Slice(matching, func(i, j int) bool {
		mi := matches(runtime.GOARCH, matching[i].Name)
		mj := matches(runtime.GOARCH, matching[j].Name)

		if mi == mj {
			mi = strings.Contains(matching[i].Name, runtime.GOARCH)
			mj = strings.Contains(matching[j].Name, runtime.GOARCH)
		}
		if mi == mj {
			// prefer non archive files
			mi = !strings.Contains(matching[i].Name, ".")
			mj = !strings.Contains(matching[j].Name, ".")
		}
		if mi == mj {
			// prefer non archive files
			mi = strings.HasSuffix(matching[i].Name, defaultFileExtension())
			mj = strings.HasSuffix(matching[j].Name, defaultFileExtension())
		}
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
		"Version":    version,
		"VersionNum": strings.TrimPrefix(version, "v"),
		"OS":         runtime.GOOS,
		"Arch":       runtime.GOARCH,
		"ArchBIT":    fmt.Sprintf("%d", strconv.IntSize),
		"FileExt":    defaultFileExtension(),
	}
}

func fetchTool(tmp string, remoteToolName string, trueToolName string, url string, targetDir string) error {
	dir := fmt.Sprintf("%s/%s", tmp, remoteToolName)
	paths := strings.Split(url, "/")
	fileName := paths[len(paths)-1]
	path := fmt.Sprintf("%s/%s", dir, fileName)
	log.Printf("📥 Downloading %s", url)
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
	defer quietly.Close(from)

	to, err := os.OpenFile(filepath.Join(targetDir, binaryName(targetName)), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer quietly.Close(to)
	log.Printf("Copy %s to %s", from.Name(), to.Name())
	_, err = to.ReadFrom(from)
	return err
}

func readToolbox(cfgFile string) (*types.Toolbox, error) {
	var tbFile string
	if cfgFile != "" {
		if _, err := os.Stat(cfgFile); err == nil {
			tbFile = cfgFile
		}
	}

	if tbFile == "" {
		tbFile = filepath.Join(".", toolboxConfFile)
		if _, err := os.Stat(tbFile); errors.Is(err, os.ErrNotExist) {
			userHomeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, err
			}
			homePath := filepath.Join(userHomeDir, toolboxConfFile)
			if _, err := os.Stat(homePath); err == nil {
				tbFile = homePath
			}
		}
	}
	log.Printf("📒 Reading config %s\n\n", tbFile)
	b, err := os.ReadFile(tbFile)
	if err != nil {
		return nil, err
	}
	tb := &types.Toolbox{}
	err = yaml.Unmarshal(b, tb)
	if err != nil {
		return nil, err
	}

	return tb, nil
}

func readVersions(target string) (map[string]string, error) {
	ver := make(map[string]string)
	path := filepath.Join(target, toolboxVersionsFile)
	if _, err := os.Stat(path); err == nil {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		v := &types.Versions{}
		err = yaml.Unmarshal(b, v)
		if err != nil {
			return nil, err
		}
		ver = v.Versions
	}
	return ver, nil
}

func matches(info string, name string) bool {
	ln := strings.ToLower(name)
	if strings.Contains(ln, strings.ToLower(info)) {
		return true
	}

	for _, a := range aliases[info] {
		if strings.Contains(ln, a) {
			matches := true
			for _, sa := range stopAliases[info] {
				if matches && strings.Contains(ln, sa) {
					matches = false
				}
			}
			if matches {
				return matches
			}
		}
	}

	return false
}

func downloadFile(path string, url string) (err error) {
	resp, err := grab.Get(path, url)
	if err != nil {
		return http.CheckError(err)
	}

	log.Printf("Download saved to %s", resp.Filename)
	return nil
}
