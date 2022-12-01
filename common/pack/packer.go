package pack

import (
	"io"
	"math"
)

// Packer contains facilities for writing byte streams to an [io.Writer]
type Packer struct {
	out io.Writer
	Err []error
	buf [8]byte
}

// NewPacker creates a [Packer] for the given [io.Writer]
func NewPacker(out io.Writer) *Packer {
	return &Packer{
		out: out,
		Err: make([]error, 0, 16),
		buf: [8]byte{0},
	}
}

// Error adds an error to thep packer's error array (useful if you don't want to
// return errors yourself)
func (p *Packer) Error(e error) {
	p.Err = append(p.Err, e)
}

// Writes writes the provided bytes directly
func (p *Packer) Write(b []byte) *Packer {
	if _, err := p.out.Write(b); err != nil {
		p.Error(err)
	}
	return p
}

func (p *Packer) writeBuf(bc uint) *Packer {
	return p.Write(p.buf[:bc])
}

// U8 writes an unsigned 8-bit int
func (p *Packer) U8(n uint8) *Packer {
	p.buf[0] = n
	return p.writeBuf(1)
}

// U16 writes an unsigned 16-bit int
func (p *Packer) U16(n uint16) *Packer {
	p.buf[1] = byte(n & 0xFF)
	p.buf[0] = byte((n & (0xFF << 8)) >> 8)
	return p.writeBuf(2)
}

// U32 writes an unsigned 32-bit int
func (p *Packer) U32(n uint32) *Packer {
	p.buf[3] = byte(n & 0xFF)
	p.buf[2] = byte((n & (0xFF << 8)) >> 8)
	p.buf[1] = byte((n & (0xFF << 16)) >> 16)
	p.buf[0] = byte((n & (0xFF << 24)) >> 24)
	return p.writeBuf(4)
}

// U64 writes an unsigned 64-bit int
func (p *Packer) U64(n uint64) *Packer {
	p.buf[7] = byte(n & 0xFF)
	p.buf[6] = byte((n & (0xFF << 8)) >> 8)
	p.buf[5] = byte((n & (0xFF << 16)) >> 16)
	p.buf[4] = byte((n & (0xFF << 24)) >> 24)
	p.buf[3] = byte((n & (0xFF << 32)) >> 32)
	p.buf[2] = byte((n & (0xFF << 40)) >> 40)
	p.buf[1] = byte((n & (0xFF << 48)) >> 48)
	p.buf[0] = byte((n & (0xFF << 56)) >> 56)
	return p.writeBuf(8)
}

// I8 writes an signed 8-bit int
func (p *Packer) I8(n int8) *Packer {
	return p.U8(uint8(n))
}

// I16 writes an signed 16-bit int
func (p *Packer) I16(n int16) *Packer {
	return p.U16(uint16(n))
}

// I32 writes an signed 32-bit int
func (p *Packer) I32(n int32) *Packer {
	return p.U32(uint32(n))
}

// I64 writes an signed 64-bit int
func (p *Packer) I64(n int64) *Packer {
	return p.U64(uint64(n))
}

// F32 writes a IEE 754 representation of a 32-bit floating point value
func (p *Packer) F32(n float32) *Packer {
	return p.U32(math.Float32bits(n))
}

// F64 writes a IEE 754 representation of a 64-bit floating point value
func (p *Packer) F64(n float64) *Packer {
	return p.U64(math.Float64bits(n))
}

// Str writes the string's length followed by the string's bytes
func (p *Packer) Str(s string) *Packer {
	return p.U8(uint8(len(s))).Write([]byte(s))
}
