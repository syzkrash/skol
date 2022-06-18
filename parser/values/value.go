package values

import (
	"fmt"
)

type Value struct {
	Type *Type
	Data any
}

func NewValue(v any) *Value {
	if v == nil {
		return &Value{Nothing, nil}
	}
	if i, ok := v.(int); ok {
		return &Value{Int, int32(i)}
	}
	if i, ok := v.(int32); ok {
		return &Value{Int, i}
	}
	if b, ok := v.(bool); ok {
		return &Value{Bool, b}
	}
	if f, ok := v.(float32); ok {
		return &Value{Float, f}
	}
	if s, ok := v.(string); ok {
		return &Value{String, s}
	}
	if c, ok := v.(byte); ok {
		return &Value{Char, c}
	}
	panic(fmt.Sprintf("unexpected value for NewValue: %v", v))
}

func Default(t *Type) *Value {
	switch t.Prim {
	case PNothing:
		return &Value{t, nil}
	case PInt:
		return &Value{t, 0}
	case PBool:
		return &Value{t, false}
	case PFloat:
		return &Value{t, 0}
	case PChar:
		return &Value{t, 0}
	case PString:
		return &Value{t, ""}
	case PStruct:
		v := map[string]*Value{}
		for _, f := range t.Structure {
			v[f.Name] = Default(f.Type)
		}
		return &Value{t, v}
	case PArray:
		return &Value{t, []*Value{}}
	}
	panic(fmt.Sprintf("invalid value primitive: %d", t.Prim))
}

func (v *Value) ToBool() bool {
	switch v.Type.Prim {
	case PBool:
		return v.Data.(bool)
	case PInt:
		return v.Data.(int32) > 0
	case PNothing:
		return false
	default:
		return true
	}
}

func (v *Value) String() string {
	switch v.Type.Prim {
	case PNothing:
		return "Nothing"
	case PInt:
		return fmt.Sprint(v.Data.(int32))
	case PBool:
		return fmt.Sprint(v.Data.(bool))
	case PFloat:
		return fmt.Sprint(v.Data.(float32))
	case PString:
		return v.Data.(string)
	case PChar:
		return string(v.Data.(byte))
	}
	return fmt.Sprintf("<%s value>", v.Type)
}

func (v *Value) Int() int32 {
	if v.Type.Prim != PInt {
		panic("Int() call to value of type " + v.Type.String())
	}
	return v.Data.(int32)
}

func (v *Value) Bool() bool {
	if v.Type.Prim != PBool {
		panic("Bool() call to value of type " + v.Type.String())
	}
	return v.Data.(bool)
}

func (v *Value) Float() float32 {
	if v.Type.Prim != PFloat {
		panic("Float() call to value of type " + v.Type.String())
	}
	return v.Data.(float32)
}

func (v *Value) Char() byte {
	if v.Type.Prim != PChar {
		panic("Char() call to value of type " + v.Type.String())
	}
	return v.Data.(byte)
}

func (v *Value) Struct() map[string]*Value {
	if v.Type.Prim != PStruct {
		panic("Struct() call to value of type " + v.Type.String())
	}
	return v.Data.(map[string]*Value)
}

func (v *Value) Array() []*Value {
	if v.Type.Prim != PArray {
		panic("Array() call to value of type " + v.Type.String())
	}
	return v.Data.([]*Value)
}
