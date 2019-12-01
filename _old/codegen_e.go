// 15 july 2012
package main

import (
//	"fmt"
)

// eor dN,<ea>
func o_eor(suffix rune, src Operand, dest Operand) error {
	opmodes := map[rune][]byte{
		'b':	{ 1, 0, 0 },
		'w':	{ 1, 0, 1 },
		'l':	{ 1, 1, 0 },
	}

	WriteBits(1, 0, 1, 1)
	WriteRegNum(src.Reg)
	WriteBits(opmodes[suffix]...)
	WriteEANow(dest)
	return nil
}

// TODO eori

// exg dN,dN
// exg aN,aN
// exg dN,aN
func o_exg(suffix rune, src Operand, dest Operand) error {
	_ = suffix		// unused
	WriteBits(1, 1, 0, 0)
	switch {
	case src.Type == 'd' && dest.Type == 'd':		// exg dN,dN
		WriteRegNum(src.Reg)
		WriteBits(1)
		WriteBits(0, 1, 0, 0, 0)		// opmode
		WriteRegNum(dest.Reg)
	case src.Type == 'a' && dest.Type == 'a':		// exg aN,aN
		WriteRegNum(src.Reg)
		WriteBits(1)
		WriteBits(0, 1, 0, 0, 1)		// opmode
		WriteRegNum(dest.Reg)
	case src.Type == 'a' && dest.Type == 'd':		// exg aN,dN - allow this as an alternative to exg dN,aN (TODO or should I not)
		src, dest = dest, src
		fallthrough
	case src.Type == 'd' && dest.Type == 'a':		// exg dN,aN
		WriteRegNum(src.Reg)
		WriteBits(1)
		WriteBits(1, 0, 0, 0, 1)		// opmode
		WriteRegNum(dest.Reg)
	default:								// sanity check
		FATAL_BUG("exg opcode on invalid source/destination types %c/%c",
			src.Type, dest.Type)
	}
	return nil
}

// ext.w dN
// ext.l dN
func o_ext(suffix rune, src Operand, dest Operand) error {
	opmodes := map[rune][]byte{
		'w':	{ 0, 1, 0 },
		'l':	{ 0, 1, 1 },
	}

	_ = src		// unused
	WriteBits(0, 1, 0, 0)
	WriteBits(1, 0, 0)
	WriteBits(opmodes[suffix]...)
	WriteBits(0, 0)
	WriteBits(0)
	WriteRegNum(dest.Reg)
	return nil
}
