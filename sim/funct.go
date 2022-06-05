package sim

import (
	"fmt"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type ArgMap map[string]*values.Value
type NativeFunc func(*Simulator, ArgMap) (*values.Value, error)

type Funct struct {
	Args     map[string]values.ValueType
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

var DefaultFuncs = map[string]*Funct{
	"print": {
		Args:     map[string]values.ValueType{"a": values.VtString},
		Ret:      values.VtNothing,
		IsNative: true,
		Native:   NativePrint,
	},
	"to_str": {
		Args:     map[string]values.ValueType{"a": values.VtAny},
		Ret:      values.VtString,
		IsNative: true,
		Native:   NativeToString,
	},
	"add_i": {
		Args:     map[string]values.ValueType{"a": values.VtInteger, "b": values.VtInteger},
		Ret:      values.VtInteger,
		IsNative: true,
		Native:   NativeAddI,
	},
}
