# The Skol Programming Language

skol (for **S**yz**k**rash **O**rdinary **L**anguage) is a minimal, explicit
and otherwise ordinary programming language.

Designed with multi-target compilation in mind, skol can be transpiled to
Python or interpreted directly via Simulation mode.

## Hello World

Here's the traditional Hello World program in skol:

```rust
// define function the main function
$Main
(
  // print to standard output
  print! "Hello World!"
)
// call the main function:
// most compilation targets don't call the main function automatically
Main!
```

## Learn Skol

Check the [Introduction to Skol](intro) if you know nothing, or the [Skol Syntax](syntax)
for a full reference.
