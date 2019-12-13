// 13 december 2019
package cpu

type Operand interface {
	operand()
}

type DataRegisterOperand uint
func (DataRegisterOperand) operand() {}

type AddressRegisterOperand uint
func (AddressRegisterOperand) operand() {}

type AbsoluteWordOperand uint16
func (AbsoluteWordOperand) operand() {}

type AbsoluteLongOperand uint32
func (AbsoluteLongOperand) operand() {}

type PCRelativeWithOffsetOperand uint16
func (PCRelativeWithOffsetOperand) operand() {}

type IndexRegister uint
const (
	D0Word IndexRegister = iota
	D1Word
	D2Word
	D3Word
	D4Word
	D5Word
	D6Word
	D7Word
	A0Word
	A1Word
	A2Word
	A3Word
	A4Word
	A5Word
	A6Word
	A7Word
	D0Long
	D1Long
	D2Long
	D3Long
	D4Long
	D5Long
	D6Long
	D7Long
	A0Long
	A1Long
	A2Long
	A3Long
	A4Long
	A5Long
	A6Long
	A7Long
)

type PCRelativeWithIndexAndOffsetOperand struct {
	Index	IndexRegister
	Offset	uint8
}
func (PCRelativeWithIndexAndOffsetOperand) operand() {}

type AddressRegisterIndirectOperand uint
func (AddressRegisterIndirectOperand) operand() {}

type AddressRegisterIndirectPostincrementOperand uint
func (AddressRegisterIndirectPostincrementOperand) operand() {}

type AddressRegisterIndirectPredecrementOperand uint
func (AddressRegisterIndirectPredecrementOperand) operand() {}

type AddressRegisterIndirectWithOffsetOperand struct {
	Register	xxxx
	Offset	uint16
}
func (AddressRegisterIndirectWithOffsetOperand) operand() {}

type AddressRegisterIndirectWithIndexAndOffsetOperand struct {
	Register	xxxxx
	Index	IndexRegister
	Offset	uint8
}
func (AddressRegisterIndirectWithIndexAndOffsetOperand) operand() {}

type ImmediateOperand uint32
func (ImmediateOperand) operand() {}

type CCROperand struct{}
func (CCROperand) operand() {}

type SROperand struct{}
func (SROperand) operand() {}

type USPOperand struct{}
func (USPOperand) operand() {}

type MovemOperand uint16
func (MovemOperand) operand() {}
