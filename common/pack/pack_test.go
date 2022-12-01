package pack_test

import (
	"bytes"
	"testing"

	"github.com/syzkrash/skol/common/pack"
)

var testData = []byte{
	0x01,
	0x01, 0x02,
	0x01, 0x02, 0x03, 0x04,
	0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
	0xFF,
	0xFF, 0xFE,
	0xFF, 0xFF, 0xFF, 0xFD,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFC,
	0x0B, 'h', 'e', 'l', 'l', 'o', ' ', 'w', 'o', 'r', 'l', 'd',
}

// TestUnpack ensures that an [pack.Unpacker] can decode known-good data
func TestUnpack(t *testing.T) {
	u := pack.NewUnpacker(bytes.NewReader(testData))

	expect := func(msg string, a, b any) {
		if a != b {
			t.Fatalf("%s: %+v != %+v", msg, a, b)
		}
		if len(u.Err) > 0 {
			t.Fatal(u.Err[0])
		}
	}

	expect("U8", u.U8(), uint8(0x01))
	expect("U16", u.U16(), uint16(0x01_02))
	expect("U32", u.U32(), uint32(0x01_02_03_04))
	expect("U64", u.U64(), uint64(0x01_02_03_04_05_06_07_08))
	expect("I8", u.I8(), int8(-1))
	expect("I16", u.I16(), int16(-2))
	expect("I32", u.I32(), int32(-3))
	expect("I64", u.I64(), int64(-4))
	expect("Str", u.Str(), "hello world")
}

// TestPack ensures that a [pack.Packer] can encode data correctly, compared
// against known-good data
func TestPack(t *testing.T) {
	buf := bytes.Buffer{}

	p := pack.NewPacker(&buf)
	p.U8(0x01)
	p.U16(0x01_02)
	p.U32(0x01_02_03_04)
	p.U64(0x01_02_03_04_05_06_07_08)
	p.I8(-1)
	p.I16(-2)
	p.I32(-3)
	p.I64(-4)
	p.Str("hello world")

	if len(p.Err) > 0 {
		t.Fatal(p.Err[0])
	}

	b := buf.Bytes()
	for i, b := range b {
		if e := testData[i]; e != b {
			t.Fatalf(
				"packed data does not match test data\n  offset %02X: expected %02X, got %02x",
				i, e, b)
		}
	}
}
