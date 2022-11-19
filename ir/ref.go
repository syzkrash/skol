package ir

import "fmt"

// RefType represents the unique type identifier of a reference
type RefType byte

// Reference type constants
const (
	RefLocal RefType = iota
	RefGlobal
	RefLocalIdx
	RefGlobalIdx
)

var refNames = []string{
	"LOCAL",
	"GLOBAL",
	"LOCAL.IDX",
	"GLOBAL.IDX",
}

func (rt RefType) String() string {
	return refNames[rt]
}

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

func (r SingleRef) String() string {
	return fmt.Sprintf("%s %02X", r.RefType, r.Idx)
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

func (r DoubleRef) String() string {
	return fmt.Sprintf("%s %02X$%08X", r.RefType, r.Val, r.Idx)
}

var _ Ref = DoubleRef{}
