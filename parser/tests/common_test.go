package parser_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/values/types"
)

func makeParser(test string) (*parser.Parser, *strings.Reader) {
	r := strings.NewReader("")
	return parser.NewParser("Test"+test, r, "test"), r
}

func expect(t *testing.T, p *parser.Parser, exp ast.MetaNode) {
	got, err := p.TopLevel()
	if err != nil {
		if pe, ok := err.(common.Printable); ok {
			pe.Print()
		}
		t.Fatal(err)
	}
	t.Logf("%+v", got)
	compare(t, "expect", exp, got)
}

func expectValue(t *testing.T, p *parser.Parser, k ast.NodeKind) ast.MetaNode {
	mn, err := p.Value()
	if err != nil {
		if pe, ok := err.(common.Printable); ok {
			pe.Print()
		}
		t.Fatal(err)
	}
	if mn.Node.Kind() != k {
		t.Fatalf("expected %s, got %s", k, mn.Node.Kind())
	}
	return mn
}

func compareLiteral(t *testing.T, note string, mexp, mgot ast.MetaNode) {
	exp := mexp.Node
	got := mgot.Node
	if exp.Kind() != got.Kind() {
		t.Fatalf("%s: expected %s, got %s", note, exp.Kind(), got.Kind())
	}
	var ev, gv any
	switch exp.Kind() {
	case ast.NBool:
		eb := exp.(ast.BoolNode)
		gb := got.(ast.BoolNode)
		ev = eb.Value
		gv = gb.Value
	case ast.NChar:
		ec := exp.(ast.CharNode)
		gc := got.(ast.CharNode)
		ev = ec.Value
		gv = gc.Value
	case ast.NInt:
		ei := exp.(ast.IntNode)
		gi := got.(ast.IntNode)
		ev = ei.Value
		gv = gi.Value
	case ast.NFloat:
		ef := exp.(ast.FloatNode)
		gf := got.(ast.FloatNode)
		ev = ef.Value
		gv = gf.Value
	case ast.NString:
		es := exp.(ast.StringNode)
		gs := got.(ast.StringNode)
		ev = es.Value
		gv = gs.Value
	}
	if ev != gv {
		t.Fatalf("expected %v, got %v", ev, gv)
	}
}

