package codegen

import (
	"fmt"
	"io"

	"github.com/syzkrash/skol/parser"
)

type Generator interface {
	CanGenerate() bool
	Generate(io.StringWriter) error
	CanRun() bool
	Ext() string
	Run(string) error
}

type ASTGenerator struct {
	parser *parser.Parser
}

func NewAST(fn string, src io.RuneScanner) Generator {
	return &ASTGenerator{
		parser: parser.NewParser(fn, src),
	}
}

func (*ASTGenerator) CanGenerate() bool {
	return true
}

func (g *ASTGenerator) Generate(output io.StringWriter) error {
	n, err := g.parser.Next()
	if err != nil {
		return err
	}
	_, err = output.WriteString(fmt.Sprint(n) + "\n")
	return err
}

func (*ASTGenerator) CanRun() bool {
	return false
}

func (*ASTGenerator) Ext() string {
	return ".ast_dump.txt"
}

func (*ASTGenerator) Run(string) error {
	panic("not supposed to call Run() here")
}
