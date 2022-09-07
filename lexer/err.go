package lexer

import "fmt"

// LexerError represents an error that occured while the lexer was trying to
// read a token. It also includes information on the position in the file at
// which the error occured and the exact character that caused the error.
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

// Print nicely formats the error and prints it to stdout. This is required by
// the [common.Printable] interface.
func (e *LexerError) Print() {
	fmt.Println("\x1b[31m\x1b[1mLexer error:\x1b[0m")
	fmt.Println("   ", e.msg)
	fmt.Println("\x1b[1mCaused by:\x1b[0m")
	fmt.Println("   ", "`"+string(e.Char)+"`", "at", e.Where)
	fmt.Println()
}
