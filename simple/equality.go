package simple

import (
	"github.com/raumanzug/gr-unify/unif"
)

func EqualsTermVar[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S, E],
	S unif.ITermWithEquality[V, S, E],
	E unif.IEquator[V, S],
](vt VT, equ E, term S) bool {
	v := vt.Variable()
	value, isNotEmpty := equ.Get(v)
	if isNotEmpty {
		return value.EqualsTerm(equ, term)
	}

	var normalizedTerm unif.ITerm[V, S, E]
	normalizedTerm = term.Normalform(equ)
	switch lbox := normalizedTerm.(type) {
	case VT:
		return v.Equals(lbox.Variable())
	default:
		return false
	}
}
