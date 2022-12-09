package parser

import (
	"errors"
	"io"
	"strconv"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

// Parser consumes tokens from its internal lexer and constructs nodes out of
// them.
type Parser struct {
	lexer  *lexer.Lexer
	errs   chan error
	Tree   ast.AST
	Engine string
	Scope  *Scope
}

// NewParser creates a new parser for the given engine, creating a [Lexer] with
// the given input stream.
func NewParser(fn string, src io.RuneScanner, eng string, errOut chan error) *Parser {
	return &Parser{
		lexer:  lexer.NewLexer(src, fn),
		errs:   errOut,
		Tree:   ast.NewAST(),
		Engine: eng,
		Scope:  NewScope(nil),
	}
}

// Parse constructs nodes using the internal lexer's tokens and compiles them
// into an [ast.AST].
func (p *Parser) Parse() ast.AST {
	var (
		n    ast.MetaNode
		skip bool
	)

	p.Tree = ast.AST{
		Vars:     make(map[string]ast.Var),
		Typedefs: make(map[string]ast.Typedef),
		Funcs:    make(map[string]ast.Func),
		Exerns:   make(map[string]ast.Extern),
		Structs:  make(map[string]ast.Structure),
	}

	for {
		tok, err := p.lexer.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			p.errs <- err
			continue
		}

		n, skip, err = p.next(tok)
		if skip {
			continue
		}
		if err != nil {
			debug.Log(debug.AttrParser, "Error %s", err)
		} else {
			debug.Log(debug.AttrParser, "%s node at %s", n.Node.Kind(), n.Where)
		}
		if err != nil {
			p.errs <- err
			continue
		}

		switch n.Node.Kind() {
		case ast.NVarSet:
			nvs := n.Node.(ast.VarSetNode)
			p.Tree.Vars[nvs.Var] = ast.Var{
				Name:  nvs.Var,
				Value: nvs.Value,
				Node:  n,
			}
			delete(p.Tree.Typedefs, nvs.Var)
		case ast.NVarDef:
			nvd := n.Node.(ast.VarDefNode)
			p.Tree.Typedefs[nvd.Var] = ast.Typedef{
				Name: nvd.Var,
				Type: nvd.Type,
				Node: n,
			}
		case ast.NVarSetTyped:
			nvst := n.Node.(ast.VarSetTypedNode)
			p.Tree.Vars[nvst.Var] = ast.Var{
				Name:  nvst.Var,
				Value: nvst.Value,
				Node:  n,
			}
		case ast.NFuncDef:
			nfd := n.Node.(ast.FuncDefNode)
			p.Tree.Funcs[nfd.Name] = ast.Func{
				Name: nfd.Name,
				Args: nfd.Proto,
				Ret:  nfd.Ret,
				Body: nfd.Body,
				Node: n,
			}
			delete(p.Tree.Exerns, nfd.Name)
		case ast.NFuncShorthand:
			nfs := n.Node.(ast.FuncShorthandNode)
			body := ast.Block{{Where: nfs.Body.Where}}
			if nfs.Body.Node.Kind().IsValue() {
				body[0].Node = ast.ReturnNode{Value: nfs.Body}
			} else {
				body[0].Node = nfs.Body.Node
			}
			p.Tree.Funcs[nfs.Name] = ast.Func{
				Name: nfs.Name,
				Args: nfs.Proto,
				Ret:  nfs.Ret,
				Body: body,
				Node: n,
			}
			delete(p.Tree.Exerns, nfs.Name)
		case ast.NFuncExtern:
			nfe := n.Node.(ast.FuncExternNode)
			p.Tree.Exerns[nfe.Alias] = ast.Extern{
				Name:  nfe.Name,
				Alias: nfe.Alias,
				Ret:   nfe.Ret,
				Args:  nfe.Proto,
				Node:  n,
			}
		case ast.NStructDef:
			nsd := n.Node.(ast.StructDefNode)
			p.Tree.Structs[nsd.Name] = ast.Structure{
				Name:   nsd.Name,
				Fields: nsd.Fields,
				Node:   n,
			}
		default:
			p.errs <- nodeErr(pe.EIllegalTopLevelNode, n)
			continue
		}
	}

	return p.Tree
}

// TopLevel parses a top-level statement. One of:
//   - Function/Extern definition
//   - Variable defintion and/or assignment
//   - Structure type definition
func (p *Parser) TopLevel() (mn ast.MetaNode) {
	tok, err := p.lexer.Next()
	if err != nil {
		p.errs <- err
		return
	}

	var skip bool
	for {
		mn, skip, err = p.next(tok)
		if err != nil {
			p.errs <- err
			return
		}
		if !skip {
			break
		}
	}

	return
}

