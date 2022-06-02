package parser

type ValueType uint8

const (
	VtNothing ValueType = iota
	VtInteger
	VtBool
	VtFloat
	VtString
	VtChar
	VtPointer
	VtAny
)

var types = [...]string{
	"Nothing",
	"Integer",
	"Boolean",
	"Float",
	"String",
	"Char",
	"Pointer",
	"Any",
}

func (t ValueType) String() string {
	return types[t]
}

func ParseType(raw string) (t ValueType, ok bool) {
	ok = true
	switch raw {
	case "i", "int", "integer":
		t = VtInteger
	case "b", "bool", "boolean":
		t = VtBool
	case "f", "float", "real":
		t = VtFloat
	case "s", "str", "string":
		t = VtString
	case "c", "char", "rune":
		t = VtChar
	case "a", "any":
		t = VtAny
	case "n", "null", "none", "nothing", "v", "void":
		t = VtNothing
	default:
		ok = false
	}
	return
}

func (p *Parser) TypeOf(n Node) (t ValueType, ok bool) {
	switch n.Kind() {
	case NdInteger:
		t = VtInteger
	case NdFloat:
		t = VtFloat
	case NdChar:
		t = VtChar
	case NdString:
		t = VtString
	case NdReturn:
		t, ok = p.TypeOf(n.(*ReturnNode).Value)
	case NdVarRef:
		var v *VarDefNode
		v, ok = p.Scope.FindVar(n.(*VarRefNode).Var)
		if !ok {
			return
		}
		t = v.VarType
	case NdFuncCall:
		var f *Function
		f, ok = p.Scope.FindFunc(n.(*FuncCallNode).Func)
		if !ok {
			return
		}
		t = f.Ret
	default:
		ok = false
	}
	return
}
