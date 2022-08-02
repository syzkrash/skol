package codegen

import (
	"fmt"
	"io"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values/types"
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

func (g *ASTGenerator) internalGenerate(w io.Writer, n nodes.Node) error {
	for i := 0; i < int(g.indent); i++ {
		fmt.Fprint(w, indent)
	}
	switch n.Kind() {
	case nodes.NdInteger:
		fmt.Fprint(w, "Integer")
	case nodes.NdBoolean:
		fmt.Fprint(w, "Boolean")
	case nodes.NdFloat:
		fmt.Fprint(w, "Float")
	case nodes.NdString:
		fmt.Fprint(w, "String")
	case nodes.NdChar:
		fmt.Fprint(w, "Char")
	case nodes.NdVarDef:
		vdn := n.(*nodes.VarDefNode)
		fmt.Fprintf(w, "VarDef (%s):\n", vdn.Var)
		g.indent++
		g.internalGenerate(w, vdn.Value)
		g.indent--
	case nodes.NdFuncCall:
		fcn := n.(*nodes.FuncCallNode)
		fmt.Fprintf(w, "FuncCall (%s):\n", fcn.Func)
		g.indent++
		for _, an := range fcn.Args {
			g.internalGenerate(w, an)
		}
		g.indent--
	case nodes.NdFuncDef:
		fdn := n.(*nodes.FuncDefNode)
		fmt.Fprintf(w, "FuncDef (%s; ", fdn.Name)
		for _, a := range fdn.Args {
			fmt.Fprintf(w, "%s %s  ", a.Type, a.Name)
		}
		fmt.Fprintf(w, "; %s):\n", fdn.Ret)
		g.indent++
		for _, bn := range fdn.Body {
			g.internalGenerate(w, bn)
		}
		g.indent--
	case nodes.NdFuncExtern:
		fdn := n.(*nodes.FuncExternNode)
		fmt.Fprintf(w, "FuncExtern (%s; ", fdn.Name)
		for _, a := range fdn.Args {
			fmt.Fprintf(w, "%s %s  ", a.Type, a.Name)
		}
		fmt.Fprintf(w, "; %s)", fdn.Ret)
	case nodes.NdReturn:
		fmt.Fprintln(w, "Return:")
		g.indent++
		g.internalGenerate(w, n.(*nodes.ReturnNode).Value)
		g.indent--
	case nodes.NdIf:
		in := n.(*nodes.IfNode)
		fmt.Fprintln(w, "If:")
		g.indent++
		g.internalGenerate(w, in.Condition)
		g.indent++
		for _, bn := range in.IfBlock {
			g.internalGenerate(w, bn)
		}
		g.indent -= 2
	case nodes.NdWhile:
		wn := n.(*nodes.WhileNode)
		fmt.Fprintln(w, "While:")
		g.indent++
		g.internalGenerate(w, wn.Condition)
		g.indent++
		for _, bn := range wn.Body {
			g.internalGenerate(w, bn)
		}
		g.indent -= 2
	case nodes.NdStruct:
		sn := n.(*nodes.StructNode)
		fmt.Fprintf(w, "Struct (%s):\n", sn.Name)
		g.indent++
		for _, f := range sn.Type.(types.StructType).Fields {
			for i := 0; i < int(g.indent); i++ {
				fmt.Fprint(w, indent)
			}
			fmt.Fprintf(w, "%s %s\n", f.Type, f.Name)
		}
		g.indent--
	case nodes.NdNewStruct:
		nsn := n.(*nodes.NewStructNode)
		s := nsn.Type.(types.StructType)
		fmt.Fprintf(w, "NewStruct (%s):\n", s.Name)
		g.indent++
		for _, an := range nsn.Args {
			g.internalGenerate(w, an)
		}
		g.indent--
	case nodes.NdArray:
		an := n.(*nodes.ArrayNode)
		fmt.Fprintf(w, "Array (%s):\n", an.Type)
		g.indent++
		for _, en := range an.Elements {
			g.internalGenerate(w, en)
		}
		g.indent--
	default:
		if sel, ok := n.(nodes.Selector); ok {
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
				} else {
					fmt.Fprintf(w, "Access index %d\n", e.Idx)
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
