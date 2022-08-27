package codegen

import (
	"fmt"
	"io"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/typecheck"
)

const indent = "  "

type ASTGenerator struct {
	checker *typecheck.Typechecker
	indent  uint
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
	return g.internalGenerate(output.(io.Writer), n)
}

func (g *ASTGenerator) internalGenerate(w io.Writer, mn ast.MetaNode) error {
	for i := 0; i < int(g.indent); i++ {
		fmt.Fprint(w, indent)
	}
	n := mn.Node
	switch n.Kind() {
	case ast.NInt:
		fmt.Fprint(w, "Int")
	case ast.NBool:
		fmt.Fprint(w, "Bool")
	case ast.NFloat:
		fmt.Fprint(w, "Float")
	case ast.NString:
		fmt.Fprint(w, "String")
	case ast.NChar:
		fmt.Fprint(w, "Char")
	case ast.NVarSet:
		vdn := n.(ast.VarSetNode)
		fmt.Fprintf(w, "VarSet (%s):\n", vdn.Var)
		g.indent++
		g.internalGenerate(w, vdn.Value)
		g.indent--
	case ast.NFuncCall:
		fcn := n.(ast.FuncCallNode)
		fmt.Fprintf(w, "FuncCall (%s):\n", fcn.Func)
		g.indent++
		for _, an := range fcn.Args {
			g.internalGenerate(w, an)
		}
		g.indent--
	case ast.NFuncDef:
		fdn := n.(ast.FuncDefNode)
		fmt.Fprintf(w, "FuncDef (%s)", fdn.Name)
		g.indent++
		for _, bn := range fdn.Body {
			g.internalGenerate(w, bn)
		}
		g.indent--
	case ast.NFuncExtern:
		fdn := n.(ast.FuncExternNode)
		fmt.Fprintf(w, "FuncExtern (%s)", fdn.Name)
	case ast.NReturn:
		fmt.Fprintln(w, "Return:")
		g.indent++
		g.internalGenerate(w, n.(ast.ReturnNode).Value)
		g.indent--
	case ast.NIf:
		in := n.(ast.IfNode)
		fmt.Fprintln(w, "If:")
		g.indent++
		g.internalGenerate(w, in.Main.Cond)
		g.indent++
		for _, bn := range in.Main.Block {
			g.internalGenerate(w, bn)
		}
		g.indent -= 2
	case ast.NWhile:
		wn := n.(ast.WhileNode)
		fmt.Fprintln(w, "While:")
		g.indent++
		g.internalGenerate(w, wn.Cond)
		g.indent++
		for _, bn := range wn.Block {
			g.internalGenerate(w, bn)
		}
		g.indent -= 2
	case ast.NStructDef:
		sn := n.(ast.StructDefNode)
		fmt.Fprintf(w, "StructDef (%s):\n", sn.Name)
		g.indent++
		for _, f := range sn.Fields {
			for i := 0; i < int(g.indent); i++ {
				fmt.Fprint(w, indent)
			}
			fmt.Fprintf(w, "%s %s\n", f.Type, f.Name)
		}
		g.indent--
	case ast.NStruct:
		nsn := n.(ast.StructNode)
		s := nsn.Type
		fmt.Fprintf(w, "Struct (%s):\n", s.Name)
		g.indent++
		for _, an := range nsn.Args {
			g.internalGenerate(w, an)
		}
		g.indent--
	case ast.NArray:
		an := n.(ast.ArrayNode)
		fmt.Fprintf(w, "Array (%s):\n", an.Type)
		g.indent++
		for _, en := range an.Elems {
			g.internalGenerate(w, en)
		}
		g.indent--
	default:
		if sel, ok := n.(ast.Selector); ok {
			fmt.Fprintln(w, "Selector:")
			g.indent++
			for _, e := range sel.Path() {
				for i := 0; i < int(g.indent); i++ {
					fmt.Fprint(w, indent)
				}
				if e.Cast != nil {
					fmt.Fprintf(w, "Cast to %s\n", e.Cast)
				} else if e.Name != "" {
					fmt.Fprintf(w, "Access field '%s'\n", e.Name)
				} else if e.IdxS != nil {
					fmt.Fprintf(w, "Access index %+v\n", e.IdxS)
				} else {
					fmt.Fprintf(w, "Access index %d\n", e.IdxC)
				}
			}
			g.indent--
		}
	}
	fmt.Fprint(w, "\n")
	return nil
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
