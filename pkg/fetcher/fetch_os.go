package fetcher

import (
	"os"
	"runtime"
)

const (
	envToolboxOS   = "TOOLBOX_OS"
	envToolboxArch = "TOOLBOX_ARCH"
)

// goOs allow overwriting os
func goOs() string {
	if e, ok := os.LookupEnv(envToolboxOS); ok {
		return e
	}
	return runtime.GOOS
}

// goArch allow overwriting arch
func goArch() string {
	if e, ok := os.LookupEnv(envToolboxArch); ok {
		return e
	}
	return runtime.GOOS
}
