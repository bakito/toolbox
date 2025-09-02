## Generated with https://github.com/bakito/toolbox

## Current working directory
TB_LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
TB_LOCALBIN ?= $(TB_LOCALDIR)/bin
$(TB_LOCALBIN):
	if [ ! -e $(TB_LOCALBIN) ]; then mkdir -p $(TB_LOCALBIN); fi
{{- if $.Toolchain }}

## Go Version
TB_GO_VERSION ?= $(shell grep -E '^go [0-9]+\.[0-9]+' go.mod | awk '{print $$2}')
{{- end }}

## Tool Binaries
{{- range .Tools }}
TB_{{.UpperName}} ?= $(TB_LOCALBIN)/{{.Name}}
{{- end }}
{{- if .WithVersions }}

## Tool Versions
{{- range .Tools }}
{{- if .Version }}
{{- if $.Renovate }}
# renovate: packageName={{.RepoURL}}
{{- end }}
TB_{{.UpperName}}_VERSION ?= {{.Version}}
{{- if .VersionParam }}
TB_{{.UpperName}}_VERSION_NUM ?= {{.VersionNumeric}}
{{- end }}
{{- end }}
{{- end }}
{{- end }}

## Tool Installer
{{- range .Tools }}
.PHONY: tb.{{.Name}}
tb.{{.Name}}: ## Download {{.Name}} locally if necessary.
	@test -s $(TB_{{.UpperName}}) {{ if .VersionParam }}&& $(TB_{{.UpperName}}) {{ .VersionParam }} | grep -q $(TB_{{.UpperName}}_VERSION_NUM) {{ end }}|| \
		GOBIN=$(TB_LOCALBIN) {{ if $.Toolchain }}GOTOOLCHAIN=go$(TB_GO_VERSION) {{ end }}go install {{.ToolName}}{{- if .Version }}@$(TB_{{.UpperName}}_VERSION){{- end }}
{{- end }}

## Reset Tools
.PHONY: tb.reset
tb.reset:
	@rm -f{{- range .Tools }} \
		$(TB_{{.UpperName}})
{{- end }}

## Update Tools
.PHONY: tb.update
tb.update: tb.reset
	toolbox makefile {{ if $.Renovate }}--renovate {{ end }}{{ if $.Toolchain }}--toolchain {{ end }}-f $(TB_LOCALDIR)/Makefile{{- range .Tools }}{{- if not .FromToolsGo }} \
		{{.Tool}}{{ if .VersionParam }}?{{ .VersionParam }}{{ end }}{{- end }}
{{- end }}
