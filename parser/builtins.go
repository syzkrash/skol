package parser

// builtin represents a function built into the skol language itself.
// the skol parser only needs an argument count of a function to properly parse
// it. hence this struct only contains such information. parsed functions retain
// all other information as it is later passed on to things like the typechecker
// and optimizers, which will need the extra information. these passes will have
// their unique sets of information on builtin functions that they need
type builtin struct {
	ArgCount int
}

var builtins = map[string]builtin{
	"add": {ArgCount: 2},
	"sub": {ArgCount: 2},
	"mul": {ArgCount: 2},
	"div": {ArgCount: 2},
	"pow": {ArgCount: 2},
	"mod": {ArgCount: 2},

	"eq": {ArgCount: 2},
	"gt": {ArgCount: 2},
	"lt": {ArgCount: 2},

	"not": {ArgCount: 1},
	"and": {ArgCount: 2},
	"or":  {ArgCount: 2},

	"append": {ArgCount: 2},
	"concat": {ArgCount: 2},
	"slice":  {ArgCount: 3},
	"at":     {ArgCount: 2},
	"len":    {ArgCount: 1},

	"str":  {ArgCount: 1},
	"bool": {ArgCount: 1},

	"parse_bool": {ArgCount: 1},
	"char":       {ArgCount: 1},
	"int":        {ArgCount: 1},
	"float":      {ArgCount: 1},

	"print": {ArgCount: 1},
}
