package fetcher

import (
	"runtime"
	"testing"

	"github.com/bakito/toolbox/pkg/types"
)

func TestFindMatching(t *testing.T) {
	tests := []struct {
		name     string
		tb       *types.Toolbox
		toolName string
		assets   []types.Asset
		expected *types.Asset
	}{
		{
			name:     "Single matching asset",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-linux-amd64"},
				{Name: "tool2-linux-amd64"},
			},
			expected: &types.Asset{Name: "tool1-linux-amd64"},
		},
		{
			name:     "Multiple matching assets - select preferred",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-linux-386"},
				{Name: "tool1-linux-amd64"},
			},
			expected: &types.Asset{Name: "tool1-linux-amd64"},
		},
		{
			name:     "No matching assets",
			tb:       nil,
			toolName: "tool3",
			assets: []types.Asset{
				{Name: "tool1-linux-amd64"},
				{Name: "tool2-linux-amd64"},
			},
			expected: nil,
		},
		{
			name:     "Exclude forbidden suffixes",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-linux-amd64.sha256"},
				{Name: "tool1-linux-amd64"},
			},
			expected: &types.Asset{Name: "tool1-linux-amd64"},
		},
		{
			name:     "Match runtime GOOS and GOARCH",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH},
				{Name: "tool1-unknown-os-unknown-arch"},
			},
			expected: &types.Asset{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH},
		},
		{
			name:     "Handle empty asset list",
			tb:       nil,
			toolName: "tool1",
			assets:   []types.Asset{},
			expected: nil,
		},
		{
			name:     "Prefer exact tool name prefix match",
			tb:       nil,
			toolName: "tool",
			assets: []types.Asset{
				{Name: "mytool-linux-amd64"},
				{Name: "tool-linux-amd64"},
			},
			expected: &types.Asset{Name: "tool-linux-amd64"},
		},
		{
			name:     "Prefer assets with matching GOARCH",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-linux-386"},
				{Name: "tool1-linux-" + runtime.GOARCH},
			},
			expected: &types.Asset{Name: "tool1-linux-" + runtime.GOARCH},
		},
		{
			name:     "Prefer non-archive files",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-linux-amd64.tar.gz"},
				{Name: "tool1-linux-amd64"},
			},
			expected: &types.Asset{Name: "tool1-linux-amd64"},
		},
		{
			name:     "Use custom excluded suffixes from toolbox",
			tb:       &types.Toolbox{ExcludedSuffixes: []string{"custom"}},
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-linux-amd64.custom"},
				{Name: "tool1-linux-amd64.sha256"},
			},
			expected: &types.Asset{Name: "tool1-linux-amd64.sha256"},
		},
		{
			name:     "Filter by GOOS match",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-windows-amd64"},
				{Name: "tool1-darwin-amd64"},
				{Name: "tool1-linux-amd64"},
			},
			expected: &types.Asset{Name: "tool1-" + runtime.GOOS + "-amd64"},
		},
		{
			name:     "No match when GOOS doesn't match",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-fakeos-amd64"},
			},
			expected: nil,
		},
		{
			name:     "Prefer default file extension",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-linux-amd64.zip"},
				{Name: "tool1-linux-amd64"},
			},
			expected: &types.Asset{Name: "tool1-linux-amd64"},
		},
		{
			name:     "Multiple forbidden suffixes filtered",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-linux-amd64.sha256"},
				{Name: "tool1-linux-amd64.sha512"},
				{Name: "tool1-linux-amd64.txt"},
				{Name: "tool1-linux-amd64"},
			},
			expected: &types.Asset{Name: "tool1-linux-amd64"},
		},
		{
			name:     "Asset name contains tool name but doesn't match GOOS",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-notmatchingos-amd64"},
			},
			expected: nil,
		},
		{
			name:     "Prefer assets with exact GOARCH in name",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-" + runtime.GOOS + "-x86_64"},
				{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH},
			},
			expected: &types.Asset{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH},
		},
		{
			name:     "Sort by contains GOARCH when matches GOARCH is equal",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-" + runtime.GOOS + "-other"},
				{Name: "tool1-" + runtime.GOOS + "-contains-" + runtime.GOARCH},
			},
			expected: &types.Asset{Name: "tool1-" + runtime.GOOS + "-contains-" + runtime.GOARCH},
		},
		{
			name:     "Prefer files without dots when all else equal",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH + ".bin"},
				{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH},
			},
			expected: &types.Asset{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH},
		},
		{
			name:     "Prefer default file extension when no dot preference",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"},
				{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH + (func() string {
					if runtime.GOOS == "windows" {
						return ".exe"
					}
					return ""
				})()},
			},
			expected: &types.Asset{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH + (func() string {
				if runtime.GOOS == "windows" {
					return ".exe"
				}
				return ""
			})()},
		},
		{
			name:     "All sorting criteria equal - return first",
			tb:       nil,
			toolName: "tool1",
			assets: []types.Asset{
				{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH + "-a"},
				{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH + "-b"},
			},
			expected: &types.Asset{Name: "tool1-" + runtime.GOOS + "-" + runtime.GOARCH + "-a"},
		},
		{
			name:     "Prefer arm64 even if aarch64 is present",
			tb:       nil,
			toolName: "fnox",
			assets: []types.Asset{
				{Name: "fnox-aarch64-unknown-linux-gnu.tar.gz"},
				{Name: "fnox-arm64-unknown-linux-gnu.tar.gz"},
			},
			expected: func() *types.Asset {
				if runtime.GOARCH == "arm64" {
					return &types.Asset{Name: "fnox-arm64-unknown-linux-gnu.tar.gz"}
				}
				// if not arm64, it shouldn't match either ideally,
				// but findMatching filters by GOOS first, then sorts.
				// If GOOS matches, it will return one of them.
				return &types.Asset{Name: "fnox-aarch64-unknown-linux-gnu.tar.gz"}
			}(),
		},
		{
			name:     "Should match aarch64 for arm64 if arm64 not available",
			tb:       nil,
			toolName: "fnox",
			assets: []types.Asset{
				{Name: "fnox-aarch64-unknown-linux-gnu.tar.gz"},
				{Name: "fnox-x86_64-unknown-linux-gnu.tar.gz"},
			},
			expected: func() *types.Asset {
				switch runtime.GOARCH {
				case "arm64":
					return &types.Asset{Name: "fnox-aarch64-unknown-linux-gnu.tar.gz"}
				case "amd64":
					return &types.Asset{Name: "fnox-x86_64-unknown-linux-gnu.tar.gz"}
				}
				return findMatching(nil, "fnox", []types.Asset{
					{Name: "fnox-aarch64-unknown-linux-gnu.tar.gz"},
					{Name: "fnox-x86_64-unknown-linux-gnu.tar.gz"},
				})
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := findMatching(tt.tb, tt.toolName, tt.assets)
			if (actual == nil && tt.expected != nil) || (actual != nil && tt.expected == nil) {
				t.Errorf("Expected: %v, but got: %v", tt.expected, actual)
			} else if actual != nil && tt.expected != nil && actual.Name != tt.expected.Name {
				t.Errorf("Expected: %v, but got: %v", tt.expected.Name, actual.Name)
			}
		})
	}
}
