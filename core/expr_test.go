// 22 december 2019
package core

import (
	"testing"
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
	raw:		[]byte{byte(ExprInt), 5},
	mk:		func(t *testing.T) *Expr {
		e := NewExpr()
		mustAddInt(t, e, 5)
		mustFinish(t, e)
		return e
	},
	value:	5,
}}

func TestExprs(t *testing.T) {
	for _, tc := range goodExprCases {
		t.Run(tc.name, func(t *testing.T) {
			e := tc.mk(t)
			testEval(t, e, tc.value, tc.valerrs)
		})
	}
}
