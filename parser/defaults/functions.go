package defaults

import (
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

// Functions contains definitions of built-in functions for the parser.
var Functions = map[string]*values.Function{
	"print": {
		Name: "print",
		Args: []values.FuncArg{{"a", types.String}},
		Ret:  types.Nothing,
	},
	"to_str": {
		Name: "to_str",
		Args: []values.FuncArg{{"a", types.Any}},
		Ret:  types.String,
	},
	"to_bool": {
		Name: "to_bool",
		Args: []values.FuncArg{{"a", types.Any}},
		Ret:  types.Bool,
	},
	"add_i": {
		Name: "add_i",
		Args: []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:  types.Int,
	},
	"add_f": {
		Name: "add_f",
		Args: []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:  types.Float,
	},
	"add_c": {
		Name: "add_c",
		Args: []values.FuncArg{{"a", types.Char}, {"b", types.Char}},
		Ret:  types.Char,
	},
	"sub_i": {
		Name: "sub_i",
		Args: []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:  types.Int,
	},
	"sub_f": {
		Name: "sub_f",
		Args: []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:  types.Float,
	},
	"sub_c": {
		Name: "sub_c",
		Args: []values.FuncArg{{"a", types.Char}, {"b", types.Char}},
		Ret:  types.Char,
	},
	"mul_i": {
		Name: "mul_i",
		Args: []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:  types.Int,
	},
	"mul_f": {
		Name: "mul_f",
		Args: []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:  types.Float,
	},
	"div_i": {
		Name: "div_i",
		Args: []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:  types.Int,
	},
	"div_f": {
		Name: "div_f",
		Args: []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:  types.Float,
	},
	"mod_i": {
		Name: "mod_i",
		Args: []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:  types.Int,
	},
	"mod_f": {
		Name: "mod_f",
		Args: []values.FuncArg{{"a", types.Float}, {"b", types.Int}},
		Ret:  types.Int,
	},
	"concat": {
		Name: "concat",
		Args: []values.FuncArg{{"a", types.String}, {"b", types.String}},
		Ret:  types.String,
	},
	"not": {
		Name: "not",
		Args: []values.FuncArg{{"a", types.Bool}},
		Ret:  types.Bool,
	},
	"or": {
		Name: "or",
		Args: []values.FuncArg{{"a", types.Bool}, {"b", types.Bool}},
		Ret:  types.Bool,
	},
	"and": {
		Name: "and",
		Args: []values.FuncArg{{"a", types.Bool}, {"b", types.Bool}},
		Ret:  types.Bool,
	},
	"eq": {
		Name: "eq",
		Args: []values.FuncArg{{"a", types.Any}, {"b", types.Any}},
		Ret:  types.Bool,
	},
	"gt_i": {
		Name: "gt_i",
		Args: []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:  types.Bool,
	},
	"gt_f": {
		Name: "gt_f",
		Args: []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:  types.Bool,
	},
	"gt_c": {
		Name: "gt_c",
		Args: []values.FuncArg{{"a", types.Char}, {"b", types.Char}},
		Ret:  types.Bool,
	},
	"lt_i": {
		Name: "lt_i",
		Args: []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:  types.Bool,
	},
	"lt_f": {
		Name: "lt_f",
		Args: []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:  types.Bool,
	},
	"lt_c": {
		Name: "lt_c",
		Args: []values.FuncArg{{"a", types.Char}, {"b", types.Char}},
		Ret:  types.Bool,
	},
	"char_at": {
		Name: "char_at",
		Args: []values.FuncArg{{"s", types.String}, {"i", types.Int}},
		Ret:  types.Char,
	},
	"substr": {
		Name: "substr",
		Args: []values.FuncArg{{"s", types.String}, {"a", types.Int}, {"b", types.Int}},
		Ret:  types.String,
	},
	"char_append": {
		Name: "char_append",
		Args: []values.FuncArg{{"s", types.String}, {"c", types.Char}},
		Ret:  types.String,
	},
	"str_len": {
		Name: "str_len",
		Args: []values.FuncArg{{"s", types.String}},
		Ret:  types.Int,
	},
	"skol": {
		Name: "skol",
		Args: []values.FuncArg{
			{"engine", types.String},
			{"ver", types.Float},
		},
		Ret: types.Nothing,
	},
	"ctoi": {
		Name: "ctoi",
		Args: []values.FuncArg{{"c", types.Char}},
		Ret:  types.Int,
	},
	"open": {
		Name: "open",
		Args: []values.FuncArg{{"fn", types.String}},
		Ret:  file_descriptor_result,
	},
	"fgetc": {
		Name: "fgetc",
		Args: []values.FuncArg{{"fd", file_descriptor}},
		Ret:  char_result,
	},
	"close": {
		Name: "close",
		Args: []values.FuncArg{{"fd", file_descriptor}},
		Ret:  types.Nothing,
	},
}
