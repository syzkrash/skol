package pe

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

type seq string

const (
	ESC = "\x1b"
	CSI = ESC + "["

	Bold      seq = CSI + "1m"
	Underline seq = CSI + "4m"
	Negative  seq = CSI + "7m"
	Dim       seq = CSI + "2m"

	NoBold      seq = CSI + "22m"
	NoUnderline seq = CSI + "24m"
	Positive    seq = CSI + "27m"
	NoDim       seq = NoBold

	FgBlack   seq = CSI + "30m"
	FgRed     seq = CSI + "31m"
	FgGreen   seq = CSI + "32m"
	FgYellow  seq = CSI + "33m"
	FgBlue    seq = CSI + "34m"
	FgMagenta seq = CSI + "35m"
	FgCyan    seq = CSI + "36m"
	FgWhite   seq = CSI + "37m"
	FgDefault seq = CSI + "39m"

	BgBlack   seq = CSI + "40m"
	BgRed     seq = CSI + "41m"
	BgGreen   seq = CSI + "42m"
	BgYellow  seq = CSI + "43m"
	BgBlue    seq = CSI + "44m"
	BgMagenta seq = CSI + "45m"
	BgCyan    seq = CSI + "46m"
	BgWhite   seq = CSI + "47m"
	BgDefault seq = CSI + "49m"
)

func (s seq) reverse() seq {
	switch s {
	case Bold:
		return NoBold
	case Underline:
		return NoUnderline
	case Negative:
		return Positive
	case Dim:
		return NoDim

	case FgBlack, FgRed, FgGreen, FgYellow, FgBlue, FgMagenta, FgCyan, FgWhite:
		return FgDefault

	case BgBlack, BgRed, BgGreen, BgYellow, BgBlue, BgMagenta, BgCyan, BgWhite:
		return BgDefault
	}

	return ""
}

func Pprintln(msg string, g ...seq) {
	if isatty.IsTerminal(os.Stderr.Fd()) {
		for _, s := range g {
			fmt.Fprint(os.Stderr, s)
		}
		defer func() {
			for _, s := range g {
				fmt.Fprint(os.Stderr, s.reverse())
			}
		}()
	}
	fmt.Fprint(os.Stderr, msg+"\n")
}
