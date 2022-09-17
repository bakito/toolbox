package main

import (
	"bytes"
	"fmt"
	"github.com/bakito/toolbox/pkg/extract"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/cavaliergopher/grab/v3"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {

	tb, err := readToolbox()
	if err != nil {
		panic(err)
	}

	tmp, err := os.MkdirTemp("", "toolbox")
	if err != nil {
		panic(err)
	}

	defer os.RemoveAll(tmp)

	client := resty.New()

	for _, tool := range tb.Tools {
		log.Printf("Download %s\n", tool.Name)
		if tool.Github != "" {
			ghr := &types.GithubRelease{}
			_, err := client.R().
				EnableTrace().
				SetResult(ghr).
				Get(tool.LatestURL())
			//Get(os.Args[0])

			if err != nil {
				panic(err)
			}

			if tool.Version == "" {
				tool.Version = ghr.Name
				log.Printf("Latest Version: %s", tool.Version)
			}

			for _, a := range ghr.Assets {
				if strings.Contains(a.Name, tool.Name) && matches(runtime.GOOS, a.Name) && matches(runtime.GOARCH, a.Name) {
					if err := fetchTool(tmp, tool.Name, a.Name, a.BrowserDownloadUrl); err != nil {
						panic(err)
					}
				}
				for _, add := range tool.Additional {

					if strings.Contains(a.Name, add) && matches(runtime.GOOS, a.Name) && matches(runtime.GOARCH, a.Name) {
						if err := fetchTool(tmp, add, a.Name, a.BrowserDownloadUrl); err != nil {
							panic(err)
						}
					}
				}
			}
		} else if tool.Version != "" {
			if strings.HasPrefix(tool.Version, "http") {
				resp, err := client.R().
					EnableTrace().
					Get(tool.Version)
				if err != nil {
					panic(err)
				}
				tool.Version = string(resp.Body())
			}
			ut, err := template.New("url").Parse(tool.DownloadURL)
			if err != nil {
				panic(err)
			}

			var b bytes.Buffer
			err = ut.Execute(&b, map[string]string{"Version": tool.Version, "OS": runtime.GOOS, "Arch": runtime.GOARCH})

			fetchTool(tmp, tool.Name, tool.Name, b.String())
		}
	}
}

func fetchTool(tmp string, toolName string, fileName string, url string) error {
	dir := fmt.Sprintf("%s/%s", tmp, toolName)
	path := fmt.Sprintf("%s/%s", dir, fileName)
	log.Printf("Downloading %s", url)
	if err := downloadFile(path, url); err != nil {
		return err
	}
	if err := extract.File(path, dir); err != nil {
		return err
	}
	if _, err := copyTool(dir, toolName); err != nil {
		return err
	}
	return nil
}

func copyTool(dir string, tool string) (bool, error) {
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

			if err := copyFile(dir, file); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	for _, d := range dirs {
		ok, err := copyTool(filepath.Join(dir, d.Name()), tool)
		if ok || err != nil {
			return ok, err
		}
	}
	return false, nil
}

func copyFile(dir string, file os.DirEntry) error {
	from, err := os.Open(filepath.Join(dir, file.Name()))
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.OpenFile(filepath.Join("download", file.Name()), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer to.Close()

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

	fmt.Println("Download saved to", resp.Filename)
	return nil
}
