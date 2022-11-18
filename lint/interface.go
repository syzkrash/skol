package lint

import (
	"github.com/qeaml/all/slices"
	"github.com/syzkrash/skol/ast"
)

type Rule func(w chan *Warn, n ast.MetaNode) error

var newArrayFuncs = []string{"append", "concat"}

func newArrayRule(w chan *Warn, n ast.MetaNode) error {
	if n.Node.Kind() != ast.NFuncCall {
		return nil
	}
	fcn := n.Node.(ast.FuncCallNode)
	// fmt.Printf("newArray check: %s\n", fcn.Func)
	if slices.Contains(newArrayFuncs, fcn.Func) {
		w <- warning(WNewArray).NodeCause(n)
	}
	return nil
}

var Rules = map[string]Rule{
	"newArray": newArrayRule,
}
