// 30 june 2012
package main

type ExprOpcode int

const (
	eNumber		ExprOpcode = iota
	eIdentifier
	eLocal
	eCharacter
	eCall
	eAdd
	eSubtract
	eMultiply
	eDivide
	eRemainder
	eBitAnd
	eBitOr
	eBitXor
	eLeftShift
	eRightShift
	eNegate
	eComplement
	eEqual
	eNotEqual
	eLess
	eLessEqual
	eGreater
	eGreaterEqual
	eAnd
	eOr
	eNot
	// ...
)

type exprop struct {
	opcode	ExprOpcode
	arg		string
}

type Expression []exprop
const EBUFSIZ = 500			// TODO meter this for optimal memory use?

func NewExpr(o ExprOpcode, v string) Expression {
	e := make(Expression, 0, EBUFSIZ)
	e.Add(o, v)
	return e
}

func (e *Expression) Add(o ExprOpcode, v string) {
	if len(*e) + 1 >= cap(*e) {
		// reallocate more just in case we need to add more later
		// grow linearly to keep memory use at a minimum
		s2 := make([]exprop, len(*e), cap(*e) + EBUFSIZ)
		copy(s2, *e)
		*e = s2
	}
	*e = append(*e, exprop{
		opcode:	o,
		arg:		v,
	})
}

func (e *Expression) Concatenate(e2 *Expression) {
	for i := 0; i < len(*e2); i++ {
		e.Add((*e2)[i].opcode, (*e2)[i].arg)
	}
}

func (e Expression) CanEvaluateNow() bool {
	// TODO
	return true
}

type eval_stack []uint32

func newEvalStack() eval_stack {
	return make(eval_stack, 0, EBUFSIZ)
}

func (e *eval_stack) push(v uint32) {
	if len(*e) + 1 >= cap(*e) {
		// reallocate like above
		s2 := make(eval_stack, len(*e), cap(*e) + EBUFSIZ)
		copy(s2, *e)
		*e = s2
	}
	*e = append(*e, v)
}

func (e *eval_stack) pop() uint32 {
	if len(*e) == 0 {
		FATAL_BUG("pop from empty evaluation stack")
	}
	pos := len(*e) - 1
	v := (*e)[pos]
	*e = (*e)[:pos - 1]
	return v
}

func (e Expression) Evaluate() uint32 {
	s := newEvalStack()

	for pc := 0; pc < len(e); pc++ {
		switch e[pc].opcode {
		case eNumber:
//(TODO)			s.push(getNumber(e[pc].arg))
		case eIdentifier:
			// TODO
		case eLocal:
			// TODO
		case eCharacter:
			//  TODO
		case eCall:
			// TODO
		case eAdd:
			// TODO
		case eSubtract:
			// TODO
		case eMultiply:
			// TODO
		case eDivide:
			// TODO
		case eRemainder:
			// TODO
		case eBitAnd:
			// TODO
		case eBitOr:
			// TODO
		case eBitXor:
			// TODO
		case eLeftShift:
			// TODO
		case eRightShift:
			// TODO
		case eNegate:
			// TODO
		case eComplement:
			// TODO
		case eEqual:
			// TODO
		case eNotEqual:
			// TODO
		case eLess:
			// TODO
		case eLessEqual:
			// TODO
		case eGreater:
			// TODO
		case eGreaterEqual:
			// TODO
		case eAnd:
			// TODO
		case eOr:
			// TODO
		case eNot:
			// TODO
		default:
			FATAL_BUG(
				"undefined opcode %d evaluating expression!", e[pc].opcode)
			// TODO print more information about this
		}
	}

	if len(s) > 1 {
		POTENTIAL_BUG(
			"evaluating expression resulted in leftover values on stack")
		// TODO print more information about this?
	} else if len(s) == 0 {
		FATAL_BUG(
			"evaluating expression left no result somehow!")
		// TODO definitely print more information about this
	}
	return s.pop()
}

// for debugging
func (e Expression) String() string {
	// TODO
	return "<stub Expression.String()>"
}
