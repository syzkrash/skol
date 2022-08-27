package python

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/values/types"
)

const indent = "  "

func (p *pythonState) vt2pt(t types.Type) string {
	switch {
	case types.Bool.Equals(t):
		return "bool"
	case types.Int.Equals(t), types.Char.Equals(t):
		return "int"
	case types.Float.Equals(t):
		return "float"
	case types.String.Equals(t):
		return "string"
	case t.Prim() == types.PStruct:
		return t.(types.StructType).Name
	case t.Prim() == types.PArray:
		return "list"
	}
	return ""
}

func (p *pythonState) class(s ast.StructDefNode) (err error) {
	w := p.out.(io.Writer)

	// begin class definition
	//		class x:
	fmt.Fprintf(w, "class %s:\n", s.Name)

	// add __slots__
	//		__slots__ = ('a', 'b', 'c', ...,)
	fmt.Fprint(w, indent+"__slots__ = (")
	for _, f := range s.Fields {
		fmt.Fprintf(w, "\"%s\", ", f.Name)
	}
	fmt.Fprint(w, ")\n")

	// add type hints for fields
	//		a: int
	//		b: str
	//		...
	for _, f := range s.Fields {
		fmt.Fprintf(w, indent+"%s: %s\n", f.Name, p.vt2pt(f.Type))
	}

	// add constructor
	//	def __init__(self, a: int, b: str, ...,):
	fmt.Fprint(w, indent+"def __init__(")
	for _, f := range s.Fields {
		fmt.Fprintf(w, "%s: %s, ", f.Name, p.vt2pt(f.Type))
	}
	fmt.Fprint(w, "):\n")

	// constructor body
	//		self.a = a
	//		self.b = b
	//		...
	for _, f := range s.Fields {
		fmt.Fprintf(w, indent+indent+"self.%s = %s\n", f.Name, f.Name)
	}

	// add a final newline to ensure no indentation fuckery
	fmt.Fprint(w, "\n")
	return
}

func (p *pythonState) statement(n ast.Node) (err error) {
	defer p.out.WriteString("\n")

	for i := 0; i < int(p.ind); i++ {
		p.out.WriteString(indent)
	}

	switch n.Kind() {
	case ast.NVarSet:
		return p.varSet(n.(ast.VarSetNode))
	case ast.NFuncCall:
		return p.callOrExpr(n.(ast.FuncCallNode))
	case ast.NFuncDef:
		return p.funcDef(n.(ast.FuncDefNode))
	case ast.NReturn:
		return p.ret(n.(ast.ReturnNode))
	case ast.NIf:
		return p.ifn(n.(ast.IfNode))
	case ast.NWhile:
		return p.while(n.(ast.WhileNode))
	case ast.NStruct:
		return p.class(n.(ast.StructDefNode))
	case ast.NFuncExtern:
		// special case for externs
		return nil
	}
	return fmt.Errorf("%s node is not a statement", n.Kind())
}

func (p *pythonState) integer(n ast.IntNode) (err error) {
	_, err = p.out.WriteString(strconv.Itoa(int(n.Value)))
	return
}

func (p *pythonState) boolean(n ast.BoolNode) (err error) {
	if n.Value {
		_, err = p.out.WriteString("True")
	} else {
		_, err = p.out.WriteString("False")
	}
	return
}

func (p *pythonState) float(n ast.FloatNode) (err error) {
	_, err = p.out.WriteString(strconv.FormatFloat(float64(n.Value), 'g', 10, 64))
	return
}

func (p *pythonState) string(n ast.StringNode) (err error) {
	_, err = p.out.WriteString("\"" + strings.ReplaceAll(n.Value, "\"", "\\\"") + "\"")
	return
}

func (p *pythonState) char(n ast.CharNode) (err error) {
	_, err = p.out.WriteString(strconv.FormatInt(int64(n.Value), 10))
	return
}

func (p *pythonState) callOrExpr(n ast.FuncCallNode) (err error) {
	if n.Func == "import" {
		return p.impt(n.Args[0].Node.(ast.StringNode).Value)
	}
	if oper, ok := operators[n.Func]; ok {
		return p.expr(oper, n.Args[0].Node, n.Args[1].Node)
	}
	if sgen, ok := specialGenerators[n.Func]; ok {
		return sgen(p, n)
	}
	return p.funcCall(n)
}

