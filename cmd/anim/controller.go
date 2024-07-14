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

const (
    queueSize = 200
)

// Der Controller ist das Bindeglied zwischen dem
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
    animPit   time.Time
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
	c.scaler = draw.CatmullRom.NewScaler(c.ledGrid.Rect.Dx(), c.ledGrid.Rect.Dy(), c.canvas.Rect.Dx(), c.canvas.Rect.Dy())
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

func (c *Controller) DelAll() {
    c.objList = c.objList[:0]
}

func (c *Controller) AddAnim(anims ...Animation) {
    c.animMutex.Lock()
	c.animList = append(c.animList, anims...)
    c.animMutex.Unlock()
}

func (c *Controller) DelAllAnim() {
    c.animMutex.Lock()
    c.animList = c.animList[:0]
    c.animMutex.Unlock()
}

func (c *Controller) DelAnim(anim Animation) {
    c.animMutex.Lock()
    for i, a := range c.animList {
        if a == anim {
            c.animList[i] = nil
            return
        }
    }
    c.animMutex.Unlock()
}

func (c *Controller) Stop() {
	c.ticker.Stop()
}

func (c *Controller) Continue() {
	c.ticker.Reset(refreshRate)
}

func (c *Controller) backgroundThread() {
	backColor := color.Black.Alpha(backAlpha)
    numCores := runtime.NumCPU()
    animChan := make(chan int, queueSize)
    doneChan := make(chan bool, queueSize)

    for range numCores {
        go c.animationThread(animChan, doneChan)
    }

	for c.animPit = range c.ticker.C {
        numAnims := 0
		for i, anim := range c.animList {
			if anim == nil || anim.IsStopped() {
				continue
			}
            numAnims++
            animChan <- i
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

func (c *Controller) animationThread(animChan <-chan int, doneChan chan<- bool) {
    for animId := range animChan {
        if !c.animList[animId].Update(c.animPit) {
            c.animList[animId] = nil
        }
        doneChan <- true
    }
}
