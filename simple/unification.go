package simple

import (
	"fmt"

	"github.com/raumanzug/gr-unify/unif"
)

func unifyVar[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S],
	NT unif.INonVarTerm[V, NT, S],
	S unif.ITerm[V, S],
](
	vt VT,
	equ unif.IEquator[V, S],
	term S,
) bool {
	v := vt.Variable()

	var normalizedTerm unif.ITerm[V, S]
	normalizedTerm = Normalform[V, VT, NT, S](equ, term)
	switch box := normalizedTerm.(type) {
	case VT:
		if v.Equals(box.Variable()) {
			return true
		}
		value, isNotEmpty := equ.Get(v)
		if isNotEmpty {
			return Unify[V, VT, NT, S](equ, value, box.Term())
		}
		return equ.Set(v, box.Term())
	case S:
		value, isNotEmpty := equ.Get(v)
		if isNotEmpty {
			return Unify[V, VT, NT, S](equ, value, box)
		} else {
			return equ.Set(v, box)
		}
	}

	return false
}

func unifyNonVar[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S],
	NT unif.INonVarTerm[V, NT, S],
	S unif.ITerm[V, S],
](
	left NT,
	equ unif.IEquator[V, S],
	right S,
) bool {
	var normalizedTerm unif.ITerm[V, S]
	normalizedTerm = Normalform[V, VT, NT, S](equ, right)
	switch box := normalizedTerm.(type) {
	case VT:
		return equ.Set(box.Variable(), left.Term())
	case NT:
		return box.UnifyNonVar(equ, left)
	default:
		return false
	}
}

func normalformVar[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S],
	NT unif.INonVarTerm[V, NT, S],
	S unif.ITerm[V, S],
](
	vt VT,
	equ unif.IEquator[V, S],
) S {
	v := vt.Variable()
	value, isNotEmpty := equ.Get(v)
	if isNotEmpty {
		return Normalform[V, VT, NT, S](equ, value)
	} else {
		return vt.Term()
	}
}

func Normalform[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S],
	NT unif.INonVarTerm[V, NT, S],
	S unif.ITerm[V, S],
](
	equ unif.IEquator[V, S],
	term unif.ITerm[V, S],
) S {
	switch box := term.(type) {
	case VT:
		return normalformVar[V, VT, NT, S](box, equ)
	case NT:
		return box.Term()
	}

	panic(fmt.Sprintf(
		"type %T of term arg in simple.Normalform not permitted.",
		term,
	))
}

func Unify[
	V unif.IEquality[V],
	VT unif.IVarTerm[V, S],
	NT unif.INonVarTerm[V, NT, S],
	S unif.ITerm[V, S],
](
	equ unif.IEquator[V, S],
	left unif.ITerm[V, S],
	right S,
) bool {
	switch box := left.(type) {
	case VT:
		return unifyVar[V, VT, NT, S](box, equ, right)
	case NT:
		return unifyNonVar[V, VT, NT, S](box, equ, right)
	}

	return false
}
