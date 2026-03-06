package makefile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func assertEqualDiff(t *testing.T, actual, expected string) {
	t.Helper()
	diff, err := unifiedDiff(expected, "Expected", actual, "Actual")
	if err != nil {
		t.Fatal(err)
	}
	if diff != "" {
		t.Errorf("Result mismatch:\n%s", diff)
	}
}

func assertEqualFileDiff(t *testing.T, actualFile string, expectedPath ...string) {
	t.Helper()
	expectedFile := filepath.Join(expectedPath...)
	actualContent := readFile(t, actualFile)
	expectedContent := readFile(t, expectedFile)

	diff, err := unifiedDiff(expectedContent, expectedFile, actualContent, actualFile)
	if err != nil {
		t.Fatal(err)
	}
	if diff != "" {
		t.Errorf("File %s mismatch:\n%s", actualFile, diff)
	}
}

func unifiedDiff(a, nameA, b, nameB string) (string, error) {
	ud := difflib.UnifiedDiff{
		FromFile: nameA,
		A:        difflib.SplitLines(a),
		ToFile:   nameB,
		B:        difflib.SplitLines(b),
		Context:  3,
	}
	return difflib.GetUnifiedDiffString(ud)
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	return string(b)
}
