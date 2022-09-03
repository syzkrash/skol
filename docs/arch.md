# Architecture

## Component Breakdown

Component   | Completeness | Source Code                      | Purpose
------------|--------------|----------------------------------|---------
CLI         | Incomplete   | [`main.go`][main]                | Parses command line-flags and prepares files for compilation or analysis.
Lexer       | Complete     | [`lexer` package][lexer]         | Breaks meaningless plain text into a sequence of meaningful tokens.
Parser      | Incomplete   | [`parser` package][parser]       | Consumes tokens from the lexer and constructs an [AST][wiki_ast].
AST         | Incomplete   | [`ast` package][ast]             | Represents the structure of a source code file.
Simulator   | Incomplete   | [`sim` package][sim]             | Simulates code from AST nodes. (not an interpreter; used for constant evaluation)
Typechecker | Incomplete   | [`typecheck` package][typecheck] | Ensures that everything in the AST has the correct typing. (ensures a int variable isn't set to a string, etc.)

## Usual flow for compilation

1. The CLI reads the source code.
2. A parser is created, which contains it's own lexer for the previously read
   source code.
3. The parser consumes tokens, creates the adequate nodes and constructs an AST
   out of them.
4. The typechecker ensures type correctness in the program.

[main]: https://github.com/syzkrash/skol/blob/nightly/main.go
[lexer]: https://github.com/syzkrash/skol/tree/nightly/lexer
[parser]: https://github.com/syzkrash/skol/tree/nightly/parser
[ast]: https://github.com/syzkrash/skol/tree/nightly/ast
[sim]: https://github.com/syzkrash/skol/tree/nightly/sim
[typecheck]: https://github.com/syzkrash/skol/tree/nightly/typecheck

[wiki_ast]: https://en.wikipedia.org/wiki/Abstract_syntax_tree
