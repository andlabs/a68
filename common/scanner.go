// 5 december 2019
package common

import (
	"fmt"
	"unicode"
	"unicode/utf8"
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

type PosRune struct {
	Pos		Pos
	Rune		rune
}

type Scanner struct {
	ErrorCount	int
	handler		ErrorHandler

	f			*File
	r			*runeReader

	unread		[]PosRune
}

func NewScanner(f *File, data []byte, handler ErrorHandler) *Scanner {
	if f.Size() != len(data) {
		panic(fmt.Sprintf("size mismatch in NewScanner(): file size %d != data size %d", f.Size(), len(data))
	}
	return &Scanner{
		handler:		handler,
		f:			f,
		r:			newRuneReader(data),
		unread:		make([]PosRune, 0, 16),
	}
}

func (s *Scanner) Err(p token.Pos, format string, args ...interface{}) {
	s.ErrorCount++
	if s.handler != nil {
		s.handler(s.f.Position(p), fmt.Sprintf(format, args...))
	}
}

func (s *Scanner) Read() (pr PosRune, ok bool) {
	if len(s.unread) != 0 {
		i := len(s.unread) - 1
		pr = s.unread[i]
		s.unread = s.unread[:i]
		return pr, true
	}
	for {
		ok = s.r.next()
		pr.Pos = s.f.Pos(s.r.off())
		pr.Rune = s.r.rune()
		if !ok || s.r.isValid() {
			break
		}
		// report error and try next byte
		s.Err(pr.Pos, "invalid byte 0x%X in UTF-8 stream", s.r.firstByte())
	}
	return pr, ok
}

func (s *Scanner) Unread(pr PosRune) {
	s.unread = append(s.unread, pr)
}

func (s *Scanner) MarkEOL(p Pos) {
	s.f.AddLine(s.f.Offset(p))
}
