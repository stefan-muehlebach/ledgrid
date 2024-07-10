package main

import (
	"runtime"
	"image"
	"sync"
	"time"

	"github.com/stefan-muehlebach/gg"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/ledgrid"
	"golang.org/x/image/draw"
)

//----------------------------------------------------------------------------

type Controller struct {
	pixCtrl   ledgrid.PixelClient
	ledGrid   *ledgrid.LedGrid
	canvas    *image.RGBA
	gc        *gg.Context
	objList   []CanvasObject
	scaler    draw.Scaler
	ticker    *time.Ticker
	animList  []Animation
	animMutex *sync.Mutex
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
    c.animList = make([]Animation, 0)
    c.animMutex = &sync.Mutex{}
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

type animJobType struct {
    id int
    pit time.Time
}

func (c *Controller) backgroundThread() {
	backColor := color.Black.Alpha(backAlpha)
    numCores := runtime.NumCPU()
    animChan := make(chan animJobType, 2*numCores)
    doneChan := make(chan bool, 2*numCores)

    for range numCores {
        go c.updateWorker(animChan, doneChan)
    }

	for pit := range c.ticker.C {
        numAnims := 0
		for i, anim := range c.animList {
			if anim == nil || anim.IsStopped() {
				continue
			}
            numAnims++
            animChan <- animJobType{i, pit}
			// if !anim.Update(pit) {
			// 	c.animList[i] = nil
			// }
        }
        for range numAnims {
            <- doneChan
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

func (c *Controller) updateWorker(animChan <-chan animJobType, doneChan chan<- bool) {
    for animJob := range animChan {
        if !c.animList[animJob.id].Update(animJob.pit) {
            c.animList[animJob.id] = nil
        }
        doneChan <- true
    }
}
