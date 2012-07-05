// 5 july 2012
package main

import (
	"io"
	"errors"
)

// BitWriter is the interface through which the code generation system is implemented.
// It allows me to write data in groups of bits rather than bytes, so I can write out the bit arrangement of the opcodes directly.
// I'll probably wind up extracting this and making it a standalone library in the future.

type BitWriter struct {
	underlying	io.Writer
	curByte		byte
	bitCount		int
	err			error
	pos			int64
}

// these errors are nonfatal
var (
	ErrSeekInsideByte		= errors.New("attempted to seek while inside a byte")
	ErrUnknownFailedWrite	= errors.New("actually writing a byte failed but no error was returned")
)

func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{
		underlying:	w,
	}
}

func (w *BitWriter) Write(p []byte) (int, error) {
	for i, x := range p {
		err := w.WriteBit(x)
		if err != nil {
			return i, err
		}
	}
	return len(p), nil
}

func (w *BitWriter) Pos() int64 {
	if w.bitCount != 0 {		// don't say anything if in the middle of a byte
		return -1
	}
	return w.pos
}

// this does the actual work
func (w *BitWriter) WriteBit(bit byte) error {
	if w.err != nil {			// don't continue after an error
		return w.err
	}
	w.curByte = (w.curByte << 1) | (bit & 1)
	w.bitCount++
	if w.bitCount == 8 {
		n, err := w.underlying.Write([]byte{w.curByte})
		if err != nil {
			w.err = err
			return err
		}
		if n != 1 {
			w.err = ErrUnknownFailedWrite
			return w.err
		}
		w.pos++
		w.bitCount = 0
	}
	return nil
}

func (w *BitWriter) InsideByte() bool {
	return w.bitCount != 0
}

// shorthand
func (w *BitWriter) WriteBits(bits ...byte) (int, error) {
	return w.Write(bits)
}

/* testing
func main() {
	b := new(bytes.Buffer)
	bits := NewBitWriter(b)
	bits.WriteBits(0,1,1,0,0,1,0,0, 0,1,1,0,1,1,1,0)
	fmt.Println(b)
	fmt.Println(bits.InsideByte())
	fmt.Println(bits.Pos())
	bits.WriteBit(0)
	fmt.Println(bits.Pos())
}
output:
dn
false
2
-1
*/
