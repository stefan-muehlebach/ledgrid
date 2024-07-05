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
	oversize         = 5.0
	pixelHost        = "raspi-2"
	pixelPort   uint = 5333
	gammaValue       = 3.0
	refreshRate      = 30 * time.Millisecond
	backAlpha        = 1.0
	c1Size1          = geom.Point{2.0, 2.0}.Mul(oversize)
	c1Size2          = geom.Point{6.0, 6.0}.Mul(oversize)
	c1Pos1           = geom.Point{17.0, 5.0}.Mul(oversize)
	c1Pos2           = geom.Point{3.0, 5.0}.Mul(oversize)
	c2Size           = geom.Point{3.0, 3.0}.Mul(oversize)
	c2Pos            = geom.Point{3.0, 5.0}.Mul(oversize)
	c2PosSize        = geom.Point{14.0, 10.0}.Mul(oversize)

	AnimCtrl *Controller
)

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

func ChasingCircles(ctrl *Controller) {
	c1 := &Ellipse{c1Pos1, c1Size1, 0.0, color.OrangeRed, color.Black}
	c1pos := NewPositionAnimation(&c1.Pos, c1Pos2, 2*time.Second)
	c1pos.AutoReverse = true
	c1radius := NewSizeAnimation(&c1.Size, c1Size2, 1*time.Second)
	c1radius.AutoReverse = true
	c1radius.RepeatCount = 1
	c1color := NewColorAnimation(&c1.FillColor, color.Gold, time.Second)
	c1color.AutoReverse = true
	c1color.RepeatCount = 1

	c2 := &Ellipse{c2Pos, c2Size, 0.0, color.Lime, color.DeepSkyBlue}
	c2pos := NewPathAnimation(&c2.Pos, circlePathFunc, c2PosSize, 4*time.Second)
	c2color := NewColorAnimation(&c2.FillColor, color.DeepSkyBlue, 2*time.Second)
	c2color.AutoReverse = true
	c2color.RepeatCount = AnimationRepeatForever

	ctrl.Add(c2, c1)

	c1pos.Start()
	c1radius.Start()
	c1color.Start()
	c2pos.Start()
	c2color.Start()

	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			c1pos.Start()
			c1radius.Start()
			c1color.Start()
			c2pos.Start()
		}
	}()
}

func CirclingCircles(ctrl *Controller) {
    pos1 := geom.Point{2.0, 2.0}.Mul(oversize)
    pos2 := geom.Point{18.0, 2.0}.Mul(oversize)
    pos3 := geom.Point{18.0, 8.0}.Mul(oversize)
    pos4 := geom.Point{2.0, 8.0}.Mul(oversize)
    cSize := geom.Point{3.0, 3.0}.Mul(oversize)

    c1 := &Ellipse{pos1, cSize, 0.0, color.OrangeRed, color.Black}
    c2 := &Ellipse{pos2, cSize, 0.0, color.MediumSeaGreen, color.Black}
    c3 := &Ellipse{pos3, cSize, 0.0, color.SkyBlue, color.Black}
    c4 := &Ellipse{pos4, cSize, 0.0, color.Gold, color.Black}

    ctrl.Add(c1, c2, c3, c4)
}

func GrowingCircles(ctrl *Controller) {
	go func() {
		for {
			rnd := 3 + rand.Intn(3)
			pos := geom.Point{rand.Float64() * float64(width), rand.Float64() * float64(height)}
			size1 := geom.Point{0.1, 0.1}
			size2 := geom.Point{float64(rnd) * 6.0, float64(rnd) * 6.0}
			col1 := color.RandColor()
			col2 := col1.Alpha(0.0)
			dur := time.Duration(rnd) * time.Second
			c := &Ellipse{pos, size1, 0.0, col1, col1}
			ctrl.Add(c)
			cRad := NewSizeAnimation(&c.Size, size2, dur)
			cColor := NewColorAnimation(&c.FillColor, col2, dur)
			cRad.Start()
			cColor.Start()
			time.Sleep(time.Duration(rnd) * time.Second)
		}
	}()
}

func ChasingRectangles(ctrl *Controller) {
	r1 := &Rectangle{geom.Point{17.0, 5.0}, geom.Point{2.0, 6.0}, 0.0, color.GreenYellow, color.Black}
	r2 := &Rectangle{geom.Point{3.0, 5.0}, geom.Point{6.0, 6.0}, 0.0, color.Red, color.Black}
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

    // CirclingCircles(ctrl)
	// GrowingCircles(ctrl)
	ChasingCircles(ctrl)
	// ChasingRectangles(ctrl)

	SignalHandler()

	ctrl.Stop()
	ledGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(ledGrid)
	pixCtrl.Close()

}
