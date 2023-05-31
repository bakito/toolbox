//go:build !windows

package fetcher

func binaryName(name string) string {
	return name
}

func defaultFileExtension() string {
	return ""
}
