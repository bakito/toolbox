# Run go golanci-lint
lint: golangci-lint
	$(GOLANGCI_LINT) run --fix

# Run go mod tidy
tidy:
	go mod tidy

# Run tests
test: ginkgo tidy lint
	$(GINKGO) -r --cover --coverprofile=coverage.out

release: goreleaser semver
	@version=$$($(SEMVER)); \
	git tag -s $$version -m"Release $$version"
	$(GORELEASER) --clean

test-release: goreleaser
	$(GORELEASER) --skip=publish --snapshot --clean

## toolbox - start
## Current working directory
LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
LOCALBIN ?= $(LOCALDIR)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
DEEPCOPY_GEN ?= $(LOCALBIN)/deepcopy-gen
GINKGO ?= $(LOCALBIN)/ginkgo
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
GORELEASER ?= $(LOCALBIN)/goreleaser
SEMVER ?= $(LOCALBIN)/semver

## Tool Versions
# renovate: packageName=k8s.io/code-generator/cmd/deepcopy-gen
DEEPCOPY_GEN_VERSION ?= v0.31.1
# renovate: packageName=github.com/golangci/golangci-lint/cmd/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.61.0
# renovate: packageName=github.com/goreleaser/goreleaser/v2
GORELEASER_VERSION ?= v2.3.2
# renovate: packageName=github.com/bakito/semver
SEMVER_VERSION ?= v1.1.3

## Tool Installer
.PHONY: deepcopy-gen
deepcopy-gen: $(DEEPCOPY_GEN) ## Download deepcopy-gen locally if necessary.
$(DEEPCOPY_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/deepcopy-gen || GOBIN=$(LOCALBIN) go install k8s.io/code-generator/cmd/deepcopy-gen@$(DEEPCOPY_GEN_VERSION)
.PHONY: ginkgo
ginkgo: $(GINKGO) ## Download ginkgo locally if necessary.
$(GINKGO): $(LOCALBIN)
	test -s $(LOCALBIN)/ginkgo || GOBIN=$(LOCALBIN) go install github.com/onsi/ginkgo/v2/ginkgo
.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	test -s $(LOCALBIN)/golangci-lint || GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
.PHONY: goreleaser
goreleaser: $(GORELEASER) ## Download goreleaser locally if necessary.
$(GORELEASER): $(LOCALBIN)
	test -s $(LOCALBIN)/goreleaser || GOBIN=$(LOCALBIN) go install github.com/goreleaser/goreleaser/v2@$(GORELEASER_VERSION)
.PHONY: semver
semver: $(SEMVER) ## Download semver locally if necessary.
$(SEMVER): $(LOCALBIN)
	test -s $(LOCALBIN)/semver || GOBIN=$(LOCALBIN) go install github.com/bakito/semver@$(SEMVER_VERSION)

## Update Tools
.PHONY: update-toolbox-tools
update-toolbox-tools:
	@rm -f \
		$(LOCALBIN)/deepcopy-gen \
		$(LOCALBIN)/ginkgo \
		$(LOCALBIN)/golangci-lint \
		$(LOCALBIN)/goreleaser \
		$(LOCALBIN)/semver
	toolbox makefile --renovate -f $(LOCALDIR)/Makefile \
		k8s.io/code-generator/cmd/deepcopy-gen@github.com/kubernetes/code-generator \
		github.com/golangci/golangci-lint/cmd/golangci-lint \
		github.com/goreleaser/goreleaser/v2 \
		github.com/bakito/semver
## toolbox - end
