package ast

import (
	"github.com/syzkrash/skol/parser/values/types"
)

type VarSetNode struct {
	Var   string
	Value MetaNode
}

var _ Node = VarSetNode{}

func (VarSetNode) Kind() NodeKind {
	return NVarSet
}

type VarDefNode struct {
	Var  string
	Type types.Type
}

var _ Node = VarDefNode{}

func (VarDefNode) Kind() NodeKind {
	return NVarDef
}

type VarSetTypedNode struct {
	Var   string
	Type  types.Type
	Value MetaNode
}

var _ Node = VarSetTypedNode{}

func (VarSetTypedNode) Kind() NodeKind {
	return NVarSetTyped
}

type FuncProtoArg struct {
	Name string
	Type types.Type
}

type FuncDefNode struct {
	Name  string
	Proto []FuncProtoArg
	Ret   types.Type
	Body  Block
}

var _ Node = FuncDefNode{}

func (FuncDefNode) Kind() NodeKind {
	return NFuncDef
}

type FuncExternNode struct {
	Alias string
	Proto []FuncProtoArg
	Ret   types.Type
	Name  string
}

var _ Node = FuncExternNode{}

func (FuncExternNode) Kind() NodeKind {
	return NFuncExtern
}

type StructProtoField struct {
	Name string
	Type types.Type
}

type StructDefNode struct {
	Name   string
	Fields []StructProtoField
}

var _ Node = StructDefNode{}

func (StructDefNode) Kind() NodeKind {
	return NStructDef
}
