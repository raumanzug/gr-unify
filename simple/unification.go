package simple

import (
	"github.com/raumanzug/gr-unify/unif"
)

func UnifyVar[
	V unif.IVarTerm[V, S, E],
	S unif.ITerm[V, S, E],
	E unif.IEquator[V, S],
](
	v V,
	equ E,
	term S,
) bool {
	var normalizedTerm unif.ITerm[V, S, E]
	normalizedTerm = term.Normalform(equ)
	switch box := normalizedTerm.(type) {
	case V:
		if v.Equals(box) {
			return true
		}
		value, isNotEmpty := equ.Get(v)
		if isNotEmpty {
			return value.Unify(equ, box.Term())
		}
		return equ.Set(v, box.Term())
	case S:
		value, isNotEmpty := equ.Get(v)
		if isNotEmpty {
			return value.Unify(equ, box)
		} else {
			return equ.Set(v, box)
		}
	}

	return false
}

func UnifyNonVar[
	V unif.IVarTerm[V, S, E],
	N unif.INonVarTerm[V, N, S, E],
	S unif.ITerm[V, S, E],
	E unif.IEquator[V, S],
](
	left N,
	equ E,
	right S,
) bool {
	var normalizedTerm unif.ITerm[V, S, E]
	normalizedTerm = right.Normalform(equ)
	switch box := normalizedTerm.(type) {
	case V:
		return equ.Set(box, left.Term())
	case N:
		return box.UnifyNonVar(equ, left)
	default:
		return false
	}
}

func Normalform[
	V unif.IVarTerm[V, S, E],
	S unif.ITerm[V, S, E],
	E unif.IEquator[V, S],
](
	v V,
	equ E,
) S {
	value, isNotEmpty := equ.Get(v)
	if isNotEmpty {
		return value.Normalform(equ)
	} else {
		return v.Term()
	}
}
