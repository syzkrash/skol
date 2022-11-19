package ir

import (
	"fmt"
	"strings"
)

// Block is a series of instructions
type Block []Instr

func (b Block) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "BLOCK (%d):\n", len(b))
	for _, i := range b {
		fmt.Fprintf(&str, "  %s\n", strings.ReplaceAll(fmt.Sprint(i), "\n", "\n  "))
	}
	return str.String()
}

// Program represents a full program in IR form
type Program struct {
	Entrypoint byte
	Globals    []Value
	Funcs      []Block
}

func (p Program) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "ENTRY %02X\n", p.Entrypoint)
	fmt.Fprintf(&str, "GLOBALS (%d):\n", len(p.Globals))
	for i, g := range p.Globals {
		fmt.Fprintf(&str, "  %02X: %s\n", i, strings.ReplaceAll(fmt.Sprint(g), "\n", "\n  "))
	}
	fmt.Fprintf(&str, "FUNCS (%d):\n", len(p.Funcs))
	for i, f := range p.Funcs {
		fmt.Fprintf(&str, "  %02X: %s\n", i, strings.ReplaceAll(fmt.Sprint(f), "\n", "\n  "))
	}
	return str.String()
}