func (p *pythonState) expr(oper string, lhs, rhs ast.Node) (err error) {
	_, err = p.out.WriteString("(")
	if err != nil {
		return
	}
	err = p.value(lhs)
	if err != nil {
		return
	}
	_, err = p.out.WriteString(oper)
	if err != nil {
		return
	}
	err = p.value(rhs)
	if err != nil {
		return
	}
	_, err = p.out.WriteString(")")
	return
}

func (p *pythonState) funcCall(n ast.FuncCallNode) (err error) {
	if nn, rename := renames[n.Func]; rename {
		_, err = p.out.WriteString(nn)
		if err != nil {
			return
		}
	} else {
		f, ok := p.parser.Scope.FindFunc(n.Func)
		if !ok {
			err = fmt.Errorf("unknown function: %s", n.Func)
			return
		}
		_, err = p.out.WriteString(f.Name)
		if err != nil {
			return
		}
	}
	_, err = p.out.WriteString("(")
	if err != nil {
		return
	}

	if len(n.Args) == 1 {
		err = p.value(n.Args[0].Node)
		if err != nil {
			return
		}
	} else if len(n.Args) > 1 {
		for _, a := range n.Args[:len(n.Args)-1] {
			err = p.value(a.Node)
			if err != nil {
				return
			}
			_, err = p.out.WriteString(", ")
			if err != nil {
				return
			}
		}
		err = p.value(n.Args[len(n.Args)-1].Node)
		if err != nil {
			return
		}
	}

	_, err = p.out.WriteString(")")
	return
}

func (p *pythonState) impt(mod string) (err error) {
	_, err = p.out.WriteString("import ")
	if err != nil {
		return
	}
	_, err = p.out.WriteString(mod)
	if err != nil {
		return
	}
	return
}

func (p *pythonState) instantiate(n ast.StructNode) (err error) {
	_, err = p.out.WriteString(n.Type.Name)
	if err != nil {
		return
	}
	_, err = p.out.WriteString("(")
	if err != nil {
		return
	}
	if len(n.Args) == 1 {
		err = p.value(n.Args[0].Node)
		if err != nil {
			return
		}
	} else {
		for _, a := range n.Args[:len(n.Args)-1] {
			err = p.value(a.Node)
			if err != nil {
				return
			}
			p.out.WriteString(", ")
		}
		err = p.value(n.Args[len(n.Args)-1].Node)
		if err != nil {
			return
		}
	}
	_, err = p.out.WriteString(")")
	return
}

func (p *pythonState) selector(s ast.Selector) (err error) {
	w := p.out.(io.Writer)
	path := s.Path()

	// preemtively write the first element since the parser and typechecker will
	// ensure it is valid by this point
	fmt.Fprint(w, path[0].Name)
	if len(path) == 1 {
		return
	}

	// iterate the rest of the path
	for _, e := range path[1:] {
		// the algorithm here works in the same way as in (*Parser).TypeOf()
		// see /parser/types.go, line 137
		if e.Cast != nil {
			// don't do anything since Python is dynamically typed anyway
		} else if e.Name != "" {
			fmt.Fprintf(w, ".%s", e.Name)
		} else if e.IdxS != nil {
			fmt.Fprintf(w, "[%s]", e.IdxS.(ast.SelectorNode).Child)
		} else {
			// this is not correct, since this will return a regular value, whereas
			// skol arrays return a result type when indexing
			// ¯\_(ツ)_/¯
			fmt.Fprintf(w, "[%d]", e.IdxC)
		}
	}
	return
}

func (p *pythonState) value(n ast.Node) error {
	switch n.Kind() {
	case ast.NInt:
		return p.integer(n.(ast.IntNode))
	case ast.NBool:
		return p.boolean(n.(ast.BoolNode))
	case ast.NFloat:
		return p.float(n.(ast.FloatNode))
	case ast.NString:
		return p.string(n.(ast.StringNode))
	case ast.NChar:
		return p.char(n.(ast.CharNode))
	case ast.NFuncCall:
		return p.callOrExpr(n.(ast.FuncCallNode))
	case ast.NStruct:
		return p.instantiate(n.(ast.StructNode))
	case ast.NArray:
		return p.list(n.(ast.ArrayNode))
	default:
		if s, ok := n.(ast.Selector); ok {
			return p.selector(s)
		}
	}
	return fmt.Errorf("%s node is not a value", n.Kind())
}

func (p *pythonState) stringOrVar(n ast.Node) error {
	switch n.Kind() {
	case ast.NString:
		return p.string(n.(ast.StringNode))
	}
	return fmt.Errorf("expected string or variable, got %s", n.Kind())
}

