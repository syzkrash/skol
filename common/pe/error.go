package pe

import "fmt"

type ErrorCode uint16

const (
	ECLI ErrorCode = 100 + iota

	EUnknownAction
	ENoInput
	EBadInput
	EBadAST
	EBadOutput
	EUnknownEngine
	EBadDebugFlag
	EUnimplemented
)

const (
	ELexer ErrorCode = 200 + iota

	EIllegalChar
	EInvalidCharLit
	EInvalidEscape
)

const (
	EParser ErrorCode = 300 + iota

	EBadFloatLit
	EBadIntLit
	EBadSelectorRoot
	EBadSelectorParent
	EBadIndexParent
	EConstantRedefined
	EBadFuncArgType
	EIllegalTopLevelNode

	EUnknownFunction
	EUnknownVariable
	EUnknownType
	EUnknownField

	EExpectedValue
	EExpectedType
	ENeedTypeOrValue
	EExpectedName
	ENeedBodyOrExtern
	EExpectedSelectorElem

	EExpectedLParen
	EExpectedColon
	EExpectedRBrack

	EUnexpectedToken
)

const (
	ETypeCheck ErrorCode = 400 + iota

	ETypeMismatch
	EVarTypeChanged
	ENeedMoreArgs
)

var emsgs = map[ErrorCode]string{
	EUnknownAction: "Unknown action.",
	ENoInput:       "Provide an input file.",
	EBadInput:      "Could not read input file.",
	EBadAST:        "Malformed AST.",
	EBadOutput:     "Could not write output file.",
	EUnknownEngine: "Unknown engine.",
	EBadDebugFlag:  "Unknown debug flag.",
	EUnimplemented: "Unimplemented.",

	EIllegalChar:    "Illegal character.",
	EInvalidCharLit: "Invalid character literal.",
	EInvalidEscape:  "Invalid escape sequence.",

	EBadFloatLit:          "Invalid float literal.",
	EBadIntLit:            "Invalid integer literal.",
	EBadSelectorRoot:      "Selectors must start with a variable name.",
	EBadSelectorParent:    "Can only select fields on structures.",
	EBadIndexParent:       "Can only index arrays.",
	EConstantRedefined:    "Constants cannot be redefined.",
	EBadFuncArgType:       "Illegal function argument type.",
	EIllegalTopLevelNode:  "Illegal top-level node.",
	EUnknownFunction:      "Unknown function.",
	EUnknownVariable:      "Unknown variable.",
	EUnknownType:          "Unknown type.",
	EUnknownField:         "Unknown field.",
	EExpectedValue:        "Expected a value.",
	EExpectedType:         "Expected a type.",
	ENeedTypeOrValue:      "Need an explicit type or value for implicit type.",
	EExpectedName:         "Expected a name.",
	ENeedBodyOrExtern:     "Need function body or '?' for extern.",
	EExpectedSelectorElem: "Expected a selector element.",
	EExpectedLParen:       "Expected '('.",
	EExpectedColon:        "Expected ':'.",
	EExpectedRBrack:       "Expected ']'.",
	EUnexpectedToken:      "Unexpected token.",

	ETypeMismatch:   "Type mismatch.",
	EVarTypeChanged: "Variable type cannot change.",
	ENeedMoreArgs:   "Need more arguments.",
}

type section struct {
	title   string
	message string
}

type PrettyError struct {
	Code     ErrorCode
	cause    error
	sections []section
}

func (p *PrettyError) Error() string {
	return emsgs[p.Code]
}

func (p *PrettyError) Unwrap() error {
	return p.cause
}

func (p *PrettyError) Print() {
	pprintln(fmt.Sprintf("Error %d", p.Code), gBold, cFgRed)
	pprintln("  " + emsgs[p.Code])
	for _, s := range p.sections {
		pprintln(s.title, gBold)
		pprintln("  " + s.message)
	}
}

func (p *PrettyError) Section(title, mfmt string, margs ...any) *PrettyError {
	p.sections = append(p.sections, section{
		title:   title,
		message: fmt.Sprintf(mfmt, margs...),
	})
	return p
}

func (p *PrettyError) Cause(e error) *PrettyError {
	p.cause = e
	return p.Section("Cause", e.Error())
}

func New(c ErrorCode) *PrettyError {
	return &PrettyError{
		Code:     c,
		sections: []section{},
	}
}
