package unif

type IEquality[V any] interface {
	Equals(V) bool
}

type IEquator[V any, S any] interface {
	Get(V) (S, bool)

	Set(V, S) bool
}

type IVarTerm[
	V any,
	S ITerm[V, S, E],
	E IEquator[V, S],
] interface {
	IEquality[V]
	ITerm[V, S, E]

	Term() S
}

type INonVarTerm[
	V any,
	N any,
	S ITerm[V, S, E],
	E IEquator[V, S],
] interface {
	ITerm[V, S, E]

	Term() S
	UnifyNonVar(E, N) bool
}

type ITerm[
	V any,
	S ITerm[V, S, E],
	E IEquator[V, S],
] interface {
	Normalform(E) S
	Unify(E, S) bool
}

type ITermWithEquality[
	V any,
	S ITermWithEquality[V, S, E],
	E IEquator[V, S],
] interface {
	ITerm[V, S, E]

	EqualsTerm(E, S) bool
}
