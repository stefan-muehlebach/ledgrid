//go:build cameraV4L2

package main

import (
	"context"
	"errors"
	"image"
	"log"
	"math"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"

	"golang.org/x/image/draw"
)

const (
	camDevName    = "/dev/video0"
	camDevId      = 0
	camWidth      = 160
	camHeight     = 120
	camFrameRate  = 30
	camBufferSize = 1
)

type CameraCtrl struct {
	ID  v4l2.CtrlID
	Val v4l2.CtrlValue
}

type camera struct {
	ledgrid.CanvasObjectEmbed
	Pos, Size        geom.Point
	dev              *device.Device
	scaler           draw.Scaler
	srcRect, dstRect image.Rectangle
	running          bool
	ctx              context.Context
	ctrlList         []CameraCtrl
	capture          func(frame *device.Frame)
}

func (c *camera) Init(pos, size geom.Point, ctx context.Context) {
	c.Pos = pos
	c.Size = size
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
	c.ctx = ctx
	ledgrid.AnimCtrl.Add(c)
}

func (c *camera) SetDuration(dur time.Duration) {}

func (c *camera) Duration() time.Duration {
	return time.Duration(0)
}

func (c *camera) Start() {
	c.StartAt(time.Now())
}

func (c *camera) StartAt(t time.Time) {
	var err error

	if c.running {
		return
	}
	c.dev, err = device.Open(
		camDevName,
		device.WithPixFormat(v4l2.PixFormat{
			Width:       uint32(camWidth),
			Height:      uint32(camHeight),
			PixelFormat: v4l2.PixelFmtRGB24,
		}),
		device.WithFPS(uint32(camFrameRate)),
	)
	if err != nil {
		log.Fatalf("failed to open device: %v", err)
	}

	for _, ctrl := range c.ctrlList {
		if err := v4l2.SetControlValue(c.dev.Fd(), ctrl.ID,
			ctrl.Val); err != nil {
			log.Printf("failed to set control %v: %v", ctrl.ID, err)
		}
	}

	if err := c.dev.Start(c.ctx); err != nil {
		log.Fatalf("failed to start stream: %s", err)
	}

	go c.captureThread()
	c.running = true
}

func (c *camera) Suspend() {
	if !c.running {
		return
	}
	if err := c.dev.Stop(); err != nil {
		log.Fatalf("failed to suspend stream: %s", err)
	}
	c.running = false
}

func (c *camera) Continue() {
	if c.running {
		return
	}
	if err := c.dev.Start(c.ctx); err != nil {
		log.Fatalf("failed to continue stream: %s", err)
	}
	c.running = true
}

func (c *camera) IsRunning() bool {
	return c.running
}

func (c *camera) Update(pit time.Time) bool {
	return true
}

func (c *camera) captureThread() {
	for frame := range c.dev.GetFrames() {
		c.capture(frame)
		frame.Release()
	}
	log.Printf("captureThread() is terminating")
	c.dev.Close()
}

// ---------------------------------------------------------------------------

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

func copyGray(dst *image.Gray, src []byte) {
	numPixels := dst.Rect.Dx() * dst.Rect.Dy()
	for pix := 0; pix < numPixels; pix++ {
		idxSrc := 3 * pix
		idxDst := pix
		dst.Pix[idxDst] = src[idxSrc]
	}
}

const (
	grayAmpl = 3
	grayMax  = 255 / grayAmpl
)

func diffRGB(dst *image.Alpha, src0, src1 *image.RGBA) error {
	if dst.Rect.Size() != src0.Rect.Size() {
		return errors.New("dst and src0 don't have the same dimension")
	}
	if dst.Rect.Size() != src1.Rect.Size() {
		return errors.New("dst and src1 don't have the same dimension")
	}
	numPixels := dst.Rect.Dx() * dst.Rect.Dy()
	for i := range numPixels {
		d := byte(0)
		for j := range 3 {
			v0 := src0.Pix[4*i+j]
			v1 := src1.Pix[4*i+j]
			if v0 > v1 {
				v0, v1 = v1, v0
			}
			d = max(d, (v1 - v0))
		}
		if d < grayMax {
			d *= grayAmpl
		} else {
			d = 0xff
		}
		v2 := dst.Pix[i]
		if d > v2 {
			v2 = d
		} else if v2 >= 2 {
			v2 = v2 - 2
		}
		dst.Pix[i] = v2
	}
	return nil
}

func diffGray(dst *image.Alpha, src0, src1 *image.Gray) error {
	if dst.Rect.Size() != src0.Rect.Size() {
		return errors.New("dst and src0 don't have the same dimension")
	}
	if dst.Rect.Size() != src1.Rect.Size() {
		return errors.New("dst and src1 don't have the same dimension")
	}
	for i, val0 := range src0.Pix {
		val1 := src1.Pix[i]
		val2 := dst.Pix[i]
		if val0 > val1 {
			val0, val1 = val1, val0
		}
		d := val1 - val0
		if d < grayMax {
			d *= grayAmpl
		} else {
			d = 0xff
		}
		if d > val2 {
			val2 = d
		} else if val2 >= 2 {
			val2 = val2 - 2
		}
		dst.Pix[i] = val2
	}
	return nil
}
