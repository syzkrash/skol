package parser

import "github.com/syzkrash/skol/parser/values/types"

type Function struct {
	Name string
	Args []types.Descriptor
	Ret  types.Type
}
