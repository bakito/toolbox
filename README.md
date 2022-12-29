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
## Current working directory
LOCALDIR ?= $(shell which cygpath > /dev/null 2>&1 && cygpath -m $$(pwd) || pwd)
## Location to install dependencies to
LOCALBIN ?= $(LOCALDIR)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
SEMVER ?= $(LOCALBIN)/semver
TOOLBOX ?= $(LOCALBIN)/toolbox
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen

## Tool Versions
SEMVER_VERSION ?= v1.1.3
TOOLBOX_VERSION ?= v0.2.4
CONTROLLER_GEN_VERSION ?= v0.10.0

## Tool Installer
.PHONY: semver
semver: $(SEMVER) ## Download semver locally if necessary.
$(SEMVER): $(LOCALBIN)
	test -s $(LOCALBIN)/semver || GOBIN=$(LOCALBIN) go install github.com/bakito/semver@$(SEMVER_VERSION)
.PHONY: toolbox
toolbox: $(TOOLBOX) ## Download toolbox locally if necessary.
$(TOOLBOX): $(LOCALBIN)
	test -s $(LOCALBIN)/toolbox || GOBIN=$(LOCALBIN) go install github.com/bakito/toolbox@$(TOOLBOX_VERSION)
.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_GEN_VERSION)

## Update Tools
.PHONY: update-toolbox-tools
update-toolbox-tools:
	@rm -f \
		$(LOCALBIN)/semver \
		$(LOCALBIN)/toolbox \
		$(LOCALBIN)/controller-gen
	toolbox makefile -f $(LOCALDIR)/Makefile \
		github.com/bakito/semver \
		github.com/bakito/toolbox \
		sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools
## toolbox - end

```
