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
