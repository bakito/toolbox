package types

import (
	"fmt"
	"sort"
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

func (t *Toolbox) GetTools() []*Tool {
	var tools []*Tool
	for n := range t.Tools {
		tool := t.Tools[n]
		if tool.Name == "" {
			tool.Name = n
		}
		tools = append(tools, tool)
	}

	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name < tools[j].Name
	})

	return tools
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
