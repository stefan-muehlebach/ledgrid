//go:build ignore
// +build ignore

package main

import (
	"math"

	"github.com/stefan-muehlebach/ledgrid"
	"golang.org/x/image/math/fixed"
)

func ipart(x float64) float64 {
	return math.Floor(x)
}

func round(x float64) float64 {
	return ipart(x + .5)
}

func fpart(x float64) float64 {
	return x - ipart(x)
}

func rfpart(x float64) float64 {
	return 1 - fpart(x)
}

// AaLine plots anti-aliased line by Xiaolin Wu's line algorithm.
func DrawLine(g *ledgrid.LedGrid, p0, p1 fixed.Point26_6, col ledgrid.LedColor) {
	var x0, y0, x1, y1 float64
	var bgCol ledgrid.LedColor

	x0, y0 = fix2float(p0.X), fix2float(p0.Y)
	x1, y1 = fix2float(p1.X), fix2float(p1.Y)

	// straight translation of WP pseudocode
	dx := x1 - x0
	dy := y1 - y0
	ax := dx
	if ax < 0 {
		ax = -ax
	}
	ay := dy
	if ay < 0 {
		ay = -ay
	}
	// plot function set here to handle the two cases of slope
	var plot func(int, int, float64)

	if ax < ay {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		dx, dy = dy, dx
		plot = func(x, y int, c float64) {
			bgCol = g.LedColorAt(y, x)
			g.SetLedColor(y, x, col.Alpha(c).Mix(bgCol, ledgrid.Blend)) //uint16(c*math.MaxUint16))
		}
	} else {
		plot = func(x, y int, c float64) {
			bgCol = g.LedColorAt(x, y)
			g.SetLedColor(x, y, col.Alpha(c).Mix(bgCol, ledgrid.Blend)) //uint16(c*math.MaxUint16))
		}
	}
	if x1 < x0 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	gradient := dy / dx

	// handle first endpoint
	xend := round(x0)
	yend := y0 + gradient*(xend-x0)
	xgap := rfpart(x0 + 0.5)
	xpxl1 := int(xend) // this will be used in the main loop
	ypxl1 := int(ipart(yend))
	plot(xpxl1, ypxl1, rfpart(yend)*xgap)
	plot(xpxl1, ypxl1+1, fpart(yend)*xgap)
	intery := yend + gradient // first y-intersection for the main loop

	// handle second endpoint
	xend = round(x1)
	yend = y1 + gradient*(xend-x1)
	xgap = fpart(x1 + 0.5)
	xpxl2 := int(xend) // this will be used in the main loop
	ypxl2 := int(ipart(yend))
	plot(xpxl2, ypxl2, rfpart(yend)*xgap)
	plot(xpxl2, ypxl2+1, fpart(yend)*xgap)

	// main loop
	for x := xpxl1 + 1; x <= xpxl2-1; x++ {
		plot(x, int(ipart(intery)), rfpart(intery))
		plot(x, int(ipart(intery))+1, fpart(intery))
		intery = intery + gradient
	}
}
