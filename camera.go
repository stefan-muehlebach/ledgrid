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
    camZoomCtrlID = 10094861
)

type Camera struct {
	VisualizableEmbed
	lg     *LedGrid
	img    image.Image
    imgRect image.Rectangle
	scaler draw.Scaler
	dev    *device.Device
    params []*Bounded[float64]
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

    c.params = make([]*Bounded[float64], 3)
    for i, id := range []v4l2.CtrlID{v4l2.CtrlBrightness, v4l2.CtrlContrast, v4l2.CtrlSaturation} {
        ctrl, err := c.dev.GetControl(id)
        if err != nil {
            log.Fatalf("couldn't get control: %v", err)
        }
        c.params[i] = NewBounded[float64](ctrl.Name, float64(ctrl.Default),
            float64(ctrl.Minimum), float64(ctrl.Maximum), float64(ctrl.Step))
        c.params[i].SetCallback(func (oldVal, newVal float64) {
            c.dev.SetControlValue(id, int32(newVal))
        })
    }

/*
    c.params[0] = NewBounded[float64]("Brightness", 128, 0, 255, 1)
    c.params[0].SetCallback(func (oldVal, newVal float64) {
        c.dev.SetControlBrightness(int32(newVal))
    })
    c.params[1] = NewBounded[float64]("Contrast", 128, 0, 255, 1)
    c.params[1].SetCallback(func (oldVal, newVal float64) {
        c.dev.SetControlContrast(int32(newVal))
    })
    c.params[2] = NewBounded[float64]("Saturation", 128, 0, 255, 1)
    c.params[2].SetCallback(func (oldVal, newVal float64) {
        c.dev.SetControlSaturation(int32(newVal))
    })
*/
/*
    camCtrl, err := c.dev.GetControl(camZoomCtrlID)
    if err != nil {
        log.Fatalf("couldn't get zoom control: %v", err)
    }
    c.params[3] = NewBounded[float64](camCtrl.Name, float64(camCtrl.Default),
        float64(camCtrl.Minimum), float64(camCtrl.Maximum),
        float64(camCtrl.Step))
    c.params[3].SetCallback(func (oldVal, newVal float64) {
        c.dev.SetControlValue(camZoomCtrlID, int32(newVal))
    })
    log.Printf("name: %s", camCtrl.Name)
*/

	return c
}

func (c *Camera) ParamList() ([]*Bounded[float64]) {
    return c.params
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
