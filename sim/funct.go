package sim

import (
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type Funct struct {
	Args map[string]values.ValueType
	Ret  values.ValueType
	Body []nodes.Node
}
