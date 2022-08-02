package types

type ArrayType struct {
	Element Type
}

func (ArrayType) Prim() Primitive {
	return PArray
}

func (a ArrayType) Equals(b Type) bool {
	if b.Prim() != PArray {
		return false
	}
	return a.Element.Equals(b.(ArrayType).Element)
}

func (t ArrayType) String() string {
	return "Array of " + t.Element.String()
}
