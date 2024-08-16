//go:build cameraOpenCV

package ledgrid

import (
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
	cut       image.Rectangle
	mat       gocv.Mat
	running   bool
}

func NewCamera(pos, size geom.Point) *Camera {
	c := &Camera{Pos: pos, Size: size, cut: image.Rect(0, 80, 320, 160)}
	c.CanvasObjectEmbed.ExtendCanvasObject(c)
	c.mat = gocv.NewMatWithSize(c.cut.Dx(), c.cut.Dy(), gocv.MatTypeCV8UC3)
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
	c.dev.Set(gocv.VideoCaptureZoom, 0)
	c.running = true
}

func (c *Camera) Stop() {
	var err error

	if !c.running {
		return
	}
	err = c.dev.Close()
	if err != nil {
		log.Fatalf("failed to close device: %v", err)
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
		log.Fatal("Device closed")
	}
	return true
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
	draw.CatmullRom.Scale(canv.img, rect.Add(refPt).Int(), c.img, c.cut, draw.Over, nil)
}
