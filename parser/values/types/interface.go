package types

type Primitive uint8

const (
	PBool Primitive = iota
	PChar
	PInt
	PFloat
	PString
	PStruct
	PAny
	PNothing
	PUndefined
)

type Type interface {
	Prim() Primitive
	Equals(Type) bool
	String() string
}

type PrimType struct {
	prim Primitive
}

func (p PrimType) Prim() Primitive {
	return p.prim
}

func (p PrimType) Equals(t Type) bool {
	return t.Prim() == p.prim
}

func (p PrimType) String() string {
	switch p.prim {
	case PBool:
		return "Bool"
	case PChar:
		return "Char"
	case PInt:
		return "Int"
	case PFloat:
		return "Float"
	case PString:
		return "String"
	default:
		return "???"
	}
}
