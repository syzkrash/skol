package python

import (
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/values"
)

func (p *pythonState) initEnv() {
	newScope := parser.NewScope(p.parser.Scope)
	newScope.Funcs = map[string]*parser.Function{
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
	}
	p.parser.Scope = newScope
}
