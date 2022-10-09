package parser

import (
	"strconv"
	"strings"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

// ParseValue parses a value.
//
// Boolean literal:
//
//   - /
//
// Character literal:
//
//	'h'  'i'
//
// Numeric literal:
//
//	1234  12_34  0x1234 0x12_34
//	12.34  123_456.789
//
// String literal:
//
//	"Hello world"  "hi\nthere"  "how\tare\tyou?"
//
// Structure literal:
//
//	@Vec2i(12 34)
//	@Vec3f(1.23 4.56 7.68)
//
// Array literal:
//
//	[int](0 1 2 3 4 5 6 7 8 9)
//	   [](0.1 2.3 4.5 6.7 8.9)
//	[string]()
//
// Any selector:
//
//	Someone
//	Someone#Age
//	Someone#@Employee#Employer#Age
//	People#0#@Employee#JobHistory#0#Owner#Name
//
// Function call:
//
//	DoSomething!
//	DontDoAnything!
//	Say! MyName
//	add_i! 12 34
func (p *Parser) ParseValue() (mn ast.MetaNode, err error) {
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
				err = tokErr(pe.EBadFloatLit, tok)
				return
			}

			mn.Node = ast.FloatNode{
				Value: float32(f),
			}
		} else {
			var i int64
			i, err = strconv.ParseInt(tok.Raw, 0, 32)
			if err != nil {
				err = tokErr(pe.EBadIntLit, tok)
				return
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
			f, ok := p.Tree.Funcs[fn]
			if !ok {
				err = tokErr(pe.EUnknownFunction, tok)
				return
			}
			mn.Node, err = p.parseCall(fn, f, tok.Where)
		} else if _, ok := p.Scope.FindVar(tok.Raw); ok {
			mn.Node, err = p.parseSelector(tok)
		} else if v, ok := p.Scope.FindConst(tok.Raw); ok {
			mn.Node = v
		} else {
			err = tokErr(pe.EUnknownVariable, tok)
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
				err = tokErr(pe.EExpectedName, tok)
				return
			}
			t, ok := p.Scope.FindType(tok.Raw)
			if !ok {
				err = tokErr(pe.EUnknownType, tok)
				return
			}
			s := t.(types.StructType)
			args := make([]ast.MetaNode, len(s.Fields))
			for i := range s.Fields {
				mn, err = p.ParseValue()
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
					err = tokErr(pe.EUnknownType, tok)
					return
				}
				tok, err = p.lexer.Next()
				if err != nil {
					return
				}
			}
			if tok.Kind != lexer.TkPunct || tok.Raw[0] != ']' {
				err = tokErr(pe.EExpectedType, tok)
				return
			}
			tok, err = p.lexer.Next()
			if err != nil {
				return
			}
			if tok.Kind != lexer.TkPunct || tok.Raw[0] != '(' {
				err = tokErr(pe.EExpectedLParen, tok)
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
				elem, err = p.ParseValue()
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
				err = tokErr(pe.ENeedTypeOrValue, begin)
				return
			}
			mn.Node = ast.ArrayNode{
				Type:  types.ArrayType{Element: elemtype},
				Elems: elems,
			}
		default:
			err = tokErr(pe.EUnexpectedToken, tok)
		}
	default:
		err = tokErr(pe.EExpectedValue, tok)
	}

	return
}
