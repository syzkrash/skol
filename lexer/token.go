package lexer

import "fmt"

// TokenKind classifes a token from one of the constants
type TokenKind uint8

// TokenKind constants
const (
	// Identifier
	TkIdent TokenKind = iota
	// Numeric constant
	TkConstant
	// Quoted string literal
	TkString
	// Single-character string literal
	TkChar
	// Punctuator
	TkPunct
)

// used in (TokenKind).String()
var tokenKinds = []string{
	"Ident",
	"Constant",
	"String",
	"Char",
	"Punct",
	"Oper",
}

// String returns the name of this kind of token
func (k TokenKind) String() string {
	return tokenKinds[k]
}

// Token represents a single token of any kind
type Token struct {
	Kind  TokenKind
	Where Position
	Raw   string
}

func (t *Token) String() string {
	return fmt.Sprintf("%s `%s` at %s", t.Kind, t.Raw, t.Where)
}
