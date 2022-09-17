// Package parser defines the Skol parser.
//
// The [Parser] consumes [lexer.Token]s and produces [ast.Node]s out of them,
// wrapped in an [ast.MetaNode]. This is the most complicated part of Skol.
package parser
