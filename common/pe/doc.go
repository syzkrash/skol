// Package pe defines the pretty errors Skol uses.
//
// A [PrettyError] consists of an [ErrorCode] and a slice of [section]s giving
// more detailed information about the error. The error does not contain the
// error message itself. The message is retrieved from the [emsgs] map instead.
package pe
