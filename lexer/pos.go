package lexer

import "fmt"

// Position represents a location within a file.
type Position struct {
	File string
	Line uint
	Col  uint
}

// String returns a formatted string of this Position
func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Col)
}
