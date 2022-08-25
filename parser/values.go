package parser

import (
	"strconv"
	"strings"

	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

// Value parses a nodes.Node that has a Value
//
// Example values:
//
//	123        // nodes.IntegerNode
//	45.67      // nodes.FloatNode
//	"hello"    // nodes.StringNode
//	'E'        // nodes.CharNode
//	add! 1 2   // nodes.FuncCallNode
//	age        // nodes.VarRefNode
//	@VectorTwo 1.2 3.4 // nodes.NewStructNode
//	pos#x      // nodes.SelectorNode
//
func (p *Parser) Value() (n nodes.Node, err error) {
	tok, err := p.lexer.Next()
	if err != nil {
		return nil, err
	}

	switch tok.Kind {
	case lexer.TkConstant:
		if strings.ContainsRune(tok.Raw, '.') {
			var f float64
			f, err = strconv.ParseFloat(tok.Raw, 32)
			if err != nil {
				err = p.otherError(tok, "invalid floating-point constant", err)
			}

			n = &nodes.FloatNode{
				Float: float32(f),
				Pos:   tok.Where,
			}
		} else {
			var i int64
			i, err = strconv.ParseInt(tok.Raw, 0, 32)
			if err != nil {
				err = p.otherError(tok, "invalid integer constant", err)
			}

			n = &nodes.IntegerNode{
				Int: int32(i),
				Pos: tok.Where,
			}
		}
	case lexer.TkString:
		n = &nodes.StringNode{
			Str: tok.Raw,
			Pos: tok.Where,
		}
	case lexer.TkChar:
		rdr := strings.NewReader(tok.Raw)
		r, _, _ := rdr.ReadRune()
		n = &nodes.CharNode{
			Char: byte(r),
			Pos:  tok.Where,
		}
	case lexer.TkIdent:
		if tok.Raw[len(tok.Raw)-1] == '!' {
			fn := tok.Raw[:len(tok.Raw)-1]
			f, ok := p.Scope.FindFunc(fn)
			if !ok {
				err = p.selfError(tok, "unknown function: "+fn)
				return
			}
			n, err = p.funcCall(fn, f, tok.Where)
		} else if _, ok := p.Scope.FindVar(tok.Raw); ok {
			return p.selector(tok)
		} else if v, ok := p.Scope.FindConst(tok.Raw); ok {
			n = p.ToNode(v, tok.Where)
		} else {
			err = p.selfError(tok, "unknown variable: "+tok.Raw)
		}
	case lexer.TkPunct:
		switch tok.Raw[0] {
		case '*':
			n = &nodes.BooleanNode{
				Bool: true,
				Pos:  tok.Where,
			}
		case '/':
			n = &nodes.BooleanNode{
				Bool: false,
				Pos:  tok.Where,
			}
		case '@':
			tok, err = p.lexer.Next()
			if err != nil {
				return
			}
			if tok.Kind != lexer.TkIdent {
				err = p.selfError(tok, "expected Identifier, got "+tok.Kind.String())
				return
			}
			t, ok := p.Scope.FindType(tok.Raw)
			if !ok {
				err = p.selfError(tok, "unknown type: "+tok.Raw)
				return
			}
			s := t.(types.StructType)
			args := make([]nodes.Node, len(s.Fields))
			for i := range s.Fields {
				n, err = p.Value()
				if err != nil {
					return
				}
				args[i] = n
			}
			n = &nodes.NewStructNode{
				Type: t,
				Args: args,
				Pos:  tok.Where,
			}
		case '[':
			begin := tok
			var elemtype types.Type = types.Undefined
			tok, err = p.lexer.Next()
			if err != nil {
				return
			}
			if tok.Kind == lexer.TkIdent {
				var ok bool
				elemtype, ok = p.typeByName(tok.Raw)
				if !ok {
					err = p.selfError(tok, "unknown type: "+tok.Raw)
					return
				}
				tok, err = p.lexer.Next()
				if err != nil {
					return
				}
			}
			if tok.Kind != lexer.TkPunct || tok.Raw[0] != ']' {
				err = p.selfError(tok, "expected type name or ']'")
				return
			}
			tok, err = p.lexer.Next()
			if err != nil {
				return
			}
			if tok.Kind != lexer.TkPunct || tok.Raw[0] != '(' {
				err = p.selfError(tok, "expected '('")
				return
			}
			elems := []nodes.Node{}
			var elem nodes.Node
			for {
				tok, err = p.lexer.Next()
				if err != nil {
					return
				}
				if tok.Kind == lexer.TkPunct && tok.Raw[0] == ')' {
					break
				} else {
					p.lexer.Rollback(tok)
				}
				elem, err = p.Value()
				if err != nil {
					return
				}
				if elemtype.Prim() == types.PUndefined {
					elemtype, err = p.TypeOf(elem)
					if err != nil {
						return
					}
				}
				elems = append(elems, elem)
			}
			if elemtype.Prim() == types.PUndefined {
				err = p.selfError(begin, "array literal must have a type or at least one element")
				return
			}
			n = &nodes.ArrayNode{
				Type:     elemtype,
				Elements: elems,
				Pos:      begin.Where,
			}
		default:
			err = p.selfError(tok, "unexpected punctuator")
		}
	default:
		err = p.selfError(tok, "expected value")
	}

	return
}

func (p *Parser) ToNode(v *values.Value, pos lexer.Position) nodes.Node {
	switch v.Type.Prim() {
	case types.PBool:
		return &nodes.BooleanNode{v.Data.(bool), pos}
	case types.PChar:
		return &nodes.CharNode{v.Data.(byte), pos}
	case types.PFloat:
		return &nodes.FloatNode{v.Data.(float32), pos}
	case types.PInt:
		return &nodes.IntegerNode{v.Data.(int32), pos}
	case types.PString:
		return &nodes.StringNode{v.Data.(string), pos}
	}
	panic(v.Type.String())
}
