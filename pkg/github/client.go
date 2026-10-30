// Package github API functions
package github

import (
	"fmt"
	"log"
	http2 "net/http"
	"os"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/bakito/toolbox/pkg/http"
	"github.com/bakito/toolbox/pkg/types"
)

const (
	EnvGithubToken = "GITHUB_TOKEN" // #nosec G101: variable name for token

	toolboxRepo           = "bakito/toolbox"
	toolboxReleasePattern = `^v\d+\.\d+\.\d+$`
)

var (
	releaseURLPattern       = "https://api.github.com/repos/%s/releases/tags/%s"
	latestReleaseURLPattern = "https://api.github.com/repos/%s/releases/latest"
	latestTagURLPattern     = "https://api.github.com/repos/%s/tags"
	releasesURLPattern      = "https://api.github.com/repos/%s/releases?per_page=100"
)

// LatestRelease returns the latest release with a valid semver version.
// If releasePattern is defined, only versions matching this regex pattern are considered.
// For the toolbox itself, a default pattern of plain vX.Y.Z versions is used if not otherwise defined.
func LatestRelease(
	client *resty.Client,
	repo, token string,
	quiet bool,
	releasePattern string,
) (*types.GithubRelease, error) {
	if releasePattern == "" && strings.EqualFold(repo, toolboxRepo) {
		releasePattern = toolboxReleasePattern
	}

	var pattern *regexp.Regexp
	if releasePattern != "" {
		var err error
		pattern, err = regexp.Compile(releasePattern)
		if err != nil {
			return nil, fmt.Errorf("invalid releasePattern %q for %s: %w", releasePattern, repo, err)
		}
	}

	ghr := &types.GithubRelease{}
	ghErr := &types.GithubError{}
	ghc := client.R().
		SetError(ghErr).
		SetHeader("Accept", "application/json")
	handleGithubToken(ghc, token, quiet)

	var url string
	if pattern == nil {
		url = latestReleaseURL(repo)
		ghc.SetResult(ghr)
	} else {
		url = releasesURL(repo)
		ghc.SetResult(&types.GithubReleases{})
	}
	resp, err := ghc.Get(url)
	if err != nil {
		return nil, http.CheckError(err)
	}
	if resp.IsError() && resp.StatusCode() != http2.StatusNotFound {
		return nil, fmt.Errorf("github request was not successful: %s (%d) %s", url, resp.StatusCode(), ghErr.Message)
	}

	if pattern != nil {
		if releases, ok := resp.Result().(*types.GithubReleases); ok && !resp.IsError() {
			if latest := releases.GetLatest(pattern); latest != nil {
				ghr = latest
			}
		}
	}

	if !types.IsReleaseVersion(ghr.TagName, pattern) {
		ghr = &types.GithubRelease{}
		ght := &types.GithubTags{}
		ghc.SetResult(ght).
			SetError(ghErr).
			SetHeader("Accept", "application/json")
		resp, err := ghc.Get(latestTagURL(repo))
		if err != nil {
			return nil, http.CheckError(err)
		}
		if resp.IsError() {
			return nil, fmt.Errorf(
				"github request was not successful: %s (%d) %s",
				url,
				resp.StatusCode(),
				ghErr.Message,
			)
		}

		if latest := ght.GetLatest(pattern); latest != nil {
			ghr.TagName = latest.Name
		}
	}

	return ghr, nil
}

func TokenSet(token string) bool {
	if token != "" {
		return true
	}
	t, ok := os.LookupEnv(EnvGithubToken)
	return ok && strings.TrimSpace(t) != ""
}

func handleGithubToken(ghc *resty.Request, token string, quiet bool) {
	t := token
	if t == "" {
		if v, ok := os.LookupEnv(EnvGithubToken); ok && strings.TrimSpace(v) != "" {
			t = v
		}
	}

	if t != "" {
		if !quiet {
			log.Print("🔑 Using github token\n")
		}
		ghc.SetAuthToken(t)
	}
}

func Release(client *resty.Client, repo, version, token string, quiet bool) (*types.GithubRelease, error) {
	ghr := &types.GithubRelease{}
	ghErr := &types.GithubError{}

	ghc := client.R().
		SetResult(ghr).
		SetError(ghErr).
		SetHeader("Accept", "application/json")

	handleGithubToken(ghc, token, quiet)

	url := releaseURL(repo, version)
	resp, err := ghc.Get(releaseURL(repo, version))
	if err != nil {
		return nil, http.CheckError(err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("github request was not successful: %s (%d) %s", url, resp.StatusCode(), ghErr.Message)
	}

	if ghr.TagName == "" {
		ght := &types.GithubTags{}
		ghc.SetResult(ght)
		_, err := ghc.Get(latestTagURL(repo))
		if err != nil {
			return nil, http.CheckError(err)
		}

		if latest := ght.GetLatest(nil); latest != nil {
			ghr.TagName = latest.Name
		}
	}

	return ghr, nil
}

func latestReleaseURL(repo string) string {
	if repo != "" {
		return fmt.Sprintf(latestReleaseURLPattern, repo)
	}
	return ""
}

func releasesURL(repo string) string {
	if repo != "" {
		return fmt.Sprintf(releasesURLPattern, repo)
	}
	return ""
}

func releaseURL(repo, version string) string {
	if repo != "" {
		return fmt.Sprintf(releaseURLPattern, repo, version)
	}
	return ""
}

func latestTagURL(repo string) string {
	if repo != "" {
		return fmt.Sprintf(latestTagURLPattern, repo)
	}
	return ""
}
