package ir

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

// InstrNames holds the instructions' names
var InstrNames = []string{
	"SET",
	"CALL",
	"RET",
	"BRANCH",
	"LOOP",
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

var _ Instr = CallInstr{}

// RetInstr holds the data of a RET instruction
type RetInstr struct {
	Value Value
}

// Op returns OpRet
func (RetInstr) Op() Opcode {
	return OpRet
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

var _ Instr = LoopInstr{}
