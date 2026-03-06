package quietly_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bakito/toolbox/pkg/quietly"
)

var _ = Describe("Close", func() {
	Context("CheckError", func() {
		It("Should close the Closer", func() {
			cl := &closer{}
			quietly.Close(cl)
			Ω(cl.closed).Should(BeTrue())
		})
		It("Should not fail on nil", func() {
			quietly.Close(nil)
		})
		It("Should not fail when close return an error", func() {
			cl := &closer{fail: true}
			quietly.Close(cl)
			Ω(cl.closed).Should(BeFalse())
		})
	})
})

type closer struct {
	closed bool
	fail   bool
}

func (c *closer) Close() error {
	if c.fail {
		return errors.New("failed")
	}
	c.closed = true
	return nil
}
