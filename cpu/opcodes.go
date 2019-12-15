// 12 december 2019
package cpu

type Opcode interface {
	Name() string
/*TODO
	ValidSuffix(suffix string) bool
	NumOperands() int
	ValidOperand(operand Operand, which int) bool
*/
	// TODO opcode encoding
}

var Opcodes = []Opcode{
/*TODO
	Abcd{},
	Add{},
	Adda{},
	Addi{},
	Addq{},
	Addx{},
	And{},
	Andi{},
	Asl{},
	Asr{},
	Bcc{},
	Bchg{},
	Bclr{},
	Bcs{},
	Beq{},
	Bge{},
	Bgt{},
	Bhi{},
	Ble{},
	Bls{},
	Blt{},
	Bmi{},
	Bne{},
	Bpl{},
	Bra{},
	Bset{},
	Bsr{},
	Btst{},
	Bvc{},
	Bvs{},
	Chk{},
	Clr{},
	Cmp{},
	Cmpa{},
	Cmpi{},
	Cmpm{},
	Dbcc{},
	Dbcs{},
	Dbeq{},
	Dbf{},
	Dbge{},
	Dbgt{},
	Dbhi{},
	Dble{},
	Dbls{},
	Dblt{},
	Dbmi{},
	Dbne{},
	Dbpl{},
	Dbt{},
	Dbvc{},
	Dbvs{},
	Divs{},
	Divu{},
	Eor{},
	Eori{},
	Exg{},
	Ext{},
	Illegal{},		// $4AFC specifically
	Jmp{},
	Jsr{},
	Lea{},
	Link{},
	Lsl{},
	Lsr{},
	Move{},
	Movea{},
	Movem{},
	Movep{},
	Moveq{},
	Muls{},
	Mulu{},
	Nbcd{},
	Neg{},
	Negx{},
	Nop{},
	Not{},
	Or{},
	Ori{},
	Pea{},
	Reset{},
	Rol{},
	Ror{},
	Roxl{},
	Roxr{},
	Rte{},
	Rtr{},
	Rts{},
	Sbcd{},
	Scc{},
	Scs{},
	Seq{},
	Sf{},
	Sge{},
	Sgt{},
	Shi{},
	Sle{},
	Sls{},
	Slt{},
	Smi{},
	Sne{},
	Spl{},
	St{},
	Stop{},
	Sub{},
	Suba{},
	Subi{},
	Subq{},
	Subx{},
	Svc{},
	Svs{},
	Swap{},
	Tas{},
	Trap{},
	Trapv{},
	Tst{},
	Unlk{},
*/
}
