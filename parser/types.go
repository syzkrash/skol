package parser

import (
	"fmt"
	"strings"

	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values/types"
)

func (p *Parser) parseType() (t types.Type, err error) {
	tk, err := p.lexer.Next()
	if err != nil {
		return
	}
	isArray := false
	if tk.Kind == lexer.TkPunct && tk.Raw[0] == '[' {
		isArray = true
		tk, err = p.lexer.Next()
		if err != nil {
			return
		}
	}
	if tk.Kind != lexer.TkIdent {
		err = p.selfError(tk, "expected identifier")
		return
	}
	switch strings.ToLower(tk.Raw) {
	case "integer", "int32", "int", "i32", "i":
		t = types.Int
	case "boolean", "bool", "b":
		t = types.Bool
	case "float32", "float", "f32", "f":
		t = types.Float
	case "char", "ch", "c":
		t = types.Char
	case "string", "str", "s":
		t = types.String
	case "any", "a":
		t = types.Any
	default:
		var ok bool
		t, ok = p.Scope.FindType(tk.Raw)
		if !ok {
			err = p.selfError(tk, "unknown type: "+tk.Raw)
			return
		}
	}
	if isArray {
		tk, err = p.lexer.Next()
		if err != nil {
			return
		}
		if tk.Kind != lexer.TkPunct || tk.Raw[0] != ']' {
			err = p.selfError(tk, "expected ']'")
			return
		}
		t = types.ArrayType{t}
	}
	return
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
