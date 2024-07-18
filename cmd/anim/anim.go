package main

import (
	"flag"
	"fmt"
	"image"
	"math"
	"math/rand/v2"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colornames"
)

var (
	width            = 30
	height           = 10
	gridSize         = image.Point{width, height}
	pixelHost        = "raspi-3"
	pixelPort   uint = 5333
	gammaValue       = 3.0
	refreshRate      = 30 * time.Millisecond
	backAlpha        = 1.0

	AnimCtrl Animator
)

//----------------------------------------------------------------------------

func RegularPolygonTest(c *Canvas) {
    posList := []geom.Point{
        ConvertPos(geom.Point{-5.5, 4.5}),
        ConvertPos(geom.Point{34.5, 4.5}),
    }
    posCenter := ConvertPos(geom.Point{14.5, 4.5})
    smallSize := ConvertSize(geom.Point{9.0, 9.0})
    largeSize := ConvertSize(geom.Point{80.0, 80.0})

    polyList := make([]*RegularPolygon, 9)

    aSeq := NewSequence(0.0)
    for n := 3; n <= 8; n++ {
        col := colornames.RandColor()
        polyList[n] = NewRegularPolygon(n, posList[n%2], smallSize, col)
        c.Add(polyList[n])
        dur := 2*time.Second + rand.N(time.Second)
        sign := []float64{+1.0, -1.0}[n%2]
        angle := sign * 2 * math.Pi
        animPos := NewPositionAnimation(&polyList[n].Pos, posCenter, dur)
        animAngle := NewFloatAnimation(&polyList[n].Angle, angle, dur)
        animSize := NewSizeAnimation(&polyList[n].Size, largeSize, 4*time.Second)
        animSize.Cont = true
        animFade := NewColorAnimation(&polyList[n].BorderColor, col.Alpha(0.0), time.Second)

        aGrp := NewGroup(dur, animPos, animAngle)
        aObjSeq := NewSequence(0, aGrp, animSize, animFade)
        aSeq.Add(aObjSeq)
    }
    aSeq.RepeatCount = AnimationRepeatForever
    aSeq.Start()
}

func GroupTest(ctrl *Canvas) {
	ctrl.Stop()

	rPos1 := ConvertPos(geom.Point{4.5, 4.5})
	rPos2 := ConvertPos(geom.Point{27.5, 4.5})
	rSize1 := ConvertSize(geom.Point{7.0, 7.0})
	rSize2 := ConvertSize(geom.Point{1.0, 1.0})

	r := NewRectangle(rPos1, rSize1, colornames.SkyBlue)

	aPos := NewPositionAnimation(&r.Pos, rPos2, time.Second)
	aPos.AutoReverse = true
	aSize := NewSizeAnimation(&r.Size, rSize2, time.Second)
	aSize.AutoReverse = true
	aColor := NewColorAnimation(&r.BorderColor, colornames.GreenYellow, time.Second)
	aColor.AutoReverse = true
	aAngle := NewFloatAnimation(&r.Angle, math.Pi, time.Second)
	aAngle.AutoReverse = true

	aGroup := NewGroup(4*time.Second, aPos, aSize, aColor, aAngle)
	aGroup.RepeatCount = AnimationRepeatForever

	ctrl.Add(r)

	ctrl.Save("gobs/GroupTest.gob")
	ctrl.Continue()
	aGroup.Start()
}

func ReadGroupTest(ctrl *Canvas) {
	ctrl.Stop()

	ctrl.Load("gobs/GroupTest.gob")
	ctrl.Continue()

	// fh, err := os.Open("AnimationProgram.gob")
	// if err != nil {
	// 	log.Fatalf("Couldn't create file: %v", err)
	// }
	// gobDecoder := gob.NewDecoder(fh)
	// err = gobDecoder.Decode(&c)
	// if err != nil {
	// 	log.Fatalf("Couldn't decode data: %v", err)
	// }
	// fh.Close()

	// log.Printf("Controller : %+v\n", c)
	// log.Printf("ObjList[0] : (%T) %+v\n", c.ObjList[0], c.ObjList[0])
	// log.Printf("AnimList[0]: (%T) %+v\n", c.AnimList[0], c.AnimList[0])

	// ctrl.Continue()
}

