package extract_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bakito/toolbox/pkg/extract"
)

func TestExtract(t *testing.T) {
	tempDir := t.TempDir()

	testFile, err := os.ReadFile("../../testdata/testfile")
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	tests := []struct {
		name string
		file string
		want bool
	}{
		{name: "It should extract a simple zip file", file: "testfile.zip", want: true},
		{name: "It should extract a zip file with directories", file: "testfile-dirs.zip", want: true},
		{name: "It should extract a simple tar.gz file", file: "testfile.tar.gz", want: true},
		{name: "It should extract a tar.gz file with directories", file: "testfile-dirs.tar.gz", want: true},
		{name: "It should extract a simple tar.xz file", file: "testfile.tar.xz", want: true},
		{name: "It should extract a tar.xz file with directories", file: "testfile-dirs.tar.xz", want: true},
		{name: "should not know the extension", file: "a.txt", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean tempDir for each run
			runTempDir := filepath.Join(tempDir, tt.file)
			_ = os.MkdirAll(runTempDir, 0o755)

			ok, err := extract.File("../../testdata/"+tt.file, runTempDir)
			if err != nil {
				t.Errorf("extract.File() error = %v", err)
				return
			}
			if ok != tt.want {
				t.Errorf("extract.File() ok = %v, want %v", ok, tt.want)
				return
			}

			if ok {
				files, err := findFiles(runTempDir)
				if err != nil {
					t.Fatalf("findFiles() error = %v", err)
				}
				if len(files) != 1 {
					t.Fatalf("expected 1 file, got %v", len(files))
				}

				archiveFile, err := os.ReadFile(files[0])
				if err != nil {
					t.Fatalf("failed to read archived file: %v", err)
				}
				if diff := cmp.Diff(testFile, archiveFile); diff != "" {
					t.Errorf("archiveFile mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

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
