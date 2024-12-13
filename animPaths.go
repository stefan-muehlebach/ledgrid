package ledgrid

import (
	"math"

	"github.com/stefan-muehlebach/gg/geom"
)

var (
    CirclePath = NewGeomPath(circlePathFunc, 0, 1)
    LinearPath = NewGeomPath(linearPathFunc, 0, 1)
    RectanglePath = NewGeomPath(rectPathFunc, 0, 1)
)

// Ein Pfad ist im Grunde nichts anderes, als eine Funktion der Form
//   [0,1] -> (x,y)
type Path interface {
    Pos(t float64) geom.Point
}

// Geometrische Pfade heissen deshalb so, weil ihre Grundlage zum Punkte
// bauen irgendwo in der Geometrie zu suchen ist.
type GeomPath struct {
    t0, l float64
    p0 geom.Point
    fnc PathFunctionType
}

// Alle geom. Pfade beruhen auf einer Pfad-Funktion (erster Parameter), der
// beim Erstellen eines Pfades zwingend angegeben werden muss. Mit t0 und l
// kann der Startpunkt, bzw. die Laenge der Animation im Vergleich zur
// gesamten Laenge von fnc bestimmt werden.
func NewGeomPath(fnc PathFunctionType, t0, l float64) *GeomPath {
    p := &GeomPath{t0: t0, l: l, p0: fnc(t0), fnc: fnc}
    return p
}

// Damit implementiert GeomPath das Interface Path... was die Absicht unseres
// ganzen Planes ist.
func (p *GeomPath) Pos(t float64) geom.Point {
    t *= p.l
    t += p.t0
    if t > 1.0 {
        t -= 1.0
    }
    return p.fnc(t).Sub(p.p0)
}

// Baut aus dem geom. Pfad p einen neuen, indem der Anfangspunkt auf t0
// verlegt wird.
func (p *GeomPath) NewStart(t0 float64) *GeomPath {
    q := NewGeomPath(p.fnc, t0, p.l)
    return q
}

// Baut aus dem geom. Pfad p einen neuen, indem der Anfangspunkt auf t0
// und die totale Laenge des Pfades auf l gesetzt wird.
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

// Beschreibt eine Gerade vom Punkt (0,0) zum Punkt (1,1)
func linearPathFunc(t float64) geom.Point {
	return geom.Point{t, t}
}

// Beschreibt einen Kreis, der eine Breite/Hoehe von 1.0 hat, am Ursprung
// zentriert ist, oben in der Mitte beginnt und dann im Uhrzeigersinn
// verlaeuft.
func circlePathFunc(t float64) geom.Point {
	phi := t * 2 * math.Pi
	return geom.Point{0.5*math.Sin(phi), -0.5*math.Cos(phi)}
}

// Beschreibt ein Rechteck. Start ist links oben und der Verlauf ist im
// Uhrzeigersinn.
func rectPathFunc(t float64) geom.Point {
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

// Neben den vorhandenen Pfaden (Kreise, Halbkreise, Viertelkreise) koennen
// Positions-Animationen auch entlang komplett frei definierten Pfaden
// erfolgen. Der Schluessel dazu ist der Typ PolygonPath.
type PolygonPath struct {
	rect     geom.Rectangle
	stopList []polygonStop
}

type polygonStop struct {
	dist float64
	pos geom.Point
}

// Erstellt ein neues PolygonPath-Objekt und verwendet die Punkte in points
// als Eckpunkte eines offenen Polygons.
func NewPolygonPath(points ...geom.Point) *PolygonPath {
	p := &PolygonPath{}
	p.stopList = make([]polygonStop, len(points))

	origin := geom.Point{}
	for i, point := range points {
		if i == 0 {
			origin = point
			p.stopList[i] = polygonStop{0.0, geom.Point{0, 0}}
			continue
		}
		pt := point.Sub(origin)
		dist := p.stopList[i-1].dist + pt.Distance(p.stopList[i-1].pos)
		p.stopList[i] = polygonStop{dist, pt}

		p.rect.Min = p.rect.Min.Min(pt)
		p.rect.Max = p.rect.Max.Max(pt)
	}
	return p
}

// Diese Methode ist bei der Erstellung einer Pfad-Animation als Parameter
// fnc anzugeben.
func (p *PolygonPath) Pos(t float64) geom.Point {
	dist := t * p.stopList[len(p.stopList)-1].dist
	for i, stop := range p.stopList[1:] {
		if dist < stop.dist {
			p1 := p.stopList[i].pos
			p2 := stop.pos
			relDist := dist - p.stopList[i].dist
			f := relDist / (stop.dist - p.stopList[i].dist)
			return p1.Interpolate(p2, f)
		}
	}
	return p.stopList[len(p.stopList)-1].pos
}
