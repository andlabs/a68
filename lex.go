// 29 june 2012
package main

import (
	// FileLexer

	// Lexer
	"fmt"
	"os"
)

type FileLexer struct {
	Filename	string
	Lineno	uint64
	File		*os.File
	Tokens	chan int
}

type lexState func(*FileLexer) lexState

type Lexer struct {
	files		 []*FileLexer
	Errors	uint64
}

func (l *Lexer) Lex(lval *yySymType) int {
	if len(l.files) <= 0 {
		POTENTIAL_BUG(
			"attempted to lex before any files are open")
		return -1
	}
	cf := l.files[len(l.files) - 1]		// TODO: replace this with a curfile pointer?
	tok := <-cf.Tokens
	if tok == -1 {			// EOF
		cf.File.Close()
		l.files = l.files[:len(l.files) - 2]
		if len(l.files) <= 0 {	// no more files
			return -1
		}
		return l.Lex(lval)
	}
	return tok
}

func (l *Lexer) Error(e string) {
	if len(l.files) > 0 {
		cf := l.files[len(l.files) - 1]
		fmt.Fprintf(os.Stderr,
			"%s:%d %s\n",
			cf.Filename, cf.Lineno, e)
		l.Errors++
		return
	}
	POTENTIAL_BUG(
		"error reported before any files were opened/after all were closed: %s",
		e)
}

func (l *Lexer) Open(s string) error {
	f, err := os.Open(s)
	if err != nil {
		return err
	}
	// do not defer f.Close(); we do that when we're done reading input
	l.files = append(l.files, &FileLexer{
		Filename:		s,
		File:			f,
		Tokens:		make(chan int)
	})
	return nil
}
