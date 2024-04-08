package ledgrid

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"log"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
	"golang.org/x/image/draw"
)

//----------------------------------------------------------------------------

const (
	camDevName    = "/dev/video0"
	camWidth      = 320
	camHeight     = 240
	camFrameRate  = 30
	camBufferSize = 4
)

type Camera struct {
	VisualizableEmbed
	lg     *LedGrid
	img    image.Image
    imgRect image.Rectangle
	scaler draw.Scaler
	dev    *device.Device
}

func NewCamera(lg *LedGrid) *Camera {
	var err error

	c := &Camera{}
	c.VisualizableEmbed.Init("Camera")
	c.lg = lg
	c.dev, err = device.Open(camDevName,
		device.WithIOType(v4l2.IOTypeMMAP),
		device.WithPixFormat(v4l2.PixFormat{PixelFormat: v4l2.PixelFmtMJPEG,
			Width: uint32(camWidth), Height: uint32(camHeight)}),
		device.WithFPS(uint32(camFrameRate)),
		device.WithBufferSize(uint32(camBufferSize)),
	)
	if err != nil {
		log.Fatalf("failed to open device: %v", err)
	}
	if err := c.dev.Start(context.TODO()); err != nil {
		log.Fatalf("failed to start stream: %v", err)
	}
    c.imgRect = image.Rect(40, 0, 280, 240)
	c.scaler = draw.BiLinear.NewScaler(10, 10, camHeight, camHeight)
	return c
}

func (c *Camera) Update(dt time.Duration) bool {
	var err error

	dt = c.AnimatableEmbed.Update(dt)
	frame := <-c.dev.GetOutput()
	reader := bytes.NewReader(frame)
	c.img, err = jpeg.Decode(reader)
	if err != nil {
		log.Fatalf("failed to decode data: %v", err)
	}
	return true
}

func (c *Camera) Draw() {
	c.scaler.Scale(c.lg, c.lg.Bounds(), c.img, c.imgRect, draw.Src, nil)
}
