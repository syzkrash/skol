# Syntax

## Comments

skol uses C-like comments:

```c
// this is a line comment.

/* this is a block comment.
   block comments cannot be nested. */
```

## Variable Definition

```hs
%name/string
%name: "value"
%name/string: "value"
```

A variable definition can either use an explicit value or a type, or both.
The first line defines the type of the variable, initializing it to that type's
zero value. The second line defines a variable named `name` with an explicit
`String` value `"value"`. Note that defining the type of the variable is not
required ahead of time. So both the following codeblocks are equivalent:

```hs
%my_var/string
%my_var: "hello"
print! my_var
```

```hs
%my_var: "hello"
print! my_var
```

## Variable Reference

```hs
my_var
```

A reference to a variable is simply the variable's name. Quite ordinary.

## Function Definition

```hs
$greet who/s (
  %greeting: concat! "Hello, " who
  print! greeting
  >greeting
)
```

A function definition is quite simple. Simply use the `$` punctuator, followed
by the function name and the fuction arguments, and top it off with the function
body. The return value will be automatically inferred from the function body.

The function in the above example concatenates the string `"Hello, "` and the
provided `String` argument `who`, into the `greeting` variable. It then `print`s
the `greeting` variable and returns it.

## Function Call

```hs
greet! "John"
```

A function call differs from a variable reference in that it has an additional
`!` at the end (inspired by Rust macros).

When a function is called it's arguments are looked up and values are consumed
from whatever comes after the call until enough argument to satisfy the function
are found. The below example defines and call multiple functions to show this.

```hs
$greet who/s (
  >concat! "Hello, " who
)
print! greet! "Joe"

$exclaim greeting/s who/s (
  >concat! concat! greeting ", " who
)
$hello (
  >"Hello"
)
print! exclaim! hello! "Joe"
```

## Conditional

```hs
? SomeCondition! (
  print! "Condition #1 met."
) :? OtherCondition! (
  print! "Condition #1 not met, but condition #2 met."
) : (
  print! "Neither condition was met."
)
```

Skol now features full-blown control flow and branching!

An if statement is written just like any other language, except you use `?` for
`if`, `:?` for `else if` and `:` for `else`.

For example, to determine the relationship of two numbers:

```hs
$ NumCompare A/f B/f (
  ?  gtr! A B
    (> 1)
  ?: gtr! B A
    (> -1)
  >0
)
```

The above function returns `1` is `A` is greater than `B`, `-1` if `B` is grater
than `A`, and `0` if the numbers are equal.

## Loop

```hs
$fact/i n/i
(
  %acc: 1
  %i: 0
  *lt_i! i n
  (
    %i: add_i! i 1
    %acc: mul_i! acc i
  )
  >acc
)

#ten: 10
print!
  concat! "The factorial of 10 is: "
  to_str! fact! ten
```

In skol, there only exists a `while` loop in the form of `*`. As long as the
condition after `*` evaluates to `true` (or `*` in actual skol code), the code
inside the block after it will be repeated. The example above shows not only the
`while` loop, but also function definitions and calls, variable definitions and
references as well as a constant definition.

## Literals

The literals are quite similar to other languages. Here's a quick rundown:

* `'a'` is a __character__ literal, __not__ a string.
* `"hello"` is a string literal.
* `123`, `12.3` and `0xD34D` are all numeric literals.
* `*` is the boolean `true` and `/` is `false`.

Note that skol does not use `true` and `false` for boolean literals. Use `*`
and `/` instead.

A cool side note is that due to the boolean literal syntax and loop syntax
using the same punctuator, an infinite loop is simply defined with `**`:

```hs
**
(
  print! "You have been hax!"
)
```

## Structure

```hs
@Vec2i(
  x/int
  y/int
)

$AddVec2i/Vec2i a/Vec2i b/Vec2i(
  >@Vect2i add_i! a#x b#x add_i! a#y b#y
)
```

Like C, skol does not have classes or any OOP concepts. Instead, everything is
done with structures. Naturally, structures have fields of certain types.
Polymorphism is possible in skol due to loose typechecks:

```hs
// considering the Vec2i and AddVec2i from above

@Vec3i(
  x/int
  y/int
  z/int
)

$Main(
  %my_vec: @Vec3i 1 2 3
  %other_vec: @Vec2i 2 1
  %final_vec: AddVec2i! my_vec other_vec
)
```

The above code will not produce a type error. Type checks in skol only consider
whether all the required fields are present in a type, not checking the identity
of the type itself. That means: a Vec3i can act as a Vec2i, as it contains all
the fields Vec2i contains.

## Typecast

```hs
// again, considering Vec2i from above

@Vec2or3(
  is_vec3/bool
  x/int
  y/int
)

@Vec2or3Result(
  ok/bool
  vec/Vec2or3
)

// Vec3i doesn't have the is_vec3 field
@Vec3O(
  is_vec3/bool
  x/int
  y/int
  z/int
)

$AddVec/Vec2or3Result a/Vec2or3 b/Vec2or3(
  // fail if a and b are not the same size
  ?not! eq! a#is_vec3 b#is_vec3(
    >@Vec2or3Result / @Vec2or3 / 0 0
  )
  // call the appropriate function based on the type
  ?a#is_vec3(
    %v2r: AddVec2O! a b
    >@Vec2or3Result * v2r
  ):(
    // convert to Vec3 type: possible, since Vec3O contains all the fields
    // Vec2or3 contais
    %v3a: a#@Vec3O
    %v3b: b#@Vec3O
    // add Vec3
    %v3r: AddVec3O! v3a v3b
    // convert back to Vec2or3
    %v3o: v3r#@Vec2or3
    // and return the result
    >@Vec2or3Result * v3o
  )
)
```

Typecasts are similar to selectors, which use the `#` punctuator. Instead of
using an identifier to denote a field of a structure, use `@` to denote a
typecast, followed by the name of the type to cast to.

**Note** that primitive types, like booleans, integers, floats and strings are
*not* castable. Typecasts can only be performed on structures.

## Array

Well, a better word would be 'list' or 'slice', but if JavaScript can call
these 'array's so can we!

```hs
%CoolPeople: [str]("Joe" "Jah" "Jen")
%UncoolPeople:  []("Bob" "Ben" "Bil")
%EmptyTyped: [int]()
```

An array literal is very simple: type type of the array's elements surrounded
by square brackets, followed by a list of values of that type. Said values are
not limited to just literals. Function calls, constants, variable references are
allowed. The type of the elements can be omitted to allow skol to deduce the
element type. If the array doesn't have an explicit type declaration and doesn't
contain any elements, skol cannot determine what type it's supposed to be and
throws an error.

## Conclusion

If you wish to see skol in action, feel free to view the [JSON][json] and
[calculator][calc] examples. You can also see the more recent, but simpler
[document parser][doc] example.

[json]: /examples/json.skol
[calc]: /examples/calculator.skol
[doc]:  /examples/Doc.skol
