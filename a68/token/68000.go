// 11 december 2019
package token

import (
	"github.com/andlabs/a68/common"
)

var (
	// Normal registers.
	D0 = common.AddKeyword("d0")
	D1 = common.AddKeyword("d1")
	D2 = common.AddKeyword("d2")
	D3 = common.AddKeyword("d3")
	D4 = common.AddKeyword("d4")
	D5 = common.AddKeyword("d5")
	D6 = common.AddKeyword("d6")
	D7 = common.AddKeyword("d7")
	A0 = common.AddKeyword("a0")
	A1 = common.AddKeyword("a1")
	A2 = common.AddKeyword("a2")
	A3 = common.AddKeyword("a3")
	A4 = common.AddKeyword("a4")
	A5 = common.AddKeyword("a5")
	A6 = common.AddKeyword("a6")
	A7 = common.AddKeyword("a7")
	SP = common.AddKeyword("sp")

	// Word indexes.
	D0_W = common.AddKeyword("d0.w")
	D1_W = common.AddKeyword("d1.w")
	D2_W = common.AddKeyword("d2.w")
	D3_W = common.AddKeyword("d3.w")
	D4_W = common.AddKeyword("d4.w")
	D5_W = common.AddKeyword("d5.w")
	D6_W = common.AddKeyword("d6.w")
	D7_W = common.AddKeyword("d7.w")
	A0_W = common.AddKeyword("a0.w")
	A1_W = common.AddKeyword("a1.w")
	A2_W = common.AddKeyword("a2.w")
	A3_W = common.AddKeyword("a3.w")
	A4_W = common.AddKeyword("a4.w")
	A5_W = common.AddKeyword("a5.w")
	A6_W = common.AddKeyword("a6.w")
	A7_W = common.AddKeyword("a7.w")
	SP_W = common.AddKeyword("sp.w")

	// Long indexes.
	D0_L = common.AddKeyword("d0.l")
	D1_L = common.AddKeyword("d1.l")
	D2_L = common.AddKeyword("d2.l")
	D3_L = common.AddKeyword("d3.l")
	D4_L = common.AddKeyword("d4.l")
	D5_L = common.AddKeyword("d5.l")
	D6_L = common.AddKeyword("d6.l")
	D7_L = common.AddKeyword("d7.l")
	A0_L = common.AddKeyword("a0.l")
	A1_L = common.AddKeyword("a1.l")
	A2_L = common.AddKeyword("a2.l")
	A3_L = common.AddKeyword("a3.l")
	A4_L = common.AddKeyword("a4.l")
	A5_L = common.AddKeyword("a5.l")
	A6_L = common.AddKeyword("a6.l")
	A7_L = common.AddKeyword("a7.l")
	SP_L = common.AddKeyword("sp.l")

	// Special registers.
	PC = common.AddKeyword("pc")
	USP = common.AddKeyword("usp")
	CCR = common.AddKeyword("ccr")
	SR = common.AddKeyword("sr")

	// Memory direct suffixes.
	DOT_W = common.AddKeyword(".w")
	DOT_L = common.AddKeyword(".l")
)
