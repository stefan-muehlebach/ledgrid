package main

import (
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/geom"
)

//----------------------------------------------------------------------------

type AnimationCurve func(t float64) float64

func LinearAnimationCurve(t float64) float64 {
	return t
}

func CubicAnimationCurve(t float64) float64 {
	return 3*t*t - 2*t*t*t
}

const (
	AnimationRepeatForever = -1
)

// Interface fuer jede Animation.
type Animation interface {
	Start()
	Stop()
	Continue()
	IsStopped() bool
	Update(t time.Time) bool
}

// Embeddable mit den fuer alle Animationen verbindlichen Methoden.
type AnimationEmbed struct {
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
    a.Curve = CubicAnimationCurve
    a.Duration = d
    a.RepeatCount = 0
}

func (a *AnimationEmbed) Start() {
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
	a.tick = a.UpdateColor
	return a
}

func (a *ColorAnimation) UpdateColor(t float64) {
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
	a.tick = a.UpdatePosition
	return a
}

func (a *PositionAnimation) UpdatePosition(t float64) {
	p := a.p1.Interpolate(a.p2, t)
	*a.pp = p
}

var (
    NewSizeAnimation = NewPositionAnimation
)

// Animation fuer das Fahren entlang eines Pfades.
type PathAnimation struct {
	AnimationEmbed
	pp     *geom.Point
	p1 geom.Point
    size geom.Point
    fnc PathFunction
}

type PathFunction func(t float64) geom.Point

func NewPathAnimation(pp *geom.Point, fnc PathFunction, size geom.Point, d time.Duration) *PathAnimation {
	a := &PathAnimation{}
    a.Init(d)
	a.pp = pp
	a.p1 = *pp
    a.size = size
    a.fnc = fnc
	a.tick = a.UpdatePosition
	return a
}

func (a *PathAnimation) UpdatePosition(t float64) {
    pos := a.fnc(t)
    dp := geom.Point{pos.X*a.size.X, pos.Y*a.size.Y}
	*a.pp = a.p1.Add(dp)
}

func circlePathFunc(t float64) geom.Point {
    phi := t * 2 * math.Pi
    return geom.Point{0.5*(1-math.Cos(phi)), 0.5*math.Sin(phi)}
}

// Allgemeine Parameter-Animation.
type FloatAnimation struct {
	AnimationEmbed
	fp     *float64
	f1, f2 float64
}

func NewFloatAnimation(fp *float64, f2 float64, d time.Duration) *FloatAnimation {
	a := &FloatAnimation{}
    a.Init(d)
	a.fp = fp
	a.f1, a.f2 = *fp, f2
	a.tick = a.UpdateFloat
	return a
}

func (a *FloatAnimation) UpdateFloat(t float64) {
	f := (1.0-t)*a.f1 + t*a.f2
	*a.fp = f
}
