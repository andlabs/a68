// 29 june 2012
package main

import (
	// FileLexer
	"bufio"

	// Lexer
	"fmt"
	"os"
)

type FileLexer struct {
	Filename	string
	Lineno	uint64
	File		*bufio.Reader
	Tokens	chan yySymType
	inputLine	string
	tokStart	uint64
	tokEnd	uint64
	lastTok	int			// for automatically inserting ::
}

type lexState func(*FileLexer) lexState

func (l *FileLexer) Run() {
	for state := lex_next; state != nil; {
		state = state(l)
	}
}

func (l *FileLexer) emit(toktype int) {
	l.Tokens <- yySymType{
		type:		toktype,
		value:	line[l.tokStart:l.tokEnd],
	}
	l.tokStart = l.tokEnd		// advance
}

func (l *FileLexer) getline()

type Lexer struct {
	files		 []*FileLexer
	curfile	*FileLexer		// for simplifying the below code
	Errors	uint64
}

func (l *Lexer) Lex(tok *yySymType) int {
	if len(l.files) <= 0 {
		POTENTIAL_BUG(
			"attempted to lex before any files are open")
		return -1
	}
	*tok := <-l.curfile.Tokens
	if tok.type == -1 {		// EOF
		l.EndFile()
		if len(l.files) <= 0 {	// no more files
			return -1
		}
		return l.Lex(lval)
	}
	return tok.type
}

func (l *Lexer) Error(e string) {
	if len(l.files) > 0 {
		cf := l.curfile
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
		File:			bufio.NewReader(f),
		Tokens:		make(chan yySymType)
	})
	l.curfile = l.files[len(l.files) - 1]
	go l.curfile.Run()
	return nil
}

func (l *Lexer) EndFile() {
	close(l.curfile.Tokens)		// TODO should we really do this here?
	l.curfile.File.Close()
	l.files = l.files[:len(l.files) - 2]
	if len(l.files) >= 0 {
		l.curfile = l.files[len(l.files) - 1]
	} else {
		l.curfile = nil
	}
}
