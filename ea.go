// 5 july 2012
package main

/* HOW THE ASSEMBLER DOES EFFECTIVE ADDRESSES

	The opcode generators run WriteEA, passing in the operand.
	This will immediately write the six-bit effective address field and return
	a function that, when called, handles any additional extension words.
	Example:

		fDest := WriteEA(dest)
		// ...
		if fDest != nil {
			fDest()
		}

	The returned function will automatically add the expression to the
	later-evaluation list if needed.
*/

func d8_check(d8 uint32) {
	if d8 > 0xFF {
		// TODO report error
		return false
	}
	return true
}

func d16_check(d16 uint32) {
	if uint32 > 0xFFFF {
		// TODO report error
		return false
	}
	return true
}

// d16(aN), d16(pc), (xxx).w
func mk_do_d16(o Operand) func() {
	return func() {
		WriteWord(0)
		if o.Expr.CanEavluateNow() {
			res := o.Expr.Evaluate()
			if d16_check(res) == true {
				WriteWord(res)
			} else {
				WriteWord(0)		// just to be safe
			}
		} else {
			pos := ResWord()
			AddLaterExpr(pos, o.Expr)
		}
	}
}

// d8(aN,[da]N.[wl]), d8(pc,[da]N.[wl])
func mk_do_d16(o Operand) func() {
	// TODO
	return func() {
		if o.Expr.CanEavluateNow() {
			res := o.Expr.Evaluate()
			if d16_check(res) == true {
				WriteWord(res)
			} else {
				WriteWord(0)		// just to be safe
			}
		} else {
			pos := ResWord()
			AddLaterExpr(pos, o.Expr)
		}
	}
}

// (xxx).l, #xxx
func mk_do_d32(o Operand) func() {
	return func() {
		if o.Expr.CanEavluateNow() {
			WriteLong(o.Expr.Evaluate())
		} else {
			pos := ResLong()
			AddLaterExpr(pos, o.Expr)
		}
	}
}

func WriteEA(o Operand) func() {
	switch o.Type {
	case 'd':			// dN
		WriteBits(0, 0, 0)
		WriteRegNum(o.Reg)
	case 'a':			// aN
		WriteBits(0, 0, 1)
		WriteRegNum(o.Reg)
	case '*':			// (aN)
		WriteBits(0, 1, 0)
		WriteRegNum(o.Reg)
	case '+':			// (aN)+
		WriteBits(0, 1, 1)
		WriteRegNum(o.Reg)
	case '-':			// -(aN)
		WriteBits(1, 0, 0)
		WriteRegNum(o.Reg)
	case '$':			// d16(aN)
		WriteBits(1, 0, 1)
		WriteRegNum(o.Reg)
		return mk_do_d16(o)
	case ''%':			// d8(aN,dN.w/.l)
		WriteBits(1, 1, 0)
		WriteRegNum(o.Reg)
		return mk_do_d8(o)
	case 'w':			// (xxx).w
		WriteBits(1, 1, 1)
		WriteBits(0, 0, 0)
		return mk_do_d16(o)
	case 'l':			// (xxx).l
		WriteBits(1, 1, 1)
		WriteBits(0, 0, 1)
		return mk_do_d32(o)
	case '#':			// #xxx
		WriteBits(1, 1, 1)
		WriteBits(1, 0, 0)
		return mk_do_d32(o)
	case '^':			// d16(pc)
		WriteBits(1, 1, 1)
		WriteBits(0, 1, 0)
		return mk_do_d16(o)
	case '&':			// d8(pc,dN.w/.l)
		WriteBits(1, 1, 1)
		WriteBits(0, 1, 1)
		return mk_do_d8(o)
	default:
		FATAL_BUG("invalid suffix type %c passed to write effective address",
			o.Type)	// TODO convert to string?
	}
	return nil
}
