//go:build windows

package py

import "github.com/syzkrash/skol/common"

func (e executor) Execute(fn string) error {
	return common.Cmd("py", "-3", fn)
}
