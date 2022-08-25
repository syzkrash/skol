package nodes

import (
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

type SelectorElem struct {
	Idx  Node
	Name string
	Cast types.Type
}

type Selector interface {
	Node
	Path() []SelectorElem
}

type SelectorNode struct {
	Parent Selector
	Child  string
	Pos    lexer.Position
}

func (*SelectorNode) Kind() NodeKind {
	return NdSelector
}

func (n *SelectorNode) Where() lexer.Position {
	return n.Pos
}

func (n *SelectorNode) Path() []SelectorElem {
	if n.Parent == nil {
		return []SelectorElem{{Name: n.Child}}
	}
	return append(n.Parent.Path(), SelectorElem{Name: n.Child})
}

type TypecastNode struct {
	Parent Selector
	Type   types.Type
	Pos    lexer.Position
}

func (*TypecastNode) Kind() NodeKind {
	return NdTypecast
}

func (n *TypecastNode) Where() lexer.Position {
	return n.Pos
}

func (n *TypecastNode) Path() []SelectorElem {
	if n.Parent == nil {
		return []SelectorElem{{Cast: n.Type}}
	}
	return append(n.Parent.Path(), SelectorElem{Cast: n.Type})
}

type IndexNode struct {
	Parent Selector
	Idx    Node
	Pos    lexer.Position
}

func (*IndexNode) Kind() NodeKind {
	return NdIndex
}

func (n *IndexNode) Where() lexer.Position {
	return n.Pos
}

func (n *IndexNode) Path() []SelectorElem {
	if n.Parent == nil {
		return []SelectorElem{{Idx: n.Idx}}
	}
	return append(n.Parent.Path(), SelectorElem{Idx: n.Idx})
}
