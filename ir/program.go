package ir

// Block is a series of instructions
type Block []Instr

// Program represents a full program in IR form
type Program struct {
	Entrypoint byte
	Globals    []Value
	Funcs      []Block
}
