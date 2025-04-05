//go:build cameraV4L2

package main

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"log"
	"math"
	"sync"
	"time"

	"github.com/vladimirvivien/go4vl/device"
	"github.com/vladimirvivien/go4vl/v4l2"

	// "github.com/korandiz/v4l"
	// "github.com/korandiz/v4l/fmt/mjpeg"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colors"
	"golang.org/x/image/draw"
)

type HistCamera struct {
	ledgrid.CanvasObjectEmbed
	Pos, Size        geom.Point
	Rect             image.Rectangle
	Color            *ledgrid.UniformPalette
	dev              *device.Device
	imgIdx, histLen  int
	rawImg           *image.RGBA
	grayImg          []*image.Gray
	srcMask, dstMask *image.Alpha
	imgMutex         []*sync.RWMutex
	scaler           draw.Scaler
	srcRect          image.Rectangle
	running          bool
	cancel           context.CancelFunc
}

func NewHistCamera(pos, size geom.Point, histLen int, col colors.LedColor) *HistCamera {
	var err error

	c := &HistCamera{Pos: pos, Size: size}
	c.CanvasObjectEmbed.Extend(c)
	c.Color = ledgrid.NewUniformPalette("Uniform", col)
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

	ledgrid.AnimCtrl.Add(c)

	c.dev, err = device.Open(
		camDevName,
		device.WithPixFormat(v4l2.PixFormat{
			Width:       camWidth,
			Height:      camHeight,
			PixelFormat: v4l2.PixelFmtMJPEG,
		}))
	if err != nil {
		log.Fatalf("failed to open device: %v", err)
	}

	if err = c.dev.SetFrameRate(camFrameRate); err != nil {
		log.Fatalf("failed to set frame rate: %v", err)
	}
	return c
}

func (c *HistCamera) Duration() time.Duration {
	return time.Duration(0)
}

func (c *HistCamera) SetDuration(dur time.Duration) {}

func (c *HistCamera) StartAt(t time.Time) {
	var err error
	var ctx context.Context

	if c.running {
		return
	}

	ctx, c.cancel = context.WithCancel(context.TODO())
	if err = c.dev.Start(ctx); err != nil {
		log.Fatalf("Failed to start stream: %v", err)
	}

	go c.captureThread()
	c.running = true
}
func (c *HistCamera) Start() {
	c.StartAt(time.Now())
}

func (c *HistCamera) Suspend() {
	if !c.running {
		return
	}
	c.cancel()
	c.dev.Close()
	c.running = false
}

func (c *HistCamera) captureThread() {
	var err error
	var img image.Image
	var srcMask, dstMask *image.Uniform
	var srcVal uint8 = 0x0f
	var dstVal uint8 = 0xff

	srcMask = image.NewUniform(color.Alpha{srcVal})
	dstMask = image.NewUniform(color.Alpha{dstVal})

	for frame := range c.dev.GetOutput() {
		if len(frame) == 0 {
			log.Printf("Received frame size 0")
			continue
		}
		img, _, err = image.Decode(bytes.NewReader(frame))
		if err != nil {
			log.Fatalf("Failed to decode data: %v", err)
		}
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
	c.scaler.Scale(canv.Img, c.Rect, c.Color, c.srcRect,
		draw.Over, &draw.Options{
			SrcMask: c.srcMask,
		})
}
