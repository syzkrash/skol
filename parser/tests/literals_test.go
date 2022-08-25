package parser_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestLiteralBool(t *testing.T) {
	p, src := makeParser("LiteralBool")

	src.Reset("*")
	n := expectValue(t, p, nodes.NdBoolean)
	bn := n.(*nodes.BooleanNode)
	if bn.Bool != true {
		t.Fatalf("expected true, got false")
	}

	src.Reset("/")
	n = expectValue(t, p, nodes.NdBoolean)
	bn = n.(*nodes.BooleanNode)
	if bn.Bool != false {
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
		n := expectValue(t, p, nodes.NdChar)
		ch := n.(*nodes.CharNode)
		if ch.Char != out {
			t.Fatalf("expected '%c', got '%c'", out, ch.Char)
		}
		t.Log("OK")
	}
}

func TestLiteralInteger(t *testing.T) {
	p, src := makeParser("LiteralInteger")

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
		n := expectValue(t, p, nodes.NdInteger)
		i := n.(*nodes.IntegerNode)
		if i.Int != out {
			t.Fatalf("expected %d, got %d", out, i.Int)
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
		n := expectValue(t, p, nodes.NdFloat)
		f := n.(*nodes.FloatNode)
		if f.Float != out {
			t.Fatalf("expected %f, got %f", out, f.Float)
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
		n := expectValue(t, p, nodes.NdString)
		s := n.(*nodes.StringNode)
		if s.Str != out {
			t.Fatalf("expected \"%s\", got \"%s\"", out, s.Str)
		}
		t.Log("OK")
	}
}

func TestLiteralArray(t *testing.T) {
	p, src := makeParser("LiteralArray")

	cases := map[string]*nodes.ArrayNode{
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
		t.Logf("%s -> %s (%d elems)", in, out.Type, len(out.Elements))
		src.Reset(in)
		n := expectValue(t, p, nodes.NdArray)
		a := n.(*nodes.ArrayNode)
		if !out.Type.Equals(a.Type) {
			t.Fatalf("expected %s, got %s", out.Type, a.Type)
		}
		if len(out.Elements) != len(a.Elements) {
			t.Fatalf("expected %d elements, got %d", len(out.Elements), len(a.Elements))
		}
		for i, ee := range out.Elements {
			compareLiteral(t, "element", ee, a.Elements[i])
		}
		t.Log("OK")
	}
}
