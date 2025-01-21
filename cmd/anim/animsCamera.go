package main

import (
	"context"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

func init() {
	programList.Add("Ordinary camera", OrdinaryCamera)
	programList.Add("Differential camera", DiffCamera)
}

func OrdinaryCamera(ctx context.Context, canv *ledgrid.Canvas) {
	pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	size := geom.Point{float64(width), float64(height)}

	cam := NewCamera(pos, size)
	canv.Add(cam)
	cam.Start()
}

func DiffCamera(ctx context.Context, c *ledgrid.Canvas) {
	pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	size := geom.Point{float64(width), float64(height)}

	cam := NewHistCamera(pos, size, 100, color.SkyBlue)
	c.Add(cam)
	cam.Start()
}
