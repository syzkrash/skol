package ast

import (
	"io"

	"github.com/syzkrash/skol/common/pack"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser/values/types"
)

// Decode reads a binary representation of an AST and returns it. This function
// assumes the input data begins with the [FormatMagic] string, followed by a
// one-byte version of the format. See [FormatVersion] for the current version.
func Decode(src io.Reader) (tree AST, err error) {
	u := pack.NewUnpacker(src)

	magic := u.Bytes(uint(len(FormatMagic)))
	if string(magic) != FormatMagic {
		err = pe.New(pe.EBadMagic)
		return
	}
	ver := u.U8()
	if ver != FormatVersion {
		err = pe.New(pe.EBadEncoderVer).Section("Caused By", "%02X at $%08X", ver, u.Offset-1)
		return
	}

	tree = NewAST()

	count := u.U8()
	for i := uint8(0); i < count; i++ {
		v := decodeVar(u)
		tree.Vars[v.Name] = v
	}

	count = u.U8()
	for i := uint8(0); i < count; i++ {
		t := decodeTypedef(u)
		tree.Typedefs[t.Name] = t
	}

	count = u.U8()
	for i := uint8(0); i < count; i++ {
		f := decodeFunc(u)
		tree.Funcs[f.Name] = f
	}

	count = u.U8()
	for i := uint8(0); i < count; i++ {
		e := decodeExtern(u)
		tree.Exerns[e.Alias] = e
	}

	count = u.U8()
	for i := uint8(0); i < count; i++ {
		s := decodeStruct(u)
		tree.Structs[s.Name] = s
	}

	if len(u.Err) > 0 {
		err = u.Err[0]
	}

	return
}

func decodeVar(u *pack.Unpacker) (v Var) {
	v.Name = u.Str()
	v.Value = decodeNode(u)
	return
}

func decodeTypedef(u *pack.Unpacker) (t Typedef) {
	t.Name = u.Str()
	t.Type = decodeType(u)
	return
}

func decodeFunc(u *pack.Unpacker) (f Func) {
	f.Name = u.Str()
	f.Ret = decodeType(u)
	f.Args = decodeDescriptorSlice(u)
	f.Body = decodeNodeSlice(u)
	return
}

func decodeExtern(u *pack.Unpacker) (e Extern) {
	e.Alias = u.Str()
	e.Name = u.Str()
	e.Ret = decodeType(u)
	e.Args = decodeDescriptorSlice(u)
	return
}

func decodeStruct(u *pack.Unpacker) (s Structure) {
	s.Name = u.Str()
	s.Fields = decodeDescriptorSlice(u)
	return
}

func decodeNode(u *pack.Unpacker) (mn MetaNode) {
	mn.Where = decodePos(u)
	k := NodeKind(u.U8())

	switch k {
	case NBool:
		mn.Node = BoolNode{
			Value: u.U8() > 0,
		}
	case NChar:
		mn.Node = CharNode{
			Value: u.U8(),
		}
	case NInt:
		mn.Node = IntNode{
			Value: u.I64(),
		}
	case NFloat:
		mn.Node = FloatNode{
			Value: u.F64(),
		}
	case NString:
		mn.Node = StringNode{
			Value: u.Str(),
		}
	case NArray:
		t := decodeType(u)
		e := decodeNodeSlice(u)
		mn.Node = ArrayNode{
			Type: types.ArrayType{
				Element: t,
			},
			Elems: e,
		}

	case NIf:
		m := decodeBranch(u)
		o := decodeBranchSlice(u)
		e := decodeNodeSlice(u)
		mn.Node = IfNode{
			Main:  m,
			Other: o,
			Else:  e,
		}
	case NWhile:
		c := decodeNode(u)
		b := decodeNodeSlice(u)
		mn.Node = WhileNode{
			Cond:  c,
			Block: b,
		}
	case NReturn:
		v := decodeNode(u)
		mn.Node = ReturnNode{
			Value: v,
		}

	case NVarSet:
		n := u.Str()
		v := decodeNode(u)
		mn.Node = VarSetNode{
			Var:   n,
			Value: v,
		}
	case NVarDef:
		n := u.Str()
		t := decodeType(u)
		mn.Node = VarDefNode{
			Var:  n,
			Type: t,
		}
	case NVarSetTyped:
		n := u.Str()
		t := decodeType(u)
		v := decodeNode(u)
		mn.Node = VarSetTypedNode{
			Var:   n,
			Type:  t,
			Value: v,
		}

	case NFuncCall:
		n := u.Str()
		a := decodeNodeSlice(u)
		mn.Node = FuncCallNode{
			Func: n,
			Args: a,
		}

	default:
		u.Error(pe.New(pe.EBadNodeKind).Section("Caused By", "%02X at $%08X", k, u.Offset-1))
	}

	return
}

func decodeNodeSlice(u *pack.Unpacker) (mns []MetaNode) {
	count := u.U8()
	mns = make([]MetaNode, count)
	for i := uint8(0); i < count; i++ {
		mns[i] = decodeNode(u)
	}
	return
}

func decodeType(u *pack.Unpacker) (t types.Type) {
	p := types.Primitive(u.U8())

	switch p {
	case types.PBool:
		t = types.Bool
	case types.PChar:
		t = types.Char
	case types.PInt:
		t = types.Int
	case types.PFloat:
		t = types.Float
	case types.PString:
		t = types.String
	case types.PStruct:
		n := u.Str()
		f := decodeDescriptorSlice(u)
		t = types.StructType{
			Name:   n,
			Fields: f,
		}
	case types.PArray:
		t = types.ArrayType{
			Element: decodeType(u),
		}
	case types.PAny:
		t = types.Any
	case types.PNothing:
		t = types.Nothing
	case types.PUndefined:
		t = types.Undefined

	default:
		u.Error(pe.New(pe.EBadTypePrim).Section("Caused By", "%02X at $%08X", p, -1))
	}

	return
}

func decodeDescriptor(u *pack.Unpacker) (d types.Descriptor) {
	d.Name = u.Str()
	d.Type = decodeType(u)
	return
}

func decodeDescriptorSlice(u *pack.Unpacker) (ds []types.Descriptor) {
	count := u.U8()
	ds = make([]types.Descriptor, count)
	for i := uint8(0); i < count; i++ {
		ds[i] = decodeDescriptor(u)
	}
	return
}

func decodeBranch(u *pack.Unpacker) (b Branch) {
	b.Cond = decodeNode(u)
	b.Block = decodeNodeSlice(u)
	return
}

func decodeBranchSlice(u *pack.Unpacker) (bs []Branch) {
	count := u.U8()
	bs = make([]Branch, count)
	for i := uint8(0); i < count; i++ {
		bs[i] = decodeBranch(u)
	}
	return
}

func decodePos(u *pack.Unpacker) (p lexer.Position) {
	p.Col = uint(u.U32())
	p.Line = uint(u.U32())
	p.File = u.Str()
	return
}
