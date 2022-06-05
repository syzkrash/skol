package sim

import (
	"fmt"

	"github.com/syzkrash/skol/parser"
)

type Simulator struct {
	scope *Scope
}

func NewSimulator() *Simulator {
	return &Simulator{
		scope: &Scope{
			parent: nil,
			Vars:   map[string]*Value{},
			Funcs:  map[string]*Funct{},
		},
	}
}

func (s *Simulator) Stmt(n parser.Node) error {
	switch n.Kind() {
	case parser.NdVarDef:
		vdn := n.(*parser.VarDefNode)
		val, err := s.Expr(vdn.Value)
		if err != nil {
			return err
		}
		s.scope.Vars[vdn.Var] = val
		return nil
	case parser.NdFuncDef:
		fdn := n.(*parser.FuncDefNode)
		s.scope.Funcs[fdn.Name] = &Funct{
			Args: fdn.Args,
			Ret:  fdn.Ret,
			Body: fdn.Body,
		}
		return nil
	case parser.NdFuncExtern:
		fen := n.(*parser.FuncExternNode)
		s.scope.Funcs[fen.Name] = &Funct{
			Args: fen.Args,
			Ret:  fen.Ret,
			Body: []parser.Node{},
		}
		return nil
	case parser.NdIf:
		ifn := n.(*parser.IfNode)
		val, err := s.Expr(ifn.Condition)
		if err != nil {
			return err
		}
		if val.ToBool() {
			return s.block(ifn.IfBlock)
		}
		for _, branch := range ifn.ElseIfNodes {
			val, err = s.Expr(branch.Condition)
			if err != nil {
				return err
			}
			if val.ToBool() {
				return s.block(branch.Block)
			}
		}
		return s.block(ifn.ElseBlock)
	case parser.NdWhile:
		whn := n.(*parser.WhileNode)
		val, err := s.Expr(whn.Condition)
		if err != nil {
			return err
		}
		if !val.ToBool() {
			return nil
		}
		for {
			err = s.block(whn.Body)
			if err != nil {
				return err
			}
			val, err = s.Expr(whn.Condition)
			if err != nil {
				return err
			}
			if !val.ToBool() {
				return nil
			}
		}
	case parser.NdFuncCall:
		fcn := n.(*parser.FuncCallNode)
		if fcn.Func == "print" {
			strs := []any{}
			for _, a := range fcn.Args {
				v, err := s.Expr(a)
				if err != nil {
					return err
				}
				strs = append(strs, v.String())
			}
			fmt.Println(strs...)
			return nil
		}
		funct, ok := s.scope.FindFunc(fcn.Func)
		if !ok {
			return fmt.Errorf("unknown function: %s", fcn.Func)
		}
		argn := []string{}
		for n := range funct.Args {
			argn = append(argn, n)
		}
		argv := map[string]*Value{}
		for i := 0; i < len(fcn.Args); i++ {
			val, err := s.Expr(fcn.Args[i])
			if err != nil {
				return err
			}
			argv[argn[i]] = val
		}
		s.scope = &Scope{
			parent: s.scope,
			Vars:   argv,
			Funcs:  map[string]*Funct{},
		}
		for _, n := range funct.Body {
			err := s.Stmt(n)
			if err != nil {
				return err
			}
		}
		s.scope = s.scope.parent
		return nil
	}
	return fmt.Errorf("%s node is not a statement", n.Kind())
}

func (s *Simulator) Expr(n parser.Node) (*Value, error) {
	switch n.Kind() {
	case parser.NdInteger:
		return NewValue(n.(*parser.IntegerNode).Int), nil
	case parser.NdBoolean:
		return NewValue(n.(*parser.BooleanNode).Bool), nil
	case parser.NdFloat:
		return NewValue(n.(*parser.FloatNode).Float), nil
	case parser.NdString:
		return NewValue(n.(*parser.StringNode).Str), nil
	case parser.NdChar:
		return NewValue(n.(*parser.CharNode).Char), nil
	case parser.NdVarRef:
		vrn := n.(*parser.VarRefNode)
		val, ok := s.scope.FindVar(vrn.Var)
		if !ok {
			return nil, fmt.Errorf("unknown variable: %s", vrn.Var)
		}
		return val, nil
	case parser.NdFuncCall:
		fcn := n.(*parser.FuncCallNode)
		if fcn.Func == "print" {
			fmt.Println(fcn.Args)
		}
		funct, ok := s.scope.FindFunc(fcn.Func)
		if !ok {
			return nil, fmt.Errorf("unknown function: %s", fcn.Func)
		}
		if funct.Ret == parser.VtNothing {
			return nil, fmt.Errorf("function %s does not return a value", fcn.Func)
		}
		argn := []string{}
		for n := range funct.Args {
			argn = append(argn, n)
		}
		argv := map[string]*Value{}
		for i := 0; i < len(fcn.Args); i++ {
			val, err := s.Expr(fcn.Args[i])
			if err != nil {
				return nil, err
			}
			argv[argn[i]] = val
		}
		s.scope = &Scope{
			parent: s.scope,
			Vars:   argv,
			Funcs:  map[string]*Funct{},
		}
		var val *Value = Default(funct.Ret)
		var err error
		for _, n := range funct.Body {
			if n.Kind() == parser.NdReturn {
				val, err = s.Expr(n)
				if err != nil {
					return nil, err
				}
				break
			}
			err := s.Stmt(n)
			if err != nil {
				return nil, err
			}
		}
		s.scope = s.scope.parent
		return val, nil
	}
	return nil, fmt.Errorf("%s node is not a value", n.Kind())
}

func (s *Simulator) block(b []parser.Node) error {
	s.scope = &Scope{
		parent: s.scope,
		Vars:   map[string]*Value{},
		Funcs:  map[string]*Funct{},
	}
	for _, n := range b {
		if err := s.Stmt(n); err != nil {
			return err
		}
	}
	s.scope = s.scope.parent
	return nil
}
