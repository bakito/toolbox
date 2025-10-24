# Include toolbox tasks
include ./.toolbox.mk

# Run go golanci-lint
lint: tb.golangci-lint
	$(TB_GOLANGCI_LINT) run --fix

# Run go mod tidy
tidy:
	go mod tidy

# Run tests
test: tb.ginkgo
	$(TB_GINKGO) -r --cover --coverprofile=coverage.out

release: tb.goreleaser tb.semver tb.syft
	@version=$$($(TB_SEMVER)); \
	git tag -s $$version -m"Release $$version"
	PATH=$(TB_LOCALBIN):$${PATH} $(TB_GORELEASER) --clean

test-release: tb.goreleaser tb.syft
	PATH=$(TB_LOCALBIN):$${PATH} $(TB_GORELEASER) --skip=publish --snapshot --clean
