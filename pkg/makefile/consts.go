package makefile

import (
	_ "embed"
	"fmt"
)

const (
	markerStart     = "## toolbox - start"
	markerEnd       = "## toolbox - end"
	includeFileName = ".toolbox.mk"
)

var (
	//go:embed .toolbox.mk.tpl
	tpl              string
	makefileTemplate = fmt.Sprintf("%s\n%s%s\n", markerStart, tpl, markerEnd)
)
