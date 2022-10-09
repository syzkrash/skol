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
	Tree   ast.AST
	Engine string
	Scope  *Scope
}

// NewParser creates a new parser for the given engine, creating a [Lexer] with
// the given input stream.
func NewParser(fn string, src io.RuneScanner, eng string) *Parser {
	return &Parser{
		lexer:  lexer.NewLexer(src, fn),
		Tree:   ast.NewAST(),
		Engine: eng,
		Scope:  NewScope(nil),
	}
}

func (p *Parser) getFunc(fn string) (f ast.Func, ok bool) {
	if f, ok = p.Tree.Funcs[fn]; ok {
		return
	}
	f, ok = defaultFunctions[fn]
	return
}

// parseCall parser a function call.
//
//	print! concat! "Hello " World
func (p *Parser) parseCall(fn string, f ast.Func, pos lexer.Position) (n ast.Node, err error) {
	args := make([]ast.MetaNode, len(f.Args))
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
// Index selector and field selector:
//
//	People#0#Name
//	People#PersonNo#Age
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
		if tok.Kind != lexer.TkPunct || tok.Raw[0] != '#' {
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
		case lexer.TkIdent:
			// determine if the parent selector is an array
			var pt types.Type
			pt, err = p.TypeOf(n)
			if err != nil {
				return
			}
			// append to the chain of selectors
			if pt.Prim() == types.PArray {
				n = ast.IndexSelectorNode{
					Parent: n.(ast.Selector),
					Idx: ast.SelectorNode{
						Parent: nil,
						Child:  tok.Raw,
					},
				}
			} else {
				n = ast.SelectorNode{
					Parent: n.(ast.Selector),
					Child:  tok.Raw,
				}
			}
		// constant: index into an array
		//	indexes are always unsigned integers, but base prefixes are allowed
		case lexer.TkConstant:
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
		// punct: can be a typecast
		//	typecasts use the @ punctuator
		case lexer.TkPunct:
			if tok.Raw[0] != '@' {
				// error out if the punctuator is not @
				err = tokErr(pe.EExpectedSelectorElem, tok)
				return
			}
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

	if tok.Kind != lexer.TkPunct || tok.Raw[0] != '(' {
		err = tokErr(pe.EExpectedLParen, tok)
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

// next constructs whatever node is next. If skip is true, no node was produced
// by the given code and next() needs to be called again to retrieve a node.
func (p *Parser) next(tok *lexer.Token) (mn ast.MetaNode, skip bool, err error) {
	var (
		n ast.Node
	)

	switch tok.Kind {
	case lexer.TkPunct:
		switch tok.Raw[0] {
		case '$':
			n, err = p.parseFunc()
		case '%':
			n, err = p.parseVar()
		case '?':
			n, err = p.parseIf()
		case '*':
			n, err = p.parseWhile()
		case '@':
			n, err = p.parseStruct()
		case '>':
			n, err = p.parseReturn()
		case '#':
			err = p.parseConst()
			if err != nil {
				return
			}
			skip = true
		default:
			err = tokErr(pe.EUnexpectedToken, tok)
		}
	case lexer.TkIdent:
		if tok.Raw[len(tok.Raw)-1] == '!' {
			fnm := tok.Raw[:len(tok.Raw)-1]
			f, ok := p.getFunc(fnm)
			if !ok {
				err = tokErr(pe.EUnknownFunction, tok)
				return
			}
			n, err = p.parseCall(fnm, f, tok.Where)
		} else {
			err = tokErr(pe.EUnexpectedToken, tok)
		}
	default:
		err = tokErr(pe.EUnexpectedToken, tok)
	}

	mn.Node = n
	mn.Where = tok.Where

	return
}

// TopLevel parses a top-level statement. One of:
//   - Function/Extern definition
//   - Variable defintion and/or assignment
//   - Structure type definition
func (p *Parser) TopLevel() (mn ast.MetaNode, err error) {
	tok, err := p.lexer.Next()
	if err != nil {
		return
	}

	var skip bool
	for {
		mn, skip, err = p.next(tok)
		if err != nil {
			return
		}
		if !skip {
			break
		}
	}

	return
}

// Parse constructs nodes using the internal lexer's tokens and compiles them
// into an [ast.AST].
func (p *Parser) Parse() (tree ast.AST, err error) {
	var (
		tok  *lexer.Token
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
		tok, err = p.lexer.Next()
		if errors.Is(err, io.EOF) {
			err = nil
			break
		}
		if err != nil {
			return
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
			return
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
			err = nodeErr(pe.EIllegalTopLevelNode, n)
			return
		}
	}

	tree = p.Tree
	return
}

func tokErr(c pe.ErrorCode, cause *lexer.Token) *pe.PrettyError {
	return pe.New(c).Section("Caused by", "%s \"%s\" at %s", cause.Kind, cause.Raw, cause.Where)
}

func nodeErr(c pe.ErrorCode, cause ast.MetaNode) *pe.PrettyError {
	return pe.New(c).Section("Caused by", "`%s` node at %s", cause.Node.Kind(), cause.Where)
}
