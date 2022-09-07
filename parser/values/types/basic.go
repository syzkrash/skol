package types

// Definitions of built-in and internal types.
var (
	Bool      = PrimType{PBool}
	Char      = PrimType{PChar}
	Int       = PrimType{PInt}
	Float     = PrimType{PFloat}
	String    = PrimType{PString}
	Any       = AnyType{}
	Nothing   = NothingType{false}
	Undefined = NothingType{true}
)
