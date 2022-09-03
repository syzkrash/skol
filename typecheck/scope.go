package typecheck

import (
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/values/types"
)

type funcproto func(ast.MetaNode, []types.Type) (types.Type, *TypeError)

type scope struct {
	parent *scope
	vars   map[string]types.Type
	funcs  map[string]funcproto
}

func (s *scope) sub() *scope {
	return &scope{
		parent: s,
		vars:   make(map[string]types.Type),
		funcs:  make(map[string]funcproto),
	}
}

func (s *scope) getVar(name string) (types.Type, bool) {
	if t, ok := s.vars[name]; ok {
		return t, true
	} else if s.parent != nil {
		return s.parent.getVar(name)
	} else {
		return nil, false
	}
}

func (s *scope) setVar(name string, t types.Type) {
	var target *scope
	current := s
	for {
		if _, ok := current.vars[name]; ok || current.parent == nil {
			target = current
			break
		}
		current = current.parent
	}
	if target == nil {
		s.vars[name] = t
	} else {
		target.vars[name] = t
	}
}

func (s *scope) getFunc(name string) (funcproto, bool) {
	if f, ok := s.funcs[name]; ok {
		return f, true
	} else if s.parent != nil {
		return s.parent.getFunc(name)
	} else {
		return nil, false
	}
}
