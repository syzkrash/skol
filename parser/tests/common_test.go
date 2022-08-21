package parser_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values/types"
)

func makeParser(test string) (*parser.Parser, *strings.Reader) {
	r := strings.NewReader("")
	return parser.NewParser("Test"+test, r, "test"), r
}

func expect(t *testing.T, p *parser.Parser, exp nodes.Node) {
	got, err := p.Next()
	if err != nil {
		if pe, ok := err.(common.Printable); ok {
			pe.Print()
		}
		t.Fatal(err)
	}
	t.Logf("%+v", got)
	compare(t, "expect", exp, got)
}

func expectValue(t *testing.T, p *parser.Parser, k nodes.NodeKind) nodes.Node {
	n, err := p.Value()
	if err != nil {
		if pe, ok := err.(common.Printable); ok {
			pe.Print()
		}
		t.Fatal(err)
	}
	if n.Kind() != k {
		t.Fatalf("expected %s, got %s", k, n.Kind())
	}
	return n
}

func compareLiteral(t *testing.T, note string, exp, got nodes.Node) {
	if exp.Kind() != got.Kind() {
		t.Fatalf("%s: expected %s, got %s", note, exp.Kind(), got.Kind())
	}
	var ev, gv any
	switch exp.Kind() {
	case nodes.NdBoolean:
		eb := exp.(*nodes.BooleanNode)
		gb := got.(*nodes.BooleanNode)
		ev = eb.Bool
		gv = gb.Bool
	case nodes.NdChar:
		ec := exp.(*nodes.CharNode)
		gc := got.(*nodes.CharNode)
		ev = ec.Char
		gv = gc.Char
	case nodes.NdInteger:
		ei := exp.(*nodes.IntegerNode)
		gi := got.(*nodes.IntegerNode)
		ev = ei.Int
		gv = gi.Int
	case nodes.NdFloat:
		ef := exp.(*nodes.FloatNode)
		gf := got.(*nodes.FloatNode)
		ev = ef.Float
		gv = gf.Float
	case nodes.NdString:
		es := exp.(*nodes.StringNode)
		gs := got.(*nodes.StringNode)
		ev = es.Str
		gv = gs.Str
	}
	if ev != gv {
		t.Fatalf("expected %v, got %v", ev, gv)
	}
}

