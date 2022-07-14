package defaults

import "github.com/syzkrash/skol/parser/values"

var file_descriptor = values.MakeStruct("file_descriptor",
	"id", values.Int,
	"fn", values.String)
var file_descriptor_result = values.MakeStruct("file_descriptor_result",
	"fd", file_descriptor,
	"ok", values.Bool,
	"err", values.String)

var char_result = values.MakeStruct("char_result",
	"char", values.Char,
	"ok", values.Bool,
	"err", values.String)

var Types = map[string]*values.Type{
	"file_descriptor":        file_descriptor,
	"file_descriptor_result": file_descriptor_result,
	"char_result":            char_result,
}
