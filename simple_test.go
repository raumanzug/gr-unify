package gr

import (
	"slices"
	"testing"

	"github.com/raumanzug/gr-unify/simple"
	"github.com/raumanzug/gr-unify/unif"
)

type simpleEquator_t map[string]iTerm

type iTerm interface {
	unif.ITermWithEquality[
		variable_t,
		iTerm,
		simpleEquator_t,
	]
}

type variable_t struct {
	symbol string
}

func (recv simpleEquator_t) Get(v variable_t) (iTerm, bool) {
	retval, isOK := recv[v.symbol]

	return retval, isOK
}

func newSimpleEquator() simpleEquator_t {
	return make(map[string]iTerm)
}

func (recv simpleEquator_t) Set(v variable_t, x iTerm) bool {
	recv[v.symbol] = x

	return true
}

type operationTerm_t struct {
	operation string
	args      []iTerm
}

func (recv variable_t) Equals(v variable_t) bool {
	return recv.symbol == v.symbol
}

func (recv variable_t) EqualsTerm(
	equ simpleEquator_t,
	term iTerm,
) bool {
	return simple.EqualsTermVar(recv, equ, term)
}

func (recv variable_t) Normalform(
	equ simpleEquator_t,
) iTerm {
	return simple.Normalform(recv, equ)
}

func (recv variable_t) Term() iTerm {
	return recv
}

func (recv variable_t) Unify(
	equ simpleEquator_t,
	term iTerm,
) bool {
	return simple.UnifyVar(recv, equ, term)
}

func (recv operationTerm_t) EqualsTerm(
	equ simpleEquator_t,
	term iTerm,
) bool {
	var normalizedTerm iTerm
	normalizedTerm = term.Normalform(equ)
	switch box := normalizedTerm.(type) {
	case variable_t:
		return false
	case operationTerm_t:
		if recv.operation != box.operation {
			return false
		}

		return slices.EqualFunc(
			recv.args,
			box.args,
			func(left iTerm, right iTerm) bool {
				return left.EqualsTerm(equ, right)
			},
		)
	}

	return false
}

func (recv operationTerm_t) Normalform(
	equ simpleEquator_t,
) iTerm {
	return recv
}

func (recv operationTerm_t) Term() iTerm {
	return recv
}

func (recv operationTerm_t) Unify(
	equ simpleEquator_t,
	term iTerm,
) bool {
	return simple.UnifyNonVar(recv, equ, term)
}

func (recv operationTerm_t) UnifyNonVar(
	equ simpleEquator_t,
	term operationTerm_t,
) bool {
	if recv.operation != term.operation {
		return false
	}

	return slices.EqualFunc(
		recv.args,
		term.args,
		func(left iTerm, right iTerm) bool {
			return left.Unify(equ, right)
		},
	)
}

func Test_Unify(t *testing.T) {
	myEquator := newSimpleEquator()
	termA := operationTerm_t{
		operation: "f",
		args: []iTerm{
			variable_t{
				symbol: "x",
			},
			variable_t{
				symbol: "x",
			},
		},
	}
	termB := operationTerm_t{
		operation: "f",
		args: []iTerm{
			variable_t{
				symbol: "x",
			},
			variable_t{
				symbol: "y",
			},
		},
	}

	if termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if termB.EqualsTerm(myEquator, termA) {
		t.Fail()
	}

	if !termB.Unify(myEquator, termA) {
		t.Fail()
	}

	if !termA.EqualsTerm(myEquator, termB) {
		t.Fail()
	}

	if !termB.EqualsTerm(myEquator, termA) {
		t.Fail()
	}
}

func Test_ReflexivityEqualsTerm(t *testing.T) {
	myEquator := newSimpleEquator()
	term := operationTerm_t{
		operation: "f",
		args: []iTerm{
			variable_t{
				symbol: "x",
			},
			variable_t{
				symbol: "y",
			},
		},
	}

	if !term.EqualsTerm(myEquator, term) {
		t.Fail()
	}
}
