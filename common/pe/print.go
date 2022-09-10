package pe

import (
	"fmt"
	"os"
)

type seq string

const (
	ESC = "\x1b"
	CSI = ESC + "["

	gBold      seq = CSI + "1m"
	gUnderline seq = CSI + "4m"
	gNegative  seq = CSI + "7m"
	gDim       seq = CSI + "2m"

	gNoBold      seq = CSI + "22m"
	gNoUnderline seq = CSI + "24m"
	gPositive    seq = CSI + "27m"
	gNoDim       seq = gNoBold

	cFgBlack   seq = CSI + "30m"
	cFgRed     seq = CSI + "31m"
	cFgGreen   seq = CSI + "32m"
	cFgYellow  seq = CSI + "33m"
	cFgBlue    seq = CSI + "34m"
	cFgMagenta seq = CSI + "35m"
	cFgCyan    seq = CSI + "36m"
	cFgWhite   seq = CSI + "37m"
	cFgDefault seq = CSI + "39m"

	cBgBlack   seq = CSI + "40m"
	cBgRed     seq = CSI + "41m"
	cBgGreen   seq = CSI + "42m"
	cBgYellow  seq = CSI + "43m"
	cBgBlue    seq = CSI + "44m"
	cBgMagenta seq = CSI + "45m"
	cBgCyan    seq = CSI + "46m"
	cBgWhite   seq = CSI + "47m"
	cBgDefault seq = CSI + "49m"
)

func (s seq) reverse() seq {
	switch s {
	case gBold:
		return gNoBold
	case gUnderline:
		return gNoUnderline
	case gNegative:
		return gPositive
	case gDim:
		return gNoDim

	case cFgBlack, cFgRed, cFgGreen, cFgYellow, cFgBlue, cFgMagenta, cFgCyan, cFgWhite:
		return cFgDefault

	case cBgBlack, cBgRed, cBgGreen, cBgYellow, cBgBlue, cBgMagenta, cBgCyan, cBgWhite:
		return cBgDefault
	}

	return ""
}

func pprintln(msg string, g ...seq) {
	if isatty() {
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
