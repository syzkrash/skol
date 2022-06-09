package parser

import (
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type Scope struct {
	Parent *Scope
	Funcs  map[string]*Function
	Vars   map[string]*nodes.VarDefNode
	Consts map[string]*values.Value
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent: parent,
		Funcs:  DefaultFuncs,
		Vars:   make(map[string]*nodes.VarDefNode),
		Consts: make(map[string]*values.Value),
	}
}

func (s *Scope) FindFunc(name string) (*Function, bool) {
	f, ok := s.Funcs[name]
	if s.Parent != nil && !ok && name[0] != '_' {
		return s.Parent.FindFunc(name)
	}
	return f, ok
}

func (s *Scope) SetFunc(n string, f *Function) {
	if _, ok := s.Funcs[n]; ok || s.Parent == nil {
		s.Funcs[n] = f
	} else {
		s.Parent.SetFunc(n, f)
	}
}

func (s *Scope) FindVar(name string) (*nodes.VarDefNode, bool) {
	v, ok := s.Vars[name]
	if s.Parent != nil && !ok && name[0] != '_' {
		return s.Parent.FindVar(name)
	}
	return v, ok
}

func (s *Scope) SetVar(n string, v *nodes.VarDefNode) {
	if _, ok := s.Vars[n]; ok || s.Parent == nil {
		s.Vars[n] = v
	} else {
		s.Parent.SetVar(n, v)
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
