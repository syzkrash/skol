package debug

// Attribute defines a bitmask for debug message visibility.
type Attribute uint8

// Debug message attribute constnts.
const (
	AttrDummy Attribute = 1 << iota
	AttrLexer
	AttrParser
	AttrScope
)

func (a Attribute) Name() string {
	switch a {
	case AttrLexer:
		return "Lexer"
	case AttrParser:
		return "Parser"
	case AttrScope:
		return "Scope"
	default:
		return "Debug"
	}
}

// GlobalAttr contains the currently enabled debug messages.
var GlobalAttr = AttrDummy
