package cli

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/typecheck"
)

// AstCommand defines the `skol ast` command.
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
	Run: runAst,
}

func runAst(args []string) error {
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

	tree, err := parseOrCacheAST(input)
	if err != nil {
		return err
	}

	errs := make(chan error)

	go func() {
		typecheck.NewChecker(errs).Check(tree)
		close(errs)
	}()

	var errOne error

	for err := range errs {
		if err == nil {
			continue
		}

		if errOne == nil {
			errOne = err
			continue
		}

		if perr, ok := err.(common.Printable); ok {
			perr.Print()
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s", err)
		}
	}

	if errOne != nil {
		return err
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

func parseOrCacheAST(input string) (tree ast.AST, err error) {
	debug.Log(debug.AttrCache, "Checking for AST cache for %s", input)
	cacheName := common.CachedASTName(input)
	asts, err := os.Stat(cacheName)
	if err == nil {
		srcs, err := os.Stat(input)
		if err == nil {
			if srcs.ModTime().Before(asts.ModTime()) {
				return loadCachedAST(cacheName)
			}
		}
	}
	return parseAST(input, cacheName)
}

func loadCachedAST(input string) (tree ast.AST, err error) {
	debug.Log(debug.AttrCache, "Loading cached AST from %s", input)
	// we don't need to check if the file exists as this function would not get
	// called if it didn't (due to the os.Stat call)
	f, _ := os.Open(input)
	defer f.Close()
	return ast.Decode(f)
}

func parseAST(input, cacheName string) (tree ast.AST, err error) {
	debug.Log(debug.AttrCache, "No cached AST found for %s", input)
	f, err := os.Open(input)
	if err != nil {
		err = pe.New(pe.EBadInput).Cause(err)
		return
	}
	data, err := io.ReadAll(f)
	f.Close()
	if err != nil {
		err = pe.New(pe.EBadInput).Cause(err)
		return
	}

	errs := make(chan error)
	var errOne error

	go func() {
		for err := range errs {
			if err == nil {
				continue
			}

			if errOne == nil {
				errOne = err
				continue
			}

			if perr, ok := err.(common.Printable); ok {
				perr.Print()
			} else {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			}
		}
	}()

	src := bytes.NewReader(data)
	p := parser.NewParser(input, src, "ast", errs)

	tree = p.Parse()
	close(errs)
	if errOne != nil {
		err = errOne
		return
	}

	f, err = os.Create(cacheName)
	if err != nil {
		err = pe.New(pe.EBadAST).Cause(err)
		return
	}
	err = ast.Encode(f, tree)
	if err == nil {
		debug.Log(debug.AttrCache, "Cached AST for %s", input)
	}
	f.Close()
	return
}
