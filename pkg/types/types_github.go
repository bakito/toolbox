package types

import (
	"regexp"
	"slices"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

type GithubError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

type GithubRelease struct {
	TagName    string  `json:"tag_name"`
	Name       string  `json:"name"`
	Assets     []Asset `json:"assets"`
	Body       string  `json:"body"`
	Draft      bool    `json:"draft"`
	Prerelease bool    `json:"prerelease"`
}

type GithubReleases []GithubRelease

// GetLatest returns the latest non draft / non prerelease release with a valid semver tag matching the optional pattern.
func (ghr GithubReleases) GetLatest(pattern *regexp.Regexp) *GithubRelease {
	var latest *GithubRelease
	for i := range ghr {
		rel := &ghr[i]
		if rel.Draft || rel.Prerelease || !IsReleaseVersion(rel.TagName, pattern) {
			continue
		}
		if latest == nil || semver.Compare(rel.TagName, latest.TagName) > 0 {
			latest = rel
		}
	}
	return latest
}

// IsReleaseVersion checks if the version is a valid semver release version (no prerelease) matching the optional pattern.
func IsReleaseVersion(version string, pattern *regexp.Regexp) bool {
	if !semver.IsValid(version) || strings.Contains(version, "-") {
		return false
	}
	return pattern == nil || pattern.MatchString(version)
}

type Asset struct {
	URL                string    `json:"url"`
	ID                 int       `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              any       `json:"label"`
	Uploader           Uploader  `json:"uploader"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}

type Uploader struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type GithubTags []GithubTag

// GetLatest returns the latest valid semver release tag matching the optional pattern.
func (ght GithubTags) GetLatest(pattern *regexp.Regexp) *GithubTag {
	slices.SortFunc(ght, func(a, b GithubTag) int {
		return semver.Compare(b.Name, a.Name)
	})

	for _, tag := range ght {
		if IsReleaseVersion(tag.Name, pattern) {
			return &tag
		}
	}
	return nil
}

type GithubTag struct {
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     Commit `json:"commit"`
	NodeID     string `json:"node_id"`
}
type Commit struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}
