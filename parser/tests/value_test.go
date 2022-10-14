package parser_test

import (
	"fmt"
	"testing"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestFuncCall(t *testing.T) {
	p, src := makeParser("FuncCall")

	p.Tree.Funcs["hi"] = ast.Func{
		Name: "hi",
		Args: []types.Descriptor{},
		Ret:  types.Nothing,
	}
	p.Tree.Funcs["hello"] = ast.Func{
		Name: "hello",
		Args: []types.Descriptor{
			{Name: "who", Type: types.String},
		},
		Ret: types.String,
	}
	p.Tree.Funcs["join"] = ast.Func{
		Name: "join",
		Args: []types.Descriptor{
			{Name: "sep", Type: types.String},
			{Name: "elem", Type: types.ArrayType{Element: types.String}},
		},
		Ret: types.String,
	}
	p.Tree.Funcs["world"] = ast.Func{
		Name: "world",
		Args: []types.Descriptor{},
		Ret:  types.String,
	}

	cases := map[string]ast.FuncCallNode{
		"hi!": {
			Func: "hi",
			Args: []ast.MetaNode{}},

		"hello! \"Joe\"": {
			Func: "hello",
			Args: []ast.MetaNode{
				{Node: ast.StringNode{Value: "Joe"}},
			}},

		"join! \";\" [](\"Hello\" \"world\")": {
			Func: "join",
			Args: []ast.MetaNode{
				{Node: ast.StringNode{Value: ";"}},
				{Node: arrOf(types.String, "Hello", "world")},
			}},

		"join! \";\" [](hello! \"Joe\" hello! world!)": {
			Func: "join",
			Args: []ast.MetaNode{
				{Node: ast.StringNode{Value: ";"}},
				{Node: ast.ArrayNode{
					Type: types.ArrayType{Element: types.String},
					Elems: []ast.MetaNode{
						{Node: ast.FuncCallNode{
							Func: "hello",
							Args: []ast.MetaNode{
								{Node: ast.StringNode{Value: "Joe"}},
							}}},
						{Node: ast.FuncCallNode{
							Func: "hello",
							Args: []ast.MetaNode{
								{Node: ast.FuncCallNode{
									Func: "world",
									Args: []ast.MetaNode{},
								}}}},
						}}}}}},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, ast.MetaNode{Node: out})
		t.Log("OK")
	}
}

func TestStruct(t *testing.T) {
	p, src := makeParser("Struct")

	V2I := types.MakeStruct("V2I",
		"x", types.Int,
		"y", types.Int).(types.StructType)
	V3I := types.MakeStruct("V3I",
		"x", types.Int,
		"y", types.Int,
		"z", types.Int).(types.StructType)
	Rect := types.MakeStruct("Rect",
		"pos", V2I,
		"extents", V2I).(types.StructType)
	Cuboid := types.MakeStruct("Cuboid",
		"pos", V3I,
		"extents", V3I).(types.StructType)

	p.Scope.Types["V2I"] = V2I
	p.Scope.Types["V3I"] = V3I
	p.Scope.Types["Rect"] = Rect
	p.Scope.Types["Cuboid"] = Cuboid

	cases := map[string]ast.StructNode{
		"@V2I 6 9": {
			Type: V2I,
			Args: []ast.MetaNode{
				{Node: ast.IntNode{Value: 6}},
				{Node: ast.IntNode{Value: 9}},
			}},

		"@V3I 4 2 0": {
			Type: V3I,
			Args: []ast.MetaNode{
				{Node: ast.IntNode{Value: 4}},
				{Node: ast.IntNode{Value: 2}},
				{Node: ast.IntNode{Value: 0}},
			}},

		"@Rect @V2I 12 34 @V2I 1 2": {
			Type: Rect,
			Args: []ast.MetaNode{
				{Node: ast.StructNode{
					Type: V2I,
					Args: []ast.MetaNode{
						{Node: ast.IntNode{Value: 12}},
						{Node: ast.IntNode{Value: 34}},
					}}},
				{Node: ast.StructNode{
					Type: V2I,
					Args: []ast.MetaNode{
						{Node: ast.IntNode{Value: 1}},
						{Node: ast.IntNode{Value: 2}},
					}}}}},

		"@Cuboid @V3I 12 34 56 @V3I 3 3 3": {
			Type: Cuboid,
			Args: []ast.MetaNode{
				{Node: ast.StructNode{
					Type: V3I,
					Args: []ast.MetaNode{
						{Node: ast.IntNode{Value: 12}},
						{Node: ast.IntNode{Value: 34}},
						{Node: ast.IntNode{Value: 56}},
					}}},
				{Node: ast.StructNode{
					Type: V3I,
					Args: []ast.MetaNode{
						{Node: ast.IntNode{Value: 3}},
						{Node: ast.IntNode{Value: 3}},
						{Node: ast.IntNode{Value: 3}},
					}}}}},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		n := expectValue(t, p, ast.NStruct)
		compare(t, "NewStruct", ast.MetaNode{Node: out}, n)
		t.Log("OK")
	}
}

