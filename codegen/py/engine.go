package py

import (
	"github.com/syzkrash/skol/codegen"
)

var Engine = codegen.Engine{
	Name:       "Python",
	Desc:       "Transpile Skol code to Python code.",
	Gen:        &generator{},
	Ephemeral:  false,
	Extension:  ".py",
	Exec:       executor{},
	Executable: false,
}

type executor struct{}

var _ codegen.FilenameExecutor = executor{}
