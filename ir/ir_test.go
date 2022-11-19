package ir_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/syzkrash/skol/ir"
	"golang.org/x/exp/slices"
)

var (
	//go:embed example.skir
	encodedExample []byte

	exampleProgram = ir.Program{
		Entrypoint: 0,
		Globals: []ir.Value{
			ir.IntegerValue{
				Value: 123,
			},
		},
		Funcs: []ir.Block{
			{
				ir.SetInstr{
					Target: ir.SingleRef{
						RefType: ir.RefGlobal,
						Idx:     0,
					},
					Value: ir.IntegerValue{
						Value: 321,
					},
				},
			},
		},
	}
)

func TestEncode(t *testing.T) {
	encodedBuf := bytes.Buffer{}
	err := ir.Encode(&encodedBuf, exampleProgram)
	if err != nil {
		t.Fatalf("encoding error: %s", err)
	}
	if !slices.Equal(encodedBuf.Bytes(), encodedExample) {
		t.Fatalf("encoding difference:\n%+v\n!=\n%+v", encodedBuf.Bytes(), encodedExample)
	}
}

func TestDecode(t *testing.T) {
	encodedRdr := bytes.NewReader(encodedExample)
	p, err := ir.Decode(encodedRdr)
	if err != nil {
		t.Fatalf("decoding error: %s", err)
	}
	if p.Entrypoint != exampleProgram.Entrypoint {
		t.Fatalf("entrypoint mismatch: %02X != %02X", p.Entrypoint, exampleProgram.Entrypoint)
	}
	if len(p.Globals) != len(exampleProgram.Globals) {
		t.Fatalf("global count mismatch: %02X != %02X", len(p.Globals), len(exampleProgram.Globals))
	}
	for i := 0; i < len(p.Globals); i++ {
		pg := p.Globals[i]
		eg := exampleProgram.Globals[i]
		if pg.Type() != eg.Type() {
			t.Fatalf("global %02X type mismatch: %02X != %02X", i, pg.Type(), eg.Type())
		}
	}
	if len(p.Funcs) != len(exampleProgram.Funcs) {
		t.Fatalf("func count mismatch: %02x != %02X", len(p.Funcs), len(exampleProgram.Funcs))
	}
	for i := 0; i < len(p.Funcs); i++ {
		pf := p.Funcs[i]
		ef := exampleProgram.Funcs[i]
		if len(pf) != len(ef) {
			t.Fatalf("func %02X len mismatch: %02X != %02X", i, len(pf), len(ef))
		}
	}

	t.Logf("Program:\n%s", p)
}

func TestRecode(t *testing.T) {
	encodedBuf := bytes.Buffer{}
	err := ir.Encode(&encodedBuf, exampleProgram)
	if err != nil {
		t.Fatalf("encoding error: %s", err)
	}
	p, err := ir.Decode(bytes.NewBuffer(encodedBuf.Bytes()))
	if err != nil {
		t.Fatalf("decoding error: %s", err)
	}
	if p.Entrypoint != exampleProgram.Entrypoint {
		t.Fatalf("entrypoint mismatch: %02X != %02X", p.Entrypoint, exampleProgram.Entrypoint)
	}
	if len(p.Globals) != len(exampleProgram.Globals) {
		t.Fatalf("global count mismatch: %02X != %02X", len(p.Globals), len(exampleProgram.Globals))
	}
	for i := 0; i < len(p.Globals); i++ {
		pg := p.Globals[i]
		eg := exampleProgram.Globals[i]
		if pg.Type() != eg.Type() {
			t.Fatalf("global %02X type mismatch: %02X != %02X", i, pg.Type(), eg.Type())
		}
	}
	if len(p.Funcs) != len(exampleProgram.Funcs) {
		t.Fatalf("func count mismatch: %02x != %02X", len(p.Funcs), len(exampleProgram.Funcs))
	}
	for i := 0; i < len(p.Funcs); i++ {
		pf := p.Funcs[i]
		ef := exampleProgram.Funcs[i]
		if len(pf) != len(ef) {
			t.Fatalf("func %02X len mismatch: %02X != %02X", i, len(pf), len(ef))
		}
	}
}
