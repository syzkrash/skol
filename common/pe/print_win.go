//go:build windows

package pe

import "golang.org/x/sys/windows"

func isatty() bool {
	t, err := windows.GetFileType(windows.Stderr)
	if err != nil {
		return false
	}
	return t == windows.FILE_TYPE_CHAR
}
