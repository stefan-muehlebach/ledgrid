package ledgrid

import (
	"image"
	"log"

	"gocv.io/x/gocv"
	"golang.org/x/image/draw"
)

//----------------------------------------------------------------------------

const (
	camDevId      = 0
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
	scaler  draw.Scaler
	webcam  *gocv.VideoCapture
	mat     gocv.Mat
	anim    Animation
}

func NewCamera(lg *LedGrid) *Camera {
	var err error

	c := &Camera{}
	c.VisualEmbed.Init("Camera")
	c.lg = lg

	c.imgRect = image.Rect(40, 0, 280, 240)
	c.scaler = draw.BiLinear.NewScaler(10, 10, camHeight, camHeight)

	c.webcam, err = gocv.VideoCaptureDevice(camDevId)
	if err != nil {
		log.Fatalf("Couldn't open device: %v", err)
	}
	c.webcam.Set(gocv.VideoCaptureFrameWidth, camWidth)
	c.webcam.Set(gocv.VideoCaptureFrameHeight, camHeight)
	c.mat = gocv.NewMat()
	c.anim = NewInfAnimation(c.Update)
	// c.anim.Start()
	theAnimator.AddAnimations(c.anim)
	return c
}

func (c *Camera) ParamList() []*Bounded[float64] {
	return nil
}

func (c *Camera) Update(t float64) {
	var err error

	if !c.webcam.Read(&c.mat) {
		log.Fatal("Device closed")
	}
	c.img, err = c.mat.ToImage()
	if err != nil {
		log.Fatalf("Couldn't convert image: %v", err)
	}
}

func (c *Camera) Draw() {
	c.scaler.Scale(c.lg, c.lg.Bounds(), c.img, c.imgRect, draw.Src, nil)
}

func (c *Camera) SetVisible(visible bool) {
	c.VisualEmbed.SetVisible(visible)
	if visible {
		c.anim.Start()
	} else {
		c.anim.Stop()
	}
}
