package sim

import (
	"fmt"

	"github.com/syzkrash/skol/lexer"
)

type Call struct {
	Root  bool
	Func  string
	Where lexer.Position
}

func (c *Call) String() string {
	if c.Root {
		return "(top-level)"
	}
	return fmt.Sprintf("%s at %s", c.Func, c.Where)
}
