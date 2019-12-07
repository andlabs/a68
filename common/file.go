// 5 december 2019
package common

import (
	"go/token"
)

type Pos = token.Pos
const NoPos Pos = token.NoPos

type Position = token.Position

type File = token.File

type FileSet = token.FileSet
func NewFileSet() *FileSet {
	return token.NewFileSet()
}
