package fetcher

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadToolbox_GithubAPIToken(t *testing.T) {
	content := `
githubApiToken: my-secret-token
tools:
  gh:
    github: cli/cli
`
	tmpDir := t.TempDir()
	cfgFile := filepath.Join(tmpDir, ".toolbox.yaml")
	err := os.WriteFile(cfgFile, []byte(content), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	tb, _, err := ReadToolbox(cfgFile)
	if err != nil {
		t.Fatalf("ReadToolbox() error = %v", err)
	}

	if tb.GithubAPIToken != "my-secret-token" {
		t.Errorf("expected GithubAPIToken to be 'my-secret-token', got %q", tb.GithubAPIToken)
	}
}
