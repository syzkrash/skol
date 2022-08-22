package parser_test

import (
	"fmt"
	"testing"

	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestFuncCall(t *testing.T) {
	p, src := makeParser("FuncCall")

	p.Scope.Funcs["hi"] = &values.Function{
		Name: "hi",
		Args: []values.FuncArg{},
		Ret:  types.Nothing,
	}
	p.Scope.Funcs["hello"] = &values.Function{
		Name: "hello",
		Args: []values.FuncArg{
			{Name: "who", Type: types.String},
		},
		Ret: types.String,
	}
	p.Scope.Funcs["join"] = &values.Function{
		Name: "join",
		Args: []values.FuncArg{
			{Name: "sep", Type: types.String},
			{Name: "elem", Type: types.ArrayType{Element: types.String}},
		},
		Ret: types.String,
	}
	p.Scope.Funcs["world"] = &values.Function{
		Name: "world",
		Args: []values.FuncArg{},
		Ret:  types.String,
	}

	cases := map[string]*nodes.FuncCallNode{
		"hi!": {
			Func: "hi",
			Args: []nodes.Node{}},

		"hello! \"Joe\"": {
			Func: "hello",
			Args: []nodes.Node{
				&nodes.StringNode{Str: "Joe"},
			}},

		"join! \";\" [](\"Hello\" \"world\")": {
			Func: "join",
			Args: []nodes.Node{
				&nodes.StringNode{Str: ";"},
				arrOf(types.String, "Hello", "world"),
			}},

		"join! \";\" [](hello! \"Joe\" hello! world!)": {
			Func: "join",
			Args: []nodes.Node{
				&nodes.StringNode{Str: ";"},
				&nodes.ArrayNode{
					Type: types.String,
					Elements: []nodes.Node{
						&nodes.FuncCallNode{
							Func: "hello",
							Args: []nodes.Node{
								&nodes.StringNode{Str: "Joe"},
							}},
						&nodes.FuncCallNode{
							Func: "hello",
							Args: []nodes.Node{
								&nodes.FuncCallNode{
									Func: "world",
									Args: []nodes.Node{},
								}},
						}}}}},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, out)
		t.Log("OK")
	}
}

func TestNewStruct(t *testing.T) {
	p, src := makeParser("NewStruct")

	V2I := types.MakeStruct("V2I",
		"x", types.Int,
		"y", types.Int)
	V3I := types.MakeStruct("V3I",
		"x", types.Int,
		"y", types.Int,
		"z", types.Int)
	Rect := types.MakeStruct("Rect",
		"pos", V2I,
		"extents", V2I)
	Cuboid := types.MakeStruct("Cuboid",
		"pos", V3I,
		"extents", V3I)

	p.Scope.Types["V2I"] = V2I
	p.Scope.Types["V3I"] = V3I
	p.Scope.Types["Rect"] = Rect
	p.Scope.Types["Cuboid"] = Cuboid

	cases := map[string]*nodes.NewStructNode{
		"@V2I 6 9": {
			Type: V2I,
			Args: []nodes.Node{
				&nodes.IntegerNode{Int: 6},
				&nodes.IntegerNode{Int: 9},
			}},

		"@V3I 4 2 0": {
			Type: V3I,
			Args: []nodes.Node{
				&nodes.IntegerNode{Int: 4},
				&nodes.IntegerNode{Int: 2},
				&nodes.IntegerNode{Int: 0},
			}},

		"@Rect @V2I 12 34 @V2I 1 2": {
			Type: Rect,
			Args: []nodes.Node{
				&nodes.NewStructNode{
					Type: V2I,
					Args: []nodes.Node{
						&nodes.IntegerNode{Int: 12},
						&nodes.IntegerNode{Int: 34},
					}},
				&nodes.NewStructNode{
					Type: V2I,
					Args: []nodes.Node{
						&nodes.IntegerNode{Int: 1},
						&nodes.IntegerNode{Int: 2},
					}}}},

		"@Cuboid @V3I 12 34 56 @V3I 3 3 3": {
			Type: Cuboid,
			Args: []nodes.Node{
				&nodes.NewStructNode{
					Type: V3I,
					Args: []nodes.Node{
						&nodes.IntegerNode{Int: 12},
						&nodes.IntegerNode{Int: 34},
						&nodes.IntegerNode{Int: 56},
					}},
				&nodes.NewStructNode{
					Type: V3I,
					Args: []nodes.Node{
						&nodes.IntegerNode{Int: 3},
						&nodes.IntegerNode{Int: 3},
						&nodes.IntegerNode{Int: 3},
					}}}},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		n := expectValue(t, p, nodes.NdNewStruct)
		compare(t, "NewStruct", out, n)
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

	p.Scope.Vars["a"] = &nodes.VarDefNode{
		VarType: TA,
		Var:     "a",
	}
	p.Scope.Vars["z"] = &nodes.VarDefNode{
		VarType: Z,
		Var:     "z",
	}

	cases := map[string][]nodes.SelectorElem{
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
		"z#y": {
			{Name: "z"},
			{Idx: &nodes.SelectorNode{
				Parent: nil,
				Child:  "y",
			}},
		},
		"z#y#x": {
			{Name: "z"},
			{Idx: &nodes.SelectorNode{
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
		n, err := p.Value()
		if err != nil {
			if pe, ok := err.(common.Printable); ok {
				pe.Print()
			}
			t.Fatal(err)
		}
		s, ok := n.(nodes.Selector)
		if !ok {
			t.Fatalf("%d node is not a selector", n.Kind())
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
			compare(t, fmt.Sprintf("selector element %d", i), ee.Idx, ge.Idx)
		}
	}
}
