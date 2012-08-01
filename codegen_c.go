// 9 july 2012
package main

// TODO chk

// clr <ea>
func o_clr(suffix rune, src Operand, dest Operand) error {
	sizes := map[rune][]byte{
		'b':	{ 0, 0 },
		'w':	{ 0, 1 },
		'l':	{ 1, 0 },
	}

	_ = src		// unused
	WriteBits(0, 1, 0, 0)
	WriteBits(0, 0, 1, 0)
	WriteBits(sizes[suffix]...)
	WriteEANow(dest)
	return nil
}

// cmp <ea>,dN
func o_cmp(suffix rune, src Operand, dest Operand) error {
	opmodes := map[rune][]byte{		// not Opmodes because there's only one setting
		'b':	{ 0, 0, 0 },
		'w':	{ 0, 0, 1 },
		'l':	{ 0, 1, 0 },
	}

	if err := ochk_noByteReadFromAddrReg(suffix, src, "cmp"); err != nil {
		return err
	}
	WriteBits(1, 0, 1, 1)
	WriteRegNum(dest.Reg)
	WriteBits(opmodes[suffix]...)
	WriteEANow(src)
	return nil
}

// cmpa <ea>,aN
func o_cmpa(suffix rune, src Operand, dest Operand) error {
	opmodes := map[rune][]byte{
		'w':	{ 0, 1, 1 },
		'l':	{ 1, 1, 1 },
	}

	WriteBits(1, 0, 1, 1)
	WriteRegNum(dest.Reg)
	WriteBits(opmodes[suffix]...)
	WriteEANow(src)
	return nil
}

// cmpi #xxx,<ea>
func o_cmpi(suffix rune, src Operand, dest Operand) error {
	sizes := map[rune][]byte{
		'b':	{ 0, 0 },
		'w':	{ 0, 1 },
		'l':	{ 1, 0 },
	}

	WriteBits(0, 0, 0, 0)
	WriteBits(1, 1, 0, 0)
	WriteBits(sizes[suffix]...)
	fDest := WriteEA(dest)
	WriteImmediate(src, suffix)
	if fDest != nil {
		fDest()
	}
	return nil
}

// cmpm (aN)+,(aN)+
func o_cmpm(suffix rune, src Operand, dest Operand) error {
	sizes := map[rune][]byte{
		'b':	{ 0, 0 },
		'w':	{ 0, 1 },
		'l':	{ 1, 0 },
	}

	WriteBits(1, 0, 1, 1)
	WriteRegNum(dest.Reg)
	WriteBits(1)
	WriteBits(sizes[suffix]...)
	WriteBits(0, 0)
	WriteBits(1)
	WriteRegNum(src.Reg)
	return nil
}
