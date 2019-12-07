// 5 december 2019
package scanner

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/andlabs/a68/common"
)

type Token struct {
	Pos		common.Pos
	Tok		common.Token
	Lit		string
}

type Scanner struct {
	s		*common.Scanner
	tok		chan Token
	errs		*common.ErrorList
}

func NewScanner(f *common.File, data []byte) *Scanner {
	s := &Scanner{
		tok:		make(chan Token),
		errs:		&common.ErrorList{},
	}
	s.s = common.NewScanner(f, data, s.errs.Add)
	go s.run()
	return s
}

func (s *Scanner) Next() (tok Token, ok bool) {
	return <-s.tok
}

func (s *Scanner) send(p common.Pos, tok common.Token, lit []rune) {
	s.tok <- Token{
		Pos:		p,
		Tok:		tok,
		Lit:		string(lit),
	}
}

type statefunc func(s *Scanner) statefunc

func (s *Scanner) run() {
	var sf statefunc := sf.nextInit
	for sf != nil {
		sf(s)
	}
	close(s.tok)
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
	pr, ok := s.s.Read()
	if !ok {
		return nil					// stop scanning
	}
	if pr.Rune == '\n' {
		s.s.MarkEOL(pr.Pos)			// mark end of line
		reutrn (*Scanner).nextInit		// skip whitespace
	}
	if pr.Rune == ' ' || pr.Rune == '\t' || pr.Rune == '\r' {
		return (*Scanner).nextInit		// skip whitespace
	}
	if pr.Rune >= '0' && pr.Rune <= '9' {
		s.Unread(pr)
		return (*Scanner).nextInteger
	}
	if pr.Rune == '$' {
		return (*Scanner).nextHexInteger
	}
	if pr.Rune == '.' || pr.Rune == '_' || unicode.IsLetter(pr.Rune) {
		s.Unread(pr)
		return (*Scanner).nextIdentifier
	}
	// TODO characters and strings
	if f, ok := multibyteTokens[pr.Rune]; ok {
		s.Unread(pr)
		return f
	}
	tok, ok := singlebyteTokens[pr.Rune]
	if !ok {
		tok = common.ILLEGAL
	}
	s.send(pr.Pos, tok, []rune{pr.Rune})
	return (*Scanner).nextInit
}
