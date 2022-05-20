package parser

type Scope struct {
	parent *Scope
	Funcs  map[string]*FuncDefNode
	Vars   map[string]*VarDefNode
}

func (s *Scope) FindFunc(name string) (*FuncDefNode, bool) {
	f, ok := s.Funcs[name]
	if s.parent != nil && !ok && name[0] != '_' {
		return s.parent.FindFunc(name)
	}
	return f, ok
}

func (s *Scope) FindVar(name string) (*VarDefNode, bool) {
	v, ok := s.Vars[name]
	if s.parent != nil && !ok && name[0] != '_' {
		return s.parent.FindVar(name)
	}
	return v, ok
}
