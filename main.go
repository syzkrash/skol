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

var input string

func main() {
	flag.StringVar(&input, "input", "", "File to compile/transpile/interpret. (Leave blank for REPL)")
	flag.Func("engine", "Which interpreter/compiler to use.", engineFlag)
	flag.Parse()

	fmt.Printf("Skol v%s\n", common.Version)

	if theEngine == EnUndefined {
		fmt.Println("Specify an engine to use.")
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
		fmt.Println(err)
		return
	}
	defer inFile.Close()
	code, err := io.ReadAll(inFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	gen := gen(input, bytes.NewReader(code))
	outFile, err := os.Create(input + gen.Ext())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outFile.Close()

	if gen.CanGenerate() {
		fmt.Println("Compiling using engine:", engines[theEngine])
		compStart := time.Now()

		for {
			err = gen.Generate(outFile)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Println("Error -", err)
				return
			}
		}

		fmt.Println("Compiled in", time.Since(compStart))
	}

	if gen.CanRun() {
		fmt.Println("Running...")
		fmt.Println("----------")
		if err = gen.Run(input + ".py"); err != nil {
			fmt.Println(err)
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
	fmt.Println("Type a line of Skol code and hit Enter.")
	fmt.Println("Code generated for the given engine will be printed.")
	fmt.Println("Press ^C at any time to exit.")

	stdin := bufio.NewReader(os.Stdin)
	src := strings.NewReader("")
	gen := gen("REPL", src)

	for {
		fmt.Print(">> ")
		line, err := stdin.ReadString('\n')
		if errors.Is(err, io.EOF) {
			fmt.Print("\n")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		src.Reset(line + "\n")
		if err = gen.Generate(os.Stdout); err != nil {
			fmt.Println("Error -", err)
		}
	}
}
