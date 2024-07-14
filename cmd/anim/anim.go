package main

import (
	"fmt"
	"image"
	"math"
	"math/rand/v2"
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
	pixelHost        = "raspi-3"
	pixelPort   uint = 5333
	gammaValue       = 3.0
	refreshRate      = 30 * time.Millisecond
	backAlpha        = 1.0

	AnimCtrl *Controller
)

//----------------------------------------------------------------------------

func GroupTest(ctrl *Controller) {
	rPos1 := ConvertPos(geom.Point{4.5, 4.5})
	rPos2 := ConvertPos(geom.Point{17.5, 4.5})
	rSize1 := ConvertSize(geom.Point{7.0, 7.0})
	rSize2 := ConvertSize(geom.Point{1.0, 1.0})

	r := NewRectangle(rPos1, rSize1, color.SkyBlue)

	aPos := NewPositionAnimation(&r.Pos, rPos2, time.Second)
	aPos.AutoReverse = true
	aSize := NewSizeAnimation(&r.Size, rSize2, time.Second)
	aSize.AutoReverse = true
	aColor := NewColorAnimation(&r.BorderColor, color.GreenYellow, time.Second)
	aColor.AutoReverse = true
	aAngle := NewFloatAnimation(&r.Angle, math.Pi, time.Second)
	aAngle.AutoReverse = true

	aGroup := NewGroup(4*time.Second, aPos, aSize, aColor, aAngle)
	aGroup.RepeatCount = AnimationRepeatForever

	ctrl.Add(r)
	aGroup.Start()
}

func SequenceTest(ctrl *Controller) {
	rPos := ConvertPos(geom.Point{9.5, 4.5})
	rSize1 := ConvertSize(geom.Point{19.0, 9.0})
	rSize2 := ConvertSize(geom.Point{1.0, 1.0})

	r := NewRectangle(rPos, rSize1, color.SkyBlue)

	aSize1 := NewSizeAnimation(&r.Size, rSize2, time.Second)
	aColor1 := NewColorAnimation(&r.BorderColor, color.OrangeRed, time.Second/2)
	aColor1.AutoReverse = true
	aColor1.RepeatCount = 2
	aColor2 := NewColorAnimation(&r.BorderColor, color.OrangeRed, time.Second/2)
	aSize2 := NewSizeAnimation(&r.Size, rSize1, time.Second)
	aSize2.Cont = true
	aColor3 := NewColorAnimation(&r.BorderColor, color.SkyBlue, 3*time.Second/2)
	aColor3.Cont = true

	aSeq := NewSequence(8*time.Second, aSize1, aColor1, aColor2, aSize2, aColor3)
	aSeq.RepeatCount = AnimationRepeatForever

	ctrl.Add(r)
	aSeq.Start()
}

func TimelineTest(ctrl *Controller) {
	r3Pos1 := ConvertPos(geom.Point{4.5, 4.5})
	r3Size1 := ConvertSize(geom.Point{7.0, 5.0})
	r4Pos1 := ConvertPos(geom.Point{14.5, 4.5})
	r4Size1 := ConvertSize(geom.Point{7.0, 5.0})

	r3 := NewRectangle(r3Pos1, r3Size1, color.GreenYellow)
	r4 := NewRectangle(r4Pos1, r4Size1, color.SkyBlue)

	aAngle1 := NewFloatAnimation(&r3.Angle, math.Pi, 2*time.Second)
	aAngle2 := NewFloatAnimation(&r4.Angle, -math.Pi, 2*time.Second)
	aAngle3 := NewFloatAnimation(&r3.Angle, -math.Pi, time.Second)
	aAngle4 := NewFloatAnimation(&r4.Angle, math.Pi, time.Second)

	r3Color := NewColorAnimation(&r3.BorderColor, color.OrangeRed, 200*time.Millisecond)
	r3Color.AutoReverse = true
	r4Color := NewColorAnimation(&r4.BorderColor, color.OrangeRed, 200*time.Millisecond)
	r4Color.AutoReverse = true

	tl := NewTimeline(5 * time.Second)
	tl.RepeatCount = AnimationRepeatForever
	tl.Add(0*time.Millisecond, aAngle1, aAngle2)
	tl.Add(2200*time.Millisecond, r3Color)
	tl.Add(2400*time.Millisecond, r4Color)
	tl.Add(2600*time.Millisecond, aAngle3)
	tl.Add(2900*time.Millisecond, aAngle4)

	ctrl.Add(r3, r4)
	tl.Start()
}

