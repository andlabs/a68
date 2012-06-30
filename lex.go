// 29 june 2012
package main

import (
	// FileLexer
	"utf8"
	"bufio"

	// Lexer
	"fmt"
	"os"
)

// TODO move to another file
var nErrors uint64

const (
	lexEOF rune = utf8.MaxRune + 1 + iota
	lexError
)

type FileLexer struct {
	Filename	string
	Lineno	uint64
	File		*bufio.Reader
	Tokens	chan yySymType
	inputLine	string
	tokStart	uint64
	tokEnd	uint64
	runeLen	uint64		// for ignoring the current character
	lastTok	int			// for automatically inserting ::
}

type lexState func(*FileLexer) lexState

func (l *FileLexer) Run() {
	for state := lex_next; state != nil; {
		state = state(l)
	}
}

func (l *FileLexer) Error(e string) {
	fmt.Fprintf(os.Stderr,
		"%s:%d %s\n",
		l.Filename, l.Lineno, e)
	nErrors++
}

func (l *FileLexer) Emit(toktype int) {
	l.Tokens <- yySymType{
		type:		toktype,
		value:	line[l.tokStart:l.tokEnd],
	}
	l.lastTok = toktype
	l.tokStart = l.tokEnd		// advance
}

func (l *FileLexer) read() (rune, error) {
	l.advance()
	if l.tokStart >= len(l.inputLine) {
		err := l.getline()
		if err == io.EOF {
			l.runeLen = 0			// don't unget an EOF
			return lexEOF, nil
		} else if err != nil {
			l.Error(fmt.Sprintf("error reading from file: %v", err))
			return lexError, err		// TODO more proper return?
		}
	}
	r, l.runeLen := utf8.DecodeRuneInString(len[i.tokStart:])
	l.tokEnd = l.tokStart + l.runeLen
	return r, nil
}

func (l *FileLexer) advance() {
	l.tokStart = l.tokEnd
	l.runeLen = 0
}

func (l *FileLexer) unget() {
	l.tokEnd -= l.runeLen
}

func (l *FileLexer) peek() (r rune, err error) {
	r, err = l.read()
	if err == nil {
		// TODO just quit on error?
		// l.read() gave us the right values anyway
		l.unget()
	}
	return
}

func (l *FileLexer) getline() error {
	line, err := l.File.ReadString('\n')
	if err != nil {
		return err
	}
	l.inputLine = line
	l.tokStart = 0
	l.tokEnd = 0
	l.runeLen = 0
	return nil
}

type Lexer struct {
	files		 []*FileLexer
	curfile	*FileLexer		// for simplifying the below code
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
		l.curfile.Error(e)
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
