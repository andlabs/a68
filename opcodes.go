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
		(space)	optional
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
		an empty string means no argument
	Generator - a function that, when called, will actually generate the given opcode
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

var Opcodes = [...]Opcode{
	{ "abcd", " ", "d-", "d-", o_abcd },
	{ "add", "bwl", "da*+-$%^&wl#", "d*+-$%wl", o_add },
	{ "adda", "wl", "da*+-$%^&wl#", "a", o_adda },
	{ "addi", "bwl", "#", "d*+-$%wl", o_addi },
	{ "addq", "bwl", "#", "da*+-$%wl", o_addq },
	{ "addx", "bwl", "d-", "d-", o_addx },
	{ "and", "bwl", "d*+-$%^&wl#", "d*+-$%wl", o_and },
	{ "andi", " bwl", "#", "d*+-$%wlc", o_andi },
	{ "asl", " bwl", " d#", "d*+-$%wl", o_asl },
	{ "asr", " bwl", " d#", "d*+-$%wl", o_asr },
}