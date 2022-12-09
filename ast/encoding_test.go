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
	errs := make(chan error)
	var parseError error

	go func() {
		for err := range errs {
			if err == nil {
				continue
			}

			close(errs)
			parseError = err
		}
	}()

	p := parser.NewParser(
		"TestEncode",
		strings.NewReader(`
			#greeting: "Hello, world!"

			$hello/int(
				print! greeting
				>123
			)
		`),
		"test", errs)

	tree := p.Parse()
	if parseError != nil {
		t.Fatal(parseError)
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

func TestRecode(t *testing.T) {
	errs := make(chan error)
	var parseError error

	go func() {
		for err := range errs {
			if err == nil {
				continue
			}

			close(errs)
			parseError = err
		}
	}()

	p := parser.NewParser(
		"TestRecode",
		strings.NewReader(`
			#greeting: "Hello, world!"

			$hello/int(
				print! greeting
				>123
			)
		`),
		"test", errs)

	tree := p.Parse()
	if parseError != nil {
		t.Fatal(parseError)
	}

	out := bytes.Buffer{}
	if err := ast.Encode(&out, tree); err != nil {
		if p, ok := err.(common.Printable); ok {
			p.Print()
		}
		t.Fatal(err)
	}

	decTree, err := ast.Decode(bytes.NewReader(out.Bytes()))
	if err != nil {
		if p, ok := err.(common.Printable); ok {
			p.Print()
		}
		t.Fatal(err)
	}

	if len(decTree.Vars) != len(tree.Vars) {
		t.Fatalf("incorrect # of vars: %d != %d", len(decTree.Vars), len(tree.Vars))
	}
	if len(decTree.Typedefs) != len(tree.Typedefs) {
		t.Fatalf("incorrect # of typedefs: %d != %d", len(decTree.Typedefs), len(tree.Typedefs))
	}
	if len(decTree.Funcs) != len(tree.Funcs) {
		t.Fatalf("incorrect # of funcs: %d != %d", len(decTree.Funcs), len(tree.Funcs))
	}
	if len(decTree.Exerns) != len(tree.Exerns) {
		t.Fatalf("incorrect # of externs: %d != %d", len(decTree.Exerns), len(tree.Exerns))
	}
	if len(decTree.Structs) != len(tree.Structs) {
		t.Fatalf("incorrect # of structs: %d != %d", len(decTree.Structs), len(tree.Structs))
	}
}
