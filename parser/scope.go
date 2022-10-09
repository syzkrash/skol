package parser

import (
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/values/types"
)

// Scope represents the current lexical scope and its parent scope if any.
type Scope struct {
	Parent *Scope
	Vars   map[string]ast.Node
	Consts map[string]ast.Node
	Types  map[string]types.Type
}

// NewScope creates a new lexical scope, e.g. when entering a function. If
// parent is nil the scope can be considered the global scope.
func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent: parent,
		Vars:   make(map[string]ast.Node),
		Consts: make(map[string]ast.Node),
		Types:  make(map[string]types.Type),
	}
}

// FindVar searches this and all parent scopes for the given variable.
func (s *Scope) FindVar(name string) (ast.Node, bool) {
	v, ok := s.Vars[name]
	if s.Parent != nil && !ok && name[0] != '_' {
		return s.Parent.FindVar(name)
	}
	return v, ok
}

// SetVar updates the given variable value in this scope or one of the parent
// scopes. If this variable is found in a parent scope, it is updated in that
// parent scope. If it isn't found, it is created in this scope.
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

// FindConst searches this and all parent scopes for the given constant.
func (s *Scope) FindConst(n string) (ast.Node, bool) {
	v, ok := s.Consts[n]
	if s.Parent != nil && !ok {
		return s.Parent.FindConst(n)
	}
	return v, ok
}

// SetConst sets a constant value, ensuring it is not changed.
func (s *Scope) SetConst(n string, v ast.Node) bool {
	if _, exists := s.FindConst(n); exists {
		return false
	}
	s.Consts[n] = v
	return true
}

// FindType searches this and all parent scopes for the given structure type.
func (s *Scope) FindType(n string) (types.Type, bool) {
	t, ok := s.Types[n]
	if s.Parent != nil && !ok {
		return s.Parent.FindType(n)
	}
	return t, ok
}
