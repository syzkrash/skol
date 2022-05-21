# skol Syntax

This is a reference guide for the syntax of the skol language.

## Variable Definition

```rust
%name/string
%name: "value"
```

A variable definition can either use an explicit value or a type, *but not both.*
The first line defines the type of the variable, initializing it to that type's
zero value. The second line defines a variable named `name` with an explicit
`String` value `"value"`. Note that defining the type of the variable is not
required ahead of time. So both the following codeblocks are equivalent:

```rust
%my_var/string
%my_var: "hello"
print! my_var
```

```rust
%my_var: "hello"
print! my_var
```

## Variable Reference

```rust
my_var
```

A reference to a variable is simply the variable's name. Quite ordinary.

## Function Definition

```rust
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

```rust
greet! "John"
```

A function call differs from a variable reference in that it has an additional
`!` at the end (inspired by Rust macros).

When a function is called it's arguments are looked up and values are consumed
from whatever comes after the call until enough argument to satisfy the function
are found. The below example defines and call multiple functions to show this.

```rust
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
