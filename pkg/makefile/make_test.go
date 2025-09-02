package makefile

import (
	"os"
	"path/filepath"

	"github.com/bakito/toolbox/pkg/github"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

const testDataDir = "../../testdata"

var _ = Describe("Make", func() {
	var (
		tempDir         string
		makeFilePath    string
		includeFilePath string
	)
	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "toolbox_make_test_")
		Ω(err).ShouldNot(HaveOccurred())
		getRelease = func(client *resty.Client, repo string, quiet bool) (*types.GithubRelease, error) {
			return &types.GithubRelease{TagName: "v0.2.1"}, nil
		}
		makeFilePath = copyFile("Makefile.content", tempDir)
		includeFilePath = filepath.Join(tempDir, includeFileName)

		format.TruncatedDiff = false
	})
	AfterEach(func() {
		_ = os.RemoveAll(tempDir)
		getRelease = github.LatestRelease
		format.TruncatedDiff = true
	})
	Context("Generate", func() {
		It("should generateForTools a correct output", func() {
			err := Generate(resty.New(), makeFilePath, false, false, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint?--version",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(readFile(makeFilePath)).Should(Equal(readFile(testDataDir, "Makefile.content.expected")))
			Ω(readFile(includeFilePath) + "\n").Should(Equal(readFile(testDataDir, ".toolbox.mk.content.expected")))
		})
		It("should migrate to include correct output", func() {
			makeFilePath = copyFile("Makefile.content.migrate", tempDir)
			err := Generate(resty.New(), makeFilePath, false, false, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(readFile(makeFilePath)).Should(Equal(readFile(testDataDir, "Makefile.content.expected")))
			Ω(readFile(includeFilePath) + "\n").Should(Equal(readFile(testDataDir, ".toolbox.mk.content.expected")))
		})
		It("should generateForTools a correct output with hybrid tools", func() {
			err := Generate(resty.New(), makeFilePath, false, false,
				filepath.Join(testDataDir, "tools.go.tst"),
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/bakito/toolbox",
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readFile(makeFilePath)).Should(Equal(readFile(testDataDir, "Makefile.content.expected")))
			Ω(readFile(includeFilePath) + "\n").Should(Equal(readFile(testDataDir, ".toolbox.mk.hybrid.expected")))
		})
		It("should generateForTools a correct output with renovate enabled", func() {
			err := Generate(resty.New(), makeFilePath, true, false, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readFile(makeFilePath)).Should(Equal(readFile(testDataDir, "Makefile.content.expected")))
			Ω(readFile(includeFilePath) + "\n").Should(Equal(readFile(testDataDir, ".toolbox.mk.renovate.expected")))
		})

		It("should generateForTools a correct output with toolchain enabled", func() {
			err := Generate(resty.New(), makeFilePath, false, true, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			Ω(err).ShouldNot(HaveOccurred())

			Ω(readFile(makeFilePath)).Should(Equal(readFile(testDataDir, "Makefile.content.expected")))
			Ω(readFile(includeFilePath) + "\n").Should(Equal(readFile(testDataDir, ".toolbox.mk.toolchain.expected")))
		})
	})

	Context("generateForTools", func() {
		It("should generateForTools a correct output", func() {
			td := []toolData{
				dataForTool(true, "sigs.k8s.io/controller-tools/cmd/controller-gen"),
				dataForTool(true, "github.com/bakito/semver"),
				dataForTool(true, "github.com/bakito/toolbox"),
			}
			err := generateForTools(resty.New(), makeFilePath, false, false, nil, td)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(readFile(makeFilePath)).Should(Equal(readFile(testDataDir, "Makefile.content.expected")))

			Ω(err).ShouldNot(HaveOccurred())
			Ω(readFile(includeFilePath) + "\n").Should(EqualDiff(readFile(testDataDir, ".toolbox.mk.tools.go.expected")))
		})
	})
	Context("updateRenovateConfInternal", func() {
		It("should add a customManagers section", func() {
			withRenovate, cfg, err := updateRenovateConfInternal(
				filepath.Join(testDataDir, "renovate.no-managers.json"),
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(withRenovate).Should(BeTrue())
			Ω(string(cfg)).Should(Equal(readFile(testDataDir, "renovate.no-managers.expected.json")))
		})
		It("should add the toolbox customManager", func() {
			withRenovate, cfg, err := updateRenovateConfInternal(
				filepath.Join(testDataDir, "renovate.other-managers.json"),
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(withRenovate).Should(BeTrue())
			Ω(string(cfg)).Should(Equal(readFile(testDataDir, "renovate.other-managers.expected.json")))
		})
		It("should add the toolbox customManager with fileMatch", func() {
			withRenovate, cfg, err := updateRenovateConfInternal(
				filepath.Join(testDataDir, "renovate.other-managers-fileMatch.json"),
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(withRenovate).Should(BeTrue())
			Ω(string(cfg)).Should(Equal(readFile(testDataDir, "renovate.other-managers.expected.json")))
		})
		It("should update the toolbox customManager", func() {
			withRenovate, cfg, err := updateRenovateConfInternal(
				filepath.Join(testDataDir, "renovate.incorrect-managers.json"),
			)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(withRenovate).Should(BeTrue())
			Ω(string(cfg)).Should(Equal(readFile(testDataDir, "renovate.incorrect-managers.expected.json")))
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
