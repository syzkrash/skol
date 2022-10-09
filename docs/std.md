# Standard Library

*Note:* Read `*` as "true" and `/` as "false".

*Note:* For "generic" types (e.g. two function parameters of the same type),
`T` is used.

All of these functions are always present in all scopes in your Skol programs.

## Math and logic

* `$add/T a/T b/T`, `$sub/T a/T b/T`, `$mul/T a/T b/T`, `$div/T a/T b/T`,
  `$pow/T a/T b/T`

  Performs the given operation on two numeric values. The "generic" type `T`
  may be one of: `char`, `int`, `float`.

* `$mod/i a/T b/i`

  Returns the remainder of `div! a b`. The "generic" type `T` may be any numeric
  type.

* `$eq/bool a/any b/any`

  Returns `*` if the two values are equal.

* `$gt/bool a/T b/T`, `$lt/bool a/T b/T`

  Compare `a` and `b`. For `gt`, `*` is returned if `a > b`. For `lt`, `*` is
  returned if `b > a`.

* `$not/bool v/bool`, `$and/bool a/bool b/bool`, `$or/bool a/bool b/bool`

  Self-explainatory logical operations.

## Strings and arrays

* `$append/[T] a/[T] b/T`

  Appends the given element to the back of the given array and returns it. The
  "generic" type `T` may be any type. If `a` is a string, `b` may be a char.
  In such case, this function returns a string.

* `$concat/[T] a/[T] b/[T]`

  Concatenates array `b` onto the end of array `a` and returns it. The "generic"
  type `T` may be any type. If `a` and `b` are string, they will be returned as
  a string.

* `$slice/[T] arr/[T] start/i end/i`

  Returns a slice of array `arr` starting at `start` (inclusive) and ending at
  `end` (exclusive). If `end` is below 0, then all values including the last
  value in the array are returned. If `arr` is a string, then a string will be
  returned.

* `$len/int a/[any]`

  Returns the length of an array or a string.

## Conversions

* `$str/string val/any`

  Turns any value into a string. For basic types, this returns their value as
  a string. (eg. `char` with a value of `0x41` will turn to the string `"A"`)
  For strucutres, this will return the name of the structure concatenated with
  it's contained values. (eg. a `CharResult` structure with values `*` and
  `'E'` will return `"CharResult(* 'E')"`)

* `$parse_bool/BoolResult s/str`, `$char/CharResult s/str`,
  `$int/IntResult s/str`, `$float/FloatResult s/str`

  Parses a given type from the given string. This will fail if the given string
  is not a valid literal for that type.

* `$bool/bool v/any`

  Returns `*` if the given value is truthy. Except values equal to 0 and `/`
  itself, every value is truthy. (this includes empty arrays and strings!)
