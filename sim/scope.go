package sim

import "github.com/syzkrash/skol/parser/values"

type Scope struct {
	parent *Scope
	Vars   map[string]*values.Value
	Funcs  map[string]*Funct
}

func (s *Scope) FindVar(n string) (*values.Value, bool) {
	if v, ok := s.Vars[n]; ok {
		return v, true
	}
	if s.parent != nil && n[0] != '_' {
		return s.parent.FindVar(n)
	}
	return nil, false
}

func (s *Scope) SetVar(n string, v *values.Value) {
	if _, ok := s.Vars[n]; ok || s.parent == nil {
		s.Vars[n] = v
	} else {
		s.parent.SetVar(n, v)
	}
}

func (s *Scope) FindFunc(n string) (*Funct, bool) {
	if v, ok := s.Funcs[n]; ok {
		return v, true
	}
	if s.parent != nil && n[0] != '_' {
		return s.parent.FindFunc(n)
	}
	return nil, false
}

func (s *Scope) SetFunc(n string, f *Funct) {
	if _, ok := s.Funcs[n]; ok || s.parent == nil {
		s.Funcs[n] = f
	} else {
		s.parent.SetFunc(n, f)
	}
}
