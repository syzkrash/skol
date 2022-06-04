package sim

type Scope struct {
	parent *Scope
	Vars   map[string]*Value
	Funcs  map[string]*Funct
}

func (s *Scope) FindVar(n string) (*Value, bool) {
	if v, ok := s.Vars[n]; ok {
		return v, true
	}
	if s.parent != nil && n[0] != '_' {
		return s.parent.FindVar(n)
	}
	return nil, false
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
