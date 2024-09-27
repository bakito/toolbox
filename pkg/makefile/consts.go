package makefile

import (
	_ "embed"
	"fmt"
)

const (
	markerStart = "## toolbox - start"
	markerEnd   = "## toolbox - end"

	renovateConfig = `{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "customManagers": [
    {
      "customType": "regex",
      "description": "Update toolbox _VERSION variables in Makefile",
      "fileMatch": [
        "^Makefile$"
      ],
      "matchStrings": [
        "# renovate: packageName=(?<packageName>.+?)\\s+.+?_VERSION \\?= (?<currentValue>.+?)\\s"
      ],
      "datasourceTemplate": "go"
    }
  ]
}
`
)

var (
	//go:embed Makefile.tpl
	tpl              string
	makefileTemplate = fmt.Sprintf("%s\n%s%s", markerStart, tpl, markerEnd)
)
