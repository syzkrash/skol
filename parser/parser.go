package parser

import (
	"strconv"
	"strings"

	"github.com/syzkrash/skol/lexer"
)

type Parser struct {
	lexer *lexer.Lexer
	scope *Scope
}

func (p *Parser) selfError(where *lexer.Token, msg string) error {
	return &ParserError{
		msg:   msg,
		cause: nil,
		Where: where,
	}
}

func (p *Parser) otherError(where *lexer.Token, msg string, cause error) error {
	return &ParserError{
		msg:   msg + ": " + cause.Error(),
		cause: cause,
		Where: where,
	}
}

func (p *Parser) funcCall(f *FuncDefNode) (n Node, err error) {
	panic("unimplemented")
}

// value parses a node that has a value
//
// Example values:
//
//	123        // IntegerNode
//	45.67      // FloatNode
//	"hello"    // StringNode
//	'E'        // CharNode
//	add 1 2    // FuncCallNode
//	age        // VarRefNode
//
func (p *Parser) value() (n Node, err error) {
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

			n = &FloatNode{
				Float: float32(f),
			}
		} else {
			var i int64
			i, err = strconv.ParseInt(tok.Raw, 0, 32)
			if err != nil {
				err = p.otherError(tok, "invalid integer constant", err)
			}

			n = &IntegerNode{
				Int: int32(i),
			}
		}
	case lexer.TkString:
		n = &StringNode{
			Str: tok.Raw,
		}
	case lexer.TkChar:
		rdr := strings.NewReader(tok.Raw)
		r, _, _ := rdr.ReadRune()
		n = &CharNode{
			Char: r,
		}
	case lexer.TkIdent:
		if _, ok := p.scope.FindVar(tok.Raw); ok {
			n = &VarRefNode{
				Var: tok.Raw,
			}
		} else if f, ok := p.scope.FindFunc(tok.Raw); ok {
			n, err = p.funcCall(f)
		} else {
			err = p.selfError(tok, "unknown identifier")
		}
	default:
		err = p.selfError(tok, "expected value")
	}

	return
}

// varDef parses a variable definition node (VarDefNode)
//
// Example variable definition:
//
//	%i: 123
//	%f	:45.67
//	%s: "hello"
//	%	r	:	'E'
//
func (p *Parser) varDef() (n Node, err error) {
	nameToken, err := p.lexer.Next()
	if err != nil {
		return
	}
	if nameToken.Kind != lexer.TkIdent {
		err = p.selfError(nameToken, "expected an identifier")
		return
	}

	sept, err := p.lexer.Next()
	if err != nil {
		return nil, err
	}
	if sept.Kind != lexer.TkPunct {
		err = p.selfError(sept, "expected a punctuator")
		return
	}
	if sept.Raw != ":" {
		err = p.selfError(sept, "expected ':'")
		return
	}

	val, err := p.value()
	if err != nil {
		return nil, err
	}

	n = &VarDefNode{
		Var:   nameToken.Raw,
		Value: val,
	}

	return
}

func (p *Parser) funcDef() (n Node, err error) {
	panic("unimplemented")
}

func (p *Parser) internalNext(tok *lexer.Token) (n Node, err error) {
	switch tok.Kind {
	case lexer.TkPunct:
		if tok.Raw == "$" {
			return p.funcDef()
		}
		if tok.Raw == "%" {
			return p.varDef()
		}
		err = p.selfError(tok, "unexpected punctuator")
	default:
		err = p.selfError(tok, "unexpected token")
	}

	return
}

func (p *Parser) Next() (n Node, err error) {
	tok, err := p.lexer.Next()
	if err != nil {
		return nil, err
	}
	n, err = p.internalNext(tok)
	return
}
