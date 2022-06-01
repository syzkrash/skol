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

	"github.com/syzkrash/skol/codegen"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/lexer"
	"github.com/syzkrash/skol/parser"
	"github.com/syzkrash/skol/python"
)

type Engine uint8

const (
	EnUndefined Engine = iota
	EnPy
	EnAST
)

var theEngine = EnUndefined

func engineFlag(arg string) error {
	switch arg {
	case "py", "python":
		theEngine = EnPy
	case "ast":
		theEngine = EnAST
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
	outFile, err := os.Create(input + ".py")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer outFile.Close()
	gen := gen(input, bytes.NewReader(code))
	err = gen.Generate(outFile)
	if err != nil {
		var perr *parser.ParserError
		var lerr *lexer.LexerError
		if errors.As(err, &perr) {
			fmt.Println("Parser error at", perr.Where, "-", perr.Error())
		} else if errors.As(err, &lerr) {
			fmt.Println("Lexer error at", lerr.Where, "-", lerr.Error())
		} else {
			fmt.Println("Error -", err)
		}
		return
	}
	if gen.CanRun() {
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
