package defaults

import "github.com/syzkrash/skol/parser/values/types"

var file_descriptor = types.MakeStruct("file_descriptor",
	"id", types.Int,
	"fn", types.String)
var file_descriptor_result = types.MakeStruct("file_descriptor_result",
	"fd", file_descriptor,
	"ok", types.Bool,
	"err", types.String)

var char_result = types.MakeStruct("char_result",
	"char", types.Char,
	"ok", types.Bool,
	"err", types.String)

var Types = map[string]types.Type{
	"file_descriptor":        file_descriptor,
	"file_descriptor_result": file_descriptor_result,
	"char_result":            char_result,
}
