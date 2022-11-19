package ir

import (
	"fmt"
	"strings"
)

// Opcode represents the unique number of an instruction
type Opcode byte

// Instruction constants
const (
	OpSet Opcode = iota
	OpCall
	OpRet
	OpBranch
	OpLoop
)

var opNames = []string{
	"SET",
	"CALL",
	"RET",
	"BRANCH",
	"LOOP",
}

func (op Opcode) String() string {
	return opNames[op]
}

// Instr represents any instruction. Each instruction has it's own type and
// stores the information about that instruction and may contian other
// instructions inside.
type Instr interface {
	// Op returns the unique opcode number for this opcode. The opcode number and
	// the factual type of this instruction must always match.
	Op() Opcode
}

// SetInstr holds the data of a SET instruction
type SetInstr struct {
	Target Ref
	Value  Value
}

// Op returns OpSet
func (SetInstr) Op() Opcode {
	return OpSet
}

func (i SetInstr) String() string {
	return fmt.Sprintf("%s %s, %s", OpSet, i.Target, i.Value)
}

var _ Instr = SetInstr{}

// CallInstr holds the data of a CALL instruction
type CallInstr struct {
	Func byte
	Args []Value
}

// Op returns OpCall
func (CallInstr) Op() Opcode {
	return OpCall
}

func (i CallInstr) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "%s %02x, [%02X](", OpCall, i.Func, len(i.Args))
	for n := 0; n < len(i.Args)-1; n++ {
		fmt.Fprintf(&str, "%s, ", i.Args[n])
	}
	fmt.Fprintf(&str, "%s)", i.Args[len(i.Args)-1])
	return str.String()
}

var _ Instr = CallInstr{}

// RetInstr holds the data of a RET instruction
type RetInstr struct {
	Value Value
}

// Op returns OpRet
func (RetInstr) Op() Opcode {
	return OpRet
}

func (i RetInstr) String() string {
	return fmt.Sprintf("%s %s", OpRet, i.Value)
}

var _ Instr = RetInstr{}

type branch struct {
	Cond Value
	Body []Instr
}

// BranchInstr holds the data of a BRANCH instruction (IR equivalent of an if
// statement)
type BranchInstr struct {
	branches []branch
}

// Op returns OpBranch
func (BranchInstr) Op() Opcode {
	return OpBranch
}

func (i BranchInstr) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "%s:\n", OpBranch)
	for _, b := range i.branches {
		fmt.Fprintf(&str, "  CASE %s:\n", b.Cond)
		for _, i := range b.Body {
			fmt.Fprintf(&str, "    %s\n", strings.ReplaceAll(fmt.Sprint(i), "\n", "\n    "))
		}
	}
	return str.String()
}

var _ Instr = BranchInstr{}

// LoopInstr holds the data of a LOOP instruction (IR equivalent of a while
// loop)
type LoopInstr struct {
	Cond Value
	Body []Instr
}

// Op returns OpLoop
func (LoopInstr) Op() Opcode {
	return OpLoop
}

func (i LoopInstr) String() string {
	str := strings.Builder{}
	fmt.Fprintf(&str, "%s %s:\n", OpLoop, i.Cond)
	for _, i := range i.Body {
		fmt.Fprintf(&str, "  %s\n", strings.ReplaceAll(fmt.Sprint(i), "\n", "\n  "))
	}
	return str.String()
}

var _ Instr = LoopInstr{}
