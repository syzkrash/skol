package parser

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/syzkrash/skol/common"
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
//	@VectorTwo 1.2 3.4 // nodes.NewStructNode
//	pos#x      // nodes.SelectorNode
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
				Pos:   tok.Where,
			}
		} else {
			var i int64
			i, err = strconv.ParseInt(tok.Raw, 0, 32)
			if err != nil {
				err = p.otherError(tok, "invalid integer constant", err)
			}

			n = &nodes.IntegerNode{
				Int: int32(i),
				Pos: tok.Where,
			}
		}
	case lexer.TkString:
		n = &nodes.StringNode{
			Str: tok.Raw,
			Pos: tok.Where,
		}
	case lexer.TkChar:
		rdr := strings.NewReader(tok.Raw)
		r, _, _ := rdr.ReadRune()
		n = &nodes.CharNode{
			Char: byte(r),
			Pos:  tok.Where,
		}
	case lexer.TkIdent:
		if tok.Raw[len(tok.Raw)-1] == '!' {
			fn := tok.Raw[:len(tok.Raw)-1]
			f, ok := p.Scope.FindFunc(fn)
			if !ok {
				err = p.selfError(tok, "unknown function: "+fn)
				return
			}
			n, err = p.funcCall(fn, f, tok.Where)
		} else if _, ok := p.Scope.FindVar(tok.Raw); ok {
			return p.selectorOrTypecast(tok)
		} else if v, ok := p.Scope.FindConst(tok.Raw); ok {
			n = p.ToNode(v, tok.Where)
		} else {
			err = p.selfError(tok, "unknown variable: "+tok.Raw)
		}
	case lexer.TkPunct:
		if tok.Raw == "*" {
			n = &nodes.BooleanNode{true, tok.Where}
		} else if tok.Raw == "/" {
			n = &nodes.BooleanNode{false, tok.Where}
		} else if tok.Raw == "@" {
			tok, err = p.lexer.Next()
			if err != nil {
				return
			}
			if tok.Kind != lexer.TkIdent {
				err = p.selfError(tok, "expected Identifier, got "+tok.Kind.String())
				return
			}
			t, ok := p.Scope.FindType(tok.Raw)
			if !ok {
				err = p.selfError(tok, "unknown type: "+tok.Raw)
				return
			}
			s := t.(types.StructType)
			args := make([]nodes.Node, len(s.Fields))
			for i := range s.Fields {
				n, err = p.value()
				if err != nil {
					return
				}
				args[i] = n
			}
			n = &nodes.NewStructNode{
				Type: t,
				Args: args,
				Pos:  tok.Where,
			}
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

	vt, err := p.TypeOf(val)
	if err != nil {
		err = fmt.Errorf("could not deduce type of %s: %s", val, err)
		return
	}

	if old, ok := p.Scope.FindVar(nameToken.Raw); ok {
		if !old.VarType.Equals(vt) {
			err = p.selfError(nameToken, fmt.Sprintf(
				"variable redefinition with incorrect type: %s (expected %s)",
				vt.String(), old.VarType.String()))
			return
		}
	}

	n = &nodes.VarDefNode{
		Var:     nameToken.Raw,
		Value:   val,
		VarType: vt,
		Pos:     nameToken.Where,
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
	var funcType types.Type = types.Undefined
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
		funcType, ok = p.ParseType(typeToken.Raw)
		if !ok {
			err = p.selfError(typeToken, "unknown type: "+typeToken.Raw)
			return
		}
	} else {
		p.lexer.Rollback(sepToken)
	}

	var body []nodes.Node
	var argName *lexer.Token
	args := []values.FuncArg{}
	for {
		argName, err = p.lexer.Next()
		if err != nil {
			return nil, err
		}
		if argName.Kind == lexer.TkPunct && argName.Raw[0] == '(' {
			newScope := NewScope(p.Scope)
			for _, a := range args {
				newScope.SetVar(a.Name, &nodes.VarDefNode{
					VarType: a.Type,
					Var:     a.Name,
				})
			}
			debug.Log(debug.AttrScope, "Entering new scope")
			p.Scope = newScope

			p.lexer.Rollback(argName)
			body, err = p.block()
			if err != nil {
				return
			}
			for _, bn := range body {
				if bn.Kind() == nodes.NdReturn {
					rn := bn.(*nodes.ReturnNode)
					t, _ := p.TypeOf(rn.Value)
					if funcType == types.Undefined {
						funcType = t
					}
				}
			}
			if funcType.Prim() == types.PUndefined {
				funcType = types.Nothing
			}
			fdn := &nodes.FuncDefNode{
				Name: nameToken.Raw,
				Args: args,
				Ret:  funcType,
				Body: body,
				Pos:  nameToken.Where,
			}
			n = fdn
			p.Scope.Funcs[nameToken.Raw] =
				values.DefinedFunction(fdn.Name, fdn.Args, fdn.Ret)
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
			if funcType.Prim() == types.PUndefined {
				funcType = types.Nothing
			}
			fen := &nodes.FuncExternNode{
				Name:   nameToken.Raw,
				Intern: intern,
				Args:   args,
				Ret:    funcType,
				Pos:    nameToken.Where,
			}
			n = fen
			p.Scope.Funcs[nameToken.Raw] =
				values.ExternFunction(fen.Name, fen.Intern, fen.Args, fen.Ret)
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
		t, ok := p.ParseType(argType.Raw)
		if !ok {
			return nil, p.selfError(argType, "unknown type: "+argType.Raw)
		}
		args = append(args, values.FuncArg{Name: argName.Raw, Type: t})
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

func (p *Parser) structn() (n nodes.Node, err error) {
	nameTk, err := p.lexer.Next()
	if err != nil {
		return
	}
	if nameTk.Kind != lexer.TkIdent {
		err = p.selfError(nameTk, "expected Identifier, got "+nameTk.Kind.String())
		return
	}
	startTk, err := p.lexer.Next()
	if err != nil {
		return
	}
	if startTk.Kind != lexer.TkPunct {
		err = p.selfError(startTk, "expected Punctuator, got "+startTk.Kind.String())
		return
	}
	if startTk.Raw != "(" {
		err = p.selfError(startTk, "expected '(', got '"+startTk.Raw+"'")
		return
	}
	fields := []types.Field{}
	var (
		fNameTk *lexer.Token
		sepTk   *lexer.Token
		typeTk  *lexer.Token
		fType   types.Type
		ok      bool
	)
	for {
		fNameTk, err = p.lexer.Next()
		if err != nil {
			return
		}
		if fNameTk.Kind == lexer.TkPunct && fNameTk.Raw == ")" {
			break
		}
		if fNameTk.Kind != lexer.TkIdent {
			err = p.selfError(fNameTk, "expected Identifier, got "+fNameTk.Kind.String())
			return
		}
		sepTk, err = p.lexer.Next()
		if err != nil {
			return
		}
		if sepTk.Kind != lexer.TkPunct {
			err = p.selfError(sepTk, "expected Punctuator, got "+sepTk.Kind.String())
			return
		}
		if sepTk.Raw != "/" {
			err = p.selfError(sepTk, "expected '/', got '"+sepTk.Raw+"'")
			return
		}
		typeTk, err = p.lexer.Next()
		if err != nil {
			return
		}
		if typeTk.Kind != lexer.TkIdent {
			err = p.selfError(typeTk, "expected Identifier, got "+typeTk.Kind.String())
			return
		}
		fType, ok = p.ParseType(typeTk.Raw)
		if !ok {
			err = p.selfError(typeTk, "unknown type: "+typeTk.Raw)
			return
		}
		fields = append(fields, types.Field{fNameTk.Raw, fType})
	}
	t := types.StructType{nameTk.Raw, fields}
	n = &nodes.StructNode{
		Name: nameTk.Raw,
		Type: t,
		Pos:  nameTk.Where,
	}
	p.Scope.Types[nameTk.Raw] = t
	return
}

func (p *Parser) ParseType(raw string) (types.Type, bool) {
	switch strings.ToLower(raw) {
	case "integer", "int32", "int", "i32", "i":
		return types.Int, true
	case "boolean", "bool", "b":
		return types.Bool, true
	case "float32", "float", "f32", "f":
		return types.Float, true
	case "char", "ch", "c":
		return types.Char, true
	case "string", "str", "s":
		return types.String, true
	case "any", "a":
		return types.Any, true
	}
	if stype, ok := p.Scope.FindType(raw); ok {
		return stype, true
	}
	return nil, false
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

func (p *Parser) TypeOf(n nodes.Node) (t types.Type, err error) {
	switch n.Kind() {
	case nodes.NdBoolean:
		t = types.Bool
	case nodes.NdInteger:
		t = types.Int
	case nodes.NdFloat:
		t = types.Float
	case nodes.NdChar:
		t = types.Char
	case nodes.NdString:
		t = types.String
	case nodes.NdNewStruct:
		t = n.(*nodes.NewStructNode).Type
	case nodes.NdFuncCall:
		fn := n.(*nodes.FuncCallNode).Func
		f, ok := p.Scope.FindFunc(fn)
		if !ok {
			err = fmt.Errorf("unknown function: %s", fn)
			return
		}
		t = f.Ret
	case nodes.NdSelector:
		s := n.(*nodes.SelectorNode)
		path := s.Path()
		v, ok := p.Scope.FindVar(path[0])
		if !ok {
			err = fmt.Errorf("unknown variable: %s", path[0])
			return
		}
		t = v.VarType
		if len(path) == 1 {
			return
		}
	outer:
		for _, e := range path[1:] {
			if t.Prim() != types.PStruct {
				err = common.Error(n, "can only select fields on structures")
				return
			}
			for _, f := range t.(types.StructType).Fields {
				if f.Name == e {
					t = f.Type
					continue outer
				}
			}
			err = common.Error(n, "%s does not contain field '%s'", t.String(), e)
			return
		}
	case nodes.NdTypecast:
		return n.(*nodes.TypecastNode).Target, nil
	default:
		err = fmt.Errorf("%s node is not a value", n.Kind())
	}
	return
}

func (p *Parser) ToNode(v *values.Value, pos lexer.Position) nodes.Node {
	switch v.Type.Prim() {
	case types.PBool:
		return &nodes.BooleanNode{v.Data.(bool), pos}
	case types.PChar:
		return &nodes.CharNode{v.Data.(byte), pos}
	case types.PFloat:
		return &nodes.FloatNode{v.Data.(float32), pos}
	case types.PInt:
		return &nodes.IntegerNode{v.Data.(int32), pos}
	case types.PString:
		return &nodes.StringNode{v.Data.(string), pos}
	}
	panic(v.Type.String())
}
