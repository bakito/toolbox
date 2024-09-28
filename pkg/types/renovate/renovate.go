package renovate

const (
	CustomType         = "regex"
	Description        = "Update toolbox tools in Makefile"
	FileMatch          = "^Makefile$"
	MatchString        = `# renovate: packageName=(?<packageName>.+?)\s+.+?_VERSION \?= (?<currentValue>.+?)\s`
	DatasourceTemplate = "go"
)

func Config() CustomManager {
	return CustomManager{
		CustomType:         CustomType,
		Description:        Description,
		FileMatch:          []string{FileMatch},
		MatchStrings:       []string{MatchString},
		DatasourceTemplate: DatasourceTemplate,
	}
}

type CustomManagers []CustomManager

type CustomManager struct {
	CustomType         string   `json:"customType"`
	Description        string   `json:"description"`
	FileMatch          []string `json:"fileMatch"`
	MatchStrings       []string `json:"matchStrings"`
	DatasourceTemplate string   `json:"datasourceTemplate"`
}

func (m *CustomManager) UpdateParams() {
	m.FileMatch = []string{FileMatch}
	m.MatchStrings = []string{MatchString}
	m.DatasourceTemplate = DatasourceTemplate
}

func (m *CustomManager) IsToolbox() bool {
	return m.CustomType == CustomType && m.Description == Description
}