func TestSelector(t *testing.T) {
	p, src := makeParser("Selector")

	X := types.MakeStruct("X")
	Y := types.MakeStruct("Y",
		"x", X)
	Z := types.ArrayType{Element: Y}

	TC := types.MakeStruct("TC")
	TB := types.MakeStruct("TB",
		"c", TC)
	TA := types.MakeStruct("TA",
		"b", TB,
		"c", TC)

	p.Scope.Types["TA"] = TA
	p.Scope.Types["TB"] = TB
	p.Scope.Types["TC"] = TC
	p.Scope.Types["Y"] = Y
	p.Scope.Types["X"] = X

	p.Scope.Vars["a"] = ast.StructNode{
		Type: TA.(types.StructType),
		Args: []ast.MetaNode{},
	}
	p.Scope.Vars["z"] = ast.ArrayNode{
		Type:  Z,
		Elems: []ast.MetaNode{},
	}

	cases := map[string][]ast.SelectorElem{
		"a": {
			{Name: "a"},
		},
		"a#b": {
			{Name: "a"},
			{Name: "b"},
		},
		"a#b#c": {
			{Name: "a"},
			{Name: "b"},
			{Name: "c"},
		},
		"z#[y]": {
			{Name: "z"},
			{IdxS: ast.SelectorNode{
				Parent: nil,
				Child:  "y",
			}},
		},
		"z#[y]#x": {
			{Name: "z"},
			{IdxS: ast.SelectorNode{
				Parent: nil,
				Child:  "y",
			}},
			{Name: "x"},
		},
		"a#@TB": {
			{Name: "a"},
			{Cast: TB},
		},
		"a#@TB#c": {
			{Name: "a"},
			{Cast: TB},
			{Name: "c"},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		mn, err := p.ParseValue()
		if err != nil {
			if pe, ok := err.(common.Printable); ok {
				pe.Print()
			}
			t.Fatal(err)
		}
		s, ok := mn.Node.(ast.Selector)
		if !ok {
			t.Fatalf("%d node is not a selector", mn.Node.Kind())
		}
		p := s.Path()
		if len(out) != len(p) {
			t.Fatalf("expected %d path elements, got %d", len(out), len(p))
		}
		for i, ee := range out {
			ge := p[i]
			if ee.Name != "" {
				if ge.Name != ee.Name {
					t.Fatalf("selector element %d: expected `%s` name, got `%s`", i, ee.Name, ge.Name)
				}
				continue
			}
			if ee.Cast != nil {
				if !ee.Cast.Equals(ge.Cast) {
					t.Fatalf("selector element %d: expected %s cast, got %s", i, ee.Cast, ge.Cast)
				}
				continue
			}
			compare(t, fmt.Sprintf("selector element %d", i), ast.MetaNode{Node: ee.IdxS}, ast.MetaNode{Node: ge.IdxS})
		}
	}
}
