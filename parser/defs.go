package parser

import (
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type Function struct {
	Name string
	Args map[string]values.ValueType
	Ret  values.ValueType
}

func DefinedFunction(n *nodes.FuncDefNode) *Function {
	return &Function{
		Name: n.Name,
		Args: n.Args,
		Ret:  n.Ret,
	}
}

func ExternFunction(n *nodes.FuncExternNode) *Function {
	var name string
	if n.Intern != "" {
		name = n.Intern
	} else {
		name = n.Name
	}
	return &Function{
		Name: name,
		Args: n.Args,
		Ret:  n.Ret,
	}
}

var DefaultFuncs = map[string]*Function{
	"print": {
		Name: "print",
		Args: map[string]values.ValueType{"a": values.VtString},
		Ret:  values.VtNothing,
	},
	"to_str": {
		Name: "to_str",
		Args: map[string]values.ValueType{"a": values.VtAny},
		Ret:  values.VtString,
	},
	"add_i": {
		Name: "add_i",
		Args: map[string]values.ValueType{"a": values.VtInteger, "b": values.VtInteger},
		Ret:  values.VtInteger,
	},
}
