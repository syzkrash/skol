package ast

import "github.com/syzkrash/skol/parser/values/types"

// Var represents a global variable along which has a statically known value.
type Var struct {
	Name  string
	Value MetaNode
	Node  MetaNode
}

// Typedef represents a global variable that only has a known type.
type Typedef struct {
	Name string
	Type types.Type
	Node MetaNode
}

// Func represents a global function definition with it's body.
type Func struct {
	Name string
	Args []types.Descriptor
	Ret  types.Type
	Body Block
	Node MetaNode
}

// Extern represents a global external function with an unknown body.
type Extern struct {
	Name  string
	Alias string
	Ret   types.Type
	Args  []types.Descriptor
	Node  MetaNode
}

// Structure represents a global structure type definition.
type Structure struct {
	Name   string
	Fields []types.Descriptor
	Node   MetaNode
}

// AST is the complete Abstract Syntax Tree of a Skol source file.
type AST struct {
	Vars     map[string]Var
	Typedefs map[string]Typedef
	Funcs    map[string]Func
	Exerns   map[string]Extern
	Structs  map[string]Structure
}

func NewAST() AST {
	return AST{
		Vars:     make(map[string]Var),
		Typedefs: make(map[string]Typedef),
		Funcs:    make(map[string]Func),
		Exerns:   make(map[string]Extern),
		Structs:  make(map[string]Structure),
	}
}
