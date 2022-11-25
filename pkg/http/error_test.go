package http

import (
	"errors"
	"net"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error", func() {
	var err error
	BeforeEach(func() {
		err = errors.New("test")
	})
	Context("CheckError", func() {
		It("should return the same error", func() {
			err2 := CheckError(err)
			Ω(err2).Should(BeIdenticalTo(err))
		})
		It("should return the same url.Error", func() {
			urlErr := &url.Error{Err: err}
			err2 := CheckError(urlErr)
			Ω(err2).Should(BeIdenticalTo(urlErr))
		})
		It("should return the url.Error if wrong OpError", func() {
			urlErr := &url.Error{Err: &net.OpError{Op: "foo"}}
			err2 := CheckError(urlErr)
			Ω(err2).Should(BeIdenticalTo(urlErr))
		})

		It("should log fatal error", func() {
			var logformat string
			logFatalf = func(format string, v ...any) {
				logformat = format
			}
			urlErr := &url.Error{Err: &net.OpError{Op: dialOperation}}
			err2 := CheckError(urlErr)
			Ω(err2).Should(BeNil())
			Ω(logformat).Should(Equal(msgFormat))
		})
	})
})
