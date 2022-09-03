package ast

import (
	"github.com/syzkrash/skol/parser/values/types"
)

type Var struct {
	Name  string
	Value MetaNode
	Node  MetaNode
}

type Typedef struct {
	Name string
	Type types.Type
	Node MetaNode
}

type Func struct {
	Name string
	Args []FuncProtoArg
	Ret  types.Type
	Body Block
	Node MetaNode
}

type Extern struct {
	Name  string
	Alias string
	Ret   types.Type
	Args  []FuncProtoArg
	Node  MetaNode
}

type Structure struct {
	Name   string
	Fields []StructProtoField
	Node   MetaNode
}

type AST struct {
	Vars     map[string]Var
	Typedefs map[string]Typedef
	Funcs    map[string]Func
	Exerns   map[string]Extern
	Structs  map[string]Structure
}
