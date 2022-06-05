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
	t := []any{}
	for _, a := range args {
		t = append(t, a.String())
	}
	fmt.Println(t...)
	return nil, nil
}
