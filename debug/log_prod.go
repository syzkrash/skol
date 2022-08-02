//go:build !DEBUG

package debug

func Log(attr Attribute, format string, args ...any) {}
