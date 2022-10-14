package parser

import (
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/values/types"
)

var defaultFunctions = map[string]ast.Func{
	"add": {
		Name: "add",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Any},
			{Name: "b", Type: types.Any},
		},
		Ret: types.Any,
	},
	"sub": {
		Name: "sub",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Any},
			{Name: "b", Type: types.Any},
		},
		Ret: types.Any,
	},
	"mul": {
		Name: "mul",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Any},
			{Name: "b", Type: types.Any},
		},
		Ret: types.Any,
	},
	"div": {
		Name: "div",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Any},
			{Name: "b", Type: types.Any},
		},
		Ret: types.Any,
	},
	"pow": {
		Name: "pow",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Any},
			{Name: "b", Type: types.Any},
		},
		Ret: types.Any,
	},
	"mod": {
		Name: "mod",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Any},
			{Name: "b", Type: types.Int},
		},
		Ret: types.Any,
	},

	"eq": {
		Name: "eq",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Any},
			{Name: "b", Type: types.Any},
		},
		Ret: types.Bool,
	},
	"gt": {
		Name: "gt",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Any},
			{Name: "b", Type: types.Any},
		},
		Ret: types.Bool,
	},
	"lt": {
		Name: "lt",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Any},
			{Name: "b", Type: types.Any},
		},
		Ret: types.Bool,
	},

	"not": {
		Name: "not",
		Args: []types.Descriptor{
			{Name: "v", Type: types.Bool},
		},
		Ret: types.Bool,
	},
	"and": {
		Name: "and",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Bool},
			{Name: "b", Type: types.Bool},
		},
		Ret: types.Bool,
	},
	"or": {
		Name: "or",
		Args: []types.Descriptor{
			{Name: "a", Type: types.Bool},
			{Name: "b", Type: types.Bool},
		},
		Ret: types.Bool,
	},

	"append": {
		Name: "append",
		Args: []types.Descriptor{
			{Name: "a", Type: types.ArrayType{Element: types.Any}},
			{Name: "b", Type: types.Any},
		},
		Ret: types.ArrayType{Element: types.Any},
	},
	"concat": {
		Name: "concat",
		Args: []types.Descriptor{
			{Name: "a", Type: types.ArrayType{Element: types.Any}},
			{Name: "b", Type: types.ArrayType{Element: types.Any}},
		},
		Ret: types.ArrayType{Element: types.Any},
	},
	"slice": {
		Name: "slice",
		Args: []types.Descriptor{
			{Name: "arr", Type: types.ArrayType{Element: types.Any}},
			{Name: "start", Type: types.Int},
			{Name: "end", Type: types.Int},
		},
		Ret: types.ArrayType{Element: types.Any},
	},
	"at": {
		Name: "at",
		Args: []types.Descriptor{
			{Name: "arr", Type: types.ArrayType{Element: types.Any}},
			{Name: "idx", Type: types.Int},
		},
		Ret: types.Any,
	},
	"len": {
		Name: "len",
		Args: []types.Descriptor{
			{Name: "arr", Type: types.ArrayType{Element: types.Any}},
		},
		Ret: types.Int,
	},

	"str": {
		Name: "str",
		Args: []types.Descriptor{
			{Name: "v", Type: types.Any},
		},
		Ret: types.String,
	},
	"bool": {
		Name: "bool",
		Args: []types.Descriptor{
			{Name: "v", Type: types.Any},
		},
		Ret: types.Bool,
	},

	"parse_bool": {
		Name: "parse_bool",
		Args: []types.Descriptor{
			{Name: "raw", Type: types.String},
		},
		Ret: types.Result(types.Bool),
	},
	"char": {
		Name: "char",
		Args: []types.Descriptor{
			{Name: "raw", Type: types.String},
		},
		Ret: types.Result(types.Char),
	},
	"int": {
		Name: "int",
		Args: []types.Descriptor{
			{Name: "raw", Type: types.String},
		},
		Ret: types.Result(types.Int),
	},
	"float": {
		Name: "float",
		Args: []types.Descriptor{
			{Name: "raw", Type: types.String},
		},
		Ret: types.Result(types.Float),
	},
}
