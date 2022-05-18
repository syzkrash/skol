package lexer

import "io"

// Source wraps a RuneScanner and tracks it's position within the file, with
// Position always pointing to the last read rune
type Source struct {
	Scanner  io.RuneScanner
	Position Position
	lastCol  uint
	lastRune rune
}

// NewSource wraps any RuneScanner with the given filename for it's Position
// value
func NewSource(s io.RuneScanner, fn string) *Source {
	return &Source{
		Scanner: s,
		Position: Position{
			File: fn,
			Line: 1,
			Col:  0,
		},
	}
}

// ReadRune reads 1 rune from the underlying RuneScanner and advances the
// Position value accordingly
func (s *Source) ReadRune() (c rune, l int, err error) {
	c, l, err = s.Scanner.ReadRune()
	if c == '\n' {
		s.Position.Line++
		s.lastCol = s.Position.Col
		s.Position.Col = 0
	} else {
		s.Position.Col++
	}
	s.lastRune = c
	return
}

// UnreadRune unreads 1 rune from the underlying RuneScanner and properly
// backs the Position value
func (s *Source) UnreadRune() (err error) {
	err = s.Scanner.UnreadRune()
	if s.lastRune == '\n' {
		s.Position.Line--
		s.Position.Col = s.lastCol
		s.lastRune = 0
	}
	return err
}
