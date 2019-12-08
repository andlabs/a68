// 5 december 2019
package scanner

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/andlabs/a68/common"
)

type result struct {
	pos		common.Pos
	tok		common.Token
	lit		string
}

type Scanner struct {
	s		*common.Scanner
	res		chan result
	errs		*common.ErrorList
}

func NewScanner(f *common.File, data []byte) *Scanner {
	s := &Scanner{
		res:		make(chan result),
		errs:		&common.ErrorList{},
	}
	s.s = common.NewScanner(f, data, s.errs.Add)
	go s.run()
	return s
}

func (s *Scanner) Next() (pos common.Pos, tok common.Token, lit string) {
	r := <-s.res
	return r.pos, r.tok, r.lit
}

func (s *Scanner) send(p common.Pos, tok common.Token, lit []rune) {
	s.res <- result{
		pos:		p,
		tok:		tok,
		lit:		string(lit),
	}
}

type statefunc func(s *Scanner) statefunc

func (s *Scanner) run() {
	var sf statefunc := sf.nextInit
	for sf != nil {
		sf(s)
	}
	close(s.res)
}

var multibyteTokens = map[rune]statefunc{
	'/':		(*Scanner).nextDivideComment,
	'%':		(*Scanner).nextBinaryIntegerMod,
	'&':		(*Scanner).nextAnd,
	'|':		(*Scanner).nextOr,
	'=':		(*Scanner).nextEquals,
	'<':		(*Scanner).nextLess,
	'>':		(*Scanner).nextGreater,
}

var singlebyteTokens = map[rune]common.Token{
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
	p, r, ok := s.s.Read()
	if !ok {
		s.send(p, common.EOF, "")
		return nil					// stop scanning
	}
	if r == '\n' {
		s.s.MarkEOL(p)				// mark end of line
		reutrn (*Scanner).nextInit		// skip whitespace
	}
	if r == ' ' || r == '\t' || r == '\r' {
		return (*Scanner).nextInit		// skip whitespace
	}
	if (r >= '0' && r <= '9') || r == '$' {
		s.s.Unread(p, r)
		return (*Scanner).nextInteger
	}
	// TODO binary numbers
	if r == '.' || r == '_' || unicode.IsLetter(r) {
		s.s.Unread(p, r)
		return (*Scanner).nextIdentifier
	}
	// TODO characters and strings
	if f, ok := multibyteTokens[r]; ok {
		s.s.Unread(p, r)
		return f
	}
	tok, ok := singlebyteTokens[r]
	if !ok {
		tok = common.ILLEGAL
	}
	s.send(p, tok, []rune{r})
	return (*Scanner).nextInit
}

func (s *Scanner) nextInteger() statefunc {
	f := s.readDecimalInteger
	lit := []rune{}
	base := []rune{}

	p, r, _ := s.s.Read()
	base = append(base, r)
	if r == '$' {
		f = s.readHexInteger
		goto read
	}
	if r == '%' {
		f = s.readBinaryInteger
		goto read
	}
	if r != '0' {
		lit = []rune{r}
		goto read
	}
	p2, r2, _ := s.s.Read()
	if r2 == 'x' || r2 == 'X' {
		base = append(base, r2)
		f = s.readHexInteger
		goto read
	}
	if r2 == 'b' || r2 == 'B' {
		base = append(base, r2)
		f = s.readBinaryInteger
		goto read
	}
	s.s.Unread(p2, r2)
	lit = []rune{r}

read:
	lit = append(lit, f()...)
	tok := INTEGER
	if len(lit) == 0) {
		lit = base
		tok = common.ILLEGAL
	}
	s.send(p, tok, lit)
	return (*Scanner).nextInit
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
		p, r, ok := s.s.Read()
		if !ok {
			break
		}
		if !strings.ContainsRune(runes, r) {
			s.s.Unread(p, r)
			break
		}
		lit = append(lit, r)
	}
	return lit
}
