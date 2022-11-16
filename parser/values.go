package parser

import (
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
	case lexer.TInt:
		i, ok := tok.Int()
		if !ok {
			err = tokErr(pe.EBadIntLit, tok)
			return
		}

		mn.Node = ast.IntNode{
			Value: i,
		}
	case lexer.TFloat:
		f, ok := tok.Float()
		if !ok {
			err = tokErr(pe.EBadFloatLit, tok)
			return
		}

		mn.Node = ast.FloatNode{
			Value: f,
		}
	case lexer.TString:
		mn.Node = ast.StringNode{
			Value: tok.Raw,
		}
	case lexer.TChar:
		mn.Node = ast.CharNode{
			Value: tok.Raw[0],
		}
	case lexer.TIdent:
		if tok.Raw[len(tok.Raw)-1] == '!' {
			fn := tok.Raw[:len(tok.Raw)-1]
			var argc int
			f, ok := p.Tree.Funcs[fn]
			if !ok {
				bf, ok := builtins[fn]
				if !ok {
					err = tokErr(pe.EUnknownFunction, tok)
					return
				}
				argc = bf.ArgCount
			} else {
				argc = len(f.Args)
			}
			mn.Node, err = p.parseCall(fn, argc, tok.Where)
		} else if _, ok := p.Scope.FindVar(tok.Raw); ok {
			mn.Node, err = p.parseSelector(tok)
		} else if v, ok := p.Scope.FindConst(tok.Raw); ok {
			mn.Node = v
		} else {
			err = tokErr(pe.EUnknownVariable, tok)
		}
	case lexer.TPunct:
		pn, _ := tok.Punct()
		switch pn {
		case lexer.PLoop:
			mn.Node = ast.BoolNode{
				Value: true,
			}
		case lexer.PType:
			mn.Node = ast.BoolNode{
				Value: false,
			}
		case lexer.PStruct:
			tok, err = p.lexer.Next()
			if err != nil {
				return
			}
			if tok.Kind != lexer.TIdent {
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
		case lexer.PLBrack:
			begin := tok
			var elemtype types.Type = types.Undefined
			tok, err = p.lexer.Next()
			if err != nil {
				return
			}
			if tok.Kind == lexer.TIdent {
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
			if pn, ok := tok.Punct(); !ok || pn != lexer.PRBrack {
				err = tokErr(pe.EExpectedType, tok)
				return
			}
			tok, err = p.lexer.Next()
			if err != nil {
				return
			}
			if pn, ok := tok.Punct(); !ok || pn != lexer.PLParen {
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
				if pn, ok := tok.Punct(); ok && pn == lexer.PRParen {
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
