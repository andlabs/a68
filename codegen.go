// 5 july 2012
package main

// TODO move to another file?
type Operand struct {
	Type				rune
	Reg				int
	Expr				*Expression	// d8/d16/immediates
	IndexReg			int
	IndexRegAddress	bool
	IndexRegLong		bool
}
