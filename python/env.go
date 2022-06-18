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
			Args: []values.FuncArg{{"a", values.Any}},
			Ret:  values.Nothing,
		},
		"import": {
			Name: "import",
			Args: []values.FuncArg{{"module", values.String}},
			Ret:  values.Nothing,
		},
	}
	p.parser.Scope = newScope
}
