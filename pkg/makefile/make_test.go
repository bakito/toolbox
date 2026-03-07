package makefile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-cmp/cmp"

	"github.com/bakito/toolbox/pkg/types"
)

const testDataDir = "../../testdata"

func TestMake_Generate(t *testing.T) {
	tempDir := t.TempDir()

	originalGetRelease := getRelease
	getRelease = func(*resty.Client, string, bool) (*types.GithubRelease, error) {
		return &types.GithubRelease{TagName: "v0.2.1"}, nil
	}
	defer func() { getRelease = originalGetRelease }()

	tests := []struct {
		name            string
		setup           func(t *testing.T, tempDir string) (string, string) // returns makeFilePath, includeFilePath
		renovate        bool
		toolchain       bool
		toolsGo         string
		args            []string
		expectedMake    string
		expectedInclude string
	}{
		{
			name: "should generateForTools a correct output",
			setup: func(t *testing.T, tempDir string) (string, string) {
				t.Helper()
				return copyFile(t, "Makefile.content", tempDir), filepath.Join(tempDir, includeFileName)
			},
			args: []string{
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			},
			expectedMake:    "Makefile.content.expected",
			expectedInclude: ".toolbox.mk.content.expected",
		},
		{
			name: "should migrate to include correct output",
			setup: func(t *testing.T, tempDir string) (string, string) {
				t.Helper()
				return copyFile(t, "Makefile.content.migrate", tempDir), filepath.Join(tempDir, includeFileName)
			},
			args: []string{
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			},
			expectedMake:    "Makefile.content.expected",
			expectedInclude: ".toolbox.mk.content.expected",
		},
		{
			name: "should generateForTools a correct output with hybrid tools",
			setup: func(t *testing.T, tempDir string) (string, string) {
				t.Helper()
				return copyFile(t, "Makefile.content", tempDir), filepath.Join(tempDir, includeFileName)
			},
			toolsGo: filepath.Join(testDataDir, "tools.go.tst"),
			args: []string{
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/bakito/toolbox",
			},
			expectedMake:    "Makefile.content.expected",
			expectedInclude: ".toolbox.mk.hybrid.expected",
		},
		{
			name: "should generateForTools a correct output with renovate enabled",
			setup: func(t *testing.T, tempDir string) (string, string) {
				t.Helper()
				return copyFile(t, "Makefile.content", tempDir), filepath.Join(tempDir, includeFileName)
			},
			renovate: true,
			args: []string{
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			},
			expectedMake:    "Makefile.content.expected",
			expectedInclude: ".toolbox.mk.renovate.expected",
		},
		{
			name: "should generateForTools a correct output with toolchain enabled",
			setup: func(t *testing.T, tempDir string) (string, string) {
				t.Helper()
				return copyFile(t, "Makefile.content", tempDir), filepath.Join(tempDir, includeFileName)
			},
			toolchain: true,
			args: []string{
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			},
			expectedMake:    "Makefile.content.expected",
			expectedInclude: ".toolbox.mk.toolchain.expected",
		},
		{
			name: "should generateForTools a correct output with versioned tool",
			setup: func(t *testing.T, tempDir string) (string, string) {
				t.Helper()
				return copyFile(t, "Makefile.content", tempDir), filepath.Join(tempDir, includeFileName)
			},
			args: []string{
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint?--version",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
				"go.uber.org/mock/mockgen@github.com/uber-go/mock?--version",
			},
			expectedMake:    "Makefile.content.expected",
			expectedInclude: ".toolbox.mk.version.expected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runTempDir := filepath.Join(tempDir, tt.name)
			_ = os.MkdirAll(runTempDir, 0o755)
			makeFilePath, includeFilePath := tt.setup(t, runTempDir)
			err := Generate(resty.New(), makeFilePath, tt.renovate, tt.toolchain, tt.toolsGo, tt.args...)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}

			assertEqualFiles(t, makeFilePath, filepath.Join(testDataDir, tt.expectedMake))
			assertEqualFiles(t, includeFilePath, filepath.Join(testDataDir, tt.expectedInclude))
		})
	}
}

func TestMake_generateForTools(t *testing.T) {
	tempDir := t.TempDir()

	makeFilePath := copyFile(t, "Makefile.content", tempDir)
	includeFilePath := filepath.Join(tempDir, includeFileName)

	td := []toolData{
		dataForTool(true, "sigs.k8s.io/controller-tools/cmd/controller-gen"),
		dataForTool(true, "github.com/bakito/semver"),
		dataForTool(true, "github.com/bakito/toolbox"),
	}
	err := generateForTools(resty.New(), makeFilePath, false, false, nil, td)
	if err != nil {
		t.Fatalf("generateForTools() error = %v", err)
	}
	assertEqualFiles(t, makeFilePath, filepath.Join(testDataDir, "Makefile.content.expected"))
	assertEqualFiles(t, includeFilePath, filepath.Join(testDataDir, ".toolbox.mk.tools.go.expected"))
}

func TestMake_updateRenovateConfInternal(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected string
	}{
		{
			name:     "should add a customManagers section",
			file:     "renovate.no-managers.json",
			expected: "renovate.no-managers.expected.json",
		},
		{
			name:     "should add the toolbox customManager",
			file:     "renovate.other-managers.json",
			expected: "renovate.other-managers.expected.json",
		},
		{
			name:     "should add the toolbox customManager with fileMatch",
			file:     "renovate.other-managers-fileMatch.json",
			expected: "renovate.other-managers.expected.json",
		},
		{
			name:     "should update the toolbox customManager",
			file:     "renovate.incorrect-managers.json",
			expected: "renovate.incorrect-managers.expected.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			withRenovate, cfg, err := updateRenovateConfInternal(filepath.Join(testDataDir, tt.file))
			if err != nil {
				t.Fatalf("updateRenovateConfInternal() error = %v", err)
			}
			if !withRenovate {
				t.Error("updateRenovateConfInternal() withRenovate = false, want true")
			}

			expectedContent, err := os.ReadFile(filepath.Join(testDataDir, tt.expected))
			if err != nil {
				t.Fatalf("failed to read expected file: %v", err)
			}

			if diff := cmp.Diff(string(expectedContent), string(cfg)); diff != "" {
				t.Errorf("updateRenovateConfInternal() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func copyFile(t *testing.T, name, targetDir string) string {
	t.Helper()
	bytesRead, err := os.ReadFile(filepath.Join(testDataDir, name))
	if err != nil {
		t.Fatalf("failed to read file %s: %v", name, err)
	}

	path := filepath.Join(targetDir, name)
	err = os.WriteFile(path, bytesRead, 0o600)
	if err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
	return path
}

func assertEqualFiles(t *testing.T, actualFile, expectedFile string) {
	t.Helper()
	actual, err := os.ReadFile(actualFile)
	if err != nil {
		t.Fatalf("failed to read actual file %s: %v", actualFile, err)
	}
	expected, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("failed to read expected file %s: %v", expectedFile, err)
	}

	if diff := cmp.Diff(string(expected), string(actual)); diff != "" {
		t.Errorf("File %s mismatch (-want +got):\n%s", actualFile, diff)
	}
}
