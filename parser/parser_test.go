package parser

import (
	"strings"
	"testing"
)

func TestVarDef(t *testing.T) {
	code := ` %a: 'E'  `
	src := strings.NewReader(code)
	p := NewParser("TestVarDef", src)
	n, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != NdVarDef {
		t.Fatalf("expected %s node, got %s", NdVarDef, n.Kind())
	}
	v := n.(*VarDefNode)
	if v.Var != "a" {
		t.Fatalf("expected variable name %s, got %s", "a", v.Var)
	}
	if v.Value.Kind() != NdChar {
		t.Fatalf("expected value to be %s node, got %s", NdChar, v.Value.Kind())
	}
	c := v.Value.(*CharNode)
	if c.Char != 'E' {
		t.Fatalf("expected vale to be rune %c, got %c", 'E', c.Char)
	}
}

func TestVarRef(t *testing.T) {
	code := ` a `
	src := strings.NewReader(code)
	p := NewParser("TestVarRef", src)
	p.scope.Vars["a"] = &VarDefNode{} // to prevent 'unknown variable' error
	n, err := p.value()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != NdVarRef {
		t.Fatalf("expected %s node, got %s", NdVarRef, n.Kind())
	}
	v := n.(*VarRefNode)
	if v.Var != "a" {
		t.Fatalf("expected variable name %s, got %s", "a", v.Var)
	}
}

func TestFuncCall(t *testing.T) {
	code := `	add!	1.2 3.4`
	src := strings.NewReader(code)
	p := NewParser("TestFuncCall", src)
	// to prevent 'unknown function' error
	p.scope.Funcs["add"] = &FuncDefNode{
		Func: "add",
		Args: map[string]ValueType{
			"a": VtFloat,
			"b": VtFloat,
		},
		Ret:  VtFloat,
		Body: []Node{},
	}
	n, err := p.value()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != NdFuncCall {
		t.Fatalf("expected %s node, got %s", NdFuncCall, n.Kind())
	}
	c := n.(*FuncCallNode)
	if c.Func != "add" {
		t.Fatalf("expected function name %s, got %s", "add", c.Func)
	}
	if len(c.Args) != 2 {
		t.Fatalf("expected %d arguments, got %d", 2, len(c.Args))
	}
	a := c.Args[0]
	if a.Kind() != NdFloat {
		t.Fatalf("expected argument %d to be %s, got %s", 0, NdFloat, a.Kind())
	}
	a = c.Args[1]
	if a.Kind() != NdFloat {
		t.Fatalf("expected argument %d to be %s, got %s", 1, NdFloat, a.Kind())
	}
}
