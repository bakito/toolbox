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

func (t *Toolbox) Versions() Versions {
	v := Versions{Versions: map[string]string{}}
	for _, t := range t.Tools {
		v.Versions[t.Name] = t.Version
	}
	return v
}

type Tool struct {
	Name        string   `yaml:"name"`
	Github      string   `yaml:"github"`
	Google      string   `yaml:"google"`
	DownloadURL string   `yaml:"downloadURL"`
	Version     string   `yaml:"version"`
	Additional  []string `yaml:"additional"`
}

type ToolVersion struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Versions struct {
	Versions map[string]string `yaml:"versions"`
}

func (t *Tool) LatestURL() string {
	if t.Github != "" {
		return fmt.Sprintf(latestURLPattern, t.Github)
	}
	return ""
}
