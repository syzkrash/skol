# About

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

Check the [Introduction](intro) if you know nothing, or the [Syntax](syntax)
for a full reference. If you wish to know more about skol's inner workings,
check the [Architecture](arch).
