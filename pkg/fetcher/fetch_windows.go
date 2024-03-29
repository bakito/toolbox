//go:build windows

package fetcher

import "strings"

func binaryName(name string) string {
	if strings.HasSuffix(name, defaultFileExtension()) {
		return name
	}
	return name + defaultFileExtension()
}

func defaultFileExtension() string {
	return ".exe"
}
