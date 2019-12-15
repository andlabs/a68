// 4 december 2019
package token

import (
	"sync"
)

// Token is the set of lexical tokens in the assembler.
type Token int

// Predefined tokens.
const (
	// Special tokens.
	// These are not literals, operators, or keywords.
	ILLEGAL Token = iota
	EOF
	COMMENT

	literalBegin
	// Literal tokens.
	IDENT
	INT
	CHAR
	STRING
	literalEnd

	operatorBegin
	// Operators and delimiters.
	ADD		// +
	SUB		// -
	MUL		// *
	DIV		// /
	// MOD is listed as a keyword because % is reserved for binary integers.

	BAND	// &
	BOR		// |
	BXOR	// ^
	SHL		// <<
	SHR		// >>
	CMPL	// ~

	EQ		// ==
	NE		// !=
	LT		// <
	LE		// <=
	GT		// >
	GE		// >=

	LAND	// &&
	LOR		// ||
	NOT		// !

	COMMA	// ,
	SEMI		// ;
	COLON	// :

	AT		// @ (denotes local labels)
	NEXT	// :+ (reference to next nameless label in scope)
	PREV		// :- (reference to previous nameless label in scope)

	LPAREN	// (
	RPAREN	// )
	operatorEnd

	keywordBegin
	// Keywords and entire classes of keywords for 68000-specific tokens.
	OPCODE		// all opcodes
	DATAREG		// d0 .. d7
	ADDRREG		// a0 .. a7, sp
	DATAREG_W	// d0.w .. d7.w
	ADDRREG_W	// a0.w .. a7.w, sp.w
	DATAREG_L	// d0.l .. d7.l
	ADDRREG_L	// a0.l .. a7.l, sp.l
	PC			// pc
	USP			// usp
	CCR			// ccr
	SR			// sr
	DOT_W		// .w (absolute addressing suffix)
	DOT_L		// .l (absolute addressing suffix)

	DOT			// . (the current position; equivalent to $ or * in other assemblers)
	MOD			// .mod
	keywordEnd
)

// Lookup returns the token type for the identifier or keyword stored in str. If str does not store a keyword, IDENT is returned.
func Lookup(str string) Token {
	TODO
	return IDENT
}

// IsLiteral returns whether t is a literal.
func (t Token) IsLiteral() bool {
	return t > literalBegin && t < literalEnd
}

// IsOperator returns whether t is an operator.
// MOD is not considered an operator for the purposes of this test.
func (t Token) IsOperator() bool {
	return t > operatorBegin && t < operatorEnd
}

// IsKeyword returns whether t is a keyword.
func (t Token) IsKeyword() bool {
	return t > keywordBegin && t < keywordEnd
}

// Precedence returns the binary-operator precedence for t. TODO
func (t Token) Precedence() int {
	TODO
}

// String returns the string for t. TODO
func (t Token) String() string {
	TODO
}
