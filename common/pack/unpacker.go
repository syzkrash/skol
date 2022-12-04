package pack

import (
	"io"
	"math"
)

// Unpacker contains facilities for reading concrete values from byte streams
// read from an [io.Reader]
type Unpacker struct {
	in     io.Reader
	Err    []error
	buf    [8]byte
	Offset uint32
}

// NewUnpacker creates a [Unpacker] for the given [io.Reader]
func NewUnpacker(in io.Reader) *Unpacker {
	return &Unpacker{
		in:  in,
		Err: make([]error, 0, 16),
		buf: [8]byte{0},
	}
}

// Error adds an error to thep unpacker's error array (useful if you don't want
// to return errors yourself)
func (u *Unpacker) Error(e error) {
	u.Err = append(u.Err, e)
}

func (u *Unpacker) read(p []byte) {
	if _, err := u.in.Read(p); err != nil {
		u.Error(err)
	} else {
		u.Offset += uint32(len(p))
	}
}

func (u *Unpacker) readN(n uint) {
	u.read(u.buf[:n])
}

// U8 reads an unsigned 8-bit int
func (u *Unpacker) U8() uint8 {
	u.readN(1)
	return u.buf[0]
}

// U16 reads an unsigned 16-bit int
func (u *Unpacker) U16() uint16 {
	u.readN(2)
	return ((uint16(u.buf[0]) << 8) |
		uint16(u.buf[1]))
}

// U32 reads an unsigned 32-bit int
func (u *Unpacker) U32() uint32 {
	u.readN(4)
	return ((uint32(u.buf[0]) << 24) |
		(uint32(u.buf[1]) << 16) |
		(uint32(u.buf[2]) << 8) |
		uint32(u.buf[3]))
}

// U64 reads an unsigned 64-bit int
func (u *Unpacker) U64() uint64 {
	u.readN(8)
	return ((uint64(u.buf[0]) << 56) |
		(uint64(u.buf[1]) << 48) |
		(uint64(u.buf[2]) << 40) |
		(uint64(u.buf[3]) << 32) |
		(uint64(u.buf[4]) << 24) |
		(uint64(u.buf[5]) << 16) |
		(uint64(u.buf[6]) << 8) |
		uint64(u.buf[7]))
}

// I8 reads a signed 8-bit int
func (u *Unpacker) I8() int8 {
	return int8(u.U8())
}

// I16 reads a signed 16-bit int
func (u *Unpacker) I16() int16 {
	return int16(u.U16())
}

// I32 reads a signed 32-bit int
func (u *Unpacker) I32() int32 {
	return int32(u.U32())
}

// I64 reads a signed 64-bit int
func (u *Unpacker) I64() int64 {
	return int64(u.U64())
}

// F32 reads a 32-bit floating point value from it's IEEE 754 representation
func (u *Unpacker) F32() float32 {
	return math.Float32frombits(u.U32())
}

// F64 reads a 64-bit floating point value from it's IEEE 754 representation
func (u *Unpacker) F64() float64 {
	return math.Float64frombits(u.U64())
}

// Bytes reads at most N bytes into a new slice
func (u *Unpacker) Bytes(n uint) []byte {
	b := make([]byte, n)
	u.read(b)
	return b
}

// Str reads a string's length and then reads that amount of bytes as a string
func (u *Unpacker) Str() string {
	length := u.U8()
	buf := make([]byte, length)
	u.read(buf)
	return string(buf)
}
