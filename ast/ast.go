// 1 december 2019
package ast

import (
	"github.com/andlabs/a68/token"
)

// All node types implement the Node interface.
type Node interface {
	Pos() token.Pos		// Position of the first character of the node.
	End() token.Pos		// Position of the first character after the node.
}

// All instruction operands implement the Operand interface.
type Operand interface {
	Node
	operand()
}

// A DataRegisterDirectOperand node represents a data register as an instruction operand.
// TODO(andlabs): register aliases
type DataRegisterDirectOperand struct {
	TokPos	token.Pos
	Tok		token.Token	// D[0..7]
	NextPos	token.Pos
}

func (d *DataRegisterDirectOperand) Pos() token.Pos { return d.TokPos }
func (d *DataRegisterDirectOperand) End() token.Pos { return d.EndPos }
func (*DataRegisterDirectOperand) operand() {}

// An AddressRegisterDirectOperand node represents an address register as an instruction operand.
// TODO(andlabs): register aliases
type AddressRegisterDirectOperand struct {
	TokPos	token.Pos
	Tok		token.Token	// A[0..7], SP
	NextPos	token.Pos
}

func (a *AddressRegisterDirectOperand) Pos() token.Pos { return a.TokPos }
func (a *AddressRegisterDirectOperand) End() token.Pos { return a.EndPos }
func (*AddressRegisterDirectOperand) operand() {}

// An AbsoluteWordDataOperand node represents an instruction operand of the form (xxx).w.
type AbsoluteWordDataOperand struct {
	Lparen	token.Pos
	X		Expr
	Rparen	token.Pos
	DotWPos	token.Pos
}

func (a *AbsoluteWordDataOperand) Pos() token.Pos { return a.Lparen }
func (a *AbsoluteWordDataOperand) End() token.Pos { return a.DotWPos + 2 }
func (*AbsoluteWordDataOperand) operand() {}

// An AbsoluteLongDataOperand node represents an instruction operand of the form (xxx).l.
type AbsoluteLongDataOperand struct {
	Lparen	token.Pos
	X		Expr
	Rparen	token.Pos
	DotLPos	token.Pos
}

func (a *AbsoluteLongDataOperand) Pos() token.Pos { return a.Lparen }
func (a *AbsoluteLongDataOperand) End() token.Pos { return a.DotLPos + 2 }
func (*AbsoluteLongDataOperand) operand() {}

// A PCRelativeWithOffsetOperand represents an instruction operand of the form expr(pc).
type PCRelativeWithOffsetOperand struct {
	X		Expr			// or nil if no offset given
	Lparen	token.Pos
	PCPos	token.Pos
	Rparen	token.Pos
}

func (a *PCRelativeWithOffsetOperand) Pos() token.Pos {
	if a.X != nil {
		return a.X.Pos()
	}
	return a.Lparen
}
func (a *PCRelativeWithOffsetOperand) End() token.Pos { return a.Rparen + 1 }
func (*PCRelativeWithOffsetOperand) operand() {}

// A PCRelativeWithIndexAndOffsetOperand represents an instruction operand of the form expr(pc,reg.size).
type PCRelativeWithIndexAndOffsetOperand struct {
	X		Expr			// or nil if no offset given
	Lparen	token.Pos
	PCPos	token.Pos
	Comma	token.Pos
	// TODO register spec
	Rparen	token.Pos
}

func (a *PCRelativeWithIndexAndOffsetOperand) Pos() token.Pos {
	if a.X != nil {
		return a.X.Pos()
	}
	return a.Lparen
}
func (a *PCRelativeWithIndexAndOffsetOperand) End() token.Pos { return a.Rparen + 1 }
func (*PCRelativeWithIndexAndOffsetOperand) operand() {}

// An AddressRegisterIndirectOperand node represents an instruction operand of the form (aX).
// TODO(andlabs): register aliases
type AddressRegisterIndirectOperand struct {
	Lparen	token.Pos
	TokPos	token.Pos
	Tok		token.Token	// A[0..7], SP
	Rparen	token.Pos
}

