package parser

import (
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type Function struct {
	Name string
	Args []values.FuncArg
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
		Args: []values.FuncArg{{"a", values.VtString}},
		Ret:  values.VtNothing,
	},
	"to_str": {
		Name: "to_str",
		Args: []values.FuncArg{{"a", values.VtAny}},
		Ret:  values.VtString,
	},
	"to_bool": {
		Name: "to_bool",
		Args: []values.FuncArg{{"a", values.VtAny}},
		Ret:  values.VtBool,
	},
	"add_i": {
		Name: "add_i",
		Args: []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:  values.VtInteger,
	},
	"add_f": {
		Name: "add_f",
		Args: []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:  values.VtFloat,
	},
	"add_c": {
		Name: "add_c",
		Args: []values.FuncArg{{"a", values.VtChar}, {"b", values.VtChar}},
		Ret:  values.VtChar,
	},
	"sub_i": {
		Name: "sub_i",
		Args: []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:  values.VtInteger,
	},
	"sub_f": {
		Name: "sub_f",
		Args: []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:  values.VtFloat,
	},
	"sub_c": {
		Name: "sub_c",
		Args: []values.FuncArg{{"a", values.VtChar}, {"b", values.VtChar}},
		Ret:  values.VtChar,
	},
	"mul_i": {
		Name: "mul_i",
		Args: []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:  values.VtInteger,
	},
	"mul_f": {
		Name: "mul_f",
		Args: []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:  values.VtFloat,
	},
	"div_i": {
		Name: "div_i",
		Args: []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:  values.VtInteger,
	},
	"div_f": {
		Name: "div_f",
		Args: []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:  values.VtFloat,
	},
	"mod_i": {
		Name: "mod_i",
		Args: []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:  values.VtInteger,
	},
	"mod_f": {
		Name: "mod_f",
		Args: []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtInteger}},
		Ret:  values.VtInteger,
	},
	"concat": {
		Name: "concat",
		Args: []values.FuncArg{{"a", values.VtString}, {"b", values.VtString}},
		Ret:  values.VtString,
	},
	"not": {
		Name: "not",
		Args: []values.FuncArg{{"a", values.VtBool}},
		Ret:  values.VtBool,
	},
	"or": {
		Name: "or",
		Args: []values.FuncArg{{"a", values.VtBool}, {"b", values.VtBool}},
		Ret:  values.VtBool,
	},
	"and": {
		Name: "and",
		Args: []values.FuncArg{{"a", values.VtBool}, {"b", values.VtBool}},
		Ret:  values.VtBool,
	},
	"eq": {
		Name: "eq",
		Args: []values.FuncArg{{"a", values.VtAny}, {"b", values.VtAny}},
		Ret:  values.VtBool,
	},
	"gt_i": {
		Name: "gt_i",
		Args: []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:  values.VtBool,
	},
	"gt_f": {
		Name: "gt_f",
		Args: []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:  values.VtBool,
	},
	"gt_c": {
		Name: "gt_c",
		Args: []values.FuncArg{{"a", values.VtChar}, {"b", values.VtChar}},
		Ret:  values.VtBool,
	},
	"lt_i": {
		Name: "lt_i",
		Args: []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:  values.VtBool,
	},
	"lt_f": {
		Name: "lt_f",
		Args: []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:  values.VtBool,
	},
	"lt_c": {
		Name: "lt_c",
		Args: []values.FuncArg{{"a", values.VtChar}, {"b", values.VtChar}},
		Ret:  values.VtBool,
	},
	"char_at": {
		Name: "char_at",
		Args: []values.FuncArg{{"s", values.VtString}, {"i", values.VtInteger}},
		Ret:  values.VtChar,
	},
	"substr": {
		Name: "substr",
		Args: []values.FuncArg{{"s", values.VtString}, {"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:  values.VtString,
	},
	"char_append": {
		Name: "char_append",
		Args: []values.FuncArg{{"s", values.VtString}, {"c", values.VtChar}},
		Ret:  values.VtString,
	},
	"str_len": {
		Name: "str_len",
		Args: []values.FuncArg{{"s", values.VtString}},
		Ret:  values.VtInteger,
	},
	"skol": {
		Name: "skol",
		Args: []values.FuncArg{
			{"engine", values.VtString},
			{"ver", values.VtFloat},
		},
		Ret: values.VtNothing,
	},
	"ctoi": {
		Name: "ctoi",
		Args: []values.FuncArg{{"c", values.VtChar}},
		Ret:  values.VtInteger,
	},
}
