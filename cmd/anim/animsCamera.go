package main

import (
	"image"
	gocolor "image/color"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

func init() {
	programList.Add("Ordinary camera", OrdinaryCamera)
	programList.Add("Differential camera", DiffCamera)
}

func OrdinaryCamera(c1 *ledgrid.Canvas) {
	c2, _ := ledGrid.NewCanvas()
	pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	size := geom.Point{float64(width), float64(height)}

	cam := NewCamera(pos, size)
	c2.Add(cam)

	mask := image.NewAlpha(c2.Rect)
	for y := range c2.Rect.Dy() {
		for x := range c2.Rect.Dx() / 2 {
			mask.Set(x, y, gocolor.Alpha{0xff})
		}
	}
	c2.Mask = mask
	cam.Start()
}

func DiffCamera(c *ledgrid.Canvas) {
	pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	size := geom.Point{float64(width), float64(height)}

	cam := NewHistCamera(pos, size, 100, color.SkyBlue)
	c.Add(cam)
	cam.Start()
}
