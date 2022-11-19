package ir

import (
	"fmt"
	"strings"
)

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

var typeNames = []string{
	"INTEGER",
	"FLOAT",
	"CALL",
	"STRUCT",
	"ARRAY",
	"REF",
}

func (t Type) String() string {
	return typeNames[t]
}

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

func (v IntegerValue) String() string {
	return fmt.Sprintf("%s %d", TypeInteger, v.Value)
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

func (v FloatValue) String() string {
	return fmt.Sprintf("%s %f", TypeFloat, v.Value)
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

func (v CallValue) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "%s %02X [%02X](", TypeCall, v.Func, len(v.Args))
	for n := 0; n < len(v.Args)-1; n++ {
		fmt.Fprintf(&str, "%s, ", v.Args[n])
	}
	fmt.Fprintf(&str, "%s)", v.Args[len(v.Args)-1])
	return str.String()
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

func (v StructValue) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "%s [%02X](", TypeStruct, len(v.Fields))
	for n := 0; n < len(v.Fields)-1; n++ {
		fmt.Fprintf(&str, "%s, ", v.Fields[n])
	}
	fmt.Fprintf(&str, "%s)", v.Fields[len(v.Fields)-1])
	return str.String()
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

func (v ArrayValue) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "%s [%02X](", TypeArray, len(v.Elements))
	for n := 0; n < len(v.Elements)-1; n++ {
		fmt.Fprintf(&str, "%s, ", v.Elements[n])
	}
	fmt.Fprintf(&str, "%s)", v.Elements[len(v.Elements)-1])
	return str.String()
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

func (v RefValue) String() string {
	return fmt.Sprintf("%s %s", TypeRef, v.Ref)
}

var _ Value = RefValue{}
