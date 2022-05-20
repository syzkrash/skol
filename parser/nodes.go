package parser

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
	Char rune
}

func (*CharNode) Kind() NodeKind {
	return NdChar
}

type VarRefNode struct {
	Var string
}

func (*VarRefNode) Kind() NodeKind {
	return NdVarRef
}

type VarDefNode struct {
	VarType ValueType
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
	Func string
	Args map[string]ValueType
	Ret  ValueType
	Body []Node
}

func (*FuncDefNode) Kind() NodeKind {
	return NdFuncDef
}

type ReturnNode struct {
	Value Node
}

func (*ReturnNode) Kind() NodeKind {
	return NdReturn
}
