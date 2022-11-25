package types_test

import (
	"github.com/bakito/toolbox/pkg/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	Context("Toolbox", func() {
		var tb *types.Toolbox
		BeforeEach(func() {
			tb = &types.Toolbox{}
		})
		Context("GetTools", func() {
			It("should return an empty slice", func() {
				tools := tb.GetTools()
				Ω(tools).Should(BeEmpty())
			})
			It("should return an a sorted slice", func() {
				t1 := &types.Tool{Name: "xyz"}
				t2 := &types.Tool{Name: "abc"}
				tb.Tools = map[string]*types.Tool{t1.Name: t1, t2.Name: t2, "foo": {}}
				tools := tb.GetTools()
				Ω(tools).Should(HaveLen(3))
				Ω(tools[0].Name).Should(Equal("abc"))
				Ω(tools[1].Name).Should(Equal("foo"))
				Ω(tools[2].Name).Should(Equal("xyz"))
			})
		})
		Context("Versions", func() {
			It("should return an empty map", func() {
				versions := tb.Versions()
				Ω(versions.Versions).Should(BeEmpty())
			})
			It("should return an a sorted slice", func() {
				t1 := &types.Tool{Name: "xyz", CouldNotBeFound: true}
				t2 := &types.Tool{Name: "abc", Version: "v1.0.0"}
				tb.Tools = map[string]*types.Tool{t1.Name: t1, t2.Name: t2, "foo": {Version: "v1.2.3"}}
				versions := tb.Versions()
				Ω(versions.Versions).Should(HaveLen(2))
				Ω(versions.Versions).Should(HaveKey("abc"))
				Ω(versions.Versions["abc"]).Should(Equal("v1.0.0"))
				Ω(versions.Versions).Should(HaveKey("foo"))
				Ω(versions.Versions["foo"]).Should(Equal("v1.2.3"))
			})
		})
	})
})
