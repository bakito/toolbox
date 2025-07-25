## toolbox - start
## Generated with https://github.com/bakito/toolbox

## Current working directory
TB_LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
TB_LOCALBIN ?= $(TB_LOCALDIR)/bin
$(TB_LOCALBIN):
	if [ ! -e $(TB_LOCALBIN) ]; then mkdir -p $(TB_LOCALBIN); fi

## Tool Binaries
TB_DEEPCOPY_GEN ?= $(TB_LOCALBIN)/deepcopy-gen
TB_GINKGO ?= $(TB_LOCALBIN)/ginkgo
TB_GOLANGCI_LINT ?= $(TB_LOCALBIN)/golangci-lint
TB_GORELEASER ?= $(TB_LOCALBIN)/goreleaser
TB_OAPI_CODEGEN ?= $(TB_LOCALBIN)/oapi-codegen
TB_SEMVER ?= $(TB_LOCALBIN)/semver

## Tool Versions
# renovate: packageName=github.com/kubernetes/code-generator
TB_DEEPCOPY_GEN_VERSION ?= v0.33.3
# renovate: packageName=github.com/golangci/golangci-lint/v2
TB_GOLANGCI_LINT_VERSION ?= v2.3.0
# renovate: packageName=github.com/goreleaser/goreleaser/v2
TB_GORELEASER_VERSION ?= v2.11.0
# renovate: packageName=github.com/deepmap/oapi-codegen/v2
TB_OAPI_CODEGEN_VERSION ?= v2.5.0
# renovate: packageName=github.com/bakito/semver
TB_SEMVER_VERSION ?= v1.1.3

## Tool Installer
.PHONY: tb.deepcopy-gen
tb.deepcopy-gen: $(TB_DEEPCOPY_GEN) ## Download deepcopy-gen locally if necessary.
$(TB_DEEPCOPY_GEN): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/deepcopy-gen || GOBIN=$(TB_LOCALBIN) go install k8s.io/code-generator/cmd/deepcopy-gen@$(TB_DEEPCOPY_GEN_VERSION)
.PHONY: tb.ginkgo
tb.ginkgo: $(TB_GINKGO) ## Download ginkgo locally if necessary.
$(TB_GINKGO): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/ginkgo || GOBIN=$(TB_LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo
.PHONY: tb.golangci-lint
tb.golangci-lint: $(TB_GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(TB_GOLANGCI_LINT): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/golangci-lint || GOBIN=$(TB_LOCALBIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(TB_GOLANGCI_LINT_VERSION)
.PHONY: tb.goreleaser
tb.goreleaser: $(TB_GORELEASER) ## Download goreleaser locally if necessary.
$(TB_GORELEASER): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/goreleaser || GOBIN=$(TB_LOCALBIN) go install github.com/goreleaser/goreleaser/v2@$(TB_GORELEASER_VERSION)
.PHONY: tb.oapi-codegen
tb.oapi-codegen: $(TB_OAPI_CODEGEN) ## Download oapi-codegen locally if necessary.
$(TB_OAPI_CODEGEN): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/oapi-codegen || GOBIN=$(TB_LOCALBIN) go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@$(TB_OAPI_CODEGEN_VERSION)
.PHONY: tb.semver
tb.semver: $(TB_SEMVER) ## Download semver locally if necessary.
$(TB_SEMVER): $(TB_LOCALBIN)
	test -s $(TB_LOCALBIN)/semver || GOBIN=$(TB_LOCALBIN) go install github.com/bakito/semver@$(TB_SEMVER_VERSION)

## Reset Tools
.PHONY: tb.reset
tb.reset:
	@rm -f \
		$(TB_LOCALBIN)/deepcopy-gen \
		$(TB_LOCALBIN)/ginkgo \
		$(TB_LOCALBIN)/golangci-lint \
		$(TB_LOCALBIN)/goreleaser \
		$(TB_LOCALBIN)/oapi-codegen \
		$(TB_LOCALBIN)/semver

## Update Tools
.PHONY: tb.update
tb.update: tb.reset
	toolbox makefile --renovate -f $(TB_LOCALDIR)/Makefile \
		k8s.io/code-generator/cmd/deepcopy-gen@github.com/kubernetes/code-generator \
		github.com/golangci/golangci-lint/v2/cmd/golangci-lint \
		github.com/goreleaser/goreleaser/v2 \
		github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen \
		github.com/bakito/semver
## toolbox - end