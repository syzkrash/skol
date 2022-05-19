package parser

type ValueType uint8

const (
	VtInteger ValueType = iota
	VtFloat
	VtString
	VtChar
	VtPointer ValueType = 1 << 7
)

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
