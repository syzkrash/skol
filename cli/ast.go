package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/typecheck"
)

var AstCommand = Command{
	Name:  "ast",
	Short: "Dump file AST",
	Long: `
Usage: skol ast <file> [arguments...]
Where arguments can be any combination of:
  -json   :: Print AST as JSON
  -pretty :: If -json is specified, make the JSON human-readable.

By default, this parses and typechecks the file, then prints a summary of the
resulting AST. If -json is provided, the AST is encoded as JSON and printed to
stdout. If -pretty is also provided, additional whitespace is added to the JSON
to increase readability.`,
	Run: ast,
}

func ast(args []string) error {
	if len(args) < 1 {
		return pe.New(pe.ENoInput)
	}

	input := args[0]

	var (
		asJSON     bool
		prettyJSON bool
	)

	flags := flag.NewFlagSet("skol ast", flag.ContinueOnError)
	flags.BoolVar(&asJSON, "json", false, "")
	flags.BoolVar(&prettyJSON, "pretty", false, "")
	flags.Parse(args[1:])

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

	if asJSON {
		var data []byte
		var err error
		if prettyJSON {
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
