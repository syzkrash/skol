package parser

import (
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/defaults"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

type Scope struct {
	Parent *Scope
	Funcs  map[string]*values.Function
	Vars   map[string]ast.Node
	Consts map[string]ast.Node
	Types  map[string]types.Type
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent: parent,
		Funcs:  defaults.Functions,
		Vars:   make(map[string]ast.Node),
		Consts: make(map[string]ast.Node),
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

func (s *Scope) FindVar(name string) (ast.Node, bool) {
	v, ok := s.Vars[name]
	if s.Parent != nil && !ok && name[0] != '_' {
		return s.Parent.FindVar(name)
	}
	return v, ok
}

func (s *Scope) SetVar(n string, v ast.Node) {
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

func (s *Scope) FindConst(n string) (ast.Node, bool) {
	v, ok := s.Consts[n]
	if s.Parent != nil && !ok {
		return s.Parent.FindConst(n)
	}
	return v, ok
}

func (s *Scope) SetConst(n string, v ast.Node) bool {
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
