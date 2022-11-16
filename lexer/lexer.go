package lexer

import (
	"errors"
	"io"

	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/debug"
)

// Lexer reads individual runes from a source stream and turns them into tokens.
type Lexer struct {
	src  *Source
	prev *Token
}

// NewLexer creates and prepares a new lexer with the given source stream.
func NewLexer(src io.RuneScanner, fn string) *Lexer {
	return &Lexer{
		src: NewSource(src, fn),
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
			err = pe.New(pe.EBadInput).Cause(err)
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
		Kind:  TIdent,
		Where: pos,
		Raw:   ident,
	}
	return
}

func (l *Lexer) nextConstant(c rune) (tok *Token, err error) {
	pos := l.src.Position
	num := string(c)
	isfloat := false
	for {
		c, _, err = l.src.ReadRune()
		if errors.Is(err, io.EOF) {
			err = nil
			goto finish
		}
		if err != nil {
			err = pe.New(pe.EBadInput).Cause(err)
			return
		}
		if !isNumberTail(c) {
			if err = l.src.UnreadRune(); err != nil {
				return
			}
			break
		}
		if c == '.' {
			isfloat = true
		}
		num += string(c)
	}
finish:
	if isfloat {
		tok = &Token{
			Kind:  TFloat,
			Where: pos,
			Raw:   num,
		}
	} else {
		tok = &Token{
			Kind:  TInt,
			Where: pos,
			Raw:   num,
		}
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
			err = pe.New(pe.EBadInput).Cause(err)
			return
		}
		if c == '\\' {
			var e rune
			e, _, err = l.src.ReadRune()
			if err != nil {
				err = pe.New(pe.EBadInput).Cause(err)
				return
			}
			var ok bool
			e, ok = escapeSeq(e)
			if !ok {
				err = pe.New(pe.EInvalidEscape).Section("Caused by", "'\\%c' at %s", e, l.src.Position)
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
		Kind:  TString,
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
		err = pe.New(pe.EBadInput).Cause(err)
		return
	}
	if c == '\\' {
		var e rune
		e, _, err = l.src.ReadRune()
		if err != nil {
			err = pe.New(pe.EBadInput).Cause(err)
			return
		}
		var ok bool
		lit, ok = escapeSeq(e)
		if !ok {
			err = pe.New(pe.EInvalidEscape).Section("Caused by", "'\\%c' at %s", e, l.src.Position)
			return
		}
	} else {
		lit = c
	}
	c, _, err = l.src.ReadRune()
	if err != nil {
		err = pe.New(pe.EBadInput).Cause(err)
		return
	}
	if c != '\'' {
		err = pe.New(pe.EInvalidCharLit).Section("Caused by", "'%c' at %s", c, l.src.Position)
	}
	tok = &Token{
		Kind:  TChar,
		Where: pos,
		Raw:   string(lit),
	}
	return
}

func (l *Lexer) nextPunctuator(c rune) (tok *Token, ok bool) {
	switch c {
	case '(', ')', '[', ']', '$', '%', ':', '/', '>', '?', '*', '#', '@', '!':
		tok = &Token{
			Kind:  TPunct,
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
			err = pe.New(pe.EBadInput).Cause(err)
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
			err = pe.New(pe.EBadInput).Cause(err)
			return
		}
		if c != '*' {
			continue
		}
		c, _, err = l.src.ReadRune()
		if err != nil {
			err = pe.New(pe.EBadInput).Cause(err)
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
	if errors.Is(err, io.EOF) {
		err = nil
		comment = false
		return
	}
	if err != nil {
		err = pe.New(pe.EBadInput).Cause(err)
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
		err = pe.New(pe.EBadInput).Cause(err)
		return
	}

	for {
		if isSpace(c) {
			for isSpace(c) {
				c, _, err = l.src.ReadRune()
				if err != nil {
					err = pe.New(pe.EBadInput).Cause(err)
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
					Kind:  TPunct,
					Where: l.src.Position,
					Raw:   "/",
				}
				return
			}
			c, _, err = l.src.ReadRune()
			if err != nil {
				err = pe.New(pe.EBadInput).Cause(err)
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
			err = pe.New(pe.EIllegalChar).Section("Caused by", "'%c' at %s", c, l.src.Position)
		}
	}

	return
}

// Next will read and return the next token found in the stream or any error
// that may have occurred in the process as a [*LexerError].
func (l *Lexer) Next() (tok *Token, err error) {
	if l.prev != nil {
		tok = l.prev
		l.prev = nil
		return
	}
	tok, err = l.internalNext()
	if err != nil {
		debug.Log(debug.AttrLexer, "Error %s", err)
	} else {
		debug.Log(debug.AttrLexer, "%s token `%s` at %s", tok.Kind, tok.Raw, tok.Where)
	}
	return
}

// Rollback will save the given token as the token to return on the next call
// to [Next]
func (l *Lexer) Rollback(tok *Token) {
	l.prev = tok
}
