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
	githubPattern = regexp.MustCompile(`^github\.com/([\w-]+/[\w-]+).*$`)
	getRelease    = github.LatestRelease
)

func Generate(client *resty.Client, writer io.Writer, makefile string, toolsFile string, tools ...string) error {
	argTools, toolData := mergeWithToolsGo(toolsFile, tools)
	return generate(client, writer, makefile, argTools, toolData)
}

func generate(client *resty.Client, writer io.Writer, makefile string, argTools []string, toolData []toolData) error {
	for _, t := range argTools {
		td, err := dataForArg(client, t)
		if err != nil {
			return err
		}
		toolData = append(toolData, td)
	}

	sort.Slice(toolData, func(i, j int) bool {
		return toolData[i].Name < toolData[j].Name
	})

	withVersions := false
	for _, td := range toolData {
		if !withVersions && td.Version != "" {
			withVersions = true
		}
	}

	out := &bytes.Buffer{}
	t := template.Must(template.New("Makefile").Parse(makefileTemplate))
	if err := t.Execute(out, map[string]interface{}{
		"Tools":        toolData,
		"WithVersions": withVersions,
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

func dataForArg(client *resty.Client, tool string) (toolData, error) {
	toolRepo := strings.Split(tool, "@")
	toolName := toolRepo[0]

	td := dataForTool(false, toolName, tool)

	repo := toolRepo[len(toolRepo)-1]
	match := githubPattern.FindStringSubmatch(repo)

	if len(match) != 2 {
		return td, fmt.Errorf("invalid tool %q", tool)
	}
	ghr, err := getRelease(client, match[1], true)
	if err != nil {
		return td, err
	}
	td.Version = ghr.TagName

	return td, nil
}

func dataForTool(fromToolsGo bool, toolName string, fullTool ...string) (td toolData) {
	parts := strings.Split(toolName, "/")
	td.ToolName = toolName
	if len(fullTool) == 1 {
		td.Tool = fullTool[0]
	} else {
		td.Tool = toolName
	}
	td.Name = parts[len(parts)-1]
	td.UpperName = strings.ReplaceAll(strings.ToUpper(td.Name), "-", "_")
	td.FromToolsGo = fromToolsGo
	return
}

type toolData struct {
	Name        string `json:"Name"`
	UpperName   string `json:"UpperName"`
	Version     string `json:"Version"`
	Tool        string `json:"Tool"`
	ToolName    string `json:"ToolName"`
	FromToolsGo bool   `json:"FromToolsGo"`
}

func mergeWithToolsGo(fileName string, inTools []string) ([]string, []toolData) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return inTools, nil
	}

	t := make(map[string]bool)
	for _, tool := range inTools {
		t[tool] = true
	}

	r := regexp.MustCompile(`"(.*)"`)
	var goTools []toolData
	for _, m := range r.FindAllStringSubmatch(string(content), -1) {
		tool := m[1]
		goTools = append(goTools, dataForTool(true, tool))
		delete(t, tool)
	}

	var argTools []string
	for t := range t {
		argTools = append(argTools, t)
	}

	return argTools, goTools
}
