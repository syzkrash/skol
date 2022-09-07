package parser

import (
	"errors"
	"io"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values"
	"github.com/syzkrash/skol/parser/values/types"
)

// parseVar parses a variable definition nodes.Node (nodes.VarDefNode)
//
// Example variable definition:
//
//	%i: 123
//	%f	:45.67
//	%s: "hello"
//	%	r	:	'E'
//
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
		if errors.Is(err, io.EOF) {
			err = nil
			goto final
		}
		if err != nil {
			return
		}
	}
	if tok.Kind == lexer.TkPunct && tok.Raw[0] == ':' {
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
		err = p.selfError(tok, "variable must have an explicit type, an explicit value or both")
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

func (p *Parser) parseConst() (err error) {
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

	n, err := p.ParseValue()
	if err != nil {
		return err
	}

	if !p.Scope.SetConst(nameToken.Raw, n.Node) {
		return p.selfError(nameToken, "constant value cannot be redefined")
	}
	return nil
}

func (p *Parser) parseFunc() (n ast.Node, err error) {
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
	if err != nil {
		return
	}
	if tok.Kind != lexer.TkIdent {
		err = p.selfError(tok, "expected identifier for function name")
		return
	}

	name = tok.Raw

	tok, err = p.lexer.Next()
	if err != nil {
		return
	}
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
			fargs := make([]values.FuncArg, len(args))
			for i, a := range args {
				fargs[i] = values.FuncArg{
					Name: a.Name,
					Type: a.Type,
				}
			}
			p.Scope.Funcs[name] = &values.Function{
				Name: name,
				Args: fargs,
				Ret:  ret,
			}
			return
		}
		if tok.Kind == lexer.TkPunct && tok.Raw[0] == '(' {
			p.lexer.Rollback(tok)
			p.Scope = &Scope{
				Parent: p.Scope,
				Funcs:  make(map[string]*values.Function),
				Vars:   make(map[string]ast.Node),
				Consts: make(map[string]ast.Node),
				Types:  make(map[string]types.Type),
			}
			for _, a := range args {
				var falseVal, ok = p.NodeOf(a.Type)
				if !ok {
					err = p.selfError(tok, "bad type for function argument! "+a.Type.String())
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
			fargs := make([]values.FuncArg, len(args))
			for i, a := range args {
				fargs[i] = values.FuncArg{
					Name: a.Name,
					Type: a.Type,
				}
			}
			p.Scope.Funcs[name] = &values.Function{
				Name: name,
				Args: fargs,
				Ret:  ret,
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
			return
		}

		argType, err = p.parseType()
		if err != nil {
			return
		}

		args = append(args, ast.FuncProtoArg{Name: argName, Type: argType})
	}
}

func (p *Parser) parseStruct() (n ast.Node, err error) {
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
		err = p.selfError(tok, "expeted structure name")
		return
	}
	name = tok.Raw

	tok, err = p.lexer.Next()
	if err != nil {
		return nil, err
	}

	if tok.Kind != lexer.TkPunct || tok.Raw[0] != '(' {
		err = p.selfError(tok, "expected '(' for structure field definitions")
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
		if tok.Kind != lexer.TkIdent {
			err = p.selfError(tok, "expected structure field name")
			return
		}
		fieldName = tok.Raw

		tok, err = p.lexer.Next()
		if err != nil {
			return
		}

		if tok.Kind != lexer.TkPunct || tok.Raw[0] != '/' {
			err = p.selfError(tok, "expected '/' for structure field type")
			return
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