func compare(t *testing.T, note string, exp, got nodes.Node) {
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
	case nodes.NdBoolean, nodes.NdChar, nodes.NdInteger, nodes.NdFloat, nodes.NdString:
		compareLiteral(t, "literal", exp, got)
	case nodes.NdSelector, nodes.NdIndex, nodes.NdTypecast:
		es := exp.(nodes.Selector)
		gs := exp.(nodes.Selector)
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
			if ee.Idx == nil && ge.Idx != nil {
				t.Fatalf("%s: expected nil index, got %s", note, ge.Idx.Kind())
			}
			if ee.Idx != nil && ge.Idx == nil {
				t.Fatalf("%s: expected %s index, got nil", note, ee.Idx.Kind())
			}
			if ee.Idx != nil && ge.Idx != nil {
				if ee.Idx.Kind() != ge.Idx.Kind() {
					t.Fatalf("%s: expected %s index, got %s", note, ee.Idx.Kind(), ge.Idx.Kind())
				}
			}
			compare(t, note+": Index", ee.Idx, ge.Idx)
		}
	case nodes.NdReturn:
		er := exp.(*nodes.ReturnNode)
		gr := exp.(*nodes.ReturnNode)
		compare(t, "return", er.Value, gr.Value)
	case nodes.NdIf:
		eif := exp.(*nodes.IfNode)
		gif := got.(*nodes.IfNode)
		compareLiteral(t, note+": Condition", eif.Condition, gif.Condition)
		if len(eif.IfBlock) != len(gif.IfBlock) {
			t.Fatalf("%s: expected %d nodes in IfBlock, got %d", note, len(eif.IfBlock), len(gif.IfBlock))
		}
		for i, n := range eif.IfBlock {
			compare(t, note+": IfBlock", n, gif.IfBlock[i])
		}
		if len(eif.ElseBlock) != len(gif.ElseBlock) {
			t.Fatalf("%s: expected %d nodes in ElseBlock, got %d", note, len(eif.IfBlock), len(gif.IfBlock))
		}
		for i, n := range eif.ElseBlock {
			compare(t, note+": ElseBlock", n, gif.ElseBlock[i])
		}
		if len(eif.ElseIfNodes) != len(gif.ElseIfNodes) {
			t.Fatalf("%s: expected %d branches in ElseIfNodes, got %d", note, len(eif.IfBlock), len(gif.IfBlock))
		}
		for i, b := range eif.ElseIfNodes {
			compare(t, note+": ElseIfNode", b.Condition, gif.ElseIfNodes[i].Condition)
			for j, n := range b.Block {
				compare(t, note+": ElseIfNode Block", n, gif.ElseIfNodes[i].Block[j])
			}
		}
	case nodes.NdFuncCall:
		efc := exp.(*nodes.FuncCallNode)
		gfc := exp.(*nodes.FuncCallNode)
		if efc.Func != gfc.Func {
			t.Fatalf("%s: expected `%s` function, got `%s`", note, efc.Func, gfc.Func)
		}
		if len(efc.Args) != len(gfc.Args) {
			t.Fatalf("%s: expected %d arguments, got %d", note, len(efc.Args), len(gfc.Args))
		}
		for i, ea := range efc.Args {
			compare(t, note+": arguments", ea, gfc.Args[i])
		}
	case nodes.NdWhile:
		ew := exp.(*nodes.WhileNode)
		gw := exp.(*nodes.WhileNode)
		compare(t, note+": Condition", ew.Condition, gw.Condition)
		for i, en := range ew.Body {
			compare(t, note+": Body", en, gw.Body[i])
		}
	case nodes.NdVarDef:
		ev := exp.(*nodes.VarDefNode)
		gv := got.(*nodes.VarDefNode)
		if ev.Var != gv.Var {
			t.Fatalf("%s: expected `%s` variable, got `%s`", note, ev.Var, gv.Var)
		}
		if !ev.VarType.Equals(gv.VarType) {
			t.Fatalf("%s: expected %s variable type, got %s", note, ev.VarType, gv.VarType)
		}
		compare(t, note+": variable value", ev.Value, gv.Value)
	case nodes.NdFuncDef:
		ef := exp.(*nodes.FuncDefNode)
		gf := exp.(*nodes.FuncDefNode)
		if ef.Name != gf.Name {
			t.Fatalf("%s: expected `%s` function, got `%s`", note, ef.Name, gf.Name)
		}
		if len(ef.Args) != len(gf.Args) {
			t.Fatalf("%s: expected %d arguments, got %d", note, len(ef.Args), len(gf.Args))
		}
		for i, ea := range ef.Args {
			ga := gf.Args[i]
			if ea.Name != ga.Name {
				t.Fatalf("%s: argument %d: expected `%s` argument, got `%s`", note, i, ea.Name, ga.Name)
			}
			if !ea.Type.Equals(ga.Type) {
				t.Fatalf("%s: argument %d: expected %s type, got %s", note, i, ea.Type, ga.Type)
			}
		}
		if !ef.Ret.Equals(gf.Ret) {
			t.Fatalf("%s: expected %s return type, got %s", note, ef.Ret, gf.Ret)
		}
		if len(ef.Body) != len(gf.Body) {
			t.Fatalf("%s: expected %d body nodes, got %d", note, len(ef.Body), len(gf.Body))
		}
		for i, en := range ef.Body {
			gn := gf.Body[i]
			compare(t, fmt.Sprintf("%s: body node %d", note, i), en, gn)
		}
	case nodes.NdFuncExtern:
		ee := exp.(*nodes.FuncExternNode)
		ge := exp.(*nodes.FuncExternNode)
		if ee.Name != ge.Name {
			t.Fatalf("%s: expected `%s` extern, got `%s`", note, ee.Name, ge.Name)
		}
		if ee.Intern != ge.Intern {
			t.Fatalf("%s: expected `%s` intern, got `%s`", note, ee.Intern, ge.Intern)
		}
		if len(ee.Args) != len(ge.Args) {
			t.Fatalf("%s: expected %d arguments, got %d", note, len(ee.Args), len(ge.Args))
		}
		for i, ea := range ee.Args {
			ga := ge.Args[i]
			if ea.Name != ga.Name {
				t.Fatalf("%s: argument %d: expected `%s` name, got `%s`", note, i, ea.Name, ga.Name)
			}
			if !ea.Type.Equals(ga.Type) {
				t.Fatalf("%s: argument %d: expected %s type, got %s", note, i, ea.Type, ga.Type)
			}
		}
		if !ee.Ret.Equals(ge.Ret) {
			t.Fatalf("%s: expected %s return type, got %s", note, ee.Ret, ge.Ret)
		}
	case nodes.NdStruct:
		es := exp.(*nodes.StructNode)
		gs := got.(*nodes.StructNode)
		if es.Name != gs.Name {
			t.Fatalf("%s: expected `%s` structure name, got `%s`", note, es.Name, gs.Name)
		}
		if es.Type.Prim() != gs.Type.Prim() {
			t.Fatalf("%s: expected %v type primitive, got %v", note, es.Type.Prim(), gs.Type.Prim())
		}
		est := es.Type.(types.StructType)
		gst := gs.Type.(types.StructType)
		if est.Name != gst.Name {
			t.Fatalf("%s: type: expected `%s` name, got %s", note, est.Name, gst.Name)
		}
		if len(est.Fields) != len(gst.Fields) {
			t.Fatalf("%s: type: expected %d fields, got %d", note, len(est.Fields), len(gst.Fields))
		}
		for i, ef := range est.Fields {
			gf := gst.Fields[i]
			if ef.Name != gf.Name {
				t.Fatalf("%s: field `%s`: expected `%s` name", note, gf.Name, ef.Name)
			}
			if !ef.Type.Equals(gf.Type) {
				t.Fatalf("%s: field `%s`: expected %s type, got %s", note, gf.Name, ef.Type, gf.Type)
			}
		}
	default:
		panic(fmt.Sprintf("compare() call on unexpected node: %s", exp.Kind()))
	}
}

func randRange(min, max int) int {
	return min + rand.Intn(max-min)
}

func arrOf(element types.Type, elements ...any) *nodes.ArrayNode {
	enodes := []nodes.Node{}
	for _, e := range elements {
		switch e := e.(type) {
		case bool:
			enodes = append(enodes, &nodes.BooleanNode{Bool: e})
		case rune:
			enodes = append(enodes, &nodes.CharNode{Char: byte(e)})
		case int:
			enodes = append(enodes, &nodes.IntegerNode{Int: int32(e)})
		case float64:
			enodes = append(enodes, &nodes.FloatNode{Float: float32(e)})
		case string:
			enodes = append(enodes, &nodes.StringNode{Str: e})
		default:
			panic(fmt.Sprintf("unhandled value of type %s", reflect.ValueOf(e).Type().Name()))
		}
	}
	return &nodes.ArrayNode{
		Type:     element,
		Elements: enodes,
	}
}
