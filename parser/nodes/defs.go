package nodes

import (
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

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
