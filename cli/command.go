package cli

type Command struct {
	Name  string
	Short string
	Long  string
	Run   func(args []string) error
}

var Commands []Command

func init() {
	Commands = []Command{
		HelpCommand,
		AstCommand,
		CompileCommand,
		ReplCommand,
	}
}
