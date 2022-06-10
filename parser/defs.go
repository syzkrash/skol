package parser

import (
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type Function struct {
	Name string
	Args map[string]values.ValueType
	Ret  values.ValueType
}

func DefinedFunction(n *nodes.FuncDefNode) *Function {
	return &Function{
		Name: n.Name,
		Args: n.Args,
		Ret:  n.Ret,
	}
}

func ExternFunction(n *nodes.FuncExternNode) *Function {
	var name string
	if n.Intern != "" {
		name = n.Intern
	} else {
		name = n.Name
	}
	return &Function{
		Name: name,
		Args: n.Args,
		Ret:  n.Ret,
	}
}

var DefaultFuncs = map[string]*Function{
	"print": {
		Name: "print",
		Args: map[string]values.ValueType{"a": values.VtString},
		Ret:  values.VtNothing,
	},
	"to_str": {
		Name: "to_str",
		Args: map[string]values.ValueType{"a": values.VtAny},
		Ret:  values.VtString,
	},
	"to_bool": {
		Name: "to_bool",
		Args: map[string]values.ValueType{"a": values.VtAny},
		Ret:  values.VtBool,
	},
	"add_i": {
		Name: "add_i",
		Args: map[string]values.ValueType{"a": values.VtInteger, "b": values.VtInteger},
		Ret:  values.VtInteger,
	},
	"add_f": {
		Name: "add_f",
		Args: map[string]values.ValueType{"a": values.VtFloat, "b": values.VtFloat},
		Ret:  values.VtFloat,
	},
	"sub_i": {
		Name: "sub_i",
		Args: map[string]values.ValueType{"a": values.VtInteger, "b": values.VtInteger},
		Ret:  values.VtInteger,
	},
	"sub_f": {
		Name: "sub_f",
		Args: map[string]values.ValueType{"a": values.VtFloat, "b": values.VtFloat},
		Ret:  values.VtFloat,
	},
	"mul_i": {
		Name: "mul_i",
		Args: map[string]values.ValueType{"a": values.VtInteger, "b": values.VtInteger},
		Ret:  values.VtInteger,
	},
	"mul_f": {
		Name: "mul_f",
		Args: map[string]values.ValueType{"a": values.VtFloat, "b": values.VtFloat},
		Ret:  values.VtFloat,
	},
	"div_i": {
		Name: "div_i",
		Args: map[string]values.ValueType{"a": values.VtInteger, "b": values.VtInteger},
		Ret:  values.VtInteger,
	},
	"div_f": {
		Name: "div_f",
		Args: map[string]values.ValueType{"a": values.VtFloat, "b": values.VtFloat},
		Ret:  values.VtFloat,
	},
	"mod_i": {
		Name: "mod_i",
		Args: map[string]values.ValueType{"a": values.VtInteger, "b": values.VtInteger},
		Ret:  values.VtInteger,
	},
	"mod_f": {
		Name: "mod_f",
		Args: map[string]values.ValueType{"a": values.VtFloat, "b": values.VtInteger},
		Ret:  values.VtInteger,
	},
	"concat": {
		Name: "concat",
		Args: map[string]values.ValueType{"a": values.VtString, "b": values.VtString},
		Ret:  values.VtString,
	},
	"not": {
		Name: "not",
		Args: map[string]values.ValueType{"a": values.VtBool},
		Ret:  values.VtBool,
	},
	"or": {
		Name: "or",
		Args: map[string]values.ValueType{"a": values.VtBool, "b": values.VtBool},
		Ret:  values.VtBool,
	},
	"and": {
		Name: "and",
		Args: map[string]values.ValueType{"a": values.VtBool, "b": values.VtBool},
		Ret:  values.VtBool,
	},
	"eq": {
		Name: "eq",
		Args: map[string]values.ValueType{"a": values.VtAny, "b": values.VtAny},
		Ret:  values.VtBool,
	},
	"gt_i": {
		Name: "gt_i",
		Args: map[string]values.ValueType{"a": values.VtInteger, "b": values.VtInteger},
		Ret:  values.VtBool,
	},
	"gt_f": {
		Name: "gt_f",
		Args: map[string]values.ValueType{"a": values.VtFloat, "b": values.VtFloat},
		Ret:  values.VtBool,
	},
	"lt_i": {
		Name: "lt_i",
		Args: map[string]values.ValueType{"a": values.VtInteger, "b": values.VtInteger},
		Ret:  values.VtBool,
	},
	"lt_f": {
		Name: "lt_f",
		Args: map[string]values.ValueType{"a": values.VtFloat, "b": values.VtFloat},
		Ret:  values.VtBool,
	},
	"char_at": {
		Name: "char_at",
		Args: map[string]values.ValueType{
			"s": values.VtString,
			"i": values.VtInteger,
		},
		Ret: values.VtChar,
	},
	"substr": {
		Name: "substr",
		Args: map[string]values.ValueType{
			"s": values.VtString,
			"a": values.VtInteger,
			"b": values.VtInteger,
		},
		Ret: values.VtString,
	},
	"char_append": {
		Name: "char_append",
		Args: map[string]values.ValueType{
			"s": values.VtString,
			"c": values.VtChar,
		},
		Ret: values.VtString,
	},
}
