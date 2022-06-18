package parser

import (
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type Function struct {
	Name string
	Args []values.FuncArg
	Ret  *values.Type
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
		Args: []values.FuncArg{{"a", values.String}},
		Ret:  values.Nothing,
	},
	"to_str": {
		Name: "to_str",
		Args: []values.FuncArg{{"a", values.Any}},
		Ret:  values.String,
	},
	"to_bool": {
		Name: "to_bool",
		Args: []values.FuncArg{{"a", values.Any}},
		Ret:  values.Bool,
	},
	"add_i": {
		Name: "add_i",
		Args: []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:  values.Int,
	},
	"add_f": {
		Name: "add_f",
		Args: []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:  values.Float,
	},
	"add_c": {
		Name: "add_c",
		Args: []values.FuncArg{{"a", values.Char}, {"b", values.Char}},
		Ret:  values.Char,
	},
	"sub_i": {
		Name: "sub_i",
		Args: []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:  values.Int,
	},
	"sub_f": {
		Name: "sub_f",
		Args: []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:  values.Float,
	},
	"sub_c": {
		Name: "sub_c",
		Args: []values.FuncArg{{"a", values.Char}, {"b", values.Char}},
		Ret:  values.Char,
	},
	"mul_i": {
		Name: "mul_i",
		Args: []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:  values.Int,
	},
	"mul_f": {
		Name: "mul_f",
		Args: []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:  values.Float,
	},
	"div_i": {
		Name: "div_i",
		Args: []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:  values.Int,
	},
	"div_f": {
		Name: "div_f",
		Args: []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:  values.Float,
	},
	"mod_i": {
		Name: "mod_i",
		Args: []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:  values.Int,
	},
	"mod_f": {
		Name: "mod_f",
		Args: []values.FuncArg{{"a", values.Float}, {"b", values.Int}},
		Ret:  values.Int,
	},
	"concat": {
		Name: "concat",
		Args: []values.FuncArg{{"a", values.String}, {"b", values.String}},
		Ret:  values.String,
	},
	"not": {
		Name: "not",
		Args: []values.FuncArg{{"a", values.Bool}},
		Ret:  values.Bool,
	},
	"or": {
		Name: "or",
		Args: []values.FuncArg{{"a", values.Bool}, {"b", values.Bool}},
		Ret:  values.Bool,
	},
	"and": {
		Name: "and",
		Args: []values.FuncArg{{"a", values.Bool}, {"b", values.Bool}},
		Ret:  values.Bool,
	},
	"eq": {
		Name: "eq",
		Args: []values.FuncArg{{"a", values.Any}, {"b", values.Any}},
		Ret:  values.Bool,
	},
	"gt_i": {
		Name: "gt_i",
		Args: []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:  values.Bool,
	},
	"gt_f": {
		Name: "gt_f",
		Args: []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:  values.Bool,
	},
	"gt_c": {
		Name: "gt_c",
		Args: []values.FuncArg{{"a", values.Char}, {"b", values.Char}},
		Ret:  values.Bool,
	},
	"lt_i": {
		Name: "lt_i",
		Args: []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:  values.Bool,
	},
	"lt_f": {
		Name: "lt_f",
		Args: []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:  values.Bool,
	},
	"lt_c": {
		Name: "lt_c",
		Args: []values.FuncArg{{"a", values.Char}, {"b", values.Char}},
		Ret:  values.Bool,
	},
	"char_at": {
		Name: "char_at",
		Args: []values.FuncArg{{"s", values.String}, {"i", values.Int}},
		Ret:  values.Char,
	},
	"substr": {
		Name: "substr",
		Args: []values.FuncArg{{"s", values.String}, {"a", values.Int}, {"b", values.Int}},
		Ret:  values.String,
	},
	"char_append": {
		Name: "char_append",
		Args: []values.FuncArg{{"s", values.String}, {"c", values.Char}},
		Ret:  values.String,
	},
	"str_len": {
		Name: "str_len",
		Args: []values.FuncArg{{"s", values.String}},
		Ret:  values.Int,
	},
	"skol": {
		Name: "skol",
		Args: []values.FuncArg{
			{"engine", values.String},
			{"ver", values.Float},
		},
		Ret: values.Nothing,
	},
	"ctoi": {
		Name: "ctoi",
		Args: []values.FuncArg{{"c", values.Char}},
		Ret:  values.Int,
	},
}
