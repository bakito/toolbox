//go:build !windows

package main

func binaryName(name string) string {
	return name
}
