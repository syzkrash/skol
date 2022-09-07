package types

// Primitive represents the unique identifying number of a type.
type Primitive uint8

// Primitive constants.
const (
	PBool Primitive = iota
	PChar
	PInt
	PFloat
	PString
	PStruct
	PArray
	PAny
	PNothing
	PUndefined
)

// Type represents a Skol type.
type Type interface {
	// Prim returns the unique identifying primitive of this type.
	Prim() Primitive
	// Equals checks whether the other Type is compatible with this type. In
	// simpler terms, checks whether a value of the provided type can be cast
	// to this type.
	Equals(Type) bool
	// String returns a user-friendy name for this type.
	String() string
}

// PrimType represents a simple type that is only identified via its primitive.
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
