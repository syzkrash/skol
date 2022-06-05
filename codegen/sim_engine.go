package codegen

import (
	"errors"
	"io"

	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/nodes"
)

type SimEngine struct {
	p *parser.Parser
}

func NewSimEngine(fn string, src io.RuneScanner) Generator {
	return &SimEngine{
		p: parser.NewParser(fn, src),
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
	var n nodes.Node
	for {
		n, err = s.p.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return
		}
		err = s.p.Sim.Stmt(n)
		if err != nil {
			return
		}
	}
}
