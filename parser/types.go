package parser

import (
	"fmt"
	"strings"

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

// EnsureResultType makes sure that a given result type exists in the current
// scope, creating one if it doesn't exist
func (p *Parser) EnsureResultType(inner types.Type) types.Type {
	// result type name
	tn := inner.String() + "Result"

	// check if it already exists and return it if it does
	if rt, ok := p.Scope.FindType(tn); ok {
		return rt
	}

	// otherwise, create a new one, add it to the current scope and return it
	rt := types.MakeStruct(tn,
		"ok", types.Bool,
		"val", inner)
	p.Scope.Types[tn] = rt
	return rt
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

		// first, ensure the path is correct
		if len(path) < 1 {
			err = fmt.Errorf("selector has path length of 0; this should never happen!")
			return
		}

		// get the root of the path, which should be a variable name
		root := path[0]
		if root.Name == "" {
			err = fmt.Errorf("expected first element of selector to be a name")
			return
		}

		// get the variable
		v, ok := p.Scope.FindVar(root.Name)
		if !ok {
			err = fmt.Errorf("unknown variable: %s", root.Name)
		}

		// set the type to the variable's type (if it is found) and return if it is
		// the only element of the path
		t = v.VarType
		if len(path) == 1 {
			return
		}

		// iterate through the rest of the elements
		for _, e := range path[1:] {
			// typecasts and fields take priority over indexing because an index does
			// not have an invalid/empty value:
			//	0   is a valid index
			//	""  is NOT a valid name
			//  nil is NOT a valid type

			// typecast because the Cast is not empty
			if e.Cast != nil {
				// checking type compatibility is not our job, let the typechecker
				// figure that out
				t = e.Cast
				continue
			}
			// field selection because the Name is not empty
			if e.Name != "" {
				// make sure we are selecting fields on a structure
				if t.Prim() != types.PStruct {
					err = fmt.Errorf("can only select fields on structures (you are selecting field '%s' on %s)", e.Name, t.String())
					return
				}
				// now, ensure the structure contains the given field and update our
				// current type accordingly
				s := t.(types.StructType)
				ok := false
				for _, f := range s.Fields {
					if f.Name == e.Name {
						ok = true
						t = f.Type
						break
					}
				}
				// exit if the field was not found
				if !ok {
					err = fmt.Errorf("%s does not contain field '%s'", t.String(), e.Name)
					return
				}
				// otherwise we can just move on to the next element (or exit the loop)
				continue
			}
			// finally, if neither are specified, process array index
			if t.Prim() != types.PArray {
				err = fmt.Errorf("can only index arrays (you are getting index %d of %s)", e.Idx, t.String())
				return
			}
			// return a result type for the array's element type (because s a f e t y)
			a := t.(types.ArrayType)
			t = p.EnsureResultType(a.Element)
		}
	case nodes.NdTypecast:
		return n.(*nodes.TypecastNode).Type, nil
	case nodes.NdArray:
		return types.ArrayType{n.(*nodes.ArrayNode).Type}, nil
	case nodes.NdIndex:
		i := n.(*nodes.IndexNode)
		ptype, err := p.TypeOf(i.Parent)
		if err != nil {
			return nil, err
		}
		if ptype.Prim() != types.PArray {
			return nil, fmt.Errorf("cannot index %s value", ptype.String())
		}
		atype := ptype.(types.ArrayType)
		rtype := types.MakeStruct(atype.Element.String()+"Result",
			"val", atype.Element,
			"ok", types.Bool)
		p.Scope.Types[rtype.(types.StructType).Name] = rtype
		return rtype, nil
	default:
		err = fmt.Errorf("%s node is not a value", n.Kind())
	}
	return
}
