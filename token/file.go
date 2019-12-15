// 5 december 2019
package token

import (
	gotoken "go/token"
)

type Pos = gotoken.Pos
const NoPos Pos = gotoken.NoPos

type Position = gotoken.Position

type File = gotoken.File

type FileSet = gotoken.FileSet
func NewFileSet() *FileSet {
	return gotoken.NewFileSet()
}
