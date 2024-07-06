package fetcher

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
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
	"gopkg.in/yaml.v3"
)

const (
	toolboxConfFile     = ".toolbox.yaml"
	toolboxVersionsFile = ".toolbox-versions.yaml"
	oldExecutablePrefix = ".toolbox-old."
)

var (
	aliases = map[string][]string{
		"amd64":   {"x86_64", "64", "64bit"},
		"windows": {"win", "win64"},
		"linux":   {"linux64"},
	}
	stopAliases = map[string][]string{
		"amd64":   {"arm", "mips", "ppc", "risc", "s390"},
		"windows": {"darwin"},
	}

	excludedSuffixes = []string{"sum", "sha256", "sbom", "pem", "sig", "rpm", "txt", "deb", "json", "asc"}
)

func New() Fetcher {
	return &fetcher{}
}

type Fetcher interface {
	Fetch(string, ...string) error
}
type fetcher struct {
	executablePath string
	upx            bool
}

func (f *fetcher) Fetch(cfgFile string, selectedTools ...string) error {
	var err error
	f.executablePath, err = os.Executable()
	if err != nil {
		return err
	}

	log.Printf("🧰 toolbox %s", version.Version)

	client := resty.New()

	tbRel, err := github.LatestRelease(client, "bakito/toolbox", true)
	if err != nil {
		return err
	}
	if tbRel.TagName != version.Version {
		log.Printf("🌟 A new toolbox version is available %s (current: %s)\n", tbRel.TagName, version.Version)
	}

	tb, _, err := ReadToolbox(cfgFile)
	if err != nil {
		return err
	}
	if tb.Target == "" {
		tb.Target = "./tools"
	}

	if tb.Upx {
		cmd := exec.Command("upx", "--version")
		_, err := cmd.Output()
		if err == nil {
			log.Printf("🗜️ upx is available")
			f.upx = true
		}
	}

	err = filepath.Walk(tb.Target, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			if strings.HasPrefix(f.Name(), oldExecutablePrefix) {
				toolPath := filepath.Join(tb.Target, f.Name())
				if err := os.Remove(toolPath); err != nil {
					return err
				}
				log.Printf("🗑️  Delete old tool %s\n", toolPath)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	println()

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
		if contains(selectedTools, tool.Name) {
			if err := f.handleTool(client, ver, tmp, tb, tool); err != nil {
				return err
			}
		}
	}

	// save versions
	return SaveYamlFile(filepath.Join(tb.Target, toolboxVersionsFile), tb.Versions())
}

func SaveYamlFile(path string, obj interface{}) error {
	var b bytes.Buffer
	env := yaml.NewEncoder(&b)
	env.SetIndent(2)

	if err := env.Encode(obj); err != nil {
		return err
	}
	return os.WriteFile(path, b.Bytes(), 0o600)
}

func (f *fetcher) handleTool(client *resty.Client, ver map[string]string, tmp string, tb *types.Toolbox, tool *types.Tool) error {
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
		return f.downloadFromURL(client, tb, ver, tmp, tool)
	} else if ghr != nil {
		return f.downloadViaGithub(tb, tool, ghr, tmp)
	}
	return nil
}

func (f *fetcher) downloadViaGithub(tb *types.Toolbox, tool *types.Tool, ghr *types.GithubRelease, tmp string) error {
	matching := findMatching(tb, tool.Name, ghr.Assets)
	tool.CouldNotBeFound = true
	if matching != nil {
		tool.CouldNotBeFound = false
		if err := f.fetchTool(tmp, tool.Name, tool.Name, matching.BrowserDownloadURL, tb.Target, tool.Check); err != nil {
			return err
		}
	}
	for _, add := range tool.Additional {
		matching := findMatching(nil, add, ghr.Assets)
		if matching != nil {
			tool.CouldNotBeFound = false
			if err := f.fetchTool(tmp, add, add, matching.BrowserDownloadURL, tb.Target, tool.Check); err != nil {
				return err
			}
		}
	}
	if tool.CouldNotBeFound {
		log.Printf("❌ Couldn't find a file here!\n")
	}
	return nil
}

func (f *fetcher) downloadFromURL(client *resty.Client, tb *types.Toolbox, ver map[string]string, tmp string, tool *types.Tool) error {
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
	return f.fetchTool(tmp, tool.Name, tool.Name, parseTemplate(tool.DownloadURL, tool.Version), tb.Target, tool.Check)
}

