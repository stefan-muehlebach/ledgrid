package ledgrid

import (
	"image"
	"image/color"
	"time"
)

// Transitions can be used to blend from one canvas to another. They implement
// image.Image and can be used as a value for the field Mask of every Canvas.
// This masks are used when LedGrid combines the contents of all canvases
// to one single image: they define what part and how much of each canvas is
// drawn to the final image. Beside that, they implement the NormAnimation
// interface... they can be seen as animated masks.

// This is the simplest form of a transition: it will mask the whole canvas
// with the same alpha value. With it you can realize smooth fadings from one
// canvas to another.
type FadeTransition struct {
	NormAnimationEmbed
	c             *Canvas
	a0, a1, alpha uint8
}

func NewFadeTransition(c *Canvas, a0, a1 uint8, d time.Duration) *FadeTransition {
	a := &FadeTransition{}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(d)
	a.c = c
	a.a0, a.a1 = a0, a1
	return a
}

func (a *FadeTransition) ColorModel() color.Model {
	return color.AlphaModel
}

func (a *FadeTransition) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{-1e9, -1e9}, image.Point{1e9, 1e9}}
}

func (a *FadeTransition) At(x, y int) color.Color {
	return color.Alpha{a.alpha}
}

func (a *FadeTransition) Init() {
	a.alpha = a.a0
	a.c.Mask = a
}

func (a *FadeTransition) Tick(t float64) {
	a.alpha = uint8((1-t)*float64(a.a0) + t*float64(a.a1))
}

// Wipe transitions offer an effect which is very common in film an tv. It
// gives the illusion of a cover which is pulled away in any of the four
// possible ways.
type WipeDirection int

const (
	WipeL2R WipeDirection = iota
	WipeR2L
	WipeT2B
	WipeB2T
)

type WipeType int

const (
	WipeIn WipeType = iota
	WipeOut
)

type WipeTransition struct {
	NormAnimationEmbed
	c                                      *Canvas
	rect                                   image.Rectangle
	dir                                    WipeDirection
	typ                                    WipeType
	edgeWidth, edgeStart, edgeEnd, edgePos float64
	lb, ub                                 int
	alphaStep                              uint8
}

func NewWipeTransition(c *Canvas, dir WipeDirection, typ WipeType, d time.Duration) *WipeTransition {
	a := &WipeTransition{}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(d)
	a.c = c
	a.rect = c.Bounds()
	a.dir = dir
	a.typ = typ
	a.edgeWidth = 3.0
	switch a.dir {
	case WipeL2R:
		a.edgeStart = -(a.edgeWidth / 2.0) - 1.0
		a.edgeEnd = float64(a.rect.Dx()) + (a.edgeWidth / 2.0)
	case WipeR2L:
		a.edgeStart = float64(a.rect.Dx()) + (a.edgeWidth / 2.0)
		a.edgeEnd = -(a.edgeWidth / 2.0) - 1.0
	case WipeT2B:
		a.edgeStart = -(a.edgeWidth / 2.0) - 1.0
		a.edgeEnd = float64(a.rect.Dy()) + (a.edgeWidth / 2.0)
	case WipeB2T:
		a.edgeStart = float64(a.rect.Dy()) + (a.edgeWidth / 2.0)
		a.edgeEnd = -(a.edgeWidth / 2.0) - 1.0
	}
	a.alphaStep = uint8(255.0 / a.edgeWidth)
	return a
}

func (a *WipeTransition) ColorModel() color.Model {
	return color.AlphaModel
}

func (a *WipeTransition) Bounds() image.Rectangle {
	return a.rect
}

func (a *WipeTransition) At(x, y int) color.Color {
	t := x
	if a.dir == WipeT2B || a.dir == WipeB2T {
		t = y
	}
	if a.typ == WipeIn {
		if t < a.lb {
			return color.Alpha{0xff}
		}
		if t > a.ub {
			return color.Alpha{0x00}
		}
		return color.Alpha{0xff - uint8(t-a.lb)*a.alphaStep}
	} else {
		if t < a.lb {
			return color.Alpha{0x00}
		}
		if t > a.ub {
			return color.Alpha{0xff}
		}

		return color.Alpha{uint8(t-a.lb) * a.alphaStep}
	}
}

func (a *WipeTransition) Init() {
	a.edgePos = a.edgeStart
	a.c.Mask = a
}

func (a *WipeTransition) Tick(t float64) {
	a.edgePos = (1-t)*(a.edgeStart) + (t)*(a.edgeEnd)
	a.lb = int(a.edgePos - (a.edgeWidth / 2.0))
	a.ub = int(a.edgePos + (a.edgeWidth / 2.0))
}

// type TunnelTransition struct {
//     NormAnimationEmbed
//     rect image.Rectangle
//     mp geom.Point
//     edgeWidth, edgeStart, edgeEnd, edgePos float64
//     lb, ub int
//     alphaStep uint8
// }
