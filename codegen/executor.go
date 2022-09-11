package codegen

import (
	"io"
)

type Executor interface{}

type EphemeralExecutor interface {
	Executor
	Execute(io.Reader) error
}

type FilenameExecutor interface {
	Executor
	Execute(string) error
}
