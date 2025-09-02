## toolbox - start
## Generated with https://github.com/bakito/toolbox

## Current working directory
TB_LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
TB_LOCALBIN ?= $(TB_LOCALDIR)/bin
$(TB_LOCALBIN):
	if [ ! -e $(TB_LOCALBIN) ]; then mkdir -p $(TB_LOCALBIN); fi

## Tool Binaries
TB_CONTROLLER_GEN ?= $(TB_LOCALBIN)/controller-gen
TB_SEMVER ?= $(TB_LOCALBIN)/semver
TB_TOOLBOX ?= $(TB_LOCALBIN)/toolbox

## Tool Installer
.PHONY: tb.controller-gen
tb.controller-gen: ## Download controller-gen locally if necessary.
	@test -s $(TB_CONTROLLER_GEN) || \
		GOBIN=$(TB_LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen
.PHONY: tb.semver
tb.semver: ## Download semver locally if necessary.
	@test -s $(TB_SEMVER) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/bakito/semver
.PHONY: tb.toolbox
tb.toolbox: ## Download toolbox locally if necessary.
	@test -s $(TB_TOOLBOX) || \
		GOBIN=$(TB_LOCALBIN) go install github.com/bakito/toolbox

## Reset Tools
.PHONY: tb.reset
tb.reset:
	@rm -f \
		$(TB_CONTROLLER_GEN) \
		$(TB_SEMVER) \
		$(TB_TOOLBOX)

## Update Tools
.PHONY: tb.update
tb.update: tb.reset
	toolbox makefile -f $(TB_LOCALDIR)/Makefile
## toolbox - end
