package codegen

import (
	"io"

	"github.com/syzkrash/skol/ast"
)

// Generator represents any abstract code generator. Note that this interface
// is not meant to be implemented on it's own. You should implement ASTGenerator
// or IRGenerator (when that becomes a thing) instead.
type Generator interface {
	// Output sets the output Writer for the next call to Generate. It is OK to
	// write the file header in this call.
	Output(io.Writer)
	// Generate consumes whatever input was given to this Generator before this
	// call. It is OK to panic in case input was not provided before this call.
	Generate() error
}

// ASTGenerator is a generator based on the AST
type ASTGenerator interface {
	Generator // require Generator to also be implemented
	// Input sets the input to the next Generate call to the given AST.
	Input(ast.AST)
}
