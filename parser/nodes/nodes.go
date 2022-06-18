package nodes

import (
	"fmt"
	"strings"

	"github.com/syzkrash/skol/parser/values"
)

func body(n []Node) (text string) {
	if len(n) == 0 {
		text = "(nothing?)"
	}
	if len(n) == 1 {
		text = fmt.Sprint(n[0])
	}
	if len(n) > 1 {
		text = fmt.Sprintf("[... %s]", n[len(n)-1])
	}
	return text
}

type NodeKind uint8

const (
	NdInteger NodeKind = iota
	NdBoolean
	NdFloat
	NdString
	NdChar
	NdVarRef
	NdVarDef
	NdFuncCall
	NdFuncDef
	NdFuncExtern
	NdReturn
	NdIf
	NdWhile
	NdStruct
)

var nodeKinds = []string{
	"Integer",
	"Boolean",
	"Float",
	"String",
	"Char",
	"VarRef",
	"VarDef",
	"FuncCall",
	"FuncDef",
	"FuncExtern",
	"Return",
	"If",
	"While",
	"Struct",
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

type BooleanNode struct {
	Bool bool
}

func (*BooleanNode) Kind() NodeKind {
	return NdBoolean
}

func (n *BooleanNode) String() string {
	return fmt.Sprintf("Boolean{%v}", n.Bool)
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
	Char byte
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
	VarType *values.Type
	Var     string
	Value   Node
}

func (*VarDefNode) Kind() NodeKind {
	return NdVarDef
}

func (n *VarDefNode) String() string {
	return fmt.Sprintf("VarDef{%s %s = %s}", n.VarType, n.Var, n.Value)
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
	Name string
	Args []values.FuncArg
	Ret  *values.Type
	Body []Node
}

func (*FuncDefNode) Kind() NodeKind {
	return NdFuncDef
}

func (n *FuncDefNode) String() string {
	argText := ""
	for _, a := range n.Args {
		argText += a.Type.String()
		argText += " "
		argText += a.Name
		argText += " "
	}
	argText = strings.TrimSuffix(argText, " ")
	bodyText := body(n.Body)
	return fmt.Sprintf("FuncDef{%s %s(%s) = %s}", n.Ret, n.Name, argText, bodyText)
}

type FuncExternNode struct {
	Name   string
	Intern string
	Args   []values.FuncArg
	Ret    *values.Type
}

func (*FuncExternNode) Kind() NodeKind {
	return NdFuncExtern
}

func (n *FuncExternNode) String() string {
	argText := ""
	for _, a := range n.Args {
		argText += a.Type.String()
		argText += " "
		argText += a.Name
		argText += " "
	}
	argText = strings.TrimSuffix(argText, " ")
	return fmt.Sprintf("FuncExtern{%s %s(%s)}", n.Ret, n.Name, argText)
}

type ReturnNode struct {
	Value Node
}

func (*ReturnNode) Kind() NodeKind {
	return NdReturn
}

func (n *ReturnNode) String() string {
	return fmt.Sprintf("Return{%s}", n.Value)
}

type IfSubNode struct {
	Condition Node
	Block     []Node
}

type IfNode struct {
	Condition   Node
	IfBlock     []Node
	ElseIfNodes []*IfSubNode
	ElseBlock   []Node
}

func (*IfNode) Kind() NodeKind {
	return NdIf
}

func (n *IfNode) String() string {
	return fmt.Sprintf(
		"If{%s(%s) else(%s)}",
		n.Condition,
		body(n.IfBlock),
		body(n.ElseBlock))
}

type WhileNode struct {
	Condition Node
	Body      []Node
}

func (*WhileNode) Kind() NodeKind {
	return NdWhile
}

func (n *WhileNode) String() string {
	return fmt.Sprintf("While{%s(%s)}", n.Condition, body(n.Body))
}

type StructNode struct {
	Name string
	Type *values.Type
}

func (*StructNode) Kind() NodeKind {
	return NdStruct
}
