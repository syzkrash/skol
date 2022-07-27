package typecheck

import (
	"fmt"

	"github.com/syzkrash/skol/parser/nodes"
	"github.com/syzkrash/skol/parser/values/types"
)

type TypeError struct {
	Got  types.Type
	Want types.Type
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
	if e.Want.Prim() != types.PNothing {
		fmt.Println("\x1b[1mWanted type:\x1b[0m")
		fmt.Println("   ", e.Want.String())
		if e.Want.Prim() == types.PStruct {
			for _, f := range e.Want.(types.StructType).Fields {
				fmt.Println("       ", f.Name, f.Type.String())
			}
		}
	}
	if e.Got.Prim() != types.PNothing {
		fmt.Println("\x1b[1mFound type:\x1b[0m")
		fmt.Println("   ", e.Got.String())
		if e.Got.Prim() == types.PStruct {
			for _, f := range e.Got.(types.StructType).Fields {
				fmt.Println("       ", f.Name, f.Type.String())
			}
		}
	}
	fmt.Println()
}
