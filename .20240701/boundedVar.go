package ledgrid

import (
	"log"

	"fyne.io/fyne/v2/data/binding"
)

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
	NameableEmbed
	// Mit Cycle kann festgelegt werden, ob beim Erreichen, resp. Ueber- oder
	// Unterschreiten der Intervallgrenzen der Wert fix bleibt (false) oder
	// auf der anderen Seite des Intervalls beginnt (true).
	Cycle bool
	// Das sind die Variablen fuer den aktuellen Wert (val), den initialen
	// Wert (init), die untere und obere Schranke (lb und ub) sowie die
	// Schrittweite der Methoden Inc() und Dec().
	val, init, lb, ub, step T
	// Ist die Variable mit einer externen Variable verbunden, ist dies der
	// Pointer auf die externe Variable.
	valPtr *T
	// Dies ist der Pointer auf eine Funktion, welche bei Aenderung der
	// Variable aufgerufen wird. oldVal und newVal sind die Werte der Variable
	// vor, resp. nach der Aenderung.
	callback func(oldVal, newVal T)

	listeners []binding.DataListener
}

// Erstellt einen neuen eingeschraenkten Wert. Mit name kann der Variable
// einen Namen verliehen werden. Mit init, lb und ub kann der initiale Wert
// sowie die untere, resp. obere Schranke festgelegt werden und mit inc wird
// die Schrittweite der Methoden Inc() und Dec() spezifiziert.
func NewBounded[T Boundable](name string, init, lb, ub, inc T) *Bounded[T] {
	if lb > ub {
		log.Fatalf("lower bound must not be greater than upper bound (are '%v' and '%v' now)", lb, ub)
	}
	b := &Bounded[T]{}
	b.NameableEmbed.Init(name)
	b.Cycle = false
	b.init = init
	b.lb = lb
	b.ub = ub
	b.step = inc
	b.SetVal(init)
	b.valPtr = nil
	b.callback = nil
	b.listeners = make([]binding.DataListener, 0)
	return b
}

// Die folgenden vier Methoden wurden implementiert, um das DataItem-
// Interface aus dem binding-Package von fyne.io zu implementieren.
func (b *Bounded[T]) Get() (T, error) {
	return b.Val(), nil
}

func (b *Bounded[T]) Set(v T) error {
	b.SetVal(v)
	return nil
}

func (b *Bounded[T]) AddListener(l binding.DataListener) {
	// log.Printf("AddListener(%T)", l)
	b.listeners = append(b.listeners, l)
}

func (b *Bounded[T]) RemoveListener(l binding.DataListener) {
	var idx int
	var lsnr binding.DataListener

	// log.Printf("RemoveListener(%T)", l)
	for idx, lsnr = range b.listeners {
		if l == lsnr {
			break
		}
	}
	b.listeners = append(b.listeners[:idx], b.listeners[idx+1:]...)
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

// Min, Max und Step retournieren die untere, resp. obere Schranke sowie die
// Schrittweite.
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
	for _, lsnr := range b.listeners {
		if lsnr != nil {
			lsnr.DataChanged()
		}
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
