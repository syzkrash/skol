package lexer

import "fmt"

type LexerError struct {
	Where Position
	Char  rune
	msg   string
	cause error
}

func (e *LexerError) Error() string {
	return e.msg
}

func (e *LexerError) Unwrap() error {
	return e.cause
}

func (e *LexerError) Print() {
	fmt.Println("\x1b[31m\x1b[1mLexer error:\x1b[0m")
	fmt.Println("   ", e.msg)
	fmt.Println("\x1b[1mCaused by:\x1b[0m")
	fmt.Println("   ", "`"+string(e.Char)+"`", "at", e.Where)
	fmt.Println()
}
