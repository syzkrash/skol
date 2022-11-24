package ast

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

// Encode writes a binary representation of the given AST into the provided
// [io.Writer]
func Encode(w io.Writer, tree AST) (err error) {
	w.Write([]byte(FormatMagic))
	w.Write([]byte{FormatVersion})

	if err := encodeUint(w, uint(len(tree.Vars))); err != nil {
		return err
	}
	for _, v := range tree.Vars {
		if err := encodeVar(w, v); err != nil {
			return err
		}
	}

	if err := encodeUint(w, uint(len(tree.Typedefs))); err != nil {
		return err
	}
	for _, v := range tree.Typedefs {
		if err := encodeTypedef(w, v); err != nil {
			return err
		}
	}

	if err := encodeUint(w, uint(len(tree.Funcs))); err != nil {
		return err
	}
	for _, f := range tree.Funcs {
		if err := encodeFunc(w, f); err != nil {
			return err
		}
	}

	if err := encodeUint(w, uint(len(tree.Exerns))); err != nil {
		return err
	}
	for _, e := range tree.Exerns {
		if err := encodeExtern(w, e); err != nil {
			return err
		}
	}

	if err := encodeUint(w, uint(len(tree.Structs))); err != nil {
		return err
	}
	for _, s := range tree.Structs {
		if err := encodeStruct(w, s); err != nil {
			return err
		}
	}

	return
}

func encodeVar(w io.Writer, v Var) error {
	if err := encodeString(w, v.Name); err != nil {
		return err
	}
	return encodeNode(w, v.Value)
}

func encodeTypedef(w io.Writer, v Typedef) error {
	if err := encodeString(w, v.Name); err != nil {
		return err
	}
	return encodeType(w, v.Type)
}

func encodeFunc(w io.Writer, f Func) error {
	if err := encodeString(w, f.Name); err != nil {
		return err
	}
	if err := encodeType(w, f.Ret); err != nil {
		return err
	}
	if err := encodeDescriptorSlice(w, f.Args); err != nil {
		return err
	}
	return encodeNodeSlice(w, f.Body)
}

func encodeExtern(w io.Writer, e Extern) error {
	if err := encodeString(w, e.Name); err != nil {
		return err
	}
	if err := encodeType(w, e.Ret); err != nil {
		return err
	}
	return encodeDescriptorSlice(w, e.Args)
}

func encodeStruct(w io.Writer, s Structure) error {
	if err := encodeString(w, s.Name); err != nil {
		return err
	}
	return encodeDescriptorSlice(w, s.Fields)
}

func encodeNode(w io.Writer, mn MetaNode) (err error) {
	if err = encodePos(w, mn.Where); err != nil {
		return err
	}
	k := mn.Node.Kind()
	if _, err = w.Write([]byte{byte(k)}); err != nil {
		return err
	}

	switch k {
	case NBool:
		var b byte
		if mn.Node.(BoolNode).Value {
			b = 1
		}
		_, err = w.Write([]byte{b})
	case NChar:
		_, err = w.Write([]byte{mn.Node.(CharNode).Value})
	case NInt:
		err = encodeInt(w, mn.Node.(IntNode).Value)
	case NFloat:
		err = encodeFloat(w, mn.Node.(FloatNode).Value)
	case NString:
		err = encodeString(w, mn.Node.(StringNode).Value)
	case NStruct:
		sn := mn.Node.(StructNode)
		if err = encodeType(w, sn.Type); err != nil {
			return err
		}
		err = encodeNodeSlice(w, mn.Node.(StructNode).Args)
	case NArray:
		an := mn.Node.(ArrayNode)
		if err = encodeType(w, an.Type.Element); err != nil {
			return err
		}
		err = encodeNodeSlice(w, an.Elems)

	case NIf:
		in := mn.Node.(IfNode)
		if err = encodeBranch(w, in.Main); err != nil {
			return err
		}
		if err = encodeBranchSlice(w, in.Other); err != nil {
			return err
		}
		err = encodeNodeSlice(w, in.Else)
	case NWhile:
		wn := mn.Node.(WhileNode)
		err = encodeBranchOf(w, wn.Cond, wn.Block)
	case NReturn:
		err = encodeNode(w, mn.Node.(ReturnNode).Value)

	case NVarSet:
		vsn := mn.Node.(VarSetNode)
		if err = encodeString(w, vsn.Var); err != nil {
			return err
		}
		err = encodeNode(w, vsn.Value)
	case NVarDef:
		vdn := mn.Node.(VarDefNode)
		if err = encodeString(w, vdn.Var); err != nil {
			return err
		}
		err = encodeType(w, vdn.Type)
	case NVarSetTyped:
		vstn := mn.Node.(VarSetTypedNode)
		if err = encodeString(w, vstn.Var); err != nil {
			return err
		}
		if err = encodeType(w, vstn.Type); err != nil {
			return err
		}
		err = encodeNode(w, vstn.Value)

	case NFuncCall:
		fcn := mn.Node.(FuncCallNode)
		if err = encodeString(w, fcn.Func); err != nil {
			return err
		}
		err = encodeNodeSlice(w, fcn.Args)

	default:
		return pe.New(pe.EUnencodableNode).Section("Caused By", "%s Node at %s", k, mn.Where)
	}
	return
}

