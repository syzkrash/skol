package python

import (
	"fmt"
	"strconv"

	"github.com/syzkrash/skol/parser"
)

const indent = "  "

var ops = map[string]string{
	"add": "+",
	"sub": "-",
	"mul": "*",
	"pow": "**",
	"div": "/",
	"mod": "%",

	"eq":  "==",
	"neq": "!=",
	"gt":  ">",
	"lt":  "<",
	"geq": ">=",
	"leq": "<=",
}

func (p *pythonState) vt2pt(t parser.ValueType) string {
	switch t {
	case parser.VtBool:
		return "bool"
	case parser.VtChar:
		return "int"
	case parser.VtFloat:
		return "float"
	case parser.VtInteger:
		return "int"
	case parser.VtString:
		return "str"
	}
	return ""
}

func (p *pythonState) statement(n parser.Node) (err error) {
	defer p.out.WriteString("\n")

	for i := 0; i < int(p.ind); i++ {
		p.out.WriteString(indent)
	}

	switch n.Kind() {
	case parser.NdVarDef:
		return p.varDef(n.(*parser.VarDefNode))
	case parser.NdFuncCall:
		return p.callOrExpr(n.(*parser.FuncCallNode))
	case parser.NdFuncDef:
		return p.funcDef(n.(*parser.FuncDefNode))
	case parser.NdReturn:
		return p.ret(n.(*parser.ReturnNode))
	case parser.NdIf:
		return p.ifn(n.(*parser.IfNode))
	case parser.NdWhile:
		return p.while(n.(*parser.WhileNode))

	case parser.NdFuncExtern:
		// special case for externs
		return nil
	}
	return fmt.Errorf("%s node is not a statement", n.Kind())
}

func (p *pythonState) integer(n *parser.IntegerNode) (err error) {
	_, err = p.out.WriteString(strconv.Itoa(int(n.Int)))
	return
}

func (p *pythonState) float(n *parser.FloatNode) (err error) {
	_, err = p.out.WriteString(strconv.FormatFloat(float64(n.Float), 'g', 10, 64))
	return
}

func (p *pythonState) string(n *parser.StringNode) (err error) {
	_, err = p.out.WriteString("\"" + n.Str + "\"")
	return
}

func (p *pythonState) char(n *parser.CharNode) (err error) {
	_, err = p.out.WriteString(strconv.FormatInt(int64(n.Char), 10))
	return
}

func (p *pythonState) varRef(n *parser.VarRefNode) (err error) {
	_, err = p.out.WriteString(n.Var)
	return
}

func (p *pythonState) callOrExpr(n *parser.FuncCallNode) (err error) {
	if oper, ok := ops[n.Func]; ok {
		return p.expr(oper, n.Args[0], n.Args[1])
	}
	return p.funcCall(n)
}

func (p *pythonState) expr(oper string, lhs, rhs parser.Node) (err error) {
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

func (p *pythonState) funcCall(n *parser.FuncCallNode) (err error) {
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

func (p *pythonState) value(n parser.Node) error {
	switch n.Kind() {
	case parser.NdInteger:
		return p.integer(n.(*parser.IntegerNode))
	case parser.NdFloat:
		return p.float(n.(*parser.FloatNode))
	case parser.NdString:
		return p.string(n.(*parser.StringNode))
	case parser.NdChar:
		return p.char(n.(*parser.CharNode))
	case parser.NdVarRef:
		return p.varRef(n.(*parser.VarRefNode))
	case parser.NdFuncCall:
		return p.callOrExpr(n.(*parser.FuncCallNode))
	}
	return fmt.Errorf("%s node is not a value", n.Kind())
}

func (p *pythonState) varDef(n *parser.VarDefNode) (err error) {
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

func (p *pythonState) block(n []parser.Node) (err error) {
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

func (p *pythonState) funcDef(n *parser.FuncDefNode) (err error) {
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

func (p *pythonState) ret(n *parser.ReturnNode) (err error) {
	_, err = p.out.WriteString("return ")
	if err != nil {
		return
	}
	return p.value(n.Value)
}

func (p *pythonState) ifn(n *parser.IfNode) (err error) {
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

func (p *pythonState) while(n *parser.WhileNode) (err error) {
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
