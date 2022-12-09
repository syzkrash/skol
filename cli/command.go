package cli

// Command represents a single command of the `skol` CLI.
type Command struct {
	Name  string
	Short string
	Long  string
	Run   func(args []string) error
}

// Commands contains all the currently available CLI commands.
var Commands []Command

func init() {
	Commands = []Command{
		HelpCommand,
		AstCommand,
		CompileCommand,
		ReplCommand,
		LintCommand,
	}
}
