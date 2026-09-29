# gr-unify

Hier ist das Werkzeug für die Go-Programmiersprache, um Lösungen für Unifikationsprobleme umzusetzen.  Derartige Probleme müssen zum Bau von Computeralgebra-Systemen, funktionalen Programmiersprachen, Theorembeweisern und allen Systemen gelöst werden, die mathematische Ausdrücke verarbeiten.

Ein Unifikationsproblem besteht darin, für zwei Terme die Variablen in ihnen so zu ersetzen, dass sie syntaktisch gleich werden.  Diese beiden Terme

	f(g(u, v), y) ~ f(x, h(w))

werden z.B. gleich, wenn man für x auf der rechten Seite g(u, v) einsetzt und für das y auf der linken Seite h(w), so daß auf beiden Seiten die Terme zu

	f(g(u, v), h(w)) ~ f(g(u, v), h(w))

werden.  Eine derartige Variablenersetzung bedeutet _Unifikation_.


## Der Weg zur Umsetzung

Wir behandeln nun Schritte zur Umsetzung:  Wir nehmen hier an, wir hätten schon Strukturen geschaffen, die Terme repräsentieren.  Wir verwenden einerseits eine Struktur, um Variablen zu repräsentieren etwa diese hier:

	type variable_t struct {
		symbol string
	}

so daß o.g. Variablen etwa in folgender Weise repräsentiert werden können:

	variable_t{
		symbol: "x",
	}

für eine Variable x.  Andererseits haben wir eine Struktur eingeführt, die Terme repräsentieren, die keine Variablen sind, zum Beispiel:

	operationTerm_t {
		operation: string
		args:      []iTerm
	}

`iTerm` stellt ein Go-Interface dar, das sowohl Variablen als auch Nicht-Variablen darstellen kann.  Um ein derartiges Interface einzuführen, verwenden wir folgendes aus unserer Bibliothek:

	type iTerm interface {
		unif.ITermWithEquality[
			variable_t, iTerm, simpleEquator_t,
		]
	}

`unif.ITermWithEquality` verwendet man für eine Termrepräsentation, die eine syntaktische Gleichheitsrelation `EqualsTerm` unterstützen soll.  In allen anderen Fällen steht an derselben Stelle. `unif.ITerm`.


## Equatoren

_Equatoren_ sind Go-Strukturen, die eine Zuordnung der Variablen zu den zu ihnen zugewiesenen Termen repräsentieren.  Der hier verwendete _Equator_ ist einfach ein Go-Map:

	 type simpleEquator_t map[string]iTerm

Die Lösung unseres Unifikationsproblems besteht einfach darin, die im Package `unif` enthaltenen Go-Interfaces zu implementieren.

Der Equator-Typ muß das `unif.IEquator`-Interface implementieren.

	func (recv simpleEquator_t) Get(v variable_t) (iTerm, bool) {
		retval, isOK := recv[v.symbol]

		return retval, isOK
	}

	func (recv simpleEquator_t) Set(v variable_t, x iTerm) bool {
		recv[v.symbol] = x

		return true
	}

Man kann damit Variablen Werte zuweisen und diese Werte auslesen.  Man beachte, daß die `Set`-Methode einen `bool`-Wert zurückgibt, der anzeigt, ob die `Set`-Methode wirklich den Wert zuweist oder die Zuweisung ablehnt.  Es ist manchmal notwendig, die Zuweisung abzulehnen, weil ansonsten Unifikationen wie die der Terme

	f(x) ~ f(f(x)

zu Endlosschleifen führen würden. `gr-unify` würde versuchen, `x` durch einen unendlichen Term `f(f(f(f(...x)))...)` zu ersetzen, wenn man die Zuweisung nicht verhindert, wenn auf beiden Seiten gleiche Variablen vorkommen.


# Implementation von Interfaces

`variable_t` und `operationTerm_t` müssen Interfaces implementieren, teilweise mit Hilfe von in dieser Bibliothek extra für diesen Zweck definierten Funktionen.

Zuerst benötigen wir eine Gleichheitsrelation auf dem Variablen repräsentierenden Typ:

	func (recv variable_t) Equals(v variable_t) bool {
		return recv.symbol == v.symbol
	}

Das implementiert das Interface `unif.IEquality`.

Die schon erwähnten Hilfsfunktionen nehmen uns die Arbeit zur Implementation der Interfaces zu einem gewissen Teil ab:

	func (recv variable_t) EqualsTerm(
		equ simpleEquator_t, term iTerm,
	) bool {
		return simple.EqualsTermVar(recv, equ, term)
	}

	func (recv variable_t) Normalform(
		equ simpleEquator_t,
	) iTerm {
		return simple.Normalform(recv, equ)
	}

	func (recv variable_t) Unify(
		equ simpleEquator_t, term iTerm,
	) bool {
		return simple.UnifyVar(recv, equ, term)
	}

	func (recv operationTerm_t) Unify(
		equ simpleEquator_t, term iTerm,
	) bool {
		return simple.UnifyNonVar(recv, equ, term)
	}

Im Wurzelverzeichnis dieses Projektes sind in der Datei `simple_test.go` auch alle anderen Definitionen enthalten, die der Entwickler selbst implementieren muß.


## Normalform

Normalform bedeutet, in einem Equator, die Kette der Variablen bis zu ihrem Endpunkt nachzuverfolgen.  Equatoren weisen Variablen andere Variablen zu.  Variablen stellen ja Terme dar und die kann man in Equatoren Variablen zuweisen. Die vorliegende Bibliothek stellt hierfür die Funktion `simple.Normalform` zur Verfügung, mit der die `Normalform`-Methode implementiert wird, s.o.


## `EqualsTerm`

`EqualsTerm` stellt eine Relation dar, die Terme auf syntaktische Gleichheit überprüft.  Für Nicht-Variablen muß das entsprechende Interface von Hand implementiert werden, s. o.g. Datei `simple_test.go`.


## `UnifyNonVar`-Methode

Diese Methode stellt den Kern der Unifikation dar.  Sie unifiziert zwei Nicht-Variablen.  Für alle anderen Fälle gibt es in der vorliegenden Bibliothek Funktionen.  Diese Methode muß vom Entwickler selbst implementiert werden.  An der Datei `simple_test.go` kann man sich ein Beispiel nehmen.


## `Term`-Methode

Die `Term`-Methode bettet Variablen und Nicht-Variablen repräsentierende Typen in den allgemeinen Terme repräsentierenden Typ ein.  Die Implementierung ist in den meisten Fällen trivial.


## Das Resultat

Wenn wir alle Interfaces implementiert haben, können wir Terme unifizieren.  Das geht mit der `Unify`-Methode.  `simple_test.go` enthält Testfälle, in denen diese Methode angewendet wird.
