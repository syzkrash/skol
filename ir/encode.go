package ir

import (
	"fmt"
	"io"

	"github.com/syzkrash/skol/common/pack"
)

// Encode writes a full IR representation of the program to the given writer,
// WITH the magic string and version.
func Encode(w io.Writer, p Program) (err error) {
	pk := pack.NewPacker(w)
	pk.Write([]byte(magic)).U8(ver).U8(p.Entrypoint)
	encodeValueArray(pk, p.Globals)
	encodeBlockArray(pk, p.Funcs)
	if len(pk.Err) > 0 {
		return pk.Err[0]
	}
	return
}

func encodeValue(pk *pack.Packer, v Value) {
	pk.U8(uint8(v.Type()))
	switch v.Type() {
	case TypeInteger:
		pk.I64(v.(IntegerValue).Value)
	case TypeFloat:
		pk.F64(v.(FloatValue).Value)
	case TypeCall:
		cv := v.(CallValue)
		pk.U8(cv.Func)
		encodeValueArray(pk, cv.Args)
	case TypeStruct:
		encodeValueArray(pk, v.(StructValue).Fields)
	case TypeArray:
		encodeValueArray(pk, v.(ArrayValue).Elements)
	case TypeRef:
		encodeRef(pk, v.(RefValue).Ref)
	default:
		pk.Error(fmt.Errorf("unknown value type: %02X", v.Type()))
	}
}

func encodeValueArray(pk *pack.Packer, va []Value) {
	pk.U8(uint8(len(va)))
	for _, v := range va {
		encodeValue(pk, v)
	}
}

func encodeRef(pk *pack.Packer, r Ref) {
	pk.U8(uint8(r.Type()))
	switch r.Type() {
	case RefLocal, RefGlobal:
		pk.U8(r.(SingleRef).Idx)
	case RefLocalIdx, RefGlobalIdx:
		dr := r.(DoubleRef)
		pk.U8(dr.Val).U32(dr.Idx)
	default:
		pk.Error(fmt.Errorf("unknown reference type: %02X", r.Type()))
	}
}

func encodeBlock(pk *pack.Packer, b Block) {
	pk.U8(uint8(len(b)))
	for _, i := range b {
		encodeInstr(pk, i)
	}
}

func encodeInstr(pk *pack.Packer, i Instr) {
	pk.U8(uint8(i.Op()))
	switch i.Op() {
	case OpSet:
		si := i.(SetInstr)
		encodeRef(pk, si.Target)
		encodeValue(pk, si.Value)
	case OpCall:
		ci := i.(CallInstr)
		pk.U8(ci.Func)
		encodeValueArray(pk, ci.Args)
	case OpRet:
		encodeValue(pk, i.(RetInstr).Value)
	case OpBranch:
		encodeBranchArray(pk, i.(BranchInstr).branches)
	case OpLoop:
		encodeBranchOf(pk, i.(LoopInstr).Cond, i.(LoopInstr).Body)
	default:
		pk.Error(fmt.Errorf("unknown instruction: %02X", i.Op()))
	}
}

func encodeBlockArray(pk *pack.Packer, ba []Block) {
	pk.U8(uint8(len(ba)))
	for _, b := range ba {
		encodeBlock(pk, b)
	}
}

func encodeBranch(pk *pack.Packer, b branch) {
	encodeValue(pk, b.Cond)
	encodeBlock(pk, b.Body)
}

func encodeBranchOf(pk *pack.Packer, cond Value, body Block) {
	encodeBranch(pk, branch{
		Cond: cond,
		Body: body,
	})
}

func encodeBranchArray(pk *pack.Packer, ba []branch) {
	pk.U8(uint8(len(ba)))
	for _, b := range ba {
		encodeBranch(pk, b)
	}
}
