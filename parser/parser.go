package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/syzkrash/skol/lexer"
)

type Parser struct {
	lexer *lexer.Lexer
	Scope *Scope
}

func NewParser(fn string, src io.RuneScanner) *Parser {
	return &Parser{
		lexer: lexer.NewLexer(src, fn),
		Scope: &Scope{
			Parent: nil,
			Funcs:  make(map[string]*Function),
			Vars:   make(map[string]*VarDefNode),
		},
	}
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

// funcCall parses a function call, reading values until enough values for this
// function are found
//
// Example function calls:
//
//	add! 1 2               // add(1, 2)
//	add! sqr! 2 sqr! 2     // add(sqr(2), sqr(2))
//	add! a b               // add(a, b)
//	add! mul! a a mul! bb  // add(mul(a, a), mul(b, b))
func (p *Parser) funcCall(f *Function) (n Node, err error) {
	args := make([]Node, len(f.Args))
	for i := 0; i < len(args); i++ {
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		args[i] = v
	}
	n = &FuncCallNode{
		Func: f.Name,
		Args: args,
	}
	return
}

// value parses a node that has a value
//
// Example values:
//
//	123        // IntegerNode
//	45.67      // FloatNode
//	"hello"    // StringNode
//	'E'        // CharNode
//	add! 1 2   // FuncCallNode
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
		if tok.Raw[len(tok.Raw)-1] == '!' {
			f, ok := p.Scope.FindFunc(tok.Raw[:len(tok.Raw)-1])
			if !ok {
				err = p.selfError(tok, "unknown function")
				return
			}
			n, err = p.funcCall(f)
		} else if _, ok := p.Scope.FindVar(tok.Raw); ok {
			n = &VarRefNode{
				Var: tok.Raw,
			}
		} else {
			err = p.selfError(tok, "unknown variable")
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

	vt, _ := p.TypeOf(val)
	n = &VarDefNode{
		Var:     nameToken.Raw,
		Value:   val,
		VarType: vt,
	}
	p.Scope.Vars[nameToken.Raw] = n.(*VarDefNode)

	return
}

func (p *Parser) funcDef(name string, args map[string]ValueType) (n Node, err error) {
	newScope := &Scope{
		Parent: p.Scope,
		Funcs:  make(map[string]*Function),
		Vars:   make(map[string]*VarDefNode),
	}
	for n, t := range args {
		newScope.Vars[n] = &VarDefNode{
			VarType: t,
			Var:     n,
		}
	}
	p.Scope = newScope
	retType := VtNothing
	body := []Node{}
	for {
		tok, err := p.lexer.Next()
		if err != nil {
			return nil, err
		}
		if tok.Kind == lexer.TkPunct && tok.Raw[0] == ')' {
			break
		}
		n, err := p.internalNext(tok)
		if err != nil {
			return nil, err
		}
		if n.Kind() == NdReturn {
			if retType != VtNothing {
				return nil, &ParserError{
					Where: tok,
					msg:   "duplicate return",
				}
			}
			retType, _ = p.TypeOf(n)
		}
		body = append(body, n)
	}
	p.Scope = p.Scope.Parent

	n = &FuncDefNode{
		Name: name,
		Args: args,
		Body: body,
		Ret:  retType,
	}
	p.Scope.Funcs[name] = DefinedFunction(n.(*FuncDefNode))
	return
}

func (p *Parser) funcOrExtern() (n Node, err error) {
	nameToken, err := p.lexer.Next()
	if err != nil {
		return
	}
	if nameToken.Kind != lexer.TkIdent {
		err = p.selfError(nameToken, "expected an identifier")
		return
	}

	args := map[string]ValueType{}
	for {
		argName, err := p.lexer.Next()
		if err != nil {
			return nil, err
		}
		if argName.Kind == lexer.TkPunct && argName.Raw[0] == '(' {
			return p.funcDef(nameToken.Raw, args)
		}
		if argName.Kind == lexer.TkPunct && argName.Raw[0] == '?' {
			n = &FuncExternNode{
				Name: nameToken.Raw,
				Args: args,
				Ret:  VtNothing,
			}
			p.Scope.Funcs[nameToken.Raw] = ExternFunction(n.(*FuncExternNode))
			return n, nil
		}
		if argName.Kind != lexer.TkIdent {
			return nil, p.selfError(argName, "expected an identifier")
		}
		sept, err := p.lexer.Next()
		if err != nil {
			return nil, err
		}
		if sept.Kind != lexer.TkPunct {
			return nil, p.selfError(sept, "expected a punctuator")
		}
		if sept.Raw[0] != '/' {
			return nil, p.selfError(sept, "expected '/'")
		}
		argType, err := p.lexer.Next()
		if err != nil {
			return nil, err
		}
		if argType.Kind != lexer.TkIdent {
			return nil, p.selfError(argType, "expected an identifier")
		}
		t, ok := ParseType(argType.Raw)
		if !ok {
			return nil, p.selfError(argType, "unknown type")
		}
		args[argName.Raw] = t
	}
}

func (p *Parser) ret() (n Node, err error) {
	v, err := p.value()
	if err != nil {
		return
	}
	n = &ReturnNode{
		Value: v,
	}
	return
}

func (p *Parser) condition() (n Node, err error) {
	condition, err := p.value()
	if err != nil {
		return
	}
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
	newScope := &Scope{
		Parent: p.Scope,
		Funcs:  make(map[string]*Function),
		Vars:   make(map[string]*VarDefNode),
	}
	p.Scope = newScope
	ifBlock := []Node{}
	elseIfNodes := []*IfSubNode{}
	subBlock := []Node{}
	elseBlock := []Node{}
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
		ifBlock = append(ifBlock, n)
	}

	for {
		tok, err = p.lexer.Next()
		if errors.Is(err, io.EOF) {
			err = nil
			goto finish
		}
		if err != nil {
			return
		}
		if tok.Kind != lexer.TkPunct {
			p.lexer.Rollback(tok)
			goto finish
		}
		if tok.Raw != ":" {
			err = p.selfError(tok, "expected '(', got '"+tok.Raw+"'")
			return
		}

		tok, err = p.lexer.Next()
		if errors.Is(err, io.EOF) {
			err = nil
			goto finish
		}
		if err != nil {
			return
		}

		if tok.Kind != lexer.TkPunct {
			err = p.selfError(tok, "expected Punctuator, got "+tok.Kind.String())
			return
		}
		if tok.Raw == "(" {
			p.lexer.Rollback(tok)
			break
		}
		if tok.Raw != "?" {
			err = p.selfError(tok, "expected '(' or '?', got '"+tok.Raw+"'")
			return
		}

		var subCond Node
		subCond, err = p.value()
		if err != nil {
			return
		}

		tok, err = p.lexer.Next()
		if errors.Is(err, io.EOF) {
			err = nil
			goto finish
		}
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
				return
			}
			if tok.Kind == lexer.TkPunct && tok.Raw[0] == ')' {
				elseIfNodes = append(elseIfNodes, &IfSubNode{
					Condition: subCond,
					Block:     subBlock,
				})
				subBlock = []Node{}
				continue
			}
			n, err = p.internalNext(tok)
			if err != nil {
				return
			}
			subBlock = append(subBlock, n)
		}
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
		elseBlock = append(elseBlock, n)
	}

finish:
	p.Scope = p.Scope.Parent

	n = &IfNode{
		Condition:   condition,
		IfBlock:     ifBlock,
		ElseIfNodes: elseIfNodes,
		ElseBlock:   elseBlock,
	}
	return
}

func (p *Parser) internalNext(tok *lexer.Token) (n Node, err error) {
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
		if p.Scope.Parent != nil && tok.Raw[0] == '>' {
			return p.ret()
		}
		err = p.selfError(tok, "unexpected punctuator: "+tok.Raw)
	case lexer.TkIdent:
		if tok.Raw[len(tok.Raw)-1] == '!' {
			f, ok := p.Scope.FindFunc(tok.Raw[:len(tok.Raw)-1])
			if !ok {
				err = p.selfError(tok, "unknown function: "+tok.Raw)
				return
			}
			n, err = p.funcCall(f)
		} else {
			err = p.selfError(tok, "unexpected identifier: "+tok.Raw)
		}
	default:
		err = p.selfError(tok, "unexpected token: "+tok.Raw)
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
