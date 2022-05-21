package lexer

import "fmt"

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isIdent(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c == '!'
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isNumberHead(c rune) bool {
	return isDigit(c) || c == '-'
}

func isNumberTail(c rune) bool {
	return isDigit(c) ||
		c == '.' || // decimals
		c == 'b' || c == 'o' || c == 'x' // 0b-, 0o-, 0x-
}

func escapeSeq(e rune) (c rune, err error) {
	switch e {
	case '"', '\'':
		c = e
	case 'n':
		c = '\n'
	case 'r':
		c = '\r'
	case 't':
		c = '\t'
	default:
		err = fmt.Errorf("unknown escape sequence: \\%c", e)
	}
	return
}
