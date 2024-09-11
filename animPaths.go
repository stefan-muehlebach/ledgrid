package ledgrid

import (
	"math"

	"github.com/stefan-muehlebach/gg/geom"
)

// Im Folgenden sind einige Pfad-generierende Funktionen zusammengestellt, die
// als Parameter [pathFunc] bei NewPathAnimation verwendet werden können.
// Eigene Pfad-Funktionen sind ebenfalls möglich, die Bedingungen an eine
// solche Funktion sind beim Funktionstyp [PathFunctionType] beschrieben.

// Die PathFunctionType muss folgende Bedingungen erfuellen:
//  1. Wird mit einer Fliesskommazahl (t) aufgerufen und retourniert einen
//     2D-Punkt
//  2. t ist in [0,1]
//  3. f(0) muss (0,0) sein
//  4. max(f(t).X) - min(f(t).X) = 1.0 und
//     max(f(t).Y) - min(f(t).Y) = 1.0
type PathFunctionType func(t float64) geom.Point

// Beschreibt eine Gerade
func LinearPath(t float64) geom.Point {
	return geom.Point{t, t}
}

// Beschreibt ein Rechteck im Uhrzeigersinn.
// Startpunkt ist auf 12 Uhr.
func RectanglePathA(t float64) geom.Point {
	switch {
	case t < 1.0/8.0:
		return geom.Point{0.5 * 8.0 * t, 0.0}
	case t < 3.0/8.0:
		return geom.Point{0.5, 4.0 * (t - 1.0/8.0)}
	case t < 5.0/8.0:
		return geom.Point{0.5 - 4.0*(t-3.0/8.0), 1.0}
	case t < 7.0/8.0:
		return geom.Point{-0.5, 1.0 - 4.0*(t-5.0/8.0)}
	default:
		return geom.Point{-0.5 + 0.5*8.0*(t-7.0/8.0), 0.0}
	}
}

// Beschreibt ein Rechteck im Uhrzeigersinn.
// Startpunkt ist auf 9 Uhr.
func RectanglePathB(t float64) geom.Point {
	switch {
	case t < 1.0/8.0:
		return geom.Point{0.0, -0.5 * 8.0 * t}
	case t < 3.0/8.0:
		return geom.Point{4.0 * (t - 1.0/8.0), -0.5}
	case t < 5.0/8.0:
		return geom.Point{1.0, 4.0*(t-3.0/8.0) - 0.5}
	case t < 7.0/8.0:
		return geom.Point{1.0 - 4.0*(t-5.0/8.0), 0.5}
	default:
		return geom.Point{0, 0.5 * (1.0 - 8.0*(t-7.0/8.0))}
	}
}

// Beschreibt einen Kreis oder Ellipse im Uhrzeigersinn.
// Startpunkt ist auf 12 Uhr.
func FullCirclePathA(t float64) geom.Point {
	phi := t * 2 * math.Pi
	return geom.Point{0.5 * math.Sin(phi), 0.5 - 0.5*math.Cos(phi)}
}

// Startpunkt ist auf 9 Uhr.
func FullCirclePathB(t float64) geom.Point {
	phi := t * 2 * math.Pi
	return geom.Point{0.5 - 0.5*math.Cos(phi), -(0.5 * math.Sin(phi))}
}

// Beschreibt einen Halbkreis.
func HalfCirclePathA(t float64) geom.Point {
	phi := t * math.Pi
	return geom.Point{math.Sin(phi), (1.0 - math.Cos(phi)) / 2.0}
}

func HalfCirclePathB(t float64) geom.Point {
	phi := t * math.Pi
	return geom.Point{(1.0 - math.Cos(phi)) / 2.0, math.Sin(phi)}
}

// Beschreibt einen Viertelkreis.
// Horizontaler Start.
func QuarterCirclePathA(t float64) geom.Point {
	phi := t * math.Pi / 2.0
	return geom.Point{math.Sin(phi), 1.0 - math.Cos(phi)}
}

// Vertikaler Start.
func QuarterCirclePathB(t float64) geom.Point {
	phi := t * math.Pi / 2.0
	return geom.Point{1.0 - math.Cos(phi), math.Sin(phi)}
}
