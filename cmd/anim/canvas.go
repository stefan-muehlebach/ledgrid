package main

import (
	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/geom"
)

//----------------------------------------------------------------------------

type CanvasObject interface {
	Draw(gc *gg.Context)
}

//----------------------------------------------------------------------------

type Ellipse struct {
	Pos, Size              geom.Point
	BorderWidth    float64
	FillColor, BorderColor color.Color
}

func (e *Ellipse) Draw(gc *gg.Context) {
	gc.DrawEllipse(e.Pos.X, e.Pos.Y, e.Size.X, e.Size.Y)
	gc.SetStrokeWidth(e.BorderWidth)
	gc.SetStrokeColor(e.BorderColor)
	gc.SetFillColor(e.FillColor)
	gc.FillStroke()
}

//----------------------------------------------------------------------------

type Rectangle struct {
	Pos                    geom.Point
	Size                   geom.Point
	BorderWidth            float64
	FillColor, BorderColor color.Color
}

func (r *Rectangle) Draw(gc *gg.Context) {
	gc.DrawRectangle(r.Pos.X-r.Size.X/2, r.Pos.Y-r.Size.Y/2, r.Size.X, r.Size.Y)
	gc.SetStrokeWidth(r.BorderWidth)
	gc.SetStrokeColor(r.BorderColor)
	gc.SetFillColor(r.FillColor)
	gc.FillStroke()
}

//----------------------------------------------------------------------------

type Line struct {
	Pos1, Pos2 geom.Point
	Width      float64
	Color      color.Color
}

func (l *Line) Draw(gc *gg.Context) {
	gc.SetStrokeWidth(l.Width)
	gc.SetStrokeColor(l.Color)
	gc.DrawLine(l.Pos1.X, l.Pos1.Y, l.Pos2.X, l.Pos2.Y)
	gc.Stroke()
}
