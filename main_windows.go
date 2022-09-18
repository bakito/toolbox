//go:build windows

package main

import "strings"

func binaryName(name string) string {
	if strings.HasSuffix(name, ".exe") {
		return name
	}
	return name + ".exe"
}
