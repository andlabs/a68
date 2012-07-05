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
	underlying	io.WriteSeeker
	curByte		byte
	bitCount		int
	err			error
}

// these errors are nonfatal
var (
	ErrSeekInsideByte		= errors.New("attempted to seek while inside a byte")
	ErrUnknownFailedWrite	= errors.New("actually writing a byte failed but no error was returned")
)

func New(w io.WriteSeeker) *BitWriter {
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

func (w *BitWriter) Seek(offset int64, whence int) (ret int64, err error) {
	if w.bitCount != 0 {		// don't let us seek until we've written out a full byte
		return 0, ErrSeekInsideByte
	}
	return w.underlying.Seek(offset, whence)
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
