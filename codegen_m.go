// 30 july 2012
package main

// movea <ea>,aN
func o_movea(suffix rune, src Operand, dest Operand) error {
	sizes := map[rune][]byte{
		'w':	{ 1, 1 },
		'l':	{ 1, 0 },
	}

	WriteBits(0, 0)
	WriteBits(sizes[suffix]...)
	WriteRegNum(dest.Reg)
	WriteBits(0)
	WriteBits(0, 1)
	WriteEANow(src)
	return nil
}

// TODO movem
// TODO movep
// TODO moveq
// TODO muls
// TODO mulu
