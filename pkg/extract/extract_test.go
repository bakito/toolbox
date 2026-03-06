package extract_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/bakito/toolbox/pkg/extract"
)

func TestFile(t *testing.T) {
	testFile, err := os.ReadFile("../../testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		file string
	}{
		{"It should extract a simple zip file", "testfile.zip"},
		{"It should extract a zip file with directories", "testfile-dirs.zip"},
		{"It should extract a simple tar.gz file", "testfile.tar.gz"},
		{"It should extract a tar.gz file with directories", "testfile-dirs.tar.gz"},
		{"It should extract a simple tar.xz file", "testfile.tar.xz"},
		{"It should extract a tar.xz file with directories", "testfile-dirs.tar.xz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			defer os.RemoveAll(tempDir)

			ok, err := extract.File("../../testdata/"+tt.file, tempDir)
			if !ok {
				t.Error("expected ok to be true")
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			files, err := findFiles(tempDir)
			if err != nil {
				t.Fatal(err)
			}
			if len(files) != 1 {
				t.Fatalf("expected 1 file, got %d", len(files))
			}

			archiveFile, err := os.ReadFile(files[0])
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(archiveFile, testFile) {
				t.Error("extracted file content does not match testfile")
			}
		})
	}

	t.Run("should not know the extension", func(t *testing.T) {
		tempDir := t.TempDir()
		defer os.RemoveAll(tempDir)

		ok, err := extract.File("a.txt", tempDir)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if ok {
			t.Error("expected ok to be false")
		}
	})
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
