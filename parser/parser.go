package parser

import (
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/sim"
)

type Parser struct {
	lexer *lexer.Lexer
	Sim   *sim.Simulator
	Scope *Scope
}

func NewParser(fn string, src io.RuneScanner) *Parser {
	return &Parser{
		lexer: lexer.NewLexer(src, fn),
		Sim:   sim.NewSimulator(),
		Scope: NewScope(nil),
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
func (p *Parser) funcCall(f *Function) (n nodes.Node, err error) {
	args := make([]nodes.Node, len(f.Args))
	for i := 0; i < len(args); i++ {
		v, err := p.value()
		if err != nil {
			return nil, err
		}
		args[i] = v
	}
	n = &nodes.FuncCallNode{
		Func: f.Name,
		Args: args,
	}
	return
}

// value parses a nodes.Node that has a value
//
// Example values:
//
//	123        // nodes.IntegerNode
//	45.67      // nodes.FloatNode
//	"hello"    // nodes.StringNode
//	'E'        // nodes.CharNode
//	add! 1 2   // nodes.FuncCallNode
//	age        // nodes.VarRefNode
//
func (p *Parser) value() (n nodes.Node, err error) {
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

			n = &nodes.FloatNode{
				Float: float32(f),
			}
		} else {
			var i int64
			i, err = strconv.ParseInt(tok.Raw, 0, 32)
			if err != nil {
				err = p.otherError(tok, "invalid integer constant", err)
			}

			n = &nodes.IntegerNode{
				Int: int32(i),
			}
		}
	case lexer.TkString:
		n = &nodes.StringNode{
			Str: tok.Raw,
		}
	case lexer.TkChar:
		rdr := strings.NewReader(tok.Raw)
		r, _, _ := rdr.ReadRune()
		n = &nodes.CharNode{
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
			n = &nodes.VarRefNode{
				Var: tok.Raw,
			}
		} else if v, ok := p.Scope.FindConst(tok.Raw); ok {
			n = p.ToNode(v)
		} else {
			err = p.selfError(tok, "unknown variable")
		}
	case lexer.TkPunct:
		if tok.Raw == "*" {
			n = &nodes.BooleanNode{true}
		} else if tok.Raw == "/" {
			n = &nodes.BooleanNode{false}
		} else {
			err = p.selfError(tok, "unexpected punctuator")
		}
	default:
		err = p.selfError(tok, "expected value")
	}

	return
}

// varDef parses a variable definition nodes.Node (nodes.VarDefNode)
//
// Example variable definition:
//
//	%i: 123
//	%f	:45.67
//	%s: "hello"
//	%	r	:	'E'
//
func (p *Parser) varDef() (n nodes.Node, err error) {
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
	n = &nodes.VarDefNode{
		Var:     nameToken.Raw,
		Value:   val,
		VarType: vt,
	}
	p.Scope.SetVar(nameToken.Raw, n.(*nodes.VarDefNode))

	return
}

