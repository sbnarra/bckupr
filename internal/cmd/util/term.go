package util

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func TermClear() {
	if true || term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Print("\033[H\033[2J")
	}
}
