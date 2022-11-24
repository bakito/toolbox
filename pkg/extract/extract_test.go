package extract_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bakito/toolbox/pkg/extract"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		_ = os.Remove(tempDir)
	})

	Context("File", func() {
		DescribeTable("Extracting the testfile",
			func(file string) {
				ok, err := extract.File(fmt.Sprintf("../../testdata/%s", file), tempDir)
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
