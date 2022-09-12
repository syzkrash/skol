# Architecture

## Component Breakdown

Component   | Completeness               | Source Code                      | Purpose
------------|----------------------------|----------------------------------|---------
CLI         | [Incomplete](#cli)         | [`cli` paclage][cli]             | Parses command line-flags and prepares files for compilation or analysis.
Lexer       | [Complete](#lexer)         | [`lexer` package][lexer]         | Breaks meaningless plain text into a sequence of meaningful tokens.
Parser      | [Incomplete](#parser)      | [`parser` package][parser]       | Consumes tokens from the lexer and constructs an [AST][astw].
AST         | [Incomplete](#ast)         | [`ast` package][ast]             | Represents the structure of a source code file.
Typechecker | [Incomplete](#typechecker) | [`typecheck` package][typecheck] | Ensures that everything in the AST has the correct typing. (ensures a int variable isn't set to a string, etc.)
[IR][irw]   | [Incomplete](#ir)          | N/A                              | Breaks a program down into simple instructions to allow for easier compilation to binary formats.
Codegen     | [Incomplete](#codegen)     | [`codegen` package][codegen]     | a) Transpiles Skol code into another language from the AST. <br/> b) Compiles into executables from IR.

## Usual flow for compilation

1. The CLI reads the source code.
2. A parser is created, which contains it's own lexer for the previously read
   source code.
3. The parser consumes tokens, creates the adequate nodes and constructs an AST
   out of them.
4. The typechecker ensures type correctness in the program.

## Component Completeness Breakdown

### CLI

- [x] Is able to parse actions and arguments separately.
- [x] Is able to build a file using any engine.
- [ ] Is able to start a REPL using any engine.
- [ ] Has a way to access additional tools (e.g. linter).

### Lexer

- [x] Ignores comments.
- [x] Reads punctuators.
- [x] Reads identifiers.
- [x] Reads integer literals.
- [x] Reads float literals.
- [x] Reads character literals.
- [x] Reads string literals.

### Parser

- [x] Parses basic constructs:
   * [x] Variable definition/assignment and both.
   * [x] Function/extern definition.
   * [x] Basic control flow.
   * [x] Structured types.
   * [x] Array types.
- [x] Properly handles expected lexer errors (e.g. EOF).
- [ ] 100% proven waterproof.

### AST

- [x] Is properly constructed by the parser.
- [ ] Can correctly reflect the structure of a Skol program.
   * [x] Global variables
   * [x] Global functons/externs
   * [x] Global types
   * [ ] Multi-file compilation.
   * [ ] Top-level code.

### Typechecker

- [x] Can check variables.
- [x] Can check functions.
- [x] Can check structure types.
- [x] Can check array types.
- [ ] Can determine value types.
- [ ] Supports built-in functions.
- [ ] Supports generic functions.

### IR

- [ ] Can be constructed from any valid AST.
- [ ] Can be cached as a file.

### Codegen

- [ ] Distingushes between AST-based codegen (for transpilers) and IR-based
      codegen (for compilers).
- [ ] Transpiles valid code from AST.
- [ ] Compiles valid code from IR.

[cli]: https://github.com/syzkrash/skol/tree/nightly/cli
[lexer]: https://github.com/syzkrash/skol/tree/nightly/lexer
[parser]: https://github.com/syzkrash/skol/tree/nightly/parser
[ast]: https://github.com/syzkrash/skol/tree/nightly/ast
[sim]: https://github.com/syzkrash/skol/tree/nightly/sim
[typecheck]: https://github.com/syzkrash/skol/tree/nightly/typecheck
[codegen]: https://github.com/syzkrash/skol/tree/nightly/codegen

[astw]: https://en.wikipedia.org/wiki/Abstract_syntax_tree
[irw]: https://en.wikipedia.org/wiki/Intermediate_representation
