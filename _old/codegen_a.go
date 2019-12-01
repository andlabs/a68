// 5 july 2012
package main

import "fmt"

// TODO move to another file?
type Operand struct {
	Type				rune
	Reg				int
	Expr				*Expression	// d8/d16/immediates
	IndexReg			int
	IndexRegAddress	bool
	IndexRegLong		bool
}

// abcd dN,dN
// abcd -(aN),-(aN)
func o_abcd(suffix rune, src Operand, dest Operand) error {
	_ = suffix		// unused
	if err := ochk_sameTypeOperands_bcdx(src, dest, "abcd"); err != nil {
		return err
	}
	WriteBits(1, 1, 0, 0)
	WriteRegNum(dest.Reg)
	WriteBits(1)
	WriteBits(0, 0, 0, 0)
	if src.Type == 'd' {		// R/M bit; dN is R
		WriteBits(0)
	} else {				// -(aN) is M
		WriteBits(1)
	}
	WriteRegNum(src.Reg)
	return nil
}

// add <ea>,dN
// add dN,<ea>
func o_add(suffix rune, src Operand, dest Operand) error {
	opmodes := Opmodes{
		{	// add <ea>,dN
			'b':	{ 0, 0, 0 },
			'w':	{ 0, 0, 1 },
			'l':	{ 0, 1, 0 },
		}, {	// add dN,<ea>
			'b':	{ 1, 0, 0 },
			'w':	{ 1, 0, 1 },
			'l':	{ 1, 1, 0 },
		},
	}

	if err := ochk_needOneDataReg(src, dest, "add"); err != nil {
		return err
	}
	if err := ochk_noByteReadFromAddrReg(suffix, src, "add"); err != nil {
		return err
	}
	WriteBits(1, 1, 0, 1)
	if dest.Type == 'd' {		// add <ea>,dN
		WriteRegNum(dest.Reg)
		WriteBits(opmodes[0][suffix]...)
		WriteEANow(src)
	} else {				// add dN,<ea>
		WriteRegNum(src.Reg)
		WriteBits(opmodes[1][suffix]...)
		WriteEANow(dest)
	}
	return nil
}

// adda <ea>,aN
func o_adda(suffix rune, src Operand, dest Operand) error {
	WriteBits(1, 1, 0, 1)
	WriteRegNum(dest.Reg)
	if suffix == 'w' {		// opmode
		WriteBits(0, 1, 1)
	} else {			// .l
		WriteBits(1, 1, 1)
	}
	WriteEANow(src)
	return nil
}

// addi #xxx,<ea>
func o_addi(suffix rune, src Operand, dest Operand) error {
	sizes := map[rune][]byte{
		'b':	{ 0, 0 },
		'w':	{ 0, 1 },
		'l':	{ 1, 0 },
	}

	WriteBits(0, 0, 0, 0)
	WriteBits(0, 1, 1, 0)
	WriteBits(sizes[suffix]...)
	fDest := WriteEA(dest)
	WriteImmediate(src, suffix)
	if fDest != nil {
		fDest()
	}
	return nil
}

// addq #xxx,<ea>
func o_addq(suffix rune, src Operand, dest Operand) error {
	sizes := map[rune][]byte{
		'b':	{ 0, 0 },
		'w':	{ 0, 1 },
		'l':	{ 1, 0 },
	}

	if err := ochk_noByteWriteToAddrReg(suffix, dest, "addq"); err != nil {
		return err
	}
	// MAJOR TODO
	if !src.Expr.CanEvaluateNow() {
		return fmt.Errorf("sorry, technical restrictions require arguments to addq be evaluatable at code generation time; this will be fixed later")
	}
	n := src.Expr.Evaluate()
	if n > 8 {
		return fmt.Errorf("addq immediate argument must be in the range 0 <= n <= 8; received $%X", n)
	}
	if n == 8 {
		n = 0
	}
	WriteBits(0, 1, 0, 1)
	WriteRegNum(int(n))		// just reuse this because it does what we want
	WriteBits(0)
	WriteBits(sizes[suffix]...)
	WriteEANow(dest)
	return nil
}

// addx dN,dN
// addx -(aN),-(aN)
func o_addx(suffix rune, src Operand, dest Operand) error {
	sizes := map[rune][]byte{
		'b':	{ 0, 0 },
		'w':	{ 0, 1 },
		'l':	{ 1, 0 },
	}

	if err := ochk_sameTypeOperands_bcdx(src, dest, "addx"); err != nil {
		return err
	}
	WriteBits(1, 1, 0, 1)
	WriteRegNum(dest.Reg)
	WriteBits(1)
	WriteBits(sizes[suffix]...)
	WriteBits(0, 0)
	if src.Type == 'd' {		// R/M bit; dN is R
		WriteBits(0)
	} else {				// -(aN) is M
		WriteBits(1)
	}
	WriteRegNum(src.Reg)
	return nil
}

// and <ea>,dN
// and dN,<ea>
func o_and(suffix rune, src Operand, dest Operand) error {
	opmodes := Opmodes{
		{	// add <ea>,dN
			'b':	{ 0, 0, 0 },
			'w':	{ 0, 0, 1 },
			'l':	{ 0, 1, 0 },
		}, {	// add dN,<ea>
			'b':	{ 1, 0, 0 },
			'w':	{ 1, 0, 1 },
			'l':	{ 1, 1, 0 },
		},
	}

	if err := ochk_needOneDataReg(src, dest, "and"); err != nil {
		return err
	}
	WriteBits(1, 1, 0, 0)
	if dest.Type == 'd' {		// add <ea>,dN
		WriteRegNum(dest.Reg)
		WriteBits(opmodes[0][suffix]...)
		WriteEANow(src)
	} else {				// add dN,<ea>
		WriteRegNum(src.Reg)
		WriteBits(opmodes[1][suffix]...)
		WriteEANow(dest)
	}
	return nil
}

// TODO andi

// TODO asl/asr
