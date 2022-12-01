// Package pack provides a simple interface for writing Go values into byte
// streams and reading Go values from byte streams.
//
// This is a simple mostly zero-allocation implementation. Allocations are only
// made by string-related functions.
package pack
