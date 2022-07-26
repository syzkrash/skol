package parser

import (
	"github.com/syzkrash/skol/parser/defaults"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

type Scope struct {
	Parent *Scope
	Funcs  map[string]*values.Function
	Vars   map[string]*nodes.VarDefNode
	Consts map[string]*values.Value
	Types  map[string]types.Type
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent: parent,
		Funcs:  defaults.Functions,
		Vars:   make(map[string]*nodes.VarDefNode),
		Consts: make(map[string]*values.Value),
		Types:  make(map[string]types.Type),
	}
}

func (s *Scope) FindFunc(name string) (*values.Function, bool) {
	f, ok := s.Funcs[name]
	if s.Parent != nil && !ok && name[0] != '_' {
		return s.Parent.FindFunc(name)
	}
	return f, ok
}

func (s *Scope) FindVar(name string) (*nodes.VarDefNode, bool) {
	v, ok := s.Vars[name]
	if s.Parent != nil && !ok && name[0] != '_' {
		return s.Parent.FindVar(name)
	}
	return v, ok
}

func (s *Scope) SetVar(n string, v *nodes.VarDefNode) {
	var target *Scope
	current := s
	for {
		if _, ok := current.Vars[n]; ok || current.Parent == nil {
			target = current
			break
		}
		current = current.Parent
	}
	if target == nil {
		s.Vars[n] = v
	} else {
		target.Vars[n] = v
	}
}

func (s *Scope) FindConst(n string) (*values.Value, bool) {
	v, ok := s.Consts[n]
	if s.Parent != nil && !ok {
		return s.Parent.FindConst(n)
	}
	return v, ok
}

func (s *Scope) SetConst(n string, v *values.Value) bool {
	if _, exists := s.FindConst(n); exists {
		return false
	}
	s.Consts[n] = v
	return true
}

func (s *Scope) FindType(n string) (types.Type, bool) {
	t, ok := s.Types[n]
	if s.Parent != nil && !ok {
		return s.Parent.FindType(n)
	}
	return t, ok
}
