//go:build cameraOpenCV

package main

import (
	"context"
	"image"
	"log"
	"math"
	"sync"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"

	"gocv.io/x/gocv"
	"golang.org/x/image/draw"
)

const (
	camDevName    = "/dev/video2"
	camDevId      = 0
	camWidth      = 160
	camHeight     = 120
	camFrameRate  = 30
	camBufferSize = 1
)


// Die zweite Kamera-Umsetzung verwendet OpenCV und kann/wird/sollte spaeter
// auch dazu verwendet werden, wenn statt der Kamera-Bilder eine Interpretation
// davon angezeigt werden soll.
type Camera struct {
	ledgrid.CanvasObjectEmbed
	Pos, Size geom.Point
	dev       *gocv.VideoCapture
	img       image.Image
	srcRect, dstRect  image.Rectangle
	Mask      *image.Alpha
	mat       [2]gocv.Mat
	matMutex  [2]*sync.RWMutex
	matIdx    int
	running   bool
	scaler    draw.Scaler
	doneChan  chan bool
	ctx context.Context
}

func NewCamera(pos, size geom.Point, ctx context.Context) *Camera {
	c := &Camera{Pos: pos, Size: size}
	c.CanvasObjectEmbed.Extend(c)
	dstRatio := size.X / size.Y
	srcRatio := float64(camWidth) / float64(camHeight)
	if dstRatio > srcRatio {
		h := camWidth / dstRatio
		m := (camHeight - h) / 2.0
		c.srcRect = image.Rect(0, int(math.Round(m)),
			camWidth, int(math.Round(m+h)))
	} else {
		w := camHeight * dstRatio
		m := (camWidth - w) / 2.0
		c.srcRect = image.Rect(int(math.Round(m)), 0,
			int(math.Round(m+w)), camHeight)
	}
	c.dstRect = geom.Rectangle{Max: c.Size}.Add(c.Pos).Int()
	c.scaler = draw.CatmullRom.NewScaler(c.dstRect.Dx(), c.dstRect.Dy(),
		c.srcRect.Dx(), c.srcRect.Dy())
	c.Mask = image.NewAlpha(image.Rectangle{Max: size.Int()})
	for i := range c.Mask.Pix {
		c.Mask.Pix[i] = 0xff
	}
	for i := range 2 {
		c.mat[i] = gocv.NewMatWithSize(camWidth, camHeight, gocv.MatTypeCV8UC3)
		c.matMutex[i] = &sync.RWMutex{}
	}
	c.matIdx = -1
	c.doneChan = make(chan bool)
	c.ctx = ctx
	ledgrid.AnimCtrl.Add(c)
	return c
}

func (c *Camera) Duration() time.Duration {
	return time.Duration(0)
}

func (c *Camera) SetDuration(dur time.Duration) {}

func (c *Camera) Start() {
	c.StartAt(time.Now())
}

func (c *Camera) StartAt(t time.Time) {
	var err error

	if c.running {
		return
	}
 	c.dev, err = gocv.VideoCaptureFile(camDevName)
	if err != nil {
		log.Fatalf("Couldn't open device: %v", err)
	}
	c.dev.Set(gocv.VideoCaptureFrameWidth, camWidth)
	c.dev.Set(gocv.VideoCaptureFrameHeight, camHeight)
	c.dev.Set(gocv.VideoCaptureFPS, camFrameRate)
	go c.captureThread(c.doneChan)
	c.running = true
}

func (c *Camera) Suspend() {
	var err error

	if !c.running {
		return
	}
	c.doneChan <- true
	err = c.dev.Close()
	if err != nil {
		log.Fatalf("Failed to close device: %v", err)
	}
	c.dev = nil
	c.img = nil
	c.running = false
}

func (c *Camera) Continue() {}

func (c *Camera) IsRunning() bool {
	return c.running
}

func (c *Camera) Update(pit time.Time) bool {
	return true
}

func (c *Camera) captureThread(done <-chan bool) {
	ticker := time.NewTicker((camFrameRate + 10) * time.Millisecond)
ML:
	for {
		select {
		case <-ticker.C:
			idx := (c.matIdx + 1) % 2
			c.matMutex[idx].Lock()
			if !c.dev.Read(&c.mat[idx]) {
				c.matMutex[idx].Unlock()
				log.Println("Failed to grab and decode frames")
				continue
			}
			gocv.Flip(c.mat[idx], &c.mat[idx], 1)
			c.matMutex[idx].Unlock()
			c.matIdx = idx
		case <-done:
			break ML
		}
	}
}

func (c *Camera) Get(prop gocv.VideoCaptureProperties) float64 {
	return c.dev.Get(prop)
}

func (c *Camera) Set(prop gocv.VideoCaptureProperties, param float64) {
	c.dev.Set(prop, param)
}

func (c *Camera) Draw(canv *ledgrid.Canvas) {
	var err error
	idx := c.matIdx
	if idx < 0 {
		return
	}
	c.matMutex[idx].RLock()
	c.img, err = c.mat[idx].ToImage()
	c.matMutex[idx].RUnlock()
	if err != nil {
		log.Fatalf("Couldn't convert image: %v", err)
	}
	c.scaler.Scale(canv.Img, c.dstRect, c.img, c.srcRect,
		draw.Over, &draw.Options{DstMask: c.Mask})
}
