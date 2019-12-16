// 15 december 2019
package core

import (
	"encoding/binary"
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
	int		uint64
	str		string
}

func (e *exprOp) MarshalBinary() ([]byte, error) {
	bcode := byte(e.code)
	if e.code == ExprInt {
		b := make([]byte, 1 + binary.MaxVarintLen64)
		b[0] = bcode
		n := binary.PutUvarint(b[1:], e.int)
		return b[:1 + n], nil
	}
	if e.code == ExprName {
		l := len(e.str)
		b := make([]byte, 1 + binary.MaxVarintLen64 + l)
		b[0] = bcode
		n := binary.PutUvarint(b[1:], uint64(l))
		copy(b[1 + n:], e.str)
		return b[:1 + n + l], nil
	}
	return []byte{bcode}, nil
}

func binaryUvarintErr(b []byte) (x uint64, n int, err error) {
	x, n = binary.Uvarint(b)
	if n == 0 {
		err = fmt.Errorf("incomplete uvarint")
	} else if n < 0 {
		err = fmt.Errorf("uvarint overflows uint64; %d bytes read", -n)
	}
	return
}

func (e *exprOp) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("exprOp.UnmarshalBinary: no data")
	}
	code := exprOpcode(data[0]
	if code >= nExprOpcodes {
		return fmt.Errorf("exprOp.UnmarshalBinary: bad opcode 0x%X", code)
	}
	data = data[1:]
	if code == ExprInt {
		if len(data) == 0 {
			return fmt.Errorf("exprOp.UnmarshalBinary: ExprInt without integer")
		}
		x, n, err := binaryUvarintError(data)
		if err != nil {
			return fmt.Errorf("exprOp.UnmarshalBinary: ExprInt error: %v", err)
		}
		if len(data) != n {
			return fmt.Erorrf("exprOp.UnmarshalBinary: ExprInt value length mismatch: got %d, want %d", len(data), n)
		}
		e.code = code
		e.int = x
		return nil
	}
	if code == ExprName {
		if len(data) == 0 {
			return fmt.Errorf("exprOp.UnmarshalBinary: ExprName without length or name")
		}
		x, n, err := binaryUvarintError(data)
		if err != nil {
			return fmt.Errorf("exprOp.UnmarshalBinary: ExprName length error: %v", err)
		}
		data = data[n:]
		if len(data) != x {
			return fmt.Errorf("exprOp.UnmarshalBinary: ExprName name length mismatch: got %d, want %d", len(data), x)
		}
		e.code = code
		e.str = string(data)
		return nil
	}
	if len(data) != 0 {
		return fmt.Errorf("exprOp.UnmarshalBinary: unexpected %d extra byte(s) at end", len(data))
	}
	e.code = code
	return nil
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
	ops	[]exprOp
}

func NewExpr() *Expr {
	return &Expr{
		ops:		make([]exprOp, 0, 16),
	}
}

func (e *Expr) finished() bool {
	return len(e.ops) > 0 && e.ops[len(e.ops) - 1].code == ExprStop
}

func (e *Expr) Add(code ExprOpcode) error {
	if e.finished() {
		return fmt.Errorf("cannot add to finished expression")
	}
	if code == ExprInt || code == ExprName {
		return fmt.Errorf("cannot add %v without argument", code)
	}
	e.ops = append(e.ops, exprOp{
		code:		code,
	})
	return nil
}

func (e *Expr) AddInt(n uint64) error {
	if e.finished() {
		return fmt.Errorf("cannot add to finished expression")
	}
	e.ops = append(e.ops, exprOp{
		code:		ExprInt,
		int:			n,
	})
}

func (e *Expr) AddName(name string) error {
	if e.finished() {
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
