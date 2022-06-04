package sim

import (
	"fmt"

	"github.com/syzkrash/skol/parser"
)

type Value struct {
	ValueType parser.ValueType
	Int       int32
	Bool      bool
	Float     float32
	Str       string
	Char      rune
}

func NewValue(v any) *Value {
	if v == nil {
		return &Value{
			ValueType: parser.VtNothing,
		}
	}
	if i, ok := v.(int32); ok {
		return &Value{
			ValueType: parser.VtInteger,
			Int:       i,
		}
	}
	if b, ok := v.(bool); ok {
		return &Value{
			ValueType: parser.VtBool,
			Bool:      b,
		}
	}
	if f, ok := v.(float32); ok {
		return &Value{
			ValueType: parser.VtFloat,
			Float:     f,
		}
	}
	if s, ok := v.(string); ok {
		return &Value{
			ValueType: parser.VtString,
			Str:       s,
		}
	}
	if r, ok := v.(rune); ok {
		return &Value{
			ValueType: parser.VtChar,
			Char:      r,
		}
	}
	panic(fmt.Sprintf("unexpected value for NewValue: %v", v))
}

func Default(t parser.ValueType) *Value {
	v := &Value{
		ValueType: t,
	}
	switch t {
	case parser.VtInteger:
		v.Int = 0
	case parser.VtBool:
		v.Bool = false
	case parser.VtFloat:
		v.Float = 0.0
	case parser.VtString:
		v.Str = ""
	case parser.VtChar:
		v.Char = '\x00'
	}
	return v
}

func (v *Value) ToBool() bool {
	switch v.ValueType {
	case parser.VtBool:
		return v.Bool
	case parser.VtInteger:
		return v.Int > 0
	case parser.VtNothing:
		return false
	default:
		return true
	}
}

func (v *Value) String() string {
	switch v.ValueType {
	case parser.VtNothing:
		return "Nothing"
	case parser.VtInteger:
		return fmt.Sprint(v.Int)
	case parser.VtBool:
		return fmt.Sprint(v.Bool)
	case parser.VtFloat:
		return fmt.Sprint(v.Float)
	case parser.VtString:
		return v.Str
	case parser.VtChar:
		return string(v.Char)
	}
	return fmt.Sprintf("<%s value>", v.ValueType)
}
