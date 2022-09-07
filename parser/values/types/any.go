package types

// AnyType is a type that is completely compatible with any other type. This
// is a very dangerous type and it is only allowed to be used by built-in
// functions for generics.
type AnyType struct{}

func (AnyType) Prim() Primitive {
	return PAny
}

func (AnyType) Equals(Type) bool {
	return true
}

func (AnyType) String() string {
	return "Any"
}

// NothingType is equivalent to a nil, null or undefined. The undefined type is
// only used by the parser as a placeholder.
type NothingType struct {
	undef bool
}

func (t NothingType) Prim() Primitive {
	if t.undef {
		return PUndefined
	}
	return PNothing
}

func (NothingType) Equals(Type) bool {
	return false
}

func (t NothingType) String() string {
	if t.undef {
		return "Undefined"
	}
	return "Nothing"
}
