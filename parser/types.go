package parser

import (
	"fmt"
	"strings"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

// typeByName retrieves a type for the given name. If the name refers to a
// built-in type, that type always takes priority over user-defined types.
func (p *Parser) typeByName(name string) (t types.Type, ok bool) {
	ok = true
	switch strings.ToLower(name) {
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
		t, ok = p.Scope.FindType(name)
	}
	return
}

// parseType parses a type name.
//
// Built-in type:
//
//	b bool boolean
//	c ch char
//	i i32 int int32 integer
//	f f32 float float32
//	s str string
//	a any
//
// User-defined structure type (assuming its name is Vec2i):
//
//	Vec2i
//
// Array type:
//
//	[integer]
//	[Vec2i]
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
		err = tokErr(pe.EExpectedName, tk)
		return
	}
	t, ok := p.typeByName(tk.Raw)
	if !ok {
		err = tokErr(pe.EUnknownType, tk)
		return
	}
	if isArray {
		tk, err = p.lexer.Next()
		if err != nil {
			return
		}
		if tk.Kind != lexer.TkPunct || tk.Raw[0] != ']' {
			err = tokErr(pe.EExpectedRBrack, tk)
			return
		}
		t = types.ArrayType{Element: t}
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

// TypeOf determines the type of the value denoted by the given function in
// the parser's current scope.
func (p *Parser) TypeOf(n ast.Node) (t types.Type, err error) {
	switch n.Kind() {
	case ast.NBool:
		t = types.Bool
	case ast.NInt:
		t = types.Int
	case ast.NFloat:
		t = types.Float
	case ast.NChar:
		t = types.Char
	case ast.NString:
		t = types.String
	case ast.NStruct:
		t = n.(ast.StructNode).Type
	case ast.NFuncCall:
		fn := n.(ast.FuncCallNode).Func
		f, ok := p.getFunc(fn)
		if !ok {
			err = fmt.Errorf("unknown function: %s", fn)
			return
		}
		t = f.Ret
	case ast.NSelector:
		s := n.(ast.SelectorNode)
		path := s.Path()

		// first, ensure the path is correct
		if len(path) < 1 {
			err = fmt.Errorf("selector has path length of 0; this should never happen")
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
			return
		}

		// set the type to the variable's type (if it is found) and return if it is
		// the only element of the path
		t, err = p.TypeOf(v)
		if err != nil {
			return
		}
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
				err = fmt.Errorf("can only index arrays")
				return
			}
			// return a result type for the array's element type (because s a f e t y)
			a := t.(types.ArrayType)
			t = p.EnsureResultType(a.Element)
		}
	case ast.NTypecast:
		return n.(ast.TypecastNode).Cast, nil
	case ast.NArray:
		return n.(ast.ArrayNode).Type, nil
	case ast.NIndexSelector:
		i := n.(ast.IndexSelectorNode)
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
	case ast.NIndexConst:
		i := n.(ast.IndexConstNode)
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

// NodeOf creates a node for the given type. This node is not correctly set up
// and cannot be used like a regular, parsed node!
func (p *Parser) NodeOf(t types.Type) (n ast.Node, ok bool) {
	ok = true
	if types.Bool.Equals(t) {
		n = ast.BoolNode{}
	} else if types.Char.Equals(t) {
		n = ast.CharNode{}
	} else if types.Int.Equals(t) {
		n = ast.IntNode{}
	} else if types.Float.Equals(t) {
		n = ast.FloatNode{}
	} else if types.String.Equals(t) {
		n = ast.StringNode{}
	} else if t.Prim() == types.PArray {
		n = ast.ArrayNode{
			Type: t.(types.ArrayType),
		}
	} else if t.Prim() == types.PStruct {
		n = ast.StructNode{
			Type: t.(types.StructType),
		}
	} else {
		ok = false
	}
	return
}
