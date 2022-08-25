package nodes

import (
	"github.com/syzkrash/skol/lexer"
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
