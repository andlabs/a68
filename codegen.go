// 31 july 2012
package main

import "fmt"

// this stores stuff common to the code generator

type Opmodes [2]map[rune][]byte		// WriteBits(opmodes[read/write][suffix]...)
// TODO are all Opmodes instances the same?
// TODO same for suffixes

// for the BCD and X opcodes
func ochk_sameTypeOperands_bcdx(src Operand, dest Operand, opcode string) error {
	if src.Type != dest.Type {
		return fmt.Errorf("%s operand types must be the same (either both dN or both -(aN))", opcode)
	}
	return nil
}

// for arithmetic and logical opcodes
func ochk_needOneDataReg(src Operand, dest Operand, opcode string) error {
	if src.Type != 'd' && dest.Type != 'd' {		// at least one operand must be a data register
		// TODO print more information?
		return fmt.Errorf("at least one operand of %s must be a data register", opcode)
	}
	return nil
}

// for those opcodes which support both address register sources and byte operations but not both at the same time
func ochk_noByteReadFromAddrReg(suffix rune, src Operand, opcode string) error {
	if src.Type == 'a' && suffix == 'b' {		// no byte reads from address registers
		// TODO print the register?
		return fmt.Errorf("%s.b cannot be used with an address register source", opcode)
	}
	return nil
}

// for those opcodes which support both address register destinations and byte operations but not both at the same time
func ochk_noByteWriteToAddrReg(suffix rune, dest Operand, opcode string) error {
	if dest.Type == 'a' && suffix == 'b' {		// no byte writes to address registers
		// TODO print the register?
		return fmt.Errorf("%s.b cannot be used with an address register destination", opcode)
	}
	return nil
}
