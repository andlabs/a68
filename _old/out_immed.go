// 8 july 2012
package main

// this function is also used by out_ea.go for its mk_ea_d8() function (and probably needs a differnet name for that case)
func d8_check(d8 uint32) bool {
	if d8 > 0xFF {
		// TODO report error
		return false
	}
	return true
}

func d16_check(d16 uint32) bool {
	if d16 > 0xFFFF {
		// TODO report error
		return false
	}
	return true
}

func WriteImmed_8(o Operand) func() {
	return func() {
		WriteByte(0)				// must align to 16 bits
		if o.Expr.CanEvaluateNow() {
			res := o.Expr.Evaluate()
			if d8_check(res) == true {
				WriteByte(byte(res))
			} else {
				WriteByte(0)		// just to be safe
			}
		} else {
			pos := ResByte()
			AddLaterExpr(pos, o.Expr, d8_check)
		}
	}
}

func WriteImmed_16(o Operand) func() {
	return func() {
		if o.Expr.CanEvaluateNow() {
			res := o.Expr.Evaluate()
			if d16_check(res) == true {
				WriteWord(uint16(res))
			} else {
				WriteWord(0)		// just to be safe
			}
		} else {
			pos := ResWord()
			AddLaterExpr(pos, o.Expr, d16_check)
		}
	}
}

func WriteImmed_32(o Operand) func() {
	return func() {
		if o.Expr.CanEvaluateNow() {
			WriteLong(o.Expr.Evaluate())
		} else {
			pos := ResLong()
			AddLaterExpr(pos, o.Expr, nil)
		}
	}
}

func WriteImmediate(s Operand, suffix rune) {
	if s.Type != '#' {		// sanity check
		// TODO write the operand?
		FATAL_BUG("attempted to write non-immediate in WriteImmediate")
	}
	switch suffix {
	case 'b':
		WriteImmed_8(s)()		// call them now
		return
	case 'w':
		WriteImmed_16(s)()
		return
	case 'l':
		WriteImmed_32(s)()
		return
	}
	FATAL_BUG("attempted to write immediate with invalid suffix '%c'", suffix)
	panic("FATAL_BUG returned")
}
