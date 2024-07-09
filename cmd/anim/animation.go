package main

import (
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/geom"
)

//----------------------------------------------------------------------------

type AnimationCurve func(t float64) float64

func AnimationLinear(t float64) float64 {
	return t
}

func AnimationEaseInOutSelf(t float64) float64 {
	return 3*t*t - 2*t*t*t
}

func AnimationEaseIn(t float64) float64 {
	return t * t
}

func AnimationEaseOut(t float64) float64 {
	return t * (2 - t)
}

func AnimationEaseInOut(t float64) float64 {
	if t <= 0.5 {
		return t * t * 2
	}
	return -1 + (4-t*2)*t
}

const (
	AnimationRepeatForever = -1
)

// Das Interface fuer jede Art von Animation (bis jetzt zumindest).
type Animation interface {
	// Startet die Animation mit jenen Parametern, die zum Startzeitpunkt
	// aktuell sind. Ist die Animaton bereits am Laufen ist diese Methode
	// ein no-op.
	Start()
	// Haelt die Animation an, laesst sie jedoch in der Animations-
	// Queue.
	Stop()
	// Setzt eine mit Stop() angehaltene Animation wieder fort.
	Continue()
	// Liefert true, falls die Animation mittels Stop() angehalten wurde.
	IsStopped() bool
	// Diese Methode wird von einem Controller oder dgl. waehrend der Laufzeit
	// der Animation aufgerufen. [t] ist der Zeitpunkt, an welchem mit dem
	// Update begonnen wurde.
	Update(t time.Time) bool

}

// Embeddable mit in allen Animationen benoetigen Variablen und Code.
// Erleichert das Erstellen von neuen Animationen gewaltig.
type AnimationEmbed struct {
	// Diese Parameter koennen von einem aufrufenden Programm veraendert
	// werden, sobald jedoch die Animation gestartet ist, sind Aenderungen
	// an diesen Variablen wirkungslos.
	AutoReverse bool
	Curve       AnimationCurve
	Duration    time.Duration
	RepeatCount int

	reverse          bool
	start, stop, end time.Time
	total            int64
	repeatsLeft      int
	running          bool
	tick             func(t float64)
}

func (a *AnimationEmbed) Init(d time.Duration) {
	a.AutoReverse = false
	a.Curve = AnimationEaseInOut
	a.Duration = d
	a.RepeatCount = 0
}

func (a *AnimationEmbed) Start() {
	if a.running {
		return
	}
	a.start = time.Now()
	a.end = a.start.Add(a.Duration)
	a.total = a.end.Sub(a.start).Milliseconds()
	a.repeatsLeft = a.RepeatCount
	a.reverse = false
	a.running = true
	AnimCtrl.AddAnim(a)
}

func (a *AnimationEmbed) Stop() {
	if !a.running {
		return
	}
	a.stop = time.Now()
	a.running = false
}

func (a *AnimationEmbed) Continue() {
	if a.running {
		return
	}
	dt := time.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

func (a *AnimationEmbed) IsStopped() bool {
	return !a.running
}

// Diese Methode ist fuer die korrekte Abwicklung (Beachtung von Reverse und
// RepeatCount, etc) einer Animation zustaendig. Wenn die Animation zu Ende
// ist, retourniert Update false.
func (a *AnimationEmbed) Update(t time.Time) bool {
	if t.After(a.end) {
		if a.reverse {
			a.tick(a.Curve(0.0))
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			}
			a.reverse = false
		} else {
			a.tick(a.Curve(1.0))
			if a.AutoReverse {
				a.reverse = true
			}
		}
		if !a.reverse {
			if a.repeatsLeft == 0 {
				a.running = false
				return false
			}
			if a.repeatsLeft > 0 {
				a.repeatsLeft--
			}
		}
		a.start = t
		a.end = a.start.Add(a.Duration)
		return true
	}

	delta := t.Sub(a.start).Milliseconds()
	val := float64(delta) / float64(a.total)
	if a.reverse {
		a.tick(a.Curve(1.0 - val))
	} else {
		a.tick(a.Curve(val))
	}
	return true
}

//----------------------------------------------------------------------------

// Animation fuer einen Farbverlauf.
type ColorAnimation struct {
	AnimationEmbed
	cp     *color.Color
	c1, c2 color.Color
}

func NewColorAnimation(cp *color.Color, c2 color.Color, d time.Duration) *ColorAnimation {
	a := &ColorAnimation{}
	a.Init(d)
	a.cp = cp
	a.c1, a.c2 = *cp, c2
    a.tick = a.Tick
	return a
}

func (a *ColorAnimation) Tick(t float64) {
	*a.cp = a.c1.Interpolate(a.c2, t)
}

// Animation fuer einen Wechsel der Position, resp. Veraenderung der Groesse.
type PositionAnimation struct {
	AnimationEmbed
	pp     *geom.Point
	p1, p2 geom.Point
}

func NewPositionAnimation(pp *geom.Point, p2 geom.Point, d time.Duration) *PositionAnimation {
	a := &PositionAnimation{}
	a.Init(d)
	a.pp = pp
	a.p1, a.p2 = *pp, p2
    a.tick = a.Tick
	return a
}

func (a *PositionAnimation) Tick(t float64) {
	p := a.p1.Interpolate(a.p2, t)
	*a.pp = p
}

var (
	NewSizeAnimation = NewPositionAnimation
)

// Animation fuer das Fahren entlang eines Pfades.
type PathAnimation struct {
	AnimationEmbed
	pp    *geom.Point
	refPt geom.Point
	size  geom.Point
	fnc   PathFunction
}

type PathFunction func(t float64) geom.Point

func NewPathAnimation(pp *geom.Point, fnc PathFunction, size geom.Point, d time.Duration) *PathAnimation {
	a := &PathAnimation{}
	a.Init(d)
	a.pp = pp
	a.refPt = *pp
	a.size = size
	a.fnc = fnc
    a.tick = a.Tick
	return a
}

func (a *PathAnimation) Tick(t float64) {
	pt := a.fnc(t)
	dp := geom.Point{pt.X * a.size.X, pt.Y * a.size.Y}
	*a.pp = a.refPt.Add(dp)
}

func circlePathFunc(t float64) geom.Point {
	phi := t * 2 * math.Pi
	return geom.Point{0.5 * (1 - math.Cos(phi)), 0.5 * math.Sin(phi)}
}

// Allgemeine Parameter-Animation.
// type FloatAnimation struct {
// 	AnimationEmbed
// 	fp     *float64
// 	f1, f2 float64
// }

// func NewFloatAnimation(fp *float64, f2 float64, d time.Duration) *FloatAnimation {
// 	a := &FloatAnimation{}
//     a.Init(d)
// 	a.fp = fp
// 	a.f1, a.f2 = *fp, f2
// 	a.tick = a.UpdateFloat
// 	return a
// }

// func (a *FloatAnimation) UpdateFloat(t float64) {
// 	f := (1.0-t)*a.f1 + t*a.f2
// 	*a.fp = f
// }
