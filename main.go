package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/syzkrash/skol/parser"
)

func main() {
	fmt.Println("Welcome to Skol.")
	fmt.Println("Type a line of Skol code and hit Enter.")
	fmt.Println("AST node generated from that line will be printed.")
	fmt.Println("Press ^C at any time to exit.")

	input := bufio.NewReader(os.Stdin)
	src := strings.NewReader("")
	par := parser.NewParser("stdin", src)

	for {
		fmt.Print(">> ")
		line, err := input.ReadString('\n')
		if errors.Is(err, io.EOF) {
			fmt.Print("\n")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		src.Reset(line)
		n, err := par.Next()
		if err != nil {
			var perr *parser.ParserError
			if errors.As(err, &perr) {
				fmt.Println("Error on token", perr.Where, "-", perr.Error())
			} else {
				fmt.Println("Error -", err)
			}
		} else {
			fmt.Println(n)
		}
	}
}
