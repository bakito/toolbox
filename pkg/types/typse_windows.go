//go:build windows

package types

func (t *Tool) FileNameForOS() string {
	if t.FileNames != nil {
		return t.FileNames.Windows
	}
	return ""
}
