package ast

import (
	"github.com/syzkrash/skol/parser/values/types"
)

// VarSetNode represents a variable assignment:
//
//	%MyVar: "New Value"
type VarSetNode struct {
	Var   string
	Value MetaNode
}

var _ Node = VarSetNode{}

func (VarSetNode) Kind() NodeKind {
	return NVarSet
}

// VarDef represents a type definition for a variable:
//
//	%MyVar/string
type VarDefNode struct {
	Var  string
	Type types.Type
}

var _ Node = VarDefNode{}

func (VarDefNode) Kind() NodeKind {
	return NVarDef
}

// VarSetTypedNode represents a variable assignment with a type definition:
//
//	%MyVar/string: "New Value"
type VarSetTypedNode struct {
	Var   string
	Type  types.Type
	Value MetaNode
}

var _ Node = VarSetTypedNode{}

func (VarSetTypedNode) Kind() NodeKind {
	return NVarSetTyped
}

type FuncDefNode struct {
	Name  string
	Proto []types.Descriptor
	Ret   types.Type
	Body  Block
}

var _ Node = FuncDefNode{}

func (FuncDefNode) Kind() NodeKind {
	return NFuncDef
}

type FuncExternNode struct {
	Alias string
	Proto []types.Descriptor
	Ret   types.Type
	Name  string
}

var _ Node = FuncExternNode{}

func (FuncExternNode) Kind() NodeKind {
	return NFuncExtern
}

type StructDefNode struct {
	Name   string
	Fields []types.Descriptor
}

var _ Node = StructDefNode{}

func (StructDefNode) Kind() NodeKind {
	return NStructDef
}
