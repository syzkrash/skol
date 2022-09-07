//go:build DEBUG

package debug

import (
	"fmt"
	"os"
	"runtime"
)

// Log prints the given formatted debug message if the given attribute is
// enabled in [GlobalAttr].
func Log(attr Attribute, format string, args ...any) {
	if GlobalAttr&attr != attr {
		return
	}
	_, file, line, ok := runtime.Caller(1)
	var fileline string
	if !ok {
		fileline = "???"
	} else {
		fileline = fmt.Sprintf("%s:%d", file, line)
	}
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[%s] %s\n  %s\n", attr.Name(), fileline, msg)
}
