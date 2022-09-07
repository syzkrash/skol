package ast

import "github.com/syzkrash/skol/lexer"

// NodeKind is the unique identifying number of an abstract AST node.
type NodeKind byte

// Node kind constants.
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

// for [NodeKid.String]
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

// Ensure checks if this is a valid NodeKind, returning NInvalid if it's not.
func (k NodeKind) Ensure() NodeKind {
	if k >= NMax {
		return NInvalid
	}
	return k
}

// String returns the name of this NodeKind if it is valid, "Invalid" otherwise.
func (k NodeKind) String() string {
	return nodeKindNames[k.Ensure()]
}

// Node represents an abstract AST node.
type Node interface {
	Kind() NodeKind
}

// MetaNode wraps an abstract node with position information.
type MetaNode struct {
	Node  Node
	Where lexer.Position
}

// Block represents a list of MetaNodes, typically for multiple statements.
type Block []MetaNode
