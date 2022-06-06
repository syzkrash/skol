package python

var operators = map[string]string{
	"add_i":  "+",
	"add_f":  "+",
	"sub_i":  "-",
	"sub_f":  "-",
	"mul_i":  "*",
	"mul_f":  "*",
	"div_i":  "//",
	"div_f":  "/",
	"mod_i":  "%",
	"mod_f":  "%",
	"concat": "+",
	"or":     "or",
	"and":    "and",
	"eq":     "==",
	"gt_i":   ">",
	"gt_f":   ">",
	"lt_i":   "<",
	"lt_f":   "<",
}

var renames = map[string]string{
	"to_str":  "str",
	"to_bool": "bool",
}
