package lexer

import (
	"errors"
	"fmt"
	"io"
)

type Lexer struct {
	src  *Source
	prev *Token
}

func NewLexer(src io.RuneScanner, fn string) *Lexer {
	return &Lexer{
		src: NewSource(src, fn),
	}
}

func (l *Lexer) selfError(msg string) error {
	fmt.Println(l.src.Position)
	return &LexerError{
		msg:   msg,
		cause: nil,
		Where: l.src.Position,
	}
}

func (l *Lexer) otherError(cause error) error {
	return &LexerError{
		msg:   cause.Error(),
		cause: cause,
		Where: l.src.Position,
	}
}

func (l *Lexer) nextIdent(c rune) (tok *Token, err error) {
	pos := l.src.Position
	ident := string(c)
	for {
		c, _, err = l.src.ReadRune()
		if errors.Is(err, io.EOF) {
			err = nil
			goto finish
		}
		if err != nil {
			return
		}
		if !isIdent(c) && !isDigit(c) {
			if err = l.src.UnreadRune(); err != nil {
				return
			}
			break
		}
		ident += string(c)
	}
finish:
	tok = &Token{
		Kind:  TkIdent,
		Where: pos,
		Raw:   ident,
	}
	return
}

func (l *Lexer) nextConstant(c rune) (tok *Token, err error) {
	pos := l.src.Position
	num := string(c)
	for {
		c, _, err = l.src.ReadRune()
		if errors.Is(err, io.EOF) {
			err = nil
			goto finish
		}
		if err != nil {
			return
		}
		if !isNumberTail(c) {
			if err = l.src.UnreadRune(); err != nil {
				return
			}
			break
		}
		num += string(c)
	}
finish:
	tok = &Token{
		Kind:  TkConstant,
		Where: pos,
		Raw:   num,
	}
	return
}

func (l *Lexer) nextString() (tok *Token, err error) {
	var c rune
	pos := l.src.Position
	str := ""
	for {
		c, _, err = l.src.ReadRune()
		if err != nil {
			return
		}
		if c == '\\' {
			var e rune
			e, _, err = l.src.ReadRune()
			if err != nil {
				return
			}
			e, err = escapeSeq(e)
			if err != nil {
				return
			}
			str += string(e)
			continue
		}
		if c == '"' {
			break
		}
		str += string(c)
	}
	tok = &Token{
		Kind:  TkString,
		Where: pos,
		Raw:   str,
	}
	return
}

func (l *Lexer) nextChar() (tok *Token, err error) {
	var c rune
	pos := l.src.Position
	var lit rune
	c, _, err = l.src.ReadRune()
	if err != nil {
		return
	}
	if c == '\\' {
		var e rune
		e, _, err = l.src.ReadRune()
		if err != nil {
			return
		}
		lit, err = escapeSeq(e)
		if err != nil {
			return
		}
	} else {
		lit = c
	}
	c, _, err = l.src.ReadRune()
	if err != nil {
		return
	}
	if c != '\'' {
		err = errors.New("invalid character literal")
	}
	tok = &Token{
		Kind:  TkChar,
		Where: pos,
		Raw:   string(lit),
	}
	return
}

func (l *Lexer) nextPunctuator(c rune) (tok *Token, ok bool) {
	switch c {
	case '(', ')', '$', '%', ':', '/', '>', '?', '*', '#', '@':
		tok = &Token{
			Kind:  TkPunct,
			Where: l.src.Position,
			Raw:   string(c),
		}
		ok = true
	}
	return
}

func (l *Lexer) ignoreLineComment() (err error) {
	var c rune
	for c != '\n' {
		c, _, err = l.src.ReadRune()
		if err != nil {
			return
		}
	}
	return
}

func (l *Lexer) ignoreBlockComment() (err error) {
	var c rune
	for {
		c, _, err = l.src.ReadRune()
		if err != nil {
			return
		}
		if c != '*' {
			continue
		}
		c, _, err = l.src.ReadRune()
		if err != nil {
			return
		}
		if c == '/' {
			break
		}
	}
	return
}

func (l *Lexer) commentOrSlash() (comment bool, err error) {
	var c rune
	c, _, err = l.src.ReadRune()
	if err != nil {
		return
	}

	switch c {
	case '/':
		comment = true
		err = l.ignoreLineComment()
	case '*':
		comment = true
		err = l.ignoreBlockComment()
	default:
		comment = false
		err = l.src.UnreadRune()
	}

	return
}

func (l *Lexer) internalNext() (tok *Token, err error) {
	c, _, err := l.src.ReadRune()
	if err != nil {
		return
	}

	for {
		if isSpace(c) {
			for isSpace(c) {
				c, _, err = l.src.ReadRune()
				if err != nil {
					return
				}
			}
			continue
		}
		if c == '/' {
			var cmt bool
			cmt, err = l.commentOrSlash()
			if err != nil {
				return
			}
			if !cmt {
				tok = &Token{
					Kind:  TkPunct,
					Where: l.src.Position,
					Raw:   "/",
				}
				return
			}
			c, _, err = l.src.ReadRune()
			if err != nil {
				return
			}
			continue
		}
		break
	}

	switch {
	case isIdent(c):
		tok, err = l.nextIdent(c)
	case isNumberHead(c):
		tok, err = l.nextConstant(c)
	case c == '"':
		tok, err = l.nextString()
	case c == '\'':
		tok, err = l.nextChar()
	default:
		var ok bool
		tok, ok = l.nextPunctuator(c)
		if !ok {
			err = l.selfError("illegal token: " + string(c))
		}
	}

	return
}

func (l *Lexer) Next() (tok *Token, err error) {
	if l.prev != nil {
		tok = l.prev
		l.prev = nil
		return
	}
	tok, err = l.internalNext()
	var lerr *LexerError
	if err != nil && !errors.As(err, &lerr) {
		err = l.otherError(err)
	}
	return
}

func (l *Lexer) Rollback(tok *Token) {
	l.prev = tok
}
