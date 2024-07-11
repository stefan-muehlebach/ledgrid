package main

import (
	"image"
	"math"
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
	p1Pos1           = ConvertPos(geom.Point{1.0, 1.0})
	p1Pos2           = ConvertPos(geom.Point{18.0, 1.0})
	p2Pos1           = ConvertPos(geom.Point{18.0, 8.0})
	p2Pos2           = ConvertPos(geom.Point{1.0, 8.0})
	l1Pos1           = ConvertPos(geom.Point{2.0, 2.0})
	l1Pos2           = ConvertPos(geom.Point{16.0, 3.0})
	r1Pos1           = ConvertPos(geom.Point{0.5, 4.5})
	r1Pos2           = ConvertPos(geom.Point{8.5, 4.5})
	r2Pos1           = ConvertPos(geom.Point{18.5, 4.5})
	r2Pos2           = ConvertPos(geom.Point{10.5, 4.5})
	rSize1           = ConvertSize(geom.Point{17.0, 1.0})
	rSize2           = ConvertSize(geom.Point{1.0, 9.0})

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

var (
	c1Pos1    = ConvertPos(geom.Point{16.5, 4.5})
	c1Pos2    = ConvertPos(geom.Point{2.5, 4.5})
	c1Size1   = ConvertSize(geom.Point{10.0, 10.0})
	c1Size2   = ConvertSize(geom.Point{3.0, 3.0})
	c2Pos     = ConvertPos(geom.Point{2.5, 4.5})
	c2Size1   = ConvertSize(geom.Point{10.0, 10.0})
	c2Size2   = ConvertSize(geom.Point{3.0, 3.0})
	c2PosSize = ConvertSize(geom.Point{14.0, 10.0})
)

func ChasingCircles(ctrl *Controller) {
	pal := ledgrid.NewGradientPaletteByList("Palette", true,
		ledgrid.ColorMap["DeepSkyBlue"].Color(0),
		ledgrid.ColorMap["Lime"].Color(0),
		ledgrid.ColorMap["Teal"].Color(0),
		ledgrid.ColorMap["SkyBlue"].Color(0),
	)

	c1 := NewEllipse(c1Pos1, c1Size1, color.Gold)

	c1pos := NewPositionAnimation(&c1.Pos, c1Pos2, 2*time.Second)
	c1pos.AutoReverse = true
	c1pos.RepeatCount = AnimationRepeatForever

	c1size := NewSizeAnimation(&c1.Size, c1Size2, time.Second)
	c1size.AutoReverse = true
	c1size.RepeatCount = AnimationRepeatForever

	c1color := NewColorAnimation(&c1.BorderColor, color.OrangeRed, time.Second)
	c1color.AutoReverse = true
	c1color.RepeatCount = AnimationRepeatForever

	c2 := NewEllipse(c2Pos, c2Size1, color.Lime)

	c2pos := NewPathAnimation(&c2.Pos, EllipsePath, c2PosSize, 4*time.Second)
	c2pos.RepeatCount = AnimationRepeatForever
	c2pos.Curve = AnimationLinear

    c2size := NewSizeAnimation(&c2.Size, c2Size2, time.Second)
    c2size.AutoReverse = true
    c2size.RepeatCount = AnimationRepeatForever

	c2color := NewPaletteAnimation(&c2.BorderColor, pal, 2*time.Second)
	c2color.RepeatCount = AnimationRepeatForever
	c2color.Curve = AnimationLinear

	ctrl.Add(c2, c1)

	c1pos.Start()
	c1size.Start()
	c1color.Start()
	c2pos.Start()
    c2size.Start()
	c2color.Start()
}

func CircleAnimation(ctrl *Controller) {
	c1 := NewEllipse(c1Pos1, c1Size1, color.OrangeRed)

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

	c1 := NewEllipse(pos1, cSize, color.OrangeRed)
	c2 := NewEllipse(pos2, cSize, color.MediumSeaGreen)
	c3 := NewEllipse(pos3, cSize, color.SkyBlue)
	c4 := NewEllipse(pos4, cSize, color.Gold)

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
			c := NewEllipse(pos, size1, borderColor1)
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

func PushingRectangles(ctrl *Controller) {
	r1 := NewRectangle(r1Pos1, rSize2, color.GreenYellow)
	r2 := NewRectangle(r2Pos2, rSize1, color.SkyBlue)

	r1Pos := NewPositionAnimation(&r1.Pos, r1Pos2, 2*time.Second)
	r1Pos.AutoReverse = true
	r1Pos.RepeatCount = AnimationRepeatForever

	r1Size := NewSizeAnimation(&r1.Size, rSize1, 2*time.Second)
	r1Size.AutoReverse = true
	r1Size.RepeatCount = AnimationRepeatForever

	r2Pos := NewPositionAnimation(&r2.Pos, r2Pos1, 2*time.Second)
	r2Pos.AutoReverse = true
	r2Pos.RepeatCount = AnimationRepeatForever

	r2Size := NewSizeAnimation(&r2.Size, rSize2, 2*time.Second)
	r2Size.AutoReverse = true
	r2Size.RepeatCount = AnimationRepeatForever

	ctrl.Add(r1, r2)
	r1Pos.Start()
	r1Size.Start()
	r2Pos.Start()
	r2Size.Start()

}

var (
	r3Pos1  = ConvertPos(geom.Point{4.5, 4.5})
	r3Pos2  = ConvertPos(geom.Point{14.5, 4.5})
	r3Size1 = ConvertSize(geom.Point{7.0, 5.0})
	r4Pos1  = ConvertPos(geom.Point{14.5, 4.5})
	r4Pos2  = ConvertPos(geom.Point{4.5, 4.5})
	r4Size1 = ConvertSize(geom.Point{7.0, 5.0})
)

func MovingRectangles(ctrl *Controller) {
	r3 := NewRectangle(r3Pos1, r3Size1, color.GreenYellow)
	r4 := NewRectangle(r4Pos1, r4Size1, color.SkyBlue)

	r3Angle := NewFloatAnimation(&r3.Angle, math.Pi, 2*time.Second)
	r4Angle := NewFloatAnimation(&r4.Angle, -math.Pi, 2*time.Second)

	r3Color := NewColorAnimation(&r3.BorderColor, color.GreenYellow.Bright(0.5), time.Second)
	r4Color := NewColorAnimation(&r4.BorderColor, color.SkyBlue.Bright(0.5), time.Second)

	tl := NewTimeline(4 * time.Second)
	tl.RepeatCount = AnimationRepeatForever
	tl.Add(1*time.Second, r3Angle, r4Angle)
	tl.Add(2*time.Second, r3Color, r4Color)

	// seq := NewSequence(r3Angle, r4Angle)
	// seq.RepeatCount = AnimationRepeatForever

	ctrl.Add(r3, r4)
	tl.Start()
}

//----------------------------------------------------------------------------

func main() {
	pixCtrl := ledgrid.NewNetPixelClient(pixelHost, pixelPort)
	pixCtrl.SetGamma(gammaValue, gammaValue, gammaValue)
	pixCtrl.SetMaxBright(255, 255, 255)

	ledGrid := ledgrid.NewLedGrid(gridSize)
	ctrl := NewController(pixCtrl, ledGrid)

	// Piiiiixels(ctrl)
	// CirclingCircles(ctrl)
	// GrowingCircles(ctrl)
	ChasingCircles(ctrl)
	// CircleAnimation(ctrl)
	// PushingRectangles(ctrl)
	// MovingRectangles(ctrl)

	SignalHandler()

	ctrl.Stop()
	ledGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(ledGrid)
	pixCtrl.Close()

}
