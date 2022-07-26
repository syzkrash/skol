package types

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
