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
	a0, a1, alpha uint8
}

func NewFadeTransition(a0, a1 uint8, d time.Duration) *FadeTransition {
	a := &FadeTransition{}
	a.NormAnimationEmbed.Extend(a)
	a.SetDuration(d)
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
}

func (a *FadeTransition) Tick(t float64) {
	a.alpha = uint8((1-t)*float64(a.a0) + t*float64(a.a1))
}



type WipeTransition struct {
    NormAnimationEmbed
    rect image.Rectangle
    edgeWidth, edgeStart, edgeEnd, edgePos float64
    lb, ub int
    alphaStep uint8
}

func NewWipeTransition(r image.Rectangle, d time.Duration) *WipeTransition {
    a := &WipeTransition{}
    a.NormAnimationEmbed.Extend(a)
    a.SetDuration(d)
    a.rect = r
    a.edgeWidth = 4.0
    a.edgeStart = -(a.edgeWidth/2.0)
    a.edgeEnd = float64(r.Dx()) + (a.edgeWidth/2.0)
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
    if x < a.lb {
        return color.Alpha{0xff}
    }
    if x > a.ub {
        return color.Alpha{0x00}
    }
    return color.Alpha{0xff - uint8(x-a.lb)*a.alphaStep}
}

func (a *WipeTransition) Init() {
	a.edgePos = a.edgeStart
}

func (a *WipeTransition) Tick(t float64) {
	a.edgePos = (1-t)*(a.edgeStart) + (t)*(a.edgeEnd)
    a.lb = int(a.edgePos - (a.edgeWidth/2.0))
    a.ub = int(a.edgePos + (a.edgeWidth/2.0))
}
