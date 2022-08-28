package parser

import (
	"errors"
	"io"
	"strconv"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/lexer"
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
func (p *Parser) funcCall(fn string, f *values.Function, pos lexer.Position) (n ast.Node, err error) {
	args := make([]ast.MetaNode, len(f.Args))
	for i := 0; i < len(args); i++ {
		v, err := p.Value()
		if err != nil {
			return nil, err
		}
		args[i] = v
	}
	n = ast.FuncCallNode{
		Func: fn,
		Args: args,
	}
	return
}

func (p *Parser) selector(start *lexer.Token) (n ast.Node, err error) {
	n = ast.SelectorNode{
		Parent: nil,
		Child:  start.Raw,
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
				n = ast.IndexSelectorNode{
					Parent: n.(ast.Selector),
					Idx: ast.SelectorNode{
						Parent: nil,
						Child:  tok.Raw,
					},
				}
			} else {
				n = ast.SelectorNode{
					Parent: n.(ast.Selector),
					Child:  tok.Raw,
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
			n = ast.IndexConstNode{
				Parent: n.(ast.Selector),
				Idx:    int(idx),
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
			n = ast.TypecastNode{
				Parent: n.(ast.Selector),
				Cast:   t,
			}
		// any other token is not allowed
		default:
			err = p.selfError(tok, "expected identifer, constant or '@'")
			return
		}
	}
}

func (p *Parser) ret() (n ast.Node, err error) {
	v, err := p.Value()
	if err != nil {
		return
	}
	n = ast.ReturnNode{
		Value: v,
	}
	return
}

func (p *Parser) block() (block ast.Block, err error) {
	var (
		n ast.MetaNode

		tok *lexer.Token
	)

	tok, err = p.lexer.Next()
	if err != nil {
		return
	}

	if tok.Kind != lexer.TkPunct || tok.Raw[0] != '(' {
		err = p.selfError(tok, "expected '(' to start block")
		return
	}

	for {
		tok, err = p.lexer.Next()
		if err != nil {
			return
		}

		if tok.Kind == lexer.TkPunct && tok.Raw[0] == ')' {
			break
		}

		n, err = p.internalNext(tok)
		if err != nil {
			return
		}

		block = append(block, n)
	}

	return
}

func (p *Parser) internalNext(tok *lexer.Token) (mn ast.MetaNode, err error) {
	var (
		n ast.Node
	)

	switch tok.Kind {
	case lexer.TkPunct:
		switch tok.Raw[0] {
		case '$':
			n, err = p.funcOrExtern()
		case '%':
			n, err = p.varDef()
		case '?':
			n, err = p.condition()
		case '*':
			n, err = p.loop()
		case '@':
			n, err = p.structn()
		case '>':
			n, err = p.ret()
		case '#':
			err = p.constant()
			if err != nil {
				return
			}
			return p.Next()
		default:
			err = p.selfError(tok, "unexpected top-level punctuator: "+tok.Raw)
		}
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

	mn.Node = n
	mn.Where = tok.Where

	return
}

func (p *Parser) Next() (n ast.MetaNode, err error) {
	tok, err := p.lexer.Next()
	if err != nil {
		return ast.MetaNode{}, err
	}
	n, err = p.internalNext(tok)
	if err != nil {
		debug.Log(debug.AttrParser, "Error %s", err)
	} else {
		debug.Log(debug.AttrParser, "%s node at %s", n.Node.Kind(), n.Where)
	}
	return
}
