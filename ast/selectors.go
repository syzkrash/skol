package ast

import "github.com/syzkrash/skol/parser/values/types"

// SelectorElem represents an element of a [Selector] path. The exact value
// for an element is determined by the rule of elimination:
//   - If Cast != nil, this element is a typecast.
//   - If Name != "", this element is a variable/field access.
//   - If IdxS != nil, this element is an array index using a variable.
//   - In all other cases, this element is an array index using a constant
//     integer.
type SelectorElem struct {
	Cast types.Type
	Name string
	IdxS Selector
	IdxC int
}

// Selector represents any node that can be used as a selector element.
type Selector interface {
	Node
	// Path gets the total selector path for this selector and all it's parents.
	Path() []SelectorElem
}

type TypecastNode struct {
	Parent Selector
	Cast   types.Type
}

// don't need to check if it implements Node because Selector requires Node
var _ Selector = TypecastNode{}

func (TypecastNode) Kind() NodeKind {
	return NTypecast
}

func (n TypecastNode) Path() []SelectorElem {
	if n.Parent == nil {
		panic("TypecastNode without a Parent!")
	}
	return append(n.Parent.Path(), SelectorElem{Cast: n.Cast})
}

type SelectorNode struct {
	Parent Selector
	Child  string
}

var _ Selector = SelectorNode{}

func (SelectorNode) Kind() NodeKind {
	return NSelector
}

func (n SelectorNode) Path() []SelectorElem {
	if n.Parent == nil {
		return []SelectorElem{{Name: n.Child}}
	}
	return append(n.Parent.Path(), SelectorElem{Name: n.Child})
}

type IndexSelectorNode struct {
	Parent Selector
	Idx    Selector
}

var _ Selector = IndexSelectorNode{}

func (IndexSelectorNode) Kind() NodeKind {
	return NIndexSelector
}

func (n IndexSelectorNode) Path() []SelectorElem {
	if n.Parent == nil {
		panic("IndexSelectorNode without a Parent!")
	}
	return append(n.Parent.Path(), SelectorElem{IdxS: n.Idx})
}

type IndexConstNode struct {
	Parent Selector
	Idx    int
}

var _ Selector = IndexConstNode{}

func (IndexConstNode) Kind() NodeKind {
	return NIndexConst
}

func (n IndexConstNode) Path() []SelectorElem {
	if n.Parent == nil {
		panic("IndexConstNode without a Parent!")
	}
	return append(n.Parent.Path(), SelectorElem{IdxC: n.Idx})
}
