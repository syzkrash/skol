package ast

type FuncCallNode struct {
	Func string
	Args []MetaNode
}

var _ Node = FuncCallNode{}

func (FuncCallNode) Kind() NodeKind {
	return NFuncCall
}
