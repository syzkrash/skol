package values

import "github.com/syzkrash/skol/parser/values/types"

type FuncArg struct {
	Name string
	Type types.Type
}

// Function contains function prototype information.
type Function struct {
	Name string
	Args []FuncArg
	Ret  types.Type
}

// DefinedFunction creates a new [Function] from the given arguments.
// NOTE: this function does not need to exist anymore :)
func DefinedFunction(name string, args []FuncArg, ret types.Type) *Function {
	return &Function{
		Name: name,
		Args: args,
		Ret:  ret,
	}
}

// ExternFunction creates a new [Function] based on an extern definition.
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
