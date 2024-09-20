//go:build !cameraV4L2 && !cameraOpenCV

package main

import (
	"image"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
)

type Camera struct {
	ledgrid.CanvasObjectEmbed
	Pos, Size geom.Point
	Mask   *image.Alpha
	running   bool
}

func NewCamera(pos, size geom.Point) *Camera {
	c := &Camera{Pos: pos, Size: size}
	c.CanvasObjectEmbed.Extend(c)
	c.Mask = image.NewAlpha(image.Rectangle{Max: size.Int()})
	ledgrid.AnimCtrl.Add(c)
	return c
}

func (c *Camera) Duration() time.Duration {
	return time.Duration(0)
}

func (c *Camera) SetDuration(dur time.Duration) {}

func (c *Camera) Start() {
	if c.running {
		return
	}
	// Would do starting things here.
	c.running = true
}

func (c *Camera) Stop() {
	if !c.running {
		return
	}
	// Would do the stopping things here.
	c.running = false
}

func (c *Camera) Suspend() {}

func (c *Camera) Continue() {}

func (c *Camera) IsRunning() bool {
	return c.running
}

func (c *Camera) Update(pit time.Time) bool {
	// Update whatever you might have to update
	// Remember: this method gets called every approx. 30 ms!
	// Do only what you absolutely have to do here!
	return true
}

func (c *Camera) Draw(canv *ledgrid.Canvas) {
	// Copy or Blend the image from the camera onto the canvas.
}
