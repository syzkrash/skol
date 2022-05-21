package parser

type ValueType uint8

const (
	VtUnknown ValueType = iota
	VtInteger
	VtFloat
	VtString
	VtChar
	VtPointer
)

var types = []string{
	"Unknown",
	"Integer",
	"Float",
	"String",
	"Char",
	"Pointer",
}

func (t ValueType) String() string {
	return types[t]
}

func ParseType(raw string) (t ValueType, ok bool) {
	ok = true
	switch raw {
	case "i", "int", "integer":
		t = VtInteger
	case "f", "float", "real":
		t = VtFloat
	case "s", "str", "string":
		t = VtString
	case "c", "char", "rune":
		t = VtChar
	default:
		ok = false
	}
	return
}
