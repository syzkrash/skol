package ast

import "github.com/syzkrash/skol/lexer"

// Branch represents the condition and body of an if statement branch.
type Branch struct {
	Cond  MetaNode
	Block Block
	Pos   lexer.Position
}

// IfNode represents an if statement. Contains the main branch, any other
// (else-if) branches and the else branch.
type IfNode struct {
	Main  Branch
	Other []Branch
	Else  Block
}

var _ Node = IfNode{}

func (IfNode) Kind() NodeKind {
	return NIf
}

type WhileNode struct {
	Cond  MetaNode
	Block Block
}

var _ Node = WhileNode{}

func (WhileNode) Kind() NodeKind {
	return NWhile
}

type ReturnNode struct {
	Value MetaNode
}

var _ Node = ReturnNode{}

func (ReturnNode) Kind() NodeKind {
	return NReturn
}
