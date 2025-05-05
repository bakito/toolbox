// Package renovate defines renovate types
package renovate

const (
	CustomType            = "regex"
	DescriptionDeprecated = "Update toolbox tools in Makefile"
	Description           = "Update toolbox tools in .toolbox.mk"
	ManagerFilePatterns   = `^\.toolbox\.mk$`
	MatchString           = `# renovate: packageName=(?<packageName>.+?)\s+.+?_VERSION \?= (?<currentValue>.+?)\s`
	DatasourceTemplate    = "go"
)

func Config() CustomManager {
	return CustomManager{
		CustomType:          CustomType,
		Description:         Description,
		ManagerFilePatterns: []string{ManagerFilePatterns},
		MatchStrings:        []string{MatchString},
		DatasourceTemplate:  DatasourceTemplate,
	}
}

type CustomManagers []CustomManager

type CustomManager struct {
	CustomType          string   `json:"customType"`
	Description         string   `json:"description"`
	ManagerFilePatterns []string `json:"managerFilePatterns"`
	MatchStrings        []string `json:"matchStrings"`
	DatasourceTemplate  string   `json:"datasourceTemplate"`
}

func (m *CustomManager) UpdateParams() {
	m.Description = Description
	m.ManagerFilePatterns = []string{ManagerFilePatterns}
	m.MatchStrings = []string{MatchString}
	m.DatasourceTemplate = DatasourceTemplate
}

func (m *CustomManager) IsToolbox() bool {
	return m.CustomType == CustomType && (m.Description == Description || m.Description == DescriptionDeprecated)
}
