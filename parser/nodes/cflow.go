package nodes

import "github.com/syzkrash/skol/lexer"

type IfSubNode struct {
	Condition Node
	Block     []Node
	Pos       lexer.Position
}

type IfNode struct {
	Condition   Node
	IfBlock     []Node
	ElseIfNodes []*IfSubNode
	ElseBlock   []Node
	Pos         lexer.Position
}

func (*IfNode) Kind() NodeKind {
	return NdIf
}

func (n *IfNode) Where() lexer.Position {
	return n.Pos
}

type WhileNode struct {
	Condition Node
	Body      []Node
	Pos       lexer.Position
}

func (*WhileNode) Kind() NodeKind {
	return NdWhile
}

func (n *WhileNode) Where() lexer.Position {
	return n.Pos
}

type ReturnNode struct {
	Value Node
	Pos   lexer.Position
}

func (*ReturnNode) Kind() NodeKind {
	return NdReturn
}

func (n *ReturnNode) Where() lexer.Position {
	return n.Pos
}
