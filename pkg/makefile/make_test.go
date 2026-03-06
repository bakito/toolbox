package makefile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-resty/resty/v2"

	"github.com/bakito/toolbox/pkg/github"
	"github.com/bakito/toolbox/pkg/types"
)

const testDataDir = "../../testdata"

func TestMake(t *testing.T) {
	t.Run("Generate", func(t *testing.T) {
		t.Run("should generateForTools a correct output", func(t *testing.T) {
			tempDir, makeFilePath, includeFilePath := setupTest(t)
			defer os.RemoveAll(tempDir)

			err := Generate(resty.New(), makeFilePath, false, false, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			if err != nil {
				t.Fatal(err)
			}

			assertEqualFileDiff(t, makeFilePath, testDataDir, "Makefile.content.expected")
			assertEqualFileDiff(t, includeFilePath, testDataDir, ".toolbox.mk.content.expected")
		})
		t.Run("should migrate to include correct output", func(t *testing.T) {
			tempDir, _, includeFilePath := setupTest(t)
			defer os.RemoveAll(tempDir)

			makeFilePath := copyFile(t, "Makefile.content.migrate", tempDir)
			err := Generate(resty.New(), makeFilePath, false, false, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			if err != nil {
				t.Fatal(err)
			}

			assertEqualFileDiff(t, makeFilePath, testDataDir, "Makefile.content.expected")
			assertEqualFileDiff(t, includeFilePath, testDataDir, ".toolbox.mk.content.expected")
		})
		t.Run("should generateForTools a correct output with hybrid tools", func(t *testing.T) {
			tempDir, makeFilePath, includeFilePath := setupTest(t)
			defer os.RemoveAll(tempDir)

			err := Generate(resty.New(), makeFilePath, false, false,
				filepath.Join(testDataDir, "tools.go.tst"),
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/bakito/toolbox",
			)
			if err != nil {
				t.Fatal(err)
			}
			assertEqualFileDiff(t, makeFilePath, testDataDir, "Makefile.content.expected")
			assertEqualFileDiff(t, includeFilePath, testDataDir, ".toolbox.mk.hybrid.expected")
		})
		t.Run("should generateForTools a correct output with renovate enabled", func(t *testing.T) {
			tempDir, makeFilePath, includeFilePath := setupTest(t)
			defer os.RemoveAll(tempDir)

			err := Generate(resty.New(), makeFilePath, true, false, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			if err != nil {
				t.Fatal(err)
			}
			assertEqualFileDiff(t, makeFilePath, testDataDir, "Makefile.content.expected")
			assertEqualFileDiff(t, includeFilePath, testDataDir, ".toolbox.mk.renovate.expected")
		})

		t.Run("should generateForTools a correct output with toolchain enabled", func(t *testing.T) {
			tempDir, makeFilePath, includeFilePath := setupTest(t)
			defer os.RemoveAll(tempDir)

			err := Generate(resty.New(), makeFilePath, false, true, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
			)
			if err != nil {
				t.Fatal(err)
			}

			assertEqualFileDiff(t, makeFilePath, testDataDir, "Makefile.content.expected")
			assertEqualFileDiff(t, includeFilePath, testDataDir, ".toolbox.mk.toolchain.expected")
		})
		t.Run("should generateForTools a correct output with versioned tool", func(t *testing.T) {
			tempDir, makeFilePath, includeFilePath := setupTest(t)
			defer os.RemoveAll(tempDir)

			err := Generate(resty.New(), makeFilePath, false, false, "",
				"sigs.k8s.io/controller-tools/cmd/controller-gen@github.com/kubernetes-sigs/controller-tools",
				"github.com/golangci/golangci-lint/v2/cmd/golangci-lint?--version",
				"github.com/bakito/semver",
				"github.com/bakito/toolbox",
				"go.uber.org/mock/mockgen@github.com/uber-go/mock?--version",
			)
			if err != nil {
				t.Fatal(err)
			}

			assertEqualFileDiff(t, makeFilePath, testDataDir, "Makefile.content.expected")
			assertEqualFileDiff(t, includeFilePath, testDataDir, ".toolbox.mk.version.expected")
		})
	})

	t.Run("generateForTools", func(t *testing.T) {
		t.Run("should generateForTools a correct output", func(t *testing.T) {
			tempDir, makeFilePath, includeFilePath := setupTest(t)
			defer os.RemoveAll(tempDir)

			td := []toolData{
				dataForTool(true, "sigs.k8s.io/controller-tools/cmd/controller-gen"),
				dataForTool(true, "github.com/bakito/semver"),
				dataForTool(true, "github.com/bakito/toolbox"),
			}
			err := generateForTools(resty.New(), makeFilePath, false, false, nil, td)
			if err != nil {
				t.Fatal(err)
			}
			assertEqualFileDiff(t, makeFilePath, testDataDir, "Makefile.content.expected")
			assertEqualFileDiff(t, includeFilePath, testDataDir, ".toolbox.mk.tools.go.expected")
		})
	})
	t.Run("updateRenovateConfInternal", func(t *testing.T) {
		t.Run("should add a customManagers section", func(t *testing.T) {
			withRenovate, cfg, err := updateRenovateConfInternal(
				filepath.Join(testDataDir, "renovate.no-managers.json"),
			)
			if err != nil {
				t.Fatal(err)
			}
			if !withRenovate {
				t.Error("expected withRenovate to be true")
			}
			assertEqualDiff(t, string(cfg), readFileSimple(t, testDataDir, "renovate.no-managers.expected.json"))
		})
		t.Run("should add the toolbox customManager", func(t *testing.T) {
			withRenovate, cfg, err := updateRenovateConfInternal(
				filepath.Join(testDataDir, "renovate.other-managers.json"),
			)
			if err != nil {
				t.Fatal(err)
			}
			if !withRenovate {
				t.Error("expected withRenovate to be true")
			}
			assertEqualDiff(t, string(cfg), readFileSimple(t, testDataDir, "renovate.other-managers.expected.json"))
		})
		t.Run("should add the toolbox customManager with fileMatch", func(t *testing.T) {
			withRenovate, cfg, err := updateRenovateConfInternal(
				filepath.Join(testDataDir, "renovate.other-managers-fileMatch.json"),
			)
			if err != nil {
				t.Fatal(err)
			}
			if !withRenovate {
				t.Error("expected withRenovate to be true")
			}
			assertEqualDiff(t, string(cfg), readFileSimple(t, testDataDir, "renovate.other-managers.expected.json"))
		})
		t.Run("should update the toolbox customManager", func(t *testing.T) {
			withRenovate, cfg, err := updateRenovateConfInternal(
				filepath.Join(testDataDir, "renovate.incorrect-managers.json"),
			)
			if err != nil {
				t.Fatal(err)
			}
			if !withRenovate {
				t.Error("expected withRenovate to be true")
			}
			assertEqualDiff(t, string(cfg), readFileSimple(t, testDataDir, "renovate.incorrect-managers.expected.json"))
		})
	})
}

func setupTest(t *testing.T) (tempDir, makeFilePath, includeFilePath string) {
	t.Helper()
	getRelease = func(*resty.Client, string, bool) (*types.GithubRelease, error) {
		return &types.GithubRelease{TagName: "v0.2.1"}, nil
	}
	makeFilePath = copyFile(t, "Makefile.content", tempDir)
	includeFilePath = filepath.Join(tempDir, includeFileName)

	originalGetRelease := getRelease
	t.Cleanup(func() {
		getRelease = github.LatestRelease
		_ = originalGetRelease
	})

	return tempDir, makeFilePath, includeFilePath
}

func copyFile(t *testing.T, name, targetDir string) string {
	t.Helper()
	bytesRead, err := os.ReadFile(filepath.Join(testDataDir, name))
	if err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(targetDir, name)
	err = os.WriteFile(path, bytesRead, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func readFileSimple(t *testing.T, path ...string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(path...))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
