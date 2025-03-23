package ledgrid

import (
	"fmt"
	"math"
	"time"
)

// Simple way to measure the number and duration of events, simply by enclose
// the code to be measured by Start() and Stop().
type Stopwatch struct {
	Num             int
	Total, Max, Min time.Duration
	t               time.Time
}

func NewStopwatch() *Stopwatch {
	s := &Stopwatch{}
	s.Max = math.MinInt64
	s.Min = math.MaxInt64
	return s
}

// Starts the stopwatch. If the stopwatch is already running, the previous
// starting time is cancelled.
func (s *Stopwatch) Start() {
	s.t = time.Now()
}

// Stops the stopwatch and updates the internal variables. If the stopwatch is
// stopped, this method has no effect.
func (s *Stopwatch) Stop() {
	d := time.Since(s.t)
	if d > s.Max {
		s.Max = d
	}
	if d < s.Min {
		s.Min = d
	}
	s.Total += d
	s.Num += 1
}

// Setzt die gemessene Dauer auf 0 und die Anzahl Messungen ebenfalls.
func (s *Stopwatch) Reset() {
	s.Total = 0
	s.Max = math.MinInt64
	s.Min = math.MaxInt64
	s.Num = 0
}

// Retourniert die Anzahl Messungen.
// func (s *Stopwatch) Num() int {
// 	return s.Num
// }

// Retourniert die totale Messdauer.
// func (s *Stopwatch) Total() time.Duration {
// 	return s.Total
// }

// Retourniert die totale Messdauer.
// func (s *Stopwatch) Min() time.Duration {
// 	return s.Min
// }

// Retourniert die totale Messdauer.
// func (s *Stopwatch) Max() time.Duration {
// 	return s.Max
// }

// Berechnet die durchschnittliche Messdauer (also den Quotienten von
// Total() / Num()).
func (s *Stopwatch) Avg() time.Duration {
	if s.Num == 0 {
		return 0
	}
	return s.Total / time.Duration(s.Num)
}

// func (s *Stopwatch) Stats() (int, time.Duration, time.Duration) {
// 	return s.Num, s.Total, s.Avg()
// }

func (s *Stopwatch) String() string {
	return fmt.Sprintf("%d calls; %v in total; %v per call, min: %v, max: %v",
		s.Num, s.Total, s.Avg(), s.Min, s.Max)
}
