package github

import (
	"fmt"
	"log"
	"os"

	"github.com/bakito/toolbox/pkg/http"
	"github.com/bakito/toolbox/pkg/types"
	"github.com/go-resty/resty/v2"
)

const (
	latestURLPattern = "https://api.github.com/repos/%s/releases/latest"
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
	_, err := ghc.Get(latestURL(repo))
	if err != nil {
		return nil, http.CheckError(err)
	}
	return ghr, nil
}

func latestURL(repo string) string {
	if repo != "" {
		return fmt.Sprintf(latestURLPattern, repo)
	}
	return ""
}
