package parser_test

import (
	"testing"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestVarDef(t *testing.T) {
	p, src := makeParser("VarDef")

	expectAll(t, p, src, []testCase{{
		Code: "%SomeBool/bool",
		Result: ast.VarDefNode{
			Type: types.Bool,
			Var:  "SomeBool",
		}}, {
		Code: "%SomeChar/char",
		Result: ast.VarDefNode{
			Type: types.Char,
			Var:  "SomeChar",
		}}, {
		Code: "%SomeInt/int",
		Result: ast.VarDefNode{
			Type: types.Int,
			Var:  "SomeInt",
		}}, {
		Code: "%SomeFloat/float",
		Result: ast.VarDefNode{
			Type: types.Float,
			Var:  "SomeFloat",
		}}, {
		Code: "%SomeString/string",
		Result: ast.VarDefNode{
			Type: types.String,
			Var:  "SomeString",
		}},
	})
}

func TestVarSet(t *testing.T) {
	p, src := makeParser("VarSet")

	expectAll(t, p, src, []testCase{{
		Code: "%B0: *",
		Result: ast.VarSetNode{
			Var:   "B0",
			Value: ast.MetaNode{Node: ast.BoolNode{Value: true}},
		}}, {
		Code: "%C1: 'q'",
		Result: ast.VarSetNode{
			Var:   "C1",
			Value: ast.MetaNode{Node: ast.CharNode{Value: 'q'}},
		}}, {
		Code: "%I2: 123",
		Result: ast.VarSetNode{
			Var:   "I2",
			Value: ast.MetaNode{Node: ast.IntNode{Value: 123}},
		}}, {
		Code: "%F3: 12.34",
		Result: ast.VarSetNode{
			Var:   "F3",
			Value: ast.MetaNode{Node: ast.FloatNode{Value: 12.34}},
		}}, {
		Code: "%S4: \"foo\"",
		Result: ast.VarSetNode{
			Var:   "S4",
			Value: ast.MetaNode{Node: ast.StringNode{Value: "foo"}},
		}},
	})
}

func TestVarSetTyped(t *testing.T) {
	p, src := makeParser("VarSetTyped")

	expectAll(t, p, src, []testCase{{
		Code: "%ExplicitBool/b:/",
		Result: ast.VarSetTypedNode{
			Type:  types.Bool,
			Var:   "ExplicitBool",
			Value: ast.MetaNode{Node: ast.BoolNode{Value: false}},
		}}, {
		Code: "%ExplicitChar/c:'E'",
		Result: ast.VarSetTypedNode{
			Type:  types.Char,
			Var:   "ExplicitChar",
			Value: ast.MetaNode{Node: ast.CharNode{Value: 'E'}},
		}}, {
		Code: "%ExplicitInt/i:0",
		Result: ast.VarSetTypedNode{
			Type:  types.Int,
			Var:   "ExplicitInt",
			Value: ast.MetaNode{Node: ast.IntNode{Value: 0}},
		}}, {
		Code: "%ExplicitFloat/f:0.5",
		Result: ast.VarSetTypedNode{
			Type:  types.Float,
			Var:   "ExplicitFloat",
			Value: ast.MetaNode{Node: ast.FloatNode{Value: 0.5}},
		}}, {
		Code: "%ExplicitString/s:\"pp\"",
		Result: ast.VarSetTypedNode{
			Type:  types.String,
			Var:   "ExplicitString",
			Value: ast.MetaNode{Node: ast.StringNode{Value: "pp"}},
		}},
	})
}

func TestFuncDef(t *testing.T) {
	p, src := makeParser("FuncDef")

	expectAll(t, p, src, []testCase{{
		Code: "$a()",
		Result: ast.FuncDefNode{
			Name: "a",
			Body: ast.Block{},
		}}, {
		Code: "$b(>*)",
		Result: ast.FuncDefNode{
			Name: "b",
			Body: ast.Block{
				ast.MetaNode{
					Node: ast.ReturnNode{
						Value: ast.MetaNode{Node: ast.BoolNode{Value: true}},
					},
				},
			},
		}}, {
		Code: "$c/char(>'q')",
		Result: ast.FuncDefNode{
			Name: "c",
			Body: ast.Block{
				ast.MetaNode{
					Node: ast.ReturnNode{
						Value: ast.MetaNode{Node: ast.CharNode{Value: 'q'}},
					},
				},
			},
		}}, {
		Code: "$f/float num/float(>add! num 1.0)",
		Result: ast.FuncDefNode{
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
		}}, {
		Code: "$i/i a/i b/i(>add! sub! a b sub! b a)",
		Result: ast.FuncDefNode{
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
		}},
	})
}

func TestFuncShorthand(t *testing.T) {
	p, src := makeParser("FuncDef")

	expectAll(t, p, src, []testCase{{
		Code: "$Add1/int n/int: add! n 1",
		Result: ast.FuncShorthandNode{
			Name: "Add1",
			Proto: []types.Descriptor{{
				Name: "n",
				Type: types.Int,
			}},
			Ret: types.Int,
			Body: ast.MetaNode{
				Node: ast.FuncCallNode{
					Func: "add",
					Args: []ast.MetaNode{{
						Node: ast.SelectorNode{
							Parent: nil,
							Child:  "n",
						},
					}, {
						Node: ast.IntNode{
							Value: 1,
						},
					}},
				},
			},
		},
	}, {
		Code: "$SumSqr/int a/int b/int: add! mul! a a mul! b b",
		Result: ast.FuncShorthandNode{
			Name: "SumSqr",
			Proto: []types.Descriptor{{
				Name: "a",
				Type: types.Int,
			}, {
				Name: "b",
				Type: types.Int,
			}},
			Ret: types.Int,
			Body: ast.MetaNode{
				Node: ast.FuncCallNode{
					Func: "add",
					Args: []ast.MetaNode{{
						Node: ast.FuncCallNode{
							Func: "mul",
							Args: []ast.MetaNode{{
								Node: ast.SelectorNode{
									Parent: nil,
									Child:  "a",
								},
							}, {
								Node: ast.SelectorNode{
									Parent: nil,
									Child:  "a",
								},
							}},
						},
					}, {
						Node: ast.FuncCallNode{
							Func: "mul",
							Args: []ast.MetaNode{{
								Node: ast.SelectorNode{
									Parent: nil,
									Child:  "b",
								},
							}, {
								Node: ast.SelectorNode{
									Parent: nil,
									Child:  "b",
								},
							}},
						},
					}},
				},
			},
		},
	}})
}

func TestFuncExtern(t *testing.T) {
	p, src := makeParser("FuncExtern")

	expectAll(t, p, src, []testCase{{
		Code: "$exit code/int?",
		Result: ast.FuncExternNode{
			Alias: "exit",
			Name:  "",
		}}, {
		Code: "$puts/int txt/str?",
		Result: ast.FuncExternNode{
			Alias: "puts",
			Name:  "",
		}}, {
		Code: "$os/str?\"os_id\"",
		Result: ast.FuncExternNode{
			Alias: "os",
			Name:  "os_id",
		}}, {
		Code: "$is_64bit/bool?",
		Result: ast.FuncExternNode{
			Alias: "is_64bit",
			Name:  "",
		}}, {
		Code: "$panic?",
		Result: ast.FuncExternNode{
			Alias: "panic",
			Name:  "",
		}},
	})
}

func TestStructDef(t *testing.T) {
	p, src := makeParser("StructDef")

	expectAll(t, p, src, []testCase{{
		Code: "@V2I(x/i y/i)",
		Result: ast.StructDefNode{
			Name: "V2I",
			Fields: []types.Descriptor{
				{Name: "x", Type: types.Int},
				{Name: "y", Type: types.Int},
			},
		}}, {
		Code: "@V3I(x/i y/i z/i)",
		Result: ast.StructDefNode{
			Name: "V3I",
			Fields: []types.Descriptor{
				{Name: "x", Type: types.Int},
				{Name: "y", Type: types.Int},
				{Name: "z", Type: types.Int},
			},
		}}, {
		Code: "@Position(x/f y/f z/i)",
		Result: ast.StructDefNode{
			Name: "Position",
			Fields: []types.Descriptor{
				{Name: "x", Type: types.Float},
				{Name: "y", Type: types.Float},
				{Name: "z", Type: types.Int},
			},
		}},
	})
}
