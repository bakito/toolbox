package types

import (
	"fmt"
)

const (
	latestURLPattern = "https://api.github.com/repos/%s/releases/latest"
)

type Toolbox struct {
	Tools        map[string]*Tool     `yaml:"tools"`
	Target       string               `yaml:"target"`
	CreateTarget *bool                `yaml:"createTarget"`
	Aliases      *map[string][]string `yaml:"aliases"`
}

type Tool struct {
	Name        string   `yaml:"name"`
	Github      string   `yaml:"github"`
	Google      string   `yaml:"google"`
	DownloadURL string   `yaml:"downloadURL"`
	Version     string   `yaml:"version"`
	Additional  []string `yaml:"additional"`
}

func (t *Tool) LatestURL() string {
	if t.Github != "" {
		return fmt.Sprintf(latestURLPattern, t.Github)
	}
	return ""
}
