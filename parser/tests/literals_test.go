package parser_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestLiteralBool(t *testing.T) {
	p, src := makeParser("LiteralBool")

	src.Reset("*")
	n := expectValue(t, p, ast.NBool).Node
	bn := n.(ast.BoolNode)
	if bn.Value != true {
		t.Fatalf("expected true, got false")
	}

	src.Reset("/")
	n = expectValue(t, p, ast.NBool).Node
	bn = n.(ast.BoolNode)
	if bn.Value != false {
		t.Fatalf("expected false, got true")
	}
	t.Log("OK")
}

func TestLiteralChar(t *testing.T) {
	p, src := makeParser("LiteralChar")

	cases := map[string]byte{
		"'\\\\'": '\\',
		"'\\\"'": '"',
		"'\\''":  '\'',
		"'\\n'":  '\n',
		"'\\r'":  '\r',
		"'\\t'":  '\t',
	}
	// add the entirety of printable ASCII to test cases
	for c := byte(0x20); c < byte(0x7F); c++ {
		if c != '\\' && c != '\'' {
			cases[fmt.Sprintf("'%c'", c)] = c
		}
	}
	for in, out := range cases {
		t.Logf("%s -> %c", in, out)
		src.Reset(in)
		mn := expectValue(t, p, ast.NChar)
		ch := mn.Node.(ast.CharNode)
		if ch.Value != out {
			t.Fatalf("expected '%c', got '%c'", out, ch.Value)
		}
		t.Log("OK")
	}
}

func TestLiteralInt(t *testing.T) {
	p, src := makeParser("LiteralInt")

	cases := map[string]int32{
		"-1":         -1,
		"-0":         0,
		"0":          0,
		"0xA":        0xA,
		"0x123_456":  0x123_456,
		"0b10110101": 0b10110101,
		"07171":      07171,
		"0o1726354":  0o1726354,
		"1_2":        12,
	}
	// generate a bunch of completely randomised integers to make sure no
	// weirdness happens for example due to large numbers
	for i := 0; i < 20; i++ {
		num := rand.Int31()
		if num%2 == 0 {
			num = -num
		}
		cases[fmt.Sprint(num)] = num
	}

	for in, out := range cases {
		t.Logf("%s -> %d", in, out)
		src.Reset(in)
		mn := expectValue(t, p, ast.NInt)
		i := mn.Node.(ast.IntNode)
		if i.Value != out {
			t.Fatalf("expected %d, got %d", out, i.Value)
		}
		t.Log("OK")
	}
}

func TestLiteralFloat(t *testing.T) {
	p, src := makeParser("LiteralFloat")

	cases := map[string]float32{
		"1.0":         1.0,
		"-1.0":        -1.0,
		"123_456.789": 123_456.789,
		"0.0":         0.0,
	}
	for i := 0; i < 20; i++ {
		num := rand.Float32()
		if rand.Int()%2 == 0 {
			num = -num
		}
		cases[fmt.Sprint(num)] = num
	}

	for in, out := range cases {
		t.Logf("%s -> %f", in, out)
		src.Reset(in)
		mn := expectValue(t, p, ast.NFloat)
		f := mn.Node.(ast.FloatNode)
		if f.Value != out {
			t.Fatalf("expected %f, got %f", out, f.Value)
		}
		t.Log("OK")
	}
}

func TestLiteralString(t *testing.T) {
	p, src := makeParser("LiteralString")

	escaper := strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
	cases := map[string]string{
		"\"\"":     "",
		"\"\\\"\"": "\"",
		"\"\\n\"":  "\n",
	}
	// generate a bunch of random ASCII strings to ensure it's fine
	for i := 0; i < 20; i++ {
		out := ""
		for j := 0; j < 20; j++ {
			out += string(byte(randRange(0x20, 0x7F)))
		}
		in := "\"" + escaper.Replace(out) + "\""
		cases[in] = out
	}

	for in, out := range cases {
		t.Logf("%s -> %s", in, out)
		src.Reset(in)
		mn := expectValue(t, p, ast.NString)
		s := mn.Node.(ast.StringNode)
		if s.Value != out {
			t.Fatalf("expected \"%s\", got \"%s\"", out, s.Value)
		}
		t.Log("OK")
	}
}

func TestLiteralArray(t *testing.T) {
	p, src := makeParser("LiteralArray")

	cases := map[string]ast.ArrayNode{
		"[bool]()":                 arrOf(types.Bool),
		"[char]()":                 arrOf(types.Char),
		"[int]()":                  arrOf(types.Int),
		"[float]()":                arrOf(types.Float),
		"[string]()":               arrOf(types.String),
		"[b](* * /)":               arrOf(types.Bool, true, true, false),
		"[c]('a' 'b' 'c')":         arrOf(types.Char, 'a', 'b', 'c'),
		"[i](1 2 3)":               arrOf(types.Int, 1, 2, 3),
		"[f](1.2 3.4 5.6)":         arrOf(types.Float, 1.2, 3.4, 5.6),
		"[s](\"hello\" \"world\")": arrOf(types.String, "hello", "world"),
		"[](* / / *)":              arrOf(types.Bool, true, false, false, true),
		"[]('x' 'y' 'z')":          arrOf(types.Char, 'x', 'y', 'z'),
		"[](9 8 7)":                arrOf(types.Int, 9, 8, 7),
		"[](9.8 7.6 5.4)":          arrOf(types.Float, 9.8, 7.6, 5.4),
		"[](\"foo\" \"bar\")":      arrOf(types.String, "foo", "bar"),
	}

	for in, out := range cases {
		t.Logf("%s -> %s (%d elems)", in, out.Type, len(out.Elems))
		src.Reset(in)
		mn := expectValue(t, p, ast.NArray)
		a := mn.Node.(ast.ArrayNode)
		if !out.Type.Equals(a.Type) {
			t.Fatalf("expected %s, got %s", out.Type, a.Type)
		}
		if len(out.Elems) != len(a.Elems) {
			t.Fatalf("expected %d elements, got %d", len(out.Elems), len(a.Elems))
		}
		for i, ee := range out.Elems {
			compareLiteral(t, "element", ee, a.Elems[i])
		}
		t.Log("OK")
	}
}
