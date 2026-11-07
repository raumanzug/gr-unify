package simple

import (
	"github.com/raumanzug/gr-unify/unif"
)

func equalsTermVar[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S],
	NT unif.INonVarTermWithEquality[V, NT, S],
	S unif.ITermWithEquality[V, S],
](vt VT, equ unif.IEquator[V, S], term S) bool {
	v := vt.Variable()
	value, isNotEmpty := equ.Get(v)
	if isNotEmpty {
		return EqualsTerm[V, VT, NT, S](equ, value, term)
	}

	var normalizedTerm unif.ITerm[V, S]
	normalizedTerm = Normalform[V, VT, NT, S](equ, term)
	switch lbox := normalizedTerm.(type) {
	case VT:
		return v.Equals(lbox.Variable())
	default:
		return false
	}
}

func EqualsTerm[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S],
	NT unif.INonVarTermWithEquality[V, NT, S],
	S unif.ITermWithEquality[V, S],
](
	equ unif.IEquator[V, S],
	left unif.ITermWithEquality[V, S],
	right S,
) bool {
	switch box := left.(type) {
	case VT:
		return equalsTermVar[V, VT, NT, S](box, equ, right)
	case NT:
		return box.EqualsTermNonVar(equ, right)
	}

	return false
}
