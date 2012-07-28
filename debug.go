
// 29 june 2012
// this file contains stuff useful to developing a project as fault-critical as an assembler
package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"encoding/hex"
	"log"
)

func _bug(mode string, format string, args ...interface{}) {
	out := os.Stderr

	fmt.Fprintf(out, "==========\n")
	fmt.Fprintf(out, "!! %s !! ", mode)
	fmt.Fprintf(out, format, args...)
	fmt.Fprintf(out, "\nreport this to pietro/andlabs, including the following:\n")

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Stack trace: ")
	debug.PrintStack()	// TODO write to out

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Lexer state: ")
	if len(AsmLexer.files) == 0 {
		fmt.Fprintln(out, "no open files")
	} else {
		found := false
		for i := 0; i < len(AsmLexer.files); i++ {
			if AsmLexer.files[i] == AsmLexer.curfile {
				fmt.Fprintf(out, "> ")
				found = true
			} else {
				fmt.Fprintf(out, "  ")
			}
			fmt.Fprintf(out, "%#v\n", AsmLexer.files[i])
		}
		if !found {
			fmt.Fprintln(out, "...huh, somehow the current file isn't one of the above?")
			fmt.Fprintf(out, "current file: %#v\n", AsmLexer.curfile)
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Global symbols: ")
	if len(Symbols.m) == 0 {
		fmt.Fprintln(out, "none")
	} else {
		for k, v := range Symbols.m {
			fmt.Fprintf(out, "%s :: %s\n", k, v.String())
		}
	}

	fmt.Fprintln(out)
	fmt.Fprintln(out, "Local symbols: ")
	if len(Locals.m) == 0 {
		fmt.Fprintln(out, "none")
	} else {
		for k, v := range Locals.m {
			fmt.Fprintf(out, "%s :: %s\n", k, v.String())
		}
	}

	// TODO dump later expressions

	// TODO dump more state?

	// tristanseifert suggested making this last because I couldn't decide whether or not it should be first or last
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Present assembled data: ")
	if Output.Len() > 0 {
		fmt.Fprintln(out, hex.Dump(Output.Bytes()))
	} else {
		fmt.Fprintln(out, "no output")
	}
	if OutBits.Pos() == -1 {
		// TODO if I ever split bitwriter.go out into its own library, this will need to be changed!
		fmt.Fprintf(out, "we are still inside a byte! bitCount = %d, curByte = $%02X\n", OutBits.bitCount, OutBits.curByte)
	}

	// TODO quit?
	fmt.Fprintf(out, "==========\n")
}

func POTENTIAL_BUG(format string, args ...interface{}) {
	_bug("POTENTIAL BUG?", format, args...)
}

func FATAL_BUG(format string, args ...interface{}) {
	_bug("FATAL BUG", format, args...)
	log.Fatalf("fatal bug reported above; aborting assembly")
}
