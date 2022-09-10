//go:build unix

package pe

import (
	"syscall"

	"golang.org/x/sys/unix"
)

func isatty() bool {
	_, err := unix.IoctlGetTermio(syscall.Stderr, unix.TCGETA)
	return err == nil
}
