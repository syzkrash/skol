package codegen

import (
	"io"
)

type Generator interface {
	CanGenerate() bool
	Generate(io.StringWriter) error
	CanRun() bool
	Ext() string
	Run(string) error
}