func SequenceTest(ctrl *Canvas) {
	ctrl.Stop()

	rPos := ConvertPos(geom.Point{14.5, 4.5})
	rSize1 := ConvertSize(geom.Point{29.0, 9.0})
	rSize2 := ConvertSize(geom.Point{1.0, 1.0})

	r := NewRectangle(rPos, rSize1, colornames.SkyBlue)

	aSize1 := NewSizeAnimation(&r.Size, rSize2, time.Second)
	aColor1 := NewColorAnimation(&r.BorderColor, colornames.OrangeRed, time.Second/2)
	aColor1.AutoReverse = true
	aColor1.RepeatCount = 2
	aColor2 := NewColorAnimation(&r.BorderColor, colornames.OrangeRed, time.Second/2)
	aSize2 := NewSizeAnimation(&r.Size, rSize1, time.Second)
	aSize2.Cont = true
	aColor3 := NewColorAnimation(&r.BorderColor, colornames.SkyBlue, 3*time.Second/2)
	aColor3.Cont = true

	aSeq := NewSequence(8*time.Second, aSize1, aColor1, aColor2, aSize2, aColor3)
	aSeq.RepeatCount = AnimationRepeatForever

	ctrl.Add(r)

	ctrl.Save("gobs/SequenceTest.gob")
	ctrl.Continue()
	aSeq.Start()
}

func TimelineTest(ctrl *Canvas) {
	ctrl.Stop()

	r3Pos1 := ConvertPos(geom.Point{6.5, 4.5})
	r3Size1 := ConvertSize(geom.Point{9.0, 5.0})
	r4Pos1 := ConvertPos(geom.Point{21.5, 4.5})
	r4Size1 := ConvertSize(geom.Point{9.0, 5.0})

	r3 := NewRectangle(r3Pos1, r3Size1, colornames.GreenYellow)
	r4 := NewRectangle(r4Pos1, r4Size1, colornames.SkyBlue)

	aAngle1 := NewFloatAnimation(&r3.Angle, math.Pi, 2*time.Second)
	aAngle2 := NewFloatAnimation(&r4.Angle, -math.Pi, 2*time.Second)
	aAngle3 := NewFloatAnimation(&r3.Angle, -math.Pi, time.Second)
	aAngle4 := NewFloatAnimation(&r4.Angle, math.Pi, time.Second)

	r3Color := NewColorAnimation(&r3.BorderColor, colornames.OrangeRed, 200*time.Millisecond)
	r3Color.AutoReverse = true
	r4Color := NewColorAnimation(&r4.BorderColor, colornames.OrangeRed, 200*time.Millisecond)
	r4Color.AutoReverse = true

	tl := NewTimeline(5 * time.Second)
	tl.RepeatCount = AnimationRepeatForever
	tl.Add(0*time.Millisecond, aAngle1, aAngle2)
	tl.Add(2200*time.Millisecond, r3Color)
	tl.Add(2400*time.Millisecond, r4Color)
	tl.Add(2600*time.Millisecond, aAngle3)
	tl.Add(2900*time.Millisecond, aAngle4)

	ctrl.Add(r3, r4)

	ctrl.Save("gobs/TimelineTest.gob")
	ctrl.Continue()
	tl.Start()

}

func PathTest(ctrl *Canvas) {
	ctrl.Stop()

	pos1 := ConvertPos(geom.Point{1.0, 4.0})
	pos2 := ConvertPos(geom.Point{14.0, 1.0})
	pos3 := ConvertPos(geom.Point{27.0, 4.0})
	pos4 := ConvertPos(geom.Point{14.0, 7.0})
	cSize := ConvertSize(geom.Point{2.0, 2.0})

	c1 := NewEllipse(pos1, cSize, colornames.OrangeRed)
	c2 := NewEllipse(pos2, cSize, colornames.MediumSeaGreen)
	c3 := NewEllipse(pos3, cSize, colornames.SkyBlue)
	c4 := NewEllipse(pos4, cSize, colornames.Gold)

	c1Path := NewPathAnimation(&c1.Pos, HalfCirclePathB, ConvertSize(geom.Point{26.0, 3.0}), 2*time.Second)
	c1Path.AutoReverse = true
	// c1Path.Cont = true
	c3Path := NewPathAnimation(&c3.Pos, HalfCirclePathB, ConvertSize(geom.Point{-26.0, -3.0}), 2*time.Second)
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

	ctrl.Save("gobs/PathTest.gob")
	ctrl.Continue()
	aGrp.Start()
}

