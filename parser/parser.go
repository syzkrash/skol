package parser

import (
	"errors"
	"io"

	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/sim"
)

type Parser struct {
	lexer  *lexer.Lexer
	Engine string
	Sim    *sim.Simulator
	Scope  *Scope
}

func NewParser(fn string, src io.RuneScanner, eng string) *Parser {
	return &Parser{
		lexer:  lexer.NewLexer(src, fn),
		Engine: eng,
		Sim:    sim.NewSimulator(),
		Scope:  NewScope(nil),
	}
}

// funcCall parses a function call, reading values until enough values for this
// function are found
//
// Example function calls:
//
//	add! 1 2               // add(1, 2)
//	add! sqr! 2 sqr! 2     // add(sqr(2), sqr(2))
//	add! a b               // add(a, b)
//	add! mul! a a mul! bb  // add(mul(a, a), mul(b, b))
func (p *Parser) funcCall(fn string, f *values.Function, pos lexer.Position) (n nodes.Node, err error) {
	args := make([]nodes.Node, len(f.Args))
	for i := 0; i < len(args); i++ {
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		args[i] = v
	}
	n = &nodes.FuncCallNode{
		Func: fn,
		Args: args,
		Pos:  pos,
	}
	return
}

func (p *Parser) selectorOrTypecast(start *lexer.Token) (n nodes.Node, err error) {
	n = &nodes.SelectorNode{
		Parent: nil,
		Child:  start.Raw,
		Pos:    start.Where,
	}
	var tok *lexer.Token
	for {
		tok, err = p.lexer.Next()
		if errors.Is(err, io.EOF) {
			err = nil
			return
		}
		if err != nil {
			return
		}
		if tok.Kind != lexer.TkPunct || tok.Raw != "#" {
			p.lexer.Rollback(tok)
			return
		}
		tok, err = p.lexer.Next()
		if err != nil {
			return
		}
		if tok.Kind == lexer.TkPunct && tok.Raw == "@" {
			typecastLoc := tok.Where
			tok, err = p.lexer.Next()
			if err != nil {
				return
			}
			if tok.Kind != lexer.TkIdent {
				err = p.selfError(tok, "expected an identifer")
				return
			}
			t, ok := p.ParseType(tok.Raw)
			if !ok {
				err = p.selfError(tok, "unknown type: "+t.String())
				return
			}
			n = &nodes.TypecastNode{
				Value:  n.(*nodes.SelectorNode),
				Target: t,
				Pos:    typecastLoc,
			}
			return
		}
		if tok.Kind != lexer.TkIdent {
			err = p.selfError(tok, "expected an identifer")
			return
		}
		n = &nodes.SelectorNode{
			Parent: n.(*nodes.SelectorNode),
			Child:  tok.Raw,
			Pos:    n.Where(),
		}
	}
}

func (p *Parser) ret() (n nodes.Node, err error) {
	v, err := p.value()
	if err != nil {
		return
	}
	n = &nodes.ReturnNode{
		Value: v,
		Pos:   v.Where(),
	}
	return
}

func (p *Parser) block() (ns []nodes.Node, err error) {
	var n nodes.Node
	var tok *lexer.Token
	tok, err = p.lexer.Next()
	if err != nil {
		return
	}
	if tok.Kind != lexer.TkPunct {
		err = p.selfError(tok, "expected Punctuator, got "+tok.Kind.String())
		return
	}
	if tok.Raw != "(" {
		err = p.selfError(tok, "expected '(', got '"+tok.Raw+"'")
		return
	}
	for {
		tok, err = p.lexer.Next()
		if err != nil {
			break
		}
		if tok.Kind == lexer.TkPunct && tok.Raw == ")" {
			break
		}
		n, err = p.internalNext(tok)
		if err != nil {
			break
		}
		ns = append(ns, n)
	}
	return
}

func (p *Parser) internalNext(tok *lexer.Token) (n nodes.Node, err error) {
	switch tok.Kind {
	case lexer.TkPunct:
		if tok.Raw == "$" {
			return p.funcOrExtern()
		}
		if tok.Raw == "%" {
			return p.varDef()
		}
		if tok.Raw == "?" {
			return p.condition()
		}
		if tok.Raw == "*" {
			return p.loop()
		}
		if tok.Raw == "#" {
			err = p.constant()
			if err != nil {
				return
			}
			return p.Next()
		}
		if tok.Raw == "@" {
			return p.structn()
		}
		if p.Scope.Parent != nil && tok.Raw[0] == '>' {
			return p.ret()
		}
		err = p.selfError(tok, "unexpected top-level punctuator: "+tok.Raw)
	case lexer.TkIdent:
		if tok.Raw[len(tok.Raw)-1] == '!' {
			fnm := tok.Raw[:len(tok.Raw)-1]
			f, ok := p.Scope.FindFunc(fnm)
			if !ok {
				err = p.selfError(tok, "unknown function: "+fnm)
				return
			}
			n, err = p.funcCall(fnm, f, tok.Where)
		} else {
			err = p.selfError(tok, "unexpected top-level identifier: "+tok.Raw)
		}
	default:
		err = p.selfError(tok, "unexpected top-level token: "+tok.Raw)
	}

	return
}

func (p *Parser) Next() (n nodes.Node, err error) {
	tok, err := p.lexer.Next()
	if err != nil {
		return nil, err
	}
	n, err = p.internalNext(tok)
	if err != nil {
		debug.Log(debug.AttrParser, "Error %s", err)
	} else {
		debug.Log(debug.AttrParser, "%s node at %s", n.Kind(), n.Where())
	}
	return
}
