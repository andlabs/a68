// 1 july 2012
package main

import (
	// ...
)

/* The format of an opcode table entry:

	Name - the name of an opcode without any suffixes
	Suffixes - list of suffixes as a string of characters
		(space)	no suffix
		b		.b or .s
		w		.w
		l		.l
	SrcTypes, DestTypes - the allowed types of an opcode, just like suffixes
		(space)	optional or unspecified
		d		dN
		a		aN
		*		(aN)
		+		(aN)+
		-		-(aN)
		$		d16(aN)
		%		d8(aN,dN.w/.l)
		^		d16(pc)
		&		d8(pc,dN.w/.l)
		w		(xxx).w
		l		(xxx).l
		#		#xxx
		c		ccr
		s		sr
		m		movem register list
	Generator - a function that, when called, will actually generate the given opcode

	For branching instructions, .w and .l are treated identically. I might need to think of a way to prevent saying
		bra.s (forwardLabel).l
*/
type Opcode struct {
	Name		string
	Suffixes		string
	SrcType		string
	DestType		string
	Generator		func(suffix string, src OpcodeArg, dest OpcodeArg) error
}

// GENERAL TODOs
// allow things like add xxx,a0 -> adda xxx,a0 implicitly
// aliases for branching, eor/xor, etc.

var Opcodes = [...]Opcode{
	{ "abcd", " ", "d-", "d-", o_abcd },
	{ "add", "bwl", "da*+-$%^&wl#", "d*+-$%wl", o_add },
	{ "adda", "wl", "da*+-$%^&wl#", "a", o_adda },
	{ "addi", "bwl", "#", "d*+-$%wl", o_addi },
	{ "addq", "bwl", "#", "da*+-$%wl", o_addq },
	{ "addx", "bwl", "d-", "d-", o_addx },
	{ "and", "bwl", "d*+-$%^&wl#", "d*+-$%wl", o_and },
	{ "andi", " bwl", "#", "d*+-$%wlc", o_andi },
//	{ "asl", " bwl", " d#", "d*+-$%wl", o_asl },
//	{ "asr", " bwl", " d#", "d*+-$%wl", o_asr },
	// TODO asl/asr <ea> suffixes?
	// TODO Bcc
	{ "bchg", " ", "d#", "d*+-$%wl", o_bchg },
	{ "bclr", " ", "d#", "d*+-$%wl", o_bclr },
	// newer CPUs: bfchg, bfclr, bfexts, bfextu, bfffo, bfins, bfset, bftst
	// TODO bkpt? that's in MC68EC000 but not MC68000
	{ "bra", "bw", " ", "wl", o_bra },
	{ "bset", " ", "d#", "d*+-$%wl", o_bset },
	{ "bsr", "bw", " ", "wl", o_bsr },
	[ "btst", " ", "d#", "d*+-$%^&wl#", o_btst },
	// newer CPUs: callm, cas, cas2
	// TODO chk suffixes?
	// newer CPUs: chk2
	{ "clr", "bwl", " ", "d*+-$%wl", o_clr },
	{ "cmp", "bwl", "da*+-$%^&wl#", "d", o_cmp },
	{ "cmpa", "wl", "da*+-$%^&wl#", "a", o_cmpa },
	{ "cmpi", "bwl", "#", "d*+-$%^&wl", o_cmpi },
	{ "cmpm", "bwl", "+", "+", o_cmpm },
	// newer CPUs: cmp2, cpBcc, cpDBcc cpGEN, cpScc, cpTRAPcc
	// TODO DBcc
	// TODO divs and divu suffixes?
	{ "eor", "bwl", "d", "d*+-$%wl", o_eor },
	{ "eori", " bwl", "#", "d*+-$%wlc", o_eori },
	{ "exg", " ", "da", "da", o_exg },
	{ "ext", "wl", " ", "d", o_ext },
	// newer CPUs: extb
	{ "illegal", " ", " ", " ", o_illegal },
	{ "jmp", " ", " ", "*$%^&wl", o_jmp },
	{ "jsr", " ", " ", "*$%^&wl", o_jsr },
	{ "lea", " ", "*$%^&wl", "a", o_lea },
	// TODO link suffixes?
	// TODO lsl/lsr <ea> suffixes?
	{ "move", " bwl", "da*+-$%^&wl#s", "d*+-$%wlc", o_move },
	{ "movea", "wl", "da*+-$%^&wl#", "a", o_movea },
	// newer CPUs: move from ccr, move16
	{ "movem", "wl", "*+$%^&wlm", "*-$%wlm", o_movem },		// TODO add this to the parser, and then see if I need to add d/a modes to handle lists consiting of a single register
	{ "movep", "wl", "d*$", "d*$", o_movep },		// slight breach of the rules here, but adding * allows me to elide the 0 in the case of 0(a0) — it'll be handled properly during encoding
	{ "moveq", " ", "#", "d", o_moveq },
	// TODO muls and mulu suffixes?
	{ "nbcd", " ", " ", "d*+-$%wl", o_nbcd },
	{ "neg", "bwl", " ", "d*+-$%wl", o_neg },
	{ "negx", "bwl", " ", "d*+-$%wl", o_negx },
	{ "nop", " ", " ", " ", o_nop },
	{ "not", "bwl", " ", "d*+-$%wl", o_not },
	{ "or", "bwl", "d*+-$%^&wl#", "d*+-$%wl", o_or },
	{ "ori", " bwl", "#", "d*+-$%wlc", o_ori },
	// newer CPUs: pack
	{ "pea", " ", " ", "*$%^&wl", o_pea },
	// TODO rol/ror/roxl/roxr <ea> suffixes?
	// newer CPUs: rtd, rtm
	{ "rtr", " ", " ", " ", o_rtr },
	{ "rts", " ", " ", " ", o_rts },
	{ "sbcd", " ", "d-", "d-", o_sbcd },
	// TODO Scc
	{ "sub", "bwl", "da*+-$%^&wl#", "d*+-$%wl", o_sub },
	{ "suba", "wl", "da*+-$%^&wl#", "a", o_suba },
	{ "subi", "bwl", "#", "d*+-$%wl", o_subi },
	{ "subq", "bwl", "#", "da*+-$%wl", o_subq },
	{ "subx", "bwl", "d-", "d-", o_subx },
	{ "swap", " ", " ", "d", o_swap },
	{ "tas", " ", " ", "d*+-$%wl", o_tas },		// yeah, tas on a data register is legal; I have no idea why
	{ "trap", " ", " ", "#", o_trap },
	// newer CPUs: TRAPcc
	{ "trapv", " ", " ", " ", o_trapv },
	{ "tst", "bwl", " ", "da*+-$%^&wl#", o_tst },
	{ "unlk", " ", " ", "a", o_unlk },
	// newer CPUs: unpk
}
