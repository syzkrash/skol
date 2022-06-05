package sim

import (
	"errors"
	"io"

	"github.com/syzkrash/skol/codegen"
	"github.com/syzkrash/skol/parser"
)

type SimEngine struct {
	s *Simulator
	p *parser.Parser
}

func NewEngine(fn string, src io.RuneScanner) codegen.Generator {
	return &SimEngine{
		s: NewSimulator(),
		p: parser.NewParser(fn, src),
	}
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
	var n parser.Node
	for {
		n, err = s.p.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return
		}
		err = s.s.Stmt(n)
		if err != nil {
			return
		}
	}
}
