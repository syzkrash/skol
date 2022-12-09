package cli

import "github.com/syzkrash/skol/common/pe"

// ReplCommand represents the `skol repl` command.
var ReplCommand = Command{
	Name:  "repl",
	Short: "Start a REPL session",
	Long: `
Usage: skol repl <engine>

Starts a Read-Evaluate-Print-Loop (REPL) session for the given engine. If the
engine does not generate printable code (a.k.a. isn't a transpiler), the REPL
will not start.`,
	Run: repl,
}

func repl(args []string) error {
	return pe.New(pe.EUnimplemented)
}
