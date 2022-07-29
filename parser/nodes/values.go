package nodes

import (
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

type FuncCallNode struct {
	Func string
	Args []Node
	Pos  lexer.Position
}

func (*FuncCallNode) Kind() NodeKind {
	return NdFuncCall
}

func (n *FuncCallNode) Where() lexer.Position {
	return n.Pos
}

type SelectorNode struct {
	Parent *SelectorNode
	Child  string
	Pos    lexer.Position
}

func (*SelectorNode) Kind() NodeKind {
	return NdSelector
}

func (n *SelectorNode) Where() lexer.Position {
	return n.Pos
}

func (n *SelectorNode) Path() []string {
	if n.Parent == nil {
		return []string{n.Child}
	}
	return append(n.Parent.Path(), n.Child)
}

type TypecastNode struct {
	Value  *SelectorNode
	Target types.Type
	Pos    lexer.Position
}

func (*TypecastNode) Kind() NodeKind {
	return NdTypecast
}

func (n *TypecastNode) Where() lexer.Position {
	return n.Pos
}