func (p *Parser) funcOrExtern() (n nodes.Node, err error) {
	nameToken, err := p.lexer.Next()
	if err != nil {
		return
	}
	if nameToken.Kind != lexer.TkIdent {
		err = p.selfError(nameToken, "expected an identifier")
		return
	}

	var typeToken *lexer.Token
	var funcType values.ValueType
	var ok bool
	sepToken, err := p.lexer.Next()
	if err != nil {
		return
	}
	if sepToken.Kind == lexer.TkPunct && sepToken.Raw == "/" {
		typeToken, err = p.lexer.Next()
		if err != nil {
			return
		}
		if typeToken.Kind != lexer.TkIdent {
			err = p.selfError(typeToken, "expected Identifier, got "+typeToken.Kind.String())
			return
		}
		funcType, ok = values.ParseType(typeToken.Raw)
		if !ok {
			err = p.selfError(typeToken, "unknown type: "+funcType.String())
			return
		}
	} else {
		p.lexer.Rollback(sepToken)
	}

	var body []nodes.Node
	var argName *lexer.Token
	args := map[string]values.ValueType{}
	for {
		argName, err = p.lexer.Next()
		if err != nil {
			return nil, err
		}
		if argName.Kind == lexer.TkPunct && argName.Raw[0] == '(' {
			newScope := &Scope{
				Parent: p.Scope,
				Funcs:  make(map[string]*Function),
				Vars:   make(map[string]*nodes.VarDefNode),
			}
			for n, t := range args {
				newScope.SetVar(n, &nodes.VarDefNode{
					VarType: t,
					Var:     n,
				})
			}
			p.Scope = newScope

			p.lexer.Rollback(argName)
			body, err = p.block()
			if err != nil {
				return
			}

			n = &nodes.FuncDefNode{
				Name: nameToken.Raw,
				Args: args,
				Ret:  funcType,
				Body: body,
			}
			p.Scope.SetFunc(nameToken.Raw, DefinedFunction(n.(*nodes.FuncDefNode)))
			return
		}
		if argName.Kind == lexer.TkPunct && argName.Raw[0] == '?' {
			var internToken *lexer.Token
			internToken, err = p.lexer.Next()
			if err != nil {
				return
			}
			intern := ""
			if internToken.Kind != lexer.TkString {
				p.lexer.Rollback(internToken)
			} else {
				intern = internToken.Raw
			}
			n = &nodes.FuncExternNode{
				Name:   nameToken.Raw,
				Intern: intern,
				Args:   args,
				Ret:    funcType,
			}
			p.Scope.SetFunc(nameToken.Raw, ExternFunction(n.(*nodes.FuncExternNode)))
			return
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
		t, ok := values.ParseType(argType.Raw)
		if !ok {
			return nil, p.selfError(argType, "unknown type: "+argType.Raw)
		}
		args[argName.Raw] = t
	}
}

func (p *Parser) ret() (n nodes.Node, err error) {
	v, err := p.value()
	if err != nil {
		return
	}
	n = &nodes.ReturnNode{
		Value: v,
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

func (p *Parser) condition() (n nodes.Node, err error) {
	condition, err := p.value()
	if err != nil {
		return
	}
	newScope := &Scope{
		Parent: p.Scope,
		Funcs:  make(map[string]*Function),
		Vars:   make(map[string]*nodes.VarDefNode),
	}
	p.Scope = newScope
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

	p.Scope = p.Scope.Parent

	n = &nodes.IfNode{
		Condition:   condition,
		IfBlock:     ifb,
		ElseIfNodes: elifn,
		ElseBlock:   elseb,
	}
	return
}

func (p *Parser) loop() (n nodes.Node, err error) {
	condition, err := p.value()
	if err != nil {
		return
	}

	newScope := &Scope{
		Parent: p.Scope,
		Funcs:  map[string]*Function{},
		Vars:   map[string]*nodes.VarDefNode{},
	}
	oldScope := p.Scope
	p.Scope = newScope

	body, err := p.block()
	if err != nil {
		return
	}

	p.Scope = oldScope

	n = &nodes.WhileNode{
		Condition: condition,
		Body:      body,
	}
	return
}

func (p *Parser) constant() (err error) {
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
		return err
	}
	if sept.Kind != lexer.TkPunct {
		err = p.selfError(sept, "expected a punctuator")
		return
	}
	if sept.Raw != ":" {
		err = p.selfError(sept, "expected ':'")
		return
	}

	n, err := p.value()
	if err != nil {
		return err
	}
	v, err := p.Sim.Const(n)
	if err != nil {
		return err
	}

	if !p.Scope.SetConst(nameToken.Raw, v) {
		return p.selfError(nameToken, "constant value cannot be redefined")
	}
	return nil
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
		if p.Scope.Parent != nil && tok.Raw[0] == '>' {
			return p.ret()
		}
		err = p.selfError(tok, "unexpected top-level punctuator: "+tok.Raw)
	case lexer.TkIdent:
		if tok.Raw[len(tok.Raw)-1] == '!' {
			f, ok := p.Scope.FindFunc(tok.Raw[:len(tok.Raw)-1])
			if !ok {
				err = p.selfError(tok, "unknown function: "+tok.Raw)
				return
			}
			n, err = p.funcCall(f)
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
	return
}

func (p *Parser) TypeOf(n nodes.Node) (t values.ValueType, ok bool) {
	switch n.Kind() {
	case nodes.NdInteger:
		t = values.VtInteger
	case nodes.NdFloat:
		t = values.VtFloat
	case nodes.NdChar:
		t = values.VtChar
	case nodes.NdString:
		t = values.VtString
	case nodes.NdReturn:
		t, ok = p.TypeOf(n.(*nodes.ReturnNode).Value)
	case nodes.NdVarRef:
		var v *nodes.VarDefNode
		v, ok = p.Scope.FindVar(n.(*nodes.VarRefNode).Var)
		if !ok {
			return
		}
		t = v.VarType
	case nodes.NdFuncCall:
		var f *Function
		f, ok = p.Scope.FindFunc(n.(*nodes.FuncCallNode).Func)
		if !ok {
			return
		}
		t = f.Ret
	default:
		ok = false
	}
	return
}

func (p *Parser) ToNode(v *values.Value) nodes.Node {
	switch v.ValueType {
	case values.VtBool:
		return &nodes.BooleanNode{v.Bool}
	case values.VtChar:
		return &nodes.CharNode{v.Char}
	case values.VtFloat:
		return &nodes.FloatNode{v.Float}
	case values.VtInteger:
		return &nodes.IntegerNode{v.Int}
	case values.VtString:
		return &nodes.StringNode{v.Str}
	}
	panic(v.ValueType.String())
}
