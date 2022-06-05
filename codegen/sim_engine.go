package codegen

import (
	"errors"
	"io"

	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/sim"
)

type SimEngine struct {
	s *sim.Simulator
	p *parser.Parser
}

func NewEngine(fn string, src io.RuneScanner) Generator {
	return &SimEngine{
		s: sim.NewSimulator(),
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
	var n nodes.Node
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
