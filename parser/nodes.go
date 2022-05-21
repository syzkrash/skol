package parser

import "fmt"

type NodeKind uint8

const (
	NdInteger NodeKind = iota
	NdFloat
	NdString
	NdChar
	NdVarRef
	NdVarDef
	NdFuncCall
	NdFuncDef
	NdReturn
)

var nodeKinds = []string{
	"Integer",
	"Float",
	"String",
	"Char",
	"VarRef",
	"VarDef",
	"FuncCall",
	"FuncDef",
	"Return",
}

func (k NodeKind) String() string {
	return nodeKinds[k]
}

type Node interface {
	Kind() NodeKind
}

type IntegerNode struct {
	Int int32
}

func (*IntegerNode) Kind() NodeKind {
	return NdInteger
}

func (n *IntegerNode) String() string {
	return fmt.Sprintf("Integer{%d}", n.Int)
}

type FloatNode struct {
	Float float32
}

func (*FloatNode) Kind() NodeKind {
	return NdFloat
}

func (n *FloatNode) String() string {
	return fmt.Sprintf("Float{%f}", n.Float)
}

type StringNode struct {
	Str string
}

func (*StringNode) Kind() NodeKind {
	return NdString
}

func (n *StringNode) String() string {
	return fmt.Sprintf("String{%s}", n.Str)
}

type CharNode struct {
	Char rune
}

func (*CharNode) Kind() NodeKind {
	return NdChar
}

func (n *CharNode) String() string {
	return fmt.Sprintf("Char{%c}", n.Char)
}

type VarRefNode struct {
	Var string
}

func (*VarRefNode) Kind() NodeKind {
	return NdVarRef
}

func (n *VarRefNode) String() string {
	return fmt.Sprintf("VarRef{%s}", n.Var)
}

type VarDefNode struct {
	VarType ValueType
	Var     string
	Value   Node
}

func (*VarDefNode) Kind() NodeKind {
	return NdVarDef
}

func (n *VarDefNode) String() string {
	return fmt.Sprintf("VarDef{%s/%s = %s}", n.Var, n.VarType, n.Value)
}

type FuncCallNode struct {
	Func string
	Args []Node
}

func (*FuncCallNode) Kind() NodeKind {
	return NdFuncCall
}

func (n *FuncCallNode) String() string {
	return fmt.Sprintf("FuncCall{%s(%d)}", n.Func, len(n.Args))
}

type FuncDefNode struct {
	Func string
	Args map[string]ValueType
	Ret  ValueType
	Body []Node
}

func (*FuncDefNode) Kind() NodeKind {
	return NdFuncDef
}

func (n *FuncDefNode) String() string {
	argText := ""
	for n, t := range n.Args {
		argText += n
		argText += "/"
		argText += t.String()
		argText += " "
	}
	bodyText := ""
	if len(n.Body) == 0 {
		bodyText = "(nothing?)"
	}
	if len(n.Body) == 1 {
		bodyText = fmt.Sprint(n.Body[0])
	}
	if len(n.Body) > 1 {
		bodyText = fmt.Sprintf("[... %s]", n.Body[len(n.Body)-1])
	}
	return fmt.Sprintf("FuncDef{%s(%s) = %s}", n.Func, argText, bodyText)
}

type ReturnNode struct {
	Value Node
}

func (*ReturnNode) Kind() NodeKind {
	return NdReturn
}