func PathTest(ctrl *Controller) {
	pos1 := ConvertPos(geom.Point{1.0, 4.0})
	pos2 := ConvertPos(geom.Point{9.0, 1.0})
	pos3 := ConvertPos(geom.Point{17.0, 4.0})
	pos4 := ConvertPos(geom.Point{9.0, 7.0})
	cSize := ConvertSize(geom.Point{2.0, 2.0})

	c1 := NewEllipse(pos1, cSize, color.OrangeRed)
	c2 := NewEllipse(pos2, cSize, color.MediumSeaGreen)
	c3 := NewEllipse(pos3, cSize, color.SkyBlue)
	c4 := NewEllipse(pos4, cSize, color.Gold)

	c1Path := NewPathAnimation(&c1.Pos, HalfCirclePathB, ConvertSize(geom.Point{16.0, 3.0}), 2*time.Second)
	c1Path.AutoReverse = true
	// c1Path.Cont = true
	c3Path := NewPathAnimation(&c3.Pos, HalfCirclePathB, ConvertSize(geom.Point{-16.0, -3.0}), 2*time.Second)
	c3Path.AutoReverse = true
	// c3Path.Cont = true

	c2Path := NewPathAnimation(&c2.Pos, HalfCirclePathA, ConvertSize(geom.Point{-3.0, 6.0}), 2*time.Second)
	c2Path.AutoReverse = true
	// c2Path.Cont = true
	c4Path := NewPathAnimation(&c4.Pos, HalfCirclePathA, ConvertSize(geom.Point{3.0, -6.0}), 2*time.Second)
	c4Path.AutoReverse = true
	// c4Path.Cont = true

	aGrp := NewGroup(5*time.Second, c1Path, c2Path, c3Path, c4Path)
	aGrp.RepeatCount = AnimationRepeatForever

	ctrl.Add(c1, c2, c3, c4)
	aGrp.Start()
}

func PolygonPathTest(ctrl *Controller) {
	cPos := ConvertPos(geom.Point{1, 1})
	cSize := ConvertSize(geom.Point{2, 2})

	polyPath := NewPolygonPath(
		geom.Point{0, 0},
		geom.Point{1, 0},
		geom.Point{1, 1},
		geom.Point{0, 1},
		geom.Point{0, 2.0 / 7.0},
		geom.Point{1.0 - 2.0/7.0, 2.0 / 7.0},
		geom.Point{1.0 - 2.0/7.0, 1.0 - 2.0/7.0},
		geom.Point{2.0 / 7.0, 1.0 - 2.0/7.0},
	)

	c1 := NewEllipse(cPos, cSize, color.GreenYellow)

	aPath := NewPathAnimation(&c1.Pos, polyPath.RelPoint,
		ConvertSize(geom.Point{7, 7}), 7*time.Second)
	aPath.AutoReverse = true
	aPath.RepeatCount = AnimationRepeatForever

	ctrl.Add(c1)
	aPath.Start()

	// fmt.Printf("PolygonPath: %+v\n", polyPath)
	// fmt.Printf("Point at 0.0: %v\n", polyPath.RelPoint(0.0))
	// fmt.Printf("Point at 0.25: %v\n", polyPath.RelPoint(0.25))
	// fmt.Printf("Point at 0.5: %v\n", polyPath.RelPoint(0.5))
	// fmt.Printf("Point at 0.75: %v\n", polyPath.RelPoint(0.75))
	// fmt.Printf("Point at 1.0: %v\n", polyPath.RelPoint(1.0))
}

