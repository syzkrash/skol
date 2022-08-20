package parser_test

import (
	"testing"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestIf(t *testing.T) {
	p, src := makeParser("If")

	p.Scope.Funcs["do"] = &values.Function{
		Name: "do",
		Args: []values.FuncArg{},
		Ret:  types.Bool,
	}
	p.Scope.Funcs["ok"] = &values.Function{
		Name: "ok",
		Args: []values.FuncArg{
			{
				Name: "it",
				Type: types.Int,
			},
		},
		Ret: types.Bool,
	}

	p.Scope.Vars["a"] = &nodes.VarDefNode{
		VarType: types.Bool,
		Var:     "a",
		Value:   &nodes.BooleanNode{Bool: false},
	}
	p.Scope.Vars["b"] = &nodes.VarDefNode{
		VarType: types.Bool,
		Var:     "b",
		Value:   &nodes.BooleanNode{Bool: false},
	}
	p.Scope.Vars["c"] = &nodes.VarDefNode{
		VarType: types.Bool,
		Var:     "c",
		Value:   &nodes.BooleanNode{Bool: false},
	}
	p.Scope.Vars["d"] = &nodes.VarDefNode{
		VarType: types.Bool,
		Var:     "d",
		Value:   &nodes.BooleanNode{Bool: false},
	}
	p.Scope.Vars["e"] = &nodes.VarDefNode{
		VarType: types.Bool,
		Var:     "e",
		Value:   &nodes.BooleanNode{Bool: false},
	}
	p.Scope.Vars["it"] = &nodes.VarDefNode{
		VarType: types.Int,
		Var:     "it",
		Value:   &nodes.IntegerNode{Int: 123},
	}

	cases := map[string]*nodes.IfNode{
		"?*()": {
			Condition:   &nodes.BooleanNode{Bool: true},
			IfBlock:     []nodes.Node{},
			ElseIfNodes: []*nodes.IfSubNode{},
			ElseBlock:   []nodes.Node{},
		},
		"?/()": {
			Condition:   &nodes.BooleanNode{Bool: false},
			IfBlock:     []nodes.Node{},
			ElseIfNodes: []*nodes.IfSubNode{},
			ElseBlock:   []nodes.Node{},
		},
		"?do!(print! \"Deez\")": {
			Condition: &nodes.FuncCallNode{
				Func: "do",
				Args: []nodes.Node{},
			},
			IfBlock: []nodes.Node{
				&nodes.FuncCallNode{
					Func: "print",
					Args: []nodes.Node{
						&nodes.StringNode{Str: "Deez"},
					},
				},
			},
		},
		"?ok! it(print! \"OK\"):(print! \"Not OK\")": {
			Condition: &nodes.FuncCallNode{
				Func: "ok",
				Args: []nodes.Node{
					&nodes.SelectorNode{
						Parent: nil,
						Child:  "it",
					},
				},
			},
			IfBlock: []nodes.Node{
				&nodes.FuncCallNode{
					Func: "print",
					Args: []nodes.Node{
						&nodes.StringNode{Str: "OK"},
					},
				},
			},
			ElseIfNodes: []*nodes.IfSubNode{},
			ElseBlock: []nodes.Node{
				&nodes.FuncCallNode{
					Func: "print",
					Args: []nodes.Node{
						&nodes.StringNode{Str: "Not OK"},
					},
				},
			},
		},
		"?a(>a):?b(>b):(>c)": {
			Condition: &nodes.SelectorNode{
				Parent: nil,
				Child:  "a",
			},
			IfBlock: []nodes.Node{
				&nodes.ReturnNode{
					Value: &nodes.SelectorNode{
						Parent: nil,
						Child:  "a",
					},
				},
			},
			ElseIfNodes: []*nodes.IfSubNode{
				{
					Condition: &nodes.SelectorNode{
						Parent: nil,
						Child:  "b",
					},
					Block: []nodes.Node{
						&nodes.ReturnNode{
							Value: &nodes.SelectorNode{
								Parent: nil,
								Child:  "b",
							},
						},
					},
				},
			},
			ElseBlock: []nodes.Node{
				&nodes.ReturnNode{
					Value: &nodes.SelectorNode{
						Parent: nil,
						Child:  "c",
					},
				},
			},
		},
		"?a(>a):?b(>b):?c(>c):?d(>d):(>e)": {
			Condition: &nodes.SelectorNode{
				Parent: nil,
				Child:  "a",
			},
			IfBlock: []nodes.Node{
				&nodes.ReturnNode{
					Value: &nodes.SelectorNode{
						Parent: nil,
						Child:  "a",
					},
				},
			},
			ElseIfNodes: []*nodes.IfSubNode{
				{
					Condition: &nodes.SelectorNode{
						Parent: nil,
						Child:  "b",
					},
					Block: []nodes.Node{
						&nodes.ReturnNode{
							Value: &nodes.SelectorNode{
								Parent: nil,
								Child:  "b",
							},
						},
					},
				},
				{
					Condition: &nodes.SelectorNode{
						Parent: nil,
						Child:  "c",
					},
					Block: []nodes.Node{
						&nodes.ReturnNode{
							Value: &nodes.SelectorNode{
								Parent: nil,
								Child:  "c",
							},
						},
					},
				},
				{
					Condition: &nodes.SelectorNode{
						Parent: nil,
						Child:  "d",
					},
					Block: []nodes.Node{
						&nodes.ReturnNode{
							Value: &nodes.SelectorNode{
								Parent: nil,
								Child:  "d",
							},
						},
					},
				},
			},
			ElseBlock: []nodes.Node{
				&nodes.ReturnNode{
					Value: &nodes.SelectorNode{
						Parent: nil,
						Child:  "e",
					},
				},
			},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, out)
		t.Log("OK")
	}
}

func TestWhile(t *testing.T) {
	p, src := makeParser("While")

	p.Scope.Vars["Quit"] = &nodes.VarDefNode{
		VarType: types.Bool,
		Var:     "Quit",
		Value:   &nodes.BooleanNode{Bool: false},
	}

	p.Scope.Funcs["Do"] = &values.Function{
		Name: "Do",
		Args: []values.FuncArg{},
		Ret:  types.Nothing,
	}
	p.Scope.Funcs["Dont"] = &values.Function{
		Name: "Dont",
		Args: []values.FuncArg{},
		Ret:  types.Int,
	}

	cases := map[string]*nodes.WhileNode{
		"**()": {
			Condition: &nodes.BooleanNode{Bool: true},
			Body:      []nodes.Node{},
		},
		"*/()": {
			Condition: &nodes.BooleanNode{Bool: false},
			Body:      []nodes.Node{},
		},
		"*Do!(Dont!)": {
			Condition: &nodes.FuncCallNode{
				Func: "Do",
				Args: []nodes.Node{},
			},
			Body: []nodes.Node{
				&nodes.FuncCallNode{
					Func: "Dont",
					Args: []nodes.Node{},
				},
			},
		},
		"*not! Quit (Do! >Dont!)": {
			Condition: &nodes.FuncCallNode{
				Func: "not",
				Args: []nodes.Node{
					&nodes.SelectorNode{
						Parent: nil,
						Child:  "Quit",
					},
				},
			},
			Body: []nodes.Node{
				&nodes.FuncCallNode{
					Func: "Do",
					Args: []nodes.Node{},
				},
				&nodes.ReturnNode{
					Value: &nodes.FuncCallNode{
						Func: "Dont",
						Args: []nodes.Node{},
					},
				},
			},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, out)
		t.Log("OK")
	}
}
