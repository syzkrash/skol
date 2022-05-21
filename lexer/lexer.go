package lexer

import (
	"errors"
	"io"
)

type Lexer struct {
	src *Source
}

func NewLexer(src io.RuneScanner, fn string) *Lexer {
	return &Lexer{
		src: NewSource(src, fn),
	}
}

func (l *Lexer) selfError(msg string) error {
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

func (l *Lexer) ignoreSpace(c rune) (err error) {
	for isSpace(c) {
		c, _, err = l.src.ReadRune()
		if err != nil {
			return err
		}
	}
	return l.src.UnreadRune()
}

func (l *Lexer) ignoreComment(c rune) (err error) {
	for c != '\n' {
		c, _, err = l.src.ReadRune()
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Lexer) ignoreSpacesAndComments() error {
	for {
		c, _, err := l.src.ReadRune()
		if err != nil {
			return err
		}
		if isSpace(c) {
			if err = l.ignoreSpace(c); err != nil {
				return err
			}
			continue
		}
		if c == '#' {
			if err = l.ignoreComment(c); err != nil {
				return err
			}
			continue
		}
		break
	}
	return l.src.UnreadRune()
}

func (l *Lexer) nextIdent(c rune) (tok *Token, err error) {
	pos := l.src.Position
	ident := string(c)
	for {
		c, _, err = l.src.ReadRune()
		if errors.Is(err, io.EOF) {
			goto finish
		}
		if err != nil {
			return
		}
		if !isIdent(c) {
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
	case '(', ')', '$', '%', ':', '/', '>':
		tok = &Token{
			Kind:  TkPunct,
			Where: l.src.Position,
			Raw:   string(c),
		}
		ok = true
	}
	return
}

func (l *Lexer) internalNext() (tok *Token, err error) {
	err = l.ignoreSpacesAndComments()
	if err != nil {
		return
	}
	c, _, err := l.src.ReadRune()
	if err != nil {
		return
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
			err = l.selfError("illegal token")
		}
	}

	return
}

func (l *Lexer) Next() (tok *Token, err error) {
	tok, err = l.internalNext()
	var lerr *LexerError
	if err != nil && !errors.As(err, &lerr) {
		err = l.otherError(err)
	}
	return
}
