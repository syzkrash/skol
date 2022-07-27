package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/syzkrash/skol/codegen"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/debug"
	"github.com/syzkrash/skol/python"
)

type Engine uint8

const (
	EnUndefined Engine = iota
	EnPy
	EnAST
	EnSim
)

var engines = [...]string{
	"Undefined",
	"Python",
	"AST Dump",
	"Simulation",
}

var theEngine = EnUndefined

func engineFlag(arg string) error {
	switch arg {
	case "py", "python":
		theEngine = EnPy
	case "ast":
		theEngine = EnAST
	case "sim":
		theEngine = EnSim
	default:
		return fmt.Errorf("unknown engine: %s", arg)
	}
	return nil
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

var input string

func main() {
	flag.StringVar(&input, "input", "", "File to compile/transpile/interpret. (Leave blank for REPL)")
	flag.Func("engine", "Which interpreter/compiler to use.", engineFlag)
	flag.Func("debug", "Which debug logs to show.", debugFlag)
	flag.Parse()

	fmt.Fprintf(os.Stderr, "Skol v%s\n", common.Version)

	if theEngine == EnUndefined {
		fmt.Fprintln(os.Stderr, "Specify an engine to use.")
		return
	}

	if input == "" {
		repl()
	} else {
		compile()
	}
}

func compile() {
	inFile, err := os.Open(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer inFile.Close()
	code, err := io.ReadAll(inFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	gen := gen(input, bytes.NewReader(code))

	if gen.CanGenerate() {
		outFile, err := os.Create(input + gen.Ext())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		defer outFile.Close()

		fmt.Fprintln(os.Stderr, "Compiling using engine:", engines[theEngine])
		compStart := time.Now()

		for {
			err = gen.Generate(outFile)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error -", err)
				return
			}
		}

		fmt.Fprintln(os.Stderr, "Compiled in", time.Since(compStart))
	}

	if gen.CanRun() {
		fmt.Fprintln(os.Stderr, "Running...")
		fmt.Fprintln(os.Stderr, "----------")
		if err = gen.Run(input + gen.Ext()); err != nil {
			if perr, ok := err.(common.Printable); ok {
				perr.Print()
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			return
		}
	}
}

func gen(fn string, src io.RuneScanner) codegen.Generator {
	switch theEngine {
	case EnPy:
		return python.NewPython(fn, src)
	case EnAST:
		return codegen.NewAST(fn, src)
	case EnSim:
		return codegen.NewSimEngine(fn, src)
	}
	return nil
}

func repl() {
	stdin := bufio.NewReader(os.Stdin)
	src := strings.NewReader("")
	gen := gen("REPL", src)

	if !gen.CanGenerate() {
		fmt.Fprintln(os.Stderr, "The given engine does not generate code.")
		return
	}

	fmt.Fprintln(os.Stderr, "Type a line of Skol code and hit Enter.")
	fmt.Fprintln(os.Stderr, "Code generated for the given engine will be printed.")
	fmt.Fprintln(os.Stderr, "Press ^C at any time to exit.")

	for {
		fmt.Fprint(os.Stderr, ">> ")
		line, err := stdin.ReadString('\n')
		if errors.Is(err, io.EOF) {
			fmt.Fprint(os.Stderr, "\n")
			return
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		src.Reset(line + "\n")
		if err = gen.Generate(os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, "Error -", err)
		}
	}
}
