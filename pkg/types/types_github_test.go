package types_test

import (
	"regexp"
	"testing"

	"github.com/bakito/toolbox/pkg/types"
)

func TestGithubTags_GetLatest(t *testing.T) {
	tags := func() types.GithubTags {
		return types.GithubTags{
			{Name: "latest"},
			{Name: "v1.2.0"},
			{Name: "v2.0.0-rc1"},
			{Name: "v1.10.1"},
			{Name: "release-3"},
			{Name: "v1.9.0"},
		}
	}
	tests := []struct {
		name    string
		pattern *regexp.Regexp
		want    string
	}{
		{name: "ignores non semver tags", want: "v1.10.1"},
		{name: "considers pattern", pattern: regexp.MustCompile(`^v1\.\d\.`), want: "v1.9.0"},
		{name: "no match", pattern: regexp.MustCompile(`^v3`), want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ""
			if l := tags().GetLatest(tt.pattern); l != nil {
				got = l.Name
			}
			if got != tt.want {
				t.Errorf("GetLatest() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGithubReleases_GetLatest(t *testing.T) {
	releases := types.GithubReleases{
		{TagName: "nightly"},
		{TagName: "v1.2.0"},
		{TagName: "v1.3.0", Draft: true},
		{TagName: "v1.4.0", Prerelease: true},
		{TagName: "v2.1.0"},
		{TagName: "v1.2.5"},
	}
	tests := []struct {
		name    string
		pattern *regexp.Regexp
		want    string
	}{
		{name: "latest without pattern", want: "v2.1.0"},
		{name: "latest matching pattern", pattern: regexp.MustCompile(`^v1\.`), want: "v1.2.5"},
		{name: "no match", pattern: regexp.MustCompile(`^nightly$`), want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ""
			if l := releases.GetLatest(tt.pattern); l != nil {
				got = l.TagName
			}
			if got != tt.want {
				t.Errorf("GetLatest() = %q, want %q", got, tt.want)
			}
		})
	}
}
