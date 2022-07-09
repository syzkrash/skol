package typecheck

import (
	"fmt"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values"
)

type TypeError struct {
	Got  *values.Type
	Want *values.Type
	Node nodes.Node
	msg  string
}

func (e *TypeError) Error() string {
	return e.msg
}

func (e *TypeError) Print() {
	fmt.Println("\x1b[31m\x1b[1mType error:\x1b[0m")
	fmt.Println("   ", e.msg)
	fmt.Println("\x1b[1mCaused by:\x1b[0m")
	fmt.Println("   ", e.Node.Kind(), "at", e.Node.Where())
	fmt.Println("\x1b[1mWanted type:\x1b[0m")
	fmt.Println("   ", e.Want.Name())
	if e.Want.Prim == values.PStruct {
		for _, f := range e.Want.Structure.Fields {
			fmt.Println("       ", f.Name, f.Type.Name())
		}
	}
	fmt.Println("\x1b[1mFound type:\x1b[0m")
	fmt.Println("   ", e.Got.Name())
	if e.Got.Prim == values.PStruct {
		for _, f := range e.Got.Structure.Fields {
			fmt.Println("       ", f.Name, f.Type.Name())
		}
	}
	fmt.Println()
}
