//go:build ignore

package ledgrid

import (
	"image"
	"math"
	"time"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/geom"
	"golang.org/x/image/draw"
)

const (
	// Diese Konstante legt fest, mit welchem Faktor die Groesse des LedGrid
	// multipliziert wird, um damit ein image.RGBA-Objekt zu erstellen.
	canvasScale = 10.0
)

// Dies ist ein Sammelgefaess fuer einfache geometrische Animationen.
type Canvas struct {
	VisualEmbed
	lg      *LedGrid
	img     *image.RGBA
	gc      *gg.Context
	objList []CanvasObject
	scaler  draw.Scaler
}

func NewCanvas(lg *LedGrid) *Canvas {
	a := &Canvas{}
	a.VisualEmbed.Init("Canvas")
	a.lg = lg
	a.img = image.NewRGBA(image.Rect(0, 0, canvasScale*a.lg.Rect.Dx(),
		canvasScale*a.lg.Rect.Dy()))
	a.gc = gg.NewContextForRGBA(a.img)
	a.objList = make([]CanvasObject, 0)
	a.scaler = draw.BiLinear.NewScaler(a.lg.Bounds().Dx(), a.lg.Bounds().Dy(),
		a.img.Bounds().Dx(), a.img.Bounds().Dy())
	return a
}

func (a *Canvas) Update(dt time.Duration) bool {
	dt = a.AnimatableEmbed.Update(dt)
	for _, obj := range a.objList {
		obj.Update(a.t0, dt)
	}
	return true
}

func (a *Canvas) Draw() {
	a.gc.SetFillColor(color.Transparent)
	a.gc.Clear()
	for _, obj := range a.objList {
		obj.Draw(a.gc)
	}
	a.scaler.Scale(a.lg, a.lg.Bounds(), a.img, a.img.Bounds(), draw.Over, nil)
}

func (a *Canvas) AddObjects(objs ...CanvasObject) {
	a.objList = append(a.objList, objs...)
}

//----------------------------------------------------------------------------

type CanvasObject interface {
	Update(t time.Duration, dt time.Duration)
	Draw(gc *gg.Context)
}

//----------------------------------------------------------------------------

type RotatingLine struct {
	Pos   geom.Point
	Speed float64
	Len   float64
	Color color.Color
	Width float64
	p1    geom.Point
	Angle float64
}

func (l *RotatingLine) Update(t time.Time, dt time.Duration) {
	l.Angle += l.Speed * dt.Seconds()
	l.p1 = l.Pos.Add(geom.Point{l.Len * math.Cos(l.Angle), l.Len * math.Sin(l.Angle)})
}

func (l *RotatingLine) Draw(gc *gg.Context) {
	gc.SetStrokeColor(l.Color)
	gc.SetStrokeWidth(l.Width)
	gc.DrawLine(l.Pos.X, l.Pos.Y, l.p1.X, l.p1.Y)
	gc.Stroke()
}

//----------------------------------------------------------------------------

type GlowingCircle struct {
	Pos, Dir    geom.Point
	Speed       float64
	Radius      []float64
	FillColor   []color.Color
	StrokeColor []color.Color
	StrokeWidth []float64

	GlowPeriod time.Duration
	// pit                    time.Duration
	radius                 float64
	fillColor, strokeColor color.Color
	strokeWidth            float64
}

func (c *GlowingCircle) Update(t time.Duration, dt time.Duration) {
	c.Pos = c.Pos.Add(c.Dir.Mul(c.Speed))
	if c.Pos.X < 0.0 || c.Pos.X >= 50.0 {
		c.Dir.X = -c.Dir.X
	}
	if c.Pos.Y < 0.0 || c.Pos.Y >= 50.0 {
		c.Dir.Y = -c.Dir.Y
	}
	// c.pit += dt
	t0 := t % c.GlowPeriod
	t1 := t % (c.GlowPeriod / 2)
	if t0 > t1 {
		t0 = (c.GlowPeriod / 2) - t1
	}
	f := 2.0 * t0.Seconds() / (c.GlowPeriod.Seconds())
	c.radius = PolynomInterpol(c.Radius[0], c.Radius[1], f)
	c.fillColor = c.FillColor[0].Interpolate(c.FillColor[1], f)
	c.strokeColor = c.StrokeColor[0].Interpolate(c.StrokeColor[1], f)
	c.strokeWidth = PolynomInterpol(c.StrokeWidth[0], c.StrokeWidth[1], f)
}

func (c *GlowingCircle) Draw(gc *gg.Context) {
	gc.SetFillColor(c.fillColor)
	gc.SetStrokeColor(c.strokeColor)
	gc.SetStrokeWidth(c.strokeWidth)
	gc.DrawCircle(c.Pos.X, c.Pos.Y, c.radius)
	gc.FillStroke()
}
