package defaults

import "github.com/syzkrash/skol/parser/values/types"

var FileDescriptor = types.MakeStruct("file_descriptor",
	"id", types.Int,
	"fn", types.String)
var FileDescriptorResult = types.MakeStruct("file_descriptor_result",
	"fd", FileDescriptor,
	"ok", types.Bool,
	"err", types.String)

var CharResult = types.MakeStruct("char_result",
	"char", types.Char,
	"ok", types.Bool,
	"err", types.String)

// Types contains definitions of built-in types for the parser.
var Types = map[string]types.Type{
	"file_descriptor":        FileDescriptor,
	"file_descriptor_result": FileDescriptorResult,
	"char_result":            CharResult,
}
