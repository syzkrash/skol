package parser

import (
	"fmt"

	"github.com/syzkrash/skol/lexer"
)

// ParserError defines an error that occurred while the parser was trying to
// construct a node. Contains extra information about the token causing the
// error.
type ParserError struct {
	Where *lexer.Token
	msg   string
	cause error
}

func (e *ParserError) Error() string {
	return e.msg
}

func (e *ParserError) Unwrap() error {
	return e.cause
}

func (e *ParserError) Print() {
	fmt.Println("\x1b[31m\x1b[1mParser error:\x1b[0m")
	fmt.Println("   ", e.msg)
	fmt.Println("\x1b[1mCaused by:\x1b[0m")
	fmt.Println("   ", e.Where)
	fmt.Println()
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
