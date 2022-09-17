// Package debug contains functions for togglable debug messages.
//
// The [Log] function will print messages to stderr if it's [Attribute] is
// present in the [GlobalAttrs] bitfield. This function does nothing if Skol
// is compiled without the DEBUG build tag.
package debug
