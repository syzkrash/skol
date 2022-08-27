package typecheck

import (
	"fmt"
	"io"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

type Typechecker struct {
	Parser *parser.Parser
}

func NewTypechecker(src io.RuneScanner, fn, eng string) *Typechecker {
	return &Typechecker{
		Parser: parser.NewParser(fn, src, eng),
	}
}

func typeError(mn ast.MetaNode, want, got types.Type, format string, a ...any) error {
	return &TypeError{
		msg:   fmt.Sprintf(format, a...),
		Got:   got,
		Want:  want,
		Cause: mn,
	}
}

func (t *Typechecker) checkSelector(mn ast.MetaNode, s ast.Selector) (err error) {
	path := s.Path()
	root, _ := t.Parser.Scope.FindVar(path[0].Name)
	typeNow, err := t.Parser.TypeOf(root)
	if err != nil {
		return
	}
	for _, e := range path[1:] {
		// again, the selector resolve algorithm
		// see /parser/types.go, line 137
		if e.Cast != nil {
			if !typeNow.Equals(e.Cast) {
				return typeError(mn, typeNow, e.Cast,
					"cast to incompatible type")
			}
			typeNow = e.Cast
			continue
		}
		if e.Name != "" {
			if typeNow.Prim() != types.PStruct {
				return typeError(mn, typeNow, nil,
					"can only access fields on structures (acessing '%s' on %s)",
					e.Name, typeNow)
			}
			st := typeNow.(types.StructType)
			ok := false
			for _, f := range st.Fields {
				if f.Name == e.Name {
					ok = true
					typeNow = f.Type
				}
			}
			if !ok {
				return typeError(mn, typeNow, nil,
					"unknown field: '%s'", e.Name)
			}
			continue
		}
		if typeNow.Prim() != types.PArray {
			return typeError(mn, typeNow, nil,
				"can only index arrays")
		}
		at := typeNow.(types.ArrayType)
		typeNow = t.Parser.EnsureResultType(at.Element)
	}
	return
}

func (t *Typechecker) checkNode(mn ast.MetaNode) (err error) {
	n := mn.Node
	switch n.Kind() {
	case ast.NFuncCall:
		fcn := n.(ast.FuncCallNode)
		fdef, _ := t.Parser.Scope.FindFunc(fcn.Func)
		paramTypes := []types.Type{}
		for i, param := range fdef.Args {
			atype, err := t.Parser.TypeOf(fcn.Args[i].Node)
			if err != nil {
				return err
			}
			paramTypes = append(paramTypes, param.Type)
			if !param.Type.Equals(atype) {
				return typeError(fcn.Args[i], param.Type, atype,
					"wrong argument type in function call")
			}
		}

		if fcn.Func == "eq" {
			if !paramTypes[1].Equals(paramTypes[0]) {
				return typeError(fcn.Args[1], paramTypes[0], paramTypes[1],
					"eq! arguments must be the same type")
			}
		}

	case ast.NFuncDef:
		fdn := n.(ast.FuncDefNode)
		if fdn.Ret.Equals(types.Any) {
			return typeError(mn, types.Nothing, types.Nothing,
				"functions cannot return Any")
		}
		t.Parser.Scope = &parser.Scope{
			Parent: t.Parser.Scope,
			Funcs:  map[string]*values.Function{},
			Vars:   map[string]ast.Node{},
			Consts: map[string]ast.Node{},
			Types:  map[string]types.Type{},
		}
		for _, a := range fdn.Proto {
			an, ok := t.Parser.NodeOf(a.Type)
			if !ok {
				err = common.Error(mn, "cannot typecheck value of type %s", a.Type)
				return
			}
			t.Parser.Scope.SetVar(a.Name, an)
		}
		for _, bn := range fdn.Body {
			err = t.checkNodeInFunc(bn, fdn.Ret)
			if err != nil {
				return
			}
		}
		t.Parser.Scope = t.Parser.Scope.Parent

	case ast.NIf:
		cn := n.(ast.IfNode)
		ctype, err := t.Parser.TypeOf(cn.Main.Cond.Node)
		if err != nil {
			return err
		}
		if !types.Bool.Equals(ctype) {
			return typeError(cn.Main.Cond, types.Bool, ctype,
				"conditional must be a boolean")
		}
		for _, bn := range cn.Main.Block {
			if err = t.checkNode(bn); err != nil {
				return err
			}
		}
		for _, bn := range cn.Else {
			if err = t.checkNode(bn); err != nil {
				return err
			}
		}
		for _, branch := range cn.Other {
			for _, bn := range branch.Block {
				if err = t.checkNode(bn); err != nil {
					return err
				}
			}
		}

	case ast.NWhile:
		cn := n.(ast.WhileNode)
		ctype, err := t.Parser.TypeOf(cn.Cond.Node)
		if err != nil {
			return err
		}
		if !types.Bool.Equals(ctype) {
			return typeError(cn.Cond, types.Bool, ctype,
				"conditional must be a boolean")
		}
		for _, bn := range cn.Block {
			if err = t.checkNode(bn); err != nil {
				return err
			}
		}

	case ast.NStruct:
		nsn := n.(ast.StructNode)
		for i, f := range nsn.Type.Fields {
			atype, err := t.Parser.TypeOf(nsn.Args[i].Node)
			if err != nil {
				return err
			}
			if !f.Type.Equals(atype) {
				return typeError(nsn.Args[i], f.Type, atype,
					"wrong field type in structure literal")
			}
		}

	case ast.NVarSet:
		vdn := n.(ast.VarSetNode)
		switch vdn.Value.Node.Kind() {
		case ast.NArray:
			an := vdn.Value.Node.(ast.ArrayNode)
			for _, e := range an.Elems {
				et, err := t.Parser.TypeOf(e.Node)
				if err != nil {
					return err
				}
				if !an.Type.Equals(et) {
					return typeError(vdn.Value, an.Type, et,
						"value of incompatible type in array")
				}
			}
		}
		if sel, ok := vdn.Value.Node.(ast.Selector); ok {
			return t.checkSelector(vdn.Value, sel)
		}
	}
	return nil
}

func (t *Typechecker) checkNodeInFunc(mn ast.MetaNode, ret types.Type) (err error) {
	n := mn.Node
	switch n.Kind() {
	case ast.NReturn:
		rn := n.(ast.ReturnNode)
		vtype, err := t.Parser.TypeOf(rn.Value.Node)
		if err != nil {
			return err
		}
		if !ret.Equals(vtype) {
			return typeError(rn.Value, ret, vtype,
				"incorrect or inconsistent function return type")
		}
		return nil

	case ast.NIf:
		cn := n.(ast.IfNode)
		ctype, err := t.Parser.TypeOf(cn.Main.Cond.Node)
		if err != nil {
			return err
		}
		if !types.Bool.Equals(ctype) {
			return typeError(cn.Main.Cond, types.Bool, ctype,
				"conditional must be a boolean")
		}
		for _, bn := range cn.Main.Block {
			if err = t.checkNodeInFunc(bn, ret); err != nil {
				return err
			}
		}
		for _, bn := range cn.Else {
			if err = t.checkNodeInFunc(bn, ret); err != nil {
				return err
			}
		}
		for _, branch := range cn.Other {
			for _, bn := range branch.Block {
				if err = t.checkNodeInFunc(bn, ret); err != nil {
					return err
				}
			}
		}

	case ast.NWhile:
		cn := n.(ast.WhileNode)
		ctype, err := t.Parser.TypeOf(cn.Cond.Node)
		if err != nil {
			return err
		}
		if !types.Bool.Equals(ctype) {
			return typeError(cn.Cond, types.Bool, ctype,
				"conditional must be a boolean")
		}
		for _, bn := range cn.Block {
			if err = t.checkNodeInFunc(bn, ret); err != nil {
				return err
			}
		}

	default:
		return t.checkNode(mn)
	}

	return nil
}

func (t *Typechecker) Next() (mn ast.MetaNode, err error) {
	mn, err = t.Parser.Next()
	if err != nil {
		return
	}
	n := mn.Node
	err = t.checkNode(mn)
	if err == nil && n.Kind() == ast.NFuncCall {
		fcn := n.(ast.FuncCallNode)
		// we can check for the skol function since it'd be typechecked by now
		if fcn.Func == "skol" {
			verVal, err := t.Parser.Sim.Expr(fcn.Args[1])
			if err != nil {
				return mn, err
			}
			engVal, err := t.Parser.Sim.Expr(fcn.Args[0])
			if err != nil {
				return mn, err
			}
			if verVal.Data.(float32) > common.VersionF {
				return mn, common.Error(mn,
					"this script requires skol %.1f or above",
					verVal.Data.(float32))
			}
			if engVal.Data.(string) != t.Parser.Engine {
				return mn, common.Error(mn,
					"this script required the %s engine",
					engVal.Data.(string))
			}
			return t.Next()
		}
	}
	return
}
