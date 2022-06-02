package parser

type Function struct {
	Name string
	Args map[string]ValueType
	Ret  ValueType
}

func DefinedFunction(n *FuncDefNode) *Function {
	return &Function{
		Name: n.Name,
		Args: n.Args,
		Ret:  n.Ret,
	}
}

func ExternFunction(n *FuncExternNode) *Function {
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
