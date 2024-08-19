//go:build cameraOpenCV

package ledgrid

import (
	"math"
	"image"
	"log"
	"time"

	"github.com/stefan-muehlebach/gg/geom"

	"gocv.io/x/gocv"
	"golang.org/x/image/draw"
)

// Die zweite Kamera-Umsetzung verwendet OpenCV und kann/wird/sollte spaeter
// auch dazu verwendet werden, wenn statt der Kamera-Bilder eine Interpretation
// davon angezeigt werden soll.
type Camera struct {
	CanvasObjectEmbed
	Pos, Size geom.Point
	dev       *gocv.VideoCapture
	img       image.Image
	mask      image.Rectangle
    DstMask   *image.Alpha
	mat       gocv.Mat
	running   bool
    scaler    draw.Scaler
}

func NewCamera(pos, size geom.Point) *Camera {
	c := &Camera{Pos: pos, Size: size}
	c.CanvasObjectEmbed.Extend(c)
    ratio := size.X / size.Y
    h := camWidth / ratio
    m := (camHeight - h) / 2.0
    c.mask = image.Rect(0, int(math.Round(m)), camWidth, int(math.Round(m+h)))
    c.DstMask = image.NewAlpha(image.Rectangle{Max: size.Int()})
    for i := range c.DstMask.Pix {
        c.DstMask.Pix[i] = 0xff
    }
	c.mat = gocv.NewMatWithSize(camWidth, camHeight, gocv.MatTypeCV8UC3)
    c.scaler = draw.CatmullRom.NewScaler(int(size.X), int(size.Y), c.mask.Dx(), c.mask.Dy())
	AnimCtrl.Add(c)
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
	c.dev, err = gocv.VideoCaptureDevice(camDevId)
	if err != nil {
		log.Fatalf("Couldn't open device: %v", err)
	}
	c.dev.Set(gocv.VideoCaptureFrameWidth, camWidth)
	c.dev.Set(gocv.VideoCaptureFrameHeight, camHeight)
	c.dev.Set(gocv.VideoCaptureFPS, camFrameRate)
	c.running = true
}

func (c *Camera) Stop() {
	var err error

	if !c.running {
		return
	}
	err = c.dev.Close()
	if err != nil {
		log.Fatalf("Failed to close device: %v", err)
	}
	c.dev = nil
	c.img = nil
	c.running = false
}

func (c *Camera) Continue() {}

func (c *Camera) IsStopped() bool {
	return !c.running
}

func (c *Camera) Update(pit time.Time) bool {
	if !c.dev.Read(&c.mat) {
		log.Fatalf("Failed to grab and decode frames")
	}
	return true
}

func (c *Camera) Get(prop gocv.VideoCaptureProperties) float64 {
    return c.dev.Get(prop)
}

func (c *Camera) Set(prop gocv.VideoCaptureProperties, param float64) {
    c.dev.Set(prop, param)
}

func (c *Camera) Draw(canv *Canvas) {
	var err error
	gocv.Flip(c.mat, &c.mat, 1)
	c.img, err = c.mat.ToImage()
	if err != nil {
		log.Fatalf("Couldn't convert image: %v", err)
	}
	rect := geom.Rectangle{Max: c.Size}
	refPt := c.Pos.Sub(c.Size.Div(2.0))
    c.scaler.Scale(canv.img, rect.Add(refPt).Int(), c.img, c.mask, draw.Over,
        &draw.Options{DstMask: c.DstMask})
}
