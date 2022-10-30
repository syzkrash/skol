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

	expectAll(t, p, src, []testCase{{
		Code: "?*()",
		Result: ast.IfNode{
			Main: ast.Branch{
				Cond:  ast.MetaNode{Node: ast.BoolNode{Value: true}},
				Block: ast.Block{},
			},
			Other: []ast.Branch{},
			Else:  ast.Block{},
		}}, {
		Code: "?/()",
		Result: ast.IfNode{
			Main: ast.Branch{
				Cond:  ast.MetaNode{Node: ast.BoolNode{Value: false}},
				Block: ast.Block{},
			},
			Other: []ast.Branch{},
			Else:  ast.Block{},
		}}, {
		Code: "?do!(print! \"Deez\")",
		Result: ast.IfNode{
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
		}}, {
		Code: "?ok! it(print! \"OK\"):(print! \"Not OK\")",
		Result: ast.IfNode{
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
		}}, {
		Code: "?a(>a):?b(>b):(>c)",
		Result: ast.IfNode{
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
		}}, {
		Code: "?a(>a):?b(>b):?c(>c):?d(>d):(>e)",
		Result: ast.IfNode{
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
		}},
	})
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

	expectAll(t, p, src, []testCase{{
		Code: "**()",
		Result: ast.WhileNode{
			Cond:  ast.MetaNode{Node: ast.BoolNode{Value: true}},
			Block: ast.Block{},
		}}, {
		Code: "*/()",
		Result: ast.WhileNode{
			Cond:  ast.MetaNode{Node: ast.BoolNode{Value: false}},
			Block: ast.Block{},
		}}, {
		Code: "*Do!(Dont!)",
		Result: ast.WhileNode{
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
		}}, {
		Code: "*not! Quit (Do! >Dont!)",
		Result: ast.WhileNode{
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
		}},
	})
}
