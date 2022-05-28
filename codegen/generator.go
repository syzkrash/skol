package codegen

import (
	"errors"
	"fmt"
	"io"

	"github.com/syzkrash/skol/parser"
)

type Generator interface {
	Generate(io.StringWriter) error
	CanRun() bool
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

func (g *ASTGenerator) Generate(output io.StringWriter) error {
	for {
		n, err := g.parser.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if _, err = output.WriteString(fmt.Sprint(n) + "\n"); err != nil {
			return err
		}
	}
}

func (*ASTGenerator) CanRun() bool {
	return false
}

func (*ASTGenerator) Run(string) error {
	panic("not supposed to call Run() here")
}