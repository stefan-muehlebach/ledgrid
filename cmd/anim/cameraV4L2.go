//go:build cameraV4L2

package main

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"log"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"golang.org/x/image/draw"
)

type Camera struct {
	CanvasObjectEmbed
	Pos, Size geom.Point
	dev       *device.Device
	img       image.Image
	cut       image.Rectangle
	cancel    context.CancelFunc
	running   bool
}

func NewCamera(pos, size geom.Point) *Camera {
	c := &Camera{Pos: pos, Size: size, cut: image.Rect(0, 80, 320, 160)}
	c.CanvasObjectEmbed.ExtendCanvasObject(c)
	animCtrl.Add(c)
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
	c.running = true
}

func (c *Camera) Stop() {
	var err error

	if !c.running {
		return
	}
	c.cancel()
	if err = c.dev.Close(); err != nil {
		log.Fatalf("failed to close device: %v", err)
	}
	c.dev = nil
	c.running = false
}

func (c *Camera) Continue() {}

func (c *Camera) IsStopped() bool {
	return !c.running
}

func (c *Camera) Update(pit time.Time) bool {
	var err error
	var frame []byte
	var ok bool

	if frame, ok = <-c.dev.GetOutput(); !ok {
		log.Printf("no frame to process")
		return true
	}
	reader := bytes.NewReader(frame)
	c.img, err = jpeg.Decode(reader)
	if err != nil {
		log.Fatalf("failed to decode data: %v", err)
	}
	return true
}

func (c *Camera) Draw(canv *Canvas) {
	if c.img == nil {
		return
	}
	rect := geom.Rectangle{Max: c.Size}
	refPt := c.Pos.Sub(c.Size.Div(2.0))
	draw.CatmullRom.Scale(canv.img, rect.Add(refPt).Int(), c.img, c.cut, draw.Over, nil)
}
