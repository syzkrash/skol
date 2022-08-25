package parser_test

import (
	"testing"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestVarDef(t *testing.T) {
	p, src := makeParser("VarDef")

	cases := map[string]*nodes.VarDefNode{
		"%SomeBool/bool": {
			VarType: types.Bool,
			Var:     "SomeBool",
			Value:   nil,
		},
		"%SomeChar/char": {
			VarType: types.Char,
			Var:     "SomeChar",
			Value:   nil,
		},
		"%SomeInt/int": {
			VarType: types.Int,
			Var:     "SomeInt",
			Value:   nil,
		},
		"%SomeFloat/float": {
			VarType: types.Float,
			Var:     "SomeFloat",
			Value:   nil,
		},
		"%SomeString/string": {
			VarType: types.String,
			Var:     "SomeString",
			Value:   nil,
		},
		"%ExplicitBool/b:/": {
			VarType: types.Bool,
			Var:     "ExplicitBool",
			Value:   &nodes.BooleanNode{Bool: false},
		},
		"%ExplicitChar/c:'E'": {
			VarType: types.Char,
			Var:     "ExplicitChar",
			Value:   &nodes.CharNode{Char: 'E'},
		},
		"%ExplicitInt/i:0": {
			VarType: types.Int,
			Var:     "ExplicitInt",
			Value:   &nodes.IntegerNode{Int: 0},
		},
		"%ExplicitFloat/f:0.5": {
			VarType: types.Float,
			Var:     "ExplicitFloat",
			Value:   &nodes.FloatNode{Float: 0.5},
		},
		"%ExplicitString/s:\"pp\"": {
			VarType: types.String,
			Var:     "ExplicitString",
			Value:   &nodes.StringNode{Str: "pp"},
		},
		"%B0: *": {
			VarType: types.Bool,
			Var:     "B0",
			Value:   &nodes.BooleanNode{Bool: true},
		},
		"%C1: 'q'": {
			VarType: types.Char,
			Var:     "C1",
			Value:   &nodes.CharNode{Char: 'q'},
		},
		"%I2: 123": {
			VarType: types.Int,
			Var:     "I2",
			Value:   &nodes.IntegerNode{Int: 123},
		},
		"%F3: 12.34": {
			VarType: types.Float,
			Var:     "F3",
			Value:   &nodes.FloatNode{Float: 12.34},
		},
		"%S4: \"foo\"": {
			VarType: types.String,
			Var:     "S4",
			Value:   &nodes.StringNode{Str: "foo"},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, out)
		t.Log("OK")
	}
}

func TestFuncDef(t *testing.T) {
	p, src := makeParser("FuncDef")

	cases := map[string]*nodes.FuncDefNode{
		"$a()": {
			Name: "a",
			Args: []values.FuncArg{},
			Ret:  types.Nothing,
			Body: []nodes.Node{},
		},
		"$b(>*)": {
			Name: "b",
			Args: []values.FuncArg{},
			Ret:  types.Bool,
			Body: []nodes.Node{
				&nodes.ReturnNode{
					Value: &nodes.BooleanNode{Bool: true},
				},
			},
		},
		"$c/char(>'q')": {
			Name: "c",
			Args: []values.FuncArg{},
			Ret:  types.Char,
			Body: []nodes.Node{
				&nodes.ReturnNode{
					Value: &nodes.CharNode{Char: 'q'},
				},
			},
		},
		"$f/float num/float(>add_f! num 1.0)": {
			Name: "f",
			Args: []values.FuncArg{
				{Name: "num", Type: types.Float},
			},
			Ret: types.Float,
			Body: []nodes.Node{
				&nodes.ReturnNode{
					Value: &nodes.FuncCallNode{
						Func: "add_f",
						Args: []nodes.Node{
							&nodes.SelectorNode{
								Parent: nil,
								Child:  "num",
							},
							&nodes.FloatNode{Float: 1.0},
						},
					},
				},
			},
		},
		"$i/i a/i b/i(>add_i! sub_i! a b sub_i! b a)": {
			Name: "i",
			Args: []values.FuncArg{
				{Name: "a", Type: types.Int},
				{Name: "b", Type: types.Int},
			},
			Ret: types.Int,
			Body: []nodes.Node{
				&nodes.ReturnNode{Value: &nodes.FuncCallNode{
					Func: "add_i",
					Args: []nodes.Node{
						&nodes.FuncCallNode{
							Func: "sub_i",
							Args: []nodes.Node{
								&nodes.SelectorNode{
									Parent: nil,
									Child:  "a",
								},
								&nodes.SelectorNode{
									Parent: nil,
									Child:  "b",
								},
							},
						},
						&nodes.FuncCallNode{
							Func: "sub_i",
							Args: []nodes.Node{
								&nodes.SelectorNode{
									Parent: nil,
									Child:  "b",
								},
								&nodes.SelectorNode{
									Parent: nil,
									Child:  "a",
								},
							},
						},
					},
				}},
			},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, out)
		t.Log("OK")
	}
}

func TestFuncExtern(t *testing.T) {
	p, src := makeParser("FuncExtern")

	cases := map[string]*nodes.FuncExternNode{
		"$exit code/int?": {
			Name:   "exit",
			Intern: "",
			Args: []values.FuncArg{
				{Name: "code", Type: types.Int},
			},
			Ret: types.Nothing,
		},
		"$puts/int txt/str?": {
			Name:   "puts",
			Intern: "",
			Args: []values.FuncArg{
				{Name: "txt", Type: types.String},
			},
			Ret: types.Int,
		},
		"$os/str?\"os_id\"": {
			Name:   "os",
			Intern: "os_id",
			Args:   []values.FuncArg{},
			Ret:    types.String,
		},
		"$is_64bit/bool?": {
			Name:   "is_64bit",
			Intern: "",
			Args:   []values.FuncArg{},
			Ret:    types.Bool,
		},
		"$panic?": {
			Name:   "panic",
			Intern: "",
			Args:   []values.FuncArg{},
			Ret:    types.Nothing,
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, out)
		t.Log("OK")
	}
}

func TestStruct(t *testing.T) {
	p, src := makeParser("Struct")

	cases := map[string]*nodes.StructNode{
		"@V2I(x/i y/i)": {
			Name: "V2I",
			Type: types.StructType{
				Name: "V2I",
				Fields: []types.Field{
					{Name: "x", Type: types.Int},
					{Name: "y", Type: types.Int},
				},
			},
		},
		"@V3I(x/i y/i z/i)": {
			Name: "V3I",
			Type: types.StructType{
				Name: "V3I",
				Fields: []types.Field{
					{Name: "x", Type: types.Int},
					{Name: "y", Type: types.Int},
					{Name: "z", Type: types.Int},
				},
			},
		},
		"@Position(x/f y/f z/i)": {
			Name: "Position",
			Type: types.StructType{
				Name: "Position",
				Fields: []types.Field{
					{Name: "x", Type: types.Float},
					{Name: "y", Type: types.Float},
					{Name: "z", Type: types.Int},
				},
			},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, out)
		t.Log("OK")
	}
}
