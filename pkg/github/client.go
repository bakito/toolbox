package github

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bakito/toolbox/pkg/http"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/go-resty/resty/v2"
)

const EnvGithubToken = "GITHUB_TOKEN" // #nosec G101: variable name for token

var (
	releaseURLPattern       = "https://api.github.com/repos/%s/releases/tags/%s"
	latestReleaseURLPattern = "https://api.github.com/repos/%s/releases/latest"
	latestTagURLPattern     = "https://api.github.com/repos/%s/tags"
)

func LatestRelease(client *resty.Client, repo string, quiet bool) (*types.GithubRelease, error) {
	ghr := &types.GithubRelease{}
	ghErr := &types.GithubError{}
	ghc := client.R().
		SetResult(ghr).
		SetError(ghErr).
		SetHeader("Accept", "application/json")
	handleGithubToken(ghc, quiet)
	resp, err := ghc.Get(latestReleaseURL(repo))
	if err != nil {
		return nil, http.CheckError(err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("github request was not successful: %s", ghErr.Message)
	}

	if ghr.TagName == "" {
		ght := &types.GithubTags{}
		ghc.SetResult(ght)
		_, err := ghc.Get(latestTagURL(repo))
		if err != nil {
			return nil, http.CheckError(err)
		}

		if latest := ght.GetLatest(); latest != nil {
			ghr.TagName = latest.Name
		}
	}

	return ghr, nil
}

func TokenSet() bool {
	t, ok := os.LookupEnv(EnvGithubToken)
	return ok && strings.TrimSpace(t) != ""
}

func handleGithubToken(ghc *resty.Request, quiet bool) {
	if t, ok := os.LookupEnv(EnvGithubToken); ok && strings.TrimSpace(t) != "" {
		if !quiet {
			log.Printf("🔑 Using github token\n")
		}
		ghc.SetAuthToken(t)
	}
}

func Release(client *resty.Client, repo string, version string, quiet bool) (*types.GithubRelease, error) {
	ghr := &types.GithubRelease{}
	ghErr := &types.GithubError{}

	ghc := client.R().
		SetResult(ghr).
		SetError(ghErr).
		SetHeader("Accept", "application/json")

	handleGithubToken(ghc, quiet)

	resp, err := ghc.Get(releaseURL(repo, version))
	if err != nil {
		return nil, http.CheckError(err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("github request was not successful: %s", ghErr.Message)
	}

	if ghr.TagName == "" {
		ght := &types.GithubTags{}
		ghc.SetResult(ght)
		_, err := ghc.Get(latestTagURL(repo))
		if err != nil {
			return nil, http.CheckError(err)
		}

		if latest := ght.GetLatest(); latest != nil {
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

func releaseURL(repo string, version string) string {
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
