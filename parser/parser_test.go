package parser

import (
	"strings"
	"testing"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

func TestVarDef(t *testing.T) {
	code := ` %a: 'E'  `
	src := strings.NewReader(code)
	p := NewParser("TestVarDef", src, "test")
	n, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != nodes.NdVarDef {
		t.Fatalf("expected %s node, got %s", nodes.NdVarDef, n.Kind())
	}
	v := n.(*nodes.VarDefNode)
	if v.Var != "a" {
		t.Fatalf("expected variable name %s, got %s", "a", v.Var)
	}
	if v.Value.Kind() != nodes.NdChar {
		t.Fatalf("expected value to be %s node, got %s", nodes.NdChar, v.Value.Kind())
	}
	c := v.Value.(*nodes.CharNode)
	if c.Char != 'E' {
		t.Fatalf("expected vale to be rune %c, got %c", 'E', c.Char)
	}
}

func TestVarRef(t *testing.T) {
	code := ` a `
	src := strings.NewReader(code)
	p := NewParser("TestVarRef", src, "test")
	p.Scope.SetVar("a", &nodes.VarDefNode{}) // to prevent 'unknown variable' error)
	n, err := p.value()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != nodes.NdVarRef {
		t.Fatalf("expected %s node, got %s", nodes.NdVarRef, n.Kind())
	}
	v := n.(*nodes.VarRefNode)
	if v.Var != "a" {
		t.Fatalf("expected variable name %s, got %s", "a", v.Var)
	}
}

func TestFuncCall(t *testing.T) {
	code := `	add!	1.2 3.4`
	src := strings.NewReader(code)
	p := NewParser("TestFuncCall", src, "test")
	// to prevent 'unknown function' error
	p.Scope.SetFunc("add", &Function{
		Name: "add",
		Args: []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:  values.VtFloat,
	})
	n, err := p.value()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected %s node, got %s", nodes.NdFuncCall, n.Kind())
	}
	c := n.(*nodes.FuncCallNode)
	if c.Func != "add" {
		t.Fatalf("expected function name %s, got %s", "add", c.Func)
	}
	if len(c.Args) != 2 {
		t.Fatalf("expected %d arguments, got %d", 2, len(c.Args))
	}
	a := c.Args[0]
	if a.Kind() != nodes.NdFloat {
		t.Fatalf("expected argument %d to be %s, got %s", 0, nodes.NdFloat, a.Kind())
	}
	a = c.Args[1]
	if a.Kind() != nodes.NdFloat {
		t.Fatalf("expected argument %d to be %s, got %s", 1, nodes.NdFloat, a.Kind())
	}
}

func TestIf(t *testing.T) {
	code := `	?1(print!"hello world") `
	src := strings.NewReader(code)
	p := NewParser("TestIf", src, "test")
	// to prevent 'unknown function' error
	p.Scope.SetFunc("print", &Function{
		Name: "print",
		Args: []values.FuncArg{{"a", values.VtAny}},
		Ret:  values.VtNothing,
	})
	n, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != nodes.NdIf {
		t.Fatalf("expected If, got %s", n.Kind())
	}
	ifn := n.(*nodes.IfNode)
	if ifn.Condition.Kind() != nodes.NdInteger {
		t.Fatalf("expected Integer, got %s", n.Kind())
	}
	if len(ifn.IfBlock) < 1 {
		t.Fatalf("expected 1 Node in block, got %d", len(ifn.IfBlock))
	}
	pn := ifn.IfBlock[0]
	if pn.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", pn.Kind())
	}
	fcn := pn.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}
}

func TestIfBetween(t *testing.T) {
	code := `	?1(print!"hello world")print!"bye world" `
	src := strings.NewReader(code)
	p := NewParser("TestIfBetween", src, "test")
	// to prevent 'unknown function' error
	p.Scope.SetFunc("print", &Function{
		Name: "print",
		Args: []values.FuncArg{{"a", values.VtAny}},
		Ret:  values.VtNothing,
	})

	// check for the if statement
	n, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != nodes.NdIf {
		t.Fatalf("expected If, got %s", n.Kind())
	}
	ifn := n.(*nodes.IfNode)
	if ifn.Condition.Kind() != nodes.NdInteger {
		t.Fatalf("expected Integer, got %s", n.Kind())
	}
	if len(ifn.IfBlock) < 1 {
		t.Fatalf("expected 1 Node in block, got %d", len(ifn.IfBlock))
	}
	pn := ifn.IfBlock[0]
	if pn.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", pn.Kind())
	}
	fcn := pn.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}

	// check for the print statement after it
	n, err = p.Next()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", n.Kind())
	}
	fcn = n.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}
}

func TestIfElse(t *testing.T) {
	code := `	?1(print!"hello world"):(print!"bye world") `
	src := strings.NewReader(code)
	p := NewParser("TestIfElse", src, "test")
	// to prevent 'unknown function' error
	p.Scope.SetFunc("print", &Function{
		Name: "print",
		Args: []values.FuncArg{{"a", values.VtAny}},
		Ret:  values.VtNothing,
	})

	n, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}

	// check the if itself
	if n.Kind() != nodes.NdIf {
		t.Fatalf("expected If, got %s", n.Kind())
	}
	ifn := n.(*nodes.IfNode)
	if ifn.Condition.Kind() != nodes.NdInteger {
		t.Fatalf("expected Integer, got %s", n.Kind())
	}

	// check the IfBlock
	if len(ifn.IfBlock) != 1 {
		t.Fatalf("expected 1 Node in IfBlock, got %d", len(ifn.IfBlock))
	}
	pn := ifn.IfBlock[0]
	if pn.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", pn.Kind())
	}
	fcn := pn.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}

	// check the ElseBlock
	if len(ifn.ElseBlock) != 1 {
		t.Fatalf("expected 1 Node in ElseBlock, got %d", len(ifn.ElseBlock))
	}
	pn = ifn.ElseBlock[0]
	if pn.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", pn.Kind())
	}
	fcn = pn.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}
}

