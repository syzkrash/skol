// Package typechecker defines a type checker for Skol code.
//
// The [Checker] ensures that an [ast.AST] has no type errors. Such type errors
// may include:
//   - Type mismatch: expected type X, got type Y
//   - Illegal type: a variable cannot have type Any, Nothing, a function
//     argument cannot have type Any, etc.
//   - Variable retype: a variable cannot change type after it has been defined.
package typecheck
