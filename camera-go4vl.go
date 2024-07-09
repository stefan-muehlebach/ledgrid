//go:build ignore
// go:build arm || arm64

package ledgrid

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"log"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"
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
	VisualEmbed
	lg      *LedGrid
	img     image.Image
	imgRect image.Rectangle
	dev     *device.Device
	params  []*Bounded[float64]
	cancel  context.CancelFunc
	anim    Animation
}

func NewCamera(lg *LedGrid) *Camera {
	var err error

	c := &Camera{}
	c.VisualEmbed.Init("Camera")
	c.lg = lg

	c.imgRect = image.Rect(40, 0, 280, 240)

	paramDev, err := device.Open(camDevName)
	if err != nil {
		log.Fatalf("failed to open device for parameter query: %v", err)
	}
	defer paramDev.Close()
	allCtrls, err := paramDev.QueryAllControls()
	if err != nil {
		log.Fatalf("failed to query controls: %v", err)
	}
	c.params = make([]*Bounded[float64], 0)
	for _, ctrl := range allCtrls {
		switch ctrl.Type {
		case v4l2.CtrlTypeInt:
			param := NewBounded[float64](ctrl.Name, float64(ctrl.Default),
				float64(ctrl.Minimum), float64(ctrl.Maximum), float64(ctrl.Step))
			param.SetCallback(func(oldVal, newVal float64) {
				c.SetParamValue(ctrl.ID, int32(newVal))
			})
			c.params = append(c.params, param)
		}
	}
	c.anim = NewInfAnimation(c.Update)
	return c
}

func (c *Camera) ParamList() []*Bounded[float64] {
	return c.params
}

func (c *Camera) SetParamValue(id v4l2.CtrlID, val v4l2.CtrlValue) {
	if c.dev != nil {
		c.dev.SetControlValue(id, val)
	}
}

func (c *Camera) Update(t float64) {
	var err error
	var frame []byte
	var ok bool

	if frame, ok = <-c.dev.GetOutput(); !ok {
		log.Printf("no frame to process")
		return
	}
	reader := bytes.NewReader(frame)
	c.img, err = jpeg.Decode(reader)
	if err != nil {
		log.Fatalf("failed to decode data: %v", err)
	}
	return
}

func (c *Camera) ColorModel() color.Model {
	return LedColorModel
}

func (c *Camera) Bounds() image.Rectangle {
	return c.imgRect
}

func (c *Camera) At(x, y int) color.Color {
	return c.img.At(x, y)
}

func (c *Camera) SetVisible(visible bool) {
	var ctx context.Context
	var err error

	if visible {
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
		c.VisualEmbed.SetVisible(visible)
		c.anim.Start()
	} else {
		c.anim.Stop()
		c.VisualEmbed.SetVisible(visible)
		c.cancel()
		if err = c.dev.Close(); err != nil {
			log.Fatalf("failed to close device: %v", err)
		}
		c.dev = nil
	}
}
