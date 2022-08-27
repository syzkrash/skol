package sim

import (
	"fmt"

	"github.com/syzkrash/skol/ast"
)

type SimError struct {
	msg   string
	Cause ast.MetaNode
	Calls []*Call
}

func (s *SimError) Print() {
	fmt.Println("\x1b[31m\x1b[1mSimulation error:\x1b[0m")
	fmt.Println("   ", s.msg)
	fmt.Println("\x1b[1mCall stack:\x1b[0m")
	for _, c := range s.Calls {
		fmt.Println("   ", c.String())
	}
	fmt.Println("-->", "("+s.Cause.Node.Kind().String()+")", "at", s.Cause.Where)
	fmt.Println()
}

func (s *SimError) Error() string {
	return s.msg
}
