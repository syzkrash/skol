package ast_test

import (
	"bytes"
	"testing"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestEncode(t *testing.T) {
	tree := ast.NewAST()
	tree.Vars["greeting"] = ast.Var{
		Name: "greeting",
		Value: ast.MetaNode{Node: ast.StringNode{
			Value: "Hello, world!",
		}},
	}
	tree.Funcs["hello"] = ast.Func{
		Name: "hello",
		Args: []types.Descriptor{},
		Ret:  types.Nothing,
		Body: []ast.MetaNode{{Node: ast.FuncCallNode{
			Func: "print",
			Args: []ast.MetaNode{{Node: ast.SelectorNode{
				Parent: nil,
				Child: "greeting",
			}}},
		}}},
		Node: ast.MetaNode{},
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
