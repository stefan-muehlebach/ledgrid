package main

import (
	"github.com/stefan-muehlebach/ledgrid"
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/geom"
)

// Mit dem Funktionstyp [AnimationCurve] kann der Verlauf einer Animation
// beeinflusst werden. Der Parameter [t] ist ein Wert im Intervall [0,1]
// und zeigt an, wo sich die Animation gerade befindet (t=0: Animation
// wurde eben gestartet; t=1: Animation ist zu Ende). Der Rueckgabewert
// ist ebenfalls ein Wert im Intervall [0,1] und hat die gleiche Bedeutung
// wie [t].
type AnimationCurve func(t float64) float64

// Bezeichnet eine lineare Animation.
func AnimationLinear(t float64) float64 {
	return t
}

// Beginnt langsam und nimmt immer mehr an Fahrt auf.
func AnimationEaseIn(t float64) float64 {
	return t * t
}

// Beginnt schnell und bremst zum Ende immer mehr ab.
func AnimationEaseOut(t float64) float64 {
	return t * (2 - t)
}

// Anfang und Ende der Animation werden abgebremst.
func AnimationEaseInOut(t float64) float64 {
	if t <= 0.5 {
		return t * t * 2
	}
	return -1 + (4-t*2)*t
}

// Alternative Funktion mit einer kubischen Funktion:
// f(x) = 3x^2 - 2x^3
func AnimationEaseInOutNew(t float64) float64 {
	return (3 - 2*t)*t*t
}

// Fuer Animationen, die endlos wiederholt weren sollen, kann diese Konstante
// fuer die Anzahl Wiederholungen verwendet werden.
const (
	AnimationRepeatForever = -1
)

// Das Interface fuer jede Art von Animation (bis jetzt zumindest).
type Animation interface {
	Start()
	Stop()
	Continue()
	IsStopped() bool
	Update(t time.Time) bool
}

// Jede konkrete Animation (Farben, Positionen, Groessen, etc.) muss die
// Methode [Tick] implementieren, welche die Veraenderungen am Objekt
// durchfuehrt. [t] ist ein Wert im Intervall [0,1], welcher den Verlauf der
// Animation angibt.
type AnimationImpl interface {
    Tick(t float64)
}

// Embeddable mit in allen Animationen benoetigen Variablen und Methoden.
// Erleichert das Erstellen von neuen Animationen gewaltig.
type AnimationEmbed struct {
	// Falls true, wird die Animation einmal vorwaerts und einmal rueckwerts
    // abgespielt.
	AutoReverse bool
	Curve       AnimationCurve
	Duration    time.Duration
	RepeatCount int

	reverse          bool
	start, stop, end time.Time
	total            float64
	repeatsLeft      int
	running          bool
	tick             func(t float64)
}

func (a *AnimationEmbed) Init(ai AnimationImpl, d time.Duration) {
	a.AutoReverse = false
	a.Curve = AnimationEaseInOutNew
    a.Duration = d
	a.RepeatCount = 0
    a.tick = ai.Tick
}

// Startet die Animation mit jenen Parametern, die zum Startzeitpunkt
// aktuell sind. Ist die Animaton bereits am Laufen ist diese Methode
// ein no-op.
func (a *AnimationEmbed) Start() {
	if a.running {
		return
	}
	a.start = time.Now()
	a.end = a.start.Add(a.Duration)
	a.total = a.end.Sub(a.start).Seconds()
	a.repeatsLeft = a.RepeatCount
	a.reverse = false
	a.running = true
	AnimCtrl.AddAnim(a)
}

// Haelt die Animation an, laesst sie jedoch in der Animation-Queue der
// Applikation. Mit [Continue] kann eine gestoppte Animation wieder
// fortgesetzt werden.
func (a *AnimationEmbed) Stop() {
	if !a.running {
		return
	}
	a.stop = time.Now()
	a.running = false
}

// Setzt eine mit [Stop] angehaltene Animation wieder fort.
func (a *AnimationEmbed) Continue() {
	if a.running {
		return
	}
	dt := time.Now().Sub(a.stop)
	a.start = a.start.Add(dt)
	a.end = a.end.Add(dt)
	a.running = true
}

// Liefert true, falls die Animation mittels [Stop] angehalten wurde oder
// falls die Animation zu Ende ist.
func (a *AnimationEmbed) IsStopped() bool {
	return !a.running
}

// Diese Methode ist fuer die korrekte Abwicklung (Beachtung von Reverse und
// RepeatCount, etc) einer Animation zustaendig. Wenn die Animation zu Ende
// ist, retourniert Update false. Der Parameter [t] ist ein "Point in Time".
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
		a.start = a.end
		a.end = a.start.Add(a.Duration)
		return true
	}

	delta := t.Sub(a.start).Seconds()
	val := delta / a.total
	if a.reverse {
		a.tick(a.Curve(1.0 - val))
	} else {
		a.tick(a.Curve(val))
	}
	return true
}


// Animation fuer einen Verlauf zwischen zwei Farben.
type ColorAnimation struct {
	AnimationEmbed
	cp     *color.Color
	c1, c2 color.Color
}

func NewColorAnimation(cp *color.Color, c2 color.Color, d time.Duration) *ColorAnimation {
	a := &ColorAnimation{}
	a.Init(a, d)
	a.cp = cp
	a.c1, a.c2 = *cp, c2
	return a
}

func (a *ColorAnimation) Tick(t float64) {
	*a.cp = a.c1.Interpolate(a.c2, t)
}

// Animation fuer einen Farbverlauf ueber die Farben einer Palette.
type PaletteAnimation struct {
    AnimationEmbed
    cp *color.Color
    pal ledgrid.ColorSource
}

func NewPaletteAnimation(cp *color.Color, pal ledgrid.ColorSource, d time.Duration) *PaletteAnimation {
    a := &PaletteAnimation{}
    a.Init(a, d)
    a.Curve = AnimationLinear
    a.cp = cp
    a.pal = pal
    return a
}

func (a *PaletteAnimation) Tick(t float64) {
    *a.cp = a.pal.Color(t)
}

// Animation fuer einen Wechsel der Position, resp. Veraenderung der Groesse.
type PositionAnimation struct {
	AnimationEmbed
	pp     *geom.Point
	p1, p2 geom.Point
}

func NewPositionAnimation(pp *geom.Point, p2 geom.Point, d time.Duration) *PositionAnimation {
	a := &PositionAnimation{}
	a.Init(a, d)
	a.pp = pp
	a.p1, a.p2 = *pp, p2
	return a
}

func (a *PositionAnimation) Tick(t float64) {
	p := a.p1.Interpolate(a.p2, t)
	*a.pp = p
}

// Da Positionen und Groessen mit dem gleichen Objekt aus geom realisiert
// werden (geom.Point), ist die Animation einer Groesse und einer Position
// im Wesentlichen das gleiche. Die Funktion NewSizeAnimation ist als
// syntaktische Vereinfachung zu verstehen.
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
	a.Init(a, d)
	a.pp = pp
	a.refPt = *pp
	a.size = size
	a.fnc = fnc
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
