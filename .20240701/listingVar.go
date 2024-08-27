package ledgrid

import (
	"slices"
)

type Listable interface {
	comparable
}


// Mit diesem Typ lassen sich Integer- und Float-Variablen erstellen, die auf
// ein bestimmtes Intervall eingeschraenkt sind. Die Methoden zum
// Inkrementieren (Incr) und Dekrementieren (Decr) beruecksichtigen die
// Intervallgrenzen.
type Listing[T Listable] struct {
	// Wird dieser Parameter ueber ein GUI oder TUI angezeigt, kann diese
	// Zeichenkette als Name verwendet werden.
	name string
	// Mit Cycle kann festgelegt werden, ob beim Erreichen, resp. Ueber- oder
	// Unterschreiten der Intervallgrenzen der Wert fix bleibt (false) oder
	// auf der anderen Seite des Intervalls beginnt (true).
	Cycle    bool
	list     []T
	idx      int
	val      T
	valPtr   *T
	callback func(oldVal, newVal T)
}

// Erstellt einen neuen eingeschraenkten Wert. Mit init, lb und ub kann
// der initiale Wert sowie die untere, resp. obere Schranke festgelegt werden.
func NewListing[T Listable](name string, list []T) *Listing[T] {
	l := &Listing[T]{}
	l.name = name
	l.Cycle = false
	l.list = slices.Clone(list)
	l.idx = 0
	l.valPtr = nil
	l.callback = nil
	return l
}

func (l *Listing[T]) Name() string {
	return l.name
}

func (l *Listing[T]) Val() T {
    return l.list[l.idx]
}

func (l *Listing[T]) SetVal(v T) {
    if slices.Contains(l.list, v) {
        l.idx = slices.Index(l.list, v)
    }
}

func (l *Listing[T]) Next() T {
    l.idx++
    if l.idx >= len(l.list) {
        if l.Cycle {
            l.idx %= len(l.list)
        } else {
            l.idx = len(l.list)-1
        }
    }
    return l.Val()
}

func (l *Listing[T]) Prev() T {
    l.idx--
    if l.idx < 0 {
        if l.Cycle {
            l.idx += len(l.list)
        } else {
            l.idx = 0
        }
    }
    return l.Val()
}
