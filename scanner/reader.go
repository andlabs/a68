// 5 december 2019
package scanner

import (
	"fmt"
	"unicode/utf8"

	"github.com/andlabs/a68/token"
)

type reader struct {
	f	*token.File
	b	[]byte
	off	int
	r	rune
	n	int

	errorCount	int
	handler		ErrorHandler
}

func newReader(f *File, data []byte, handler ErrorHandler) *reader {
	return &reader{
		f:			f,
		b:			data,
		handler:		handler,
	}
}

func (r *reader) err(off int, format string, args ...interface{}) {
	r.errorCount++
	if r.handler != nil {
		r.handler(r.f.Position(r.f.Pos(off)), fmt.Sprintf(format, args...))
	}
}

func (r *reader) read() (off int, ru rune) {
	r.off += r.n
	if r.off < len(r.b) {
		r.r, r.n = utf8.DecodeRune(r.b[r.off:])
		if r.r == utf8.RuneError && r.n == 1 {
			r.err(r.off, "invalid byte 0x%X in UTF-8 stream", r.b[r.off])
			// return the utf8.RuneError anyway, so as to not allow invalid UTF-8 mid-token
		} else if r.r == '\n' {
			r.f.AddLine(r.off + r.n)
		}
		return r.off, r.r
	}
	r.n = 0		// don't advance r.off past len(r.b)
	r.r = -1
	return r.off, r.r
}

func (r *reader) cur() (off int, ru rune) {
	return r.off, r.r
}

func (r *reader) peekbyteasrune() rune {
	off := r.off + r.n
	if off >= len(r.b) {
		return -1
	}
	return rune(r.b[off])
}

func (r *reader) pos(off int) token.Pos {
	return r.f.Pos(off)
}
