package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/typecheck"
)

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	flags := flag.NewFlagSet("skol", flag.ExitOnError)
	flags.StringVar(&input, "input", "", "")
	flags.Func("debug", "", debugFlag)
	flags.BoolVar(&asJson, "json", false, "")
	flags.BoolVar(&prettyJson, "pretty-json", false, "")

	flags.Parse(os.Args[2:])

	fmt.Fprintf(os.Stderr, "Skol v%s\n", common.Version)

	act := os.Args[1]

	switch act {
	case "ast":
		dumpAst()
	case "build":
		compile()
	case "repl":
		repl()
	case "":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "I don't know what '%s' means!\n", act)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <action> [arguments...]\n", os.Args[0])
	fmt.Println("Where action can be one of the following:")
	fmt.Println("  build  - Compile/run a file using a given engine")
	fmt.Println("  repl   - Start the skol Read And Evaluate Loop (REPL)")
	fmt.Println("  ast    - Print the Abstract Syntax Tree (AST) of a file")
	fmt.Println("Where arguments can be any of the following:")
	fmt.Println("  -input <filename>  - File to use as input")
	fmt.Println("  -engine <engine>   - Engine to compile/run with")
	fmt.Println("  -debug <kinds>     - Enable specified kinds of debug messages")
	fmt.Println("  -json              - For applicable actions, output as")
}

var (
	input      string
	asJson     bool
	prettyJson bool
)

func dumpAst() {
	if input == "" {
		fmt.Fprintln(os.Stderr, "Specify an input file.")
		return
	}

	f, err := os.Open(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening input file - %s\n", err)
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file - %s\n", err)
		return
	}

	src := bytes.NewReader(data)
	p := parser.NewParser(input, src, "ast")

	tree, err := p.Parse()
	if err != nil {
		if perr, ok := err.(common.Printable); ok {
			perr.Print()
		} else {
			fmt.Fprintf(os.Stderr, "Error parsing file - %s\n", err)
		}
		return
	}

	errs := typecheck.NewChecker().Check(tree)
	if len(errs) > 0 {
		for _, e := range errs {
			e.Print()
		}
		fmt.Fprintf(os.Stderr, "Found %d type error(s)\n", len(errs))
		return
	}

	if asJson {
		var data []byte
		var err error
		if prettyJson {
			data, err = json.MarshalIndent(tree, "", "  ")
		} else {
			data, err = json.Marshal(tree)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error converting to JSON - %s\n", err)
			return
		}
		os.Stdout.Write(data)
		return
	}

	fmt.Println("AST summary:")
	fmt.Printf("  %d global variables with explicit values\n", len(tree.Vars))
	fmt.Printf("  %d global variables with default values\n", len(tree.Typedefs))
	fmt.Printf("  %d global functions\n", len(tree.Funcs))
	fmt.Printf("  %d external functions\n", len(tree.Exerns))
	fmt.Printf("  %d structures\n", len(tree.Structs))

	fmt.Println()

	fmt.Println("Global variables with explicit values:")
	for _, v := range tree.Vars {
		fmt.Printf("  Variable %s: Node: %s\n", v.Name, v.Value.Node.Kind())
	}
	if len(tree.Vars) == 0 {
		fmt.Println("  (none)")
	}

	fmt.Println()

	fmt.Println("Global variables with default values:")
	for _, v := range tree.Typedefs {
		fmt.Printf("  Variable %s: Type: %s\n", v.Name, v.Type)
	}
	if len(tree.Typedefs) == 0 {
		fmt.Println("  (none)")
	}

	fmt.Println()

	fmt.Println("Global functions:")
	for _, f := range tree.Funcs {
		fmt.Printf("  Function %s -> %s\n", f.Name, f.Ret)
		fmt.Printf("    %d arguments:\n", len(f.Args))
		for _, a := range f.Args {
			fmt.Printf("      Argument %s: %s\n", a.Name, a.Type)
		}
		fmt.Printf("    %d nodes in body\n", len(f.Body))
	}
	if len(tree.Funcs) == 0 {
		fmt.Println("  (none)")
	}

	fmt.Println()

	fmt.Println("External functions:")
	for _, f := range tree.Exerns {
		fmt.Printf("  Function %s -> %s\n", f.Name, f.Ret)
		fmt.Printf("    %d arguments:\n", len(f.Args))
		for _, a := range f.Args {
			fmt.Printf("      Argument %s: %s\n", a.Name, a.Type)
		}
	}
	if len(tree.Exerns) == 0 {
		fmt.Println("  (none)")
	}

	fmt.Println()

	fmt.Println("Structures:")
	for _, s := range tree.Structs {
		fmt.Printf("  Structure %s:\n", s.Name)
		for _, f := range s.Fields {
			fmt.Printf("    Field %s: %s\n", f.Name, f.Type)
		}
	}
	if len(tree.Structs) == 0 {
		fmt.Println("  (none)")
	}
}

func compile() {
	panic("unimplemented")
}

func repl() {
	panic("unimplemented")
}

func debugFlag(arg string) error {
	for _, feat := range strings.Split(arg, ",") {
		feat = strings.ToLower(strings.TrimSpace(feat))
		switch feat {
		case "lexer":
			debug.GlobalAttr |= debug.AttrLexer
		case "parser":
			debug.GlobalAttr |= debug.AttrParser
		case "scope":
			debug.GlobalAttr |= debug.AttrScope
		default:
			return fmt.Errorf("unknown debug feature: %s", feat)
		}
	}
	return nil
}
