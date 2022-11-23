package makefile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/bakito/toolbox/pkg/github"
	"github.com/go-resty/resty/v2"
)

var pattern = regexp.MustCompile(`^github\.com\/([\w-]+\/[\w-]+)\/.*$`)

func Generate(tool string, client *resty.Client) error {
	match := pattern.FindStringSubmatch(tool)

	if len(match) != 2 {
		return fmt.Errorf("invalid tool %q", tool)
	}

	ghr, err := github.LatestRelease(client, match[1], true)
	if err != nil {
		return err
	}

	parts := strings.Split(tool, "/")
	toolName := parts[len(parts)-1]

	t := template.Must(template.New("Makefile").Parse(makefileTemplate))
	return t.Execute(os.Stdout, map[string]string{
		"Name":      toolName,
		"UpperName": strings.ToUpper(toolName),
		"Version":   ghr.TagName,
		"Tool":      tool,
	})
}

const makefileTemplate = `
## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
{{.UpperName}} ?= $(LOCALBIN)/{{.Name}}

## Tool Versions
{{.UpperName}}_VERSION ?= {{.Version}}

.PHONY: {{.Name}}
{{.Name}}: $({{.UpperName}}) ## Download {{.Name}} locally if necessary.
$({{.UpperName}}): $(LOCALBIN)
	test -s $(LOCALBIN)/{{.Name}} || GOBIN=$(LOCALBIN) go install {{.Tool}}@$({{.UpperName}}_VERSION)
`
