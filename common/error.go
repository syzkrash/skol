package common

import (
	"fmt"

	"github.com/syzkrash/skol/ast"
)

type Printable interface {
	Print()
}

// MetaError represents an error that does not have any additional information.
// Only the node that caused the error is available.
type MetaError struct {
	msg   string
	Cause ast.MetaNode
}

func (e *MetaError) Error() string {
	return e.msg
}

func (e *MetaError) Print() {
	fmt.Println("\x1b[31m\x1b[1mError:\x1b[0m")
	fmt.Println("   ", e.msg)
	fmt.Println("\x1b[1mCaused by:\x1b[0m")
	fmt.Println("   ", e.Cause.Node.Kind(), "at", e.Cause.Where)
	fmt.Println()
}

func Error(n ast.MetaNode, format string, a ...any) *MetaError {
	return &MetaError{fmt.Sprintf(format, a...), n}
}
