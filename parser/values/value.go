package values

import (
	"fmt"
)

type Value struct {
	ValueType ValueType
	Int       int32
	Bool      bool
	Float     float32
	Str       string
	Char      byte
}

func NewValue(v any) *Value {
	if v == nil {
		return &Value{
			ValueType: VtNothing,
		}
	}
	if i, ok := v.(int); ok {
		return &Value{
			ValueType: VtInteger,
			Int:       int32(i),
		}
	}
	if i, ok := v.(int32); ok {
		return &Value{
			ValueType: VtInteger,
			Int:       i,
		}
	}
	if b, ok := v.(bool); ok {
		return &Value{
			ValueType: VtBool,
			Bool:      b,
		}
	}
	if f, ok := v.(float32); ok {
		return &Value{
			ValueType: VtFloat,
			Float:     f,
		}
	}
	if s, ok := v.(string); ok {
		return &Value{
			ValueType: VtString,
			Str:       s,
		}
	}
	if c, ok := v.(byte); ok {
		return &Value{
			ValueType: VtChar,
			Char:      c,
		}
	}
	panic(fmt.Sprintf("unexpected value for NewValue: %v", v))
}

func Default(t ValueType) *Value {
	v := &Value{
		ValueType: t,
	}
	switch t {
	case VtInteger:
		v.Int = 0
	case VtBool:
		v.Bool = false
	case VtFloat:
		v.Float = 0.0
	case VtString:
		v.Str = ""
	case VtChar:
		v.Char = '\x00'
	}
	return v
}

func (v *Value) ToBool() bool {
	switch v.ValueType {
	case VtBool:
		return v.Bool
	case VtInteger:
		return v.Int > 0
	case VtNothing:
		return false
	default:
		return true
	}
}

func (v *Value) String() string {
	switch v.ValueType {
	case VtNothing:
		return "Nothing"
	case VtInteger:
		return fmt.Sprint(v.Int)
	case VtBool:
		return fmt.Sprint(v.Bool)
	case VtFloat:
		return fmt.Sprint(v.Float)
	case VtString:
		return v.Str
	case VtChar:
		return string(v.Char)
	}
	return fmt.Sprintf("<%s value>", v.ValueType)
}
