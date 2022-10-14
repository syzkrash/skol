package types

// StructType represents all structure types with the primtive [PStruct].
type StructType struct {
	Name   string
	Fields []Descriptor
}

func (s StructType) String() string {
	return "Structure " + s.Name
}

func (StructType) Prim() Primitive {
	return PStruct
}

func (s StructType) FieldType(name string) (Type, bool) {
	for _, f := range s.Fields {
		if f.Name == name {
			return f.Type, true
		}
	}
	return nil, false
}

// Equals ensures the other type is compatible with this type. This function
// is especially important for structures. If structure B contains all the
// fields that A does (A ⊂ B) then it is compatible.
func (a StructType) Equals(b Type) bool {
	if b.Prim() != PStruct {
		return false
	}
	bs := b.(StructType)
	bf := map[string]Type{}
	for _, f := range bs.Fields {
		bf[f.Name] = f.Type
	}
	for _, f := range a.Fields {
		bt, ok := bf[f.Name]
		if !ok {
			return false
		}
		if !f.Type.Equals(bt) {
			return false
		}
	}
	return true
}

// MakeStruct creates a [StructType] from the given field name/type pairs.
func MakeStruct(name string, fields ...any) Type {
	if len(fields)%2 != 0 {
		panic("MakeStruct requires an even amount of arguments")
	}
	f := make([]Descriptor, len(fields)/2)
	for i := 0; i < len(fields); i += 2 {
		f[i/2] = Descriptor{
			Name: fields[i].(string),
			Type: fields[i+1].(Type),
		}
	}
	return StructType{
		Name:   name,
		Fields: f,
	}
}

// Result makes a result type wrapping the given type
func Result(t Type) Type {
	return MakeStruct(t.String()+"Result", "ok", Bool, "value", t)
}
