// 15 december 2019
package core

import (
	"io"
	"encoding/binary"
	"bytes"
)

// TODO properly mark internal errors?

type ExprOpcode byte
const (
	ExprStop ExprOpcode = iota
	ExprInt
	ExprName
	ExprUnaryPlus
	ExprNeg
	ExprNot
	ExprCmpl
	ExprMul
	ExprDiv
	ExprMod
	ExprShl
	ExprShr
	ExprBAnd
	ExprAdd
	ExprSub
	ExprBOr
	ExprBXor
	ExprEq
	ExprNe
	ExprLt
	ExprLe
	ExprGt
	ExprGe
	ExprLAnd
	ExprLOr
	nExprOpcodes
)

var exprOpcodeStrings = [nExprOpcodes]string{
	ExprStop:		"ExprStop",
	ExprInt:		"ExprInt",
	ExprName:	"ExprName",
	ExprUnaryPlus:		"ExprUnaryPlus",
	ExprNeg:		"ExprNeg",
	ExprNot:		"ExprNot",
	ExprCmpl:	"ExprCmpl",
	ExprMul:		"ExprMul",
	ExprDiv:		"ExprDiv",
	ExprMod:		"ExprMod",
	ExprShl:		"ExprShl",
	ExprShr:		"ExprShr",
	ExprBAnd:	"ExprBAnd",
	ExprAdd:		"ExprAdd",
	ExprSub:		"ExprSub",
	ExprBOr:		"ExprBOr",
	ExprBXor:		"ExprBXor",
	ExprEq:		"ExprEq",
	ExprNe:		"ExprNe",
	ExprLt:		"ExprLt",
	ExprLe:		"ExprLe",
	ExprGt:		"ExprGt",
	ExprGe:		"ExprGe",
	ExprLAnd:	"ExprLAnd",
	ExprLOr:		"ExprLOr",
}

func (e ExprOpcode) String() string {
	if e >= nExprOpcodes {
		return fmt.Sprintf("ExprOpcode(0x%X)", e)
	}
	return exprOpcodeStrings[e]
}

type exprOp struct {
	code		ExprOpcode
	int		uint64		// (ExprInt, ExprName - see below comment)
	str		string		// (ExprName) length is stored in int to simplify below code
}

// TODO allow overriding what's returned on EOF
func readError(err error) error {
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return err
}

// binary.ReadUvarint() requires an io.ByteReader, grrr
type fakeByteReader struct {
	r	io.Reader
}
func (r *fakeByteReader) ReadByte() (byte, error) {
	b := make([]byte, 0)
	_, err := io.ReadFull(r.r, b)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}
func toByteReader(r io.Reader) io.ByteReader {
	if br, ok := r.(io.ByteReader); ok {
		return br
	}
	return &fakeByteReader{r}
}

// TODO convert this into a ReadFrom
func readExprOp(r io.Reader) (e exprOp, err error) {
	br := toByteReader(r)
	b, err := br.ReadByte()
	if err != nil {
		return exprOp{}, readError(err)
	}
	e.code = exprOpcode(b)
	if e.code >= nExprOpcodes {
		return exprOp{}, fmt.Errorf("bad opcode 0x%X", code)
	}
	if e.code == ExprInt || e.code == ExprName {
		e.int, err = binary.ReadUvarint(br)
		if err != nil {
			return exprOp{}, readError(err)
		}
	}
	if e.code == ExprName {
		buf := make([]byte, e.int)
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return exprOp{}, readError(err)
		}
		e.str = string(buf)
	}
	return e, nil
}

func (e *exprOp) WriteTo(w io.Writer) (n int64, err error) {
	var num []byte
	if e.code == ExprInt || e.code == ExprName {
		num = make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(num, e.int)
		num = num[:n]
	}
	var str string
	if e.code == ExprName {
		str = e.str
	}
	b = make([]byte, 1 + len(num) + len(str))
	b[0] = byte(e.code)
	copy(b[1:1 + len(num)], num)
	copy(b[1 + len(num):], str)
	n, err = w.Write(b)
	if n > len(b) {
		panic("exprOp.WriteTo(): invalid Write count")
	}
	if n != len(b) && err == nil {
		err = io.ErrShortWrite
	}
	return n, err
}

func (e *exprOp) String() string {
	if e.code >= nExprOpcodes {
		return fmt.Sprintf("#INVALID 0x%X", e.code)
	}
	if e.code == ExprInt {
		return fmt.Sprintf("%v 0x%08X", e.code, e.int)
	]
	if e.code == ExprName {
		return fmt.Sprintf("%v %q", e.code, e.str)
	}
	return e.code.String()
}

// Expr represents an expression that can be stored in an object file and evaluated by the assembler and linker.
// Expr implements encoding.BinaryMarshaler and encoding.BinaryUnmarshaler for dictating the format that they appear in object files.
type Expr struct {
	ops		[]exprOp
	finished	bool
}

func NewExpr() *Expr {
	return &Expr{
		ops:		make([]exprOp, 0, 16),
	}
}

func (e *Expr) Add(code ExprOpcode) error {
	if e.finished {
		return fmt.Errorf("cannot add to finished expression")
	}
	if code == ExprInt || code == ExprName {
		return fmt.Errorf("cannot add %v using Expr.Add()", code)
	}
	e.ops = append(e.ops, exprOp{
		code:		code,
	})
	return nil
}

func (e *Expr) AddInt(n uint64) error {
	if e.finished {
		return fmt.Errorf("cannot add to finished expression")
	}
	e.ops = append(e.ops, exprOp{
		code:		ExprInt,
		int:			n,
	})
}

func (e *Expr) AddName(name string) error {
	if e.finished {
		return fmt.Errorf("cannot add to finished expression")
	}
	if name == "" {
		return fmt.Errorf("cannot add empty name to expression")
	}
	e.ops = append(e.ops, exprOp{
		code:		ExprName,
		int:			len(name),
		str:			name,
	})
}

func (e *Expr) Finish() error {
	if e.finished {
		return fmt.Errorf("cannot finish finished expression again")
	}
	if len(e.ops) == 0 {
		return fmt.Errorf("cannot finish empty expression")
	}
	e.finished = true
	return nil
}

func (e *Expr) Empty() bool {
	return len(e.ops) == 0
}

// TODO read

func (e *Expr) WriteTo(w io.Writer) (n int64, err error) {
	if !e.finished {
		return 0, fmt.Errorf("cannot write unfinished expression")
	}

	// first write number of opcodes
	num = make([]byte, binary.MaxVarintLen64)
	nn := binary.PutUvarint(num, uint64(len(e.ops)))
	num = num[:nn]

	// the Write and WriteTo calls here cannot fail according to the documentation for bytes.Buffer
	buf := new(bytes.Buffer)
	buf.Write(num)
	for _, op := range ops {
		op.WriteTo(buf)
	}

	return buf.WriteTo(w)
}

// TODO evaluate
