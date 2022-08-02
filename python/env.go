package python

import (
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

func (p *pythonState) initEnv() {
	newScope := parser.NewScope(p.parser.Scope)
	newScope.Funcs = map[string]*values.Function{
		"print": {
			Name: "print",
			Args: []values.FuncArg{{"a", types.Any}},
			Ret:  types.Nothing,
		},
		"import": {
			Name: "import",
			Args: []values.FuncArg{{"module", types.String}},
			Ret:  types.Nothing,
		},
	}
	p.parser.Scope = newScope
}