func RandomWalk(ctrl *Controller) {
	rect := geom.Rectangle{Min: ConvertPos(geom.Point{1.0, 1.0}), Max: ConvertPos(geom.Point{18.0, 8.0})}
	pos1 := ConvertPos(geom.Point{1.0, 1.0})
	pos2 := ConvertPos(geom.Point{18.0, 8.0})
	cSize := ConvertSize(geom.Point{3.0, 3.0})

	c1 := NewEllipse(pos1, cSize, color.Teal)
	// c1.FillColor = c1.BorderColor
	c2 := NewEllipse(pos2, cSize, color.SeaGreen)
	// c2.FillColor = c2.BorderColor

	aPos1 := NewPositionAnimation(&c1.Pos, geom.Point{}, 1300*time.Millisecond)
	aPos1.Cont = true
	aPos1.ValFunc = RandPoint(rect)
	aPos1.RepeatCount = AnimationRepeatForever

	aPos2 := NewPositionAnimation(&c2.Pos, geom.Point{}, 901*time.Millisecond)
	aPos2.Cont = true
	aPos2.ValFunc = RandPoint(rect)
	aPos2.RepeatCount = AnimationRepeatForever

	ctrl.Add(c1, c2)
	aPos1.Start()
	aPos2.Start()
}

func Piiiiixels(ctrl *Controller) {
	p1Pos1 := ConvertPos(geom.Point{1.0, 1.0})
	// p1Pos2 := ConvertPos(geom.Point{18.0, 1.0})
	p2Pos1 := ConvertPos(geom.Point{9.0, 1.0})
	p2Range := ConvertSize(geom.Point{18.0, 7.0})
	// p2Pos2 := ConvertPos(geom.Point{1.0, 8.0})

	p1 := NewPixel(p1Pos1, color.OrangeRed)
	p2 := NewPixel(p2Pos1, color.Lime)

	// p1pos := NewPositionAnimation(&p1.Pos, p1Pos2, 3*time.Second)
	// p1pos.AutoReverse = true
	// p1pos.RepeatCount = AnimationRepeatForever

	p1color := NewPaletteAnimation(&p1.Color, ledgrid.PaletteMap["Pastell"], 2*time.Second)
	p1color.RepeatCount = AnimationRepeatForever

	p2pos := NewPathAnimation(&p2.Pos, FullCirclePathA, p2Range, 3*time.Second)
	p2pos.Curve = AnimationLinear
	p2pos.RepeatCount = AnimationRepeatForever

	// p2pos := NewPositionAnimation(&p2.Pos, p2Pos2, 3*time.Second)
	// p2pos.AutoReverse = true
	// p2pos.RepeatCount = AnimationRepeatForever

	ctrl.Add(p1, p2)

	// p1pos.Start()
	p1color.Start()
	p2pos.Start()
}

func CirclingCircles(ctrl *Controller) {
	pos1 := ConvertPos(geom.Point{2.0, 2.0})
	pos2 := ConvertPos(geom.Point{7.0, 7.0})
	pos3 := ConvertPos(geom.Point{12.0, 2.0})
	pos4 := ConvertPos(geom.Point{17.0, 7.0})
	cSize := ConvertSize(geom.Point{2.0, 2.0})

	c1 := NewEllipse(pos1, cSize, color.OrangeRed)
	c2 := NewEllipse(pos2, cSize, color.MediumSeaGreen)
	c3 := NewEllipse(pos3, cSize, color.SkyBlue)
	c4 := NewEllipse(pos4, cSize, color.Gold)

	c1Path1 := NewPathAnimation(&c1.Pos, QuarterCirclePathA, ConvertSize(geom.Point{5.0, 5.0}), time.Second)
	c1Path1.Cont = true
	c1Path2 := NewPathAnimation(&c1.Pos, QuarterCirclePathA, ConvertSize(geom.Point{5.0, -5.0}), time.Second)
	c1Path2.Cont = true
	c1Path3 := c1Path1

	c2Path1 := NewPathAnimation(&c2.Pos, QuarterCirclePathA, ConvertSize(geom.Point{-5.0, -5.0}), time.Second)
	c2Path1.Cont = true
	c2Path2 := NewPathAnimation(&c2.Pos, QuarterCirclePathA, ConvertSize(geom.Point{5.0, 5.0}), time.Second)
	c2Path2.Cont = true

	c3Path1 := NewPathAnimation(&c3.Pos, QuarterCirclePathA, ConvertSize(geom.Point{-5.0, 5.0}), time.Second)
	c3Path1.Cont = true
	c3Path2 := NewPathAnimation(&c3.Pos, QuarterCirclePathA, ConvertSize(geom.Point{5.0, -5.0}), time.Second)
	c3Path2.Cont = true

	c4Path1 := NewPathAnimation(&c4.Pos, QuarterCirclePathA, ConvertSize(geom.Point{-5.0, -5.0}), time.Second)
	c4Path1.Cont = true
	c4Path2 := NewPathAnimation(&c4.Pos, QuarterCirclePathA, ConvertSize(geom.Point{-5.0, 5.0}), time.Second)
	c4Path2.Cont = true
	c4Path3 := c4Path1

	aGrp1 := NewGroup(0, c1Path1, c2Path1)
	aGrp2 := NewGroup(0, c1Path2, c3Path1)
	aGrp3 := NewGroup(0, c1Path3, c4Path1)
	aGrp4 := NewGroup(0, c4Path2, c3Path2)
	aGrp5 := NewGroup(0, c4Path3, c2Path2)
	aSeq := NewSequence(0, aGrp1, aGrp2, aGrp3, aGrp4, aGrp5)

	ctrl.Add(c1, c2, c3, c4)
	aSeq.Start()
}

