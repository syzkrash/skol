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

type Structure struct {
	Name   string
	Fields []*Field
}

type Type struct {
	Prim      Primitive
	Structure *Structure
}

func Struct(name string, fields []*Field) *Type {
	return &Type{PStruct, &Structure{name, fields}}
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
	if a.Prim != PAny && a.Prim != b.Prim {
		return false
	}
	if a.Prim != PStruct {
		return true
	}
	// two structure types are equal if they contain the same fields, allowing for
	// semi-generic code
	bMap := map[string]*Type{}
	for _, f := range b.Structure.Fields {
		bMap[f.Name] = f.Type
	}
	for _, f := range a.Structure.Fields {
		bt, ok := bMap[f.Name]
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

func (t *Type) Name() string {
	switch t.Prim {
	case PStruct:
		return "Structure " + t.Structure.Name
	default:
		return []string{
			"Nothing",
			"Integer",
			"Boolean",
			"Float",
			"Character",
			"String",
			"Structure", // never gonna happen but whatever
			"Array",
			"Any",
			"Undefined",
		}[t.Prim]
	}
}
