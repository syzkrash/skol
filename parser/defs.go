package parser

type Function struct {
	Name string
	Args map[string]ValueType
	Ret  ValueType
}

func DefinedFunction(n *FuncDefNode) *Function {
	return &Function{
		Name: n.Func,
		Args: n.Arg,
		Ret:  n.RetType,
	}
}
