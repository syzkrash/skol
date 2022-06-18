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
	Ret      *values.Type
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
	return values.NewValue(args["a"].Int() + args["b"].Int()), nil
}

func NativeAddF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float() + args["b"].Float()), nil
}

func NativeAddC(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Char() + args["b"].Char()), nil
}

func NativeSubI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int() - args["b"].Int()), nil
}

func NativeSubF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float() - args["b"].Float()), nil
}

func NativeSubC(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Char() - args["b"].Char()), nil
}

func NativeMulI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int() * args["b"].Int()), nil
}

func NativeMulF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float() * args["b"].Float()), nil
}

func NativeDivI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int() / args["b"].Int()), nil
}

func NativeDivF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float() / args["b"].Float()), nil
}

func NativeModI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int() % args["b"].Int()), nil
}

func NativeModF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(int32(args["a"].Float()) % args["b"].Int()), nil
}

func NativeConcat(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].String() + args["b"].String()), nil
}

func NativeNot(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(!args["a"].Bool()), nil
}

func NativeOr(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Bool() || args["b"].Bool()), nil
}

func NativeAnd(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Bool() && args["b"].Bool()), nil
}

func NativeEq(s *Simulator, args ArgMap) (*values.Value, error) {
	a := args["a"]
	b := args["b"]
	return values.NewValue(a.Data == b.Data), nil
}

func NativeGtI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int() > args["b"].Int()), nil
}

func NativeGtF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Data.(float32) > args["b"].Float()), nil
}

func NativeGtC(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Char() > args["b"].Char()), nil
}

func NativeLtI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Int() < args["b"].Int()), nil
}

func NativeLtF(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Float() < args["b"].Float()), nil
}

func NativeLtC(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(args["a"].Char() < args["b"].Char()), nil
}

func NativeCharAt(s *Simulator, args ArgMap) (*values.Value, error) {
	str := args["s"].String()
	i := args["i"].Int()
	for i < 0 {
		i += int32(len(str))
	}
	for i > int32(len(str)) {
		i -= int32(len(str))
	}
	return &values.Value{values.Char, str[i]}, nil
}

func NativeSubstr(s *Simulator, args ArgMap) (*values.Value, error) {
	str := args["s"].String()
	a := args["a"].Int()
	b := args["b"].Int()
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
	return values.NewValue(args["s"].String() + string(args["c"].Char())), nil
}

func NativeStrLen(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(len(args["s"].String())), nil
}

func NativeCtoI(s *Simulator, args ArgMap) (*values.Value, error) {
	return values.NewValue(int(args["c"].Char())), nil
}

var DefaultFuncs = map[string]*Funct{
	"print": {
		Args:     []values.FuncArg{{"a", values.String}},
		Ret:      values.Nothing,
		IsNative: true,
		Native:   NativePrint,
	},
	"to_str": {
		Args:     []values.FuncArg{{"a", values.Any}},
		Ret:      values.String,
		IsNative: true,
		Native:   NativeToString,
	},
	"add_i": {
		Args:     []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:      values.Int,
		IsNative: true,
		Native:   NativeAddI,
	},
	"add_f": {
		Args:     []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:      values.Float,
		IsNative: true,
		Native:   NativeAddF,
	},
	"add_c": {
		Args:     []values.FuncArg{{"a", values.Char}, {"b", values.Char}},
		Ret:      values.Char,
		IsNative: true,
		Native:   NativeAddC,
	},
	"sub_i": {
		Args:     []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:      values.Int,
		IsNative: true,
		Native:   NativeSubI,
	},
	"sub_f": {
		Args:     []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:      values.Float,
		IsNative: true,
		Native:   NativeSubF,
	},
	"sub_c": {
		Args:     []values.FuncArg{{"a", values.Char}, {"b", values.Char}},
		Ret:      values.Char,
		IsNative: true,
		Native:   NativeSubC,
	},
	"mul_i": {
		Args:     []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:      values.Int,
		IsNative: true,
		Native:   NativeMulI,
	},
	"mul_f": {
		Args:     []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:      values.Float,
		IsNative: true,
		Native:   NativeMulF,
	},
	"div_i": {
		Args:     []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:      values.Int,
		IsNative: true,
		Native:   NativeDivI,
	},
	"div_f": {
		Args:     []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:      values.Float,
		IsNative: true,
		Native:   NativeDivF,
	},
	"mod_i": {
		Args:     []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:      values.Int,
		IsNative: true,
		Native:   NativeModI,
	},
	"mod_f": {
		Args:     []values.FuncArg{{"a", values.Float}, {"b", values.Int}},
		Ret:      values.Int,
		IsNative: true,
		Native:   NativeModF,
	},
	"concat": {
		Args:     []values.FuncArg{{"a", values.String}, {"b", values.String}},
		Ret:      values.String,
		IsNative: true,
		Native:   NativeConcat,
	},
	"not": {
		Args:     []values.FuncArg{{"a", values.Bool}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeNot,
	},
	"or": {
		Args:     []values.FuncArg{{"a", values.Bool}, {"b", values.Bool}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeOr,
	},
	"and": {
		Args:     []values.FuncArg{{"a", values.Bool}, {"b", values.Bool}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeAnd,
	},
	"eq": {
		Args:     []values.FuncArg{{"a", values.Any}, {"b", values.Any}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeEq,
	},
	"gt_i": {
		Args:     []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeGtI,
	},
	"gt_f": {
		Args:     []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeGtF,
	},
	"gt_c": {
		Args:     []values.FuncArg{{"a", values.Char}, {"b", values.Char}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeGtC,
	},
	"lt_i": {
		Args:     []values.FuncArg{{"a", values.Int}, {"b", values.Int}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeLtI,
	},
	"lt_f": {
		Args:     []values.FuncArg{{"a", values.Float}, {"b", values.Float}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeLtF,
	},
	"lt_c": {
		Args:     []values.FuncArg{{"a", values.Char}, {"b", values.Char}},
		Ret:      values.Bool,
		IsNative: true,
		Native:   NativeLtC,
	},
	"char_at": {
		Args:     []values.FuncArg{{"s", values.String}, {"i", values.Int}},
		Ret:      values.Char,
		IsNative: true,
		Native:   NativeCharAt,
	},
	"substr": {
		Args:     []values.FuncArg{{"s", values.String}, {"a", values.Int}, {"b", values.Int}},
		Ret:      values.String,
		IsNative: true,
		Native:   NativeSubstr,
	},
	"char_append": {
		Args:     []values.FuncArg{{"s", values.String}, {"c", values.Char}},
		Ret:      values.String,
		IsNative: true,
		Native:   NativeCharAppend,
	},
	"str_len": {
		Args:     []values.FuncArg{{"s", values.String}},
		Ret:      values.Int,
		IsNative: true,
		Native:   NativeStrLen,
	},
	"ctoi": {
		Args:     []values.FuncArg{{"c", values.Char}},
		Ret:      values.Int,
		IsNative: true,
		Native:   NativeCtoI,
	},
}
