## Generated with https://github.com/bakito/toolbox

## Current working directory
TB_LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
TB_LOCALBIN ?= $(TB_LOCALDIR)/bin
$(TB_LOCALBIN):
	mkdir -p $(TB_LOCALBIN)

## Tool Binaries
{{- range .Tools }}
TB_{{.UpperName}} ?= $(TB_LOCALBIN)/{{.Name}}
{{- end }}
{{- if .WithVersions }}

## Tool Versions
{{- range .Tools }}
{{- if .Version }}
{{- if $.Renovate }}
# renovate: packageName={{.ToolName}}
{{- end }}
TB_{{.UpperName}}_VERSION ?= {{.Version}}
{{- end }}
{{- end }}
{{- end }}

## Tool Installer
{{- range .Tools }}
.PHONY: tb.{{.Name}}
tb.{{.Name}}: $(TB_{{.UpperName}}) ## Download {{.Name}} locally if necessary.
$(TB_{{.UpperName}}): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/{{.Name}} || GOBIN=$(TB_LOCALBIN) go install {{.ToolName}}{{- if .Version }}@$(TB_{{.UpperName}}_VERSION){{- end }}
{{- end }}

## Update Tools
.PHONY: tb.update
tb.update:
	@rm -f{{- range .Tools }} \
		$(TB_LOCALBIN)/{{.Name}}
{{- end }}
	toolbox makefile {{ if $.Renovate }}--renovate {{ end }}-f $(TB_LOCALDIR)/Makefile{{- range .Tools }}{{- if not .FromToolsGo }} \
		{{.Tool}}{{- end }}
{{- end }}
