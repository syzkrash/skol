package codegen

import (
	"fmt"
	"io"

	"github.com/syzkrash/skol/typecheck"
)

type Generator interface {
	CanGenerate() bool
	Generate(io.StringWriter) error
	CanRun() bool
	Ext() string
	Run(string) error
}

type ASTGenerator struct {
	checker *typecheck.Typechecker
}

func NewAST(fn string, src io.RuneScanner) Generator {
	return &ASTGenerator{
		checker: typecheck.NewTypechecker(src, fn, "ast"),
	}
}

func (*ASTGenerator) CanGenerate() bool {
	return true
}

func (g *ASTGenerator) Generate(output io.StringWriter) error {
	n, err := g.checker.Next()
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
