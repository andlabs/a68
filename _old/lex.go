// 29 june 2012
package main

import (
	// FileLexer
	"unicode"
	"unicode/utf8"
	"bufio"
	"io"
	"strings"
	"log"

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

// a token that cannot mark a statement end is a "continuation token"
// isCont is a generic value for a continuation token for immediately after we emit one
const isCont = OR

// and these are the continuation tokens
var isContinuationToken = map[int]bool{
	OR:		true,
	AND:	true,
	EQ:		true,
	NE:		true,
	'<':		true,
	LE:		true,
	'>':		true,
	GE:		true,
	'+':		true,
	'-':		true,
	'|':		true,
	'^':		true,
	'*':		true,
	'/':		true,
	'%':		true,
	'&':		true,
	LSH:		true,
	RSH:		true,
	'~':		true,
	'!':		true,
	'(':		true,
	'{':		true,			// TODO this will be used later
	':':		true,
	// TODO add TERM? will this save time?
}

type FileLexer struct {
	Filename	string
	Lineno	uint64
	File		*bufio.Reader
	Closer	io.Closer
	Tokens	chan yySymType
	inputLine	string
	tokStart	uint64
	readPos	uint64		// also the position where the token ends
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

// shorthand
func (l *FileLexer) what() string {
	return l.inputLine[l.tokStart:l.readPos]
}

func (l *FileLexer) emit(toktype int) {
	l.Tokens <- yySymType{
		toktype:	toktype,
		value:	l.what(),
	}
	l.lastTok = toktype
	l.tokStart = l.readPos		// advance
}

func (l *FileLexer) getline() error {
	line, err := l.File.ReadString('\n')
	if err == io.EOF {
		return err
	} else if err != nil {
		// apparently most compiliers and assemblers just terminate on read error so let's do it too
		log.Fatalf("error reading from file %s; assembly terminated: %v\n",
			l.Filename, err)
	}
	l.inputLine = line
	l.tokStart = 0
	l.readPos = 0
	l.runeLen = 0
	return nil
}

func (l *FileLexer) read() rune {
	if l.readPos >= uint64(len(l.inputLine)) {
		err := l.getline()
		if err == io.EOF {
			l.runeLen = 0			// don't unget an EOF
			return lexEOF
		}
	}
	r, rl := utf8.DecodeRuneInString(l.inputLine[l.readPos:])
	l.runeLen = uint64(rl)
	l.readPos += l.runeLen
	// tokStart is updated either when we emit a token or when we ignore one
	return r
}

func (l *FileLexer) ignore() {
	l.tokStart = l.readPos
	l.runeLen = 0
}

func (l *FileLexer) unget() {
	l.readPos -= l.runeLen
}

func (l *FileLexer) peek() (r rune, err error) {
	r = l.read()
	l.unget()
	return
}

func (l *FileLexer) accept(r rune) bool {
	c := l.read()
	if r != c {
		l.unget()
		return false
	}
	return true
}

func (l *FileLexer) acceptAndEmit(r rune, ifSo int, ifNot int) {
	if l.accept(r) {
		l.emit(ifSo)
	} else {
		l.emit(ifNot)
	}
}

func (l *FileLexer) acceptRun(chars string) bool {
	found := false
	for {
		c := l.read()
		if !strings.ContainsRune(chars, c) {
			l.unget()
			break
		}
		found = true
	}
	return found
}

func lex_next(l *FileLexer) lexState {
	c := l.read()
	switch {
	case c == lexEOF:
		return lex_end
	case c == '\n':
		l.ignore()
		l.Lineno++
		if !isContinuationToken[l.lastTok] {		// add a terminator unless it's nonsensical to do so
			l.emit(TERM)
			l.lastTok = isCont
		}
		return lex_next
	case c == ';':			// comment; eat line
		for {
			c = l.read()
			if c == '\n' {
				l.unget()
				break
			}
			l.ignore()
		}
		return lex_next
	case unicode.IsSpace(c):
		l.ignore()
		return lex_next
	case '0' <= c && c <= '9':
		l.unget()
		return lex_decimalNumber
	case c == '%':
		// keep the % in the input; we read it out in the expression evaluation step
		return lex_binaryNumber
	case c == '$':
		// keep the $ in the input; we read it out in the expression evaluation step
		return lex_hexNumber
	case unicode.IsLetter(c) || c == '_' || c == '.':
		l.unget()
		return lex_ident
	case c == '\'':
		l.ignore()
		return lex_character
	case c == '"':
		l.ignore()
		return lex_string
	case c == ':':				// : or :: (terminator)
		l.acceptAndEmit(':', TERM, ':')
	case c == '&':				// & or &&
		l.acceptAndEmit('&', AND, '&')
	case c == '|':				// | or ||
		l.acceptAndEmit('|', OR, '|')
	case c == '=':				// = or ==
		l.acceptAndEmit('=', EQ, '=')
	case c == '!':				// ! or !=
		l.acceptAndEmit('=', NE, '!')
	case c == '<':				// < or <= or <<
		if l.accept('=') {
			l.emit(LE)
		} else if l.accept('<') {
			l.emit(LSH)
		} else {
			l.emit('<')
		}
	case c == '>':				// > or >= or >>
		if l.accept('=') {
			l.emit(GE)
		} else if l.accept('>') {
			l.emit(RSH)
		} else {
			l.emit('>')
		}
	// TODO more multi-character tokens
	default:
		l.emit(int(c))
	}
	return lex_next			// TODO tail call optimize?
}

func lex_decimalNumber(l *FileLexer) lexState {
	if !l.acceptRun("0123456789") {
		FATAL_BUG("lex_decimalNumber somehow failed to read a number (the only way it can be called is if we know there's a number already)")
	}
	l.emit(NUMBER)
	return lex_next
}

func lex_binaryNumber(l *FileLexer) lexState {
	l.acceptRun("01")
	// TODO handle error
	l.emit(NUMBER)
	return lex_next
}

func lex_hexNumber(l *FileLexer) lexState {
	l.acceptRun("0123456789ABCDEFabcdef")
	// TODO handle error
	l.emit(NUMBER)
	return lex_next
}

func lex_ident(l *FileLexer) lexState {
	for {
		c := l.read()
		if !unicode.IsLetter(c) && c != '_' && c != '.' {
			if c < '0' || c > '9' {
				break
			}
		}
	}
	l.unget()
	// TODO look up in local symbol table?
	s := Symbols.Get(l.what())
	if s != nil {
		l.emit(s.Type)
	} else {
		l.emit(IDENT)
	}
	return lex_next
}

func getStringCharacter(l *FileLexer) (r rune, isEscaped bool) {
	r = l.read()
	if r == '\\' {					// TODO have things like \u?
		isEscaped = true
		r = l.read()
		switch r {
		case 'n':
			r = '\n'
		// TODO more combinations
		}
		// default is to take that character literally
	}
	return
}

func lex_character(l *FileLexer) lexState {
	var count uint64

	// TODO worry about the length
	// TODO worry about allowing Unicode (meaning we have to worry about character encodings and ugh; probably best to worry about the length later
	for {
		r, isEscaped := getStringCharacter(l)
		if r == lexEOF {
			l.Error("EOF in character literal")
			return lex_next		// TODO
		}
		count++
		if !isEscaped && r == '\'' {
			l.unget()
			count--
			break
		}
	}
	if count == 0 {
		l.Error("empty character literal")
		return lex_next			// TODO
	} else if count > 4 {
		l.Error("character literal too long (max 4 characters)")
		return lex_next			// TODO
	}
	l.emit(CHARACTER)
	l.read(); l.ignore()				// TODO handle closing '
	return lex_next
}

func lex_string(l *FileLexer) lexState {
	for {
		r, isEscaped := getStringCharacter(l)
		if r == lexEOF {
			l.Error("EOF in string literal")
			return lex_next		// TODO
		}
		if !isEscaped && r == '"' {
			l.unget()
			break
		}
	}
	// TODO should we handle empty strings?
	l.emit(STRING)
	l.read(); l.ignore()				// TODO handle closing "
	return lex_next
}

func lex_end(l *FileLexer) lexState {
	l.emit(-1)				// EOF
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
	*tok = <-l.curfile.Tokens
	if tok.toktype == -1 {	// EOF
		l.EndFile()
		if len(l.files) <= 0 {	// no more files
			return -1
		}
		return l.Lex(tok)
	}
	return tok.toktype
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
		Closer:		f,
		Tokens:		make(chan yySymType),
	})
	l.curfile = l.files[len(l.files) - 1]
	go l.curfile.Run()
	return nil
}

func (l *Lexer) EndFile() {
	if l.curfile == nil || len(l.files) == 0 {
		FATAL_BUG("attempted to Lexer.EndFile() with no files open")
	}
	close(l.curfile.Tokens)		// TODO should we really do this here?
	l.curfile.Closer.Close()
	if len(l.files) > 1 {
		l.files = l.files[:len(l.files) - 1]
		l.curfile = l.files[len(l.files) - 1]
	} else {
		l.files = l.files[0:0]		// turn into nil slice (TOOD better way?)
		l.curfile = nil
	}
}
