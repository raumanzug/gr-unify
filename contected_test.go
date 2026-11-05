package gr

import (
	"slices"
	"testing"

	"github.com/raumanzug/gr-unify/simple"
	"github.com/raumanzug/gr-unify/unif"
)

type variableContext_t int

type contectedEquator_t map[variableContext_t]map[string]iContectedTerm

func newContectedEquator() contectedEquator_t {
	return make(map[variableContext_t]map[string]iContectedTerm)
}

type contectedVariable_t struct {
	ctxt     variableContext_t
	variable variable_t
}

type contectedOperationTerm_t struct {
	ctxt      variableContext_t
	operation string
	args      []iTerm
}

type iContectedTerm interface {
	unif.ITermWithEquality[
		contectedVariable_t,
		iContectedTerm,
		contectedEquator_t,
	]

	context() variableContext_t
}

func contect(term iTerm, ctxt variableContext_t) iContectedTerm {
	switch box := term.(type) {
	case variableTerm_t:
		return contectedVariable_t{
			ctxt:     ctxt,
			variable: box.Variable(),
		}
	case operationTerm_t:
		return contectedOperationTerm_t{
			ctxt:      ctxt,
			operation: box.operation,
			args:      box.args,
		}
	}

	return nil
}

func (recv contectedEquator_t) Get(
	v contectedVariable_t,
) (iContectedTerm, bool) {
	c1, isOK := recv[v.context()]
	if !isOK {
		return nil, false
	}

	value, hasFound := c1[string(v.variable)]

	return value, hasFound
}

func (recv contectedEquator_t) Set(
	v contectedVariable_t,
	x iContectedTerm,
) bool {
	if v.context() == x.context() {
		return false
	}

	c1, isOK := recv[v.ctxt]
	if !isOK {
		c1 = make(map[string]iContectedTerm)
		recv[v.ctxt] = c1
	}
	c1[string(v.variable)] = x

	return true
}

func (recv contectedVariable_t) context() variableContext_t {
	return recv.ctxt
}

func (recv contectedOperationTerm_t) context() variableContext_t {
	return recv.ctxt
}

func (recv contectedVariable_t) Equals(v contectedVariable_t) bool {
	return (recv.variable == v.variable) &&
		(recv.ctxt == v.ctxt)
}

func (recv contectedVariable_t) EqualsTerm(
	equ contectedEquator_t,
	term iContectedTerm,
) bool {
	return simple.EqualsTermVar[
		contectedVariable_t,
		contectedVariable_t,
		iContectedTerm,
		contectedEquator_t,
	](recv, equ, term)
}

func (recv contectedVariable_t) Normalform(
	equ contectedEquator_t,
) iContectedTerm {
	return simple.Normalform(recv, equ)
}

func (recv contectedVariable_t) Term() iContectedTerm {
	return recv
}

func (recv contectedVariable_t) Unify(
	equ contectedEquator_t,
	term iContectedTerm,
) bool {
	return simple.UnifyVar(recv, equ, term)
}

func (recv contectedVariable_t) Variable() contectedVariable_t {
	return recv
}

func (recv contectedOperationTerm_t) EqualsTerm(
	equ contectedEquator_t,
	term iContectedTerm,
) bool {
	var normalizedTerm iContectedTerm
	normalizedTerm = term.Normalform(equ)
	switch box := normalizedTerm.(type) {
	case contectedVariable_t:
		return false
	case contectedOperationTerm_t:
		if recv.operation != box.operation {
			return false
		}

		return slices.EqualFunc(
			recv.args,
			box.args,
			func(left iTerm, right iTerm) bool {
				contectedLeft := contect(left, recv.ctxt)
				contectedRight := contect(right, box.ctxt)
				return contectedLeft.EqualsTerm(
					equ,
					contectedRight,
				)
			},
		)
	}

	return false
}

func (recv contectedOperationTerm_t) Normalform(
	equ contectedEquator_t,
) iContectedTerm {
	return recv
}

func (recv contectedOperationTerm_t) Term() iContectedTerm {
	return recv
}

func (recv contectedOperationTerm_t) Unify(
	equ contectedEquator_t,
	term iContectedTerm,
) bool {
	return simple.UnifyNonVar[
		contectedVariable_t,
		contectedVariable_t,
		contectedOperationTerm_t,
		iContectedTerm,
		contectedEquator_t,
	](recv, equ, term)
}

func (recv contectedOperationTerm_t) UnifyNonVar(
	equ contectedEquator_t,
	term contectedOperationTerm_t,
) bool {
	if recv.operation != term.operation {
		return false
	}

	return slices.EqualFunc(
		recv.args,
		term.args,
		func(left iTerm, right iTerm) bool {
			contectedLeft := contect(left, recv.ctxt)
			contectedRight := contect(right, term.ctxt)
			return contectedLeft.Unify(equ, contectedRight)
		},
	)
}

func Test_ContectedUnifyA(t *testing.T) {
	myEquator := newContectedEquator()
	termA := contectedOperationTerm_t{
		ctxt:      0,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			operationTerm_t{
				operation: "h",
				args: []iTerm{
					variableTerm_t{
						symbol: "v",
					},
					variableTerm_t{
						symbol: "w",
					},
				},
			},
		},
	}
	termB := contectedOperationTerm_t{
		ctxt:      0,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}

	if termA.Unify(myEquator, termB) {
		t.Fail()
	}
}

