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
	p, src := makeParser(t, "LiteralBool")

	expectAllValues(t, p, src, []testCase{{
		Code:   "*",
		Result: ast.BoolNode{Value: true},
	}, {
		Code:   "/",
		Result: ast.BoolNode{Value: false},
	}})
}

func TestLiteralChar(t *testing.T) {
	p, src := makeParser(t, "LiteralChar")

	cases := []testCase{{
		Code:   "'\\\\'",
		Result: ast.CharNode{Value: '\\'},
	}, {
		Code:   "'\\\"'",
		Result: ast.CharNode{Value: '"'},
	}, {
		Code:   "'\\''",
		Result: ast.CharNode{Value: '\''},
	}, {
		Code:   "'\\n'",
		Result: ast.CharNode{Value: '\n'},
	}, {
		Code:   "'\\r'",
		Result: ast.CharNode{Value: '\r'},
	}, {
		Code:   "'\\t'",
		Result: ast.CharNode{Value: '\t'},
	}}
	// add the entirety of printable ASCII to test cases
	for c := byte(0x20); c < byte(0x7F); c++ {
		if c != '\\' && c != '\'' {
			cases = append(cases, testCase{
				Code:   fmt.Sprintf("'%c'", c),
				Result: ast.CharNode{Value: c},
			})
		}
	}

	expectAllValues(t, p, src, cases)
}

func TestLiteralInt(t *testing.T) {
	p, src := makeParser(t, "LiteralInt")

	cases := []testCase{{
		Code:   "-1",
		Result: ast.IntNode{Value: -1},
	}, {
		Code:   "-0",
		Result: ast.IntNode{Value: 0},
	}, {
		Code:   "0",
		Result: ast.IntNode{Value: 0},
	}, {
		Code:   "0xA",
		Result: ast.IntNode{Value: 0xA},
	}, {
		Code:   "0x123_456",
		Result: ast.IntNode{Value: 0x123_456},
	}, {
		Code:   "0b10110101",
		Result: ast.IntNode{Value: 0b10110101},
	}, {
		Code:   "07171",
		Result: ast.IntNode{Value: 07171},
	}, {
		Code:   "0o1726354",
		Result: ast.IntNode{Value: 0o1726354},
	}, {
		Code:   "1_2",
		Result: ast.IntNode{Value: 12},
	}}
	// generate a bunch of completely randomised integers to make sure no
	// weirdness happens for example due to large numbers
	for i := 0; i < 20; i++ {
		num := rand.Int63()
		if num%2 == 0 {
			num = -num
		}
		cases = append(cases, testCase{
			Code:   fmt.Sprint(num),
			Result: ast.IntNode{Value: num},
		})
	}

	expectAllValues(t, p, src, cases)
}

func TestLiteralFloat(t *testing.T) {
	p, src := makeParser(t, "LiteralFloat")

	cases := []testCase{{
		Code:   "1.0",
		Result: ast.FloatNode{Value: 1.0},
	}, {
		Code:   "-1.0",
		Result: ast.FloatNode{Value: -1.0},
	}, {
		Code:   "123_456.789",
		Result: ast.FloatNode{Value: 123_456.789},
	}, {
		Code:   "0.0",
		Result: ast.FloatNode{Value: 0.0},
	}}
	for i := 0; i < 20; i++ {
		num := rand.Float64()
		if rand.Int()%2 == 0 {
			num = -num
		}
		cases = append(cases, testCase{
			Code:   fmt.Sprint(num),
			Result: ast.FloatNode{Value: num},
		})
	}

	expectAllValues(t, p, src, cases)
}

func TestLiteralString(t *testing.T) {
	p, src := makeParser(t, "LiteralString")

	escaper := strings.NewReplacer("\\", "\\\\", "\"", "\\\"")
	cases := []testCase{{
		Code:   "\"\"",
		Result: ast.StringNode{Value: ""},
	}, {
		Code:   "\"\\\"\"",
		Result: ast.StringNode{Value: "\""},
	}, {
		Code:   "\"\\n\"",
		Result: ast.StringNode{Value: "\n"},
	}}
	// generate a bunch of random ASCII strings to ensure it's fine
	for i := 0; i < 20; i++ {
		out := ""
		for j := 0; j < 20; j++ {
			out += string(byte(randRange(0x20, 0x7F)))
		}
		in := "\"" + escaper.Replace(out) + "\""
		cases = append(cases, testCase{
			Code:   in,
			Result: ast.StringNode{Value: out},
		})
	}

	expectAllValues(t, p, src, cases)
}

func TestLiteralArray(t *testing.T) {
	p, src := makeParser(t, "LiteralArray")

	expectAllValues(t, p, src, []testCase{{
		Code:   "[bool]()",
		Result: arrOf(types.Bool),
	}, {
		Code:   "[char]()",
		Result: arrOf(types.Char),
	}, {
		Code:   "[int]()",
		Result: arrOf(types.Int),
	}, {
		Code:   "[float]()",
		Result: arrOf(types.Float),
	}, {
		Code:   "[string]()",
		Result: arrOf(types.String),
	}, {
		Code:   "[b](* * /)",
		Result: arrOf(types.Bool, true, true, false),
	}, {
		Code:   "[c]('a' 'b' 'c')",
		Result: arrOf(types.Char, 'a', 'b', 'c'),
	}, {
		Code:   "[i](1 2 3)",
		Result: arrOf(types.Int, 1, 2, 3),
	}, {
		Code:   "[f](1.2 3.4 5.6)",
		Result: arrOf(types.Float, 1.2, 3.4, 5.6),
	}, {
		Code:   "[s](\"hello\" \"world\")",
		Result: arrOf(types.String, "hello", "world"),
	}, {
		Code:   "[](* / / *)",
		Result: arrOf(types.Bool, true, false, false, true),
	}, {
		Code:   "[]('x' 'y' 'z')",
		Result: arrOf(types.Char, 'x', 'y', 'z'),
	}, {
		Code:   "[](9 8 7)",
		Result: arrOf(types.Int, 9, 8, 7),
	}, {
		Code:   "[](9.8 7.6 5.4)",
		Result: arrOf(types.Float, 9.8, 7.6, 5.4),
	}, {
		Code:   "[](\"foo\" \"bar\")",
		Result: arrOf(types.String, "foo", "bar"),
	}})
}
