// 5 july 2012
package main

import (
	"bytes"
)

var (
	Output	*bytes.Buffer
	OutBits	*BitWriter
)

func init() {
	Output = new(bytes.Buffer)
	OutBits = NewBitWriter(Output)
}

func WriteBits(bits ...byte) {
	n, err := OutBits.Write(bits)
	if err != nil {
		FATAL_BUG("write bits %v into output buffer failed: %v",
			bits, err)
	}
	if n != len(bits) {
		FATAL_BUG("short write (%d<%d) of bits %v into output buffer with no explanation (this means something's wrong in BitWriter as I have a special error for that)",
			n, len(bits), bits, err)
	}
}

// TODO make part of BitWriter?
func WriteByte(b byte) {
	if OutBits.InsideByte() {		// sanity check? TODO
		FATAL_BUG("attempt to write byte $%X resulted in write inside a byte",
			b)
	}
	WriteBits(b >> 7,
		b >> 6,
		b >> 5,
		b >> 4,
		b >> 3,
		b >> 2,
		b >> 1,
		b & 1)
}

func WriteWord(w uint16) {
	WriteByte(byte((w >> 8) & 0xFF))
	WriteByte(byte(w & 0xFF))
}

func WriteLong(l uint32) {
	WriteWord(uint16((l >> 16) & 0xFFFF))
	WriteWord(uint16(l & 0xFFFF))
}

func ResByte() int64 {
	pos := OutBits.Pos()
	if pos == -1 {			// sanity check
		FATAL_BUG("attempted to reserve space for a byte in the middle of a byte")
	}
	WriteByte(0)
	return pos
}

func ResWord() int64 {
	pos := OutBits.Pos()
	if pos == -1 {			// sanity check
		FATAL_BUG("attempted to reserve space for a word in the middle of a byte")
	}
	WriteWord(0)
	return pos
}

func ResLong() int64 {
	pos := OutBits.Pos()
	if pos == -1 {			// sanity check
		FATAL_BUG("attempted to reserve space for a long in the middle of a byte")
	}
	WriteLong(0)
	return pos
}

// now for the 68000-specific writers

func WriteRegNum(num int) {
	if num < 0 || num > 7 {		// sanity check
		FATAL_BUG("invalid register number %d written", num)
	}
	r := byte(num & 7)
	WriteBits(r >> 2,
		r >> 1,
		r & 1)
}
