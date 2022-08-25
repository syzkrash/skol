package python

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/syzkrash/skol/parser/nodes"
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

func (p *pythonState) class(s *nodes.StructNode) (err error) {
	str := s.Type.(types.StructType)
	w := p.out.(io.Writer)

	// begin class definition
	//		class x:
	fmt.Fprintf(w, "class %s:\n", s.Name)

	// add __slots__
	//		__slots__ = ('a', 'b', 'c', ...,)
	fmt.Fprint(w, indent+"__slots__ = (")
	for _, f := range str.Fields {
		fmt.Fprintf(w, "\"%s\", ", f.Name)
	}
	fmt.Fprint(w, ")\n")

	// add type hints for fields
	//		a: int
	//		b: str
	//		...
	for _, f := range str.Fields {
		fmt.Fprintf(w, indent+"%s: %s\n", f.Name, p.vt2pt(f.Type))
	}

	// add constructor
	//	def __init__(self, a: int, b: str, ...,):
	fmt.Fprint(w, indent+"def __init__(")
	for _, f := range str.Fields {
		fmt.Fprintf(w, "%s: %s, ", f.Name, p.vt2pt(f.Type))
	}
	fmt.Fprint(w, "):\n")

	// constructor body
	//		self.a = a
	//		self.b = b
	//		...
	for _, f := range str.Fields {
		fmt.Fprintf(w, indent+indent+"self.%s = %s\n", f.Name, f.Name)
	}

	// add a final newline to ensure no indentation fuckery
	fmt.Fprint(w, "\n")
	return
}

func (p *pythonState) statement(n nodes.Node) (err error) {
	defer p.out.WriteString("\n")

	for i := 0; i < int(p.ind); i++ {
		p.out.WriteString(indent)
	}

	switch n.Kind() {
	case nodes.NdVarDef:
		return p.varDef(n.(*nodes.VarDefNode))
	case nodes.NdFuncCall:
		return p.callOrExpr(n.(*nodes.FuncCallNode))
	case nodes.NdFuncDef:
		return p.funcDef(n.(*nodes.FuncDefNode))
	case nodes.NdReturn:
		return p.ret(n.(*nodes.ReturnNode))
	case nodes.NdIf:
		return p.ifn(n.(*nodes.IfNode))
	case nodes.NdWhile:
		return p.while(n.(*nodes.WhileNode))
	case nodes.NdStruct:
		return p.class(n.(*nodes.StructNode))
	case nodes.NdFuncExtern:
		// special case for externs
		return nil
	}
	return fmt.Errorf("%s node is not a statement", n.Kind())
}

func (p *pythonState) integer(n *nodes.IntegerNode) (err error) {
	_, err = p.out.WriteString(strconv.Itoa(int(n.Int)))
	return
}

func (p *pythonState) boolean(n *nodes.BooleanNode) (err error) {
	if n.Bool {
		_, err = p.out.WriteString("True")
	} else {
		_, err = p.out.WriteString("False")
	}
	return
}

func (p *pythonState) float(n *nodes.FloatNode) (err error) {
	_, err = p.out.WriteString(strconv.FormatFloat(float64(n.Float), 'g', 10, 64))
	return
}

func (p *pythonState) string(n *nodes.StringNode) (err error) {
	_, err = p.out.WriteString("\"" + strings.ReplaceAll(n.Str, "\"", "\\\"") + "\"")
	return
}

func (p *pythonState) char(n *nodes.CharNode) (err error) {
	_, err = p.out.WriteString(strconv.FormatInt(int64(n.Char), 10))
	return
}

func (p *pythonState) callOrExpr(n *nodes.FuncCallNode) (err error) {
	if n.Func == "import" {
		return p.impt(n.Args[0].(*nodes.StringNode).Str)
	}
	if oper, ok := operators[n.Func]; ok {
		return p.expr(oper, n.Args[0], n.Args[1])
	}
	if sgen, ok := specialGenerators[n.Func]; ok {
		return sgen(p, n)
	}
	return p.funcCall(n)
}

func (p *pythonState) expr(oper string, lhs, rhs nodes.Node) (err error) {
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

func (p *pythonState) funcCall(n *nodes.FuncCallNode) (err error) {
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
		err = p.value(n.Args[0])
		if err != nil {
			return
		}
	} else if len(n.Args) > 1 {
		for _, a := range n.Args[:len(n.Args)-1] {
			err = p.value(a)
			if err != nil {
				return
			}
			_, err = p.out.WriteString(", ")
			if err != nil {
				return
			}
		}
		err = p.value(n.Args[len(n.Args)-1])
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

func (p *pythonState) instantiate(n *nodes.NewStructNode) (err error) {
	_, err = p.out.WriteString(n.Type.(types.StructType).Name)
	if err != nil {
		return
	}
	_, err = p.out.WriteString("(")
	if err != nil {
		return
	}
	if len(n.Args) == 1 {
		err = p.value(n.Args[0])
		if err != nil {
			return
		}
	} else {
		for _, a := range n.Args[:len(n.Args)-1] {
			err = p.value(a)
			if err != nil {
				return
			}
			p.out.WriteString(", ")
		}
		err = p.value(n.Args[len(n.Args)-1])
		if err != nil {
			return
		}
	}
	_, err = p.out.WriteString(")")
	return
}

func (p *pythonState) selector(s nodes.Selector) (err error) {
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
		} else {
			// this is not correct, since this will return a regular value, whereas
			// skol arrays return a result type when indexing
			// ¯\_(ツ)_/¯
			if e.Idx.Kind() == nodes.NdSelector {
				fmt.Fprintf(w, "[%s]", e.Idx.(*nodes.SelectorNode).Child)
			} else {
				fmt.Fprintf(w, "[%d]", e.Idx.(*nodes.IntegerNode).Int)
			}
		}
	}
	return
}

