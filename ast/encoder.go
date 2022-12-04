package ast

import (
	"io"

	"github.com/syzkrash/skol/common/pack"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

// Encode writes a binary representation of the given AST into the provided
// [io.Writer]
func Encode(w io.Writer, tree AST) (err error) {
	pk := pack.NewPacker(w)

	pk.Write([]byte(FormatMagic)).U8(FormatVersion)

	pk.U8(uint8(len(tree.Vars)))
	for _, v := range tree.Vars {
		encodeVar(pk, v)
	}

	pk.U8(uint8(len(tree.Typedefs)))
	for _, v := range tree.Typedefs {
		encodeTypedef(pk, v)
	}

	pk.U8(uint8(len(tree.Funcs)))
	for _, f := range tree.Funcs {
		encodeFunc(pk, f)
	}

	pk.U8(uint8(len(tree.Exerns)))
	for _, e := range tree.Exerns {
		encodeExtern(pk, e)
	}

	pk.U8(uint8(len(tree.Structs)))
	for _, s := range tree.Structs {
		encodeStruct(pk, s)
	}

	if len(pk.Err) > 0 {
		return pk.Err[0]
	}

	return
}

func encodeVar(pk *pack.Packer, v Var) {
	pk.Str(v.Name)
	encodeNode(pk, v.Value)
}

func encodeTypedef(pk *pack.Packer, t Typedef) {
	pk.Str(t.Name)
	encodeType(pk, t.Type)
}

func encodeFunc(pk *pack.Packer, f Func) {
	pk.Str(f.Name)
	encodeType(pk, f.Ret)
	encodeDescriptorSlice(pk, f.Args)
	encodeNodeSlice(pk, f.Body)
}

func encodeExtern(pk *pack.Packer, e Extern) {
	pk.Str(e.Alias)
	pk.Str(e.Name)
	encodeType(pk, e.Ret)
	encodeDescriptorSlice(pk, e.Args)
}

func encodeStruct(pk *pack.Packer, s Structure) {
	pk.Str(s.Name)
	encodeDescriptorSlice(pk, s.Fields)
}

func encodeNode(pk *pack.Packer, mn MetaNode) {
	encodePos(pk, mn.Where)
	k := mn.Node.Kind()
	pk.U8(uint8(k))

	switch k {
	case NBool:
		var b byte
		if mn.Node.(BoolNode).Value {
			b = 1
		}
		pk.U8(b)
	case NChar:
		pk.U8(mn.Node.(CharNode).Value)
	case NInt:
		pk.I64(mn.Node.(IntNode).Value)
	case NFloat:
		pk.F64(mn.Node.(FloatNode).Value)
	case NString:
		pk.Str(mn.Node.(StringNode).Value)
	case NStruct:
		sn := mn.Node.(StructNode)
		encodeType(pk, sn.Type)
		encodeNodeSlice(pk, sn.Args)
	case NArray:
		an := mn.Node.(ArrayNode)
		encodeType(pk, an.Type.Element)
		encodeNodeSlice(pk, an.Elems)

	case NIf:
		in := mn.Node.(IfNode)
		encodeBranch(pk, in.Main)
		encodeBranchSlice(pk, in.Other)
		encodeNodeSlice(pk, in.Else)
	case NWhile:
		wn := mn.Node.(WhileNode)
		encodeBranchOf(pk, wn.Cond, wn.Block)
	case NReturn:
		encodeNode(pk, mn.Node.(ReturnNode).Value)

	case NVarSet:
		vsn := mn.Node.(VarSetNode)
		pk.Str(vsn.Var)
		encodeNode(pk, vsn.Value)
	case NVarDef:
		vdn := mn.Node.(VarDefNode)
		pk.Str(vdn.Var)
		encodeType(pk, vdn.Type)
	case NVarSetTyped:
		vstn := mn.Node.(VarSetTypedNode)
		pk.Str(vstn.Var)
		encodeType(pk, vstn.Type)
		encodeNode(pk, vstn.Value)

	case NFuncCall:
		fcn := mn.Node.(FuncCallNode)
		pk.Str(fcn.Func)
		encodeNodeSlice(pk, fcn.Args)

	default:
		pk.Error(pe.New(pe.EUnencodableNode).Section("Caused By", "%s Node at %s", k, mn.Where))
	}
	return
}

func encodeNodeSlice(pk *pack.Packer, ns []MetaNode) {
	pk.U8(uint8(len(ns)))
	for _, n := range ns {
		encodeNode(pk, n)
	}
}

func encodeType(pk *pack.Packer, t types.Type) {
	p := t.Prim()
	pk.U8(uint8(p))

	switch p {
	case types.PStruct:
		st := t.(types.StructType)
		pk.Str(st.Name)
		encodeDescriptorSlice(pk, st.Fields)
	case types.PArray:
		encodeType(pk, t.(types.ArrayType).Element)
	}
	return
}

func encodePos(pk *pack.Packer, p lexer.Position) {
	pk.U32(uint32(p.Col)).U32(uint32(p.Line)).Str(p.File)
}

func encodeDescriptor(pk *pack.Packer, d types.Descriptor) {
	pk.Str(d.Name)
	encodeType(pk, d.Type)
}

func encodeDescriptorSlice(pk *pack.Packer, ds []types.Descriptor) {
	pk.U8(uint8(len(ds)))
	for _, d := range ds {
		encodeDescriptor(pk, d)
	}
}

func encodeBranchOf(pk *pack.Packer, cond MetaNode, b Block) {
	encodeBranch(pk, Branch{
		Cond:  cond,
		Block: b,
	})
}

func encodeBranch(pk *pack.Packer, b Branch) {
	encodeNode(pk, b.Cond)
	encodeNodeSlice(pk, b.Block)
}

func encodeBranchSlice(pk *pack.Packer, bs []Branch) {
	pk.U8(uint8(len(bs)))
	for _, b := range bs {
		encodeBranch(pk, b)
	}
}
