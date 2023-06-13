package types

import (
	"sort"
	"strings"
)

type Toolbox struct {
	Tools        map[string]*Tool     `yaml:"tools,omitempty"`
	Target       string               `yaml:"target,omitempty"`
	CreateTarget *bool                `yaml:"createTarget,omitempty"`
	Aliases      *map[string][]string `yaml:"aliases,omitempty"`
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
		return strings.ToLower(tools[i].Name) < strings.ToLower(tools[j].Name)
	})

	return tools
}

func (t *Toolbox) Versions() *Versions {
	v := &Versions{Versions: map[string]string{}}
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
	Name            string   `yaml:"name,omitempty"`
	Github          string   `yaml:"github,omitempty"`
	Google          string   `yaml:"google,omitempty"`
	DownloadURL     string   `yaml:"downloadURL,omitempty"`
	Version         string   `yaml:"version,omitempty"`
	Additional      []string `yaml:"additional,omitempty"`
	CouldNotBeFound bool     `yaml:"-"`
}

type ToolVersion struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Versions struct {
	Versions map[string]string `yaml:"versions"`
}
