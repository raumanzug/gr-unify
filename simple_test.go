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
	]
}

type variable_t string

type variableTerm_t struct {
	symbol string
}

func (recv simpleEquator_t) Get(v variable_t) (iTerm, bool) {
	retval, isOK := recv[string(v)]

	return retval, isOK
}

func newSimpleEquator() simpleEquator_t {
	return make(map[string]iTerm)
}

func (recv simpleEquator_t) Set(v variable_t, x iTerm) bool {
	recv[string(v)] = x

	return true
}

type operationTerm_t struct {
	operation string
	args      []iTerm
}

func (recv variable_t) Equals(v variable_t) bool {
	return recv == v
}

func (recv variableTerm_t) Term() iTerm {
	return recv
}

func (recv variableTerm_t) Variable() variable_t {
	return variable_t(recv.symbol)
}

func (recv operationTerm_t) EqualsTermNonVar(
	equ unif.IEquator[variable_t, iTerm],
	term iTerm,
) bool {
	var normalizedTerm iTerm
	normalizedTerm = simple.Normalform[
		variable_t,
		variableTerm_t,
		operationTerm_t,
		iTerm,
	](equ, term)
	switch box := normalizedTerm.(type) {
	case variableTerm_t:
		return false
	case operationTerm_t:
		if recv.operation != box.operation {
			return false
		}

		return slices.EqualFunc(
			recv.args,
			box.args,
			func(left iTerm, right iTerm) bool {
				return simple.EqualsTerm[
					variable_t,
					variableTerm_t,
					operationTerm_t,
					iTerm,
				](equ, left, right)
			},
		)
	}

	return false
}

func (recv operationTerm_t) Term() iTerm {
	return recv
}

func (recv operationTerm_t) UnifyNonVar(
	equ unif.IEquator[variable_t, iTerm],
	term operationTerm_t,
) bool {
	if recv.operation != term.operation {
		return false
	}

	return slices.EqualFunc(
		recv.args,
		term.args,
		func(left iTerm, right iTerm) bool {
			return simple.Unify[
				variable_t,
				variableTerm_t,
				operationTerm_t,
				iTerm,
			](equ, left, right)
		},
	)
}

func Test_Unify(t *testing.T) {
	myEquator := newSimpleEquator()
	termA := operationTerm_t{
		operation: "f",
		args: []iTerm{
			variableTerm_t{
				symbol: "x",
			},
			variableTerm_t{
				symbol: "x",
			},
		},
	}
	termB := operationTerm_t{
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

	if simple.EqualsTerm[
		variable_t,
		variableTerm_t,
		operationTerm_t,
		iTerm,
	](myEquator, termA, termB) {
		t.Fail()
	}

	if simple.EqualsTerm[
		variable_t,
		variableTerm_t,
		operationTerm_t,
		iTerm,
	](myEquator, termB, termA) {
		t.Fail()
	}

	if !simple.Unify[
		variable_t,
		variableTerm_t,
		operationTerm_t,
		iTerm,
	](myEquator, termB, termA) {
		t.Fail()
	}

	if !simple.EqualsTerm[
		variable_t,
		variableTerm_t,
		operationTerm_t,
		iTerm,
	](myEquator, termA, termB) {
		t.Fail()
	}

	if !simple.EqualsTerm[
		variable_t,
		variableTerm_t,
		operationTerm_t,
		iTerm,
	](myEquator, termB, termA) {
		t.Fail()
	}
}

func Test_ReflexivityEqualsTerm(t *testing.T) {
	myEquator := newSimpleEquator()
	term := operationTerm_t{
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

	if !simple.EqualsTerm[
		variable_t,
		variableTerm_t,
		operationTerm_t,
		iTerm,
	](myEquator, term, term) {
		t.Fail()
	}
}
