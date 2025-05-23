package extract_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bakito/toolbox/pkg/extract"
)

var _ = Describe("Extract", func() {
	var (
		tempDir  string
		testFile []byte
	)
	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "toolbox_extract_test_")
		Ω(err).ShouldNot(HaveOccurred())
		testFile, err = os.ReadFile("../../testdata/testfile")
		Ω(err).ShouldNot(HaveOccurred())
	})
	AfterEach(func() {
		_ = os.RemoveAll(tempDir)
	})

	Context("File", func() {
		DescribeTable("Extracting the testfile",
			func(file string) {
				ok, err := extract.File("../../testdata/"+file, tempDir)
				Ω(ok).Should(BeTrue())
				Ω(err).ShouldNot(HaveOccurred())
				files, err := findFiles(tempDir)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(files).Should(HaveLen(1))

				archiveFile, err := os.ReadFile(files[0])
				Ω(err).ShouldNot(HaveOccurred())
				Ω(archiveFile).Should(Equal(testFile))
			},
			Entry("It should extract a simple zip file", "testfile.zip"),
			Entry("It should extract a zip file with directories", "testfile-dirs.zip"),
			Entry("It should extract a simple tar.gz file", "testfile.tar.gz"),
			Entry("It should extract a tar.gz file with directories", "testfile-dirs.tar.gz"),
			Entry("It should extract a simple tar.xz file", "testfile.tar.xz"),
			Entry("It should extract a tar.xz file with directories", "testfile-dirs.tar.xz"),
		)
		It("should not know the extension", func() {
			ok, err := extract.File("a.txt", tempDir)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(ok).Should(BeFalse())
		})
	})
})

func findFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}

		return nil
	})
	return files, err
}
