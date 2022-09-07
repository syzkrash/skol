package typecheck

import (
	"fmt"

	"github.com/qeaml/all/slices"
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/values/types"
)

func typeError(mn ast.MetaNode, want, got types.Type, format string, a ...any) *TypeError {
	return &TypeError{
		msg:   fmt.Sprintf(format, a...),
		Got:   got,
		Want:  want,
		Cause: mn,
	}
}

// Checker ensures type correctness of an AST.
type Checker struct {
	scope *scope
}

// NewChecker creates a blank Checker.
func NewChecker() *Checker {
	return &Checker{
		scope: &scope{
			parent: nil,
			vars:   make(map[string]types.Type),
			funcs:  make(map[string]funcproto),
		},
	}
}

// Check thoroughly inspects the provided AST for any typing-related errors
// that may have occured.
func (c *Checker) Check(tree ast.AST) (errs []*TypeError) {
	for _, v := range tree.Typedefs {
		c.scope.vars[v.Name] = v.Type
	}
	for _, v := range tree.Vars {
		t, err := c.typeOf(v.Value)
		if err != nil {
			errs = append(errs, err)
		} else {
			c.scope.vars[v.Name] = t
		}
	}
	// first loop to declare functions
	for _, f := range tree.Funcs {
		a := make([]types.Type, len(f.Args))
		for i, aa := range f.Args {
			a[i] = aa.Type
		}
		c.scope.funcs[f.Name] = basicFuncproto(a, f.Ret)
	}
	// second loop to typecheck function bodies with function type information
	for _, f := range tree.Funcs {
		args := make(map[string]types.Type)
		for _, a := range f.Args {
			args[a.Name] = a.Type
		}
		errs = append(errs, c.checkFunc(args, f.Ret, f.Body)...)
	}

	errs = slices.Filter(errs, func(e *TypeError) bool {
		return e != nil
	})
	return
}

func (c *Checker) checkNode(mn ast.MetaNode, ret types.Type) (errs []*TypeError) {
	n := mn.Node

	switch n.Kind() {
	// literals
	case ast.NStruct:
		nstruct := n.(ast.StructNode)
		for i, a := range nstruct.Args {
			at, err := c.typeOf(a)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			ft := nstruct.Type.Fields[i].Type
			if !ft.Equals(at) {
				errs = append(errs, typeError(a, ft, at,
					"incorrect struct argument type"))
			}
		}
	case ast.NArray:
		narray := n.(ast.ArrayNode)
		for _, e := range narray.Elems {
			et, err := c.typeOf(e)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if !narray.Type.Equals(et) {
				errs = append(errs, typeError(e, narray.Type, et,
					"incorrect array element type"))
			}
		}

	// control flow
	case ast.NIf:
		nif := n.(ast.IfNode)
		errs = append(errs, c.checkBranch(nif.Main, ret)...)
		for _, b := range nif.Other {
			errs = append(errs, c.checkBranch(b, ret)...)
		}
		errs = append(errs, c.checkBlock(nif.Else, ret)...)
	case ast.NWhile:
		nwhile := n.(ast.WhileNode)
		if err := c.checkCond(nwhile.Cond); err != nil {
			errs = append(errs, err)
		}
		errs = append(errs, c.checkBlock(nwhile.Block, ret)...)
	case ast.NReturn:
		nreturn := n.(ast.ReturnNode)
		rt, err := c.typeOf(nreturn.Value)
		if err != nil {
			errs = append(errs, err)
			return
		}
		if !ret.Equals(rt) {
			errs = append(errs, typeError(nreturn.Value, ret, rt,
				"incorrect type for return value"))
		}

	// definitions
	case ast.NVarSet:
		nvarset := n.(ast.VarSetNode)
		nvt, err := c.typeOf(nvarset.Value)
		if err != nil {
			errs = append(errs, err)
			return
		}
		if ovt, ok := c.scope.getVar(nvarset.Var); ok {
			if !ovt.Equals(nvt) {
				errs = append(errs, typeError(mn, ovt, nvt,
					"incorrect type for variable value"))
			}
		} else {
			c.scope.setVar(nvarset.Var, nvt)
		}
	case ast.NVarDef:
		nvardef := n.(ast.VarDefNode)
		c.scope.setVar(nvardef.Var, nvardef.Type)
	case ast.NVarSetTyped:
		nvarsettyped := n.(ast.VarSetTypedNode)
		nvt, err := c.typeOf(nvarsettyped.Value)
		if err != nil {
			errs = append(errs, err)
			return
		}
		ovt, ok := c.scope.getVar(nvarsettyped.Var)
		if ok {
			if !ovt.Equals(nvarsettyped.Type) {
				errs = append(errs, typeError(nvarsettyped.Value, ovt, nvarsettyped.Type,
					"cannot change type of variable"))
			}
			if !ovt.Equals(nvt) {
				errs = append(errs, typeError(nvarsettyped.Value, ovt, nvt,
					"incorrect type for variable value"))
			}
			if !ovt.Equals(nvarsettyped.Type) {
				errs = append(errs, typeError(nvarsettyped.Value, ovt, nvt,
					"incorrect type for variable value"))
			}
		}
	case ast.NFuncDef:
		nfuncdef := n.(ast.FuncDefNode)
		args := make(map[string]types.Type)
		argl := make([]types.Type, len(nfuncdef.Proto))
		for i, a := range nfuncdef.Proto {
			args[a.Name] = a.Type
			argl[i] = a.Type
		}
		errs = append(errs, c.checkFunc(args, nfuncdef.Ret, nfuncdef.Body)...)
		c.scope.funcs[nfuncdef.Name] = basicFuncproto(argl, nfuncdef.Ret)
	case ast.NFuncExtern:
		nfuncextern := n.(ast.FuncExternNode)
		args := make([]types.Type, len(nfuncextern.Proto))
		for i, a := range nfuncextern.Proto {
			args[i] = a.Type
		}
		c.scope.funcs[nfuncextern.Alias] = basicFuncproto(args, nfuncextern.Ret)

	// others
	case ast.NFuncCall:
		nfunccall := n.(ast.FuncCallNode)
		f, ok := c.scope.getFunc(nfunccall.Func)
		if !ok {
			errs = append(errs, typeError(mn, nil, nil,
				"unknown function"))
			return
		}
		args := make([]types.Type, len(nfunccall.Args))
		for i, a := range nfunccall.Args {
			t, err := c.typeOf(a)
			if err != nil {
				errs = append(errs, err)
				return
			}
			args[i] = t
		}
		if _, err := f(mn, args); err != nil {
			errs = append(errs, err)
		}
	}
	return
}

