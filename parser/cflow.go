package parser

import (
	"errors"
	"io"

	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/nodes"
)

func (p *Parser) condition() (n nodes.Node, err error) {
	condition, err := p.value()
	if err != nil {
		return
	}
	debug.Log(debug.AttrScope, "Entering new scope")
	p.Scope = NewScope(p.Scope)
	ifb, err := p.block()
	if err != nil {
		return
	}

	var (
		elseb []nodes.Node

		elifn    []*nodes.IfSubNode
		subcond  nodes.Node
		subblock []nodes.Node

		tok *lexer.Token
	)
	for {
		tok, err = p.lexer.Next()
		if errors.Is(err, io.EOF) {
			err = nil
			break
		}
		if err != nil {
			return
		}
		if tok.Kind != lexer.TkPunct || tok.Raw != ":" {
			p.lexer.Rollback(tok)
			break
		}
		tok, err = p.lexer.Next()
		if err != nil {
			return
		}
		if tok.Kind == lexer.TkPunct && tok.Raw == "?" {
			subcond, err = p.value()
			if err != nil {
				return
			}
			subblock, err = p.block()
			if err != nil {
				return nil, err
			}
			elifn = append(elifn, &nodes.IfSubNode{
				Condition: subcond,
				Block:     subblock,
			})
		} else {
			p.lexer.Rollback(tok)
			elseb, err = p.block()
			if err != nil {
				return
			}
		}
	}

	debug.Log(debug.AttrScope, "Exitig scope")
	p.Scope = p.Scope.Parent

	n = &nodes.IfNode{
		Condition:   condition,
		IfBlock:     ifb,
		ElseIfNodes: elifn,
		ElseBlock:   elseb,
		Pos:         condition.Where(),
	}
	return
}

func (p *Parser) loop() (n nodes.Node, err error) {
	condition, err := p.value()
	if err != nil {
		return
	}

	debug.Log(debug.AttrScope, "Entering new scope")
	p.Scope = NewScope(p.Scope)
	body, err := p.block()
	if err != nil {
		return
	}

	debug.Log(debug.AttrScope, "Exiting scope")
	p.Scope = p.Scope.Parent
	n = &nodes.WhileNode{
		Condition: condition,
		Body:      body,
		Pos:       condition.Where(),
	}
	return
}
