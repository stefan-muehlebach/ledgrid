package main

import (
	"image"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/gg/color"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
)

var (
	width            = 20
	height           = 10
	gridSize         = image.Point{width, height}
	pixelHost        = "raspi-2"
	pixelPort   uint = 5333
	gammaValue       = 3.0
	refreshRate      = 30 * time.Millisecond
	backAlpha        = 1.0
	c1Size1          = ConvertSize(geom.Point{2.0, 2.0})
	c1Size2          = ConvertSize(geom.Point{10.0, 10.0})
	c1Pos1           = ConvertPos(geom.Point{16.0, 4.0})
	c1Pos2           = ConvertPos(geom.Point{3.0, 4.0})
	c2Size           = ConvertSize(geom.Point{3.0, 3.0})
	c2Pos            = ConvertPos(geom.Point{3.0, 4.0})
	c2PosSize        = ConvertSize(geom.Point{14.0, 10.0})
	p1Pos1           = ConvertPos(geom.Point{1.0, 1.0})
	p1Pos2           = ConvertPos(geom.Point{18.0, 1.0})
	p2Pos1           = ConvertPos(geom.Point{18.0, 8.0})
	p2Pos2           = ConvertPos(geom.Point{1.0, 8.0})
	l1Pos1           = ConvertPos(geom.Point{2.0, 2.0})
	l1Pos2           = ConvertPos(geom.Point{16.0, 3.0})

	AnimCtrl *Controller
)

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

//----------------------------------------------------------------------------

func Piiiiixels(ctrl *Controller) {
	p1 := &Pixel{p1Pos1, color.OrangeRed}
	p2 := &Pixel{p2Pos1, color.Lime}

	// p1pos := NewPositionAnimation(&p1.Pos, p1Pos2, 3*time.Second)
	// p1pos.AutoReverse = true
	// p1pos.RepeatCount = AnimationRepeatForever

	p1color := NewPaletteAnimation(&p1.Color, ledgrid.PaletteMap["Pastell"], 2*time.Second)
	p1color.RepeatCount = AnimationRepeatForever

	p2pos := NewPositionAnimation(&p2.Pos, p2Pos2, 3*time.Second)
	p2pos.AutoReverse = true
	p2pos.RepeatCount = AnimationRepeatForever

	ctrl.Add(p1, p2)

	// p1pos.Start()
	p1color.Start()
	p2pos.Start()
}

func ChasingCircles(ctrl *Controller) {
	pal := ledgrid.NewGradientPaletteByList("Palette", true,
		ledgrid.ColorMap["DeepSkyBlue"].Color(0),
		ledgrid.ColorMap["Lime"].Color(0),
		ledgrid.ColorMap["Teal"].Color(0),
		ledgrid.ColorMap["SkyBlue"].Color(0),
	)

	c1 := &Ellipse{c1Pos1, c1Size1, ConvertLen(1.0), color.OrangeRed.Alpha(0.3), color.OrangeRed}

	c1pos := NewPositionAnimation(&c1.Pos, c1Pos2, 2*time.Second)
	c1pos.AutoReverse = true
	c1pos.RepeatCount = AnimationRepeatForever

	c1radius := NewSizeAnimation(&c1.Size, c1Size2, time.Second)
	c1radius.AutoReverse = true
	c1radius.RepeatCount = AnimationRepeatForever

	c1color := NewColorAnimation(&c1.BorderColor, color.Gold, time.Second)
	c1color.AutoReverse = true
	c1color.RepeatCount = AnimationRepeatForever

	c2 := &Ellipse{c2Pos, c2Size, ConvertLen(1.0), color.Lime.Alpha(0.3), color.Lime}

	c2pos := NewPathAnimation(&c2.Pos, circlePathFunc, c2PosSize, 4*time.Second)
	c2pos.RepeatCount = AnimationRepeatForever
	c2pos.Curve = AnimationLinear

	c2color := NewPaletteAnimation(&c2.BorderColor, pal, 2*time.Second)
	c2color.RepeatCount = AnimationRepeatForever
	c2color.Curve = AnimationLinear

	ctrl.Add(c2, c1)

	c1pos.Start()
	c1radius.Start()
	c1color.Start()
	c2pos.Start()
	c2color.Start()
}

