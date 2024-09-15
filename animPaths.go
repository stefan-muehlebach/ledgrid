package ledgrid

import (
	"math"

	"github.com/stefan-muehlebach/gg/geom"
)

var (
    CirclePath = NewGeomPath(circleFunc, 0, 1)
    LinearPath = NewGeomPath(linearFunc, 0, 1)
    RectanglePath = NewGeomPath(rectangleFunc, 0, 1)
)

type Path interface {
    Pos(t float64) geom.Point
}

type GeomPath struct {
    t0, l float64
    p0 geom.Point
    fnc PathFunctionType
}

func NewGeomPath(fnc PathFunctionType, t0, l float64) *GeomPath {
    p := &GeomPath{t0: t0, l: l, p0: fnc(t0), fnc: fnc}
    return p
}

func (p *GeomPath) Pos(t float64) geom.Point {
    t *= p.l
    t += p.t0
    if t > 1.0 {
        t -= 1.0
    }
    return p.fnc(t).Sub(p.p0)
}

func (p *GeomPath) NewStart(t0 float64) *GeomPath {
    q := NewGeomPath(p.fnc, t0, p.l)
    return q
}

func (p *GeomPath) NewStartLen(t0, l float64) *GeomPath {
    q := NewGeomPath(p.fnc, t0, l)
    return q
}

// Im Folgenden sind einige Pfad-generierende Funktionen zusammengestellt, die
// als Parameter [pathFunc] bei NewPathAnimation verwendet werden können.
// Eigene Pfad-Funktionen sind ebenfalls möglich, die Bedingungen an eine
// solche Funktion sind:
//
//  1. Wird mit einer Fliesskommazahl (t) aufgerufen und retourniert einen
//     2D-Punkt
//  2. t ist in [0,1]
//  3. f(0) muss (0,0) sein
//  4. max(f(t).X) - min(f(t).X) = 1.0 und
//     max(f(t).Y) - min(f(t).Y) = 1.0
//     d.h. die generierten Punkte duerfen sowohl in X- als auch in Y-Richtung
//     einen maximalen Abstand von 1.0 haben.
type PathFunctionType func(t float64) geom.Point

// Beschreibt eine Gerade
func linearFunc(t float64) geom.Point {
	return geom.Point{t, t}
}

// Beschreibt eine Gerade
func circleFunc(t float64) geom.Point {
	phi := t * 2 * math.Pi
	return geom.Point{0.5*math.Sin(phi), -0.5*math.Cos(phi)}
}

// Beschreibt ein Rechteck
func rectangleFunc(t float64) geom.Point {
	switch {
	case t < 1.0/4.0:
		return geom.Point{4.0 * t, 0.0}
	case t < 2.0/4.0:
		return geom.Point{1.0, 4.0 * (t - 1.0/4.0)}
	case t < 3.0/4.0:
		return geom.Point{1.0 - 4.0 * (t - 2.0/4.0), 1.0}
	default:
		return geom.Point{0.0, 1.0 - 4.0 * (t - 3.0/4.0)}
	}
}
