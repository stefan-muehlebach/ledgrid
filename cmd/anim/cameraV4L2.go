//go:build cameraV4L2

package main

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"log"
	"math"
	"sync"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"golang.org/x/image/draw"
)

type Camera struct {
	ledgrid.CanvasObjectEmbed
	Pos, Size geom.Point
	DstMask   *image.Alpha
	dev       *device.Device
	imgIdx    int
	img       [2]image.Image
	imgMutex  [2]*sync.RWMutex
	scaler    draw.Scaler
	doneChan  chan bool
	mask      image.Rectangle
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
		c.mask = image.Rect(0, int(math.Round(m)), camWidth, int(math.Round(m+h)))
	} else {
		w := camHeight * dstRatio
		m := (camWidth - w) / 2.0
		c.mask = image.Rect(int(math.Round(m)), 0, int(math.Round(m+w)), camHeight)
	}
	c.DstMask = image.NewAlpha(image.Rectangle{Max: size.Int()})
	for i := range c.DstMask.Pix {
		c.DstMask.Pix[i] = 0xff
	}
	c.imgIdx = -1
	c.imgMutex[0] = &sync.RWMutex{}
	c.imgMutex[1] = &sync.RWMutex{}
    c.scaler = draw.CatmullRom.NewScaler(int(size.X), int(size.Y), c.mask.Dx(), c.mask.Dy())
	c.doneChan = make(chan bool)
	ledgrid.AnimCtrl.Add(c)
	return c
}

func (c *Camera) Duration() time.Duration {
	return time.Duration(0)
}

func (c *Camera) SetDuration(dur time.Duration) {}

func (c *Camera) Start() {
	var ctx context.Context
	var err error

	if c.running {
		return
	}
	c.dev, err = device.Open(camDevName,
		device.WithIOType(v4l2.IOTypeMMAP),
		device.WithPixFormat(v4l2.PixFormat{
			PixelFormat: v4l2.PixelFmtMJPEG,
			Width:       camWidth,
			Height:      camHeight,
		}),
		device.WithFPS(camFrameRate),
		device.WithBufferSize(camBufferSize),
	)
	if err != nil {
		log.Fatalf("failed to open device: %v", err)
	}
	ctx, c.cancel = context.WithCancel(context.TODO())
	if err = c.dev.Start(ctx); err != nil {
		log.Fatalf("failed to start stream: %v", err)
	}
	go c.captureThread(c.doneChan)
	c.running = true
}

func (c *Camera) Suspend() {
	var err error

	if !c.running {
		return
	}
	c.doneChan <- true
	c.cancel()
	if err = c.dev.Close(); err != nil {
		log.Fatalf("failed to close device: %v", err)
	}
	c.dev = nil
	c.running = false
}

func (c *Camera) captureThread(done <-chan bool) {
	var err error
	var frame []byte
	var ok bool

	ticker := time.NewTicker((camFrameRate + 10) * time.Millisecond)
ML:
	for {
		select {
		case <-ticker.C:
			if frame, ok = <-c.dev.GetOutput(); !ok {
				log.Printf("no frame to process")
				continue
			}
			reader := bytes.NewReader(frame)

			idx := (c.imgIdx + 1) % 2
			c.imgMutex[idx].Lock()
			c.img[idx], err = jpeg.Decode(reader)
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
	c.scaler.Scale(canv, rect.Add(refPt).Int(), c.img[idx], c.mask, draw.Over,
		&draw.Options{DstMask: c.DstMask})
	c.imgMutex[idx].RUnlock()
}
