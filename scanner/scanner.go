// 5 december 2019
package scanner

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/andlabs/a68/token"
)

type result struct {
	pos		token.Pos
	tok		token.Token
	lit		string
}

type Scanner struct {
	r		*reader
	res		chan result
	errs		*ErrorList
}

func NewScanner(f *File, data []byte) *Scanner {
	s := &Scanner{
		res:		make(chan result),
		errs:		&ErrorList{},
	}
	s.r = newReader(f, data, s.errs.Add)
	go s.run()
	return s
}

func (s *Scanner) Next() (pos token.Pos, tok token.Token, lit string) {
	r := <-s.res
	return r.pos, r.tok, r.lit
}

func (s *Scanner) sendstr(off int, tok token.Token, lit string) {
	s.res <- result{
		pos:		s.r.pos(off),
		tok:		tok,
		lit:		lit,
	}
}

func (s *Scanner) send(off int, tok token.Token, lit []rune) {
	return s.sendstr(off, tok, string(lit))
}

type statefunc func(s *Scanner) statefunc

func (s *Scanner) run() {
	var sf statefunc = sf.next
	s.r.read()		// get things going
	for sf != nil {
		sf(s)
	}
	close(s.res)
}

var multibyteTokens = map[rune]statefunc{
}

var singlebyteTokens = map[rune]token.Token{
}

func (s *Scanner) next() statefunc {
	off, r := s.r.cur()
	if r == -1 {
		s.send(off, token.EOF, "")
		return nil					// stop scanning
	}
	if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
		s.r.read()					// skip whitespace
		return (*Scanner).next
	}
	if (r >= '0' && r <= '9') || r == '$' || r == '%' {
		return (*Scanner).nextInteger
	}
	if r == '.' || r == '_' || unicode.IsLetter(r) {
		return (*Scanner).nextIdentifier
	}
	// TODO characters and strings
	if f, ok := multibyteTokens[r]; ok {
		return f
	}
	tok, ok := singlebyteTokens[r]
	if !ok {
		tok = token.ILLEGAL
	}
	s.send(off, tok, []rune{r})
	s.r.read()
	return (*Scanner).next
}

func (s *Scanner) nextInteger() statefunc {
	lit := make([]rune, 0, 16)
	f := s.readDecimalInteger

	off, r := s.r.cur()
	lit = append(lit, r)
	if r == '$' {
		f = s.readHexInteger
		goto read
	}
	if r == '%' {
		f = s.readBinaryInteger
		goto read
	}
	if r != '0' {
		goto read
	}
	r := s.r.peekbyteasrune()
	if r == -1 {		// the last token of the file is a single 0
		goto send
	}
	if r == 'x' || r == 'X' {
		s.r.read()
		lit = append(lit, r)
		f = s.readHexInteger
		goto read
	}
	if r == 'b' || r == 'B' {
		s.r.read()
		lit = append(lit, r)
		f = s.readBinaryInteger
		goto read
	}

read:
	lit = append(lit, f()...)
send:
	s.send(off, INT, lit)
	return (*Scanner).next
}

func (s *Scanner) readBinaryInteger() []rune {
	return s.readStringOf("01")
}

func (s *Scanner) readDecimalInteger() []rune {
	return s.readStringOf("0123456789")
}

func (s *Scanner) readHexInteger() []rune {
	return s.readStringOf("0123456789ABCDEFabcdef")
}

func (s *Scanner) readStringOf(runes string) (lit []rune) {
	lit := make([]rune, 0, 8)
	for {
		_, r := s.r.read()
		if r == -1 || !strings.ContainsRune(runes, r) {
			break
		}
		lit = append(lit, r)
	}
	return lit
}

func (s *Scanner) nextIdentifier() statefunc {
	lit := make([]rune, 0, 16)
	off, r := s.r.cur()
	for r != -1 {
		if r != '.' && r != '_' && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			break
		}
		lit = append(lit, r)
		_, r = s.r.read()
	}
	strlit := string(lit)
	tok := token.Lookup(strlit)
	if tok == token.IDENT && lit[0] == '.' {
		s.r.err(firstp, "unknown keyword %q", strlit)
		return (*Scanner).next
	}
	s.sendstr(off, tok, strlit)
	return (*Scanner).next
}