// next constructs whatever node is next. If skip is true, no node was produced
// by the given code and next() needs to be called again to retrieve a node.
func (p *Parser) next(tok *lexer.Token) (mn ast.MetaNode, skip bool, err error) {
	var (
		n ast.Node
	)

	switch tok.Kind {
	case lexer.TPunct:
		pn, _ := tok.Punct()
		switch pn {
		case lexer.PFunc:
			n, err = p.parseFunc()
		case lexer.PVar:
			n, err = p.parseVar()
		case lexer.PIf:
			n, err = p.parseIf()
		case lexer.PLoop:
			n, err = p.parseWhile()
		case lexer.PStruct:
			n, err = p.parseStruct()
		case lexer.PReturn:
			n, err = p.parseReturn()
		case lexer.PField:
			err = p.parseConst()
			if err != nil {
				return
			}
			skip = true
		default:
			err = tokErr(pe.EUnexpectedToken, tok)
		}
	case lexer.TIdent:
		var maybeBang *lexer.Token
		maybeBang, err = p.lexer.Next()
		if err != nil {
			return
		}
		if pn, ok := maybeBang.Punct(); !ok || pn != lexer.PExecute {
			p.lexer.Rollback(maybeBang)
			err = tokErr(pe.EUnexpectedToken, tok)
			return
		}
		fnm := tok.Raw
		var argc int
		f, ok := p.Tree.Funcs[fnm]
		if !ok {
			bf, ok := builtins[fnm]
			if !ok {
				err = tokErr(pe.EUnknownFunction, tok)
				return
			}
			argc = bf.ArgCount
		} else {
			argc = len(f.Args)
		}
		n, err = p.parseCall(fnm, argc, tok.Where)
	default:
		err = tokErr(pe.EUnexpectedToken, tok)
	}

	mn.Node = n
	mn.Where = tok.Where

	return
}

// parseCall parses a function call. This function requires that an argument
// count be known before the call can be parsed.
//
//	print! concat! "Hello " World
func (p *Parser) parseCall(fn string, argc int, pos lexer.Position) (n ast.Node, err error) {
	args := make([]ast.MetaNode, argc)
	for i := 0; i < len(args); i++ {
		v, err := p.ParseValue()
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

// parseSelector parses a series of selector elements and compiles them into
// one selector.
//
// Basic selector:
//
//	Person
//
// Field selector:
//
//	Person#Name
//	Person#Age
//
// Index selector:
//
//	People#0
//	People#[MyID]
//
// Index selector and field selector:
//
//	People#0#Name
//	People#[PersonNo]#Age
//
// Typecast:
//
//	Person#@Employee#Employer
func (p *Parser) parseSelector(start *lexer.Token) (n ast.Node, err error) {
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
		if pn, ok := tok.Punct(); !ok || pn != lexer.PField {
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
		case lexer.TIdent:
			// append to the chain of selectors
			n = ast.SelectorNode{
				Parent: n.(ast.Selector),
				Child:  tok.Raw,
			}
		// constant: index into an array
		//	indexes are always unsigned integers, but base prefixes are allowed
		case lexer.TInt:
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
		// punct: can be a typecast or an array index
		//	typecasts use the @ punctuator and indexes use the [] punctuators
		case lexer.TPunct:
			pn, _ := tok.Punct()
			switch pn {
			case lexer.PStruct:
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
			case lexer.PLBrack:
				// get the token starting the index
				tok, err = p.lexer.Next()
				if err != nil {
					return
				}

				// parse the index itself
				var idx ast.Node
				idx, err = p.parseSelector(tok)
				if err != nil {
					return
				}

				// get the closing bracket
				tok, err = p.lexer.Next()
				if err != nil {
					return
				}
				if pn, ok := tok.Punct(); !ok || pn != lexer.PRBrack {
					err = tokErr(pe.EExpectedRBrack, tok)
					return
				}

				// append to the chain
				n = ast.IndexSelectorNode{
					Parent: n.(ast.Selector),
					Idx:    idx.(ast.Selector),
				}
			default:
				// error out if the punctuator is not @
				err = tokErr(pe.EExpectedSelectorElem, tok)
				return
			}
		// any other token is not allowed
		default:
			err = tokErr(pe.EExpectedSelectorElem, tok)
			return
		}
	}
}

// parseReturn parses a return.
//
//	>"Hello"
func (p *Parser) parseReturn() (n ast.Node, err error) {
	v, err := p.ParseValue()
	if err != nil {
		return
	}
	n = ast.ReturnNode{
		Value: v,
	}
	return
}

// parseBlock parses a block of code.
//
//	(
//		print! "Hello world!"
//	)
func (p *Parser) parseBlock() (block ast.Block, err error) {
	var (
		n    ast.MetaNode
		skip bool
		tok  *lexer.Token
	)

	tok, err = p.lexer.Next()
	if err != nil {
		return
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

		n, skip, err = p.next(tok)
		if err != nil {
			return
		}
		if skip {
			continue
		}

		block = append(block, n)
	}

	return
}

func tokErr(c pe.ErrorCode, cause *lexer.Token) *pe.PrettyError {
	return pe.New(c).Section("Caused by", "%s \"%s\" at %s", cause.Kind, cause.Raw, cause.Where)
}

func nodeErr(c pe.ErrorCode, cause ast.MetaNode) *pe.PrettyError {
	return pe.New(c).Section("Caused by", "`%s` node at %s", cause.Node.Kind(), cause.Where)
}
