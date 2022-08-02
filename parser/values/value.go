package values

import (
	"fmt"

	"github.com/syzkrash/skol/parser/values/types"
)

type Value struct {
	Type types.Type
	Data any
}

func NewValue(v any) *Value {
	switch v.(type) {
	case bool:
		return &Value{types.Bool, v}
	case byte:
		return &Value{types.Char, v}
	case int, int8, int16, int32, int64:
		return &Value{types.Char, v}
	case float32, float64:
		return &Value{types.Float, v}
	case string:
		return &Value{types.String, v}
	default:
		panic(fmt.Sprintf("Unexpected argument for NewValue: %v", v))
	}
}

func Default(t types.Type) *Value {
	switch t.Prim() {
	case types.PBool:
		return &Value{t, false}
	case types.PChar:
		return &Value{t, byte(0)}
	case types.PInt:
		return &Value{t, int(0)}
	case types.PFloat:
		return &Value{t, float32(0.0)}
	case types.PString:
		return &Value{t, ""}
	default:
		panic(fmt.Sprintf("Unexpected argument for Default: %v", t))
	}
}

func (v *Value) ToBool() bool {
	switch v.Type.Prim() {
	case types.PBool:
		return v.Data.(bool)
	case types.PInt:
		return v.Data.(int32) > 0
	default:
		return true
	}
}

func (v *Value) String() string {
	switch v.Type.Prim() {
	case types.PInt:
		return fmt.Sprint(v.Data.(int32))
	case types.PBool:
		return fmt.Sprint(v.Data.(bool))
	case types.PFloat:
		return fmt.Sprint(v.Data.(float32))
	case types.PString:
		return v.Data.(string)
	case types.PChar:
		return string(v.Data.(byte))
	}
	return fmt.Sprintf("<%s value>", v.Type)
}

func (v *Value) Int() int32 {
	if !v.Type.Equals(types.Int) {
		panic("Int() call to value of type " + v.Type.String())
	}
	return v.Data.(int32)
}

func (v *Value) Bool() bool {
	if !v.Type.Equals(types.Bool) {
		panic("Bool() call to value of type " + v.Type.String())
	}
	return v.Data.(bool)
}

func (v *Value) Float() float32 {
	if !v.Type.Equals(types.Float) {
		panic("Float() call to value of type " + v.Type.String())
	}
	return v.Data.(float32)
}

func (v *Value) Char() byte {
	if !v.Type.Equals(types.Char) {
		panic("Char() call to value of type " + v.Type.String())
	}
	return v.Data.(byte)
}

func (v *Value) Struct() map[string]*Value {
	if v.Type.Prim() != types.PStruct {
		panic("Struct() call to value of type " + v.Type.String())
	}
	return v.Data.(map[string]*Value)
}