func TestIfElseIfElse(t *testing.T) {
	code := `	?1(print!"hello world"):?0(print!"world?"):(print!"bye world") `
	src := strings.NewReader(code)
	p := NewParser("TestIfElseIfElse", src, "test")
	// to prevent 'unknown function' error
	p.Scope.SetFunc("print", &Function{
		Name: "print",
		Args: []values.FuncArg{{"a", values.VtAny}},
		Ret:  values.VtNothing,
	})

	n, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}

	// check the if itself
	if n.Kind() != nodes.NdIf {
		t.Fatalf("expected If, got %s", n.Kind())
	}
	ifn := n.(*nodes.IfNode)
	if ifn.Condition.Kind() != nodes.NdInteger {
		t.Fatalf("expected Integer, got %s", n.Kind())
	}

	// check the IfBlock
	if len(ifn.IfBlock) != 1 {
		t.Fatalf("expected 1 Node in IfBlock, got %d", len(ifn.IfBlock))
	}
	pn := ifn.IfBlock[0]
	if pn.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", pn.Kind())
	}
	fcn := pn.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}

	// check the nodes.IfSubNode
	if len(ifn.ElseIfNodes) != 1 {
		t.Fatalf("expected 1 nodes.IfSubNode in IfBlock, got %d", len(ifn.ElseIfNodes))
	}
	sn := ifn.ElseIfNodes[0]
	if len(sn.Block) != 1 {
		t.Fatalf("expected 1 Node in nodes.IfSubNode body, got %d", len(sn.Block))
	}
	pn = sn.Block[0]
	if pn.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", pn.Kind())
	}
	fcn = pn.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}

	// check the ElseBlock
	if len(ifn.ElseBlock) != 1 {
		t.Fatalf("expected 1 Node in ElseBlock, got %d", len(ifn.ElseBlock))
	}
	pn = ifn.ElseBlock[0]
	if pn.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", pn.Kind())
	}
	fcn = pn.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}
}

func TestIfElseIf(t *testing.T) {
	code := `	?1(print!"hello world"):?1(print!"bye world?") `
	src := strings.NewReader(code)
	p := NewParser("TestIfElseIf", src, "test")
	// to prevent 'unknown function' error
	p.Scope.SetFunc("print", &Function{
		Name: "print",
		Args: []values.FuncArg{{"a", values.VtAny}},
		Ret:  values.VtNothing,
	})

	n, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}

	// check the if itself
	if n.Kind() != nodes.NdIf {
		t.Fatalf("expected If, got %s", n.Kind())
	}
	ifn := n.(*nodes.IfNode)
	if ifn.Condition.Kind() != nodes.NdInteger {
		t.Fatalf("expected Integer, got %s", n.Kind())
	}

	// check the IfBlock
	if len(ifn.IfBlock) != 1 {
		t.Fatalf("expected 1 Node in IfBlock, got %d", len(ifn.IfBlock))
	}
	pn := ifn.IfBlock[0]
	if pn.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", pn.Kind())
	}
	fcn := pn.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}

	// check the nodes.IfSubNode
	if len(ifn.ElseIfNodes) != 1 {
		t.Fatalf("expected 1 nodes.IfSubNode in IfBlock, got %d", len(ifn.ElseIfNodes))
	}
	sn := ifn.ElseIfNodes[0]
	if len(sn.Block) != 1 {
		t.Fatalf("expected 1 Node in nodes.IfSubNode body, got %d", len(sn.Block))
	}
	pn = sn.Block[0]
	if pn.Kind() != nodes.NdFuncCall {
		t.Fatalf("expected FuncCall, got %s", pn.Kind())
	}
	fcn = pn.(*nodes.FuncCallNode)
	if fcn.Func != "print" {
		t.Fatalf("expected call to print, got call to %s", fcn.Func)
	}
	if len(fcn.Args) != 1 {
		t.Fatalf("expected 1 argument in call, got %d", len(fcn.Args))
	}
}

func TestConst(t *testing.T) {
	code := ` #max_int: 169 %max_int_copy: max_int  `
	src := strings.NewReader(code)
	p := NewParser("TestConst", src, "test")
	n, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}
	if n.Kind() != nodes.NdVarDef {
		t.Fatalf("expected %s node, got %s", nodes.NdVarDef, n.Kind())
	}
	v := n.(*nodes.VarDefNode)
	if v.Var != "max_int_copy" {
		t.Fatalf("expected variable name %s, got %s", "max_int_copy", v.Var)
	}
	if v.Value.Kind() != nodes.NdInteger {
		t.Fatalf("expected value to be %s node, got %s", nodes.NdInteger, v.Value.Kind())
	}
	c := v.Value.(*nodes.IntegerNode)
	if c.Int != 169 {
		t.Fatalf("expected vale to be int %d, got %c", 169, c.Int)
	}
}
