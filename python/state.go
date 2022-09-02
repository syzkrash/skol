package python

import (
	"io"
	"os"
	"os/exec"

	"github.com/syzkrash/skol/codegen"
	"github.com/syzkrash/skol/parser"
)

type pythonState struct {
	parser *parser.Parser
	ind    uint
	out    io.StringWriter
}

func NewPython(fn string, input io.RuneScanner) codegen.Generator {
	gen := &pythonState{
		parser: parser.NewParser(fn, input, "python"),
		ind:    0,
	}
	gen.initEnv()
	return gen
}

func (p *pythonState) CanGenerate() bool {
	return true
}

func (p *pythonState) Generate(output io.StringWriter) error {
	/*
		p.out = output
		n, err := p.parser.Next()
		if err != nil {
			return err
		}
		return p.statement(n.Node)
	*/
	return nil
}

func (*pythonState) Ext() string {
	return ".py"
}

func (*pythonState) CanRun() bool {
	return true
}

func (*pythonState) Run(fn string) error {
	cmd := exec.Command("py", fn)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
