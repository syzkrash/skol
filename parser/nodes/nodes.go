package nodes

import "github.com/syzkrash/skol/parser/values"

type NodeKind uint8

const (
	NdInteger NodeKind = iota
	NdBoolean
	NdFloat
	NdString
	NdChar
	NdVarDef
	NdFuncCall
	NdFuncDef
	NdFuncExtern
	NdReturn
	NdIf
	NdWhile
	NdStruct
	NdNewStruct
	NdSelector
)

var nodeKinds = []string{
	"Integer",
	"Boolean",
	"Float",
	"String",
	"Char",
	"VarDef",
	"FuncCall",
	"FuncDef",
	"FuncExtern",
	"Return",
	"If",
	"While",
	"Struct",
	"NewStruct",
	"Selector",
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

type BooleanNode struct {
	Bool bool
}

func (*BooleanNode) Kind() NodeKind {
	return NdBoolean
}

type FloatNode struct {
	Float float32
}

func (*FloatNode) Kind() NodeKind {
	return NdFloat
}

type StringNode struct {
	Str string
}

func (*StringNode) Kind() NodeKind {
	return NdString
}

type CharNode struct {
	Char byte
}

func (*CharNode) Kind() NodeKind {
	return NdChar
}

type VarDefNode struct {
	VarType *values.Type
	Var     string
	Value   Node
}

func (*VarDefNode) Kind() NodeKind {
	return NdVarDef
}

type FuncCallNode struct {
	Func string
	Args []Node
}

func (*FuncCallNode) Kind() NodeKind {
	return NdFuncCall
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

type FuncExternNode struct {
	Name   string
	Intern string
	Args   []values.FuncArg
	Ret    *values.Type
}

func (*FuncExternNode) Kind() NodeKind {
	return NdFuncExtern
}

type ReturnNode struct {
	Value Node
}

func (*ReturnNode) Kind() NodeKind {
	return NdReturn
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

type WhileNode struct {
	Condition Node
	Body      []Node
}

func (*WhileNode) Kind() NodeKind {
	return NdWhile
}

type StructNode struct {
	Name string
	Type *values.Type
}

func (*StructNode) Kind() NodeKind {
	return NdStruct
}

type NewStructNode struct {
	Type *values.Type
	Args []Node
}

func (*NewStructNode) Kind() NodeKind {
	return NdNewStruct
}

type SelectorNode struct {
	Parent *SelectorNode
	Child  string
}

func (*SelectorNode) Kind() NodeKind {
	return NdSelector
}

func (n *SelectorNode) Path() []string {
	if n.Parent == nil {
		return []string{n.Child}
	}
	return append(n.Parent.Path(), n.Child)
}
