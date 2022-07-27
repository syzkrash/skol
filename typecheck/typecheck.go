package typecheck

import (
	"fmt"
	"io"

	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/nodes"
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

func typeError(n nodes.Node, want, got types.Type, format string, a ...any) error {
	return &TypeError{
		msg:  fmt.Sprintf(format, a...),
		Got:  got,
		Want: want,
		Node: n,
	}
}

func (t *Typechecker) checkNode(n nodes.Node) (err error) {
	switch n.Kind() {
	case nodes.NdFuncCall:
		fcn := n.(*nodes.FuncCallNode)
		fdef, _ := t.Parser.Scope.FindFunc(fcn.Func)
		paramTypes := []types.Type{}
		for i, param := range fdef.Args {
			atype, err := t.Parser.TypeOf(fcn.Args[i])
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

	case nodes.NdFuncDef:
		fdn := n.(*nodes.FuncDefNode)
		if fdn.Ret.Equals(types.Any) {
			return typeError(n, types.Nothing, types.Nothing,
				"functions cannot return Any")
		}
		t.Parser.Scope = &parser.Scope{
			Parent: t.Parser.Scope,
			Funcs:  map[string]*values.Function{},
			Vars:   map[string]*nodes.VarDefNode{},
			Consts: map[string]*values.Value{},
			Types:  map[string]types.Type{},
		}
		for _, a := range fdn.Args {
			t.Parser.Scope.SetVar(a.Name, &nodes.VarDefNode{
				Var:     a.Name,
				VarType: a.Type,
			})
		}
		for _, bn := range fdn.Body {
			err = t.checkNodeInFunc(bn, fdn.Ret)
			if err != nil {
				return err
			}
		}
		t.Parser.Scope = t.Parser.Scope.Parent

	case nodes.NdIf:
		cn := n.(*nodes.IfNode)
		ctype, err := t.Parser.TypeOf(cn.Condition)
		if err != nil {
			return err
		}
		if !types.Bool.Equals(ctype) {
			return typeError(cn.Condition, types.Bool, ctype,
				"conditional must be a boolean")
		}
		for _, bn := range cn.IfBlock {
			if err = t.checkNode(bn); err != nil {
				return err
			}
		}
		for _, bn := range cn.ElseBlock {
			if err = t.checkNode(bn); err != nil {
				return err
			}
		}
		for _, branch := range cn.ElseIfNodes {
			for _, bn := range branch.Block {
				if err = t.checkNode(bn); err != nil {
					return err
				}
			}
		}

	case nodes.NdWhile:
		cn := n.(*nodes.WhileNode)
		ctype, err := t.Parser.TypeOf(cn.Condition)
		if err != nil {
			return err
		}
		if !types.Bool.Equals(ctype) {
			return typeError(cn.Condition, types.Bool, ctype,
				"conditional must be a boolean")
		}
		for _, bn := range cn.Body {
			if err = t.checkNode(bn); err != nil {
				return err
			}
		}

	case nodes.NdNewStruct:
		nsn := n.(*nodes.NewStructNode)
		for i, f := range nsn.Type.(types.StructType).Fields {
			atype, err := t.Parser.TypeOf(nsn.Args[i])
			if err != nil {
				return err
			}
			if !f.Type.Equals(atype) {
				return typeError(nsn.Args[i], f.Type, atype,
					"wrong field type in structure literal")
			}
		}

	case nodes.NdVarDef:
		vdn := n.(*nodes.VarDefNode)
		if vdn.VarType.Equals(types.Any) {
			return typeError(n, types.Nothing, types.Nothing,
				"variables cannot be of type Any")
		}
		switch vdn.Value.Kind() {
		case nodes.NdTypecast:
			tn := vdn.Value.(*nodes.TypecastNode)
			otype, err := t.Parser.TypeOf(tn.Value)
			if err != nil {
				return err
			}
			if !otype.Equals(tn.Target) {
				return typeError(tn, otype, tn.Target,
					"cannot cast to incompatible type")
			}
		}
	}
	return nil
}

func (t *Typechecker) checkNodeInFunc(n nodes.Node, ret types.Type) (err error) {
	switch n.Kind() {
	case nodes.NdReturn:
		rn := n.(*nodes.ReturnNode)
		vtype, err := t.Parser.TypeOf(rn.Value)
		if err != nil {
			return err
		}
		if !ret.Equals(vtype) {
			return typeError(rn.Value, ret, vtype,
				"incorrect or inconsistent function return type")
		}
		return nil

	case nodes.NdIf:
		cn := n.(*nodes.IfNode)
		ctype, err := t.Parser.TypeOf(cn.Condition)
		if err != nil {
			return err
		}
		if !types.Bool.Equals(ctype) {
			return typeError(cn.Condition, types.Bool, ctype,
				"conditional must be a boolean")
		}
		for _, bn := range cn.IfBlock {
			if err = t.checkNodeInFunc(bn, ret); err != nil {
				return err
			}
		}
		for _, bn := range cn.ElseBlock {
			if err = t.checkNodeInFunc(bn, ret); err != nil {
				return err
			}
		}
		for _, branch := range cn.ElseIfNodes {
			for _, bn := range branch.Block {
				if err = t.checkNodeInFunc(bn, ret); err != nil {
					return err
				}
			}
		}

	case nodes.NdWhile:
		cn := n.(*nodes.WhileNode)
		ctype, err := t.Parser.TypeOf(cn.Condition)
		if err != nil {
			return err
		}
		if !types.Bool.Equals(ctype) {
			return typeError(cn.Condition, types.Bool, ctype,
				"conditional must be a boolean")
		}
		for _, bn := range cn.Body {
			if err = t.checkNodeInFunc(bn, ret); err != nil {
				return err
			}
		}

	default:
		return t.checkNode(n)
	}

	return nil
}

func (t *Typechecker) Next() (n nodes.Node, err error) {
	n, err = t.Parser.Next()
	if err != nil {
		return
	}
	err = t.checkNode(n)
	if err == nil && n.Kind() == nodes.NdFuncCall {
		fcn := n.(*nodes.FuncCallNode)
		// we can check for the skol function since it'd be typechecked by now
		if fcn.Func == "skol" {
			verVal, err := t.Parser.Sim.Expr(fcn.Args[1])
			if err != nil {
				return n, err
			}
			engVal, err := t.Parser.Sim.Expr(fcn.Args[0])
			if err != nil {
				return n, err
			}
			if verVal.Data.(float32) > common.VersionF {
				return n, common.Error(n,
					"this script requires skol %.1f or above",
					verVal.Data.(float32))
			}
			if engVal.Data.(string) != t.Parser.Engine {
				return n, common.Error(n,
					"this script required the %s engine",
					engVal.Data.(string))
			}
			return t.Next()
		}
	}
	return
}
