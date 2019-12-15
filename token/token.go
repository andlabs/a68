// 4 december 2019
package token

import (
	gotoken "go/token"
	"strconv"

	"github.com/andlabs/a68/cpu"
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

	POUND	// #

	LPAREN	// (
	RPAREN	// )
	operatorEnd

	keywordBegin
	keywordClassBegin
	// Keywords and entire classes of keywords for 68000-specific tokens.
	OPCODE		// all opcodes
	DATAREG		// d0 .. d7
	ADDRREG		// a0 .. a7, sp
	DATAREG_W	// d0.w .. d7.w
	ADDRREG_W	// a0.w .. a7.w, sp.w
	DATAREG_L	// d0.l .. d7.l
	ADDRREG_L	// a0.l .. a7.l, sp.l
	keywordClassEnd

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

var tokens = [...]string{
	ILLEGAL:		"ILLEGAL",
	EOF:			"EOF",
	COMMENT:	"COMMENT",

	IDENT:		"IDENT",
	INT:			"INT",
	CHAR:		"CHAR",
	STRING:		"STRING",

	ADD:		"+",
	SUB:			"-",
	MUL:			"*",
	DIV:			"/",

	BAND:		"&",
	BOR:			"|",
	BXOR:		"^",
	SHL:			"<<",
	SHR:			">>",
	CMPL:		"~",

	EQ:			"==",
	NE:			"!=",
	LT:			"<",
	LE:			"<=",
	GT:			">",
	GE:			">=",

	LAND:		"&&",
	LOR:			"||",
	NOT:			"!",

	COMMA:		",",
	SEMI:		";",
	COLON:		":",

	AT:			"@",
	NEXT:		":+",
	PREV:		":-",

	POUND:		"#",

	LPAREN:		"(",
	RPAREN:		")",

	OPCODE:		"OPCODE",
	DATAREG:		"DATAREG",
	ADDRREG:	"ADDRREG",
	DATAREG_W:	"DATAREG_W",
	ADDRREG_W:	"ADDRREG_W",
	DATAREG_L:	"DATAREG_L",
	ADDRREG_L:	"ADDRREG_L",

	PC:			"pc",
	USP:			"usp",
	CCR:			"ccr",
	SR:			"sr",
	DOT_W:		".w",
	DOT_L:		".l",

	DOT:			".",
	MOD:		".mod",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token, (len(cpu.Opcodes) * 5) + (8 * 6) + 3 + (keywordEnd - keywordClassEnd))
	for _, op := range cpu.Opcodes {
		n := op.Name()
		keywords[n] = OPCODE
		keywords[n + ".b"] = OPCODE
		keywords[n + ".w"] = OPCODE
		keywords[n + ".l"] = OPCODE
		keywords[n + ".s"] = OPCODE
	}
	for i := 0; i <= 7; i++ {
		n := strconv.Itoa(i)
		dn := "d" + n
		an := "a" + n
		keywords[dn] = DATAREG
		keywords[an] = ADDRREG
		keywords[dn + ".w"] = DATAREG_W
		keywords[an + ".w"] = ADDRREG_W
		keywords[dn + ".l"] = DATAREG_L
		keywords[an + ".l"] = ADDRREG_L
	}
	keywords["sp"] = ADDRREG
	keywords["sp.w"] = ADDRREG_W
	keywords["sp.l"] = ADDRREG_L
	for t := keywordClassEnd + 1; t < keywordEnd; t++ {
		keywords[tokens[t]] = t
	}
}

// Lookup returns the token type for the identifier or keyword stored in str. If str does not store a keyword, IDENT is returned.
func Lookup(str string) Token {
	t, ok := keywords[str]
	if ok {
		return t
	}
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
	return t > keywordBegin && t < keywordEnd && t != keywordClassBegin && t != keywordClassEnd
}

const (
	LowestPrec = gotoken.LowestPrec
	UnaryPrec = gotoken.UnaryPrec
	HightestPrec = gotoken.HighestPrec
)

// Precedence returns the binary-operator precedence for t.
// The precedence of MOD is the same as that of MUL.
// Precedence rules are the same as in Go.
// If t is neither MOD nor a binary operator, LowestPrec is returned.
func (t Token) Precedence() int {
	switch t {
	case LOR:
		return 1
	case LAND:
		return 2
	case EQ, NE, LT, LE, GT, GE:
		return 3
	case ADD, SUB, BOR, BXOR:
		return 4
	case MUL, DIV, MOD, SHL, SHR, BAND:
		return 5
	}
	return LowestPrec
}

// String returns the string for t.
func (t Token) String() string {
	if t >= 0 && t < len(tokens) {
		return tokens[t]
	}
	return "Token(" + strconv.Itoa(int(t)) + ")"
}
