## Current working directory
LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
LOCALBIN ?= $(LOCALDIR)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
{{- range .Tools }}
{{.UpperName}} ?= $(LOCALBIN)/{{.Name}}
{{- end }}

## Tool Versions
{{- range .Tools }}
{{- if .Version }}
{{.UpperName}}_VERSION ?= {{.Version}}
{{- end }}
{{- end }}

## Tool Installer
{{- range .Tools }}
.PHONY: {{.Name}}
{{.Name}}: $({{.UpperName}}) ## Download {{.Name}} locally if necessary.
$({{.UpperName}}): $(LOCALBIN)
	test -s $(LOCALBIN)/{{.Name}} || GOBIN=$(LOCALBIN) go install {{.ToolName}}{{- if .Version }}@$({{.UpperName}}_VERSION){{- end }}
{{- end }}

## Update Tools
.PHONY: update-toolbox-tools
update-toolbox-tools:
	@rm -f{{- range .Tools }} \
		$(LOCALBIN)/{{.Name}}
{{- end }}
	toolbox makefile -f $(LOCALDIR)/Makefile{{- range .Tools }}{{- if not .WithDependency }} \
		{{.Tool}}{{- end }}
{{- end }}
