package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/syzkrash/skol/cli"
	"github.com/syzkrash/skol/common"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/debug"
)

func main() {
	var c string
	if len(os.Args) <= 1 {
		c = "help"
	} else {
		c = os.Args[1]
	}

	var a []string
	if len(os.Args) <= 2 {
		a = []string{}
	} else {
		a = os.Args[2:]
	}

	var cmd cli.Command
	ok := false
	for _, cmd_ := range cli.Commands {
		if cmd_.Name == c {
			cmd = cmd_
			ok = true
		}
	}
	if ok {
		if err := cmd.Run(a); err != nil {
			if p, ok := err.(common.Printable); ok {
				p.Print()
			} else {
				fmt.Fprintf(os.Stderr, "%s", err.Error())
			}
		}
	} else {
		pe.New(pe.EUnknownAction).Print()
	}
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
