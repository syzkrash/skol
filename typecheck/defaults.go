package typecheck

import (
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/parser/values/types"
)

func makeFuncproto(ret types.Type, args ...types.Type) funcproto {
	return basicFuncproto(args, ret)
}

func mathFuncproto(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
	if len(t) < 2 {
		err = pe.New(pe.ENeedMoreArgs)
		return
	}

	if !t[1].Equals(t[0]) {
		err = typeMismatch(mn, t[0], t[1])
		return
	}

	return t[0], nil
}

func cmpFuncproto(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
	if len(t) < 2 {
		err = pe.New(pe.ENeedMoreArgs)
		return
	}

	if !t[1].Equals(t[0]) {
		err = typeMismatch(mn, t[0], t[1])
		return
	}

	return types.Bool, nil
}

func result(t types.Type) types.Type {
	return types.MakeStruct(t.String()+"Result",
		"ok", types.Bool,
		"value", t)
}

var defaultFuncs = map[string]funcproto{
	"add": mathFuncproto,
	"sub": mathFuncproto,
	"mul": mathFuncproto,
	"div": mathFuncproto,
	"pow": mathFuncproto,
	"mod": func(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
		if len(t) < 2 {
			err = pe.New(pe.ENeedMoreArgs)
			return
		}

		if !types.Int.Equals(t[1]) {
			err = typeMismatch(mn, types.Int, t[1])
			return
		}

		return t[1], nil
	},

	"eq": makeFuncproto(types.Bool, types.Any, types.Any),
	"gt": cmpFuncproto,
	"lt": cmpFuncproto,

	"not": makeFuncproto(types.Bool, types.Bool),
	"and": makeFuncproto(types.Bool, types.Bool, types.Bool),
	"or":  makeFuncproto(types.Bool, types.Bool, types.Bool),

	"append": func(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
		if len(t) < 2 {
			err = pe.New(pe.ENeedMoreArgs)
			return
		}

		if types.String.Equals(t[0]) && types.Char.Equals(t[1]) {
			return types.String, nil
		}

		if t[0].Prim() != types.PArray {
			err = typeMismatch(mn, types.ArrayType{Element: t[1]}, t[0])
			return
		}

		t0a := t[0].(types.ArrayType)
		if !t[1].Equals(t0a.Element) {
			err = typeMismatch(mn, t0a.Element, t[1])
			return
		}

		return t0a, nil
	},
	"concat": func(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
		if len(t) < 2 {
			err = pe.New(pe.ENeedMoreArgs)
			return
		}

		if types.String.Equals(t[0]) && types.String.Equals(t[1]) {
			return types.String, nil
		}

		if t[0].Prim() != types.PArray {
			err = typeMismatch(mn, types.ArrayType{Element: types.Any}, t[0])
			return
		}
		if t[1].Prim() != types.PArray {
			err = typeMismatch(mn, types.ArrayType{Element: types.Any}, t[1])
			return
		}

		t0a := t[0].(types.ArrayType)
		t1a := t[1].(types.ArrayType)

		if !t1a.Element.Equals(t0a.Element) {
			err = typeMismatch(mn, t0a.Element, t1a.Element)
			return
		}

		return t0a, nil
	},
	"slice": func(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
		if len(t) < 3 {
			err = pe.New(pe.ENeedMoreArgs)
			return
		}

		if types.String.Equals(t[0]) && types.Int.Equals(t[1]) && types.Int.Equals(t[2]) {
			return types.String, nil
		}

		if t[0].Prim() != types.PArray {
			err = typeMismatch(mn, types.ArrayType{Element: types.Any}, t[0])
			return
		}
		if !types.Int.Equals(t[1]) {
			err = typeMismatch(mn, types.Int, t[1])
			return
		}
		if !types.Int.Equals(t[2]) {
			err = typeMismatch(mn, types.Int, t[2])
			return
		}

		return t[0], nil
	},
	"at": func(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
		if len(t) < 2 {
			err = pe.New(pe.ENeedMoreArgs)
			return
		}

		if types.String.Equals(t[0]) && types.Int.Equals(t[1]) {
			return types.Char, nil
		}

		if t[0].Prim() != types.PArray {
			err = typeMismatch(mn, types.ArrayType{Element: types.Any}, t[0])
			return
		}
		if !types.Int.Equals(t[1]) {
			err = typeMismatch(mn, types.Int, t[1])
			return
		}

		return t[0].(types.ArrayType).Element, nil
	},
	"len": func(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
		if len(t) < 1 {
			err = pe.New(pe.ENeedMoreArgs)
			return
		}

		if types.String.Equals(t[0]) {
			return types.Int, nil
		}

		if t[0].Prim() != types.PArray {
			err = typeMismatch(mn, types.ArrayType{Element: types.Any}, t[0])
			return
		}

		return types.Int, nil
	},

	"str":        makeFuncproto(types.String, types.Any),
	"bool":       makeFuncproto(types.Bool, types.Any),
	"parse_bool": makeFuncproto(result(types.Bool), types.String),
	"char":       makeFuncproto(result(types.Char), types.String),
	"int":        makeFuncproto(result(types.Int), types.String),
	"float":      makeFuncproto(result(types.Float), types.String),
}
