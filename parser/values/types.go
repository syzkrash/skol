package values

import "fmt"

type Primitive uint8

const (
	PNothing Primitive = iota
	PInt
	PBool
	PFloat
	PChar
	PString
	PStruct
	PArray
	PAny
	PUndefined
)

type Field struct {
	Name string
	Type *Type
}

type Type struct {
	Prim      Primitive
	Structure []*Field
}

var (
	Nothing   = &Type{PNothing, nil}
	Int       = &Type{PInt, nil}
	Bool      = &Type{PBool, nil}
	Float     = &Type{PFloat, nil}
	Char      = &Type{PChar, nil}
	String    = &Type{PString, nil}
	Any       = &Type{PAny, nil}
	Undefined = &Type{PUndefined, nil}
)

func (a *Type) Equals(b *Type) bool {
	if a.Prim != PAny && b.Prim != PAny && a.Prim != b.Prim {
		return false
	}
	// two structure types are equal if they contain the same fields, allowing for
	// semi-generic code
	aMap := map[string]*Type{}
	for _, f := range a.Structure {
		aMap[f.Name] = f.Type
	}
	for _, f := range b.Structure {
		bt, ok := aMap[f.Name]
		if !ok {
			return false
		}
		if f.Type != bt {
			return false
		}
	}
	return true
}

func (t *Type) String() string {
	switch t.Prim {
	case PNothing:
		return "Nothing"
	case PInt:
		return "Integer"
	case PBool:
		return "Boolean"
	case PFloat:
		return "Float"
	case PChar:
		return "Character"
	case PString:
		return "String"
	case PStruct:
		return "Structure"
	case PAny:
		return "Any"
	case PArray:
		return "Array"
	}
	panic(fmt.Sprintf("type with invalid primitive: %d", t.Prim))
}
