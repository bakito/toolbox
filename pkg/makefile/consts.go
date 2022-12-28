package makefile

const (
	markerStart      = "## toolbox - start"
	markerEnd        = "## toolbox - end"
	makefileTemplate = markerStart + `
## Location to install dependencies to
LOCALBIN ?= $(shell test -s "cygpath -m $$(pwd)" || pwd)/bin
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
	test -s $(LOCALBIN)/{{.Name}} || GOBIN=$(LOCALBIN) go install {{.ToolName}}@$({{.UpperName}}_VERSION)
{{- end }}

## Update Tools
.PHONY: update-toolbox-tools
update-toolbox-tools:
	@rm -f{{- range .Tools }} \
		$(LOCALBIN)/{{.Name}}
{{- end }}
	toolbox makefile -f $$(pwd)/Makefile{{- range .Tools }} \
		{{.Tool}}
{{- end }}
` + markerEnd
)