package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/syzkrash/skol/codegen"
	"github.com/syzkrash/skol/codegen/py"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/typecheck"
)

// CompileCommand represents the `skol compile` command
var CompileCommand = Command{
	Name:  "compile",
	Short: "Compile a file",
	Long: `
Usage: skol compile <engine> <file> [arguments...]
Where arguments can be any combination of:
  -run :: If the engine allows, run the result.

Depending on the engine specified, this will either:
  a) Compile the given file into an executable.
  b) Transpile it into another language.`,
	Run: compile,
}

func compile(args []string) error {
	if len(args) < 1 {
		return pe.New(pe.EUnknownEngine)
	}
	if len(args) < 2 {
		return pe.New(pe.ENoInput)
	}

	input := args[1]
	engine := args[0]

	var (
		run bool
	)

	flags := flag.NewFlagSet("skol compile", flag.ContinueOnError)
	flags.BoolVar(&run, "run", false, "")
	flags.Parse(args[2:])

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

	p := parser.NewParser(input, bytes.NewReader(srcraw), engine, errs)
	ast := p.Parse()
	if errOne != nil {
		close(errs)
		return errOne
	}

	// cool note:
	// It is neat to have a separate goroutine for typechecking and a separate one
	// for printing to stderr. Because printing to stderr is quite slow, this does
	// offer a very slight speedup, especially in case of many errors.

	errs = make(chan error)
	errOne = nil

	typecheck.NewChecker(errs).Check(ast)
	close(errs)

	if errOne != nil {
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

	if run {
		switch e.Exec.(type) {
		case codegen.EphemeralExecutor:
			e.Exec.(codegen.EphemeralExecutor).Execute(ephBuf)
		case codegen.FilenameExecutor:
			e.Exec.(codegen.FilenameExecutor).Execute(input + e.Extension)
		}
	}

	return nil
}
