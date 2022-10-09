package parser_test

import (
	"testing"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/parser/values/types"
)

func TestIf(t *testing.T) {
	p, src := makeParser("If")

	p.Tree.Funcs["do"] = ast.Func{
		Name: "do",
		Args: []types.Descriptor{},
		Ret:  types.Bool,
	}
	p.Tree.Funcs["ok"] = ast.Func{
		Name: "ok",
		Args: []types.Descriptor{
			{
				Name: "it",
				Type: types.Int,
			},
		},
		Ret: types.Bool,
	}
	p.Tree.Funcs["print"] = ast.Func{
		Name: "print",
		Args: []types.Descriptor{
			{
				Name: "what",
				Type: types.String,
			},
		},
		Ret: types.Nothing,
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

	p.Tree.Funcs["Do"] = ast.Func{
		Name: "Do",
		Args: []types.Descriptor{},
		Ret:  types.Nothing,
	}
	p.Tree.Funcs["Dont"] = ast.Func{
		Name: "Dont",
		Args: []types.Descriptor{},
		Ret:  types.Int,
	}
	p.Tree.Funcs["not"] = ast.Func{
		Name: "not",
		Args: []types.Descriptor{
			{
				Name: "what",
				Type: types.Bool,
			},
		},
		Ret: types.Bool,
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
