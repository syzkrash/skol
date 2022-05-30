package parser

type Scope struct {
	Parent *Scope
	Funcs  map[string]*Function
	Vars   map[string]*VarDefNode
}

func (s *Scope) FindFunc(name string) (*Function, bool) {
	f, ok := s.Funcs[name]
	if s.Parent != nil && !ok && name[0] != '_' {
		return s.Parent.FindFunc(name)
	}
	return f, ok
}

func (s *Scope) FindVar(name string) (*VarDefNode, bool) {
	v, ok := s.Vars[name]
	if s.Parent != nil && !ok && name[0] != '_' {
		return s.Parent.FindVar(name)
	}
	return v, ok
}
