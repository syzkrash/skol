package ir

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

const (
	magic = "SKIR"
	ver   = 1
)

// Decode reads an encoded Program from the given Reader. This assumes the
// input starts with the "SKIR" magic string.
func Decode(r io.Reader) (prog Program, err error) {
	magicBytes := make([]byte, len(magic))
	_, err = r.Read(magicBytes)
	if err != nil {
		return
	}
	if string(magicBytes) != magic {
		err = errors.New("invalid or missing magic")
	}
	p := make([]byte, 1)
	_, err = r.Read(p)
	if err != nil {
		return
	}
	if p[0] != ver {
		err = fmt.Errorf("incorrect IR version: %02X (expected %02X)", p[0], ver)
		return
	}
	_, err = r.Read(p)
	if err != nil {
		return
	}
	prog.Entrypoint = p[0]
	_, err = r.Read(p)
	if err != nil {
		return
	}
	prog.Globals = make([]Value, p[0])
	for i := 0; i < int(p[0]); i++ {
		var val Value
		val, err = decodeValue(r)
		if err != nil {
			return
		}
		prog.Globals[i] = val
	}
	_, err = r.Read(p)
	if err != nil {
		return
	}
	prog.Funcs = make([]Block, p[0])
	for i := 0; i < int(p[0]); i++ {
		var body []Instr
		body, err = decodeBlock(r)
		if err != nil {
			return
		}
		prog.Funcs[i] = body
	}

	return
}

func decodeValue(r io.Reader) (val Value, err error) {
	p := make([]byte, 1)
	_, err = r.Read(p)
	if err != nil {
		return
	}
	ty := Type(p[0])
	switch ty {
	case TypeInteger:
		rawInt := make([]byte, 8)
		_, err = r.Read(rawInt)
		if err != nil {
			return
		}
		unsigned := binary.BigEndian.Uint64(rawInt)
		val = IntegerValue{Value: int64(unsigned)}
		return
	case TypeFloat:
		rawFloat := make([]byte, 8)
		_, err = r.Read(rawFloat)
		if err != nil {
			return
		}
		floatBits := binary.BigEndian.Uint64(rawFloat)
		val = FloatValue{Value: math.Float64frombits(floatBits)}
		return
	case TypeCall:
		_, err = r.Read(p)
		if err != nil {
			return
		}
		fn := p[0]
		_, err = r.Read(p)
		if err != nil {
			return
		}
		args := make([]Value, p[0])
		var arg Value
		for i := 0; i < int(p[0]); i++ {
			arg, err = decodeValue(r)
			if err != nil {
				return
			}
			args[i] = arg
		}
		val = CallValue{
			Func: fn,
			Args: args,
		}
		return
	case TypeStruct:
		_, err = r.Read(p)
		if err != nil {
			return
		}
		fields := make([]Value, p[0])
		var field Value
		for i := 0; i < int(p[0]); i++ {
			field, err = decodeValue(r)
			if err != nil {
				return
			}
			fields[i] = field
		}
		val = StructValue{
			Fields: fields,
		}
		return
	case TypeArray:
		_, err = r.Read(p)
		if err != nil {
			return
		}
		elems := make([]Value, p[0])
		var elem Value
		for i := 0; i < int(p[0]); i++ {
			elem, err = decodeValue(r)
			if err != nil {
				return
			}
			elems[i] = elem
		}
		val = ArrayValue{
			Elements: elems,
		}
		return
	case TypeRef:
		var ref Ref
		ref, err = decodeRef(r)
		if err != nil {
			return
		}
		val = RefValue{
			Ref: ref,
		}
		return
	default:
		err = fmt.Errorf("unknown value type: %02X", ty)
		return
	}
}

func decodeRef(r io.Reader) (ref Ref, err error) {
	p := make([]byte, 1)
	_, err = r.Read(p)
	if err != nil {
		return
	}
	rt := RefType(p[0])
	switch rt {
	case RefLocal, RefGlobal:
		_, err = r.Read(p)
		if err != nil {
			return
		}
		ref = SingleRef{
			RefType: rt,
			Idx:     p[0],
		}
		return
	case RefLocalIdx, RefGlobalIdx:
		val := p[0]
		_, err = r.Read(p)
		if err != nil {
			return
		}
		rawIdx := make([]byte, 4)
		_, err = r.Read(rawIdx)
		if err != nil {
			return
		}
		idx := binary.BigEndian.Uint32(rawIdx)
		ref = DoubleRef{
			RefType: rt,
			Val:     val,
			Idx:     idx,
		}
		return
	default:
		err = fmt.Errorf("unknown reference type: %02X", rt)
		return
	}
}

func decodeBlock(r io.Reader) (block Block, err error) {
	p := make([]byte, 1)
	_, err = r.Read(p)
	if err != nil {
		return
	}
	block = make(Block, p[0])
	for i := 0; i < int(p[0]); i++ {
		var instr Instr
		instr, err = decodeInstr(r)
		if err != nil {
			return
		}
		block[i] = instr
	}
	return
}

func decodeInstr(r io.Reader) (instr Instr, err error) {
	p := make([]byte, 1)
	_, err = r.Read(p)
	if err != nil {
		return
	}
	op := Opcode(p[0])
	switch op {
	case OpSet:
		var target Ref
		target, err = decodeRef(r)
		if err != nil {
			return
		}
		var val Value
		val, err = decodeValue(r)
		if err != nil {
			return
		}
		instr = SetInstr{
			Target: target,
			Value:  val,
		}
		return
	case OpCall:
		_, err = r.Read(p)
		if err != nil {
			return
		}
		fn := p[0]
		_, err = r.Read(p)
		if err != nil {
			return
		}
		args := make([]Value, p[0])
		for i := 0; i < int(p[0]); i++ {
			var val Value
			val, err = decodeValue(r)
			if err != nil {
				return
			}
			args[i] = val
		}
		instr = CallInstr{
			Func: fn,
			Args: args,
		}
		return
	case OpRet:
		var val Value
		val, err = decodeValue(r)
		if err != nil {
			return
		}
		instr = RetInstr{
			Value: val,
		}
		return
	case OpBranch:
		_, err = r.Read(p)
		if err != nil {
			return
		}
		branches := make([]branch, p[0])
		for i := 0; i < int(p[0]); i++ {
			var branch branch
			branch, err = decodeBranch(r)
			if err != nil {
				return
			}
			branches[i] = branch
		}
		instr = BranchInstr{
			branches: branches,
		}
		return
	case OpLoop:
		var branch branch
		branch, err = decodeBranch(r)
		if err != nil {
			return
		}
		instr = LoopInstr{
			Cond: branch.Cond,
			Body: branch.Body,
		}
		return
	default:
		err = fmt.Errorf("unknown instruction: %02X", op)
		return
	}
}

func decodeBranch(r io.Reader) (branch branch, err error) {
	branch.Cond, err = decodeValue(r)
	if err != nil {
		return
	}
	branch.Body, err = decodeBlock(r)
	return
}