func Test_ContectedUnifyB(t *testing.T) {
	myEquator := newContectedEquator()
	termA := contectedOperationTerm_t{
		ctxt:      0,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "y",
			},
			variableTerm_t{
				symbol: "x",
			},
		},
	}
	termB := contectedOperationTerm_t{
		ctxt:      1,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}

	if termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if !termA.Unify(myEquator, termB) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}
}

func Test_ContectedUnifyC(t *testing.T) {
	myEquator := newContectedEquator()
	termA := contectedOperationTerm_t{
		ctxt:      0,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			operationTerm_t{
				operation: "h",
				args: []iTerm{
					variableTerm_t{
						symbol: "v",
					},
					variableTerm_t{
						symbol: "w",
					},
				},
			},
		},
	}
	termB := contectedOperationTerm_t{
		ctxt:      1,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}

	if termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if !termB.Unify(myEquator, termA) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}
}

func Test_ContectedUnifyD(t *testing.T) {
	myEquator := newContectedEquator()
	termA := contectedOperationTerm_t{
		ctxt:      0,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}
	termB := contectedOperationTerm_t{
		ctxt:      1,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}

	if termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if !termA.Unify(myEquator, termB) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}
}

func Test_ContectedUnifyE(t *testing.T) {
	myEquator := newContectedEquator()
	termA := contectedOperationTerm_t{
		ctxt:      0,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			operationTerm_t{
				operation: "h",
				args: []iTerm{
					variableTerm_t{
						symbol: "v",
					},
					variableTerm_t{
						symbol: "w",
					},
				},
			},
		},
	}

	termB := contectedOperationTerm_t{
		ctxt:      1,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			operationTerm_t{
				operation: "h",
				args: []iTerm{
					variableTerm_t{
						symbol: "v",
					},
					variableTerm_t{
						symbol: "w",
					},
				},
			},
		},
	}

	if termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if !termB.Unify(myEquator, termA) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}
}

func Test_ContectedUnifyF(t *testing.T) {
	myEquator := newContectedEquator()
	termA := contectedOperationTerm_t{
		ctxt:      0,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			operationTerm_t{
				operation: "h",
				args: []iTerm{
					variableTerm_t{
						symbol: "v",
					},
					variableTerm_t{
						symbol: "w",
					},
				},
			},
		},
	}

	termB := contectedOperationTerm_t{
		ctxt:      1,
		operation: "f",
		args: []iTerm{
			operationTerm_t{
				operation: "g",
				args: []iTerm{
					variableTerm_t{
						symbol: "u",
					},
				},
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}

	termC := contectedOperationTerm_t{
		ctxt:      2,
		operation: "f",
		args: []iTerm{
			operationTerm_t{
				operation: "g",
				args: []iTerm{
					operationTerm_t{
						operation: "k",
						args: []iTerm{
							variableTerm_t{
								symbol: "p",
							},
							variableTerm_t{
								symbol: "q",
							},
						},
					},
				},
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}

	if termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if termA.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if termB.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if !termB.Unify(myEquator, termA) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if termA.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if termB.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if !termC.Unify(myEquator, termA) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if !termB.EqualsTerm(myEquator, termC) {
		t.Fail()
	}
}

func Test_ContectedUnifyG(t *testing.T) {
	myEquator := newContectedEquator()
	termA := contectedOperationTerm_t{
		ctxt:      0,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			operationTerm_t{
				operation: "h",
				args: []iTerm{
					variableTerm_t{
						symbol: "v",
					},
					variableTerm_t{
						symbol: "w",
					},
				},
			},
		},
	}

	termB := contectedOperationTerm_t{
		ctxt:      1,
		operation: "f",
		args: []iTerm{
			operationTerm_t{
				operation: "g",
				args: []iTerm{
					variableTerm_t{
						symbol: "u",
					},
				},
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}

	termC := contectedOperationTerm_t{
		ctxt:      2,
		operation: "f",
		args: []iTerm{
			operationTerm_t{
				operation: "g",
				args: []iTerm{
					operationTerm_t{
						operation: "k",
						args: []iTerm{
							variableTerm_t{
								symbol: "p",
							},
							variableTerm_t{
								symbol: "q",
							},
						},
					},
				},
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}

	if termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if termA.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if termB.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if !termA.Unify(myEquator, termB) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if termA.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if termB.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if !termC.Unify(myEquator, termB) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termC) {
		t.Fail()
	}

	if !termB.EqualsTerm(myEquator, termC) {
		t.Fail()
	}
}

func Test_ReflexivityEqualsContectedTerm(t *testing.T) {
	myEquator := newContectedEquator()
	term := contectedOperationTerm_t{
		ctxt:      0,
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			variableTerm_t{
				symbol: "y",
			},
		},
	}

	if !term.EqualsTerm(myEquator, term) {
		t.Fail()
	}

	if !term.Unify(myEquator, term) {
		t.Fail()
	}

	if !term.EqualsTerm(myEquator, term) {
		t.Fail()
	}
}
