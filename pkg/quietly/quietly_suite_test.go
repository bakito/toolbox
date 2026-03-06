package quietly_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestQuietly(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Quietly Suite")
}
