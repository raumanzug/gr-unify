package simple

import (
	"github.com/raumanzug/gr-unify/unif"
)

func EqualsTermVar[
	V unif.IVarTerm[V, S, E],
	S unif.ITermWithEquality[V, S, E],
	E unif.IEquator[V, S],
](v V, equ E, term S) bool {
	value, isNotEmpty := equ.Get(v)
	if isNotEmpty {
		return value.EqualsTerm(equ, term)
	}

	var normalizedTerm unif.ITerm[V, S, E]
	normalizedTerm = term.Normalform(equ)
	switch lbox := normalizedTerm.(type) {
	case V:
		return v.Equals(lbox)
	default:
		return false
	}
}
