// 29 june 2012
// this file contains stuff useful to developing a project as fault-critical as an assembler
package main

import (
	"fmt"
	"os"
	"debug"
)

func _bug(mode string, format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "==========\n")
	fmt.Fprintf(os.Stderr, "!! %s !! ", mode)
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\nreport this to pietro/andlabs, including the following:\n")
	debug.PrintStack()
	// TODO dump more state?
	// TODO quit?
	fmt.Fprintf(os.Stderr, "==========\n")
}

func POTENTIAL_BUG(format string, args ...interface{}) {
	_bug("POTENTIAL BUG?", format, args...)
}
