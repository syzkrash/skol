package typecheck

import (
	"github.com/qeaml/all/slices"
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/parser/values/types"
)

func typeMismatch(mn ast.MetaNode, want, got types.Type) *pe.PrettyError {
	return pe.New(pe.ETypeMismatch).Section("Wanted type", "%s", want).Section("Got type", "%s", got).Section("Caused by", "`%s` node at %s", mn.Node.Kind(), mn.Where)
}

func varRetype(mn ast.MetaNode, old, new types.Type) *pe.PrettyError {
	return pe.New(pe.EVarTypeChanged).Section("Original type", "%s", old).Section("New type", "%s", new).Section("Caused by", "`%s` node at %s", mn.Node.Kind(), mn.Where)
}

func nodeErr(c pe.ErrorCode, mn ast.MetaNode) *pe.PrettyError {
	return pe.New(c).Section("Caused by", "`%s` node at %s", mn.Node.Kind(), mn.Where)
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
			funcs:  defaultFuncs,
		},
	}
}

// Check thoroughly inspects the provided AST for any typing-related errors
// that may have occured.
func (c *Checker) Check(tree ast.AST) (errs []*pe.PrettyError) {
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

	errs = slices.Filter(errs, func(e *pe.PrettyError) bool {
		return e != nil
	})
	return
}

func (c *Checker) checkNode(mn ast.MetaNode, ret types.Type) (errs []*pe.PrettyError) {
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
				errs = append(errs, typeMismatch(a, ft, at))
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
				errs = append(errs, typeMismatch(e, narray.Type, et))
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
			errs = append(errs, typeMismatch(nreturn.Value, ret, rt))
		}

	// definitions
	case ast.NVarSet:
		nvarset := n.(ast.VarSetNode)
		errs = append(errs, c.checkNode(nvarset.Value, ret)...)
		nvt, err := c.typeOf(nvarset.Value)
		if err != nil {
			errs = append(errs, err)
			return
		}
		if ovt, ok := c.scope.getVar(nvarset.Var); ok {
			if !ovt.Equals(nvt) {
				errs = append(errs, typeMismatch(mn, ovt, nvt))
			}
		} else {
			c.scope.setVar(nvarset.Var, nvt)
		}
	case ast.NVarDef:
		nvardef := n.(ast.VarDefNode)
		c.scope.setVar(nvardef.Var, nvardef.Type)
	case ast.NVarSetTyped:
		nvarsettyped := n.(ast.VarSetTypedNode)
		errs = append(errs, c.checkNode(nvarsettyped.Value, ret)...)
		nvt, err := c.typeOf(nvarsettyped.Value)
		if err != nil {
			errs = append(errs, err)
			return
		}
		ovt, ok := c.scope.getVar(nvarsettyped.Var)
		if ok {
			if !ovt.Equals(nvarsettyped.Type) {
				errs = append(errs, varRetype(nvarsettyped.Value, ovt, nvarsettyped.Type))
			}
			if !ovt.Equals(nvt) {
				errs = append(errs, typeMismatch(nvarsettyped.Value, ovt, nvt))
			}
			if !ovt.Equals(nvarsettyped.Type) {
				errs = append(errs, typeMismatch(nvarsettyped.Value, ovt, nvt))
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
			errs = append(errs, typeMismatch(mn, nil, nil))
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
func (c *Checker) checkFunc(args map[string]types.Type, ret types.Type, body ast.Block) (errs []*pe.PrettyError) {
	c.scope = c.scope.sub()
	for n, t := range args {
		c.scope.vars[n] = t
	}
	errs = c.checkBlock(body, ret)
	c.scope = c.scope.parent
	return
}

func (c *Checker) checkBranch(b ast.Branch, ret types.Type) (errs []*pe.PrettyError) {
	if err := c.checkCond(b.Cond); err != nil {
		errs = append(errs, err)
	}
	errs = append(errs, c.checkBlock(b.Block, ret)...)
	return
}

func (c *Checker) checkCond(mn ast.MetaNode) (err *pe.PrettyError) {
	ct, err := c.typeOf(mn)
	if err != nil {
		return
	}
	if !types.Bool.Equals(ct) {
		err = typeMismatch(mn, types.Bool, ct)
	}
	return
}

func (c *Checker) checkBlock(b ast.Block, ret types.Type) (errs []*pe.PrettyError) {
	for _, mn := range b {
		errs = append(errs, c.checkNode(mn, ret)...)
	}
	return
}

// typeOf determines the type of any abstract AST node.
func (c *Checker) typeOf(mn ast.MetaNode) (t types.Type, err *pe.PrettyError) {
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
			err = nodeErr(pe.EUnknownFunction, mn)
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
		t, err = f(mn, args)

	default:
		if sel, ok := n.(ast.Selector); !ok {
			err = nodeErr(pe.ETypeOfUnimplemented, mn)
		} else {
			p := sel.Path()
			if len(p) < 1 {
				err = nodeErr(pe.EEmptySelector, mn)
				return
			}
			root := p[0]
			if !root.IsName() {
				err = nodeErr(pe.EBadSelectorRoot, mn)
				return
			}
			rootType, ok := c.scope.getVar(root.Name)
			if !ok {
				err = nodeErr(pe.EUnknownVariable, mn)
				return
			}
			t = rootType
			if len(p) == 1 {
				return
			}
			for _, e := range p[1:] {
				switch {
				case e.IsCast():
					if !t.Equals(e.Cast) {
						err = typeMismatch(mn, e.Cast, t)
						return
					}
					t = e.Cast
				case e.IsName():
					if t.Prim() != types.PStruct {
						err = nodeErr(pe.EBadSelectorParent, mn)
						return
					}
					fieldType, ok := t.(types.StructType).FieldType(e.Name)
					if !ok {
						err = nodeErr(pe.EUnknownField, mn)
						return
					}
					t = fieldType
				default:
					if t.Prim() != types.PArray {
						err = nodeErr(pe.EBadIndexParent, mn)
						return
					}
					t = c.result(t.(types.ArrayType).Element)
				}
			}
		}

	}
	return
}

func (c Checker) result(wrapped types.Type) types.Type {
	return types.MakeStruct(wrapped.String()+" Result",
		"ok", types.Bool,
		"val", wrapped)
}

// basicFuncproto generates a funcproto for a function that accepts a certain
// combination of arguments and always returns the same type.
func basicFuncproto(exp []types.Type, ret types.Type) funcproto {
	return func(mn ast.MetaNode, got []types.Type) (r types.Type, err *pe.PrettyError) {
		r = ret
		if len(got) < len(exp) {
			err = typeMismatch(mn, nil, nil)
		}
		for i, ea := range exp {
			if !ea.Equals(got[i]) {
				err = typeMismatch(mn, ea, got[i])
				return
			}
		}
		return
	}
}
