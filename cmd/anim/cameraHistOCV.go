//go:build cameraOpenCV

package main

import (
	"image"
	gocolor "image/color"
	"log"
	"math"
	"sync"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/ledgrid"

	"gocv.io/x/gocv"
	"golang.org/x/image/draw"
)

// Die zweite Kamera-Umsetzung verwendet OpenCV und kann/wird/sollte spaeter
// auch dazu verwendet werden, wenn statt der Kamera-Bilder eine Interpretation
// davon angezeigt werden soll.
type HistCamera struct {
	ledgrid.CanvasObjectEmbed
	Pos, Size        geom.Point
	Rect             image.Rectangle
	dev              *gocv.VideoCapture
	imgIdx, histLen  int
	rawImg           *image.RGBA
	grayImg          []*image.Gray
	srcMask, dstMask *image.Alpha
	imgMutex         []*sync.RWMutex
	scaler           draw.Scaler
	srcRect          image.Rectangle
	running          bool
	mat              [2]gocv.Mat
	matMutex         [2]*sync.RWMutex
	matIdx           int
	doneChan         chan bool
}

func NewHistCamera(pos, size geom.Point, histLen int) *HistCamera {
	c := &HistCamera{Pos: pos, Size: size}
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
	rect := geom.Rectangle{Max: c.Size}
	refPt := c.Pos.Sub(c.Size.Div(2.0))
	c.Rect = rect.Add(refPt).Int()

	c.imgIdx = -1
	c.histLen = histLen

	imgRect := image.Rectangle{Max: image.Point{camWidth, camHeight}}
	c.rawImg = image.NewRGBA(imgRect)
	c.grayImg = make([]*image.Gray, histLen)
	c.imgMutex = make([]*sync.RWMutex, histLen)
	for i := range histLen {
		c.grayImg[i] = image.NewGray(imgRect)
		c.imgMutex[i] = &sync.RWMutex{}
	}
	c.srcMask = image.NewAlpha(imgRect)
	c.dstMask = image.NewAlpha(c.Rect)
	c.scaler = draw.CatmullRom.NewScaler(c.Rect.Dx(), c.Rect.Dy(), c.srcRect.Dx(), c.srcRect.Dy())

	for i := range 2 {
		c.mat[i] = gocv.NewMatWithSize(camWidth, camHeight, gocv.MatTypeCV8UC3)
		c.matMutex[i] = &sync.RWMutex{}
	}
	c.matIdx = -1

	c.doneChan = make(chan bool)

	ledgrid.AnimCtrl.Add(c)
	return c
}

func (c *HistCamera) Duration() time.Duration {
	return time.Duration(0)
}

func (c *HistCamera) SetDuration(dur time.Duration) {}

func (c *HistCamera) Start() {
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
	go c.captureThread(c.doneChan)
	c.running = true
}

func (c *HistCamera) Suspend() {
	var err error

	if !c.running {
		return
	}
	c.doneChan <- true
	err = c.dev.Close()
	if err != nil {
		log.Fatalf("Failed to close device: %v", err)
	}
	c.dev = nil
	c.running = false
}

func (c *HistCamera) captureThread(done <-chan bool) {
	var err error
	var img image.Image
	var srcMask, dstMask *image.Uniform
	var srcVal uint8 = 0x10
	var dstVal uint8 = 0xff

	srcMask = image.NewUniform(gocolor.Alpha{srcVal})
	dstMask = image.NewUniform(gocolor.Alpha{dstVal})

	ticker := time.NewTicker((camFrameRate + 10) * time.Millisecond)
ML:
	for {
		select {
		case <-ticker.C:
			idx := (c.matIdx + 1) % 2
			if !c.dev.Read(&c.mat[idx]) {
				log.Fatalf("Failed to grab and decode frames")
			}
			gocv.Flip(c.mat[idx], &c.mat[idx], 1)
			img, err = c.mat[idx].ToImage()
			if err != nil {
				log.Fatalf("Couldn't convert image: %v", err)
			}
			c.matIdx = idx

			c.imgMutex[0].Lock()
			draw.Draw(c.rawImg, img.Bounds(), img, image.Point{}, draw.Over)
			c.imgMutex[0].Unlock()

			c.imgMutex[1].Lock()
			draw.Draw(c.grayImg[0], img.Bounds(), img, image.Point{}, draw.Over)
			draw.Copy(c.grayImg[1], image.Point{}, img, img.Bounds(), draw.Over, &draw.Options{
				DstMask: dstMask,
				SrcMask: srcMask,
			})
			c.imgMutex[1].Unlock()
		case <-done:
			break ML
		}
	}
}

func (c *HistCamera) Continue() {}

func (c *HistCamera) IsRunning() bool {
	return c.running
}

// Mit der Methode Update wird das Graustufenbild der Kamera vom Bild mit den
// aufkumulierten Bildern subtrahiert und das Resultat fuer die belegung
// der Source-Mask verwendet.
func (c *HistCamera) Update(pit time.Time) bool {
	for i, val0 := range c.grayImg[0].Pix {
		val1 := c.grayImg[1].Pix[i]
		if val0 > val1 {
			val0, val1 = val1, val0
		}
		v := val1 - val0
		if v < 86 {
			v = 3 * v
		} else {
			v = 0xff
		}
		c.srcMask.Pix[i] = v
	}
	return true
}

func (c *HistCamera) Get(prop gocv.VideoCaptureProperties) float64 {
	return c.dev.Get(prop)
}

func (c *HistCamera) Set(prop gocv.VideoCaptureProperties, param float64) {
	c.dev.Set(prop, param)
}

func (c *HistCamera) Draw(canv *ledgrid.Canvas) {
	// Originales Kamerabild
	// c.scaler.Scale(canv.Img, c.Rect, c.rawImg, c.srcRect, draw.Over, nil)

	// Graustufenbild 1 (1:1 Kamerabild, aber eben Graustufig)
	// c.scaler.Scale(canv.Img, c.Rect, c.grayImg[0], c.srcRect, draw.Over, nil)

	// Graustufenbild 2 (Aufkumuliertes Bild)
	// c.scaler.Scale(canv.Img, c.Rect, c.grayImg[1], c.srcRect, draw.Over, nil)

	// Das Resultat schliesslich: originales Kamerabild, jedoch maskiert durch
	// die Bewegungserkennung.
	// c.scaler.Scale(canv.Img, c.Rect, c.rawImg, c.srcRect, draw.Over, &draw.Options{
	// 	SrcMask: c.srcMask,
	// })

	// Bewegungsbild, jedoch mit einfarbigem Hintergrund (sieht gespenstisch
	// aus).
	uniform := image.NewUniform(color.SkyBlue)
	c.scaler.Scale(canv.Img, c.Rect, uniform, c.srcRect, draw.Over, &draw.Options{
		SrcMask: c.srcMask,
	})
}
