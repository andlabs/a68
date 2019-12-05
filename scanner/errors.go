// 5 december 2019
package scanner

import (
	"io"
	goscanner "go/scanner"
)

func PrintError(w io.Writer, err error) {
	goscanner.PrintError(w, err)
}

type Error = goscanner.Error

type ErrorHandler = goscanner.ErrorHandler

type ErrorList = goscanner.ErrorList
