//go:build cameraV4L2

package main

import (
	"context"
	"image"
	"log"
	"math"
	//"sync"
	"time"

    "github.com/vladimirvivien/go4vl/device"
    "github.com/vladimirvivien/go4vl/v4l2"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"golang.org/x/image/draw"
)

type Camera struct {
	ledgrid.CanvasObjectEmbed
	Pos, Size geom.Point
	dev       *device.Device
	//imgIdx    int
	//img       [2]image.Image
	//imgMutex  [2]*sync.RWMutex
	scaler    draw.Scaler
	srcRect   image.Rectangle
	doneChan  chan bool
	ctx		  context.Context
	cancel    context.CancelFunc
	running   bool
	imgOut    *image.RGBA
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
	//c.imgIdx = -1
	//c.imgMutex[0] = &sync.RWMutex{}
	//c.imgMutex[1] = &sync.RWMutex{}
	c.scaler = draw.CatmullRom.NewScaler(int(size.X), int(size.Y),
		c.srcRect.Dx(), c.srcRect.Dy())
	c.doneChan = make(chan bool)
	c.imgOut = image.NewRGBA(image.Rect(0, 0, camWidth, camHeight))
	ledgrid.AnimCtrl.Add(c)
	return c
}

// SetDuration und Duration werden bei Kameras nicht benotigt, muessen aber
// fuer ein Interface vorhanden sein.
func (c *Camera) SetDuration(dur time.Duration) {}

func (c *Camera) Duration() time.Duration {
	return time.Duration(0)
}

func (c *Camera) Start() {
	c.StartAt(time.Now())
}

func (c *Camera) StartAt(t time.Time) {
	var err error

	if c.running {
		return
	}
	c.dev, err = device.Open(
		camDevName,
		device.WithPixFormat(v4l2.PixFormat{
			Width: uint32(camWidth),
			Height: uint32(camHeight),
			PixelFormat: v4l2.PixelFmtRGB24,
		}),
	)
	if err != nil {
		log.Fatalf("failed to open device: %v", err)
	}
    defer c.dev.Close()

	//ctrl, err := c.dev.GetControl(v4l2.CtrlRotate)
	//if err != nil {
	//	log.Fatalf("failed to get control for rotation: %v", err)
	//}
	if err := v4l2.SetControlValue(c.dev.Fd(), v4l2.CtrlRotate, 2); err != nil {
		log.Fatalf("failed to set rotation: %v", err)
	}

	c.ctx, c.cancel = context.WithCancel(context.TODO())
	if err := c.dev.Start(c.ctx); err != nil {
        log.Fatalf("failed to start stream: %s", err)
	}

	go c.captureThread()
	c.running = true
}

func (c *Camera) Suspend() {
	if !c.running {
		return
	}
	if err := c.dev.Stop(); err != nil {
        log.Fatalf("failed to suspend stream: %s", err)
	}
	c.running = false
}

func (c *Camera) Continue() {
	if c.running {
		return
	}
	if err := c.dev.Start(c.ctx); err != nil {
        log.Fatalf("failed to continue stream: %s", err)
	}
	c.running = true
}

func (c *Camera) IsRunning() bool {
	return c.running
}

func (c *Camera) Update(pit time.Time) bool {
	return true
}

func (c *Camera) Draw(canv *ledgrid.Canvas) {
	rect := geom.Rectangle{Max: c.Size}
	refPt := c.Pos.Sub(c.Size.Div(2.0))
	c.scaler.Scale(canv.Img, rect.Add(refPt).Int(), c.imgOut, c.srcRect,
		draw.Over, nil)
}

func copyRGB(dst *image.RGBA, src []byte) {
    numPixels := dst.Rect.Dx() * dst.Rect.Dy()
    for pix := 0; pix < numPixels; pix++ {
        idxSrc := 3 * pix
        idxDst := 4 * pix
        dst.Pix[idxDst+0] = src[idxSrc+0]
        dst.Pix[idxDst+1] = src[idxSrc+1]
        dst.Pix[idxDst+2] = src[idxSrc+2]
        dst.Pix[idxDst+3] = 0xFF
    }
}

func (c *Camera) captureThread() {
	for frame := range c.dev.GetFrames() {
		copyRGB(c.imgOut, frame.Data)
		frame.Release()
	}
	log.Printf("captureThread() is terminating")
}