// checkFunc ensures type correctness given the function arguments' types and
// the return type.
func (c *Checker) checkFunc(args map[string]types.Type, ret types.Type, body ast.Block) (errs []*TypeError) {
	c.scope = c.scope.sub()
	for n, t := range args {
		c.scope.vars[n] = t
	}
	errs = c.checkBlock(body, ret)
	c.scope = c.scope.parent
	return
}

func (c *Checker) checkBranch(b ast.Branch, ret types.Type) (errs []*TypeError) {
	if err := c.checkCond(b.Cond); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, c.checkBlock(b.Block, ret)...)
	return
}

func (c *Checker) checkCond(mn ast.MetaNode) (err *TypeError) {
	ct, err := c.typeOf(mn)
	if err != nil {
		return
	}
	if !types.Bool.Equals(ct) {
		err = typeError(mn, types.Bool, ct,
			"incorrect type for condition")
	}
	return
}

func (c *Checker) checkBlock(b ast.Block, ret types.Type) (errs []*TypeError) {
	for _, mn := range b {
		errs = append(errs, c.checkNode(mn, ret)...)
	}
	return
}

// typeOf determines the type of any abstract AST node.
func (c *Checker) typeOf(mn ast.MetaNode) (t types.Type, err *TypeError) {
	n := mn.Node

	switch n.Kind() {
	// literals
	case ast.NBool:
		t = types.Bool
	case ast.NChar:
		t = types.Char
	case ast.NInt:
		t = types.Int
	case ast.NFloat:
		t = types.Float
	case ast.NString:
		t = types.String
	case ast.NStruct:
		t = n.(ast.StructNode).Type
	case ast.NArray:
		t = n.(ast.ArrayNode).Type

	// others
	case ast.NFuncCall:
		nfunccall := n.(ast.FuncCallNode)
		f, ok_ := c.scope.getFunc(nfunccall.Func)
		if !ok_ {
			err = typeError(mn, nil, nil,
				"unknown function")
			return
		}
		args := make([]types.Type, len(nfunccall.Args))
		for i, a := range nfunccall.Args {
			at, err_ := c.typeOf(a)
			if err_ != nil {
				err = err_
				return
			}
			args[i] = at
		}
		if ret, err_ := f(mn, args); err_ != nil {
			err = err_
		} else {
			t = ret
		}

	default:
		err = typeError(mn, nil, nil, "typeOf() unimplemented: %s", n.Kind())
	}
	return
}

// basicFuncproto generates a funcproto for a function that accepts a certain
// combination of arguments and always returns the same type.
func basicFuncproto(exp []types.Type, ret types.Type) funcproto {
	return func(mn ast.MetaNode, got []types.Type) (r types.Type, err *TypeError) {
		r = ret
		if len(got) < len(exp) {
			err = typeError(mn, nil, nil,
				"not enough arguments")
		}
		for i, ea := range exp {
			if !ea.Equals(got[i]) {
				err = typeError(mn, ea, got[i],
					"incorrect function argument type")
				return
			}
		}
		return
	}
}
