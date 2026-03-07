package types_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bakito/toolbox/pkg/types"
)

func TestToolbox_GetTools(t *testing.T) {
	tests := []struct {
		name string
		tb   *types.Toolbox
		want []string // tool names
	}{
		{
			name: "should return an empty slice",
			tb:   &types.Toolbox{},
			want: nil,
		},
		{
			name: "should return a sorted slice",
			tb: &types.Toolbox{
				Tools: map[string]*types.Tool{
					"xyz":       {Name: "xyz", Github: "foo"},
					"abc":       {Name: "abc", Google: "bar"},
					"no-source": {Name: "no-source"},
					"foo":       {DownloadURL: "url"},
				},
			},
			want: []string{"abc", "foo", "xyz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tb.GetTools()
			var gotNames []string
			for _, tool := range got {
				gotNames = append(gotNames, tool.Name)
			}
			if tt.want == nil && len(gotNames) == 0 {
				return
			}
			if diff := cmp.Diff(tt.want, gotNames); diff != "" {
				t.Errorf("GetTools() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestToolbox_Versions(t *testing.T) {
	tests := []struct {
		name string
		tb   *types.Toolbox
		want map[string]string
	}{
		{
			name: "should return an empty map",
			tb:   &types.Toolbox{},
			want: make(map[string]string),
		},
		{
			name: "should return a map of versions",
			tb: &types.Toolbox{
				Tools: map[string]*types.Tool{
					"xyz": {Name: "xyz", CouldNotBeFound: true, Github: "foo"},
					"abc": {Name: "abc", Version: "v1.0.0", Github: "foo"},
					"foo": {Version: "v1.2.3", Github: "foo"},
				},
			},
			want: map[string]string{
				"abc": "v1.0.0",
				"foo": "v1.2.3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tb.Versions()
			if diff := cmp.Diff(tt.want, got.Versions); diff != "" {
				t.Errorf("Versions() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
