package types

import (
	"sort"
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

		if tool.Github != "" || tool.DownloadURL != "" || tool.Google != "" {
			if tool.Name == "" {
				tool.Name = n
			}
			tools = append(tools, tool)
		}
	}

	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name < tools[j].Name
	})

	return tools
}

func (t *Toolbox) Versions() Versions {
	v := Versions{Versions: map[string]string{}}
	tools := t.GetTools()
	for i := range tools {
		t := tools[i]
		if !t.CouldNotBeFound {
			v.Versions[t.Name] = t.Version
		}
	}
	return v
}

type Tool struct {
	Name            string   `yaml:"name"`
	Github          string   `yaml:"github"`
	Google          string   `yaml:"google"`
	DownloadURL     string   `yaml:"downloadURL"`
	Version         string   `yaml:"version"`
	Additional      []string `yaml:"additional"`
	CouldNotBeFound bool     `yaml:"-"`
}

type ToolVersion struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Versions struct {
	Versions map[string]string `yaml:"versions"`
}
