package values

import "github.com/syzkrash/skol/parser/values/types"

type Function struct {
	Name string
	Args []FuncArg
	Ret  types.Type
}

func DefinedFunction(name string, args []FuncArg, ret types.Type) *Function {
	return &Function{
		Name: name,
		Args: args,
		Ret:  ret,
	}
}

func ExternFunction(name, intern string, args []FuncArg, ret types.Type) *Function {
	var fn string
	if intern != "" {
		fn = intern
	} else {
		fn = name
	}
	return &Function{
		Name: fn,
		Args: args,
		Ret:  ret,
	}
}

type FuncArg struct {
	Name string
	Type types.Type
}
