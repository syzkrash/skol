package parser

import (
	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

// varDef parses a variable definition nodes.Node (nodes.VarDefNode)
//
// Example variable definition:
//
//	%i: 123
//	%f	:45.67
//	%s: "hello"
//	%	r	:	'E'
//
func (p *Parser) varDef() (n ast.Node, err error) {
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

	if tok.Kind != lexer.TkIdent {
		err = p.selfError(tok, "expected variable name")
		return
	}
	name = tok.Raw

	tok, err = p.lexer.Next()
	if err != nil {
		return
	}

	if tok.Kind == lexer.TkPunct && tok.Raw[0] == '/' {
		typed = true
		vtype, err = p.parseType()
		if err != nil {
			return
		}
		tok, err = p.lexer.Next()
		if err != nil {
			return
		}
	}
	if tok.Kind == lexer.TkPunct && tok.Raw[0] == ':' {
		valued = true
		value, err = p.Value()
		if err != nil {
			return
		}
	}

	if !typed && !valued {
		err = p.selfError(tok, "variable must have an explicit type, an explicit value or both")
		return
	}

	if typed && !valued {
		n = ast.VarDefNode{
			Var:  name,
			Type: vtype,
		}
	} else if !typed && valued {
		n = ast.VarSetNode{
			Var:   name,
			Value: value,
		}
	} else {
		n = ast.VarSetTypedNode{
			Var:   name,
			Type:  vtype,
			Value: value,
		}
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

	n, err := p.Value()
	if err != nil {
		return err
	}

	if !p.Scope.SetConst(nameToken.Raw, n.Node) {
		return p.selfError(nameToken, "constant value cannot be redefined")
	}
	return nil
}

func (p *Parser) funcOrExtern() (n ast.Node, err error) {
	var (
		name string
		ret  types.Type
		args []ast.FuncProtoArg
		body ast.Block

		argName string
		argType types.Type

		tok *lexer.Token
	)

	tok, err = p.lexer.Next()
	if tok.Kind != lexer.TkIdent {
		err = p.selfError(tok, "expected identifier for function name")
		return
	}

	name = tok.Raw

	tok, err = p.lexer.Next()
	if tok.Kind == lexer.TkPunct && tok.Raw[0] == '/' {
		ret, err = p.parseType()
		if err != nil {
			return
		}
	} else {
		p.lexer.Rollback(tok)
	}

	for {
		tok, err = p.lexer.Next()

		if tok.Kind == lexer.TkPunct && tok.Raw[0] == '?' {
			n = ast.FuncExternNode{
				Alias: name,
				Proto: args,
				Ret:   ret,
				Name:  name,
			}
			return
		}
		if tok.Kind == lexer.TkPunct && tok.Raw[0] == '(' {
			p.lexer.Rollback(tok)
			body, err = p.block()
			if err != nil {
				return
			}
			n = ast.FuncDefNode{
				Name:  name,
				Proto: args,
				Ret:   ret,
				Body:  body,
			}
			return
		}
		if tok.Kind != lexer.TkIdent {
			err = p.selfError(tok, "expected function argument, body or '?' for extern")
			return
		}

		argName = tok.Raw

		tok, err = p.lexer.Next()
		if err != nil {
			return
		}

		if tok.Kind != lexer.TkPunct || tok.Raw[0] != '/' {
			err = p.selfError(tok, "expected '/' for function argument type")
		}

		argType, err = p.parseType()
		if err != nil {
			return
		}

		args = append(args, ast.FuncProtoArg{Name: argName, Type: argType})
	}
}

func (p *Parser) structn() (n ast.Node, err error) {
	var (
		name      string
		fieldName string
		fieldType types.Type
		fields    []ast.StructProtoField

		typeFields []types.Field

		tok *lexer.Token
	)

	tok, err = p.lexer.Next()
	if err != nil {
		return
	}

	if tok.Kind != lexer.TkIdent {
		err = p.selfError(tok, "expected structure name")
	}
	name = tok.Raw

	tok, err = p.lexer.Next()
	if err != nil {
		return nil, err
	}

	if tok.Kind != lexer.TkPunct || tok.Raw[0] != '(' {
		err = p.selfError(tok, "expected '(' for structure field definitions")
	}

	for {
		tok, err = p.lexer.Next()
		if err != nil {
			return
		}

		if tok.Kind == lexer.TkPunct && tok.Raw[0] == ')' {
			break
		}
		if tok.Kind != lexer.TkIdent {
			err = p.selfError(tok, "expected structure field name")
		}
		fieldName = tok.Raw

		tok, err = p.lexer.Next()
		if err != nil {
			return
		}

		if tok.Kind != lexer.TkPunct || tok.Raw[0] != '/' {
			err = p.selfError(tok, "expected '/' for structure field type")
		}

		fieldType, err = p.parseType()
		if err != nil {
			return
		}

		fields = append(fields, ast.StructProtoField{Name: fieldName, Type: fieldType})
	}

	for _, pf := range fields {
		typeFields = append(typeFields, types.Field{Name: pf.Name, Type: pf.Type})
	}

	p.Scope.Types[name] = types.StructType{
		Name:   name,
		Fields: typeFields,
	}
	n = ast.StructDefNode{
		Name:   name,
		Fields: fields,
	}
	return
}
