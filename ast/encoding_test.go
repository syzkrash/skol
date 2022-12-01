package ast_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/parser"
)

func TestEncode(t *testing.T) {
	p := parser.NewParser(
		"TestEncode",
		strings.NewReader(`
			#greeting: "Hello, world!"

			$hello/int(
				print! greeting
				>123
			)
		`),
		"test")

	tree, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	out := bytes.Buffer{}
	if err := ast.Encode(&out, tree); err != nil {
		if p, ok := err.(common.Printable); ok {
			p.Print()
		}
		t.Fatal(err)
	}
	t.Log(out.Bytes())
}
