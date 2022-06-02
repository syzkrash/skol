package python

import "github.com/syzkrash/skol/parser"

func (p *pythonState) initEnv() {
	newScope := &parser.Scope{
		Parent: nil,
		Funcs: map[string]*parser.Function{
			"print": {
				Name: "print",
				Args: map[string]parser.ValueType{
					"a": parser.VtAny,
				},
				Ret: parser.VtNothing,
			},
		},
		Vars: map[string]*parser.VarDefNode{},
	}

	for oper, sym := range ops {
		newScope.Funcs[oper] = &parser.Function{
			Name: sym,
			Args: map[string]parser.ValueType{
				"a": parser.VtAny,
				"b": parser.VtAny,
			},
			Ret: parser.VtAny,
		}
	}

	p.parser.Scope.Parent = newScope
}
