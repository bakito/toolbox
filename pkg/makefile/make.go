// Package makefile
package makefile

import (
	"bytes"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/go-resty/resty/v2"

	"github.com/bakito/toolbox/pkg/github"
)

var (
	githubPattern = regexp.MustCompile(`^github\.com/([\w-]+/[\w-]+).*$`)
	getRelease    = github.LatestRelease
)

func Generate(client *resty.Client, makefile string, renovate, toolchain bool, toolsFile string, tools ...string) error {
	argTools, toolData := mergeWithToolsGo(toolsFile, unique(tools))
	return generateForTools(client, makefile, renovate, toolchain, argTools, toolData)
}

func generateForTools(
	client *resty.Client,
	makefile string,
	renovate, toolchain bool,
	argTools []string,
	toolData []toolData,
) error {
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
	withVersionArgs := false
	for _, td := range toolData {
		if !withVersions && td.Version != "" {
			withVersions = true
		}
		if !withVersionArgs && td.VersionArg != "" {
			withVersionArgs = true
		}
	}

	out := &bytes.Buffer{}
	t := template.Must(template.New("toolbox.mk").Parse(makefileTemplate))
	if err := t.Execute(out, map[string]any{
		"Tools":           toolData,
		"WithVersions":    withVersions,
		"WithVersionArgs": withVersionArgs,
		"Renovate":        renovate,
		"Toolchain":       toolchain,
	}); err != nil {
		return err
	}

	makefile, err := filepath.Abs(makefile)
	if err != nil {
		return err
	}

	includeFile := filepath.Join(filepath.Dir(makefile), includeFileName)

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

	var file string
	if !strings.Contains(string(data), "include ./"+includeFileName) {
		file = fmt.Sprintf("# Include toolbox tasks\ninclude ./%s\n\n", includeFileName)
	}

	file += start
	file += strings.TrimSpace(end)

	if renovate {
		if err := updateRenovateConf(); err != nil {
			return err
		}
	}
	if err := os.WriteFile(includeFile, out.Bytes(), 0o600); err != nil {
		return err
	}

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

func dataForTool(fromToolsGo bool, toolName string, fullTool ...string) toolData {
	var td toolData
	td.ToolName = toolName

	if sp := strings.Split(td.ToolName, "?"); len(sp) > 1 {
		td.ToolName = sp[0]
		td.VersionArg = sp[1]
	}

	parts := strings.Split(td.ToolName, "/")

	if len(fullTool) == 1 {
		sp := strings.Split(fullTool[0], "?")
		td.Tool = sp[0]
		if len(sp) > 1 {
			td.VersionArg = sp[1]
		}
	} else {
		td.Tool = toolName
	}
	if match, _ := regexp.MatchString(`^v\d+$`, parts[len(parts)-1]); match {
		td.Name = parts[len(parts)-2]
	} else {
		td.Name = parts[len(parts)-1]
	}

	td.UpperName = strings.ReplaceAll(strings.ToUpper(td.Name), "-", "_")
	td.FromToolsGo = fromToolsGo
	td.GoModule = extractModulePath(td.ToolName)
	td.RepoURL = td.GoModule
	if sp := strings.Split(td.Tool, "@"); len(sp) > 1 {
		td.RepoURL = sp[1]
	}
	return td
}

func extractModulePath(importPath string) string {
	// Remove quotes if present
	importPath = strings.Trim(importPath, `"`)

	re := regexp.MustCompile(`^(.*)(/v\d+)(/.*)?$`)

	matches := re.FindStringSubmatch(importPath)
	if matches == nil {
		return importPath
	}

	base := matches[1]
	version := matches[2]

	return base + version
}

type toolData struct {
	Name        string `json:"Name"`
	UpperName   string `json:"UpperName"`
	Version     string `json:"Version"`
	GoModule    string `json:"GoModule"`
	RepoURL     string `json:"RepoURL"`
	Tool        string `json:"Tool"`
	ToolName    string `json:"ToolName"`
	FromToolsGo bool   `json:"FromToolsGo"`
	VersionArg  string `json:"VersionArg"`
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

func unique(slice []string) []string {
	// create a map with all the values as key
	uniqMap := make(map[string]struct{})
	for _, v := range slice {
		uniqMap[v] = struct{}{}
	}

	return slices.Sorted(maps.Keys(uniqMap))
}