func PolygonPathTest(ctrl *Canvas) {
	ctrl.Stop()

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

	c1 := NewEllipse(cPos, cSize, colornames.GreenYellow)

	aPath := NewPathAnimation(&c1.Pos, polyPath.RelPoint,
		ConvertSize(geom.Point{27, 7}), 7*time.Second)
	aPath.AutoReverse = true
	aPath.RepeatCount = AnimationRepeatForever

	ctrl.Add(c1)
	ctrl.Save("gobs/PolygonPathTest.gob")
	ctrl.Continue()

	aPath.Start()
	// fmt.Printf("PolygonPath: %+v\n", polyPath)
	// fmt.Printf("Point at 0.0: %v\n", polyPath.RelPoint(0.0))
	// fmt.Printf("Point at 0.25: %v\n", polyPath.RelPoint(0.25))
	// fmt.Printf("Point at 0.5: %v\n", polyPath.RelPoint(0.5))
	// fmt.Printf("Point at 0.75: %v\n", polyPath.RelPoint(0.75))
	// fmt.Printf("Point at 1.0: %v\n", polyPath.RelPoint(1.0))
}

func RandomWalk(ctrl *Canvas) {
	ctrl.Stop()

	rect := geom.Rectangle{Min: ConvertPos(geom.Point{1.0, 1.0}), Max: ConvertPos(geom.Point{28.0, 8.0})}
	pos1 := ConvertPos(geom.Point{1.0, 1.0})
	pos2 := ConvertPos(geom.Point{18.0, 8.0})
	cSize := ConvertSize(geom.Point{3.0, 3.0})

	c1 := NewEllipse(pos1, cSize, colornames.SkyBlue)
	c2 := NewEllipse(pos2, cSize, colornames.GreenYellow)

	aPos1 := NewPositionAnimation(&c1.Pos, geom.Point{}, 1300*time.Millisecond)
	aPos1.Curve = AnimationLinear
	aPos1.Cont = true
	aPos1.ValFunc = RandPoint(rect)
	aPos1.RepeatCount = AnimationRepeatForever

	aPos2 := NewPositionAnimation(&c2.Pos, geom.Point{}, 901*time.Millisecond)
	aPos1.Curve = AnimationLinear
	aPos2.Cont = true
	aPos2.ValFunc = RandPoint(rect)
	aPos2.RepeatCount = AnimationRepeatForever

	ctrl.Add(c1, c2)
	ctrl.Save("gobs/RandomWalk.gob")
	ctrl.Continue()

	aPos1.Start()
	aPos2.Start()
}