func (a *AddressRegisterIndirectOperand) Pos() token.Pos { return a.Lparen }
func (a *AddressRegisterIndirectOperand) End() token.Pos { return a.Rparen + 1 }
func (*AddressRegisterIndirectOperand) operand() {}

// An AddressRegisterPostincrementOperand node represents an instruction operand of the form (aX)+.
// TODO(andlabs): register aliases
type AddressRegisterPostincrementOperand struct {
	Lparen	token.Pos
	TokPos	token.Pos
	Tok		token.Token	// A[0..7], SP
	Rparen	token.Pos
	Plus		token.Pos
}

func (a *AddressRegisterPostincrementOperand) Pos() token.Pos { return a.Lparen }
func (a *AddressRegisterPostincrementOperand) End() token.Pos { return a.Plus + 1 }
func (*AddressRegisterPostincrementOperand) operand() {}

// An AddressRegisterPredecrementOperand node represents an instruction operand of the form -(aX).
// TODO(andlabs): register aliases
type AddressRegisterPredecrementOperand struct {
	Minus	token.Pos
	Lparen	token.Pos
	TokPos	token.Pos
	Tok		token.Token	// A[0..7], SP
	Rparen	token.Pos
}

func (a *AddressRegisterPredecrementOperand) Pos() token.Pos { return a.Minus }
func (a *AddressRegisterPredecrementOperand) End() token.Pos { return a.Rparen + 1 }
func (*AddressRegisterPredecrementOperand) operand() {}

// An AddressRegisterWithOffsetOperand represents an instruction operand of the form expr(aX).
// TODO(andlabs): register aliases
type AddressRegisterWithOffsetOperand struct {
	X		Expr			// must not be nil this time
	Lparen	token.Pos
	TokPos	token.Pos
	Tok		token.Token	// A[0..7], SP
	Rparen	token.Pos
}

func (a *AddressRegisterWithOffsetOperand) Pos() token.Pos { return a.X.Pos() }
func (a *AddressRegisterWithOffsetOperand) End() token.Pos { return a.Rparen + 1 }
func (*AddressRegisterWithOffsetOperand) operand() {}

// An AddressRegisterWithIndexAndOffsetOperand represents an instruction operand of the form expr(aX,reg.size).
// TODO(andlabs): register aliases
type AddressRegisterWithIndexAndOffsetOperand struct {
	X		Expr			// or nil if no offset given
	Lparen	token.Pos
	TokPos	token.Pos
	Tok		token.Token	// A[0..7], SP
	Comma	token.Pos
	// TODO register spec
	Rparen	token.Pos
}

func (a *AddressRegisterWithIndexAndOffsetOperand) Pos() token.Pos {
	if a.X != nil {
		return a.X.Pos()
	}
	return a.Lparen
}
func (a *AddressRegisterWithIndexAndOffsetOperand) End() token.Pos { return a.Rparen + 1 }
func (*AddressRegisterWithIndexAndOffsetOperand) operand() {}

// An ImmediateDataOperand node represents an instruction operand of the form #expr.
type ImmediateDataOperand struct {
	HashSign	token.Pos
	X		Expr
}

func (i *ImmediateDataOperand) Pos() token.Pos { return a.HashSign }
func (i *ImmediateDataOperand) End() token.Pos { return a.X.End() }
func (*ImmediateDataOperand) operand() {}

// An ImpliedAddressingOperand node represents a special register as an instruction operand.
// (This name was decided by Motorola, not me.)
// TODO(andlabs): register aliases
type ImpliedAddressingOperand struct {
	TokPos	token.Pos
	Tok		token.Token	// CCR, SR, USP
	NextPos	token.Pos
}

func (i *ImpliedAddressingOperand) Pos() token.Pos { return i.TokPos }
func (i *ImpliedAddressingOperand) End() token.Pos { return i.EndPos }
func (*ImpliedAddressingOperand) operand() {}

// TODO movem register list
