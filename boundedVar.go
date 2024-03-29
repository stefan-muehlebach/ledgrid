package ledgrid

import "log"

// Auf alle Typen des Boundable-Interfaces koennen gebundene Variablen
// erstellt werden.
type Boundable interface {
	~int | ~float64 | ~float32
}

// Mit diesem Typ lassen sich Integer- und Float-Variablen erstellen, die auf
// ein bestimmtes Intervall eingeschraenkt sind. Die Methoden zum
// Inkrementieren (Incr) und Dekrementieren (Decr) beruecksichtigen die
// Intervallgrenzen.
type Bounded[T Boundable] struct {
	val, init, lb, ub, step T
	ptr                    *T
    callback func(oldVal, newVal T)
	// Mit Cycle kann festgelegt werden, ob beim Erreichen, resp. Ueber- oder
	// Unterschreiten der Intervallgrenzen der Wert fix bleibt (false) oder
	// auf der anderen Seite des Intervalls beginnt (true).
	Cycle bool
	// Wird dieser Parameter ueber ein GUI oder TUI angezeigt, kann diese
	// Zeichenkette als Name verwendet werden.
	Name string
}

// Erstellt einen neuen eingeschraenkten Wert. Mit init, lb und ub kann
// der initiale Wert sowie die untere, resp. obere Schranke festgelegt werden.
func NewBounded[T Boundable](init, lb, ub, inc T) *Bounded[T] {
	if lb >= ub {
		log.Fatalf("lower bound must be less than upper bound")
	}
	b := &Bounded[T]{}
	b.init = init
	b.lb = lb
	b.ub = ub
    b.step = inc
	b.SetVal(init)
	b.ptr = nil
    b.callback = nil
	b.Cycle = false
	b.Name = ""
	return b
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

// Mit BindVar kann eine Verbindung zu einer externen Variable hergestellt
// werden. Jede Aenderung an der Bounded-Variable wirkt sich autom. auf die
// externe Variable aus - nicht aber umgekehrt!
func (b *Bounded[T]) BindVar(ptr *T) {
	b.ptr = ptr
	b.setVal(b.val)
}

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

func (b *Bounded[T]) Inc() {
    b.add(b.step)
}

// Dekrementiert den Wert der Variable um die Groesse i. Dabei werden die
// Grenzen (lb und ub) sowie die Einstellung Cycle beruecksichtigt.
func (b *Bounded[T]) Decr(v T) {
	b.add(-v)
}

func (b *Bounded[T]) Dec() {
    b.add(-b.step)
}

func (b *Bounded[T]) setVal(v T) {
    oldVal := b.val
	b.val = v
	if b.ptr != nil {
		*b.ptr = b.val
	}
    if b.callback != nil {
        b.callback(oldVal, b.val)
    }
}

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
