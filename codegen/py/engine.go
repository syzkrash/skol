package py

import "github.com/syzkrash/skol/codegen"

var Engine = codegen.Engine{
	Name:       "Python",
	Desc:       "Transpile Skol code to Python code.",
	Gen:        &generator{},
	Ephemeral:  true,
	Extension:  ".py",
	Exec:       nil,
	Executable: false,
}