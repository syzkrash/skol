package cli

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
	"github.com/syzkrash/skol/lint"
	"github.com/syzkrash/skol/parser"
)

var LintCommand = Command{
	Name:  "lint",
	Short: "Check a file for non-critical errors.",
	Long:  ``,
	Run:   runLint,
}

func runLint(args []string) error {
	if len(args) < 1 {
		return pe.New(pe.ENoInput)
	}

	input := args[0]

	srcf, err := os.Open(input)
	if err != nil {
		return pe.New(pe.EBadInput).Cause(err)
	}
	defer srcf.Close()

	srcraw, err := io.ReadAll(srcf)
	if err != nil {
		return pe.New(pe.EBadInput).Cause(err)
	}

	p := parser.NewParser(input, bytes.NewReader(srcraw), "lint")
	tree, err := p.Parse()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(lint.Rules))
	warns := make(chan *lint.Warn)
	for _, r := range lint.Rules {
		go func(r lint.Rule) {
			for _, f := range tree.Funcs {
				for _, n := range f.Body {
					check(warns, n, r)
				}
			}
			wg.Done()
		}(r)
	}

	go func() {
		wg.Wait()
		close(warns)
	}()

	for w := range warns {
		w.Print()
	}

	return nil
}

func check(w chan *lint.Warn, n ast.MetaNode, r lint.Rule) {
	r(w, n)
	switch n.Node.Kind() {
	case ast.NIf:
		in := n.Node.(ast.IfNode)
		for _, n := range in.Main.Block {
			check(w, n, r)
		}
		for _, b := range in.Other {
			for _, n := range b.Block {
				check(w, n, r)
			}
		}
		for _, n := range in.Else {
			check(w, n, r)
		}
	case ast.NWhile:
		wn := n.Node.(ast.WhileNode)
		for _, n := range wn.Block {
			check(w, n, r)
		}
	}
}
