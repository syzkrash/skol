package parser_test

import (
	"testing"

	"github.com/syzkrash/skol/ast"
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

	p.Scope.Vars["a"] = ast.BoolNode{Value: false}
	p.Scope.Vars["b"] = ast.BoolNode{Value: false}
	p.Scope.Vars["c"] = ast.BoolNode{Value: false}
	p.Scope.Vars["d"] = ast.BoolNode{Value: false}
	p.Scope.Vars["e"] = ast.BoolNode{Value: false}
	p.Scope.Vars["it"] = ast.IntNode{Value: 123}

	cases := map[string]ast.IfNode{
		"?*()": {
			Main: ast.Branch{
				Cond:  ast.MetaNode{Node: ast.BoolNode{Value: true}},
				Block: ast.Block{},
			},
			Other: []ast.Branch{},
			Else:  ast.Block{},
		},
		"?/()": {
			Main: ast.Branch{
				Cond:  ast.MetaNode{Node: ast.BoolNode{Value: false}},
				Block: ast.Block{},
			},
			Other: []ast.Branch{},
			Else:  ast.Block{},
		},
		"?do!(print! \"Deez\")": {
			Main: ast.Branch{
				Cond: ast.MetaNode{Node: ast.FuncCallNode{
					Func: "do",
					Args: []ast.MetaNode{},
				}},
				Block: ast.Block{
					ast.MetaNode{
						Node: ast.FuncCallNode{
							Func: "print",
							Args: []ast.MetaNode{
								{Node: ast.StringNode{Value: "Deez"}},
							},
						},
					},
				},
			},
			Other: []ast.Branch{},
			Else:  ast.Block{},
		},
		"?ok! it(print! \"OK\"):(print! \"Not OK\")": {
			Main: ast.Branch{
				Cond: ast.MetaNode{Node: ast.FuncCallNode{
					Func: "ok",
					Args: []ast.MetaNode{
						{Node: ast.SelectorNode{
							Parent: nil,
							Child:  "it",
						}},
					},
				}},
				Block: ast.Block{
					ast.MetaNode{
						Node: ast.FuncCallNode{
							Func: "print",
							Args: []ast.MetaNode{
								{Node: ast.StringNode{Value: "OK"}},
							},
						},
					},
				},
			},
			Other: []ast.Branch{},
			Else: ast.Block{
				ast.MetaNode{
					Node: ast.FuncCallNode{
						Func: "print",
						Args: []ast.MetaNode{
							{Node: ast.StringNode{Value: "Not OK"}},
						},
					},
				},
			},
		},
		"?a(>a):?b(>b):(>c)": {
			Main: ast.Branch{
				Cond: ast.MetaNode{Node: ast.SelectorNode{
					Parent: nil,
					Child:  "a",
				}},
				Block: ast.Block{
					ast.MetaNode{
						Node: ast.ReturnNode{
							Value: ast.MetaNode{Node: ast.SelectorNode{
								Parent: nil,
								Child:  "a",
							}},
						},
					},
				},
			},
			Other: []ast.Branch{
				{
					Cond: ast.MetaNode{Node: ast.SelectorNode{
						Parent: nil,
						Child:  "b",
					}},
					Block: ast.Block{
						ast.MetaNode{
							Node: ast.ReturnNode{
								Value: ast.MetaNode{Node: ast.SelectorNode{
									Parent: nil,
									Child:  "b",
								}},
							},
						},
					},
				},
			},
			Else: ast.Block{
				ast.MetaNode{
					Node: ast.ReturnNode{
						Value: ast.MetaNode{Node: ast.SelectorNode{
							Parent: nil,
							Child:  "c",
						}},
					},
				},
			},
		},
		"?a(>a):?b(>b):?c(>c):?d(>d):(>e)": {
			Main: ast.Branch{
				Cond: ast.MetaNode{Node: ast.SelectorNode{
					Parent: nil,
					Child:  "a",
				}},
				Block: ast.Block{
					ast.MetaNode{
						Node: ast.ReturnNode{
							Value: ast.MetaNode{Node: ast.SelectorNode{
								Parent: nil,
								Child:  "a",
							}},
						},
					},
				},
			},
			Other: []ast.Branch{
				{
					Cond: ast.MetaNode{Node: ast.SelectorNode{
						Parent: nil,
						Child:  "b",
					}},
					Block: ast.Block{
						ast.MetaNode{
							Node: ast.ReturnNode{
								Value: ast.MetaNode{Node: ast.SelectorNode{
									Parent: nil,
									Child:  "b",
								}},
							},
						},
					},
				},
				{
					Cond: ast.MetaNode{Node: ast.SelectorNode{
						Parent: nil,
						Child:  "c",
					}},
					Block: ast.Block{
						ast.MetaNode{
							Node: ast.ReturnNode{
								Value: ast.MetaNode{Node: ast.SelectorNode{
									Parent: nil,
									Child:  "c",
								}},
							},
						},
					},
				},
				{
					Cond: ast.MetaNode{Node: ast.SelectorNode{
						Parent: nil,
						Child:  "d",
					}},
					Block: ast.Block{
						ast.MetaNode{
							Node: ast.ReturnNode{
								Value: ast.MetaNode{Node: ast.SelectorNode{
									Parent: nil,
									Child:  "d",
								}},
							},
						},
					},
				},
			},
			Else: ast.Block{
				ast.MetaNode{
					Node: ast.ReturnNode{
						Value: ast.MetaNode{Node: ast.SelectorNode{
							Parent: nil,
							Child:  "e",
						}},
					},
				},
			},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, ast.MetaNode{Node: out})
		t.Log("OK")
	}
}

func TestWhile(t *testing.T) {
	p, src := makeParser("While")

	p.Scope.Vars["Quit"] = ast.BoolNode{Value: false}

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

	cases := map[string]ast.WhileNode{
		"**()": {
			Cond:  ast.MetaNode{Node: ast.BoolNode{Value: true}},
			Block: ast.Block{},
		},
		"*/()": {
			Cond:  ast.MetaNode{Node: ast.BoolNode{Value: false}},
			Block: ast.Block{},
		},
		"*Do!(Dont!)": {
			Cond: ast.MetaNode{Node: ast.FuncCallNode{
				Func: "Do",
				Args: []ast.MetaNode{},
			}},
			Block: ast.Block{
				ast.MetaNode{
					Node: ast.FuncCallNode{
						Func: "Dont",
						Args: []ast.MetaNode{},
					},
				},
			},
		},
		"*not! Quit (Do! >Dont!)": {
			Cond: ast.MetaNode{Node: ast.FuncCallNode{
				Func: "not",
				Args: []ast.MetaNode{
					{Node: ast.SelectorNode{
						Parent: nil,
						Child:  "Quit",
					}},
				},
			}},
			Block: ast.Block{
				ast.MetaNode{
					Node: ast.FuncCallNode{
						Func: "Do",
						Args: []ast.MetaNode{},
					},
				},
				ast.MetaNode{
					Node: ast.ReturnNode{
						Value: ast.MetaNode{Node: ast.FuncCallNode{
							Func: "Dont",
							Args: []ast.MetaNode{},
						}},
					},
				},
			},
		},
	}

	for in, out := range cases {
		t.Log(in)
		src.Reset(in)
		expect(t, p, ast.MetaNode{Node: out})
		t.Log("OK")
	}
}
