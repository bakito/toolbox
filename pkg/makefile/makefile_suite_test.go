package makefile_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMakefile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Makefile Suite")
}
