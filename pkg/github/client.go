package github

import (
	"fmt"
	"log"
	"os"

	"github.com/bakito/toolbox/pkg/http"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/go-resty/resty/v2"
)

var (
	releaseURLPattern       = "https://api.github.com/repos/%s/releases/tags/%s"
	latestReleaseURLPattern = "https://api.github.com/repos/%s/releases/latest"
	latestTagURLPattern     = "https://api.github.com/repos/%s/tags"
)

func LatestRelease(client *resty.Client, repo string, quiet bool) (*types.GithubRelease, error) {
	ghr := &types.GithubRelease{}

	ghc := client.R().
		SetResult(ghr).
		SetHeader("Accept", "application/json")
	if t, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		if !quiet {
			log.Printf("🔑 Using github token\n")
		}
		ghc = ghc.SetAuthToken(t)
	}
	_, err := ghc.Get(latestReleaseURL(repo))
	if err != nil {
		return nil, http.CheckError(err)
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

func Release(client *resty.Client, repo string, version string, quiet bool) (*types.GithubRelease, error) {
	ghr := &types.GithubRelease{}

	ghc := client.R().
		SetResult(ghr).
		SetHeader("Accept", "application/json")
	if t, ok := os.LookupEnv("GITHUB_TOKEN"); ok {
		if !quiet {
			log.Printf("🔑 Using github token\n")
		}
		ghc = ghc.SetAuthToken(t)
	}
	_, err := ghc.Get(releaseURL(repo, version))
	if err != nil {
		return nil, http.CheckError(err)
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
