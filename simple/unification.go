package simple

import (
	"github.com/raumanzug/gr-unify/unif"
)

func UnifyVar[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S, E],
	S unif.ITerm[V, S, E],
	E unif.IEquator[V, S],
](
	vt VT,
	equ E,
	term S,
) bool {
	v := vt.Variable()

	var normalizedTerm unif.ITerm[V, S, E]
	normalizedTerm = term.Normalform(equ)
	switch box := normalizedTerm.(type) {
	case VT:
		if v.Equals(box.Variable()) {
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
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S, E],
	NT unif.INonVarTerm[V, NT, S, E],
	S unif.ITerm[V, S, E],
	E unif.IEquator[V, S],
](
	left NT,
	equ E,
	right S,
) bool {
	var normalizedTerm unif.ITerm[V, S, E]
	normalizedTerm = right.Normalform(equ)
	switch box := normalizedTerm.(type) {
	case VT:
		return equ.Set(box.Variable(), left.Term())
	case NT:
		return box.UnifyNonVar(equ, left)
	default:
		return false
	}
}

func Normalform[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S, E],
	S unif.ITerm[V, S, E],
	E unif.IEquator[V, S],
](
	vt VT,
	equ E,
) S {
	v := vt.Variable()
	value, isNotEmpty := equ.Get(v)
	if isNotEmpty {
		return value.Normalform(equ)
	} else {
		return vt.Term()
	}
}
