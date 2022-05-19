package parser

import "github.com/syzkrash/skol/lexer"

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
