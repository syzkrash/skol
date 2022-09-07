package sim

import (
	"fmt"

	"github.com/syzkrash/skol/lexer"
)

// Call represents 1 call in the call stack along with position information.
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
