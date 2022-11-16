package parser

import (
	"errors"
	"io"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

// parseVar parses a variable definition, assignment or both.
//
// Variable defintion:
//
//	%name/string
//
// Variable assignment:
//
//	%name: "Joe"
//
// Both:
//
//	%name/string: "Joe"
func (p *Parser) parseVar() (n ast.Node, err error) {
	var (
		name string

		typed bool
		vtype types.Type

		valued bool
		value  ast.MetaNode

		tok *lexer.Token
	)

	tok, err = p.lexer.Next()
	if err != nil {
		return
	}

	if tok.Kind != lexer.TIdent {
		err = tokErr(pe.EExpectedName, tok)
		return
	}
	name = tok.Raw

	tok, err = p.lexer.Next()
	if err != nil {
		return
	}

	if pn, ok := tok.Punct(); ok && pn == lexer.PType {
		typed = true
		vtype, err = p.parseType()
		if err != nil {
			return
		}
		tok, err = p.lexer.Next()
		if errors.Is(err, io.EOF) {
			err = nil
			goto final
		}
		if err != nil {
			return
		}
	}
	if pn, ok := tok.Punct(); ok && pn == lexer.PIs {
		valued = true
		value, err = p.ParseValue()
		if err != nil {
			return
		}
	} else if typed {
		p.lexer.Rollback(tok)
	}

final:
	if !typed && !valued {
		err = tokErr(pe.ENeedTypeOrValue, tok)
		return
	}

	if typed && !valued {
		n = ast.VarDefNode{
			Var:  name,
			Type: vtype,
		}
		p.Scope.SetVar(name, nil)
	} else if !typed && valued {
		n = ast.VarSetNode{
			Var:   name,
			Value: value,
		}
		p.Scope.SetVar(name, value.Node)
	} else {
		n = ast.VarSetTypedNode{
			Var:   name,
			Type:  vtype,
			Value: value,
		}
		p.Scope.SetVar(name, value.Node)
	}
	return
}

// parseConst parses a constant definition.
//
//	#name: "Joe"
func (p *Parser) parseConst() (err error) {
	nameToken, err := p.lexer.Next()
	if err != nil {
		return
	}
	if nameToken.Kind != lexer.TIdent {
		err = tokErr(pe.EExpectedName, nameToken)
		return
	}

	sept, err := p.lexer.Next()
	if err != nil {
		return err
	}
	if pn, ok := sept.Punct(); !ok || pn != lexer.PIs {
		err = tokErr(pe.EExpectedColon, sept)
		return
	}

	n, err := p.ParseValue()
	if err != nil {
		return err
	}

	if !p.Scope.SetConst(nameToken.Raw, n.Node) {
		return tokErr(pe.EConstantRedefined, nameToken)
	}
	return nil
}

// parseFunc parses a function definition or an extern definition.
//
// Function definition:
//
//	$SayHello Name/string(print! concat! "Hello, " Name)
//
// Extern definition:
//
//	$Exit Status/int?
//
// Aliased extern definition:
//
//	$OSExit Status/int?"exit"
//	 ^                  ^
//	 |                  |
//	alias/skol name     actual name
//
// Function shorthand:
//
//	$Add1/int n/int: add_i! n 1
func (p *Parser) parseFunc() (n ast.Node, err error) {
	var (
		name          string
		ret           types.Type
		args          []types.Descriptor
		body          ast.Block
		shorthandBody ast.MetaNode

		argName string
		argType types.Type

		tok *lexer.Token
	)

	tok, err = p.lexer.Next()
	if err != nil {
		return
	}
	if tok.Kind != lexer.TIdent {
		err = tokErr(pe.EExpectedName, tok)
		return
	}

	name = tok.Raw

	tok, err = p.lexer.Next()
	if err != nil {
		return
	}
	if pn, ok := tok.Punct(); ok && pn == lexer.PType {
		ret, err = p.parseType()
		if err != nil {
			return
		}
	} else {
		p.lexer.Rollback(tok)
	}

	for {
		tok, err = p.lexer.Next()

		if pn, ok := tok.Punct(); ok {
			switch pn {
			case lexer.PIf:
				n = ast.FuncExternNode{
					Alias: name,
					Proto: args,
					Ret:   ret,
					Name:  name,
				}
				return
			case lexer.PLParen:
				p.lexer.Rollback(tok)
				p.Scope = &Scope{
					Parent: p.Scope,
					Vars:   make(map[string]ast.Node),
					Consts: make(map[string]ast.Node),
					Types:  make(map[string]types.Type),
				}
				for _, a := range args {
					var falseVal, ok = p.NodeOf(a.Type)
					if !ok {
						err = tokErr(pe.EBadFuncArgType, tok)
						return
					}
					p.Scope.Vars[a.Name] = falseVal
				}
				body, err = p.parseBlock()
				if err != nil {
					return
				}
				p.Scope = p.Scope.Parent
				n = ast.FuncDefNode{
					Name:  name,
					Proto: args,
					Ret:   ret,
					Body:  body,
				}
				return
			case lexer.PIs:
				p.Scope = &Scope{
					Parent: p.Scope,
					Vars:   make(map[string]ast.Node),
					Consts: make(map[string]ast.Node),
					Types:  make(map[string]types.Type),
				}
				for _, a := range args {
					var falseVal, ok = p.NodeOf(a.Type)
					if !ok {
						err = tokErr(pe.EBadFuncArgType, tok)
						return
					}
					p.Scope.Vars[a.Name] = falseVal
				}
				tok, err = p.lexer.Next()
				if err != nil {
					return
				}
				// ensure that @ always means structure instatiation inside of a
				// shorthand (struct definitions are not allowed in functions anyways)
				if pn, ok := tok.Punct(); ok && pn == lexer.PStruct {
					p.lexer.Rollback(tok)
					shorthandBody, err = p.ParseValue()
				} else {
					shorthandBody, _, err = p.next(tok)
				}
				if err != nil {
					return
				}
				p.Scope = p.Scope.Parent
				n = ast.FuncShorthandNode{
					Name:  name,
					Proto: args,
					Ret:   ret,
					Body:  shorthandBody,
				}
				return
			}
		}
		if tok.Kind != lexer.TIdent {
			err = tokErr(pe.ENeedBodyOrExtern, tok)
			return
		}

		argName = tok.Raw

		tok, err = p.lexer.Next()
		if err != nil {
			return
		}

		if pn, ok := tok.Punct(); !ok || pn != lexer.PType {
			err = tokErr(pe.EExpectedType, tok)
			return
		}

		argType, err = p.parseType()
		if err != nil {
			return
		}

		args = append(args, types.Descriptor{Name: argName, Type: argType})
	}
}

// parseStruct parses a structure type definition.
//
//	@Vec2i(x/int y/int)
func (p *Parser) parseStruct() (n ast.Node, err error) {
	var (
		name      string
		fieldName string
		fieldType types.Type
		fields    []types.Descriptor

		tok *lexer.Token
	)

	tok, err = p.lexer.Next()
	if err != nil {
		return
	}

	if tok.Kind != lexer.TIdent {
		err = tokErr(pe.EExpectedName, tok)
		return
	}
	name = tok.Raw

	tok, err = p.lexer.Next()
	if err != nil {
		return nil, err
	}

	if pn, ok := tok.Punct(); !ok || pn != lexer.PLParen {
		err = tokErr(pe.EExpectedLParen, tok)
		return
	}

	for {
		tok, err = p.lexer.Next()
		if err != nil {
			return
		}

		if pn, ok := tok.Punct(); ok && pn == lexer.PRParen {
			break
		}
		if tok.Kind != lexer.TIdent {
			err = tokErr(pe.EExpectedName, tok)
			return
		}
		fieldName = tok.Raw

		tok, err = p.lexer.Next()
		if err != nil {
			return
		}

		if pn, ok := tok.Punct(); !ok || pn != lexer.PType {
			err = tokErr(pe.EExpectedType, tok)
			return
		}

		fieldType, err = p.parseType()
		if err != nil {
			return
		}

		fields = append(fields, types.Descriptor{Name: fieldName, Type: fieldType})
	}

	p.Scope.Types[name] = types.StructType{
		Name:   name,
		Fields: fields,
	}
	n = ast.StructDefNode{
		Name:   name,
		Fields: fields,
	}
	return
}
