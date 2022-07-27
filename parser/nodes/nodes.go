package nodes

import (
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

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
	NdTypecast
	NdArray
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
	"Typecast",
}

func (k NodeKind) String() string {
	return nodeKinds[k]
}

type Node interface {
	Kind() NodeKind
	Where() lexer.Position
}

type IntegerNode struct {
	Int int32
	Pos lexer.Position
}

func (*IntegerNode) Kind() NodeKind {
	return NdInteger
}

func (n *IntegerNode) Where() lexer.Position {
	return n.Pos
}

type BooleanNode struct {
	Bool bool
	Pos  lexer.Position
}

func (*BooleanNode) Kind() NodeKind {
	return NdBoolean
}

func (n *BooleanNode) Where() lexer.Position {
	return n.Pos
}

type FloatNode struct {
	Float float32
	Pos   lexer.Position
}

func (*FloatNode) Kind() NodeKind {
	return NdFloat
}

func (n *FloatNode) Where() lexer.Position {
	return n.Pos
}

type StringNode struct {
	Str string
	Pos lexer.Position
}

func (*StringNode) Kind() NodeKind {
	return NdString
}

func (n *StringNode) Where() lexer.Position {
	return n.Pos
}

type CharNode struct {
	Char byte
	Pos  lexer.Position
}

func (*CharNode) Kind() NodeKind {
	return NdChar
}

func (n *CharNode) Where() lexer.Position {
	return n.Pos
}

type VarDefNode struct {
	VarType types.Type
	Var     string
	Value   Node
	Pos     lexer.Position
}

func (*VarDefNode) Kind() NodeKind {
	return NdVarDef
}

func (n *VarDefNode) Where() lexer.Position {
	return n.Pos
}

type FuncCallNode struct {
	Func string
	Args []Node
	Pos  lexer.Position
}

func (*FuncCallNode) Kind() NodeKind {
	return NdFuncCall
}

func (n *FuncCallNode) Where() lexer.Position {
	return n.Pos
}

type FuncDefNode struct {
	Name string
	Args []values.FuncArg
	Ret  types.Type
	Body []Node
	Pos  lexer.Position
}

func (*FuncDefNode) Kind() NodeKind {
	return NdFuncDef
}

func (n *FuncDefNode) Where() lexer.Position {
	return n.Pos
}

type FuncExternNode struct {
	Name   string
	Intern string
	Args   []values.FuncArg
	Ret    types.Type
	Pos    lexer.Position
}

func (*FuncExternNode) Kind() NodeKind {
	return NdFuncExtern
}

func (n *FuncExternNode) Where() lexer.Position {
	return n.Pos
}

type ReturnNode struct {
	Value Node
	Pos   lexer.Position
}

func (*ReturnNode) Kind() NodeKind {
	return NdReturn
}

func (n *ReturnNode) Where() lexer.Position {
	return n.Pos
}

type IfSubNode struct {
	Condition Node
	Block     []Node
	Pos       lexer.Position
}

type IfNode struct {
	Condition   Node
	IfBlock     []Node
	ElseIfNodes []*IfSubNode
	ElseBlock   []Node
	Pos         lexer.Position
}

func (*IfNode) Kind() NodeKind {
	return NdIf
}

func (n *IfNode) Where() lexer.Position {
	return n.Pos
}

type WhileNode struct {
	Condition Node
	Body      []Node
	Pos       lexer.Position
}

func (*WhileNode) Kind() NodeKind {
	return NdWhile
}

func (n *WhileNode) Where() lexer.Position {
	return n.Pos
}

type StructNode struct {
	Name string
	Type types.Type
	Pos  lexer.Position
}

func (*StructNode) Kind() NodeKind {
	return NdStruct
}

func (n *StructNode) Where() lexer.Position {
	return n.Pos
}

type NewStructNode struct {
	Type types.Type
	Args []Node
	Pos  lexer.Position
}

func (*NewStructNode) Kind() NodeKind {
	return NdNewStruct
}

func (n *NewStructNode) Where() lexer.Position {
	return n.Pos
}

type SelectorNode struct {
	Parent *SelectorNode
	Child  string
	Pos    lexer.Position
}

func (*SelectorNode) Kind() NodeKind {
	return NdSelector
}

func (n *SelectorNode) Where() lexer.Position {
	return n.Pos
}

func (n *SelectorNode) Path() []string {
	if n.Parent == nil {
		return []string{n.Child}
	}
	return append(n.Parent.Path(), n.Child)
}

type TypecastNode struct {
	Value  *SelectorNode
	Target types.Type
	Pos    lexer.Position
}

func (*TypecastNode) Kind() NodeKind {
	return NdTypecast
}

func (n *TypecastNode) Where() lexer.Position {
	return n.Pos
}

type ArrayNode struct {
	Type     types.Type
	Elements []Node
	Pos      lexer.Position
}

func (*ArrayNode) Kind() NodeKind {
	return NdArray
}

func (n *ArrayNode) Where() lexer.Position {
	return n.Pos
}
