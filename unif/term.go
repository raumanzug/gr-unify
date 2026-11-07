package unif

type IEquality[V any] interface {
	Equals(V) bool
}

type IEquator[V any, S any] interface {
	Get(V) (S, bool)

	Set(V, S) bool
}

type IVarTerm[
	V IEquality[V],
	S ITerm[V, S],
] interface {
	ITerm[V, S]

	Variable() V
}

type INonVarTerm[
	V any,
	N any,
	S ITerm[V, S],
] interface {
	ITerm[V, S]

	UnifyNonVar(IEquator[V, S], N) bool
}

type ITerm[
	V any,
	S ITerm[V, S],
] interface {
	Term() S
}

type INonVarTermWithEquality[
	V any,
	N any,
	S ITermWithEquality[V, S],
] interface {
	INonVarTerm[V, N, S]

	EqualsTermNonVar(IEquator[V, S], S) bool
}

type ITermWithEquality[
	V any,
	S ITermWithEquality[V, S],
] interface {
	ITerm[V, S]
}
