package typecheck

import (
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/parser/values/types"
)

// Checker ensures type correctness of an AST.
type Checker struct {
	scope *scope
	errs  chan error
}

// NewChecker creates a blank Checker.
func NewChecker(errOut chan error) *Checker {
	return &Checker{
		scope: &scope{
			parent: nil,
			vars:   make(map[string]types.Type),
			funcs:  make(map[string]funcproto),
		},
		errs: errOut,
	}
}

// Check thoroughly inspects the provided AST for any typing-related errors
// that may have occured.
func (c *Checker) Check(tree ast.AST) {
	for _, v := range tree.Typedefs {
		c.scope.vars[v.Name] = v.Type
	}
	for _, v := range tree.Vars {
		t, ok := c.typeOf(v.Value)
		if ok {
			c.scope.vars[v.Name] = t
		}
	}
	// first loop to declare functions
	for _, f := range tree.Funcs {
		c.scope.funcs[f.Name] = funcproto{
			Args: f.Args,
			Ret:  f.Ret,
		}
	}
	// second loop to typecheck function bodies with function type information
	for _, f := range tree.Funcs {
		args := make(map[string]types.Type)
		for _, a := range f.Args {
			args[a.Name] = a.Type
		}
		c.checkFunc(args, f.Ret, f.Body)
	}
}

func (c *Checker) checkNode(mn ast.MetaNode, ret types.Type) {
	n := mn.Node

	switch n.Kind() {
	// literals
	case ast.NStruct:
		nstruct := n.(ast.StructNode)
		for i, a := range nstruct.Args {
			at, ok := c.typeOf(a)
			if !ok {
				continue
			}
			ft := nstruct.Type.Fields[i].Type
			if !ft.Equals(at) {
				c.typeMismatch(a, ft, at)
			}
		}
	case ast.NArray:
		narray := n.(ast.ArrayNode)
		for _, e := range narray.Elems {
			et, ok := c.typeOf(e)
			if !ok {
				continue
			}
			if !narray.Type.Equals(et) {
				c.typeMismatch(e, narray.Type, et)
			}
		}

	// control flow
	case ast.NIf:
		nif := n.(ast.IfNode)
		c.checkBranch(nif.Main, ret)
		for _, b := range nif.Other {
			c.checkBranch(b, ret)
		}
		c.checkBlock(nif.Else, ret)
	case ast.NWhile:
		nwhile := n.(ast.WhileNode)
		c.checkCond(nwhile.Cond)
		c.checkBlock(nwhile.Block, ret)
	case ast.NReturn:
		nreturn := n.(ast.ReturnNode)
		rt, ok := c.typeOf(nreturn.Value)
		if !ok {
			return
		}
		if !ret.Equals(rt) {
			c.typeMismatch(nreturn.Value, ret, rt)
		}

	// definitions
	case ast.NVarSet:
		nvarset := n.(ast.VarSetNode)
		c.checkNode(nvarset.Value, ret)
		nvt, ok := c.typeOf(nvarset.Value)
		if !ok {
			return
		}
		if ovt, ok := c.scope.getVar(nvarset.Var); ok {
			if !ovt.Equals(nvt) {
				c.typeMismatch(mn, ovt, nvt)
			}
		} else {
			c.scope.setVar(nvarset.Var, nvt)
		}
	case ast.NVarDef:
		nvardef := n.(ast.VarDefNode)
		c.scope.setVar(nvardef.Var, nvardef.Type)
	case ast.NVarSetTyped:
		nvarsettyped := n.(ast.VarSetTypedNode)
		c.checkNode(nvarsettyped.Value, ret)
		nvt, ok := c.typeOf(nvarsettyped.Value)
		if !ok {
			return
		}
		ovt, ok := c.scope.getVar(nvarsettyped.Var)
		if ok {
			if !ovt.Equals(nvarsettyped.Type) {
				c.varRetype(nvarsettyped.Value, ovt, nvarsettyped.Type)
			}
			if !ovt.Equals(nvt) {
				c.typeMismatch(nvarsettyped.Value, ovt, nvt)
			}
			if !ovt.Equals(nvarsettyped.Type) {
				c.typeMismatch(nvarsettyped.Value, ovt, nvt)
			}
		}
	case ast.NFuncDef:
		nfuncdef := n.(ast.FuncDefNode)
		args := make(map[string]types.Type)
		for _, a := range nfuncdef.Proto {
			args[a.Name] = a.Type
		}
		c.checkFunc(args, nfuncdef.Ret, nfuncdef.Body)
		c.scope.funcs[nfuncdef.Name] = funcproto{
			Args: nfuncdef.Proto,
			Ret:  nfuncdef.Ret,
		}
	case ast.NFuncExtern:
		nfuncextern := n.(ast.FuncExternNode)
		c.scope.funcs[nfuncextern.Alias] = funcproto{
			Args: nfuncextern.Proto,
			Ret:  nfuncextern.Ret,
		}

	// others
	case ast.NFuncCall:
		nfunccall := n.(ast.FuncCallNode)
		args := make([]types.Type, len(nfunccall.Args))
		for i, a := range nfunccall.Args {
			t, ok := c.typeOf(a)
			if !ok {
				return
			}
			args[i] = t
		}
		f, ok := c.scope.getFunc(nfunccall.Func)
		if !ok {
			bf, ok := builtins[nfunccall.Func]
			if !ok {
				c.nodeErr(pe.EUnknownFunction, mn)
				return
			}
			if _, err := bf(mn, args); err != nil {
				c.errs <- err
			}
			return
		}
		if len(args) != len(f.Args) {
			c.nodeErr(pe.ENeedMoreArgs, mn)
			return
		}
		for i := 0; i < len(args); i++ {
			if !f.Args[i].Type.Equals(args[i]) {
				c.typeMismatch(mn, f.Args[i].Type, args[i])
			}
		}
	}
	return
}

