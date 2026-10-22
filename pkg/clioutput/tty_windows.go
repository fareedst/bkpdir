//go:build windows

package clioutput

import (
	"io"
	"os"

	"golang.org/x/term"
)

func isTerminalFD(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