func (p *pythonState) integerOrVar(n ast.Node) error {
	switch n.Kind() {
	case ast.NInt:
		return p.integer(n.(ast.IntNode))
	}
	return fmt.Errorf("expected integer or variable, got %s", n.Kind())
}

func (p *pythonState) charOrVar(n ast.Node) error {
	switch n.Kind() {
	case ast.NChar:
		return p.char(n.(ast.CharNode))
	}
	return fmt.Errorf("expected char or variable, got %s", n.Kind())
}

func (p *pythonState) varSet(n ast.VarSetNode) (err error) {
	_, err = p.out.WriteString(n.Var)
	if err != nil {
		return
	}
	_, err = p.out.WriteString(" = ")
	if err != nil {
		return
	}
	err = p.value(n.Value.Node)
	return
}

func (p *pythonState) block(b ast.Block) (err error) {
	p.ind++
	for _, s := range b {
		err = p.statement(s.Node)
		if err != nil {
			return
		}
	}
	p.ind--
	return
}

func (p *pythonState) funcDef(n ast.FuncDefNode) (err error) {
	_, err = p.out.WriteString("def ")
	if err != nil {
		return
	}
	_, err = p.out.WriteString(n.Name)
	if err != nil {
		return
	}
	_, err = p.out.WriteString("(")
	if err != nil {
		return
	}

	i := len(n.Proto)
	for _, a := range n.Proto {
		_, err = p.out.WriteString(a.Name)
		if err != nil {
			return
		}
		if pyType := p.vt2pt(a.Type); pyType != "" {
			_, err = p.out.WriteString(": ")
			if err != nil {
				return
			}
			_, err = p.out.WriteString(pyType)
			if err != nil {
				return
			}
		}

		i--
		if i > 0 {
			_, err = p.out.WriteString(", ")
			if err != nil {
				return
			}
		}
	}

	_, err = p.out.WriteString(")")
	if err != nil {
		return
	}

	if pyType := p.vt2pt(n.Ret); pyType != "" {
		_, err = p.out.WriteString(" -> ")
		if err != nil {
			return
		}
		_, err = p.out.WriteString(pyType)
		if err != nil {
			return
		}
	}

	_, err = p.out.WriteString(":\n")
	if err != nil {
		return
	}

	return p.block(n.Body)
}

func (p *pythonState) ret(n ast.ReturnNode) (err error) {
	_, err = p.out.WriteString("return ")
	if err != nil {
		return
	}
	return p.value(n.Value.Node)
}

func (p *pythonState) ifn(n ast.IfNode) (err error) {
	_, err = p.out.WriteString("if ")
	if err != nil {
		return
	}
	err = p.value(n.Main.Cond.Node)
	if err != nil {
		return
	}
	_, err = p.out.WriteString(":\n")
	if err != nil {
		return
	}
	err = p.block(n.Main.Block)
	if err != nil {
		return
	}

	for _, elif := range n.Other {
		for i := 0; i < int(p.ind); i++ {
			p.out.WriteString(indent)
		}

		_, err = p.out.WriteString("elif ")
		if err != nil {
			return
		}
		err = p.value(elif.Cond.Node)
		if err != nil {
			return
		}
		_, err = p.out.WriteString(":\n")
		if err != nil {
			return
		}
		err = p.block(elif.Block)
		if err != nil {
			return
		}
	}

	if len(n.Else) > 0 {
		for i := 0; i < int(p.ind); i++ {
			p.out.WriteString(indent)
		}

		_, err = p.out.WriteString("else:\n")
		if err != nil {
			return
		}
		err = p.block(n.Else)
		if err != nil {
			return
		}
	}

	return
}

func (p *pythonState) while(n ast.WhileNode) (err error) {
	_, err = p.out.WriteString("while ")
	if err != nil {
		return
	}
	err = p.value(n.Cond.Node)
	if err != nil {
		return
	}
	_, err = p.out.WriteString(":\n")
	if err != nil {
		return
	}
	return p.block(n.Block)
}

func (p *pythonState) list(n ast.ArrayNode) (err error) {
	w := p.out.(io.Writer)
	fmt.Fprint(w, "[")
	if len(n.Elems) == 1 {
		p.value(n.Elems[0].Node)
	} else if len(n.Elems) > 1 {
		for _, e := range n.Elems[:len(n.Elems)-1] {
			p.value(e.Node)
			p.out.WriteString(", ")
		}
		p.value(n.Elems[len(n.Elems)-1].Node)
	}
	fmt.Fprint(w, "]")
	return
}
