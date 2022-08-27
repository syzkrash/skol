package ast

import (
	"github.com/syzkrash/skol/parser/values/types"
)

type BoolNode struct {
	Value bool
}

var _ Node = BoolNode{}

func (BoolNode) Kind() NodeKind {
	return NBool
}

type CharNode struct {
	Value byte
}

var _ Node = CharNode{}

func (CharNode) Kind() NodeKind {
	return NChar
}

type IntNode struct {
	Value int32
}

var _ Node = IntNode{}

func (IntNode) Kind() NodeKind {
	return NInt
}

type FloatNode struct {
	Value float32
}

var _ Node = FloatNode{}

func (FloatNode) Kind() NodeKind {
	return NFloat
}

type StringNode struct {
	Value string
}

var _ Node = StringNode{}

func (StringNode) Kind() NodeKind {
	return NString
}

type StructNode struct {
	Type types.StructType
	Args []MetaNode
}

var _ Node = StructNode{}

func (StructNode) Kind() NodeKind {
	return NStruct
}

type ArrayNode struct {
	Type  types.ArrayType
	Elems []MetaNode
}

var _ Node = ArrayNode{}

func (ArrayNode) Kind() NodeKind {
	return NArray
}
