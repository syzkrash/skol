package sim

import "github.com/syzkrash/skol/parser"

type Funct struct {
	Args map[string]parser.ValueType
	Ret  parser.ValueType
	Body []parser.Node
}
