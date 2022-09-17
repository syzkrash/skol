// Package codegen defines the engines Skol used to compile/transpile/execute
// code.
//
// An [Engine] consists of any combination of a code [Generator] and an
// [Executor].
//
// Tthere are two types of generator: an [ASTGenerator], generating
// it's output directly from the AST and an IRGenerator, generating from a
// simplified IR (note that the IR is not yet implemented). An engine may
// want to use it's own IR, rather than the generic one, so an ASTGenerator
// is not exclusive to transpilers.
//
// Executors are also split into two types: an [EphemeralExecutor], which
// executes code directly from memory, and a [FilenameExecutor], which executes
// code from a file given the file's name.
package codegen
