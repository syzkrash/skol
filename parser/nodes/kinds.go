package nodes

import (
	"github.com/syzkrash/skol/lexer"
)

type NodeKind uint8

const (
	NdInteger NodeKind = iota
	NdBoolean
	NdFloat
	NdString
	NdChar
	NdVarDef
	NdFuncCall
	NdFuncDef
	NdFuncExtern
	NdReturn
	NdIf
	NdWhile
	NdStruct
	NdNewStruct
	NdSelector
	NdTypecast
	NdArray
	NdIndex
)

var nodeKinds = []string{
	"Integer",
	"Boolean",
	"Float",
	"String",
	"Char",
	"VarDef",
	"FuncCall",
	"FuncDef",
	"FuncExtern",
	"Return",
	"If",
	"While",
	"Struct",
	"NewStruct",
	"Selector",
	"Typecast",
	"Array",
	"Index",
}

func (k NodeKind) String() string {
	return nodeKinds[k]
}

type Node interface {
	Kind() NodeKind
	Where() lexer.Position
}
