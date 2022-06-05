package python

import (
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

func (p *pythonState) initEnv() {
	newScope := &parser.Scope{
		Parent: nil,
		Funcs: map[string]*parser.Function{
			"print": {
				Name: "print",
				Args: map[string]values.ValueType{
					"a": values.VtAny,
				},
				Ret: values.VtNothing,
			},
			"import": {
				Name: "import",
				Args: map[string]values.ValueType{
					"module": values.VtString,
				},
				Ret: values.VtNothing,
			},
		},
		Vars: map[string]*nodes.VarDefNode{},
	}

	for oper, sym := range ops {
		newScope.Funcs[oper] = &parser.Function{
			Name: sym,
			Args: map[string]values.ValueType{
				"a": values.VtAny,
				"b": values.VtAny,
			},
			Ret: values.VtAny,
		}
	}

	p.parser.Scope.Parent = newScope
}
