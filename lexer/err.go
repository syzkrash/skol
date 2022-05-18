package lexer

type LexerError struct {
	Where Position
	msg   string
	cause error
}

func (e *LexerError) Error() string {
	return e.msg
}

func (e *LexerError) Unwrap() error {
	return e.cause
}
