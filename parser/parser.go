package parser

import (
	"errors"
	"io"
	"strconv"

	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
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
		v, err := p.Value()
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

func (p *Parser) selector(start *lexer.Token) (n nodes.Node, err error) {
	n = &nodes.SelectorNode{
		Parent: nil,
		Child:  start.Raw,
		Pos:    start.Where,
	}
	var tok *lexer.Token
	for {
		// first, consume the #
		tok, err = p.lexer.Next()
		// the selector *could* be the last thing in a file, so we just return
		// on EOF
		if errors.Is(err, io.EOF) {
			err = nil
			return
		}
		if err != nil {
			return
		}

		// rollback the token we read in case it isn't a #
		// we have just consumed an element of the selector, so if we don't have
		// another # that means that is the end of the selector
		if tok.Kind != lexer.TkPunct || tok.Raw[0] != '#' {
			p.lexer.Rollback(tok)
			return
		}

		// now, we consume the actual selector element
		tok, err = p.lexer.Next()
		if err != nil {
			return
		}
		switch tok.Kind {
		// ident: select a field on a structure
		case lexer.TkIdent:
			// determine if the parent selector is an array
			var pt types.Type
			pt, err = p.TypeOf(n)
			if err != nil {
				return
			}
			// append to the chain of selectors
			if pt.Prim() == types.PArray {
				n = &nodes.IndexNode{
					Parent: n.(nodes.Selector),
					Idx: &nodes.SelectorNode{
						Parent: nil,
						Child:  tok.Raw,
						Pos:    tok.Where,
					},
					Pos: n.Where(),
				}
			} else {
				n = &nodes.SelectorNode{
					Parent: n.(nodes.Selector),
					Child:  tok.Raw,
					Pos:    n.Where(),
				}
			}
		// constant: index into an array
		//	indexes are always unsigned integers, but base prefixes are allowed
		case lexer.TkConstant:
			// parse the index, this will error out if the index is not an unsigned
			// integer
			var idx uint64
			idx, err = strconv.ParseUint(tok.Raw, 0, 32)
			if err != nil {
				return
			}
			// append to the chain
			n = &nodes.IndexNode{
				Parent: n.(nodes.Selector),
				Idx: &nodes.IntegerNode{
					Int: int32(idx),
					Pos: n.Where(),
				},
				Pos: n.Where(),
			}
		// punct: can be a typecast
		//	typecasts use the @ punctuator
		case lexer.TkPunct:
			if tok.Raw[0] != '@' {
				// error out if the punctuator is not @
				err = p.selfError(tok, "expected identifer, constant or '@'")
				return
			}
			// get the type for typecast
			// this also allows arrays to be typecast (makes sense if you think about
			// it)
			var t types.Type
			t, err = p.parseType()
			if err != nil {
				return
			}
			// append to the chain
			n = &nodes.TypecastNode{
				Parent: n.(nodes.Selector),
				Type:   t,
				Pos:    n.Where(),
			}
		// any other token is not allowed
		default:
			err = p.selfError(tok, "expected identifer, constant or '@'")
			return
		}
	}
}

func (p *Parser) ret() (n nodes.Node, err error) {
	v, err := p.Value()
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
