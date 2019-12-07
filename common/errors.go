// 5 december 2019
package common

import (
	"io"
	"go/scanner"
)

func PrintError(w io.Writer, err error) {
	scanner.PrintError(w, err)
}

type Error = scanner.Error

type ErrorHandler = scanner.ErrorHandler

type ErrorList = scanner.ErrorList
