package ast

import (
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

type Var struct {
	Name  string
	Value Node
	Where lexer.Position
}

type Typedef struct {
	Name  string
	Type  types.Type
	Where lexer.Position
}

type Func struct {
	Name  string
	Args  []FuncProtoArg
	Ret   types.Type
	Body  Block
	Where lexer.Position
}

type Extern struct {
	Name  string
	Alias string
	Ret   types.Type
	Args  []FuncProtoArg
	Where lexer.Position
}

type Structure struct {
	Name   string
	Fields []StructProtoField
	Where  lexer.Position
}

type AST struct {
	Vars     map[string]Var
	Typedefs map[string]Typedef
	Funcs    map[string]Func
	Exerns   map[string]Extern
	Structs  map[string]Structure
}
