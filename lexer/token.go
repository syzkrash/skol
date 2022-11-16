package lexer

import (
	"fmt"
	"strconv"
)

type TokenKind byte

const (
	TIdent TokenKind = iota
	TInt
	TFloat
	TString
	TChar
	TPunct
)

// used in (TokenKind).String()
var tokenKinds = []string{
	"Ident",
	"Int",
	"Float",
	"String",
	"Char",
	"Punct",
}

// String returns the name of this kind of token
func (k TokenKind) String() string {
	return tokenKinds[k]
}

type Punct byte

const (
	PInvalid Punct = iota
	PLParen
	PRParen
	PLBrack
	PRBrack
	PIs
	PType
	PVar
	PFunc
	PStruct
	PField
	PReturn
	PIf
	PLoop
	PExecute
)

var punctNames = []string{
	"Invalid",
	"Left Paren",
	"Right Paren",
	"Left Bracket",
	"Right Bracket",
	"Is",
	"Type",
	"Var",
	"Func",
	"Struct",
	"Field",
	"Return",
	"If",
	"Loop",
}

func (p Punct) String() string {
	return punctNames[p]
}

type Token struct {
	Kind  TokenKind
	Where Position
	Raw   string
}

func (t Token) Int() (int64, bool) {
	i, err := strconv.ParseInt(t.Raw, 0, 64)
	if err != nil {
		return 0, false
	}
	return i, true
}

func (t Token) Float() (float64, bool) {
	f, err := strconv.ParseFloat(t.Raw, 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func (t Token) Punct() (p Punct, ok bool) {
	if t.Kind != TPunct {
		return 0, false
	}
	ok = true
	switch t.Raw[0] {
	case '(', '{':
		p = PLParen
	case ')', '}':
		p = PRParen
	case '[':
		p = PLBrack
	case ']':
		p = PRBrack
	case ':', '-':
		p = PIs
	case '/':
		p = PType
	case '%':
		p = PVar
	case '$':
		p = PFunc
	case '@':
		p = PStruct
	case '#':
		p = PField
	case '>':
		p = PReturn
	case '?':
		p = PIf
	case '*':
		p = PLoop
	case '!':
		p = PExecute
	default:
		p = PInvalid
		ok = false
	}
	return
}

func (t Token) String() string {
	tn := t.Kind.String()
	if p, ok := t.Punct(); ok {
		tn += fmt.Sprintf(" (%s)", p)
	}
	return fmt.Sprintf("%s `%s` at %s", tn, t.Raw, t.Where)
}
