package python

import (
	"github.com/syzkrash/skol/parser/nodes"
)

var operators = map[string]string{
	"add_i":  "+",
	"add_f":  "+",
	"sub_i":  "-",
	"sub_f":  "-",
	"mul_i":  "*",
	"mul_f":  "*",
	"div_i":  "//",
	"div_f":  "/",
	"mod_i":  "%",
	"mod_f":  "%",
	"concat": "+",
	"or":     "or",
	"and":    "and",
	"eq":     "==",
	"gt_i":   ">",
	"gt_f":   ">",
	"lt_i":   "<",
	"lt_f":   "<",
}

var renames = map[string]string{
	"to_str":  "str",
	"to_bool": "bool",
}

type specialGenerator func(p *pythonState, n *nodes.FuncCallNode) error

var specialGenerators = map[string]specialGenerator{
	"char_at": func(p *pythonState, n *nodes.FuncCallNode) (err error) {
		_, err = p.out.WriteString("bytes(")
		if err != nil {
			return err
		}
		err = p.stringOrVar(n.Args[0])
		if err != nil {
			return err
		}
		_, err = p.out.WriteString(",'utf8')[")
		if err != nil {
			return err
		}
		err = p.integerOrVar(n.Args[1])
		if err != nil {
			return err
		}
		_, err = p.out.WriteString("]")
		return err
	},
	"substr": func(p *pythonState, n *nodes.FuncCallNode) (err error) {
		err = p.stringOrVar(n.Args[0])
		if err != nil {
			return err
		}
		_, err = p.out.WriteString("[")
		if err != nil {
			return err
		}
		err = p.integerOrVar(n.Args[1])
		if err != nil {
			return err
		}
		_, err = p.out.WriteString(":")
		if err != nil {
			return err
		}
		err = p.integerOrVar(n.Args[2])
		if err != nil {
			return err
		}
		_, err = p.out.WriteString("]")
		return err
	},
	"char_append": func(p *pythonState, n *nodes.FuncCallNode) (err error) {
		err = p.stringOrVar(n.Args[0])
		if err != nil {
			return err
		}
		_, err = p.out.WriteString("+bytes([")
		if err != nil {
			return err
		}
		err = p.charOrVar(n.Args[1])
		if err != nil {
			return err
		}
		_, err = p.out.WriteString("]).decode()")
		return err
	},
}
