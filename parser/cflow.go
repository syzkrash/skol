package parser

import (
	"errors"
	"io"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/lexer"
)

// parseIf parses an if condition and all of it's branches.
//
// With just a condition:
//
//	?Condition!(Action!)
//
// With an else branch:
//
//	?Condition!(Action!):(OtherAction!)
//
// With 1 else-if branch and an else branch:
//
//	?Condition!(Action!):?OtherCondition!(OtherAction!):(AnotherAction!)
func (p *Parser) parseIf() (n ast.Node, err error) {
	var (
		cond  ast.MetaNode
		block ast.Block

		out ast.IfNode

		tok *lexer.Token
	)

	cond, err = p.ParseValue()
	if err != nil {
		return
	}
	debug.Log(debug.AttrScope, "Entering new scope")
	p.Scope = NewScope(p.Scope)
	block, err = p.parseBlock()
	if err != nil {
		return
	}

	out.Main.Cond = cond
	out.Main.Block = block

	for {
		tok, err = p.lexer.Next()
		if errors.Is(err, io.EOF) {
			err = nil
			break
		}
		if err != nil {
			return
		}
		if pn, ok := tok.Punct(); !ok || pn != lexer.PIs {
			p.lexer.Rollback(tok)
			break
		}
		tok, err = p.lexer.Next()
		if err != nil {
			return
		}
		if pn, ok := tok.Punct(); ok && pn == lexer.PIf {
			cond, err = p.ParseValue()
			if err != nil {
				return
			}
			block, err = p.parseBlock()
			if err != nil {
				return nil, err
			}
			out.Other = append(out.Other, ast.Branch{
				Cond:  cond,
				Block: block,
			})
		} else {
			p.lexer.Rollback(tok)
			out.Else, err = p.parseBlock()
			if err != nil {
				return
			}
		}
	}

	debug.Log(debug.AttrScope, "Exitig scope")
	p.Scope = p.Scope.Parent
	n = out

	return
}

// parseWhile parses a while loop.
//
//	*Condition!(Action!)
func (p *Parser) parseWhile() (n ast.Node, err error) {
	cond, err := p.ParseValue()
	if err != nil {
		return
	}

	debug.Log(debug.AttrScope, "Entering new scope")
	p.Scope = NewScope(p.Scope)
	block, err := p.parseBlock()
	if err != nil {
		return
	}

	debug.Log(debug.AttrScope, "Exiting scope")
	p.Scope = p.Scope.Parent
	n = ast.WhileNode{
		Cond:  cond,
		Block: block,
	}
	return
}
