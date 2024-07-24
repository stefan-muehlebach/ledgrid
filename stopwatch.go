package ledgrid

import (
	"fmt"
	"time"
)

// Dieser Typ dient der Zeitmessung.
type Stopwatch struct {
	t time.Time
	d time.Duration
	n int
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{}
}

// Mit Start wird eine neue Messung begonnen.
func (s *Stopwatch) Start() {
	s.t = time.Now()
}

// Stop beendet die Messung und aktualisiert die Variablen, welche die totale
// Messdauer als auch die Anzahl Messungen enthalten.
func (s *Stopwatch) Stop() {
	s.d += time.Since(s.t)
	s.n += 1
}

// Setzt die gemessene Dauer auf 0 und die Anzahl Messungen ebenfalls.
func (s *Stopwatch) Reset() {
	s.d = 0
	s.n = 0
}

// Retourniert die Anzahl Messungen.
func (s *Stopwatch) Num() int {
	return s.n
}

// Retourniert die totale Messdauer.
func (s *Stopwatch) Total() time.Duration {
	return s.d
}

// Berechnet die durchschnittliche Messdauer (also den Quotienten von
// Total() / Num()).
func (s *Stopwatch) Avg() time.Duration {
	if s.n == 0 {
		return 0
	}
	return s.d / time.Duration(s.n)
}

func (s *Stopwatch) Stats() (int, time.Duration, time.Duration) {
    return s.Num(), s.Total(), s.Avg()
}

func (s *Stopwatch) String() string {
    return fmt.Sprintf("%d calls; %v in total; %v per call", s.Num(), s.Total(), s.Avg())
}
