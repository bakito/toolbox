package makefile

import (
	_ "embed"
	"fmt"
)

const (
	markerStart = "## toolbox - start"
	markerEnd   = "## toolbox - end"
)

var (
	//go:embed Makefile.tpl
	tpl              string
	makefileTemplate = fmt.Sprintf("%s\n%s%s", markerStart, tpl, markerEnd)
)