func CircleAnimation(ctrl *Controller) {
	c1 := &Ellipse{c1Pos1, c1Size1, ConvertLen(1.0), color.OrangeRed.Alpha(0.3), color.OrangeRed}

	c1pos := NewPositionAnimation(&c1.Pos, c1Pos2, time.Second)
	c1pos.AutoReverse = true
	c1pos.RepeatCount = AnimationRepeatForever

	c1radius := NewSizeAnimation(&c1.Size, c1Size2, time.Second)
	c1radius.AutoReverse = true
	c1radius.RepeatCount = AnimationRepeatForever

	c1color := NewColorAnimation(&c1.BorderColor, color.Gold, 2*time.Second)
	c1color.AutoReverse = true
	c1color.RepeatCount = AnimationRepeatForever

	ctrl.Add(c1)

	c1pos.Start()
	c1radius.Start()
	c1color.Start()
}

func CirclingCircles(ctrl *Controller) {
	pos1 := ConvertPos(geom.Point{2.0, 2.0})
	pos2 := ConvertPos(geom.Point{18.0, 2.0})
	pos3 := ConvertPos(geom.Point{18.0, 8.0})
	pos4 := ConvertPos(geom.Point{2.0, 8.0})
	cSize := ConvertSize(geom.Point{3.0, 3.0})

	c1 := &Ellipse{pos1, cSize, ConvertLen(0.0), color.OrangeRed, color.Black}
	c2 := &Ellipse{pos2, cSize, ConvertLen(0.0), color.MediumSeaGreen, color.Black}
	c3 := &Ellipse{pos3, cSize, ConvertLen(0.0), color.SkyBlue, color.Black}
	c4 := &Ellipse{pos4, cSize, ConvertLen(0.0), color.Gold, color.Black}

	ctrl.Add(c1, c2, c3, c4)
}

func GrowingCircles(ctrl *Controller) {
	go func() {
		for {
			rnd := 3 + rand.Intn(3)
			pos := ConvertPos(geom.Point{rand.Float64() * float64(width), rand.Float64() * float64(height)})
			size1 := ConvertSize(geom.Point{0.1, 0.1})
			size2 := ConvertSize(geom.Point{float64(rnd) * 6.0, float64(rnd) * 6.0})
			borderColor1 := color.RandColor()
			borderColor2 := borderColor1.Alpha(0.0)
			fillColor1 := borderColor1.Alpha(0.3)
			fillColor2 := fillColor1.Alpha(0.0)
			dur := time.Duration(rnd) * time.Second
			c := &Ellipse{pos, size1, ConvertLen(1.0), fillColor1, borderColor1}
			ctrl.Add(c)
			cRad := NewSizeAnimation(&c.Size, size2, dur)
			cColor1 := NewColorAnimation(&c.FillColor, fillColor2, dur)
			cColor2 := NewColorAnimation(&c.BorderColor, borderColor2, dur)
			cRad.Start()
			cColor1.Start()
			cColor2.Start()
			time.Sleep(time.Duration(rnd) * time.Second)
		}
	}()
}

func ChasingRectangles(ctrl *Controller) {
	r1 := &Rectangle{geom.Point{17.0, 5.0}, geom.Point{2.0, 6.0}, ConvertLen(0.0), color.GreenYellow, color.Black}
	r2 := &Rectangle{geom.Point{3.0, 5.0}, geom.Point{6.0, 6.0}, ConvertLen(0.0), color.Red, color.Black}
	ctrl.Add(r1, r2)

	pa1 := NewPositionAnimation(&r1.Pos, geom.Point{3.0, 5.0}, 2*time.Second)
	pa1.AutoReverse = true
	pa2 := NewPathAnimation(&r2.Pos, circlePathFunc, c2PosSize, 4*time.Second)
	sa := NewSizeAnimation(&r1.Size, geom.Point{6.0, 2.0}, 2*time.Second)
	sa.AutoReverse = true
	// sa.RepeatCount = 1

	ticker := time.NewTicker(7 * time.Second)
	go func() {
		pa1.Start()
		pa2.Start()
		sa.Start()
		for range ticker.C {
			pa1.Start()
			pa2.Start()
			sa.Start()
		}
	}()
}

//----------------------------------------------------------------------------

func main() {
	pixCtrl := ledgrid.NewNetPixelClient(pixelHost, pixelPort)
	pixCtrl.SetGamma(gammaValue, gammaValue, gammaValue)
	pixCtrl.SetMaxBright(255, 255, 255)

	ledGrid := ledgrid.NewLedGrid(gridSize)
	ctrl := NewController(pixCtrl, ledGrid)

	Piiiiixels(ctrl)
	// CirclingCircles(ctrl)
	// GrowingCircles(ctrl)
	// ChasingCircles(ctrl)
	// ChasingRectangles(ctrl)
	// CircleAnimation(ctrl)

	SignalHandler()

	ctrl.Stop()
	ledGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(ledGrid)
	pixCtrl.Close()

}
