# skol

**skol** (for **S**yz**k**rash **O**rdinary **L**anguage) is a minimalist,
keyword-less language for computers.

## Bye-bye!

Skol is no longer in development. If you wish it was, however, feel free to
fork the repository and keep it going.

## Quickstart

Skol is not yet ready to be used in production, therefore you have to build it
manually:

1. Clone and navigate into the repository:

    ```sh
    git clone https://github.com/syzkrash/skol
    cd skol
    ```

2. Run `go build`:

    ```sh
    go build
    ```

3. Done! (yes, it's really that easy)

To compile and run a basic "Hello world!" program using the Python engine:

1. Create the file `hello.sk`:

    ```hs
    $Main(
      print! "Hello world!"
    )
    ```

2. Transpile to Python and run:

    ```sh
    skol compile py hello.sk
    ```

3. Done!

## Learn More

For practical information, visit the [Skol Documentation][doc].
For Skol's full story, visit [qeaml][qeaml]'s [Skol Project Page][project].

[doc]: https://syzkrash.github.io/skol
[project]: https://qeaml.github.io/skol
[qeaml]: https://github.com/qeaml
