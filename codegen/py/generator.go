package py

import (
	"fmt"
	"io"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/codegen"
	"github.com/syzkrash/skol/parser/values/types"
)

const indent = "  "

type generator struct {
	out    io.Writer
	in     ast.AST
	indent int
}

var _ codegen.Generator = &generator{}
var _ codegen.ASTGenerator = &generator{}

func (g *generator) Output(w io.Writer) {
	g.out = w
}

func (g *generator) Input(t ast.AST) {
	g.in = t
}

func (g *generator) Generate() error {
	for _, f := range g.in.Funcs {
		g.writeFunc_(f)
	}
	return nil
}

func (g *generator) writeIndent() error {
	for i := 0; i < g.indent; i++ {
		_, err := g.out.Write([]byte(indent))
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) write(format string, args ...any) error {
	_, err := fmt.Fprintf(g.out, format, args...)
	return err
}

func (g *generator) pyType(st types.Type) string {
	t := ""
	switch {
	case types.Bool.Equals(st):
		t = "bool"
	case types.Char.Equals(st):
		t = "int"
	case types.Int.Equals(st):
		t = "int"
	case types.Float.Equals(st):
		t = "float"
	case types.String.Equals(st):
		t = "string"
	case st.Prim() == types.PArray:
		t = "list"
	case st.Prim() == types.PStruct:
		t = st.(types.StructType).Name
	default:
		panic("pyType() call got unexpected type: " + st.String())
	}
	return t
}

func (g *generator) writeArg(a ast.FuncProtoArg) error {
	return g.write("%s: %s,", a.Name, g.pyType(a.Type))
}

func (g *generator) writeBlock(b ast.Block) error {
	g.indent++
	for _, n := range b {
		if err := g.writeStmt(n); err != nil {
			return err
		}
	}
	g.indent--
	return nil
}

func (g *generator) writeStmt(mn ast.MetaNode) error {
	if err := g.writeIndent(); err != nil {
		return err
	}
	n := mn.Node
	switch n.Kind() {
	case ast.NIf:
		return g.writeIf(n.(ast.IfNode))
	case ast.NWhile:
		return g.writeWhile(n.(ast.WhileNode))
	case ast.NReturn:
		return g.writeReturn(n.(ast.ReturnNode))
	case ast.NVarSet:
		return g.writeVarSet(n.(ast.VarSetNode))
	case ast.NVarDef:
		return g.writeVarDef(n.(ast.VarDefNode))
	case ast.NVarSetTyped:
		return g.writeVarSetTyped(n.(ast.VarSetTypedNode))
	case ast.NFuncDef:
		return g.writeFunc(n.(ast.FuncDefNode))
	case ast.NFuncExtern:
		return nil
	case ast.NStructDef:
		return g.writeClass(n.(ast.StructDefNode))
	case ast.NFuncCall:
		return g.writeCall(n.(ast.FuncCallNode), true)
	default:
		panic("writeStmt() unexpected argument: " + n.Kind().String())
	}
}

func (g *generator) writeFunc(n ast.FuncDefNode) error {
	g.write("def %s(", n.Name)
	for _, a := range n.Proto {
		g.writeArg(a)
	}
	g.write("):\n")
	return g.writeBlock(n.Body)
}

func (g *generator) writeFunc_(f ast.Func) error {
	return g.writeFunc(ast.FuncDefNode{
		Name:  f.Name,
		Proto: f.Args,
		Ret:   f.Ret,
		Body:  f.Body,
	})
}

func (g *generator) writeIf(n ast.IfNode) error {
	g.write("if ")
	g.writeValue(n.Main.Cond)
	g.write(":\n")
	g.writeBlock(n.Main.Block)
	for _, b := range n.Other {
		g.writeIndent()
		g.write("elif ")
		g.writeValue(b.Cond)
		g.write(":\n")
		g.writeBlock(b.Block)
	}
	if len(n.Else) == 0 {
		return nil
	}
	g.writeIndent()
	g.write("else:\n")
	return g.writeBlock(n.Else)
}

func (g *generator) writeWhile(n ast.WhileNode) error {
	g.write("while ")
	g.writeValue(n.Cond)
	g.write(":\n")
	return g.writeBlock(n.Block)
}

func (g *generator) writeReturn(n ast.ReturnNode) error {
	g.write("return ")
	g.writeValue(n.Value)
	return g.write("\n")
}

func (g *generator) writeVarSet(n ast.VarSetNode) error {
	g.write("%s = ", n.Var)
	g.writeValue(n.Value)
	return g.write("\n")
}

func (g *generator) writeVarDef(n ast.VarDefNode) error {
	return g.write("%s: %s\n", n.Var, g.pyType(n.Type))
}

func (g *generator) writeVarSetTyped(n ast.VarSetTypedNode) error {
	g.write("%s: %s =", n.Var, g.pyType(n.Type))
	g.writeValue(n.Value)
	return g.write("\n")
}

func (g *generator) writeClass(n ast.StructDefNode) error {
	g.write("class %s:\n", n.Name)
	g.indent++
	g.writeIndent()
	g.write("__slots__ = (")
	for _, f := range n.Fields {
		g.write(`"%s", `, f.Name)
	}
	g.write(")\n")
	for _, f := range n.Fields {
		g.writeIndent()
		g.write("%s: %s\n", f.Name, g.pyType(f.Type))
	}
	g.writeIndent()
	g.write("def __init__(self, ")
	for _, f := range n.Fields {
		g.write("%s: %s, ", f.Name, f.Type)
	}
	g.write("):\n")
	g.indent++
	for _, f := range n.Fields {
		g.writeIndent()
		g.write("self.%s = %s\n", f.Name, f.Name)
	}
	g.indent--
	g.indent--
	return nil
}

func (g *generator) writeCall(n ast.FuncCallNode, stmt bool) error {
	g.write("%s(", n.Func)
	for _, a := range n.Args {
		g.writeValue(a)
		g.write(",")
	}
	if stmt {
		return g.write(")\n")
	}
	return g.write(")")
}

func (g *generator) writeValue(mn ast.MetaNode) error {
	n := mn.Node
	switch n.Kind() {
	case ast.NBool:
		if n.(ast.BoolNode).Value {
			return g.write("True")
		} else {
			return g.write("False")
		}
	case ast.NChar:
		return g.write("%d", n.(ast.CharNode).Value)
	case ast.NInt:
		return g.write("%d", n.(ast.IntNode).Value)
	case ast.NFloat:
		return g.write("%f", n.(ast.FloatNode).Value)
	case ast.NString:
		return g.write(`"%s"`, n.(ast.StringNode).Value)
	case ast.NStruct:
		return g.writeInstance(n.(ast.StructNode))
	case ast.NArray:
		return g.writeArray(n.(ast.ArrayNode))
	case ast.NFuncCall:
		return g.writeCall(n.(ast.FuncCallNode), false)
	default:
		if sel, ok := n.(ast.Selector); ok {
			return g.writeSelector(sel)
		}
		panic("writeValue() unexpected argument: " + n.Kind().String())
	}
}

func (g *generator) writeInstance(n ast.StructNode) error {
	g.write("%s(", n.Type.Name)
	for _, f := range n.Args {
		g.writeValue(f)
		g.write(", ")
	}
	return g.write(")")
}

func (g *generator) writeArray(n ast.ArrayNode) error {
	g.write("[")
	for _, v := range n.Elems {
		g.writeValue(v)
	}
	return g.write("]")
}

func (g *generator) writeSelector(sel ast.Selector) error {
	p := sel.Path()
	g.write("%s", p[0].Name)
	if len(p) == 1 {
		return nil
	}
	for _, e := range p[1:] {
		if e.Name != "" {
			g.write(".%s", e.Name)
		} else if e.Cast != nil {
			continue
		} else if e.IdxS != nil {
			g.write("[")
			g.writeSelector(e.IdxS)
			g.write("]")
		} else {
			g.write("[%d]", e.IdxC)
		}
	}
	return nil
}