func ChasingCircles(ctrl *Controller) {
	c1Pos1 := ConvertPos(geom.Point{16.5, 4.5})
	c1Size1 := ConvertSize(geom.Point{10.0, 10.0})
	c1Size2 := ConvertSize(geom.Point{3.0, 3.0})
	c1PosSize := ConvertSize(geom.Point{-14.0, -5.0})
	c2Pos := ConvertPos(geom.Point{2.5, 4.5})
	c2Size1 := ConvertSize(geom.Point{5.0, 5.0})
	c2Size2 := ConvertSize(geom.Point{3.0, 3.0})
	c2PosSize := ConvertSize(geom.Point{14.0, 7.0})

	pal := ledgrid.NewGradientPaletteByList("Palette", true,
		ledgrid.LedColorModel.Convert(color.DeepSkyBlue).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(color.Lime).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(color.Teal).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(color.SkyBlue).(ledgrid.LedColor),
		// ledgrid.ColorMap["DeepSkyBlue"].Color(0),
		// ledgrid.ColorMap["Lime"].Color(0),
		// ledgrid.ColorMap["Teal"].Color(0),
		// ledgrid.ColorMap["SkyBlue"].Color(0),
	)

	c1 := NewEllipse(c1Pos1, c1Size1, color.Gold)

	c1pos := NewPathAnimation(&c1.Pos, FullCirclePathB, c1PosSize, 4*time.Second)
	c1pos.RepeatCount = AnimationRepeatForever
	c1pos.Curve = AnimationLinear

	// c1pos := NewPositionAnimation(&c1.Pos, c1Pos2, 2*time.Second)
	// c1pos.AutoReverse = true
	// c1pos.RepeatCount = AnimationRepeatForever

	c1size := NewSizeAnimation(&c1.Size, c1Size2, time.Second)
	c1size.AutoReverse = true
	c1size.RepeatCount = AnimationRepeatForever

	c1bcolor := NewColorAnimation(&c1.BorderColor, color.OrangeRed, time.Second)
	c1bcolor.AutoReverse = true
	c1bcolor.RepeatCount = AnimationRepeatForever

	// c1fcolor := NewColorAnimation(&c1.FillColor, color.OrangeRed.Alpha(0.4), time.Second)
	// c1fcolor.AutoReverse = true
	// c1fcolor.RepeatCount = AnimationRepeatForever

	c2 := NewEllipse(c2Pos, c2Size1, color.Lime)

	c2pos := NewPathAnimation(&c2.Pos, FullCirclePathB, c2PosSize, 4*time.Second)
	c2pos.RepeatCount = AnimationRepeatForever
	c2pos.Curve = AnimationLinear

	// c2angle := NewFloatAnimation(&c2.Angle, ConstFloat(2*math.Pi), 4*time.Second)
	// c2angle.RepeatCount = AnimationRepeatForever
	// c2angle.Curve = AnimationLinear

	c2size := NewSizeAnimation(&c2.Size, c2Size2, time.Second)
	c2size.AutoReverse = true
	c2size.RepeatCount = AnimationRepeatForever

	c2color := NewPaletteAnimation(&c2.BorderColor, pal, 2*time.Second)
	c2color.RepeatCount = AnimationRepeatForever
	c2color.Curve = AnimationLinear

	ctrl.Add(c2, c1)

	c1pos.Start()
	c1size.Start()
	c1bcolor.Start()
	// c1fcolor.Start()
	c2pos.Start()
	// c2angle.Start()
	c2size.Start()
	c2color.Start()
}