func findMatching(tb *types.Toolbox, toolName string, assets []types.Asset) *types.Asset {
	var matching []*types.Asset
	for i := range assets {
		a := assets[i]
		if strings.Contains(a.Name, toolName) &&
			matches(runtime.GOOS, a.Name) &&
			!hasForbiddenSuffix(tb, a) {
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

func hasForbiddenSuffix(tb *types.Toolbox, a types.Asset) bool {
	excl := excludedSuffixes
	if tb != nil && len(tb.ExcludedSuffixes) != 0 {
		excl = tb.ExcludedSuffixes
	}
	for _, suffix := range excl {
		if strings.HasSuffix(a.Name, suffix) {
			return true
		}
	}
	return false
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

func (f *fetcher) fetchTool(tmp string, remoteToolName string, trueToolName string, url string, targetDir string, check string) error {
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

	targetFilePath, err := filepath.Abs(filepath.Join(targetDir, trueToolName))
	if err != nil {
		return err
	}
	if f.executablePath == targetFilePath {
		renameTo := filepath.Join(targetDir, oldExecutablePrefix+trueToolName)
		err = os.Rename(targetFilePath, renameTo)
		if err != nil {
			return err
		}
		log.Printf("🔀 Rename current executable to %s", renameTo)
	}
	ok, err := f.copyTool(dir, remoteToolName, targetDir, trueToolName, check)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("could not find: %s", remoteToolName)
	}
	return nil
}

func (f *fetcher) copyTool(dir string, fileName string, targetDir string, targetName string, check string) (bool, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	var dirs []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
		if file.Name() == binaryName(fileName) ||
			file.Name() == fileName ||
			file.Name() == binaryName(fmt.Sprintf("%s_%s_%s", fileName, runtime.GOOS, runtime.GOARCH)) ||
			file.Name() == binaryName(fmt.Sprintf("%s-%s_%s", fileName, runtime.GOOS, runtime.GOARCH)) {

			sourcePath := filepath.Join(dir, file.Name())
			targetPath := filepath.Join(targetDir, binaryName(targetName))

			if err := copyFile(sourcePath, targetPath); err != nil {
				return false, err
			}
			if f.upx {
				log.Printf("🗜️ Compressing with upx")
				cmd := exec.Command("upx", "-q", "-q", targetPath)
				stdout, err := cmd.Output()
				if err == nil {
					parts := strings.Fields(string(stdout))
					size, _ := strconv.Atoi(parts[2])
					log.Printf("\tCompressed to %s (%s)", parts[3], formatBytes(int64(size)))
				} else {
					var ee *exec.ExitError
					if errors.As(err, &ee) && ee.ExitCode() == 2 {
						log.Printf("\tAlready Compressed")
					} else {
						log.Printf("\tCompression error: %v", err)
					}
				}
			}

			if len(check) > 0 {
				cmd := exec.Command(targetPath, strings.Fields(check)...)
				if _, err := cmd.Output(); err != nil {
					log.Printf("🚫 Check failed: %v", err)
				} else {
					log.Printf("👍 Check was successful")
				}

			}

			return true, nil
		}
	}
	for _, d := range dirs {
		ok, err := f.copyTool(filepath.Join(dir, d.Name()), fileName, targetDir, targetName, check)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func copyFile(sourcePath string, targetPath string) error {
	from, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	fromStat, err := from.Stat()
	if err != nil {
		return err
	}
	defer quietly.Close(from)

	to, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer quietly.Close(to)
	log.Printf("Copy %s to %s (%v)", from.Name(), to.Name(), formatBytes(fromStat.Size()))
	_, err = to.ReadFrom(from)
	return err
}

func ReadToolbox(cfgFile string) (*types.Toolbox, string, error) {
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
				return nil, "", err
			}
			homePath := filepath.Join(userHomeDir, toolboxConfFile)
			if _, err := os.Stat(homePath); err == nil {
				tbFile = homePath
			}
		}
	}
	log.Printf("📒 Reading config %s\n", tbFile)
	b, err := os.ReadFile(tbFile)
	if err != nil {
		return nil, "", err
	}
	tb := &types.Toolbox{}
	err = yaml.Unmarshal(b, tb)
	if err != nil {
		return nil, "", err
	}

	return tb, tbFile, nil
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
	req, err := grab.NewRequest(path, url)
	if err != nil {
		return err
	}
	client := grab.NewClient()
	req.HTTPRequest.Header.Set("User-Agent", fmt.Sprintf("toolbox/%s", version.Version))

	resp := client.Do(req)
	if resp.Err() != nil {
		return http.CheckError(err)
	}

	log.Printf("Download saved to %s", resp.Filename)
	return nil
}

func contains(list []string, v string) bool {
	for _, s := range list {
		if s == v {
			return true
		}
	}
	return len(list) == 0
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}
