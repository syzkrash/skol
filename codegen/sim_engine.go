package codegen

import (
	"io"

	"github.com/syzkrash/skol/typecheck"
)

type SimEngine struct {
	c *typecheck.Typechecker
}

func NewSimEngine(fn string, src io.RuneScanner) Generator {
	return &SimEngine{
		c: typecheck.NewTypechecker(src, fn, "sim"),
	}
}

func (*SimEngine) CanGenerate() bool {
	return false
}

func (*SimEngine) Generate(io.StringWriter) error {
	return io.EOF
}

func (*SimEngine) CanRun() bool {
	return true
}

func (*SimEngine) Ext() string {
	return ".sim.txt"
}

func (s *SimEngine) Run(string) (err error) {
	/*
		var n ast.MetaNode
		for {
			n, err = s.c.Next()
			if errors.Is(err, io.EOF) {
				return nil
			}
			if err != nil {
				return
			}
			err = s.c.Parser.Sim.Stmt(n)
			if err != nil {
				return
			}
		}
	*/
	return nil
}
