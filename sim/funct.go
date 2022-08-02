package sim

import (
	"fmt"
	"os"

	"github.com/syzkrash/skol/parser/defaults"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

type ArgMap map[string]*values.Value
type NativeFunc func(*Simulator, ArgMap) (*values.Value, error)

type Funct struct {
	Args     []values.FuncArg
	Ret      types.Type
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

func NativeToString(s *Simulator, args ArgMap) (v *values.Value, err error) {
	return values.NewValue(args["a"].String()), nil
}

func NativeAddI(s *Simulator, args ArgMap) (v *values.Value, err error) {
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
	return &values.Value{types.Char, str[i]}, nil
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

func NativeOpen(s *Simulator, args ArgMap) (*values.Value, error) {
	fn := args["fn"].String()
	f, err := os.Open(fn)
	if err != nil {
		res := &values.Value{
			defaults.Types["file_descriptor_result"],
			map[string]*values.Value{
				"fd":  values.NewValue(0),
				"ok":  values.NewValue(false),
				"err": values.NewValue(err.Error()),
			},
		}
		return res, nil
	}
	fd := &values.Value{
		defaults.Types["file_descriptor"],
		map[string]*values.Value{
			"fd": values.NewValue(int(f.Fd())),
			"fn": values.NewValue(fn),
		},
	}
	res := &values.Value{
		defaults.Types["file_descriptor_result"],
		map[string]*values.Value{
			"fd":  fd,
			"ok":  values.NewValue(true),
			"err": values.NewValue(""),
		},
	}
	return res, nil
}

func NativeFGetC(s *Simulator, args ArgMap) (*values.Value, error) {
	fd := args["fd"].Struct()
	f := os.NewFile(uintptr(fd["fd"].Int()), fd["fn"].String())
	if f == nil {
		res := &values.Value{
			defaults.Types["char_result"],
			map[string]*values.Value{
				"char": values.NewValue('\000'),
				"ok":   values.NewValue(false),
				"err":  values.NewValue("operation on invalid file descriptor"),
			},
		}
		return res, nil
	}
	buf := []byte{0}
	_, err := f.Read(buf)
	if err != nil {
		res := &values.Value{
			defaults.Types["char_result"],
			map[string]*values.Value{
				"char": values.NewValue(byte(0)),
				"ok":   values.NewValue(false),
				"err":  values.NewValue(err.Error()),
			},
		}
		return res, nil
	}
	res := &values.Value{
		defaults.Types["char_result"],
		map[string]*values.Value{
			"char": values.NewValue(buf[0]),
			"ok":   values.NewValue(true),
			"err":  values.NewValue(""),
		},
	}
	return res, nil
}

func NativeClose(s *Simulator, args ArgMap) (*values.Value, error) {
	fd := args["fd"].Struct()
	f := os.NewFile(uintptr(fd["fd"].Int()), fd["fn"].String())
	f.Close()
	return nil, nil
}

var DefaultFuncs = map[string]*Funct{
	"print": {
		Args:     []values.FuncArg{{"a", types.String}},
		Ret:      types.Nothing,
		IsNative: true,
		Native:   NativePrint,
	},
	"to_str": {
		Args:     []values.FuncArg{{"a", types.Any}},
		Ret:      types.String,
		IsNative: true,
		Native:   NativeToString,
	},
	"add_i": {
		Args:     []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:      types.Int,
		IsNative: true,
		Native:   NativeAddI,
	},
	"add_f": {
		Args:     []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:      types.Float,
		IsNative: true,
		Native:   NativeAddF,
	},
	"add_c": {
		Args:     []values.FuncArg{{"a", types.Char}, {"b", types.Char}},
		Ret:      types.Char,
		IsNative: true,
		Native:   NativeAddC,
	},
	"sub_i": {
		Args:     []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:      types.Int,
		IsNative: true,
		Native:   NativeSubI,
	},
	"sub_f": {
		Args:     []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:      types.Float,
		IsNative: true,
		Native:   NativeSubF,
	},
	"sub_c": {
		Args:     []values.FuncArg{{"a", types.Char}, {"b", types.Char}},
		Ret:      types.Char,
		IsNative: true,
		Native:   NativeSubC,
	},
	"mul_i": {
		Args:     []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:      types.Int,
		IsNative: true,
		Native:   NativeMulI,
	},
	"mul_f": {
		Args:     []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:      types.Float,
		IsNative: true,
		Native:   NativeMulF,
	},
	"div_i": {
		Args:     []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:      types.Int,
		IsNative: true,
		Native:   NativeDivI,
	},
	"div_f": {
		Args:     []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:      types.Float,
		IsNative: true,
		Native:   NativeDivF,
	},
	"mod_i": {
		Args:     []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:      types.Int,
		IsNative: true,
		Native:   NativeModI,
	},
	"mod_f": {
		Args:     []values.FuncArg{{"a", types.Float}, {"b", types.Int}},
		Ret:      types.Int,
		IsNative: true,
		Native:   NativeModF,
	},
	"concat": {
		Args:     []values.FuncArg{{"a", types.String}, {"b", types.String}},
		Ret:      types.String,
		IsNative: true,
		Native:   NativeConcat,
	},
	"not": {
		Args:     []values.FuncArg{{"a", types.Bool}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeNot,
	},
	"or": {
		Args:     []values.FuncArg{{"a", types.Bool}, {"b", types.Bool}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeOr,
	},
	"and": {
		Args:     []values.FuncArg{{"a", types.Bool}, {"b", types.Bool}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeAnd,
	},
	"eq": {
		Args:     []values.FuncArg{{"a", types.Any}, {"b", types.Any}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeEq,
	},
	"gt_i": {
		Args:     []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeGtI,
	},
	"gt_f": {
		Args:     []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeGtF,
	},
	"gt_c": {
		Args:     []values.FuncArg{{"a", types.Char}, {"b", types.Char}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeGtC,
	},
	"lt_i": {
		Args:     []values.FuncArg{{"a", types.Int}, {"b", types.Int}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeLtI,
	},
	"lt_f": {
		Args:     []values.FuncArg{{"a", types.Float}, {"b", types.Float}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeLtF,
	},
	"lt_c": {
		Args:     []values.FuncArg{{"a", types.Char}, {"b", types.Char}},
		Ret:      types.Bool,
		IsNative: true,
		Native:   NativeLtC,
	},
	"char_at": {
		Args:     []values.FuncArg{{"s", types.String}, {"i", types.Int}},
		Ret:      types.Char,
		IsNative: true,
		Native:   NativeCharAt,
	},
	"substr": {
		Args:     []values.FuncArg{{"s", types.String}, {"a", types.Int}, {"b", types.Int}},
		Ret:      types.String,
		IsNative: true,
		Native:   NativeSubstr,
	},
	"char_append": {
		Args:     []values.FuncArg{{"s", types.String}, {"c", types.Char}},
		Ret:      types.String,
		IsNative: true,
		Native:   NativeCharAppend,
	},
	"str_len": {
		Args:     []values.FuncArg{{"s", types.String}},
		Ret:      types.Int,
		IsNative: true,
		Native:   NativeStrLen,
	},
	"ctoi": {
		Args:     []values.FuncArg{{"c", types.Char}},
		Ret:      types.Int,
		IsNative: true,
		Native:   NativeCtoI,
	},
	"open": {
		Args:     defaults.Functions["open"].Args,
		Ret:      defaults.Functions["open"].Ret,
		IsNative: true,
		Native:   NativeOpen,
	},
	"fgetc": {
		Args:     defaults.Functions["fgetc"].Args,
		Ret:      defaults.Functions["fgetc"].Ret,
		IsNative: true,
		Native:   NativeFGetC,
	},
	"close": {
		Args:     defaults.Functions["close"].Args,
		Ret:      defaults.Functions["close"].Ret,
		IsNative: true,
		Native:   NativeClose,
	},
}
