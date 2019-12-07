// 5 december 2019
package scanner

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/andlabs/a68/token"
)

// runeReader is like bytes.Reader except that it keeps track of the offset of the current rune.
type runeReader struct {
	b	[]byte
	pos	int
	r	rune
	n	int
}

func newRuneReader(b []byte) *runeReader {
	return &runeReader{
		b:	b,
	}
}

func (r *runeReader) next() bool {
	r.pos += r.n
	b := r.b[r.pos:]
	if len(b) == 0 {
		r.r = 0		// always return a rune of 0 after the end
		r.n = 0		// don't read past end of slice
		return false
	}
	r.r, r.n = utf8.DecodeRune(b)
	return true
}

func (r *runeReader) isValid() bool {
	if r.r != utf8.RuneError {
		return true
	}
	return r.n != 1
}

func (r *runeReader) off() int {
	return r.pos
}

func (r *runeReader) rune() rune {
	return r.r
}

func (r *runeReader) firstByte() byte {
	if r.pos >= len(r.b) {
		return 0
	}
	return r.b[r.pos]
}

type Scanner struct {
	ErrorCount	int
	handler		ErrorHandler

	cur			*token.File
	r			*runeReader

	unreadp		token.Pos
	unreadr		rune
}

func (s *Scanner) err(p token.Pos, format string, args ...interface{}) {
	s.ErrorCount++
	if s.handler != nil {
		s.handler(s.cur.Position(p), fmt.Sprintf(format, args...))
	}
}

func (s *Scanner) readrune() (p token.Pos, r rune, ok bool) {
	if s.unreadp != token.NoPos {
		p, r = s.unreadp, s.unreadr
		s.unreadp = token.NoPos
		s.unreadr = 0
		return p, r, true
	}
	for {
		ok = s.r.next()
		p = s.cur.Pos(s.r.off())
		if !ok || s.r.isValid() {
			break
		}
		// report error and try next byte
		s.err(p, "invalid byte 0x%X in UTF-8 stream", s.r.firstByte())
	}
	return p, s.r.rune(), ok
}

func (s *scanner) unreadrune(p token.Pos, r rune) {
	if s.unreadp != token.NoPos {
		panic("excessive unreading")
	}
	s.unreadp = p
	s.unreadr = r
}

func (s *scanner) send(p token.Pos, tok token.Token, lit []rune) {
	// TODO
}

type statefunc func(s *scanner) statefunc

var multibyteTokens = map[rune]statefunc{
	'/':		(*scanner).nextDivideComment,
	'%':		(*scanner).nextBinaryIntegerMod,
	'&':		(*scanner).nextAnd,
	'|':		(*scanner).nextOr,
	'=':		(*scanner).nextEquals,
	'<':		(*scanner).nextLess,
	'>':		(*scanner).nextGreater,
}

var singlebyteTokens = map[rune]token.Token{
	'(':		LPAREN,
	')':		RPAREN,
	'{':		LBRACE,
	'}':		RBRACE,
	'#':		IMMEDIATE,
	'+':		ADD,
	'-':		SUBTRACT,
	'*':		MULTIPLY,
	'^':		XOR,
	'~':		CMPL,
	'!':		NOT,
	',':		COMMA,
	';':		SEMICOLON,
	':':		COLON,
}

func (s *scanner) nextInit() statefunc {
	p, r, ok := s.readrune()
	if !ok {
		return nil					// stop scanning
	}
	if r == '\n' {
		s.cur.AddLine(s.cur.Offset(p))	// mark end of line
		reutrn (*scanner).nextInit		// skip whitespace
	}
	if r == ' ' || r == '\t' || r == '\r' {
		return (*scanner).nextInit		// skip whitespace
	}
	if r >= '0' && r <= '9' {
		s.unreadrune(p, r)
		return (*scanner).nextInteger
	}
	if r == '$' {
		return (*scanner).nextHexInteger
	}
	if r == '.' || r == '_' || unicode.IsLetter(r) {
		s.unreadrune(p, r)
		return (*scanner).nextIdentifier
	}
	// TODO characters and strings
	if f, ok := multibyteTokens[r]; ok {
		s.unreadrune(p, r)
		return f
	}
	tok, ok := singlebyteTokens[r]
	if !ok {
		tok = token.ILLEGAL
	}
	s.send(p, tok, []rune{r})
	return (*scanner).nextInit
}
