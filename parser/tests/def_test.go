package parser_test

import (
	"testing"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestVarDef(t *testing.T) {
	p, src := makeParser("VarDef")

	cases := map[string]ast.VarDefNode{
		"%SomeBool/bool": {
			Type: types.Bool,
			Var:  "SomeBool",
		},
		"%SomeChar/char": {
			Type: types.Char,
			Var:  "SomeChar",
		},
		"%SomeInt/int": {
			Type: types.Int,
			Var:  "SomeInt",
		},
		"%SomeFloat/float": {
			Type: types.Float,
			Var:  "SomeFloat",
		},
		"%SomeString/string": {
			Type: types.String,
			Var:  "SomeString",
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, ast.MetaNode{Node: out})
		t.Log("OK")
	}
}

func TestVarSet(t *testing.T) {
	p, src := makeParser("VarSet")

	cases := map[string]ast.VarSetNode{
		"%B0: *": {
			Var:   "B0",
			Value: ast.MetaNode{Node: ast.BoolNode{Value: true}},
		},
		"%C1: 'q'": {
			Var:   "C1",
			Value: ast.MetaNode{Node: ast.CharNode{Value: 'q'}},
		},
		"%I2: 123": {
			Var:   "I2",
			Value: ast.MetaNode{Node: ast.IntNode{Value: 123}},
		},
		"%F3: 12.34": {
			Var:   "F3",
			Value: ast.MetaNode{Node: ast.FloatNode{Value: 12.34}},
		},
		"%S4: \"foo\"": {
			Var:   "S4",
			Value: ast.MetaNode{Node: ast.StringNode{Value: "foo"}},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, ast.MetaNode{Node: out})
		t.Log("OK")
	}
}

func TestVarSetTyped(t *testing.T) {
	p, src := makeParser("VarSetTyped")

	cases := map[string]ast.VarSetTypedNode{
		"%ExplicitBool/b:/": {
			Type:  types.Bool,
			Var:   "ExplicitBool",
			Value: ast.MetaNode{Node: ast.BoolNode{Value: false}},
		},
		"%ExplicitChar/c:'E'": {
			Type:  types.Char,
			Var:   "ExplicitChar",
			Value: ast.MetaNode{Node: ast.CharNode{Value: 'E'}},
		},
		"%ExplicitInt/i:0": {
			Type:  types.Int,
			Var:   "ExplicitInt",
			Value: ast.MetaNode{Node: ast.IntNode{Value: 0}},
		},
		"%ExplicitFloat/f:0.5": {
			Type:  types.Float,
			Var:   "ExplicitFloat",
			Value: ast.MetaNode{Node: ast.FloatNode{Value: 0.5}},
		},
		"%ExplicitString/s:\"pp\"": {
			Type:  types.String,
			Var:   "ExplicitString",
			Value: ast.MetaNode{Node: ast.StringNode{Value: "pp"}},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, ast.MetaNode{Node: out})
		t.Log("OK")
	}
}

func TestFuncDef(t *testing.T) {
	p, src := makeParser("FuncDef")

	cases := map[string]ast.FuncDefNode{
		"$a()": {
			Name: "a",
			Body: ast.Block{},
		},
		"$b(>*)": {
			Name: "b",
			Body: ast.Block{
				ast.MetaNode{
					Node: ast.ReturnNode{
						Value: ast.MetaNode{Node: ast.BoolNode{Value: true}},
					},
				},
			},
		},
		"$c/char(>'q')": {
			Name: "c",
			Body: ast.Block{
				ast.MetaNode{
					Node: ast.ReturnNode{
						Value: ast.MetaNode{Node: ast.CharNode{Value: 'q'}},
					},
				},
			},
		},
		"$f/float num/float(>add! num 1.0)": {
			Name: "f",
			Body: ast.Block{
				ast.MetaNode{
					Node: ast.ReturnNode{
						Value: ast.MetaNode{Node: ast.FuncCallNode{
							Func: "add",
							Args: []ast.MetaNode{
								{Node: ast.SelectorNode{
									Parent: nil,
									Child:  "num",
								}},
								{Node: ast.FloatNode{Value: 1.0}},
							}},
						},
					},
				},
			},
		},
		"$i/i a/i b/i(>add! sub! a b sub! b a)": {
			Name: "i",
			Body: ast.Block{
				ast.MetaNode{
					Node: ast.ReturnNode{Value: ast.MetaNode{Node: ast.FuncCallNode{
						Func: "add",
						Args: []ast.MetaNode{
							{Node: ast.FuncCallNode{
								Func: "sub",
								Args: []ast.MetaNode{
									{Node: ast.SelectorNode{
										Parent: nil,
										Child:  "a",
									}},
									{Node: ast.SelectorNode{
										Parent: nil,
										Child:  "b",
									}},
								},
							},
							},
							{Node: ast.FuncCallNode{
								Func: "sub",
								Args: []ast.MetaNode{
									{Node: ast.SelectorNode{
										Parent: nil,
										Child:  "b",
									},
									},
									{Node: ast.SelectorNode{
										Parent: nil,
										Child:  "a",
									},
									},
								},
							},
							},
						}},
					},
					},
				},
			},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, ast.MetaNode{Node: out})
		t.Log("OK")
	}
}

func TestFuncExtern(t *testing.T) {
	p, src := makeParser("FuncExtern")

	cases := map[string]ast.FuncExternNode{
		"$exit code/int?": {
			Alias: "exit",
			Name:  "",
		},
		"$puts/int txt/str?": {
			Alias: "puts",
			Name:  "",
		},
		"$os/str?\"os_id\"": {
			Alias: "os",
			Name:  "os_id",
		},
		"$is_64bit/bool?": {
			Alias: "is_64bit",
			Name:  "",
		},
		"$panic?": {
			Alias: "panic",
			Name:  "",
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, ast.MetaNode{Node: out})
		t.Log("OK")
	}
}

func TestStructDef(t *testing.T) {
	p, src := makeParser("StructDef")

	cases := map[string]ast.StructDefNode{
		"@V2I(x/i y/i)": {
			Name: "V2I",
			Fields: []types.Descriptor{
				{Name: "x", Type: types.Int},
				{Name: "y", Type: types.Int},
			},
		},
		"@V3I(x/i y/i z/i)": {
			Name: "V3I",
			Fields: []types.Descriptor{
				{Name: "x", Type: types.Int},
				{Name: "y", Type: types.Int},
				{Name: "z", Type: types.Int},
			},
		},
		"@Position(x/f y/f z/i)": {
			Name: "Position",
			Fields: []types.Descriptor{
				{Name: "x", Type: types.Float},
				{Name: "y", Type: types.Float},
				{Name: "z", Type: types.Int},
			},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, ast.MetaNode{Node: out})
		t.Log("OK")
	}
}
