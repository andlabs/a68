// 27 july 2012
package main

import (
	"os"
	"log"
)

var AsmLexer *Lexer

func init() {
	AsmLexer = new(Lexer)
}

func main() {
	for i := len(os.Args) - 1; i > 0; i-- {
		err := AsmLexer.Open(os.Args[i])
		if err != nil {
			log.Fatalf("error opening %s: %v", os.Args[i], err)
		}
	}
	yyDebug = 100
	yyParse(AsmLexer)
}
