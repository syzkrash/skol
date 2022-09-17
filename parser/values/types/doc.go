// Package types defines and handles the types in the Skol language.
//
// Every [Type] has a [Primitive]. [PrimitiveType]s, like integers or strings
// are only distinguished by their primitive.
//
// [StructType]s always have the [PStruct] primitive and are distinguished by
// the [Field]s they contain.
//
// [ArrayType]s always have the [PArray] primitive and are distinguished by
// the type of their elements.
//
// The [AnyType] always has the [PAny] primitive and is compatible with any
// other type. Because of this it is only allowed under strict conditions.
// (e.g. in builtin functions and external functions)
//
// The [NothingType] can have one of two primitives: [PNothing] or [PUndefined].
// NothingTypes are incompatible with all other types. The main difference
// between PNothing and PUndefined is that PNothing is used for functions that
// do not return anything. PUndefined is used internally as a placeholder.
package types
