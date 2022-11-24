//go:build !windows

package cmd

func binaryName(name string) string {
	return name
}

func defaultFileExtension() string {
	return ""
}
