package types

var (
	Bool      = PrimType{PBool}
	Char      = PrimType{PChar}
	Int       = PrimType{PInt}
	Float     = PrimType{PFloat}
	String    = PrimType{PString}
	Struct    = PrimType{PStruct}
	Any       = AnyType{}
	Nothing   = NothingType{false}
	Undefined = NothingType{true}
)
