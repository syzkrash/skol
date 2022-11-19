package ir

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// Encode writes a full IR representation of the program to the given writer,
// WITH the magic string and version.
func Encode(w io.Writer, p Program) (err error) {
	_, err = w.Write([]byte(magic))
	if err != nil {
		return
	}
	_, err = w.Write([]byte{ver, p.Entrypoint})
	if err != nil {
		return
	}
	err = encodeValueArray(w, p.Globals)
	if err != nil {
		return
	}
	err = encodeBlockArray(w, p.Funcs)
	return
}

func encodeValue(w io.Writer, v Value) (err error) {
	_, err = w.Write([]byte{byte(v.Type())})
	if err != nil {
		return
	}
	switch v.Type() {
	case TypeInteger:
		unsigned := make([]byte, 8)
		binary.BigEndian.PutUint64(unsigned, uint64(v.(IntegerValue).Value))
		_, err = w.Write(unsigned)
		return
	case TypeFloat:
		floatBits := math.Float64bits(v.(FloatValue).Value)
		unsigned := make([]byte, 8)
		binary.BigEndian.PutUint64(unsigned, floatBits)
		_, err = w.Write(unsigned)
		return
	case TypeCall:
		cv := v.(CallValue)
		_, err = w.Write([]byte{cv.Func})
		if err != nil {
			return
		}
		err = encodeValueArray(w, cv.Args)
		return
	case TypeStruct:
		err = encodeValueArray(w, v.(StructValue).Fields)
		return
	case TypeArray:
		err = encodeValueArray(w, v.(ArrayValue).Elements)
		return
	case TypeRef:
		err = encodeRef(w, v.(RefValue).Ref)
		return
	default:
		err = fmt.Errorf("unknown value type: %02X", v.Type())
		return
	}
}

func encodeValueArray(w io.Writer, va []Value) (err error) {
	_, err = w.Write([]byte{byte(len(va))})
	for _, v := range va {
		err = encodeValue(w, v)
		if err != nil {
			return
		}
	}
	return
}

func encodeRef(w io.Writer, r Ref) (err error) {
	_, err = w.Write([]byte{byte(r.Type())})
	if err != nil {
		return
	}
	switch r.Type() {
	case RefLocal, RefGlobal:
		_, err = w.Write([]byte{r.(SingleRef).Idx})
		return
	case RefLocalIdx, RefGlobalIdx:
		dr := r.(DoubleRef)
		_, err = w.Write([]byte{dr.Val})
		if err != nil {
			return
		}
		rawIdx := make([]byte, 4)
		binary.BigEndian.PutUint32(rawIdx, dr.Idx)
		_, err = w.Write(rawIdx)
		return
	default:
		err = fmt.Errorf("unknown reference type: %02X", r.Type())
		return
	}
}

func encodeBlock(w io.Writer, b Block) (err error) {
	_, err = w.Write([]byte{byte(len(b))})
	if err != nil {
		return
	}
	for _, i := range b {
		err = encodeInstr(w, i)
		if err != nil {
			return
		}
	}
	return
}

func encodeInstr(w io.Writer, i Instr) (err error) {
	_, err = w.Write([]byte{byte(i.Op())})
	if err != nil {
		return
	}
	switch i.Op() {
	case OpSet:
		si := i.(SetInstr)
		err = encodeRef(w, si.Target)
		if err != nil {
			return
		}
		err = encodeValue(w, si.Value)
		return
	case OpCall:
		ci := i.(CallInstr)
		_, err = w.Write([]byte{ci.Func})
		if err != nil {
			return
		}
		err = encodeValueArray(w, ci.Args)
		return
	case OpRet:
		err = encodeValue(w, i.(RetInstr).Value)
		return
	case OpBranch:
		err = encodeBranchArray(w, i.(BranchInstr).branches)
		return
	case OpLoop:
		err = encodeBranchOf(w, i.(LoopInstr).Cond, i.(LoopInstr).Body)
		return
	default:
		err = fmt.Errorf("unknown instruction: %02X", i.Op())
		return
	}
}

func encodeBlockArray(w io.Writer, ba []Block) (err error) {
	_, err = w.Write([]byte{byte(len(ba))})
	if err != nil {
		return
	}
	for _, b := range ba {
		err = encodeBlock(w, b)
		if err != nil {
			return
		}
	}
	return
}

func encodeBranch(w io.Writer, b branch) (err error) {
	err = encodeValue(w, b.Cond)
	if err != nil {
		return
	}
	err = encodeBlock(w, b.Body)
	return
}

func encodeBranchOf(w io.Writer, cond Value, body Block) (err error) {
	return encodeBranch(w, branch{
		Cond: cond,
		Body: body,
	})
}

func encodeBranchArray(w io.Writer, ba []branch) (err error) {
	_, err = w.Write([]byte{byte(len(ba))})
	if err != nil {
		return
	}
	for _, b := range ba {
		err = encodeBranch(w, b)
		if err != nil {
			return
		}
	}
	return
}
