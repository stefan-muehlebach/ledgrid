package ledgrid

import (
	"image"
	"image/color"
	"image/draw"
)

// Erster Versuch eines Filter-Objektes, das vor Canvas geschaltet werden
// kann. Wird beim Uebertragen der Pixel in das LedGrid-Objekt angewandt.
// type ImageFilter struct {
// 	Img draw.Image
// 	Flt Filter
// }

type Filter interface {
    ColorModel() color.Model
    Bounds() image.Rectangle
    At(x, y int) color.Color
    Set(x, y int, c color.Color)
	FF(x, y int) (int, int)
}

type FilterImpl interface {
	FF(x, y int) (int, int)
}

type FilterBase struct {
	img draw.Image
	flt FilterImpl
}

func (f *FilterBase) Extend(img draw.Image, flt FilterImpl) {
	f.img = img
    f.flt = flt
}

func (f *FilterBase) ColorModel() color.Model {
	return f.img.ColorModel()
}

func (f *FilterBase) Bounds() image.Rectangle {
	return f.img.Bounds()
}

func (f *FilterBase) At(x, y int) color.Color {
	x, y = f.flt.FF(x, y)
	return f.img.At(x, y)
}

func (f *FilterBase) Set(x, y int, col color.Color) {
	x, y = f.flt.FF(x, y)
	f.img.Set(x, y, col)
}

type FilterIdent struct{
    FilterBase
}

func NewFilterIdent(img draw.Image) *FilterIdent {
    f := &FilterIdent{}
    f.FilterBase.Extend(img, f)
    return f
}

func (f *FilterIdent) FF(x, y int) (int, int) {
	return x, y
}
