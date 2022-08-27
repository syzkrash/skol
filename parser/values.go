package parser

import (
	"strconv"
	"strings"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/lexer"
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
func (p *Parser) Value() (mn ast.MetaNode, err error) {
	tok, err := p.lexer.Next()
	if err != nil {
		return
	}

	mn.Where = tok.Where

	switch tok.Kind {
	case lexer.TkConstant:
		if strings.ContainsRune(tok.Raw, '.') {
			var f float64
			f, err = strconv.ParseFloat(tok.Raw, 32)
			if err != nil {
				err = p.otherError(tok, "invalid floating-point constant", err)
			}

			mn.Node = ast.FloatNode{
				Value: float32(f),
			}
		} else {
			var i int64
			i, err = strconv.ParseInt(tok.Raw, 0, 32)
			if err != nil {
				err = p.otherError(tok, "invalid integer constant", err)
			}

			mn.Node = ast.IntNode{
				Value: int32(i),
			}
		}
	case lexer.TkString:
		mn.Node = ast.StringNode{
			Value: tok.Raw,
		}
	case lexer.TkChar:
		mn.Node = ast.CharNode{
			Value: tok.Raw[0],
		}
	case lexer.TkIdent:
		if tok.Raw[len(tok.Raw)-1] == '!' {
			fn := tok.Raw[:len(tok.Raw)-1]
			f, ok := p.Scope.FindFunc(fn)
			if !ok {
				err = p.selfError(tok, "unknown function: "+fn)
				return
			}
			mn.Node, err = p.funcCall(fn, f, tok.Where)
		} else if _, ok := p.Scope.FindVar(tok.Raw); ok {
			mn.Node, err = p.selector(tok)
		} else if v, ok := p.Scope.FindConst(tok.Raw); ok {
			mn.Node = v
		} else {
			err = p.selfError(tok, "unknown variable: "+tok.Raw)
		}
	case lexer.TkPunct:
		switch tok.Raw[0] {
		case '*':
			mn.Node = ast.BoolNode{
				Value: true,
			}
		case '/':
			mn.Node = ast.BoolNode{
				Value: false,
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
			args := make([]ast.MetaNode, len(s.Fields))
			for i := range s.Fields {
				mn, err = p.Value()
				if err != nil {
					return
				}
				args[i] = mn
			}
			mn.Node = ast.StructNode{
				Type: s,
				Args: args,
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
			elems := []ast.MetaNode{}
			var elem ast.MetaNode
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
					elemtype, err = p.TypeOf(elem.Node)
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
			mn.Node = ast.ArrayNode{
				Type:  types.ArrayType{Element: elemtype},
				Elems: elems,
			}
		default:
			err = p.selfError(tok, "unexpected punctuator")
		}
	default:
		err = p.selfError(tok, "expected value")
	}

	return
}
