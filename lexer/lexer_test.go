package lexer

import (
	"strings"
	"testing"
)

func TestIdent(t *testing.T) {
	code := `  hello0  `
	read := strings.NewReader(code)
	lex := NewLexer(read, "TestIdent")
	tok, err := lex.Next()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Kind != TIdent {
		t.Fatalf("Incorrect TokenKind! Want Ident but got %s!", tok.Kind)
	}
	if tok.Raw != "hello0" {
		t.Fatalf("Incorrect string! Want `hello0` but got `%s`!", tok.Raw)
	}
}

func TestConstant(t *testing.T) {
	code := `-123.456  `
	read := strings.NewReader(code)
	lex := NewLexer(read, "TestConstant")
	tok, err := lex.Next()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Kind != TFloat {
		t.Fatalf("Incorrect TokenKind! Want Constant but got %s!", tok.Kind)
	}
	if tok.Raw != "-123.456" {
		t.Fatalf("Incorrect string! Want `-123.456` but got `%s`!", tok.Raw)
	}
}

func TestString(t *testing.T) {
	code := ` "hello\tworld\n"`
	read := strings.NewReader(code)
	lex := NewLexer(read, "TestString")
	tok, err := lex.Next()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Kind != TString {
		t.Fatalf("Incorrect TokenKind! Want String but got %s!", tok.Kind)
	}
	if tok.Raw != "hello\tworld\n" {
		t.Fatalf("Incorrect string! Got `%s`!", tok.Raw)
	}
}

func TestChar(t *testing.T) {
	code := `'\''`
	read := strings.NewReader(code)
	lex := NewLexer(read, "TestChar")
	tok, err := lex.Next()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Kind != TChar {
		t.Fatalf("Incorrect TokenKind! Want Char but got %s!", tok.Kind)
	}
	if tok.Raw != "'" {
		t.Fatalf("Incorrect string! Want `'` but got `%s`!", tok.Raw)
	}
}

func TestPunct(t *testing.T) {
	code := `	 	(`
	read := strings.NewReader(code)
	lex := NewLexer(read, "TestPunct")
	tok, err := lex.Next()
	if err != nil {
		t.Fatal(err)
	}
	pn, ok := tok.Punct()
	if !ok {
		t.Fatalf("Incorrect TokenKind! Want Punct but got %s!", tok.Kind)
	}
	if pn != PLParen {
		t.Fatalf("Incorrect string! Want `(` but got `%s`!", tok.Raw)
	}
}
