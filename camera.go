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
	lg      *LedGrid
	img     image.Image
	imgRect image.Rectangle
	scaler  draw.Scaler
	dev     *device.Device
	params  []*Bounded[float64]
	cancel  context.CancelFunc
}

func NewCamera(lg *LedGrid) *Camera {
	var err error

	c := &Camera{}
	c.VisualizableEmbed.Init("Camera")
	c.lg = lg

	c.imgRect = image.Rect(40, 0, 280, 240)
	c.scaler = draw.BiLinear.NewScaler(10, 10, camHeight, camHeight)

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

func (c *Camera) Update(dt time.Duration) bool {
	var err error
	var frame []byte
	var ok bool

	dt = c.AnimatableEmbed.Update(dt)
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

func (c *Camera) Draw() {
	c.scaler.Scale(c.lg, c.lg.Bounds(), c.img, c.imgRect, draw.Src, nil)
}

func (c *Camera) SetActive(active bool) {
	var ctx context.Context
	var err error

	if active {
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
		c.VisualizableEmbed.SetActive(active)
	} else {
		c.VisualizableEmbed.SetActive(active)
		c.cancel()
		if err = c.dev.Close(); err != nil {
			log.Fatalf("failed to close device: %v", err)
		}
		c.dev = nil
	}
}
