package sim

import (
	"fmt"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type ArgMap map[string]*values.Value
type NativeFunc func(*Simulator, ArgMap) (*values.Value, error)

type Funct struct {
	Args     []values.FuncArg
	Ret      values.ValueType
	Body     []nodes.Node
	IsNative bool
	Native   NativeFunc
}

func NativeDefault(*Simulator, ArgMap) (*values.Value, error) {
	return nil, nil
}

func NativePrint(s *Simulator, args ArgMap) (*values.Value, error) {
	fmt.Println(args["a"].String())
	return nil, nil
}

func NativeToString(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].String()), nil
}

func NativeAddI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int + args["b"].Int), nil
}

func NativeAddF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float + args["b"].Float), nil
}

func NativeSubI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int - args["b"].Int), nil
}

func NativeSubF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float - args["b"].Float), nil
}

func NativeMulI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int * args["b"].Int), nil
}

func NativeMulF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float * args["b"].Float), nil
}

func NativeDivI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int / args["b"].Int), nil
}

func NativeDivF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float / args["b"].Float), nil
}

func NativeModI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int % args["b"].Int), nil
}

func NativeModF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(int32(args["a"].Float) % args["b"].Int), nil
}

func NativeConcat(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Str + args["b"].Str), nil
}

func NativeNot(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(!args["a"].Bool), nil
}

func NativeOr(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Bool || args["b"].Bool), nil
}

func NativeAnd(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Bool && args["b"].Bool), nil
}

func NativeEq(s *Simulator, args ArgMap) (*values.Value, error) {
	a := args["a"]
	b := args["b"]
	if a.ValueType != b.ValueType {
		return values.NewValue(false), nil
	}
	switch a.ValueType {
	case values.VtInteger:
		return values.NewValue(a.Int == b.Int), nil
	case values.VtBool:
		return values.NewValue(a.Bool == b.Bool), nil
	case values.VtFloat:
		return values.NewValue(a.Float == b.Float), nil
	case values.VtString:
		return values.NewValue(a.Str == b.Str), nil
	case values.VtChar:
		return values.NewValue(a.Char == b.Char), nil
	}
	return values.NewValue(false), nil
}

func NativeGtI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int > args["b"].Int), nil
}

func NativeGtF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float > args["b"].Float), nil
}

func NativeLtI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int < args["b"].Int), nil
}

func NativeLtF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float < args["b"].Float), nil
}

func NativeCharAt(s *Simulator, args ArgMap) (*values.Value, error) {
	str := args["s"].Str
	i := args["i"].Int
	for i < 0 {
		i += int32(len(str))
	}
	for i > int32(len(str)) {
		i -= int32(len(str))
	}
	return values.NewValue(rune(str[i])), nil
}

func NativeSubstr(s *Simulator, args ArgMap) (*values.Value, error) {
	str := args["s"].Str
	a := args["a"].Int
	b := args["b"].Int
	for a < 0 {
		a += int32(len(str))
	}
	for a > int32(len(str)) {
		a -= int32(len(str))
	}
	for b < 0 {
		b += int32(len(str))
	}
	for b > int32(len(str)) {
		b -= int32(len(str))
	}
	return values.NewValue(str[a:b]), nil
}

func NativeCharAppend(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["s"].Str + string(args["c"].Char)), nil
}

func NativeStrLen(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(len(args["s"].Str)), nil
}

var DefaultFuncs = map[string]*Funct{
	"print": {
		Args:     []values.FuncArg{{"a", values.VtString}},
		Ret:      values.VtNothing,
		IsNative: true,
		Native:   NativePrint,
	},
	"to_str": {
		Args:     []values.FuncArg{{"a", values.VtAny}},
		Ret:      values.VtString,
		IsNative: true,
		Native:   NativeToString,
	},
	"add_i": {
		Args:     []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:      values.VtInteger,
		IsNative: true,
		Native:   NativeAddI,
	},
	"add_f": {
		Args:     []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:      values.VtFloat,
		IsNative: true,
		Native:   NativeAddF,
	},
	"sub_i": {
		Args:     []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:      values.VtInteger,
		IsNative: true,
		Native:   NativeSubI,
	},
	"sub_f": {
		Args:     []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:      values.VtFloat,
		IsNative: true,
		Native:   NativeSubF,
	},
	"mul_i": {
		Args:     []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:      values.VtInteger,
		IsNative: true,
		Native:   NativeMulI,
	},
	"mul_f": {
		Args:     []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:      values.VtFloat,
		IsNative: true,
		Native:   NativeMulF,
	},
	"div_i": {
		Args:     []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:      values.VtInteger,
		IsNative: true,
		Native:   NativeDivI,
	},
	"div_f": {
		Args:     []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:      values.VtFloat,
		IsNative: true,
		Native:   NativeDivF,
	},
	"mod_i": {
		Args:     []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:      values.VtInteger,
		IsNative: true,
		Native:   NativeModI,
	},
	"mod_f": {
		Args:     []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtInteger}},
		Ret:      values.VtInteger,
		IsNative: true,
		Native:   NativeModF,
	},
	"concat": {
		Args:     []values.FuncArg{{"a", values.VtString}, {"b", values.VtString}},
		Ret:      values.VtString,
		IsNative: true,
		Native:   NativeConcat,
	},
	"not": {
		Args:     []values.FuncArg{{"a", values.VtBool}},
		Ret:      values.VtBool,
		IsNative: true,
		Native:   NativeNot,
	},
	"or": {
		Args:     []values.FuncArg{{"a", values.VtBool}, {"b", values.VtBool}},
		Ret:      values.VtBool,
		IsNative: true,
		Native:   NativeOr,
	},
	"and": {
		Args:     []values.FuncArg{{"a", values.VtBool}, {"b", values.VtBool}},
		Ret:      values.VtBool,
		IsNative: true,
		Native:   NativeAnd,
	},
	"eq": {
		Args:     []values.FuncArg{{"a", values.VtAny}, {"b", values.VtAny}},
		Ret:      values.VtBool,
		IsNative: true,
		Native:   NativeEq,
	},
	"gt_i": {
		Args:     []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:      values.VtBool,
		IsNative: true,
		Native:   NativeGtI,
	},
	"gt_f": {
		Args:     []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:      values.VtBool,
		IsNative: true,
		Native:   NativeGtF,
	},
	"lt_i": {
		Args:     []values.FuncArg{{"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:      values.VtBool,
		IsNative: true,
		Native:   NativeLtI,
	},
	"lt_f": {
		Args:     []values.FuncArg{{"a", values.VtFloat}, {"b", values.VtFloat}},
		Ret:      values.VtBool,
		IsNative: true,
		Native:   NativeLtF,
	},
	"char_at": {
		Args:     []values.FuncArg{{"s", values.VtString}, {"i", values.VtInteger}},
		Ret:      values.VtChar,
		IsNative: true,
		Native:   NativeCharAt,
	},
	"substr": {
		Args:     []values.FuncArg{{"s", values.VtString}, {"a", values.VtInteger}, {"b", values.VtInteger}},
		Ret:      values.VtString,
		IsNative: true,
		Native:   NativeSubstr,
	},
	"char_append": {
		Args:     []values.FuncArg{{"s", values.VtString}, {"c", values.VtChar}},
		Ret:      values.VtString,
		IsNative: true,
		Native:   NativeCharAppend,
	},
	"str_len": {
		Args:     []values.FuncArg{{"s", values.VtString}},
		Ret:      values.VtInteger,
		IsNative: true,
		Native:   NativeStrLen,
	},
}
