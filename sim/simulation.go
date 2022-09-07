package sim

import (
	"fmt"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

// Simulator performs basic interpretationa and evaluation on abstract AST
// nodes. This structure contains the current simulation scope and call stack.
type Simulator struct {
	Scope *Scope
	Calls []*Call
}

// Error creates a [*SimError] with the current call stack pointing at the
// given node.
func (s *Simulator) Error(msg string, n ast.MetaNode) error {
	return &SimError{
		msg:   msg,
		Cause: n,
		Calls: s.Calls,
	}
}

// Errorf is the same as [Simulator.Error], except with a format.
func (s *Simulator) Errorf(n ast.MetaNode, format string, a ...any) error {
	return s.Error(fmt.Sprintf(format, a...), n)
}

// NewSimulator creates and prepares a simulator instance.
func NewSimulator() *Simulator {
	return &Simulator{
		Scope: &Scope{
			parent: nil,
			Vars:   map[string]*values.Value{},
			Funcs:  DefaultFuncs,
		},
		Calls: []*Call{},
	}
}

// pushCall appends the given call to the current call stack.
func (s *Simulator) pushCall(name string, where lexer.Position) {
	s.Calls = append(s.Calls, &Call{false, name, where})
}

// popCall pops the last call off the current call stack.
func (s *Simulator) popCall() {
	s.Calls = s.Calls[:len(s.Calls)-1]
}

// Stmt simulates the given abstact AST node.
func (s *Simulator) Stmt(mn ast.MetaNode) error {
	n := mn.Node
	switch n.Kind() {
	case ast.NVarSet:
		vdn := n.(ast.VarSetNode)
		val, err := s.Expr(vdn.Value)
		if err != nil {
			return err
		}
		s.Scope.SetVar(vdn.Var, val)
		return nil
	case ast.NFuncDef:
		fdn := n.(ast.FuncDefNode)
		args := make([]values.FuncArg, len(fdn.Proto))
		for i, a := range fdn.Proto {
			args[i] = values.FuncArg{Name: a.Name, Type: a.Type}
		}
		s.Scope.SetFunc(fdn.Name, &Funct{
			Args: args,
			Ret:  fdn.Ret,
			Body: fdn.Body,
		})
		return nil
	case ast.NFuncExtern:
		fen := n.(ast.FuncExternNode)
		args := make([]values.FuncArg, len(fen.Proto))
		for i, a := range fen.Proto {
			args[i] = values.FuncArg{Name: a.Name, Type: a.Type}
		}
		s.Scope.SetFunc(fen.Name, &Funct{
			Args: args,
			Ret:  fen.Ret,
			Body: ast.Block{},
		})
		return nil
	case ast.NIf:
		ifn := n.(ast.IfNode)
		val, err := s.Expr(ifn.Main.Cond)
		if err != nil {
			return err
		}
		if val.ToBool() {
			return s.block(ifn.Main.Block)
		}
		for _, branch := range ifn.Other {
			val, err = s.Expr(branch.Cond)
			if err != nil {
				return err
			}
			if val.ToBool() {
				return s.block(branch.Block)
			}
		}
		return s.block(ifn.Else)
	case ast.NWhile:
		whn := n.(ast.WhileNode)
		val, err := s.Expr(whn.Cond)
		if err != nil {
			return err
		}
		if !val.ToBool() {
			return nil
		}
		for {
			err = s.block(whn.Block)
			if err != nil {
				return err
			}
			val, err = s.Expr(whn.Cond)
			if err != nil {
				return err
			}
			if !val.ToBool() {
				return nil
			}
		}
	case ast.NFuncCall:
		fcn := n.(ast.FuncCallNode)
		funct, ok := s.Scope.FindFunc(fcn.Func)
		if !ok {
			return s.Errorf(mn, "unknown function: %s", fcn.Func)
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
		s.pushCall(fcn.Func, mn.Where)
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
		s.popCall()
		s.Scope = s.Scope.parent
		return nil
	case ast.NStruct:
		// do nothing -- types are handled entirely inside the parser
		return nil
	}
	return s.Errorf(mn, "%s is not a statement", n)
}

// selector evaluates the value of a selector in the current scope.
func (s *Simulator) selector(mn ast.MetaNode, sel ast.Selector) (*values.Value, error) {
	p := sel.Path()
	var v *values.Value
	var ok bool
	v, ok = s.Scope.FindVar(p[0].Name)
	if !ok {
		return nil, s.Errorf(mn, "variable %s not found", p[0].Name)
	}
	if len(p) == 1 {
		return v, nil
	}
	for _, e := range p[1:] {
		// again, the selector resolve alogrithm
		// see /parser/types.go, line 137
		if e.Cast != nil {
			v.Type = e.Cast
			continue
		}
		if e.Name != "" {
			v = v.Data.(map[string]*values.Value)[e.Name]
			continue
		}
		var idx int
		if e.IdxS != nil {
			p := e.IdxS.Path()
			vn := p[len(p)-1].Name
			val, ok := s.Scope.FindVar(vn)
			if !ok {
				return nil, s.Errorf(mn, "unknown variable: %s", vn)
			}
			if !types.Int.Equals(val.Type) {
				return nil, s.Errorf(mn, "can only index with integers")
			}
			idx = int(val.Int())
		} else {
			idx = e.IdxC
		}
		arraytype := v.Type.(types.ArrayType)
		arraydata := v.Data.([]*values.Value)
		resulttype := types.MakeStruct(arraytype.Element.String()+"Result",
			"ok", types.Bool,
			"val", arraytype.Element)
		if idx >= len(arraydata) {
			v = &values.Value{
				Type: resulttype,
				Data: map[string]*values.Value{
					"ok": values.NewValue(false),
				},
			}
		} else {
			v = &values.Value{
				Type: resulttype,
				Data: map[string]*values.Value{
					"ok":  values.NewValue(true),
					"val": arraydata[idx],
				},
			}
		}
	}
	return v, nil
}

// Expr evalues the given abstact AST node in the current scope.
func (s *Simulator) Expr(mn ast.MetaNode) (*values.Value, error) {
	n := mn.Node
	switch n.Kind() {
	case ast.NInt:
		return values.NewValue(n.(ast.IntNode).Value), nil
	case ast.NBool:
		return values.NewValue(n.(ast.BoolNode).Value), nil
	case ast.NFloat:
		return values.NewValue(n.(ast.FloatNode).Value), nil
	case ast.NString:
		return values.NewValue(n.(ast.StringNode).Value), nil
	case ast.NChar:
		return values.NewValue(n.(ast.CharNode).Value), nil
	case ast.NFuncCall:
		fcn := n.(ast.FuncCallNode)
		funct, ok := s.Scope.FindFunc(fcn.Func)
		if !ok {
			return nil, s.Errorf(mn, "unknown function: %s", fcn.Func)
		}
		if funct.Ret.Equals(types.Nothing) {
			return nil, s.Errorf(mn, "function %s does not return a values.Value", fcn.Func)
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
		s.pushCall(fcn.Func, mn.Where)
		var val *values.Value = values.Default(funct.Ret)
		var err error
		if funct.IsNative {
			val, err = funct.Native(s, argv)
			s.Scope = s.Scope.parent
			s.popCall()
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
		s.popCall()
		return nil, s.Errorf(mn, "function %s did not return", fcn.Func)
	case ast.NStruct:
		nsn := n.(ast.StructNode)
		var v *values.Value
		fields := map[string]*values.Value{}
		for i, f := range nsn.Args {
			v, err := s.Expr(f)
			if err != nil {
				return nil, err
			}
			fields[nsn.Type.Fields[i].Name] = v
		}
		v = &values.Value{Type: nsn.Type, Data: fields}
		return v, nil
	default:
		if sel, ok := n.(ast.Selector); ok {
			return s.selector(mn, sel)
		}
	}
	return nil, s.Errorf(mn, "%s node is not a value", n.Kind())
}

// Const evaluates a constant value from the given abstract AST node, limited
// to only literals.
func (s *Simulator) Const(mn ast.MetaNode) (*values.Value, error) {
	n := mn.Node
	switch n.Kind() {
	case ast.NInt, ast.NBool, ast.NFloat, ast.NString, ast.NChar:
		return s.Expr(mn)
	default:
		return nil, s.Errorf(mn, "%s node is not constant", n.Kind())
	}
}

// block simulates a block of statements in a new child scope.
func (s *Simulator) block(b ast.Block) error {
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

// stmtInFunc simulates a statement as if inside a function body.
func (s *Simulator) stmtInFunc(mn ast.MetaNode) (*values.Value, error) {
	n := mn.Node
	switch n.Kind() {
	case ast.NVarDef, ast.NFuncDef, ast.NFuncExtern, ast.NFuncCall:
		return nil, s.Stmt(mn)
	case ast.NIf:
		ifn := n.(ast.IfNode)
		val, err := s.Expr(ifn.Main.Cond)
		if err != nil {
			return nil, err
		}
		if val.ToBool() {
			return s.blockInFunc(ifn.Main.Block)
		}
		for _, branch := range ifn.Other {
			val, err = s.Expr(branch.Cond)
			if err != nil {
				return nil, err
			}
			if val.ToBool() {
				return s.blockInFunc(branch.Block)
			}
		}
		return s.blockInFunc(ifn.Else)
	case ast.NWhile:
		whn := n.(ast.WhileNode)
		val, err := s.Expr(whn.Cond)
		if err != nil {
			return nil, err
		}
		if !val.ToBool() {
			return nil, nil
		}
		for {
			val, err = s.blockInFunc(whn.Block)
			if err != nil {
				return nil, err
			}
			if val != nil {
				return val, nil
			}
			val, err = s.Expr(whn.Cond)
			if err != nil {
				return nil, err
			}
			if !val.ToBool() {
				return nil, nil
			}
		}
	case ast.NReturn:
		return s.Expr(n.(ast.ReturnNode).Value)
	}
	return nil, s.Errorf(mn, "%s is not a statement", n)
}

// blockInFunc simulates a block of statements as if inside a function body in
// a new child scope.
func (s *Simulator) blockInFunc(b ast.Block) (*values.Value, error) {
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