func compare(t *testing.T, note string, mexp, mgot ast.MetaNode) {
	exp := mexp.Node
	got := mgot.Node
	if exp == nil && got != nil {
		t.Fatalf("%s: expected nil, got %s", note, got.Kind())
	}
	if exp != nil && got == nil {
		t.Fatalf("%s: expected %s, got nil", note, exp.Kind())
	}
	if exp == nil && got == nil {
		return
	}
	if got.Kind() != exp.Kind() {
		t.Fatalf("%s: expected %s, got %s", note, exp.Kind(), got.Kind())
	}
	switch exp.Kind() {
	case ast.NBool, ast.NChar, ast.NInt, ast.NFloat, ast.NString:
		compareLiteral(t, "literal", mexp, mgot)
	case ast.NSelector, ast.NIndexConst, ast.NIndexSelector, ast.NTypecast:
		es := exp.(ast.Selector)
		gs := exp.(ast.Selector)
		ep := es.Path()
		gp := gs.Path()
		if len(ep) != len(gp) {
			t.Fatalf("%s: expected %d selector elements, got %d", note, len(ep), len(gp))
		}
		for i, ee := range ep {
			ge := gp[i]
			if ee.Cast != ge.Cast {
				t.Fatalf("%s: expected %s cast, got %s", note, ee.Cast, ge.Cast)
			}
			if ee.Name != ge.Name {
				t.Fatalf("%s: expected `%s` field, got `%s`", note, ee.Name, ge.Name)
			}
			if ee.IdxS == nil && ge.IdxS != nil {
				t.Fatalf("%s: expected nil index, got %s", note, ge.IdxS.Kind())
			}
			if ee.IdxS != nil && ge.IdxS == nil {
				t.Fatalf("%s: expected %s index, got nil", note, ee.IdxS.Kind())
			}
			if ee.IdxS != nil && ge.IdxS != nil {
				if ee.IdxS.Kind() != ge.IdxS.Kind() {
					t.Fatalf("%s: expected %s index, got %s", note, ee.IdxS.Kind(), ge.IdxS.Kind())
				}
			}
			compare(t, note+": IndexS", ast.MetaNode{Node: ee.IdxS}, ast.MetaNode{Node: ge.IdxS})
			if ee.IdxC != ge.IdxC {
				t.Fatalf("%s: expected %d index, got %d", note, ee.IdxC, ge.IdxC)
			}
		}
	case ast.NReturn:
		er := exp.(ast.ReturnNode)
		gr := exp.(ast.ReturnNode)
		compare(t, "return", er.Value, gr.Value)
	case ast.NIf:
		eif := exp.(ast.IfNode)
		gif := got.(ast.IfNode)
		compareLiteral(t, note+": Cond", eif.Main.Cond, gif.Main.Cond)
		if len(eif.Main.Block) != len(gif.Main.Block) {
			t.Fatalf("%s: expected %d nodes in IfBlock, got %d", note, len(eif.Main.Block), len(gif.Main.Block))
		}
		for i, n := range eif.Main.Block {
			compare(t, note+": IfBlock", n, gif.Main.Block[i])
		}
		if len(eif.Else) != len(gif.Else) {
			t.Fatalf("%s: expected %d nodes in Else, got %d", note, len(eif.Else), len(gif.Else))
		}
		for i, n := range eif.Else {
			compare(t, note+": Else", n, gif.Else[i])
		}
		if len(eif.Other) != len(gif.Other) {
			t.Fatalf("%s: expected %d branches in Other, got %d", note, len(eif.Other), len(gif.Other))
		}
		for i, b := range eif.Other {
			compare(t, note+": Other Cond", b.Cond, gif.Other[i].Cond)
			for j, n := range b.Block {
				compare(t, note+": Other Block", n, gif.Other[i].Block[j])
			}
		}
	case ast.NFuncCall:
		efc := exp.(ast.FuncCallNode)
		gfc := exp.(ast.FuncCallNode)
		if efc.Func != gfc.Func {
			t.Fatalf("%s: expected `%s` function, got `%s`", note, efc.Func, gfc.Func)
		}
		if len(efc.Args) != len(gfc.Args) {
			t.Fatalf("%s: expected %d arguments, got %d", note, len(efc.Args), len(gfc.Args))
		}
		for i, ea := range efc.Args {
			compare(t, note+": arguments", ea, gfc.Args[i])
		}
	case ast.NWhile:
		ew := exp.(ast.WhileNode)
		gw := exp.(ast.WhileNode)
		compare(t, note+": Cond", ew.Cond, gw.Cond)
		for i, en := range ew.Block {
			compare(t, note+": Body", en, gw.Block[i])
		}
	case ast.NVarDef:
		ev := exp.(ast.VarDefNode)
		gv := got.(ast.VarDefNode)
		if ev.Var != gv.Var {
			t.Fatalf("%s: expected `%s` variable, got `%s`", note, ev.Var, gv.Var)
		}
		if !ev.Type.Equals(gv.Type) {
			t.Fatalf("%s: expected %s variable type, got %s", note, ev.Type, gv.Type)
		}
	case ast.NVarSet:
		ev := exp.(ast.VarSetNode)
		gv := got.(ast.VarSetNode)
		if ev.Var != gv.Var {
			t.Fatalf("%s: expected `%s` variable, got `%s`", note, ev.Var, gv.Var)
		}
		compare(t, note+": variable value", ev.Value, gv.Value)
	case ast.NVarSetTyped:
		ev := exp.(ast.VarSetTypedNode)
		gv := got.(ast.VarSetTypedNode)
		if ev.Var != gv.Var {
			t.Fatalf("%s: expected `%s` variable, got `%s`", note, ev.Var, gv.Var)
		}
		if !ev.Type.Equals(gv.Type) {
			t.Fatalf("%s: expected %s variable type, got %s", note, ev.Type, gv.Type)
		}
		compare(t, note+": variable value", ev.Value, gv.Value)
	case ast.NFuncDef:
		ef := exp.(ast.FuncDefNode)
		gf := exp.(ast.FuncDefNode)
		if ef.Name != gf.Name {
			t.Fatalf("%s: expected `%s` function, got `%s`", note, ef.Name, gf.Name)
		}
		if len(ef.Body) != len(gf.Body) {
			t.Fatalf("%s: expected %d body nodes, got %d", note, len(ef.Body), len(gf.Body))
		}
		for i, en := range ef.Body {
			gn := gf.Body[i]
			compare(t, fmt.Sprintf("%s: body node %d", note, i), en, gn)
		}
	case ast.NFuncExtern:
		ee := exp.(ast.FuncExternNode)
		ge := exp.(ast.FuncExternNode)
		if ee.Alias != ge.Alias {
			t.Fatalf("%s: expected `%s` extern alias, got `%s`", note, ee.Alias, ge.Alias)
		}
		if ee.Name != ge.Name {
			t.Fatalf("%s: expected `%s` intern name, got `%s`", note, ee.Name, ge.Name)
		}
	case ast.NStructDef:
		es := exp.(ast.StructDefNode)
		gs := got.(ast.StructDefNode)
		if es.Name != gs.Name {
			t.Fatalf("%s: expected `%s` structure name, got `%s`", note, es.Name, gs.Name)
		}
		if len(es.Fields) != len(gs.Fields) {
			t.Fatalf("%s: type: expected %d fields, got %d", note, len(es.Fields), len(gs.Fields))
		}
		for i, ef := range es.Fields {
			gf := gs.Fields[i]
			if ef.Name != gf.Name {
				t.Fatalf("%s: field `%s`: expected `%s` name", note, gf.Name, ef.Name)
			}
			if !ef.Type.Equals(gf.Type) {
				t.Fatalf("%s: field `%s`: expected %s type, got %s", note, gf.Name, ef.Type, gf.Type)
			}
		}
	case ast.NArray:
		ea := exp.(ast.ArrayNode)
		ga := got.(ast.ArrayNode)
		if !ea.Type.Equals(ga.Type) {
			t.Fatalf("%s: expected array of %s, got %s", note, ea.Type, ga.Type)
		}
		if len(ea.Elems) != len(ga.Elems) {
			t.Fatalf("%s: expected %d array elements, got %d", note, len(ea.Elems), len(ga.Elems))
		}
		for i, ee := range ea.Elems {
			ge := ga.Elems[i]
			compare(t, fmt.Sprintf("%s: element %d", note, i), ee, ge)
		}
	case ast.NStruct:
		es := exp.(ast.StructNode)
		gs := exp.(ast.StructNode)
		if !es.Type.Equals(gs.Type) {
			t.Fatalf("%s: expected %s, got %s", note, es.Type, gs.Type)
		}
		if len(es.Args) != len(gs.Args) {
			t.Fatalf("%s: expected %d arguments, got %d", note, len(es.Args), len(gs.Args))
		}
		for i, ea := range es.Args {
			ga := gs.Args[i]
			compare(t, fmt.Sprintf("%s: argument %d", note, i), ea, ga)
		}
	default:
		panic(fmt.Sprintf("compare() call on unexpected node: %s", exp.Kind()))
	}
}

func randRange(min, max int) int {
	return min + rand.Intn(max-min)
}

func arrOf(element types.Type, elements ...any) ast.ArrayNode {
	enodes := []ast.MetaNode{}
	for _, e := range elements {
		switch e := e.(type) {
		case bool:
			enodes = append(enodes, ast.MetaNode{Node: ast.BoolNode{Value: e}})
		case rune:
			enodes = append(enodes, ast.MetaNode{Node: ast.CharNode{Value: byte(e)}})
		case int:
			enodes = append(enodes, ast.MetaNode{Node: ast.IntNode{Value: int32(e)}})
		case float64:
			enodes = append(enodes, ast.MetaNode{Node: ast.FloatNode{Value: float32(e)}})
		case string:
			enodes = append(enodes, ast.MetaNode{Node: ast.StringNode{Value: e}})
		default:
			panic(fmt.Sprintf("unhandled value of type %s", reflect.ValueOf(e).Type().Name()))
		}
	}
	return ast.ArrayNode{
		Type:  types.ArrayType{Element: element},
		Elems: enodes,
	}
}
