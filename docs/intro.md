# An introduction to skol

## Note

Currently, skol is not finished. While the syntax will remain _mostly_ unchanged,
it is important to keep a close eye on the [issues][issues] and
[pull request][pr] pages of the skol repository.

Also note that this introduction is not complete. It only covers a few topics
regarding the language.

## Hello world

The following codeblock shows a Hello World program in skol:

```hs
$Main(print! "Hello world")
Main!
```

That might look a bit confusing, but that's because I overcomplicated it. The
above can be rewritten as:

```hs
print! "Hello world"
```

Let's analyse the first snippet.

1. `$Main(` - This is the function definition syntax.
  The `$` is like the `function` keyword in JavaScript, `func` in Go or `fn` in
  Rust. The following identifier is the name of the function, just like in the
  previously mentioned languages. Now here's the tricky bit -- the parenthesis
  (`(`) does not start the function's argument list. It starts the function's
  body. It's a very small difference, but the parenthesis is the only kind of
  bracket skol uses. There are no curly brackets (`{}`) and most definitely no
  arrow brackets (`<>`). This introduction assumes you are somewhat familiar
  with C, so the definition `$Main()` is equivalent to `void Main() {}`.

2. `print! "Hello world"` -- This is the function call syntax.
  The `!` at the end of the `print` is what allows us to differentiate between
  variable names and function names. Otherwise, it would be very easy to mistake
  variables and functions and cause overall confusion. This call is equivalent
  to `print("Hello world");` in C.

3. `Main!` -- Again, the function call syntax.
  The `Main` function is not called automatically by most skol compilation targets.
  So, like in many scripting languages we have to call it ourselves.

We've only analysed a single snippet of code and already know quite a lot about
the language. Let's move on with another snippet.

## Variables and constants

Consider the following snippet:

```hs
#Hello: "Hello"
%World: "World"

$Hello(
  print! concat! Hello World
)

$Main(
  Hello!
  %World: "There"
  Hello!
)
```

Now we're getting into the real meat and potatoes of this introduction. Let's
get analysing, shall we:

1. `#Hello: "Hello"` -- This is the constant definition syntax.
  Constants work very similarily to variables, except their value never changes.
  Constants are evaluated at compile time and their value is simply put instead
  of a reference. Constants are the only values in skol that are passed by value.
  You can think of this as a `#define` in C.

2. `%World: "World"` -- This is the variable definition syntax.
  It is very similar to the constant definition syntax, only differing by the
  `%` instead of a `#`. The same exact syntax is used both to define a new
  variable, and to change an existing variable's value.

3. `$Hello()` -- Again, function definition.
  It is notable that the variable reference and constant reference syntax are
  the same. We can also note the difference between a constant/variable
  reference and a function call.

[issues]: https://github.com/syzkrash/skol/issues
[pr]: https://github.com/syzkrash/skol/pulls
