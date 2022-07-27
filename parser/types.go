package parser

import (
	"fmt"
	"strings"

	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values/types"
)

func (p *Parser) ParseType(raw string) (types.Type, bool) {
	switch strings.ToLower(raw) {
	case "integer", "int32", "int", "i32", "i":
		return types.Int, true
	case "boolean", "bool", "b":
		return types.Bool, true
	case "float32", "float", "f32", "f":
		return types.Float, true
	case "char", "ch", "c":
		return types.Char, true
	case "string", "str", "s":
		return types.String, true
	case "any", "a":
		return types.Any, true
	}
	if stype, ok := p.Scope.FindType(raw); ok {
		return stype, true
	}
	return nil, false
}

func (p *Parser) TypeOf(n nodes.Node) (t types.Type, err error) {
	switch n.Kind() {
	case nodes.NdBoolean:
		t = types.Bool
	case nodes.NdInteger:
		t = types.Int
	case nodes.NdFloat:
		t = types.Float
	case nodes.NdChar:
		t = types.Char
	case nodes.NdString:
		t = types.String
	case nodes.NdNewStruct:
		t = n.(*nodes.NewStructNode).Type
	case nodes.NdFuncCall:
		fn := n.(*nodes.FuncCallNode).Func
		f, ok := p.Scope.FindFunc(fn)
		if !ok {
			err = fmt.Errorf("unknown function: %s", fn)
			return
		}
		t = f.Ret
	case nodes.NdSelector:
		s := n.(*nodes.SelectorNode)
		path := s.Path()
		v, ok := p.Scope.FindVar(path[0])
		if !ok {
			err = fmt.Errorf("unknown variable: %s", path[0])
			return
		}
		t = v.VarType
		if len(path) == 1 {
			return
		}
	outer:
		for _, e := range path[1:] {
			if t.Prim() != types.PStruct {
				err = common.Error(n, "can only select fields on structures")
				return
			}
			for _, f := range t.(types.StructType).Fields {
				if f.Name == e {
					t = f.Type
					continue outer
				}
			}
			err = common.Error(n, "%s does not contain field '%s'", t.String(), e)
			return
		}
	case nodes.NdTypecast:
		return n.(*nodes.TypecastNode).Target, nil
	case nodes.NdArray:
		return types.ArrayType{n.(*nodes.ArrayNode).Type}, nil
	default:
		err = fmt.Errorf("%s node is not a value", n.Kind())
	}
	return
}
