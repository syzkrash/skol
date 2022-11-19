package ir

// Type represents the unique type identifier of a value
type Type byte

// Type constants
const (
	TypeInteger Type = iota
	TypeFloat
	TypeCall
	TypeStruct
	TypeArray
	TypeRef
)

// Value represents any valid value within the IR. Note that the type of this
// value does not correspond to a skol type. It only consists of information
// critical to memory management.
type Value interface {
	Type() Type
}

// IntegerValue holds the data of a (immediate) integer value
type IntegerValue struct {
	Value int64
}

// Type returns TypeInteger
func (IntegerValue) Type() Type {
	return TypeInteger
}

var _ Value = IntegerValue{}

// FloatValue holds the data of a (immediate) float value
type FloatValue struct {
	Value float64
}

// Type returns TypeFloat
func (FloatValue) Type() Type {
	return TypeFloat
}

var _ Value = FloatValue{}

// CallValue holds the data of a function call
type CallValue struct {
	Func byte
	Args []Value
}

// Type returns TypeCall
func (CallValue) Type() Type {
	return TypeCall
}

var _ Value = CallValue{}

// StructValue holds the data of a struct instantiation
type StructValue struct {
	Fields []Value
}

// Type returns TypeStruct
func (StructValue) Type() Type {
	return TypeStruct
}

var _ Value = StructValue{}

// ArrayValue holds the data of an array instantiation
type ArrayValue struct {
	Elements []Value
}

// Type returns TypeArray
func (ArrayValue) Type() Type {
	return TypeArray
}

var _ Value = ArrayValue{}

// RefValue holds data of any reference
type RefValue struct {
	Ref Ref
}

// Type returns TypeRef
func (RefValue) Type() Type {
	return TypeRef
}

var _ Value = RefValue{}
