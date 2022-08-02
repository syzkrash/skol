package nodes

import (
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

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