// checkFunc ensures type correctness given the function arguments' types and
// the return type.
func (c *Checker) checkFunc(args map[string]types.Type, ret types.Type, body ast.Block) {
	c.scope = c.scope.sub()
	for n, t := range args {
		c.scope.vars[n] = t
	}
	c.checkBlock(body, ret)
	c.scope = c.scope.parent
	return
}

func (c *Checker) checkBranch(b ast.Branch, ret types.Type) {
	c.checkCond(b.Cond)
	c.checkBlock(b.Block, ret)
}

func (c *Checker) checkCond(mn ast.MetaNode) bool {
	ct, ok := c.typeOf(mn)
	if !ok {
		return false
	}
	if !types.Bool.Equals(ct) {
		c.typeMismatch(mn, types.Bool, ct)
		return false
	}
	return true
}

func (c *Checker) checkBlock(b ast.Block, ret types.Type) {
	for _, mn := range b {
		c.checkNode(mn, ret)
	}
	return
}

// typeOf determines the type of any abstract AST node.
func (c *Checker) typeOf(mn ast.MetaNode) (t types.Type, ok bool) {
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
		args := make([]types.Type, len(nfunccall.Args))
		for i, a := range nfunccall.Args {
			var at types.Type
			at, ok = c.typeOf(a)
			if !ok {
				return
			}
			args[i] = at
		}
		var f funcproto
		f, ok = c.scope.getFunc(nfunccall.Func)
		if !ok {
			var bf builtin
			bf, ok = builtins[nfunccall.Func]
			if !ok {
				c.nodeErr(pe.EUnknownFunction, mn)
				return
			}
			var err error
			t, err = bf(mn, args)
			if err != nil {
				ok = false
				c.errs <- err
			}
			return
		}
		if len(args) != len(f.Args) {
			c.nodeErr(pe.ENeedMoreArgs, mn)
			ok = false
			return
		}
		for i := 0; i < len(args); i++ {
			if !f.Args[i].Type.Equals(args[i]) {
				c.typeMismatch(mn, f.Args[i].Type, args[i])
				ok = false
				return
			}
		}
		t = f.Ret

	default:
		var sel ast.Selector
		if sel, ok = n.(ast.Selector); !ok {
			c.nodeErr(pe.ETypeOfUnimplemented, mn)
			ok = false
		} else {
			p := sel.Path()
			if len(p) < 1 {
				c.nodeErr(pe.EEmptySelector, mn)
				ok = false
				return
			}
			root := p[0]
			if !root.IsName() {
				c.nodeErr(pe.EBadSelectorRoot, mn)
				ok = false
				return
			}
			var rootType types.Type
			rootType, ok = c.scope.getVar(root.Name)
			if !ok {
				c.nodeErr(pe.EUnknownVariable, mn)
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
						c.typeMismatch(mn, e.Cast, t)
						ok = false
						return
					}
					t = e.Cast
				case e.IsName():
					if t.Prim() != types.PStruct {
						c.nodeErr(pe.EBadSelectorParent, mn)
						ok = false
						return
					}
					var fieldType types.Type
					fieldType, ok = t.(types.StructType).FieldType(e.Name)
					if !ok {
						c.nodeErr(pe.EUnknownField, mn)
						return
					}
					t = fieldType
				default:
					if types.String.Equals(t) {
						t = types.Result(types.Char)
						return
					}
					if t.Prim() != types.PArray {
						c.nodeErr(pe.EBadIndexParent, mn)
						ok = false
						return
					}
					t = types.Result(t.(types.ArrayType).Element)
				}
			}
		}

	}
	return
}

func (c *Checker) typeMismatch(mn ast.MetaNode, want, got types.Type) {
	c.errs <- typeMismatch(mn, want, got)
}

func (c *Checker) varRetype(mn ast.MetaNode, old, new types.Type) {
	c.errs <- pe.New(pe.EVarTypeChanged).Section("Original type", "%s", old).Section("New type", "%s", new).Section("Caused by", "`%s` node at %s", mn.Node.Kind(), mn.Where)
}

func (c *Checker) nodeErr(e pe.ErrorCode, mn ast.MetaNode) {
	c.errs <- nodeErr(e, mn)
}

func typeMismatch(mn ast.MetaNode, want, got types.Type) *pe.PrettyError {
	return pe.New(pe.ETypeMismatch).Section("Wanted type", "%s", want).Section("Got type", "%s", got).Section("Caused by", "`%s` node at %s", mn.Node.Kind(), mn.Where)
}

func nodeErr(e pe.ErrorCode, mn ast.MetaNode) *pe.PrettyError {
	return pe.New(e).Section("Caused by", "`%s` node at %s", mn.Node.Kind(), mn.Where)
}
