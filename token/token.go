// 2 december 2019
package token

// Token is the set of lexical tokens of the assembler.
type Token int
const (
	// Special tokens.
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Primitives.
	IDENT
	OPCODE			// abcd, add.b, ...
	DATADEF			// dc.[bwl]
	DATAREG			// d0 .. d7
	ADDRREG			// a0 .. a7, sp
	DATAREG_W		// d0.w .. d7.w
	ADDRREG_W		// a0.w .. a7.w, sp.w
	DATAREG_L		// d0.l .. d7.l
	ADDRREG_L		// a0.l .. a7.l, sp.l
	PC				// pc
	SPECIALREG		// ccr, sr, usp
	INTEGER
	CHAR			// 'x'
	STRING			// "x"

	// Addressing mode modifiers.
	SUFFIX_B			// .b (not used, but reserved)
	SUFFIX_W			// .w
	SUFFIX_L			// .l

	// Operators and delimiters.
	ADD				// +
	SUB				// -
	MUL				// *
	QUO				// /
	REM				// %

	AND				// &
	OR				// |
	XOR				// ^
	SHL				// <<
	SHRL			// >>
	SHRA			// >>>
	CMPL			// ~

	LAND			// &&
	LOR				// ||
	NOT				// !

	EQL				// ==
	NEQ				// !=
	LSS				// <
	LEQ				// <=
	GTR				// >
	GEQ				// >=

	COMMA			// ,
	COLON			// :
	HASH			// #

	LPAREN			// (
	RPAREN			// )

	// TODO keywords
)
