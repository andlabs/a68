// 15 july 2012
package main

// lea <ea>,aN
func o_lea(suffix rune, src Operand, dest Operand) error {
	_ = suffix		// unused
	WriteBits(0, 1, 0, 0)
	WriteRegNum(dest.Reg)
	WriteBits(1)
	WriteBits(1, 1)
	WriteEANow(src)
	return nil
}

// TODO link
// TODO lsl/lsr