func Piiiiixels(ctrl *Canvas) {
	p1Pos1 := ConvertPos(geom.Point{1.0, 1.0})
	// p1Pos2 := ConvertPos(geom.Point{18.0, 1.0})
	p2Pos1 := ConvertPos(geom.Point{9.0, 1.0})
	p2Range := ConvertSize(geom.Point{18.0, 7.0})
	// p2Pos2 := ConvertPos(geom.Point{1.0, 8.0})

	p1 := NewPixel(p1Pos1, colornames.OrangeRed)
	p2 := NewPixel(p2Pos1, colornames.Lime)

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

func CirclingCircles(ctrl *Canvas) {
	pos1 := ConvertPos(geom.Point{2.0, 2.0})
	pos2 := ConvertPos(geom.Point{7.0, 7.0})
	pos3 := ConvertPos(geom.Point{12.0, 2.0})
	pos4 := ConvertPos(geom.Point{17.0, 7.0})
	cSize := ConvertSize(geom.Point{2.0, 2.0})

	c1 := NewEllipse(pos1, cSize, colornames.OrangeRed)
	c2 := NewEllipse(pos2, cSize, colornames.MediumSeaGreen)
	c3 := NewEllipse(pos3, cSize, colornames.SkyBlue)
	c4 := NewEllipse(pos4, cSize, colornames.Gold)

	c1Path1 := NewPathAnimation(&c1.Pos, QuarterCirclePathA, ConvertSize(geom.Point{5.0, 5.0}), time.Second)
	c1Path1.Cont = true
	c1Path2 := NewPathAnimation(&c1.Pos, QuarterCirclePathA, ConvertSize(geom.Point{5.0, -5.0}), time.Second)
	c1Path2.Cont = true
	c1Path3 := NewPathAnimation(&c1.Pos, QuarterCirclePathA, ConvertSize(geom.Point{5.0, 5.0}), time.Second)
	c1Path3.Cont = true

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
	c4Path3 := NewPathAnimation(&c4.Pos, QuarterCirclePathA, ConvertSize(geom.Point{-5.0, -5.0}), time.Second)
	c4Path3.Cont = true

	aGrp1 := NewGroup(0, c1Path1, c2Path1)
	aGrp2 := NewGroup(0, c1Path2, c3Path1)
	aGrp3 := NewGroup(0, c1Path3, c4Path1)
	aGrp4 := NewGroup(0, c4Path2, c3Path2)
	aGrp5 := NewGroup(0, c4Path3, c2Path2)
	aSeq := NewSequence(0, aGrp1, aGrp2, aGrp3, aGrp4, aGrp5)

	ctrl.Add(c1, c2, c3, c4)
	aSeq.Start()
}

func ChasingCircles(ctrl *Canvas) {
	c1Pos1 := ConvertPos(geom.Point{16.5, 4.5})
	c1Size1 := ConvertSize(geom.Point{10.0, 10.0})
	c1Size2 := ConvertSize(geom.Point{3.0, 3.0})
	c1PosSize := ConvertSize(geom.Point{-14.0, -5.0})
	c2Pos := ConvertPos(geom.Point{2.5, 4.5})
	c2Size1 := ConvertSize(geom.Point{5.0, 5.0})
	c2Size2 := ConvertSize(geom.Point{3.0, 3.0})
	c2PosSize := ConvertSize(geom.Point{14.0, 7.0})

	pal := ledgrid.NewGradientPaletteByList("Palette", true,
		ledgrid.LedColorModel.Convert(colornames.DeepSkyBlue).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(colornames.Lime).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(colornames.Teal).(ledgrid.LedColor),
		ledgrid.LedColorModel.Convert(colornames.SkyBlue).(ledgrid.LedColor),
		// ledgrid.ColorMap["DeepSkyBlue"].Color(0),
		// ledgrid.ColorMap["Lime"].Color(0),
		// ledgrid.ColorMap["Teal"].Color(0),
		// ledgrid.ColorMap["SkyBlue"].Color(0),
	)

	c1 := NewEllipse(c1Pos1, c1Size1, colornames.Gold)

	c1pos := NewPathAnimation(&c1.Pos, FullCirclePathB, c1PosSize, 4*time.Second)
	c1pos.RepeatCount = AnimationRepeatForever
	c1pos.Curve = AnimationLinear

	// c1pos := NewPositionAnimation(&c1.Pos, c1Pos2, 2*time.Second)
	// c1pos.AutoReverse = true
	// c1pos.RepeatCount = AnimationRepeatForever

	c1size := NewSizeAnimation(&c1.Size, c1Size2, time.Second)
	c1size.AutoReverse = true
	c1size.RepeatCount = AnimationRepeatForever

	c1bcolor := NewColorAnimation(&c1.BorderColor, colornames.OrangeRed, time.Second)
	c1bcolor.AutoReverse = true
	c1bcolor.RepeatCount = AnimationRepeatForever

	// c1fcolor := NewColorAnimation(&c1.FillColor, colornames.OrangeRed.Alpha(0.4), time.Second)
	// c1fcolor.AutoReverse = true
	// c1fcolor.RepeatCount = AnimationRepeatForever

	c2 := NewEllipse(c2Pos, c2Size1, colornames.Lime)

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

func CircleAnimation(ctrl *Canvas) {
	c1Pos1 := ConvertPos(geom.Point{1.5, 4.5})
	// c1Pos2 := ConvertPos(geom.Point{14.5, 4.5})
	c1Pos3 := ConvertPos(geom.Point{27.5, 4.5})

	c1Size1 := ConvertSize(geom.Point{3.0, 3.0})
	c1Size2 := ConvertSize(geom.Point{9.0, 9.0})

	// c1Pos1 := ConvertPos(geom.Point{26.5, 4.5})
	// c1Pos2 := ConvertPos(geom.Point{2.5, 4.5})

	c1 := NewEllipse(c1Pos1, c1Size1, colornames.OrangeRed)

	c1pos := NewPositionAnimation(&c1.Pos, c1Pos3, 2*time.Second)
	c1pos.AutoReverse = true
	c1pos.RepeatCount = AnimationRepeatForever

	c1radius := NewSizeAnimation(&c1.Size, c1Size2, time.Second)
	c1radius.AutoReverse = true
	c1radius.RepeatCount = AnimationRepeatForever

	c1color := NewColorAnimation(&c1.BorderColor, colornames.Gold, 4*time.Second)
	c1color.AutoReverse = true
	c1color.RepeatCount = AnimationRepeatForever

	ctrl.Add(c1)

	c1pos.Start()
	c1radius.Start()
	c1color.Start()
}

func PushingRectangles(ctrl *Canvas) {
	r1Pos1 := ConvertPos(geom.Point{0.5, 4.5})
	r1Pos2 := ConvertPos(geom.Point{13.5, 4.5})
	r2Pos1 := ConvertPos(geom.Point{28.5, 4.5})
	r2Pos2 := ConvertPos(geom.Point{15.5, 4.5})
	rSize1 := ConvertSize(geom.Point{27.0, 1.0})
	rSize2 := ConvertSize(geom.Point{1.0, 9.0})
	duration := 3 * time.Second / 2

	r1 := NewRectangle(r1Pos1, rSize2, colornames.Crimson)

	a1Pos := NewPositionAnimation(&r1.Pos, r1Pos2, duration)
	a1Pos.AutoReverse = true
	a1Pos.RepeatCount = AnimationRepeatForever

	a1Size := NewSizeAnimation(&r1.Size, rSize1, duration)
	a1Size.AutoReverse = true
	a1Size.RepeatCount = AnimationRepeatForever

	a1Color := NewColorAnimation(&r1.BorderColor, colornames.GreenYellow, duration)
	a1Color.AutoReverse = true
	a1Color.RepeatCount = AnimationRepeatForever

	r2 := NewRectangle(r2Pos2, rSize1, colornames.SkyBlue)

	a2Pos := NewPositionAnimation(&r2.Pos, r2Pos1, duration)
	a2Pos.AutoReverse = true
	a2Pos.RepeatCount = AnimationRepeatForever

	a2Size := NewSizeAnimation(&r2.Size, rSize2, duration)
	a2Size.AutoReverse = true
	a2Size.RepeatCount = AnimationRepeatForever

	a2Color := NewColorAnimation(&r2.BorderColor, colornames.Crimson, duration)
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
	return y*30 + x
}

func pixCoord(idx int) (x, y int) {
	return idx % 30, idx / 30
}

func GlowingPixels(ctrl *Canvas) {
	var pixField []*Pixel

	ctrl.Stop()

	pixField = make([]*Pixel, 10*30)

	for _, idx := range rand.Perm(300) {
		x, y := pixCoord(idx)
		pos := ConvertPos(geom.Point{float64(x), float64(y)})

		// Farbschema 1: alle Pixel haben zufaellige Farben aus dem
		// gesamten Farbvorrat.
		// idx := rand.IntN(len(colornames.Names))
		// col := colornames.Map[colornames.Names[idx]]

		// Farbschema 2: jeweils zwei nebeneinanderliegende Spalten
		// erhalten zufaellig gewaehlte Farben, die jedoch aus der gleichen
		// Farbgruppe stammen.
		// grpIdx := ledgrid.LedColorGroup(x/2) % colornames.NumColorGroups
		// col := colornames.RandGroupColor(grpIdx)

		// Farbschema 3: alle Pixel haben die gleiche Farbe. Damit laesst
		// sich die Farbanimation besonders gut beobachten. Die Animation
		// verringert den Alpha-Wert (laesst die Farben etwas transparenter
		// werden), von daher sind helle Farben etwas geeigneter.
		// col := colornames.DodgerBlue

		// Farbschema 4: es werden zwei Farben ausgewaehlt und jedes Pixel
		// erhaelt eine zufaellige Interpolation zwischen diesen beiden
		// Farben.
		col := colornames.DimGray.Interpolate(colornames.DarkGray, rand.Float64())
		// col := colornames.SeaGreen.Interpolate(colornames.SteelBlue, rand.Float64()).Dark(0.2 * rand.Float64())
		// col := colornames.BlueViolet.Interpolate(colornames.DarkMagenta, rand.Float64())
		// col := colornames.Purple.Interpolate(colornames.Indigo, rand.Float64())

		pix := NewPixel(pos, col)
		pixField[idx] = pix
		ctrl.Add(pix)

		dur := time.Second + rand.N(400*time.Millisecond)
		aAlpha := NewFloatAnimation(&pix.Alpha, 0.5, dur)
		aAlpha.AutoReverse = true
		aAlpha.RepeatCount = AnimationRepeatForever
		aAlpha.Start()
	}

	aGrpPurple := NewGroup(0)
	aGrpYellow := NewGroup(0)
	aGrpGreen := NewGroup(0)
	// aGrpGray := NewGroup(0)

	for _, pix := range pixField {
		aColor := NewColorAnimation(&pix.Color, colornames.MediumPurple.Interpolate(colornames.Fuchsia, rand.Float64()), 3*time.Second)
		// aColor.Cont = true
		aColor.AutoReverse = true
		aGrpPurple.Add(aColor)

		aColor = NewColorAnimation(&pix.Color, colornames.Gold.Interpolate(colornames.Khaki, rand.Float64()), 3*time.Second)
		// aColor.Cont = true
		aColor.AutoReverse = true
		aGrpYellow.Add(aColor)

		aColor = NewColorAnimation(&pix.Color, colornames.GreenYellow.Interpolate(colornames.LightSeaGreen, rand.Float64()), 3*time.Second)
		// aColor.Cont = true
		aColor.AutoReverse = true
		aGrpGreen.Add(aColor)

		// aColor = NewColorAnimation(&pix.Color, colornames.DimGray.Interpolate(colornames.DarkGray, rand.Float64()), 3*time.Second)
		// aColor.Cont = true
		// aGrpGray.Add(aColor)
	}

	aTimel := NewTimeline(40 * time.Second)
	aTimel.Add(10*time.Second, aGrpPurple)
	aTimel.Add(20*time.Second, aGrpYellow)
	aTimel.Add(30*time.Second, aGrpGreen)
	// aTimel.Add(40*time.Second, aGrpGray)
	aTimel.RepeatCount = AnimationRepeatForever

	ctrl.Save("gobs/GlowingPixels.gob")
	ctrl.Continue()

	aTimel.Start()
}

func MovingText(c *Canvas) {
	tPos := ConvertPos(geom.Point{14.5, 4.5})

	t := NewText(tPos, "Beni", colornames.OrangeRed)

	aAngle := NewFloatAnimation(&t.Angle, 2*math.Pi, 4*time.Second)
	aAngle.RepeatCount = AnimationRepeatForever

	c.Add(t)
	aAngle.Start()
}

/*
Canvas statistics:
  animation: 1460 calls; 574.805478ms in total; 393.702µs per call
  painting : 1460 calls; 11.590908618s in total; 7.938978ms per call
  sending  : 1460 calls; 2.62475604s in total; 1.797778ms per call

Grid statistics:
  animation: 1291 calls; 534.033261ms in total; 413.658µs per call
  painting : 1291 calls; 18.301573ms in total; 14.176µs per call
  sending  : 1291 calls; 40.68729ms in total; 31.516µs per call
*/

func GridTest(g *Grid) {
	for _, cell := range g.Cells {
		cell.Color = colornames.DimGray.Interpolate(colornames.DarkGray, rand.Float64())
		dur := time.Second + rand.N(400*time.Millisecond)
		aAlpha := NewFloatAnimation(&cell.Alpha, 0.7, dur)
		aAlpha.AutoReverse = true
		aAlpha.RepeatCount = AnimationRepeatForever
		aAlpha.Start()
	}

    	aGrpPurple := NewGroup(0)
	aGrpYellow := NewGroup(0)
	aGrpGreen := NewGroup(0)

	for _, cell := range g.Cells {
		aColor := NewColorAnimation(&cell.Color, colornames.MediumPurple.Interpolate(colornames.Fuchsia, rand.Float64()), 3*time.Second)
		aColor.AutoReverse = true
		aGrpPurple.Add(aColor)

		aColor = NewColorAnimation(&cell.Color, colornames.Gold.Interpolate(colornames.Khaki, rand.Float64()), 3*time.Second)
		aColor.AutoReverse = true
		aGrpYellow.Add(aColor)

		aColor = NewColorAnimation(&cell.Color, colornames.GreenYellow.Interpolate(colornames.LightSeaGreen, rand.Float64()), 3*time.Second)
		aColor.AutoReverse = true
		aGrpGreen.Add(aColor)
	}

	aTimel := NewTimeline(40 * time.Second)
	aTimel.Add(10*time.Second, aGrpPurple)
	aTimel.Add(20*time.Second, aGrpYellow)
	aTimel.Add(30*time.Second, aGrpGreen)
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

type canvasSceneRecord struct {
	name string
	fnc  func(canvas *Canvas)
}

type gridSceneRecord struct {
	name string
	fnc  func(grid *Grid)
}

func main() {
	var sceneId int
	var input string
	var runInteractive bool

	flag.StringVar(&input, "scene", input, "no menu: direct play")
	flag.BoolVar(&doLog, "log", doLog, "enable logging")
	flag.Parse()

	if len(input) > 0 {
		runInteractive = false
	} else {
		runInteractive = true
	}

	canvasSceneList := []canvasSceneRecord{
        {"(Regular) Polygon test", RegularPolygonTest},
		{"Group test", GroupTest},
		// {"Group test (from saved program)", ReadGroupTest},
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
		{"Moving Text", MovingText},
	}

	gridSceneList := []gridSceneRecord{
		{"Grid Test", GridTest},
	}

	pixCtrl := ledgrid.NewNetPixelClient(pixelHost, pixelPort)
	pixCtrl.SetGamma(gammaValue, gammaValue, gammaValue)
	pixCtrl.SetMaxBright(255, 255, 255)

	ledGrid := ledgrid.NewLedGrid(gridSize)
	canvas := NewCanvas(pixCtrl, ledGrid)
	canvas.Stop()
	grid := NewGrid(pixCtrl, ledGrid)
	grid.Stop()

	if runInteractive {
		sceneId = -1
		for {
			fmt.Printf("Animations:\n")
			fmt.Printf("---------------------------------------\n")
			for i, scene := range canvasSceneList {
				if i == sceneId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%c] %s\n", 'a'+i, scene.name)
			}
			fmt.Printf("---------------------------------------\n")
			for i, scene := range gridSceneList {
				if i == sceneId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%c] %s\n", 'A'+i, scene.name)
			}
			fmt.Printf("---------------------------------------\n")

			fmt.Printf("Enter a character (or '0' for quit): ")
			fmt.Scanf("%s\n", &input)
			if input[0] == '0' {
				break
			}
			if input[0] >= 'a' && input[0] <= 'z' {
				sceneId = int(input[0] - 'a')
				if sceneId < 0 || sceneId >= len(canvasSceneList) {
					continue
				}
				if AnimCtrl != nil {
					AnimCtrl.Stop()
					AnimCtrl.DelAllAnim()
				}
				AnimCtrl = canvas
				AnimCtrl.Continue()
				canvasSceneList[sceneId].fnc(canvas)
			}
			if input[0] >= 'A' && input[0] <= 'Z' {
				sceneId = int(input[0] - 'A')
				if sceneId < 0 || sceneId >= len(gridSceneList) {
					continue
				}
				if AnimCtrl != nil {
					AnimCtrl.Stop()
					AnimCtrl.DelAllAnim()
				}
				AnimCtrl = grid
				AnimCtrl.Continue()
				gridSceneList[sceneId].fnc(grid)
			}
		}
	} else {
		if input[0] >= 'a' && input[0] <= 'z' {
			sceneId = int(input[0] - 'a')
			if sceneId >= 0 && sceneId < len(canvasSceneList) {
				AnimCtrl = canvas
				canvasSceneList[sceneId].fnc(canvas)
			}
		}
		if input[0] >= 'A' && input[0] <= 'Z' {
			sceneId = int(input[0] - 'A')
			if sceneId >= 0 && sceneId < len(gridSceneList) {
				AnimCtrl = grid
				gridSceneList[sceneId].fnc(grid)
			}
		}
		AnimCtrl.Continue()
		fmt.Printf("Quit by Ctrl-C\n")
		SignalHandler()
	}

	AnimCtrl.Stop()
	ledGrid.Clear(ledgrid.Black)
	pixCtrl.Draw(ledGrid)
	pixCtrl.Close()

	fmt.Printf("Canvas statistics:\n")
	fmt.Printf("  animation: %v\n", canvas.animWatch)
	fmt.Printf("  painting : %v\n", canvas.paintWatch)
	fmt.Printf("  sending  : %v\n", canvas.sendWatch)
	fmt.Printf("Grid statistics:\n")
	fmt.Printf("  animation: %v\n", grid.animWatch)
	fmt.Printf("  painting : %v\n", grid.paintWatch)
	fmt.Printf("  sending  : %v\n", grid.sendWatch)
}
