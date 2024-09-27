//go:build cameraV4L2

package main

import (
	"context"
	"image"
	"image/jpeg"
	"log"
	"math"
	"sync"
	"time"

	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"golang.org/x/image/draw"
)

type Camera struct {
	ledgrid.CanvasObjectEmbed
	Pos, Size geom.Point
	dev       *v4l.Device
	imgIdx    int
	img       [2]image.Image
	imgMutex  [2]*sync.RWMutex
	scaler    draw.Scaler
	srcRect   image.Rectangle
	Mask      *image.Alpha
	doneChan  chan bool
	cancel    context.CancelFunc
	running   bool
}

func NewCamera(pos, size geom.Point) *Camera {
	c := &Camera{Pos: pos, Size: size}
	c.CanvasObjectEmbed.Extend(c)
	dstRatio := size.X / size.Y
	srcRatio := float64(camWidth) / float64(camHeight)
	if dstRatio > srcRatio {
		h := camWidth / dstRatio
		m := (camHeight - h) / 2.0
		c.srcRect = image.Rect(0, int(math.Round(m)), camWidth, int(math.Round(m+h)))
	} else {
		w := camHeight * dstRatio
		m := (camWidth - w) / 2.0
		c.srcRect = image.Rect(int(math.Round(m)), 0, int(math.Round(m+w)), camHeight)
	}
	c.imgIdx = -1
	c.imgMutex[0] = &sync.RWMutex{}
	c.imgMutex[1] = &sync.RWMutex{}
	c.scaler = draw.CatmullRom.NewScaler(int(size.X), int(size.Y), c.srcRect.Dx(), c.srcRect.Dy())
	c.doneChan = make(chan bool)
	c.Mask = image.NewAlpha(image.Rectangle{Max: size.Int()})
	for i := range c.Mask.Pix {
		c.Mask.Pix[i] = 0xff
	}
	ledgrid.AnimCtrl.Add(c)
	return c
}

func (c *Camera) Duration() time.Duration {
	return time.Duration(0)
}

func (c *Camera) SetDuration(dur time.Duration) {}

func (c *Camera) Start() {
	var err error

	if c.running {
		return
	}
	c.dev, err = v4l.Open(camDevName)
	if err != nil {
		log.Fatalf("failed to open device: %v", err)
	}
	cfg, err := c.dev.GetConfig()
	if err != nil {
		log.Fatalf("failed to read configuration: %v", err)
	}
	cfg.Format = mjpeg.FourCC
	cfg.Width = camWidth
	cfg.Height = camHeight
	cfg.FPS = v4l.Frac{uint32(camFrameRate), 1}
	err = c.dev.SetConfig(cfg)
	if err != nil {
		log.Fatalf("failed to write configuration back: %v", err)
	}

	err = c.dev.TurnOn()
	if err != nil {
		log.Fatalf("failed to turn on camera: %v", err)
	}

	go c.captureThread(c.doneChan)
	c.running = true
}

func (c *Camera) Suspend() {
	if !c.running {
		return
	}
	c.doneChan <- true
	c.dev.TurnOff()
	c.running = false
}

func (c *Camera) captureThread(done <-chan bool) {
	var err error
	var buf *v4l.Buffer

	ticker := time.NewTicker((camFrameRate + 10) * time.Millisecond)
ML:
	for {
		select {
		case <-ticker.C:
			buf, err = c.dev.Capture()
			if err != nil {
				log.Fatalf("failed to capture image data: %v", err)
			}
			idx := (c.imgIdx + 1) % 2
			c.imgMutex[idx].Lock()
			c.img[idx], err = jpeg.Decode(buf)
			if err != nil {
				log.Fatalf("failed to decode data: %v", err)
			}
			c.imgMutex[idx].Unlock()
			c.imgIdx = idx

		case <-done:
			break ML
		}
	}
}

func (c *Camera) Continue() {}

func (c *Camera) IsRunning() bool {
	return c.running
}

func (c *Camera) Update(pit time.Time) bool {
	return true
}

func (c *Camera) Draw(canv *ledgrid.Canvas) {
	idx := c.imgIdx
	if idx < 0 {
		return
	}
	c.imgMutex[idx].RLock()
	rect := geom.Rectangle{Max: c.Size}
	refPt := c.Pos.Sub(c.Size.Div(2.0))
	c.scaler.Scale(canv.Img, rect.Add(refPt).Int(), c.img[idx], c.srcRect,
		draw.Over, &draw.Options{DstMask: c.Mask})
	c.imgMutex[idx].RUnlock()
}
