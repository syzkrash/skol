package parser

import (
	"fmt"

	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
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

	val, err := p.Value()
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
	v, err := p.Sim.Const(n)
	if err != nil {
		return err
	}

	if !p.Scope.SetConst(nameToken.Raw, v) {
		return p.selfError(nameToken, "constant value cannot be redefined")
	}
	return nil
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

	var funcType types.Type = types.Undefined
	sepToken, err := p.lexer.Next()
	if err != nil {
		return
	}
	if sepToken.Kind == lexer.TkPunct && sepToken.Raw == "/" {
		funcType, err = p.parseType()
		if err != nil {
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
		t, err := p.parseType()
		if err != nil {
			return nil, err
		}
		args = append(args, values.FuncArg{Name: argName.Raw, Type: t})
	}
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
		fType   types.Type
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
		fType, err = p.parseType()
		if err != nil {
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
