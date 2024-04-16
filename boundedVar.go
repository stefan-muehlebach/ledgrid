package ledgrid

import "log"

// Auf alle Typen des Boundable-Interfaces koennen gebundene Variablen
// erstellt werden.
type Boundable interface {
	~int | ~int32 | ~int64 | ~float64 | ~float32
}

// Mit diesem Typ lassen sich Integer- und Float-Variablen erstellen, die auf
// ein bestimmtes Intervall eingeschraenkt sind. Die Methoden zum
// Inkrementieren (Incr) und Dekrementieren (Decr) beruecksichtigen die
// Intervallgrenzen.
type Bounded[T Boundable] struct {
	// Wird dieser Parameter ueber ein GUI oder TUI angezeigt, kann diese
	// Zeichenkette als Name verwendet werden.
	name string
	// Mit Cycle kann festgelegt werden, ob beim Erreichen, resp. Ueber- oder
	// Unterschreiten der Intervallgrenzen der Wert fix bleibt (false) oder
	// auf der anderen Seite des Intervalls beginnt (true).
	Cycle                   bool
	val, init, lb, ub, step T
	valPtr                  *T
	callback                func(oldVal, newVal T)
}

// Erstellt einen neuen eingeschraenkten Wert. Mit init, lb und ub kann
// der initiale Wert sowie die untere, resp. obere Schranke festgelegt werden.
func NewBounded[T Boundable](name string, init, lb, ub, inc T) *Bounded[T] {
	if lb > ub {
		log.Fatalf("lower bound must not be greater than upper bound (are '%v' and '%v' now)", lb, ub)
	}
	b := &Bounded[T]{}
	b.name = name
	b.Cycle = false
	b.init = init
	b.lb = lb
	b.ub = ub
	b.step = inc
	b.SetVal(init)
	b.valPtr = nil
	b.callback = nil
	return b
}

func (b *Bounded[T]) Name() string {
    return b.name
}

// Da der Zugriff auf den Wert einer Bounded-Variable immer geprueft werden
// muss, kann er nur ueber Methoden erfolgen.
func (b *Bounded[T]) Val() T {
	return b.val
}

func (b *Bounded[T]) SetVal(v T) {
	if v < b.lb || v > b.ub {
		log.Fatalf("value must be between lower and upper bound")
	}
	b.setVal(v)
}

func (b *Bounded[T]) Min() T {
    return b.lb
}

func (b *Bounded[T]) Max() T {
    return b.ub
}

func (b *Bounded[T]) Step() T {
    return b.step
}

// Mit BindVar kann eine Verbindung zu einer externen Variable hergestellt
// werden. Jede Aenderung an der Bounded-Variable wirkt sich autom. auf die
// externe Variable aus - nicht aber umgekehrt!
func (b *Bounded[T]) BindVar(ptr *T) {
	b.valPtr = ptr
	b.setVal(b.val)
}

// Mit SetCallback kann eine Funktion hinterlegt werden, die bei einer
// Aenderung der Variable aufgerufen werden soll. Als Parameter werden der
// Funktion der alte und der neue Wert der Variable uebergeben.
func (b *Bounded[T]) SetCallback(callback func(oldVal, newVal T)) {
	b.callback = callback
	b.callback(b.val, b.val)
}

// Mit Reset kann der Wert der Variable auf einen festgelegten Default
// (siehe Parameter init bei NewBounded) zurueckgesetzt werden.
func (b *Bounded[T]) Reset() {
	b.setVal(b.init)
}

// Inkrementiert den Wert der Variable um die Groesse v. Dabei werden die
// Grenzen (lb und ub) sowie die Einstellung Cycle beruecksichtigt.
func (b *Bounded[T]) Incr(v T) {
	b.add(v)
}

// Inkrementiert der Wert der Variable um den vordefinierten Wert incr.
func (b *Bounded[T]) Inc() {
	b.add(b.step)
}

// Dekrementiert den Wert der Variable um die Groesse v. Dabei werden die
// Grenzen (lb und ub) sowie die Einstellung Cycle beruecksichtigt.
func (b *Bounded[T]) Decr(v T) {
	b.add(-v)
}

// Dekrementiert der Wert der Variable um den vordefinierten Wert incr.
func (b *Bounded[T]) Dec() {
	b.add(-b.step)
}

// Diese interne Funktion setzt der Wert der Variable auf den Wert v.
// Allfällige Checks (ob v in [lb,up] liegt) müssen vorgängig gemacht werden!
// Über diese Methode wird auch die externe Variable aktualisiert und eine
// hinterlegte Callback-Methode aufgerufen.
func (b *Bounded[T]) setVal(v T) {
	oldVal := b.val
	b.val = v
	if b.valPtr != nil {
		*b.valPtr = b.val
	}
	if b.callback != nil {
		b.callback(oldVal, b.val)
	}
}

// Addiert zum aktuellen Wert der Variable den Wert v. Prüft dabei, ob die
// vordefinierten Grenzen (lb, ub) eingehalten werden, korrigiert ggf. den
// Wert und führt ggf. ein 'cycling' des Wertes durch.
func (b *Bounded[T]) add(v T) {
	if v > 0 && b.ub-v < b.val {
		if b.Cycle {
			b.setVal(b.lb)
		} else {
			b.setVal(b.ub)
		}
	} else if v < 0 && b.lb-v > b.val {
		if b.Cycle {
			b.setVal(b.ub)
		} else {
			b.setVal(b.lb)
		}
	} else {
		b.setVal(b.val + v)
	}
}
