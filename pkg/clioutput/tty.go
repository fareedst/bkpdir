package clioutput

import (
	"io"
	"os"
)

func isTerminalWriter(w io.Writer) bool {
	switch w {
	case os.Stdout:
		return isTerminalFD(os.Stdout)
	case os.Stderr:
		return isTerminalFD(os.Stderr)
	default:
		return false
	}
}
