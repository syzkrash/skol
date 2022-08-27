package ast

import "github.com/syzkrash/skol/lexer"

type NodeKind byte

const (
	NInvalid NodeKind = iota

	// literals
	NBool
	NChar
	NInt
	NFloat
	NString
	NStruct
	NArray

	// control flow
	NIf
	NWhile
	NReturn

	// definitions
	NVarSet
	NVarDef
	NVarSetTyped
	NFuncDef
	NFuncExtern
	NStructDef

	// selectors
	NSelector
	NTypecast
	NIndexConst
	NIndexSelector

	// others
	NFuncCall

	// max bound
	NMax
)

var nodeKindNames = []string{
	"Invalid",
	"Bool",
	"Char",
	"Int",
	"Float",
	"String",
	"Struct",
	"Array",
	"If",
	"While",
	"Return",
	"VarSet",
	"VarDef",
	"VarTypedSet",
	"FuncDef",
	"FuncExtern",
	"StructDef",
	"Selector",
	"Typecast",
	"IndexConst",
	"IndexSelector",
	"FuncCall",
}

func (k NodeKind) Ensure() NodeKind {
	if k >= NMax {
		return NInvalid
	}
	return k
}

func (k NodeKind) String() string {
	return nodeKindNames[k.Ensure()]
}

type Node interface {
	Kind() NodeKind
}

type MetaNode struct {
	Node  Node
	Where lexer.Position
}

type Block []MetaNode
