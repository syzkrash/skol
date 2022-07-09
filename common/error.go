package common

import (
	"fmt"

	"github.com/syzkrash/skol/parser/nodes"
)

type Printable interface {
	Print()
}

type MetaError struct {
	msg  string
	Node nodes.Node
}

func (e *MetaError) Error() string {
	return e.msg
}

func (e *MetaError) Print() {
	fmt.Println("\x1b[31m\x1b[1mError:\x1b[0m")
	fmt.Println("   ", e.msg)
	fmt.Println("\x1b[1mCaused by:\x1b[0m")
	fmt.Println("   ", e.Node.Kind(), "at", e.Node.Where())
	fmt.Println()
}

func Error(n nodes.Node, format string, a ...any) *MetaError {
	return &MetaError{fmt.Sprintf(format, a...), n}
}
