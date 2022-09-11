package common

import (
	"fmt"
	"os"
	"os/exec"
)

func Cmd(name string, args ...string) error {
	fmt.Fprintf(os.Stderr, ": %s", name)
	for _, a := range args {
		fmt.Fprintf(os.Stderr, " %s", a)
	}
	fmt.Fprint(os.Stderr, "\n")
	return exec.Command(name, args...).Run()
}
