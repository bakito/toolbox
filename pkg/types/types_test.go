package types_test

import (
	"reflect"
	"testing"

	"github.com/bakito/toolbox/pkg/types"
)

func TestToolbox(t *testing.T) {
	t.Run("GetTools", func(t *testing.T) {
		t.Run("should return an empty slice", func(t *testing.T) {
			tb := &types.Toolbox{}
			tools := tb.GetTools()
			if len(tools) != 0 {
				t.Errorf("expected 0 tools, got %d", len(tools))
			}
		})
		t.Run("should return an a sorted slice", func(t *testing.T) {
			tb := &types.Toolbox{}
			t1 := &types.Tool{Name: "xyz", Github: "foo"}
			t2 := &types.Tool{Name: "abc", Google: "bar"}
			t3 := &types.Tool{Name: "no-source"}
			tb.Tools = map[string]*types.Tool{t1.Name: t1, t2.Name: t2, t3.Name: t3, "foo": {DownloadURL: "url"}}
			tools := tb.GetTools()
			if len(tools) != 3 {
				t.Fatalf("expected 3 tools, got %d", len(tools))
			}
			if tools[0].Name != "abc" {
				t.Errorf("expected name 'abc', got %q", tools[0].Name)
			}
			if tools[1].Name != "foo" {
				t.Errorf("expected name 'foo', got %q", tools[1].Name)
			}
			if tools[2].Name != "xyz" {
				t.Errorf("expected name 'xyz', got %q", tools[2].Name)
			}
		})
	})
	t.Run("Versions", func(t *testing.T) {
		t.Run("should return an empty map", func(t *testing.T) {
			tb := &types.Toolbox{}
			versions := tb.Versions()
			if len(versions.Versions) != 0 {
				t.Errorf("expected 0 versions, got %d", len(versions.Versions))
			}
		})
		t.Run("should return an a sorted slice", func(t *testing.T) {
			tb := &types.Toolbox{}
			t1 := &types.Tool{Name: "xyz", CouldNotBeFound: true, Github: "foo"}
			t2 := &types.Tool{Name: "abc", Version: "v1.0.0", Github: "foo"}
			tb.Tools = map[string]*types.Tool{t1.Name: t1, t2.Name: t2, "foo": {Version: "v1.2.3", Github: "foo"}}
			versions := tb.Versions()
			if len(versions.Versions) != 2 {
				t.Fatalf("expected 2 versions, got %d", len(versions.Versions))
			}

			expected := map[string]string{
				"abc": "v1.0.0",
				"foo": "v1.2.3",
			}
			if !reflect.DeepEqual(versions.Versions, expected) {
				t.Errorf("expected %v, got %v", expected, versions.Versions)
			}
		})
	})
}
