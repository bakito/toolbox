# Run go golanci-lint
lint: golangci-lint
	$(LOCALBIN)/golangci-lint run --fix

# Run go mod tidy
tidy:
	go mod tidy

# Run tests
test: tidy lint
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

release: semver
	@version=$$($(LOCALBIN)/semver); \
	git tag -s $$version -m"Release $$version"
	goreleaser --rm-dist

test-release:
	goreleaser --skip-publish --snapshot --rm-dist

## toolbox - start
## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
SEMVER ?= $(LOCALBIN)/semver
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint

## Tool Versions
SEMVER_VERSION ?= v1.1.3
GOLANGCI_LINT_VERSION ?= v1.50.1

## Tool Installer
.PHONY: semver
semver: $(SEMVER) ## Download semver locally if necessary.
$(SEMVER): $(LOCALBIN)
	test -s $(LOCALBIN)/semver || GOBIN=$(LOCALBIN) go install github.com/bakito/semver@$(SEMVER_VERSION)
.PHONY: golangci-lint
golangci-lint: $(GOLANGCI_LINT) ## Download golangci-lint locally if necessary.
$(GOLANGCI_LINT): $(LOCALBIN)
	test -s $(LOCALBIN)/golangci-lint || GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

## Update Tools
.PHONY: update-toolbox-tools
update-toolbox-tools:
	@rm -f \
		$(LOCALBIN)/semver \
		$(LOCALBIN)/golangci-lint
	toolbox makefile -f $$(pwd)/Makefile \
		github.com/bakito/semver \
		github.com/golangci/golangci-lint/cmd/golangci-lint
## toolbox - end
