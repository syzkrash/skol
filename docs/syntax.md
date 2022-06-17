# skol Syntax

This is a reference guide for the syntax of the skol language.

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
```

A variable definition can either use an explicit value or a type, *but not both.*
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

## Conclusion

If you wish to see skol in action, feel free to view the [JSON][json] and
[calculator][calc] examples.

[json]: /examples/json.skol
[calc]: /examples/calculator.skol
