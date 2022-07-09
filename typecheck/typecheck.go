package typecheck

import (
	"fmt"
	"io"

	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type Typechecker struct {
	Parser *parser.Parser
}

func NewTypechecker(src io.RuneScanner, fn, eng string) *Typechecker {
	return &Typechecker{
		Parser: parser.NewParser(fn, src, eng),
	}
}

func typeError(n nodes.Node, want, got *values.Type, format string, a ...any) error {
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
		for i, param := range fdef.Args {
			atype, err := t.Parser.TypeOf(fcn.Args[i])
			if err != nil {
				return err
			}
			if !param.Type.Equals(atype) {
				return typeError(fcn.Args[i], param.Type, atype,
					"wrong argument type in function call")
			}
		}

	case nodes.NdFuncDef:
		fdn := n.(*nodes.FuncDefNode)
		t.Parser.Scope = &parser.Scope{
			Parent: t.Parser.Scope,
			Funcs:  map[string]*parser.Function{},
			Vars:   map[string]*nodes.VarDefNode{},
			Consts: map[string]*values.Value{},
			Types:  map[string]*values.Type{},
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
		if !values.Bool.Equals(ctype) {
			return typeError(cn.Condition, values.Bool, ctype,
				"conditional must be a boolean")
		}

	case nodes.NdWhile:
		cn := n.(*nodes.WhileNode)
		ctype, err := t.Parser.TypeOf(cn.Condition)
		if err != nil {
			return err
		}
		if !values.Bool.Equals(ctype) {
			return typeError(cn.Condition, values.Bool, ctype,
				"conditional must be a boolean")
		}

	case nodes.NdNewStruct:
		nsn := n.(*nodes.NewStructNode)
		for i, f := range nsn.Type.Structure.Fields {
			atype, err := t.Parser.TypeOf(nsn.Args[i])
			if err != nil {
				return err
			}
			if !f.Type.Equals(atype) {
				return typeError(nsn.Args[i], f.Type, atype,
					"wrong field type in structure literal")
			}
		}
	}
	return nil
}

func (t *Typechecker) checkNodeInFunc(n nodes.Node, ret *values.Type) (err error) {
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
	default:
		return t.checkNode(n)
	}
}

func (t *Typechecker) Next() (n nodes.Node, err error) {
	n, err = t.Parser.Next()
	if err != nil {
		return
	}
	err = t.checkNode(n)
	return
}
