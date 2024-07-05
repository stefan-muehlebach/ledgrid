package main

import (
	"image"
	"time"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/ledgrid"
	"golang.org/x/image/draw"
)

//----------------------------------------------------------------------------

type Controller struct {
	pixCtrl  ledgrid.PixelClient
	ledGrid  *ledgrid.LedGrid
	canvas   *image.RGBA
	gc       *gg.Context
	objList  []CanvasObject
	scaler   draw.Scaler
	ticker   *time.Ticker
	animList []Animation
}

func NewController(pixCtrl ledgrid.PixelClient, ledGrid *ledgrid.LedGrid) *Controller {
	if AnimCtrl != nil {
		return AnimCtrl
	}
	c := &Controller{}
	c.pixCtrl = pixCtrl
	c.ledGrid = ledGrid
	c.canvas = image.NewRGBA(image.Rectangle{Max: c.ledGrid.Rect.Max.Mul(int(oversize))})
	c.gc = gg.NewContextForRGBA(c.canvas)
	c.objList = make([]CanvasObject, 0)
	c.scaler = draw.BiLinear.NewScaler(c.ledGrid.Rect.Dx(), c.ledGrid.Rect.Dy(), c.canvas.Rect.Dx(), c.canvas.Rect.Dy())
	c.ticker = time.NewTicker(refreshRate)
	go c.backgroundThread()
	AnimCtrl = c
	return c
}

func (c *Controller) Add(objs ...CanvasObject) {
	c.objList = append(c.objList, objs...)
}

func (c *Controller) AddAnim(anims ...Animation) {
	c.animList = append(c.animList, anims...)
}

func (c *Controller) Stop() {
	c.ticker.Stop()
}

func (c *Controller) Continue() {
	c.ticker.Reset(refreshRate)
}

func (c *Controller) backgroundThread() {
	backColor := color.Black.Alpha(backAlpha)

	for pit := range c.ticker.C {
		for i, anim := range c.animList {
			if anim == nil {
				continue
			}
			if anim.IsStopped() {
				continue
			}
			if !anim.Update(pit) {
				c.animList[i] = nil
			}
		}
		c.gc.SetFillColor(backColor)
		c.gc.Clear()
		for _, obj := range c.objList {
			obj.Draw(c.gc)
		}
		c.scaler.Scale(c.ledGrid, c.ledGrid.Rect, c.canvas, c.canvas.Rect, draw.Over, nil)
		c.pixCtrl.Draw(c.ledGrid)
	}
}
