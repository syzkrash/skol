//go:build !DEBUG

package debug

// Log does nothing in production.
func Log(attr Attribute, format string, args ...any) {}
