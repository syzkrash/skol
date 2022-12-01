package ir

import (
	"errors"
	"fmt"
	"io"

	"github.com/syzkrash/skol/common/pack"
)

const (
	magic = "SKIR"
	ver   = 1
)

// Decode reads an encoded Program from the given Reader. This assumes the
// input starts with the "SKIR" magic string.
func Decode(r io.Reader) (prog Program, err error) {
	u := pack.NewUnpacker(r)
	magicBytes := u.Bytes(uint(len(magic)))
	if string(magicBytes) != magic {
		err = errors.New("invalid or missing magic")
	}
	if fv := u.U8(); fv != ver {
		err = fmt.Errorf("incorrect IR version: %02X (expected %02X)", fv, ver)
		return
	}
	prog.Entrypoint = u.U8()

	count := u.U8()
	prog.Globals = make([]Value, count)
	for i := 0; i < int(count); i++ {
		prog.Globals[i] = decodeValue(u)
	}

	count = u.U8()
	prog.Funcs = make([]Block, count)
	for i := 0; i < int(count); i++ {
		prog.Funcs[i] = decodeBlock(u)
	}

	if len(u.Err) > 0 {
		err = u.Err[0]
	}

	return
}

func decodeValue(u *pack.Unpacker) (val Value) {
	ty := Type(u.U8())
	switch ty {
	case TypeInteger:
		val = IntegerValue{Value: u.I64()}
	case TypeFloat:
		val = FloatValue{Value: u.F64()}
	case TypeCall:
		fn := u.U8()
		count := u.U8()
		args := make([]Value, count)
		for i := 0; i < int(count); i++ {
			args[i] = decodeValue(u)
		}
		val = CallValue{
			Func: fn,
			Args: args,
		}
	case TypeStruct:
		count := u.U8()
		fields := make([]Value, count)
		for i := 0; i < int(count); i++ {
			fields[i] = decodeValue(u)
		}
		val = StructValue{
			Fields: fields,
		}
	case TypeArray:
		count := u.U8()
		elems := make([]Value, count)
		for i := 0; i < int(count); i++ {
			elems[i] = decodeValue(u)
		}
		val = ArrayValue{
			Elements: elems,
		}
	case TypeRef:
		val = RefValue{
			Ref: decodeRef(u),
		}
	default:
		u.Error(fmt.Errorf("unknown value type: %02X", ty))
	}
	return
}

func decodeRef(u *pack.Unpacker) (ref Ref) {
	rt := RefType(u.U8())
	switch rt {
	case RefLocal, RefGlobal:
		ref = SingleRef{
			RefType: rt,
			Idx:     u.U8(),
		}
	case RefLocalIdx, RefGlobalIdx:
		// read these here to make sure they are read in the correct order
		v := u.U8()
		i := u.U32()
		ref = DoubleRef{
			RefType: rt,
			Val:     v,
			Idx:     i,
		}
	default:
		u.Error(fmt.Errorf("unknown reference type: %02X", rt))
	}
	return
}

func decodeBlock(u *pack.Unpacker) (block Block) {
	count := u.U8()
	block = make(Block, count)
	for i := 0; i < int(count); i++ {
		block[i] = decodeInstr(u)
	}
	return
}

func decodeInstr(u *pack.Unpacker) (instr Instr) {
	op := Opcode(u.U8())
	switch op {
	case OpSet:
		target := decodeRef(u)
		val := decodeValue(u)
		instr = SetInstr{
			Target: target,
			Value:  val,
		}
	case OpCall:
		fn := u.U8()
		count := u.U8()
		args := make([]Value, count)
		for i := 0; i < int(count); i++ {
			args[i] = decodeValue(u)
		}
		instr = CallInstr{
			Func: fn,
			Args: args,
		}
	case OpRet:
		instr = RetInstr{
			Value: decodeValue(u),
		}
	case OpBranch:
		count := u.U8()
		branches := make([]branch, count)
		for i := 0; i < int(count); i++ {
			branches[i] = decodeBranch(u)
		}
		instr = BranchInstr{
			branches: branches,
		}
	case OpLoop:
		branch := decodeBranch(u)
		instr = LoopInstr{
			Cond: branch.Cond,
			Body: branch.Body,
		}
	default:
		u.Error(fmt.Errorf("unknown instruction: %02X", op))
	}
	return
}

func decodeBranch(u *pack.Unpacker) (branch branch) {
	branch.Cond = decodeValue(u)
	branch.Body = decodeBlock(u)
	return
}
