package types

import (
	"sort"
	"strings"
)

type Toolbox struct {
	Tools            map[string]*Tool     `yaml:"tools,omitempty"`
	Target           string               `yaml:"target,omitempty"`
	Upx              bool                 `yaml:"upx,omitempty"`
	CreateTarget     *bool                `yaml:"createTarget,omitempty"`
	Aliases          *map[string][]string `yaml:"aliases,omitempty"`
	ExcludedSuffixes []string             `yaml:"excludedSuffixes,omitempty"`
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
		if !t.CouldNotBeFound && !t.Invalid {
			v.Versions[t.Name] = t.Version
		}
	}
	return v
}

func (t *Toolbox) HasGithubTools() bool {
	for _, tool := range t.Tools {
		if tool.Github != "" {
			return true
		}
	}
	return false
}

type Tool struct {
	Name            string   `yaml:"name,omitempty"`
	Github          string   `yaml:"github,omitempty"`
	Google          string   `yaml:"google,omitempty"`
	DownloadURL     string   `yaml:"downloadURL,omitempty"`
	Version         string   `yaml:"version,omitempty"`
	Additional      []string `yaml:"additional,omitempty"`
	Check           string   `yaml:"check,omitempty"`
	SkipUpx         bool     `yaml:"skipUpx,omitempty"`
	CouldNotBeFound bool     `yaml:"-"`
	Invalid         bool     `yaml:"-"`
}

type ToolVersion struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Versions struct {
	Versions map[string]string `yaml:"versions"`
}