func encodeNodeSlice(w io.Writer, ns []MetaNode) error {
	if err := encodeUint(w, uint(len(ns))); err != nil {
		return err
	}
	for _, n := range ns {
		if err := encodeNode(w, n); err != nil {
			return err
		}
	}
	return nil
}

func encodeType(w io.Writer, t types.Type) (err error) {
	p := t.Prim()
	if _, err = w.Write([]byte{byte(p)}); err != nil {
		return err
	}

	switch p {
	case types.PNothing:
		_, err = w.Write([]byte{0})
	case types.PUndefined:
		_, err = w.Write([]byte{1})
	case types.PStruct:
		st := t.(types.StructType)
		if err = encodeString(w, st.Name); err != nil {
			return err
		}
		if err = encodeUint(w, uint(len(st.Fields))); err != nil {
			return err
		}
		err = encodeDescriptorSlice(w, st.Fields)
	case types.PArray:
		err = encodeType(w, t.(types.ArrayType).Element)
	}
	return
}

func encodePos(w io.Writer, p lexer.Position) error {
	if err := encodeUint(w, p.Col); err != nil {
		return err
	}
	if err := encodeUint(w, p.Line); err != nil {
		return err
	}
	return encodeString(w, p.File)
}

func encodeDescriptor(w io.Writer, d types.Descriptor) error {
	if err := encodeString(w, d.Name); err != nil {
		return err
	}
	return encodeType(w, d.Type)
}

func encodeDescriptorSlice(w io.Writer, ds []types.Descriptor) error {
	if err := encodeUint(w, uint(len(ds))); err != nil {
		return err
	}
	for _, d := range ds {
		if err := encodeDescriptor(w, d); err != nil {
			return err
		}
	}
	return nil
}

func encodeBranchOf(w io.Writer, cond MetaNode, b Block) error {
	return encodeBranch(w, Branch{
		Cond:  cond,
		Block: b,
	})
}

func encodeBranch(w io.Writer, b Branch) error {
	if err := encodeNode(w, b.Cond); err != nil {
		return err
	}
	return encodeNodeSlice(w, b.Block)
}

func encodeBranchSlice(w io.Writer, bs []Branch) error {
	if err := encodeUint(w, uint(len(bs))); err != nil {
		return err
	}
	for _, b := range bs {
		if err := encodeBranch(w, b); err != nil {
			return err
		}
	}
	return nil
}

func encodeString(w io.Writer, s string) error {
	if err := encodeUint(w, uint(len(s))); err != nil {
		return err
	}
	_, err := w.Write([]byte(s))
	return err
}

func encodeUint(w io.Writer, n uint) error {
	p := make([]byte, 4)
	binary.BigEndian.PutUint32(p, uint32(n))
	_, err := w.Write(p)
	return err
}

func encodeInt(w io.Writer, n int64) error {
	p := make([]byte, 8)
	binary.BigEndian.PutUint64(p, uint64(n))
	_, err := w.Write(p)
	return err
}

func encodeFloat(w io.Writer, n float64) error {
	p := make([]byte, 8)
	bits := math.Float64bits(n)
	binary.BigEndian.PutUint64(p, bits)
	_, err := w.Write(p)
	return err
}
