package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/syzkrash/skol/codegen"
	"github.com/syzkrash/skol/codegen/py"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/common/pe"
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
	flags.StringVar(&engine, "engine", "", "")
	flags.StringVar(&input, "input", "", "")
	flags.Func("debug", "", debugFlag)
	flags.BoolVar(&asJson, "json", false, "")
	flags.BoolVar(&prettyJson, "pretty-json", false, "")

	flags.Parse(os.Args[2:])

	fmt.Fprintf(os.Stderr, "Skol v%s\n", common.Version)

	err := cli(os.Args[1])
	if err != nil {
		if printable, ok := err.(common.Printable); ok {
			printable.Print()
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		}
	}
}

func cli(act string) error {
	switch act {
	case "ast":
		return dumpAst()
	case "build":
		return compile()
	case "repl":
		return repl()
	case "":
		usage()
		return nil
	default:
		e := pe.New(pe.EUnknownAction)
		e.Section("Action", act)
		return e
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <action> [arguments...]\n", os.Args[0])
	fmt.Fprintln(os.Stderr, "Where action can be one of the following:")
	fmt.Fprintln(os.Stderr, "  build  - Compile/run a file using a given engine")
	fmt.Fprintln(os.Stderr, "  repl   - Start the skol Read And Evaluate Loop (REPL)")
	fmt.Fprintln(os.Stderr, "  ast    - Print the Abstract Syntax Tree (AST) of a file")
	fmt.Fprintln(os.Stderr, "Where arguments can be any of the following:")
	fmt.Fprintln(os.Stderr, "  -input <filename>  - File to use as input")
	fmt.Fprintln(os.Stderr, "  -engine <engine>   - Engine to compile/run with")
	fmt.Fprintln(os.Stderr, "  -debug <kinds>     - Enable specified kinds of debug messages")
	fmt.Fprintln(os.Stderr, "  -json              - For applicable actions, output as JSON")
}

var (
	input      string
	asJson     bool
	prettyJson bool
)

func dumpAst() error {
	if input == "" {
		return pe.New(pe.ENoInput)
	}

	f, err := os.Open(input)
	if err != nil {
		return pe.New(pe.EBadInput).Cause(err)
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return pe.New(pe.EBadInput).Cause(err)
	}

	src := bytes.NewReader(data)
	p := parser.NewParser(input, src, "ast")

	tree, err := p.Parse()
	if err != nil {
		return err
	}

	errs := typecheck.NewChecker().Check(tree)
	if len(errs) > 0 {
		for _, e := range errs {
			e.Print()
		}
		fmt.Fprintf(os.Stderr, "Found %d type error(s)\n", len(errs))
		return nil
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
			return pe.New(pe.EBadAST).Cause(err)
		}
		os.Stdout.Write(data)
		return nil
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

	return nil
}

var (
	engine string
)

func compile() error {
	if input == "" {
		return pe.New(pe.ENoInput)
	}

	srcf, err := os.Open(input)
	if err != nil {
		return pe.New(pe.EBadInput).Cause(err)
	}
	defer srcf.Close()

	srcraw, err := io.ReadAll(srcf)
	if err != nil {
		return pe.New(pe.EBadInput).Cause(err)
	}

	var e codegen.Engine
	switch engine {
	case "py":
		e = py.Engine
	default:
		return pe.New(pe.EUnknownEngine).Section("Engine", engine)
	}

	p := parser.NewParser(input, bytes.NewReader(srcraw), engine)
	ast, err := p.Parse()
	if err != nil {
		return err
	}

	var ephBuf *bytes.Buffer
	var out io.Writer
	if e.Ephemeral {
		ephBuf = &bytes.Buffer{}
		out = ephBuf
	} else {
		outf, err := os.Create(input + e.Extension)
		if err != nil {
			return pe.New(pe.EBadOutput).Cause(err)
		}
		defer outf.Close()
		out = outf
	}

	e.Gen.Output(out)

	if astgen, ok := e.Gen.(codegen.ASTGenerator); ok {
		astgen.Input(ast)
	}

	err = e.Gen.Generate()
	if err != nil {
		return pe.New(pe.EBadOutput).Cause(err)
	}

	switch e.Exec.(type) {
	case codegen.EphemeralExecutor:
		e.Exec.(codegen.EphemeralExecutor).Execute(ephBuf)
	case codegen.FilenameExecutor:
		e.Exec.(codegen.FilenameExecutor).Execute(input + e.Extension)
	}

	return nil
}

func repl() error {
	return pe.New(pe.EUnimplemented)
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