func CircleAnimation(ctrl *Controller) {
	c1Pos1 := ConvertPos(geom.Point{16.5, 4.5})
	c1Pos2 := ConvertPos(geom.Point{2.5, 4.5})
	c1Size1 := ConvertSize(geom.Point{10.0, 10.0})
	c1Size2 := ConvertSize(geom.Point{3.0, 3.0})

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

func PushingRectangles(ctrl *Controller) {
	r1Pos1 := ConvertPos(geom.Point{0.5, 4.5})
	r1Pos2 := ConvertPos(geom.Point{8.5, 4.5})
	r2Pos1 := ConvertPos(geom.Point{18.5, 4.5})
	r2Pos2 := ConvertPos(geom.Point{10.5, 4.5})
	rSize1 := ConvertSize(geom.Point{17.0, 1.0})
	rSize2 := ConvertSize(geom.Point{1.0, 9.0})
	duration := 3 * time.Second / 2

	r1 := NewRectangle(r1Pos1, rSize2, color.Crimson)

	a1Pos := NewPositionAnimation(&r1.Pos, r1Pos2, duration)
	a1Pos.AutoReverse = true
	a1Pos.RepeatCount = AnimationRepeatForever

	a1Size := NewSizeAnimation(&r1.Size, rSize1, duration)
	a1Size.AutoReverse = true
	a1Size.RepeatCount = AnimationRepeatForever

	a1Color := NewColorAnimation(&r1.BorderColor, color.GreenYellow, duration)
	a1Color.AutoReverse = true
	a1Color.RepeatCount = AnimationRepeatForever

	r2 := NewRectangle(r2Pos2, rSize1, color.SkyBlue)

	a2Pos := NewPositionAnimation(&r2.Pos, r2Pos1, duration)
	a2Pos.AutoReverse = true
	a2Pos.RepeatCount = AnimationRepeatForever

	a2Size := NewSizeAnimation(&r2.Size, rSize2, duration)
	a2Size.AutoReverse = true
	a2Size.RepeatCount = AnimationRepeatForever

	a2Color := NewColorAnimation(&r2.BorderColor, color.Crimson, duration)
	a2Color.AutoReverse = true
	a2Color.RepeatCount = AnimationRepeatForever

	ctrl.Add(r1, r2)
	a1Pos.Start()
	a1Size.Start()
	a1Color.Start()
	a2Pos.Start()
	a2Size.Start()
	a2Color.Start()
}

func pixIdx(x, y int) int {
	return y*20 + x
}

func pixCoord(idx int) (x, y int) {
	return idx % 20, idx / 20
}

func GlowingPixels(ctrl *Controller) {
	var pixField []*Pixel

	pixField = make([]*Pixel, 10*20)

	for _, idx := range rand.Perm(200) {
		x, y := pixCoord(idx)
		pos := ConvertPos(geom.Point{float64(x), float64(y)})

		// Farbschema 1: alle Pixel haben zufaellige Farben aus dem
		// gesamten Farbvorrat.
		// idx := rand.IntN(len(color.Names))
		// col := color.Map[color.Names[idx]]

		// Farbschema 2: jeweils zwei nebeneinanderliegende Spalten
		// erhalten zufaellig gewaehlte Farben, die jedoch aus der gleichen
		// Farbgruppe stammen.
		// grpIdx := color.ColorGroup(x/2) % color.NumColorGroups
		// col := color.RandGroupColor(grpIdx)

		// Farbschema 3: alle Pixel haben die gleiche Farbe. Damit laesst
		// sich die Farbanimation besonders gut beobachten. Die Animation
		// verringert den Alpha-Wert (laesst die Farben etwas transparenter
		// werden), von daher sind helle Farben etwas geeigneter.
		// col := color.DodgerBlue

		// Farbschema 4: es werden zwei Farben ausgewaehlt und jedes Pixel
		// erhaelt eine zufaellige Interpolation zwischen diesen beiden
		// Farben.
		col := color.DimGray.Interpolate(color.DarkGray, rand.Float64())
		// col := color.SeaGreen.Interpolate(color.SteelBlue, rand.Float64()).Dark(0.2 * rand.Float64())
		// col := color.BlueViolet.Interpolate(color.DarkMagenta, rand.Float64())
		// col := color.Purple.Interpolate(color.Indigo, rand.Float64())

		pix := NewPixel(pos, col)
		pixField[idx] = pix
		ctrl.Add(pix)
		// col2 := col.Alpha(0.2 + 0.3*rand.Float64())
		dur := time.Second + rand.N(400*time.Millisecond)
		aAlpha := NewFloatAnimation(&pix.Alpha, 0.5, dur)
		aAlpha.AutoReverse = true
		aAlpha.RepeatCount = AnimationRepeatForever
		aAlpha.Start()
	}

	aGrpPurple := NewGroup(0)
	aGrpYellow := NewGroup(0)
	aGrpGreen := NewGroup(0)

	for _, pix := range pixField {
		aColor := NewColorAnimation(&pix.Color, color.MediumPurple.Interpolate(color.Fuchsia, rand.Float64()), 2*time.Second)
		aColor.AutoReverse = true
		aGrpPurple.Add(aColor)
		aColor = NewColorAnimation(&pix.Color, color.Gold.Interpolate(color.Khaki, rand.Float64()), 2*time.Second)
		aColor.AutoReverse = true
		aGrpYellow.Add(aColor)
		aColor = NewColorAnimation(&pix.Color, color.GreenYellow.Interpolate(color.LightSeaGreen, rand.Float64()), 2*time.Second)
		aColor.AutoReverse = true
		aGrpGreen.Add(aColor)
	}

	aTimel := NewTimeline(30 * time.Second)
	aTimel.Add(5*time.Second, aGrpPurple)
	aTimel.Add(15*time.Second, aGrpYellow)
	aTimel.Add(25*time.Second, aGrpGreen)
	aTimel.RepeatCount = AnimationRepeatForever

	aTimel.Start()
}

//----------------------------------------------------------------------------

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

//----------------------------------------------------------------------------

type sceneRecord struct {
	name string
	fnc  func(ctrl *Controller)
}

func main() {
	var sceneId int
	sceneList := []sceneRecord{
		{"Group test", GroupTest},
		{"Sequence test", SequenceTest},
		{"Timeline test", TimelineTest},
		{"Path test", PathTest},
		{"PolygonPath Test", PolygonPathTest},
		{"Random walk", RandomWalk},
		{"Piiiiixels", Piiiiixels},
		{"Circling Circles", CirclingCircles},
		{"Chasing Circles", ChasingCircles},
		{"Circle Animation", CircleAnimation},
		{"Pushing Rectangles", PushingRectangles},
		{"Glowing Pixels", GlowingPixels},
	}

	pixCtrl := ledgrid.NewNetPixelClient(pixelHost, pixelPort)
	pixCtrl.SetGamma(gammaValue, gammaValue, gammaValue)
	pixCtrl.SetMaxBright(255, 255, 255)

	ledGrid := ledgrid.NewLedGrid(gridSize)
	ctrl := NewController(pixCtrl, ledGrid)

	go SignalHandler()

	for {
		fmt.Printf("Choose an animation:\n")
		fmt.Printf("----------------------------\n")
		for i, scene := range sceneList {
			if i == sceneId-1 {
				fmt.Printf("> ")
			} else {
				fmt.Printf("  ")
			}
			fmt.Printf("[%2d] %s\n", i+1, scene.name)
		}
		fmt.Printf("----------------------------\n")
		fmt.Printf("Enter a number (0: quit): ")
		fmt.Scanf("%d", &sceneId)
		if sceneId == 0 {
			break
		}
		if sceneId < 1 || sceneId > len(sceneList) {
			continue
		}
		ctrl.DelAllAnim()
		ctrl.DelAll()
		sceneList[sceneId-1].fnc(ctrl)
	}

	ctrl.Stop()
	ledGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(ledGrid)
	pixCtrl.Close()

}
