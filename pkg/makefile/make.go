package makefile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/bakito/toolbox/pkg/github"
	"github.com/go-resty/resty/v2"
)

var (
	pattern    = regexp.MustCompile(`^github\.com\/([\w-]+\/[\w-]+).*$`)
	getRelease = github.LatestRelease
)

func Generate(client *resty.Client, writer io.Writer, makefile string, tools ...string) error {
	var toolData []toolData

	for _, t := range tools {
		td, err := data(client, t)
		if err != nil {
			return err
		}
		toolData = append(toolData, td)
	}

	sort.Slice(toolData, func(i, j int) bool {
		return toolData[i].Name < toolData[j].Name
	})

	out := &bytes.Buffer{}
	t := template.Must(template.New("Makefile").Parse(makefileTemplate))
	if err := t.Execute(out, map[string]interface{}{
		"Tools": toolData,
	}); err != nil {
		return err
	}

	if makefile == "" {
		_, err := writer.Write(out.Bytes())
		return err
	}

	data, err := os.ReadFile(makefile)
	if err != nil {
		return err
	}

	parts := strings.Split(string(data), markerStart)

	start := parts[0]
	end := ""
	if len(parts) > 1 {
		parts = strings.Split(parts[1], markerEnd)
		if len(parts) > 1 {
			end = parts[1]
		}
	}
	file := start
	file += out.String()
	file += end

	return os.WriteFile(makefile, []byte(file), 0o600)
}

func data(client *resty.Client, tool string) (toolData, error) {
	toolRepo := strings.Split(tool, "@")

	toolName := toolRepo[0]
	repo := toolRepo[len(toolRepo)-1]
	match := pattern.FindStringSubmatch(repo)
	t := toolData{}

	if len(match) != 2 {
		return t, fmt.Errorf("invalid tool %q", tool)
	}

	ghr, err := getRelease(client, match[1], true)
	if err != nil {
		return t, err
	}

	parts := strings.Split(toolName, "/")

	t.Version = ghr.TagName
	t.ToolName = toolName
	t.Tool = tool
	t.Name = parts[len(parts)-1]
	t.UpperName = strings.ReplaceAll(strings.ToUpper(t.Name), "-", "_")
	return t, nil
}

type toolData struct {
	Name      string `json:"Name"`
	UpperName string `json:"UpperName"`
	Version   string `json:"Version"`
	Tool      string `json:"Tool"`
	ToolName  string `json:"ToolName"`
}
