package makefile

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/bakito/toolbox/pkg/github"
	"github.com/go-resty/resty/v2"
)

var pattern = regexp.MustCompile(`^github\.com\/([\w-]+\/[\w-]+).*$`)

func Generate(client *resty.Client, makefile string, tools ...string) error {
	var toolData []toolData

	for _, t := range tools {
		td, err := data(client, t)
		if err != nil {
			return err
		}
		toolData = append(toolData, td)
	}

	out := &bytes.Buffer{}
	t := template.Must(template.New("Makefile").Parse(makefileTemplate))
	if err := t.Execute(out, map[string]interface{}{
		"Tools": toolData,
	}); err != nil {
		return err
	}

	if makefile == "" {
		print(out.String())
		return nil
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
	match := pattern.FindStringSubmatch(tool)

	t := toolData{}

	if len(match) != 2 {
		return t, fmt.Errorf("invalid tool %q", tool)
	}

	ghr, err := github.LatestRelease(client, match[1], true)
	if err != nil {
		return t, err
	}

	parts := strings.Split(tool, "/")

	t.Version = ghr.TagName
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
}

const (
	markerStart      = "## toolbox - start"
	markerEnd        = "## toolbox - end"
	makefileTemplate = markerStart + `
## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
{{- range .Tools }}
{{.UpperName}} ?= $(LOCALBIN)/{{.Name}}
{{- end }}

## Tool Versions
{{- range .Tools }}
{{.UpperName}}_VERSION ?= {{.Version}}
{{- end }}

## Tool Installer
{{- range .Tools }}
.PHONY: {{.Name}}
{{.Name}}: $({{.UpperName}}) ## Download {{.Name}} locally if necessary.
$({{.UpperName}}): $(LOCALBIN)
	test -s $(LOCALBIN)/{{.Name}} || GOBIN=$(LOCALBIN) go install {{.Tool}}@$({{.UpperName}}_VERSION)
{{- end }}

## Update Tools
.PHONY: update-toolbox-tools
update-toolbox-tools:
	toolbox makefile -f $$(pwd)/Makefile{{- range .Tools }} \
		{{.Tool}}
{{- end }}
` + markerEnd
)
