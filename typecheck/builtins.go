package typecheck

import (
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/parser/values/types"
)

// builtin here represents only the prototype of a builting function. the
// prototype is responsible for ensuring argument correctness and determining
// the return value of the function. the "generic" builtin functions will have
// return types that change with the passed arguments (append or concat for
// example), while some will simply not do any extra checking at all (eq)
type builtin func(ast.MetaNode, []types.Type) (types.Type, *pe.PrettyError)

func simpleBuiltin(ret types.Type, exp ...types.Type) builtin {
	return func(mn ast.MetaNode, got []types.Type) (r types.Type, err *pe.PrettyError) {
		r = ret
		if len(got) < len(exp) {
			err = nodeErr(pe.ENeedMoreArgs, mn)
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

func mathBuiltin(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
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

func cmpBuiltin(mn ast.MetaNode, t []types.Type) (rt types.Type, err *pe.PrettyError) {
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

var builtins = map[string]builtin{
	"add": mathBuiltin,
	"sub": mathBuiltin,
	"mul": mathBuiltin,
	"div": mathBuiltin,
	"pow": mathBuiltin,
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

	"eq": simpleBuiltin(types.Bool, types.Any, types.Any),
	"gt": cmpBuiltin,
	"lt": cmpBuiltin,

	"not": simpleBuiltin(types.Bool, types.Bool),
	"and": simpleBuiltin(types.Bool, types.Bool, types.Bool),
	"or":  simpleBuiltin(types.Bool, types.Bool, types.Bool),

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

	"str":        simpleBuiltin(types.String, types.Any),
	"bool":       simpleBuiltin(types.Bool, types.Any),
	"parse_bool": simpleBuiltin(types.Result(types.Bool), types.String),
	"char":       simpleBuiltin(types.Result(types.Char), types.String),
	"int":        simpleBuiltin(types.Result(types.Int), types.String),
	"float":      simpleBuiltin(types.Result(types.Float), types.String),

	"print": simpleBuiltin(types.Nothing, types.String),
}
