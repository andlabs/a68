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

// d8(aN,[da]N.[wl]), d8(pc,[da]N.[wl]) - these require special handling that the other modes do not, using the Brief Extension Word(correct phrase?) format
func mk_ea_d8(o Operand) func() {
	return func() {
		if o.IndexRegAddress {
			WriteBits(1)
		} else {
			WriteBits(0)
		}
		WriteRegNum(o.IndexReg)
		if o.IndexRegLong {
			WriteBits(1)
		} else {
			WriteBits(0)
		}
		WriteBits(0, 0)			// scale
		WriteBits(0)
		if o.Expr.CanEvaluateNow() {
			res := o.Expr.Evaluate()
			if d8_check(res) == true {
				WriteByte(byte(res))
			} else {
				WriteByte(0)			// just to be safe
			}
		} else {
			pos := ResByte()
			AddLaterExpr(pos, o.Expr, d8_check)
		}
	}
}

// let's just put this here for now
func AddLaterExpr(...interface{}){panic("TODO AddLaterExpr")}

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
		return WriteImmed_16(o)
	case '%':			// d8(aN,dN.w/.l)
		WriteBits(1, 1, 0)
		WriteRegNum(o.Reg)
		return mk_ea_d8(o)
	case 'w':			// (xxx).w
		WriteBits(1, 1, 1)
		WriteBits(0, 0, 0)
		return WriteImmed_16(o)
	case 'l':			// (xxx).l
		WriteBits(1, 1, 1)
		WriteBits(0, 0, 1)
		return WriteImmed_32(o)
	case '#':			// #xxx
		WriteBits(1, 1, 1)
		WriteBits(1, 0, 0)
		return WriteImmed_32(o)
	case '^':			// d16(pc)
		WriteBits(1, 1, 1)
		WriteBits(0, 1, 0)
		return WriteImmed_16(o)
	case '&':			// d8(pc,dN.w/.l)
		WriteBits(1, 1, 1)
		WriteBits(0, 1, 1)
		return mk_ea_d8(o)
	default:
		FATAL_BUG("invalid suffix type %c passed to write effective address",
			o.Type)	// TODO convert to string?
	}
	return nil
}

// shorthand
func WriteEANow(o Operand) {
	f := WriteEA(o)
	if f != nil {
		f()
	}
}
