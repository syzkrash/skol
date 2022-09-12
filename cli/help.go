package cli

import (
	"fmt"
	"os"
)

var HelpCommand = Command{
	Name:  "help",
	Short: "Need help?",
	Long: `
Provides usage information for either the Skol CLI or a given command in the CLI
itself.`,
	Run: help,
}

func help(args []string) error {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: skol <command> [command arguments...]\n")
		fmt.Fprintf(os.Stderr, "Where <command> may be one of the following:\n")
		for _, cmd := range Commands {
			fmt.Fprintf(os.Stderr, "  %s - %s\n", cmd.Name, cmd.Short)
		}
		fmt.Fprintf(os.Stderr, "Use `skol help <command>' for more information.\n")
	} else {
		var cmd Command
		ok := false
		for _, c := range Commands {
			if c.Name == args[0] {
				cmd = c
				ok = true
			}
		}
		if ok {
			fmt.Fprintf(os.Stderr, "%s\n", cmd.Long)
		}
	}
	return nil
}
