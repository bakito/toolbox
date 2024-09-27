package makefile

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/bakito/toolbox/pkg/github"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testDataDir = "../../testdata"

var _ = Describe("Make", func() {
	var tempDir string
	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "toolbox_make_test_")
		Ω(err).ShouldNot(HaveOccurred())
		getRelease = func(client *resty.Client, repo string, quiet bool) (*types.GithubRelease, error) {
			return &types.GithubRelease{TagName: "v0.2.1"}, nil
		}
	})
	AfterEach(func() {
		_ = os.RemoveAll(tempDir)
		getRelease = github.LatestRelease
	})
	Context("Generate", func() {
		It("should generate a correct output", func() {
			out := &bytes.Buffer{}
			err := Generate(resty.New(), out, "", "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(out.String() + "\n").Should(Equal(readFile(testDataDir, "Makefile.expected")))
		})
		It("should generate a correct output", func() {
			out := &bytes.Buffer{}
			path := copyFile("Makefile.content", tempDir)
			err := Generate(resty.New(), out, path, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(out.Bytes()).Should(BeEmpty())

			Ω(readFile(path)).Should(Equal(readFile(testDataDir, "Makefile.content.expected")))
		})
		It("should generate a correct output", func() {
			out := &bytes.Buffer{}

			err := Generate(resty.New(), out, "",
				filepath.Join(testDataDir, "tools.go.tst"),
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/bakito/toolbox",
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(out.String() + "\n").Should(Equal(readFile(testDataDir, "Makefile.hybrid.expected")))
		})
	})
	Context("generate", func() {
		It("should generate a correct output", func() {
			out := &bytes.Buffer{}

			td := []toolData{
				dataForTool(true, "sigs.k8s.io/controller-tools/cmd/controller-gen"),
				dataForTool(true, "github.com/bakito/semver"),
				dataForTool(true, "github.com/bakito/toolbox"),
			}
			err := generate(resty.New(), out, "", nil, td)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(out.String() + "\n").Should(Equal(readFile(testDataDir, "Makefile.tools.go.expected")))
		})
	})
	Context("PrintRenovateConfig", func() {
		It("should generate a correct output", func() {
			out := &bytes.Buffer{}
			PrintRenovateConfig(out)
			Ω(out.String()).Should(Equal(renovateConfig))
		})
	})
})

func readFile(path ...string) string {
	b, err := os.ReadFile(filepath.Join(path...))
	Ω(err).ShouldNot(HaveOccurred())
	return string(b)
}

func copyFile(name, targetDir string) string {
	bytesRead, err := os.ReadFile(filepath.Join(testDataDir, name))

	Ω(err).ShouldNot(HaveOccurred())

	path := filepath.Join(targetDir, name)
	err = os.WriteFile(path, bytesRead, 0o600)

	Ω(err).ShouldNot(HaveOccurred())
	return path
}
