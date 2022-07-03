package sim

import (
	"fmt"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type Simulator struct {
	Scope *Scope
	Calls []string
}

func NewSimulator() *Simulator {
	return &Simulator{
		Scope: &Scope{
			parent: nil,
			Vars:   map[string]*values.Value{},
			Funcs:  DefaultFuncs,
		},
		Calls: []string{},
	}
}

func (s *Simulator) Stmt(n nodes.Node) error {
	switch n.Kind() {
	case nodes.NdVarDef:
		vdn := n.(*nodes.VarDefNode)
		val, err := s.Expr(vdn.Value)
		if err != nil {
			return err
		}
		s.Scope.SetVar(vdn.Var, val)
		return nil
	case nodes.NdFuncDef:
		fdn := n.(*nodes.FuncDefNode)
		s.Scope.SetFunc(fdn.Name, &Funct{
			Args: fdn.Args,
			Ret:  fdn.Ret,
			Body: fdn.Body,
		})
		return nil
	case nodes.NdFuncExtern:
		fen := n.(*nodes.FuncExternNode)
		s.Scope.SetFunc(fen.Name, &Funct{
			Args: fen.Args,
			Ret:  fen.Ret,
			Body: []nodes.Node{},
		})
		return nil
	case nodes.NdIf:
		ifn := n.(*nodes.IfNode)
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
	case nodes.NdWhile:
		whn := n.(*nodes.WhileNode)
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
	case nodes.NdFuncCall:
		fcn := n.(*nodes.FuncCallNode)
		funct, ok := s.Scope.FindFunc(fcn.Func)
		if !ok {
			return fmt.Errorf("unknown function: %s", fcn.Func)
		}
		argv := map[string]*values.Value{}
		for i := 0; i < len(fcn.Args); i++ {
			val, err := s.Expr(fcn.Args[i])
			if err != nil {
				return err
			}
			argv[funct.Args[i].Name] = val
		}
		s.Scope = &Scope{
			parent: s.Scope,
			Vars:   argv,
			Funcs:  map[string]*Funct{},
		}
		s.Calls = append(s.Calls, fcn.Func)
		if funct.IsNative {
			_, err := funct.Native(s, argv)
			if err != nil {
				return err
			}
		} else {
			for _, n := range funct.Body {
				err := s.Stmt(n)
				if err != nil {
					return err
				}
			}
		}
		s.Calls = s.Calls[:len(s.Calls)-1]
		s.Scope = s.Scope.parent
		return nil
	}
	return fmt.Errorf("%s is not a statement", n)
}

func (s *Simulator) Expr(n nodes.Node) (*values.Value, error) {
	switch n.Kind() {
	case nodes.NdInteger:
		return values.NewValue(n.(*nodes.IntegerNode).Int), nil
	case nodes.NdBoolean:
		return values.NewValue(n.(*nodes.BooleanNode).Bool), nil
	case nodes.NdFloat:
		return values.NewValue(n.(*nodes.FloatNode).Float), nil
	case nodes.NdString:
		return values.NewValue(n.(*nodes.StringNode).Str), nil
	case nodes.NdChar:
		return values.NewValue(n.(*nodes.CharNode).Char), nil
	case nodes.NdFuncCall:
		fcn := n.(*nodes.FuncCallNode)
		funct, ok := s.Scope.FindFunc(fcn.Func)
		if !ok {
			return nil, fmt.Errorf("unknown function: %s", fcn.Func)
		}
		if funct.Ret.Equals(values.Nothing) {
			return nil, fmt.Errorf("function %s does not return a values.Value", fcn.Func)
		}
		argv := map[string]*values.Value{}
		for i := 0; i < len(fcn.Args); i++ {
			val, err := s.Expr(fcn.Args[i])
			if err != nil {
				return nil, err
			}
			argv[funct.Args[i].Name] = val
		}
		s.Scope = &Scope{
			parent: s.Scope,
			Vars:   argv,
			Funcs:  map[string]*Funct{},
		}
		s.Calls = append(s.Calls, fcn.Func)
		var val *values.Value = values.Default(funct.Ret)
		var err error
		if funct.IsNative {
			val, err = funct.Native(s, argv)
			s.Scope = s.Scope.parent
			return val, err
		}
		for _, n := range funct.Body {
			val, err = s.stmtInFunc(n)
			if err != nil {
				return nil, err
			}
			if val != nil {
				s.Scope = s.Scope.parent
				return val, err
			}
		}
		s.Calls = s.Calls[:len(s.Calls)-1]
		return nil, fmt.Errorf("function %s did not return", fcn.Func)
	case nodes.NdNewStruct:
		nsn := n.(*nodes.NewStructNode)
		var v *values.Value
		fields := map[string]*values.Value{}
		for i, f := range nsn.Args {
			v, err := s.Expr(f)
			if err != nil {
				return nil, err
			}
			fields[nsn.Type.Structure.Fields[i].Name] = v
		}
		v = &values.Value{nsn.Type, fields}
		return v, nil
	case nodes.NdSelector:
		sn := n.(*nodes.SelectorNode)
		p := sn.Path()
		var v *values.Value
		var ok bool
		v, ok = s.Scope.FindVar(p[0])
		if !ok {
			return nil, fmt.Errorf("variable %s not found", p[0])
		}
		if len(p) == 1 {
			return v, nil
		}
		for _, name := range p[1:] {
			if v != nil && v.Type.Prim != values.PStruct {
				return nil, fmt.Errorf("non-struct types do not have fields")
			}
			v, ok = v.Struct()[name]
			if !ok {
				return nil, fmt.Errorf("field %s not found", name)
			}
		}
		return v, nil
	}
	fmt.Printf("Simulator: not a value: %s\n", n.Kind())
	return nil, fmt.Errorf("%s node is not a value", n.Kind())
}

func (s *Simulator) Const(n nodes.Node) (*values.Value, error) {
	switch n.Kind() {
	case nodes.NdInteger, nodes.NdBoolean, nodes.NdFloat, nodes.NdString, nodes.NdChar:
		return s.Expr(n)
	default:
		return nil, fmt.Errorf("%s node is not constant", n.Kind())
	}
}

func (s *Simulator) block(b []nodes.Node) error {
	s.Scope = &Scope{
		parent: s.Scope,
		Vars:   map[string]*values.Value{},
		Funcs:  map[string]*Funct{},
	}
	for _, n := range b {
		if err := s.Stmt(n); err != nil {
			return err
		}
	}
	s.Scope = s.Scope.parent
	return nil
}

func (s *Simulator) stmtInFunc(n nodes.Node) (*values.Value, error) {
	switch n.Kind() {
	case nodes.NdVarDef, nodes.NdFuncDef, nodes.NdFuncExtern, nodes.NdFuncCall:
		return nil, s.Stmt(n)
	case nodes.NdIf:
		ifn := n.(*nodes.IfNode)
		val, err := s.Expr(ifn.Condition)
		if err != nil {
			return nil, err
		}
		if val.ToBool() {
			return s.blockInFunc(ifn.IfBlock)
		}
		for _, branch := range ifn.ElseIfNodes {
			val, err = s.Expr(branch.Condition)
			if err != nil {
				return nil, err
			}
			if val.ToBool() {
				return s.blockInFunc(branch.Block)
			}
		}
		return s.blockInFunc(ifn.ElseBlock)
	case nodes.NdWhile:
		whn := n.(*nodes.WhileNode)
		val, err := s.Expr(whn.Condition)
		if err != nil {
			return nil, err
		}
		if !val.ToBool() {
			return nil, nil
		}
		for {
			val, err = s.blockInFunc(whn.Body)
			if err != nil {
				return nil, err
			}
			if val != nil {
				return val, nil
			}
			val, err = s.Expr(whn.Condition)
			if err != nil {
				return nil, err
			}
			if !val.ToBool() {
				return nil, nil
			}
		}
	case nodes.NdReturn:
		return s.Expr(n.(*nodes.ReturnNode).Value)
	}
	return nil, fmt.Errorf("%s is not a statement", n)
}

func (s *Simulator) blockInFunc(b []nodes.Node) (*values.Value, error) {
	s.Scope = &Scope{
		parent: s.Scope,
		Vars:   map[string]*values.Value{},
		Funcs:  map[string]*Funct{},
	}
	for _, n := range b {
		v, err := s.stmtInFunc(n)
		if err != nil {
			return nil, err
		}
		if v != nil {
			return v, nil
		}
	}
	s.Scope = s.Scope.parent
	return nil, nil
}
