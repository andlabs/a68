// 22 december 2019
package core

import (
	"testing"
	"bytes"

	// TODO this is BSD 3-clause, which is technically not MIT compatible
	"github.com/google/go-cmp/cmp"
)

func mustAdd(t *testing.T, e *Expr, op ExprOpcode) {
	err := e.Add(op)
	if err != nil {
		t.Fatalf("error adding %v to expression in mk function: %v", op, err)
	}
}

func mustAddInt(t *testing.T, e *Expr, n uint64) {
	err := e.AddInt(n)
	if err != nil {
		t.Fatalf("error adding int %d to expression in mk function: %v", n, err)
	}
}

func mustAddName(t *testing.T, e *Expr, name string) {
	err := e.AddName(name)
	if err != nil {
		t.Fatalf("error adding name %q to expression in mk function: %v", name, err)
	}
}

func mustFinish(t *testing.T, e *Expr) {
	err := e.Finish()
	if err != nil {
		t.Fatalf("error finishing expression in mk function: %v", err)
	}
}

var goodExprCases = []struct {
	name	string
	raw		[]byte
	mk		func(t *testing.T) *Expr
	value	uint64
	valerrs	[]error
}{{
	name:	"5",
	raw:		[]byte{
		1,
		byte(ExprInt), 5,
	},
	mk:		func(t *testing.T) *Expr {
		e := NewExpr()
		mustAddInt(t, e, 5)
		mustFinish(t, e)
		return e
	},
	value:	5,
}, {
	name:	"KnownName",
	raw:		[]byte{
		1,
		byte(ExprName), 9, 'K', 'n', 'o', 'w', 'n', 'N', 'a', 'm', 'e',
	},
	mk:		func(t *testing.T) *Expr {
		e := NewExpr()
		mustAddName(t, e, "KnownName")
		mustFinish(t, e)
		return e
	},
	value:	5,
}, {
	name:	"UnknownName",
	raw:		[]byte{
		1,
		byte(ExprName), 11, 'U', 'n', 'k', 'n', 'o', 'w', 'n', 'N', 'a', 'm', 'e',
	},
	mk:		func(t *testing.T) *Expr {
		e := NewExpr()
		mustAddName(t, e, "UnknownName")
		mustFinish(t, e)
		return e
	},
	valerrs:	[]error{UnknownNameError("UnknownName")},
}}

type testEvalHandler struct {
	errs		[]error
}

func (h *testEvalHandler) LookupName(name string) (val uint64, ok bool) {
	if name == "KnownName" {
		return 5, true
	}
	return 0, false
}

func (h *testEvalHandler) ReportError(err error) {
	h.errs = append(h.errs, err)
}

func testRead(t *testing.T, data []byte) *Expr {
	e := NewExpr()
	n, err := e.ReadFrom(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ReadFrom() failed: %v", err)
	} else if n != int64(len(data)) {
		t.Fatalf("ReadFrom() read wrong amount: got %d, want %d", n, len(data))
	}
	return e
}

func testWrite(t *testing.T, e *Expr, want []byte) {
	b := &bytes.Buffer{}
	n, err := e.WriteTo(b)
	if err != nil {
		t.Errorf("WriteTo() failed: %v", err)
	} else if n != int64(b.Len()) {
		t.Errorf("WriteTo() indicated it wrote wrong amount: got %d, want %d", n, b.Len())
	}
	if diff := cmp.Diff(b.Bytes(), want); diff != "" {
		t.Errorf("WriteTo() wrote wrong data: (-got +want)\n%v", diff)
	}
}

func testEval(t *testing.T, e *Expr, wantval uint64, wanterrs []error) {
	wantok := len(wanterrs) == 0
	h := &testEvalHandler{}
	gotval, gotok := e.Evaluate(h)
	if gotval != wantval || gotok != wantok {
		t.Errorf("Evaluate() return value wrong: got (%v, %v), want (%v, %v)", gotval, gotok, wantval, wantok)
	}
	if diff := cmp.Diff(h.errs, wanterrs); diff != "" {
		t.Errorf("Evaluate() returned wrong errors: (-got +want)\n%v", diff)
	}
}

func TestExprMkEval(t *testing.T) {
	for _, tc := range goodExprCases {
		t.Run(tc.name, func(t *testing.T) {
			e := tc.mk(t)
			testEval(t, e, tc.value, tc.valerrs)
		})
	}
}

func TestExprReadEval(t *testing.T) {
	for _, tc := range goodExprCases {
		t.Run(tc.name, func(t *testing.T) {
			e := testRead(t, tc.raw)
			testEval(t, e, tc.value, tc.valerrs)
		})
	}
}

func TestExprMkWrite(t *testing.T) {
	for _, tc := range goodExprCases {
		t.Run(tc.name, func(t *testing.T) {
			e := tc.mk(t)
			testWrite(t, e, tc.raw)
		})
	}
}

func TestExprReadWrite(t *testing.T) {
	for _, tc := range goodExprCases {
		t.Run(tc.name, func(t *testing.T) {
			e := testRead(t, tc.raw)
			testWrite(t, e, tc.raw)
		})
	}
}

// TODO TestExprStableInversionMkWriteReadEval
// TODO TestExprStableInversionMkWriteReadWrite
// TODO TestExprStableInversionReadWriteReadEval

// TODO all the error conditions
