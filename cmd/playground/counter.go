package main

import (
	"image"

	"github.com/stefan-muehlebach/ledgrid"
)

//----------------------------------------------------------------------------

type Counter struct {
	size  image.Point
	bits  []bool
	color ledgrid.LedColor
}

func NewCounter(size image.Point, color ledgrid.LedColor) *Counter {
	c := &Counter{}
	c.size = size
	c.bits = make([]bool, c.size.X*c.size.Y)
	c.color = color
	return c
}

func (c *Counter) Update(t float64) {
	for i, b := range c.bits {
		if !b {
			c.bits[i] = true
			break
		} else {
			c.bits[i] = false
		}
	}
}

func (c *Counter) Draw(grid *ledgrid.LedGrid) {
	for i, b := range c.bits {
		if !b {
			continue
		}
		row := i / c.size.X
		col := i % c.size.X
		grid.SetLedColor(col, row, c.color)
	}
}
