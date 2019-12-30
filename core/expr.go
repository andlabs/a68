// 15 december 2019
package core

import (
	"fmt"
	"io"
	"encoding/binary"
	"bytes"
)

// TODO properly mark internal errors?

type ExprOpcode byte
const (
	ExprInt ExprOpcode = iota
	ExprName
	ExprNeg
	ExprNot		// current implementation: 0 is false != 0 is false and !0 == 1
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
	ExprEq		// current implementation: does a signed 64-bit integer comparison
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
	ExprInt:		"ExprInt",
	ExprName:	"ExprName",
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

var exprOpcodeStackDeltas = [nExprOpcodes]int{
	ExprInt:		1,
	ExprName:	1,
	ExprNeg:		0,
	ExprNot:		0,
	ExprCmpl:	0,
	ExprMul:		-1,
	ExprDiv:		-1,
	ExprMod:		-1,
	ExprShl:		-1,
	ExprShr:		-1,
	ExprBAnd:	-1,
	ExprAdd:		-1,
	ExprSub:		-1,
	ExprBOr:		-1,
	ExprBXor:		-1,
	ExprEq:		-1,
	ExprNe:		-1,
	ExprLt:		-1,
	ExprLe:		-1,
	ExprGt:		-1,
	ExprGe:		-1,
	ExprLAnd:	-1,
	ExprLOr:		-1,
}

func (e ExprOpcode) String() string {
	if e >= nExprOpcodes {
		return fmt.Sprintf("ExprOpcode(0x%X)", byte(e))
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

// type to juggle a few varying requirements
// - ReaderFrom requires us to return the number of bytes read, which binary.ReadUvariant() does not
// - binary.ReadUvarint() requires an io.ByteReader
type trackingReader struct {
	r	io.Reader
	n	int
}
func (r *trackingReader) readFull(p []byte) (int, error) {
	n, err := io.ReadFull(r.r, p)
	r.n += n
	return n, err
}
func (r *trackingReader) ReadByte() (byte, error) {
	b := make([]byte, 1)
	_, err := r.readFull(b)
	return b[0], err
}

func (e *exprOp) ReadFrom(r io.Reader) (n int64, err error) {
	return e.readFrom(&trackingReader{r: r})
}

func (e *exprOp) readFrom(r *trackingReader) (n int64, err error) {
	e2, err := readExprOp(r)
	if err == nil {
		*e = e2
	}
	return int64(r.n), err
}

func readExprOp(r *trackingReader) (e exprOp, err error) {
	b, err := r.ReadByte()
	if err != nil {
		return e, readError(err)
	}
	e.code = ExprOpcode(b)
	if e.code >= nExprOpcodes {
		return e, fmt.Errorf("bad opcode 0x%X", e.code)
	}
	if e.code == ExprInt || e.code == ExprName {
		e.int, err = binary.ReadUvarint(r)
		if err != nil {
			return e, readError(err)
		}
	}
	if e.code == ExprName {
		buf := make([]byte, e.int)
		_, err = r.readFull(buf)
		if err != nil {
			return e, readError(err)
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
	b := make([]byte, 1 + len(num) + len(str))
	b[0] = byte(e.code)
	copy(b[1:1 + len(num)], num)
	copy(b[1 + len(num):], str)
	return bytes.NewReader(b).WriteTo(w)
}

func (e *exprOp) String() string {
	if e.code >= nExprOpcodes {
		return fmt.Sprintf("#INVALID 0x%X", e.code)
	}
	if e.code == ExprInt {
		return fmt.Sprintf("%v 0x%08X", e.code, e.int)
	}
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
	return nil
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
		int:			uint64(len(name)),
		str:			name,
	})
	return nil
}

func (e *Expr) checkValid() error {
	nStack := 0
	for i, op := range e.ops {
		delta := exprOpcodeStackDeltas[op.code]
		if (nStack + (delta - 1)) < 0 {
			// don't allow popping an empty stack
			// TODO don't assume that (delta - 1) is the number of stack pops (even though it is true for now)
			return fmt.Errorf("%v at index %d does not have enough arguments", op.code, i)
		}
		nStack += delta
	}
	if nStack != 1 {
		return fmt.Errorf("expression doesn't resolve to a single value")
	}
	return nil
}

func (e *Expr) Finish() error {
	if e.finished {
		return fmt.Errorf("cannot finish finished expression again")
	}
	if len(e.ops) == 0 {
		return fmt.Errorf("cannot finish empty expression")
	}
	err := e.checkValid()
	if err != nil {
		return fmt.Errorf("cannot finish invalid expression: %v", err)
	}
	e.finished = true
	return nil
}

func (e *Expr) Empty() bool {
	return len(e.ops) == 0
}

func (e *Expr) ReadFrom(r io.Reader) (n int64, err error) {
	return e.readFrom(&trackingReader{r: r})
}

func (e *Expr) readFrom(r *trackingReader) (n int64, err error) {
	e2, err := readExpr(r)
	if err == nil {
		*e = *e2
	}
	return int64(r.n), err
}

func readExpr(r *trackingReader) (e *Expr, err error) {
	e = NewExpr()
	n, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, readError(err)
	}
	if n == 0 {
		// Since we cannot have zero-length expressions yet, we can use whether the first byte (and thus, n) is zero to denote any future revisions to the binary expression representation.
		return nil, fmt.Errorf("unsupported expression encoding version")
	}
	e.ops = make([]exprOp, n)
	for i, _ := range e.ops {
		_, err = e.ops[i].readFrom(r)
		if err != nil {
			return nil, readError(err)
		}
	}
	err = e.Finish()
	if err != nil {
		return nil, fmt.Errorf("invalid expression read: %v", err)
	}
	return e, nil
}

func (e *Expr) WriteTo(w io.Writer) (n int64, err error) {
	if !e.finished {
		return 0, fmt.Errorf("cannot write unfinished expression")
	}

	// first write number of opcodes
	num := make([]byte, binary.MaxVarintLen64)
	nn := binary.PutUvarint(num, uint64(len(e.ops)))
	num = num[:nn]

	// the Write and WriteTo calls here cannot fail according to the documentation for bytes.Buffer
	buf := new(bytes.Buffer)
	buf.Write(num)
	for _, op := range e.ops {
		op.WriteTo(buf)
	}

	return buf.WriteTo(w)
}

type EvaluateHandler interface {
	LookupName(name string) (val uint64, ok bool)
	ReportError(err error)
}

var (
	ErrEvaluatingUnfinishedExpr = fmt.Errorf("cannot evaluate unfinished expression")
	ErrZeroDivisor = fmt.Errorf("division by zero")
	ErrZeroDivisorMod = fmt.Errorf("division by zero in modulo")
)

type UnknownNameError string

func (e UnknownNameError) Error() string {
	return fmt.Sprintf("unknown names %q", string(e))
}

func boolval(b bool) uint64 {
	if b {
		return 1
	}
	return 0
}

func (e *Expr) Evaluate(handler EvaluateHandler) (val uint64, ok bool) {
	if !e.finished {
		// this also enforces the precondition that the stack will always have the right number of entries
		handler.ReportError(ErrEvaluatingUnfinishedExpr)
		return 0, false
	}
	stack := make([]uint64, 0, 16)
	pop := func() (v uint64) {
		i := len(stack) - 1
		v = stack[i]
		stack = stack[:i]
		return v
	}
	pop2 := func() (a uint64, b uint64) {
		i := len(stack) - 2
		a = stack[i]
		b = stack[i + 1]
		stack = stack[:i]
		return a, b
	}
	noError := true
	for _, op := range e.ops {
		switch op.code {
		case ExprInt:
			stack = append(stack, op.int)
		case ExprName:
			val, ok := handler.LookupName(op.str)
			if !ok {
				handler.ReportError(UnknownNameError(op.str))
				noError = false
				val = 1		// don't stop evaluation
			}
			stack = append(stack, val)
		case ExprNeg:
			val := pop()
			val = ^val + 1
			stack = append(stack, val)
		case ExprNot:
			val := pop()
			if val != 0 {
				val = 0
			} else {
				val = 1
			}
			stack = append(stack, val)
		case ExprCmpl:
			val := pop()
			val = ^val
			stack = append(stack, val)
		case ExprMul:
			a, b := pop2()
			stack = append(stack, a * b)
		case ExprDiv:
			a, b := pop2()
			if b == 0 {
				handler.ReportError(ErrZeroDivisor)
				noError = false
				b = 1			// don't stop evaluation
			}
			stack = append(stack, a / b)
		case ExprMod:
			a, b := pop2()
			if b == 0 {
				handler.ReportError(ErrZeroDivisorMod)
				noError = false
				b = 1			// don't stop evaluation
			}
			stack = append(stack, a % b)
		case ExprShl:
			a, b := pop2()
			stack = append(stack, a << b)
		case ExprShr:
			a, b := pop2()
			stack = append(stack, a >> b)
		case ExprBAnd:
			a, b := pop2()
			stack = append(stack, a & b)
		case ExprAdd:
			a, b := pop2()
			stack = append(stack, a + b)
		case ExprSub:
			a, b := pop2()
			stack = append(stack, a - b)
		case ExprBOr:
			a, b := pop2()
			stack = append(stack, a | b)
		case ExprBXor:
			a, b := pop2()
			stack = append(stack, a ^ b)
		case ExprEq:
			a, b := pop2()
			stack = append(stack, boolval(int64(a) == int64(b)))
		case ExprNe:
			a, b := pop2()
			stack = append(stack, boolval(int64(a) != int64(b)))
		case ExprLt:
			a, b := pop2()
			stack = append(stack, boolval(int64(a) < int64(b)))
		case ExprLe:
			a, b := pop2()
			stack = append(stack, boolval(int64(a) <= int64(b)))
		case ExprGt:
			a, b := pop2()
			stack = append(stack, boolval(int64(a) > int64(b)))
		case ExprGe:
			a, b := pop2()
			stack = append(stack, boolval(int64(a) >= int64(b)))
		case ExprLAnd:
			a, b := pop2()
			val := uint64(0)
			if a != 0 && b != 0 {
				val = 1
			}
			stack = append(stack, val)
		case ExprLOr:
			a, b := pop2()
			val := uint64(1)
			if a == 0 && b == 0 {
				val = 0
			}
			stack = append(stack, val)
		default:
			panic("can't happen; likely missing new opcode implementation in Evaluate()")
		}
	}
	if !noError {
		return 0, false
	}
	return stack[0], true
}
