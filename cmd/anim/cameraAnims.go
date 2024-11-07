package main

import (
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

var (
	OrdinaryCamera = NewLedGridProgram("Ordinary camera",
		func(c *ledgrid.Canvas) {
			pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			size := geom.Point{float64(width), float64(height)}

			cam := NewCamera(pos, size)
			c.Add(cam)
			cam.Start()
		})

	SpecialCamera = NewLedGridProgram("Camera in differential mode",
		func(c *ledgrid.Canvas) {
			pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			size := geom.Point{float64(width), float64(height)}

			cam := NewHistCamera(pos, size, 100, color.SkyBlue)
			c.Add(cam)
			cam.Start()
		})
)
