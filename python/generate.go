package python

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

const indent = "  "

func (p *pythonState) vt2pt(t values.ValueType) string {
	switch t {
	case values.VtBool:
		return "bool"
	case values.VtChar:
		return "int"
	case values.VtFloat:
		return "float"
	case values.VtInteger:
		return "int"
	case values.VtString:
		return "str"
	}
	return ""
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

func (p *pythonState) varRef(n *nodes.VarRefNode) (err error) {
	_, err = p.out.WriteString(n.Var)
	return
}

func (p *pythonState) callOrExpr(n *nodes.FuncCallNode) (err error) {
	if n.Func == "import" {
		return p.impt(n.Args[0].(*nodes.StringNode).Str)
	}
	if oper, ok := operators[n.Func]; ok {
		return p.expr(oper, n.Args[0], n.Args[1])
	}
	if new_name, ok := renames[n.Func]; ok {
		// copy and change the name
		b := *n
		b.Func = new_name
		return p.funcCall(&b)
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
	_, err = p.out.WriteString(n.Func)
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
	case nodes.NdVarRef:
		return p.varRef(n.(*nodes.VarRefNode))
	case nodes.NdFuncCall:
		return p.callOrExpr(n.(*nodes.FuncCallNode))
	}
	return fmt.Errorf("%s node is not a value", n.Kind())
}

func (p *pythonState) stringOrVar(n nodes.Node) error {
	switch n.Kind() {
	case nodes.NdString:
		return p.string(n.(*nodes.StringNode))
	case nodes.NdVarRef:
		return p.varRef(n.(*nodes.VarRefNode))
	}
	return fmt.Errorf("expected string or variable, got %s", n.Kind())
}

func (p *pythonState) integerOrVar(n nodes.Node) error {
	switch n.Kind() {
	case nodes.NdInteger:
		return p.integer(n.(*nodes.IntegerNode))
	case nodes.NdVarRef:
		return p.varRef(n.(*nodes.VarRefNode))
	}
	return fmt.Errorf("expected integer or variable, got %s", n.Kind())
}

func (p *pythonState) charOrVar(n nodes.Node) error {
	switch n.Kind() {
	case nodes.NdChar:
		return p.char(n.(*nodes.CharNode))
	case nodes.NdVarRef:
		return p.varRef(n.(*nodes.VarRefNode))
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
	for a, t := range n.Args {
		_, err = p.out.WriteString(a)
		if err != nil {
			return
		}
		if pyType := p.vt2pt(t); pyType != "" {
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
