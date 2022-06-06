package parser

import "github.com/syzkrash/skol/parser/nodes"

type Scope struct {
	Parent *Scope
	Funcs  map[string]*Function
	Vars   map[string]*nodes.VarDefNode
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
