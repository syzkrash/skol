package ir

// RefType represents the unique type identifier of a reference
type RefType byte

// Reference type constants
const (
	RefLocal RefType = iota
	RefGlobal
	RefLocalIdx
	RefGlobalIdx
)

// Ref holds data of some value reference
type Ref interface {
	Type() RefType
}

// SingleRef is a reference with a single unique identifier
type SingleRef struct {
	RefType RefType
	Idx     byte
}

// Type returns this reference's underlying type
func (r SingleRef) Type() RefType {
	return r.RefType
}

var _ Ref = SingleRef{}

// DoubleRef is a reference with two unique identifiers
type DoubleRef struct {
	RefType RefType
	Val     byte
	Idx     uint32
}

// Type returns this reference's underlying type
func (r DoubleRef) Type() RefType {
	return r.RefType
}

var _ Ref = DoubleRef{}
