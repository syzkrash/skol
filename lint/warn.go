package lint

import (
	"fmt"

	"github.com/syzkrash/skol/ast"
	"github.com/syzkrash/skol/common/pe"
)

type WarnCode uint16

const (
	WDummy WarnCode = 100 + iota
	WNewArray
	WInfiniteLoop
)

var wmsgs = map[WarnCode]string{
	WDummy:        "dummy warning",
	WNewArray:     "This function does not modify the array in-place: it returns a new array.",
	WInfiniteLoop: "Avoid infinite loops.",
}

type section struct {
	title   string
	message string
}

type Warn struct {
	Code     WarnCode
	cause    error
	sections []section
}

func (w *Warn) Error() string {
	return wmsgs[w.Code]
}

func (w *Warn) Unwrap() error {
	return w.cause
}

func (w *Warn) Print() {
	pe.Pprintln(fmt.Sprintf("Warning %d", w.Code), pe.Bold, pe.FgYellow)
	pe.Pprintln("  " + wmsgs[w.Code])
	for _, s := range w.sections {
		pe.Pprintln(s.title, pe.Bold)
		pe.Pprintln("  " + s.message)
	}
}

func (w *Warn) Section(title, mfmt string, margs ...any) *Warn {
	w.sections = append(w.sections, section{
		title:   title,
		message: fmt.Sprintf(mfmt, margs...),
	})
	return w
}

func (w *Warn) Cause(e error) *Warn {
	w.cause = e
	return w.Section("Cause", e.Error())
}

func (w *Warn) NodeCause(m ast.MetaNode) *Warn {
	return w.Section("Cause", fmt.Sprintf("%s at %s", m.Node.Kind(), m.Where))
}

func warning(c WarnCode) *Warn {
	return &Warn{
		Code:     c,
		sections: []section{},
	}
}