func (p *pythonState) value(n nodes.Node) error {
	switch n.Kind() {
	case nodes.NdInteger:
		return p.integer(n.(*nodes.IntegerNode))
	case nodes.NdBoolean:
		return p.boolean(n.(*nodes.BooleanNode))
	case nodes.NdFloat:
		return p.float(n.(*nodes.FloatNode))
	case nodes.NdString:
		return p.string(n.(*nodes.StringNode))
	case nodes.NdChar:
		return p.char(n.(*nodes.CharNode))
	case nodes.NdFuncCall:
		return p.callOrExpr(n.(*nodes.FuncCallNode))
	case nodes.NdNewStruct:
		return p.instantiate(n.(*nodes.NewStructNode))
	case nodes.NdArray:
		return p.list(n.(*nodes.ArrayNode))
	default:
		if s, ok := n.(nodes.Selector); ok {
			return p.selector(s)
		}
	}
	return fmt.Errorf("%s node is not a value", n.Kind())
}

func (p *pythonState) stringOrVar(n nodes.Node) error {
	switch n.Kind() {
	case nodes.NdString:
		return p.string(n.(*nodes.StringNode))
	}
	return fmt.Errorf("expected string or variable, got %s", n.Kind())
}

func (p *pythonState) integerOrVar(n nodes.Node) error {
	switch n.Kind() {
	case nodes.NdInteger:
		return p.integer(n.(*nodes.IntegerNode))
	}
	return fmt.Errorf("expected integer or variable, got %s", n.Kind())
}

func (p *pythonState) charOrVar(n nodes.Node) error {
	switch n.Kind() {
	case nodes.NdChar:
		return p.char(n.(*nodes.CharNode))
	}
	return fmt.Errorf("expected char or variable, got %s", n.Kind())
}

func (p *pythonState) varDef(n *nodes.VarDefNode) (err error) {
	_, err = p.out.WriteString(n.Var)
	if err != nil {
		return
	}
	if pyType := p.vt2pt(n.VarType); pyType != "" {
		_, err = p.out.WriteString(": ")
		if err != nil {
			return
		}
		_, err = p.out.WriteString(pyType)
		if err != nil {
			return
		}
	}
	_, err = p.out.WriteString(" = ")
	if err != nil {
		return
	}
	err = p.value(n.Value)
	return
}

func (p *pythonState) block(n []nodes.Node) (err error) {
	p.ind++
	for _, s := range n {
		err = p.statement(s)
		if err != nil {
			return
		}
	}
	p.ind--
	return
}

func (p *pythonState) funcDef(n *nodes.FuncDefNode) (err error) {
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

	i := len(n.Args)
	for _, a := range n.Args {
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

func (p *pythonState) ret(n *nodes.ReturnNode) (err error) {
	_, err = p.out.WriteString("return ")
	if err != nil {
		return
	}
	return p.value(n.Value)
}

func (p *pythonState) ifn(n *nodes.IfNode) (err error) {
	_, err = p.out.WriteString("if ")
	if err != nil {
		return
	}
	err = p.value(n.Condition)
	if err != nil {
		return
	}
	_, err = p.out.WriteString(":\n")
	if err != nil {
		return
	}
	err = p.block(n.IfBlock)
	if err != nil {
		return
	}

	for _, elif := range n.ElseIfNodes {
		for i := 0; i < int(p.ind); i++ {
			p.out.WriteString(indent)
		}

		_, err = p.out.WriteString("elif ")
		if err != nil {
			return
		}
		err = p.value(elif.Condition)
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

	if len(n.ElseBlock) > 0 {
		for i := 0; i < int(p.ind); i++ {
			p.out.WriteString(indent)
		}

		_, err = p.out.WriteString("else:\n")
		if err != nil {
			return
		}
		err = p.block(n.ElseBlock)
		if err != nil {
			return
		}
	}

	return
}

func (p *pythonState) while(n *nodes.WhileNode) (err error) {
	_, err = p.out.WriteString("while ")
	if err != nil {
		return
	}
	err = p.value(n.Condition)
	if err != nil {
		return
	}
	_, err = p.out.WriteString(":\n")
	if err != nil {
		return
	}
	return p.block(n.Body)
}

func (p *pythonState) list(n *nodes.ArrayNode) (err error) {
	w := p.out.(io.Writer)
	fmt.Fprint(w, "[")
	if len(n.Elements) == 1 {
		p.value(n.Elements[0])
	} else if len(n.Elements) > 1 {
		for _, e := range n.Elements[:len(n.Elements)-1] {
			p.value(e)
			p.out.WriteString(", ")
		}
		p.value(n.Elements[len(n.Elements)-1])
	}
	fmt.Fprint(w, "]")
	return
}
