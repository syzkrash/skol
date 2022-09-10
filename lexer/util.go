package lexer

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isIdent(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_' || c == '!'
}

func isDigit(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func isNumberHead(c rune) bool {
	return isDigit(c) || c == '-'
}

func isNumberTail(c rune) bool {
	return isDigit(c) ||
		c == '.' || // decimals
		c == 'b' || c == 'o' || c == 'x' || // 0b-, 0o-, 0x-
		c == '_' // digit separator
}

func escapeSeq(e rune) (c rune, ok bool) {
	ok = true
	switch e {
	case '"', '\'':
		c = e
	case 'n':
		c = '\n'
	case 'r':
		c = '\r'
	case 't':
		c = '\t'
	case '\\':
		c = '\\'
	default:
		ok = false
	}
	return
}
