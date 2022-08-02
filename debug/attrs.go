package debug

type Attribute uint8

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

var GlobalAttr = AttrDummy
