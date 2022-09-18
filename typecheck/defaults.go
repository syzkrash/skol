package typecheck

import (
	"github.com/syzkrash/skol/parser/defaults"
	"github.com/syzkrash/skol/parser/values/types"
)

func makeFuncproto(ret types.Type, args ...types.Type) funcproto {
	return basicFuncproto(args, ret)
}

var defaultFuncs = map[string]funcproto{
	"print":       makeFuncproto(types.Nothing, types.String),
	"to_str":      makeFuncproto(types.String, types.Any),
	"to_bool":     makeFuncproto(types.Bool, types.Any),
	"add_i":       makeFuncproto(types.Int, types.Int, types.Int),
	"add_f":       makeFuncproto(types.Float, types.Float, types.Float),
	"add_c":       makeFuncproto(types.Char, types.Char, types.Char),
	"sub_i":       makeFuncproto(types.Int, types.Int, types.Int),
	"sub_f":       makeFuncproto(types.Float, types.Float, types.Float),
	"sub_c":       makeFuncproto(types.Char, types.Char, types.Char),
	"mul_i":       makeFuncproto(types.Int, types.Int, types.Int),
	"mul_f":       makeFuncproto(types.Float, types.Float, types.Float),
	"div_i":       makeFuncproto(types.Int, types.Int, types.Int),
	"div_f":       makeFuncproto(types.Float, types.Float, types.Float),
	"mod_i":       makeFuncproto(types.Int, types.Int, types.Int),
	"mod_f":       makeFuncproto(types.Float, types.Float, types.Int),
	"concat":      makeFuncproto(types.String, types.String, types.String),
	"not":         makeFuncproto(types.Bool, types.Bool),
	"or":          makeFuncproto(types.Bool, types.Bool, types.Bool),
	"and":         makeFuncproto(types.Bool, types.Bool, types.Bool),
	"eq":          makeFuncproto(types.Bool, types.Any, types.Any),
	"gt_i":        makeFuncproto(types.Bool, types.Int, types.Int),
	"gt_f":        makeFuncproto(types.Bool, types.Float, types.Float),
	"gt_c":        makeFuncproto(types.Bool, types.Char, types.Char),
	"lt_i":        makeFuncproto(types.Bool, types.Int, types.Int),
	"lt_f":        makeFuncproto(types.Bool, types.Float, types.Float),
	"lt_c":        makeFuncproto(types.Bool, types.Char, types.Char),
	"char_at":     makeFuncproto(types.Char, types.String, types.Int),
	"substr":      makeFuncproto(types.String, types.String, types.Int, types.Int),
	"char_append": makeFuncproto(types.String, types.String, types.Char),
	"str_len":     makeFuncproto(types.Int, types.String),
	"skol":        makeFuncproto(types.Nothing, types.String, types.Float),
	"ctoi":        makeFuncproto(types.Int, types.Char),
	"open":        makeFuncproto(defaults.FileDescriptorResult, types.String),
	"fgetc":       makeFuncproto(defaults.CharResult, defaults.FileDescriptor),
	"close":       makeFuncproto(types.Nothing, defaults.FileDescriptor),
}
