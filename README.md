[![Go Report Card](https://goreportcard.com/badge/github.com/bakito/toolbox)](https://goreportcard.com/report/github.com/bakito/toolbox)
[![Coverage Status](https://coveralls.io/repos/github/bakito/toolbox/badge.svg?branch=main&service=github)](https://coveralls.io/github/bakito/toolbox?branch=main)

# toolbox

🧰 a small toolbox helping to fetch tools

## Fetch tools

```text
Fetch all tools

Usage:
  toolbox fetch [flags]

Flags:
  -c, --config string   The config file to be used. (default 1. '.toolbox.yaml' current dir, 2. '~/.toolbox.yaml')
  -h, --help            help for fetch

```

### .toolbox.yaml

```yaml
tools:
  kubexporter:
    github: bakito/kubexporter
  upx:
    github: upx/upx
  kubectx:
    github: ahmetb/kubectx
    additional:
      - kubens
  kubectl:
    downloadURL: https://dl.k8s.io/release/{{ .Version }}/bin/{{ .OS }}/{{ .Arch }}/kubectl{{ .FileExt }}
    version: https://storage.googleapis.com/kubernetes-release/release/stable.txt
  helm:
    github: helm/helm
    downloadURL: https://get.helm.sh/helm-{{ .Version }}-{{ .OS }}-{{ .Arch }}.tar.gz
  jq:
    github: stedolan/jq
  yq:
    github: mikefarah/yq
  vault:
    github: hashicorp/vault
    downloadURL: https://releases.hashicorp.com/vault/{{ .VersionNum }}/vault_{{ .VersionNum }}_{{ .OS }}_{{ .Arch }}.zip
  terraform:
    github: hashicorp/terraform
    downloadURL: https://releases.hashicorp.com/terraform/{{ .VersionNum }}/terraform_{{ .VersionNum }}_{{ .OS }}_{{ .Arch }}.zip
  jf:
    github: jfrog/jfrog-cli
    downloadURL: https://releases.jfrog.io/artifactory/jfrog-cli/v2-jf/{{ .VersionNum }}/jfrog-cli-{{ .OS }}-{{ .Arch }}/jf{{ .FileExt }}
  kind:
    github: kubernetes-sigs/kind
  minikube:
    github: kubernetes/minikube
  gh:
    github: cli/cli
target: /home/xyz/bin

```

## Generate Makefile go tool install tasks

```text
Adds tools to a Makefile

Usage:
  toolbox makefile [tools] [flags]

Flags:
  -f, --file string   The Makefile path to generate tools in
  -h, --help          help for makefile
```

Example:

```bash
toolbox makefile -f ./Makefile \
    github.com/bakito/semver \
    github.com/golangci/golangci-lint/cmd/golangci-lint
```

```Makefile
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
	toolbox makefile -f $$(pwd)/Makefile \
		github.com/bakito/semver \
		github.com/golangci/golangci-lint/cmd/golangci-lint
## toolbox - end

```
