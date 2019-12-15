// 11 december 2019
package token

import (
	"github.com/andlabs/a68/common"
)

var (
	// the current position; equivalent to $ or * in other assemblers
	DOT = common.AddKeyword(".")

	// provided instead of % because % is used for binary numbers
	MOD = common.AddKeyword(".mod")
)
