//go:build tools
// +build tools

package tools

import (
	_ "github.com/bakito/semver"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/onsi/ginkgo/v2/ginkgo"
	_ "k8s.io/code-generator/cmd/deepcopy-gen"
)
