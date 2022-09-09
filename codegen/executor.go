package codegen

import (
	"io"
)

type Executor interface {
	Execute(io.Reader) error
}
