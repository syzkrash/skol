// Package ast defines the AST itself along with all of the nodes it contains.
//
// The main data structures of this package are the [AST], the [MetaNode] and
// [Node]s.
package ast

// FormatMagic is the magic string of the AST file format
const FormatMagic = "SKAST"

// FormatVersion is the version ordinal of the AST file format
const FormatVersion byte = 1
