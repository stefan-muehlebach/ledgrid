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
	"github.com/stefan-muehlebach/ledgrid/color"
	"gocv.io/x/gocv"
	"golang.org/x/image/math/fixed"
)

const (
	defWidth  = 40
	defHeight = 10
	defHost   = "raspi-3"
	defPort   = 5333
)

var (
	width, height int
	gridSize      image.Point
	backAlpha     = 1.0
	animCtrl      *ledgrid.AnimationController
)

//----------------------------------------------------------------------------

type LedGridProgram interface {
	Name() string
	Run(c *ledgrid.Canvas)
}

func NewLedGridProgram(name string, runFunc func(c *ledgrid.Canvas)) LedGridProgram {
	return &simpleProgram{name, runFunc}
}

type simpleProgram struct {
	name    string
	runFunc func(c *ledgrid.Canvas)
}

func (p *simpleProgram) Name() string {
	return p.name
}

func (p *simpleProgram) Run(c *ledgrid.Canvas) {
	p.runFunc(c)
}

var (
	GroupTest = NewLedGridProgram("Group test",
		func(c *ledgrid.Canvas) {
			rPos1 := geom.Point{5.0, 5.0}
			rPos2 := geom.Point{float64(width) - 2.0, 5.0}
			rSize1 := geom.Point{7.0, 7.0}
			rSize2 := geom.Point{1.0, 1.0}
			rColor1 := color.SkyBlue
			rColor2 := color.GreenYellow

			r := ledgrid.NewRectangle(rPos1, rSize1, rColor1)
			c.Add(r)

			aPos := ledgrid.NewPositionAnimation(&r.Pos, rPos2, time.Second)
			aPos.AutoReverse = true
			aSize := ledgrid.NewSizeAnimation(&r.Size, rSize2, time.Second)
			aSize.AutoReverse = true
			aColor := ledgrid.NewColorAnimation(&r.BorderColor, rColor2, time.Second)
			aColor.AutoReverse = true
			aAngle := ledgrid.NewFloatAnimation(&r.Angle, math.Pi, time.Second)
			aAngle.AutoReverse = true

			aGroup := ledgrid.NewGroup(aPos, aSize, aColor, aAngle)
			aGroup.RepeatCount = ledgrid.AnimationRepeatForever

			aGroup.Start()
		})

	SequenceTest = NewLedGridProgram("Sequence test",
		func(c *ledgrid.Canvas) {
			rPos := geom.NewPointIMG(gridSize).Mul(0.5)
			rSize1 := geom.NewPointIMG(gridSize).SubXY(1, 1)
			rSize2 := geom.Point{5.0, 3.0}

			r := ledgrid.NewRectangle(rPos, rSize1, color.SkyBlue)
			c.Add(r)

			aSize1 := ledgrid.NewSizeAnimation(&r.Size, rSize2, time.Second)
			aColor1 := ledgrid.NewColorAnimation(&r.BorderColor, color.OrangeRed, time.Second/2)
			aColor1.AutoReverse = true
			aColor2 := ledgrid.NewColorAnimation(&r.BorderColor, color.Crimson, time.Second/2)
			aColor2.AutoReverse = true
			aColor3 := ledgrid.NewColorAnimation(&r.BorderColor, color.Coral, time.Second/2)
			aColor3.AutoReverse = true
			aColor4 := ledgrid.NewColorAnimation(&r.BorderColor, color.FireBrick, time.Second/2)
			aSize2 := ledgrid.NewSizeAnimation(&r.Size, rSize1, time.Second)
			aSize2.Cont = true
			aColor5 := ledgrid.NewColorAnimation(&r.BorderColor, color.SkyBlue, time.Second)
			aColor5.Cont = true

			aSeq := ledgrid.NewSequence(aSize1, aColor1, aColor2, aColor3, aColor4, aSize2, aColor5)
			aSeq.RepeatCount = ledgrid.AnimationRepeatForever
			aSeq.Start()
		})

	TimelineTest = NewLedGridProgram("Timeline test",
		func(c *ledgrid.Canvas) {
			r1Pos := geom.Point{6, float64(height) / 2.0}
			r1Size := geom.Point{9.0, 3.0}
			r2Pos := geom.Point{float64(width)/2.0 - 3.0, float64(height) / 2.0}
			r2Size := geom.Point{3.0, 9.0}
			r3Pos := geom.Point{float64(width) - 6, float64(height) / 2.0}
			r3Size := geom.Point{9.0, 3.0}
			r4Pos := geom.Point{float64(width)/2.0 + 3.0, float64(height) / 2.0}
			r4Size := geom.Point{3.0, 9.0}

			r1 := ledgrid.NewRectangle(r1Pos, r1Size, color.GreenYellow)
			r2 := ledgrid.NewRectangle(r2Pos, r2Size, color.Gold)
			r3 := ledgrid.NewRectangle(r3Pos, r3Size, color.SkyBlue)
			r4 := ledgrid.NewRectangle(r4Pos, r4Size, color.Gold)
			c.Add(r1, r3, r2, r4)

			aAngle1 := ledgrid.NewFloatAnimation(&r1.Angle, math.Pi, time.Second)
			aAngle2 := ledgrid.NewFloatAnimation(&r1.Angle, 0.0, time.Second)
			aAngle2.Cont = true

			aColor1 := ledgrid.NewColorAnimation(&r1.BorderColor, color.OrangeRed, 200*time.Millisecond)
			aColor1.AutoReverse = true
			aColor1.RepeatCount = 3
			aColor2 := ledgrid.NewColorAnimation(&r1.BorderColor, color.Purple, 500*time.Millisecond)
			aColor3 := ledgrid.NewColorAnimation(&r1.BorderColor, color.GreenYellow, 500*time.Millisecond)
			aColor3.Cont = true

			aPos1 := ledgrid.NewPositionAnimation(&r1.Pos, r2.Pos.SubXY(r2Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos1.AutoReverse = true

			aAngle3 := ledgrid.NewFloatAnimation(&r3.Angle, -math.Pi, time.Second)
			aAngle4 := ledgrid.NewFloatAnimation(&r3.Angle, 0.0, time.Second)
			aAngle4.Cont = true

			aColor4 := ledgrid.NewColorAnimation(&r3.BorderColor, color.DarkOrange, 200*time.Millisecond)
			aColor4.AutoReverse = true
			aColor4.RepeatCount = 3
			aColor5 := ledgrid.NewColorAnimation(&r3.BorderColor, color.Purple, 500*time.Millisecond)
			aColor6 := ledgrid.NewColorAnimation(&r3.BorderColor, color.SkyBlue, 500*time.Millisecond)
			aColor6.Cont = true

			aPos2 := ledgrid.NewPositionAnimation(&r3.Pos, r4.Pos.AddXY(r4Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos2.AutoReverse = true

			aColor7 := ledgrid.NewColorAnimation(&r2.BorderColor, color.Cornsilk, 500*time.Millisecond)
			aColor7.AutoReverse = true
			aBorder1 := ledgrid.NewFloatAnimation(&r2.BorderWidth, 2.0, 500*time.Millisecond)
			aBorder1.AutoReverse = true

			aColor8 := ledgrid.NewColorAnimation(&r4.BorderColor, color.Cornsilk, 500*time.Millisecond)
			aColor8.AutoReverse = true
			aBorder2 := ledgrid.NewFloatAnimation(&r4.BorderWidth, 2.0, 500*time.Millisecond)
			aBorder2.AutoReverse = true

			tl := ledgrid.NewTimeline(5 * time.Second)
			tl.RepeatCount = ledgrid.AnimationRepeatForever

			// Timeline positions for the first rectangle
			tl.Add(300*time.Millisecond, aColor1)
			tl.Add(1800*time.Millisecond, aAngle1)
			tl.Add(2300*time.Millisecond, aColor2, aPos1)
			tl.Add(2400*time.Millisecond, aColor7, aBorder1)
			tl.Add(2900*time.Millisecond, aAngle2)
			tl.Add(3400*time.Millisecond, aColor3)

			// Timeline positions for the second rectangle
			tl.Add(500*time.Millisecond, aColor4)
			tl.Add(2000*time.Millisecond, aAngle3)
			tl.Add(2500*time.Millisecond, aColor5, aPos2)
			tl.Add(2600*time.Millisecond, aColor8, aBorder2)
			tl.Add(3100*time.Millisecond, aAngle4)
			tl.Add(3600*time.Millisecond, aColor6)

			tl.Start()
		})

	StopContShowHideTest = NewLedGridProgram("Stop/Continue, Show/Hide test",
		func(c *ledgrid.Canvas) {
			rPos1 := geom.Point{5.0, 5.0}
			rPos2 := geom.Point{float64(width) - 5.0, 5.0}
			rSize1 := geom.Point{7.0, 7.0}
			rSize2 := geom.Point{1.0, 1.0}
			rColor1 := color.SkyBlue
			rColor2 := color.GreenYellow

			r := ledgrid.NewRectangle(rPos1, rSize1, rColor1)
			c.Add(r)

			aPos := ledgrid.NewPositionAnimation(&r.Pos, rPos2, 4*time.Second)
			aPos.AutoReverse = true
			aSize := ledgrid.NewSizeAnimation(&r.Size, rSize2, 4*time.Second)
			aSize.AutoReverse = true
			aColor := ledgrid.NewColorAnimation(&r.BorderColor, rColor2, 4*time.Second)
			aColor.AutoReverse = true
			aAngle := ledgrid.NewFloatAnimation(&r.Angle, math.Pi, 4*time.Second)
			aAngle.AutoReverse = true

			aGroup := ledgrid.NewGroup(aPos, aSize, aColor, aAngle)
			aGroup.RepeatCount = ledgrid.AnimationRepeatForever

			// aShowHide := ledgrid.NewShowHideAnimation(r)
			aOnOffPos := ledgrid.NewStopContAnimation(aPos)
			// aOnOffAngle := ledgrid.NewStopContAnimation(aAngle)

			aTimeline := ledgrid.NewTimeline(4 * time.Second)
			// aTimeline.RepeatCount = ledgrid.AnimationRepeatForever

			aTimeline.Add(1000*time.Millisecond, aOnOffPos)
			aTimeline.Add(3000*time.Millisecond, aOnOffPos)

			// aTimeline.Add(500*time.Millisecond, aOnOffPos)
			// aTimeline.Add(1500*time.Millisecond, aOnOffPos)
			// aTimeline.Add(3500*time.Millisecond, aShowHide)
			// aTimeline.Add(3800*time.Millisecond, aShowHide)

			aGroup.Add(aTimeline)

			aGroup.Start()
			// aTimeline.Start()
		})

	PathTest = NewLedGridProgram("Path test",
		func(c *ledgrid.Canvas) {
			duration := 4 * time.Second
			pathA := ledgrid.FullCirclePathA
			pathB := ledgrid.FullCirclePathB

			pos1 := geom.Point{2, float64(height) / 2.0}
			pos2 := geom.Point{float64(width) / 2.0, 2}
			pos3 := geom.Point{float64(width) - 2, float64(height) / 2.0}
			pos4 := geom.Point{float64(width) / 2.0, float64(height) - 2}
			cSize := geom.Point{3.0, 3.0}

			c1 := ledgrid.NewEllipse(pos1, cSize, color.OrangeRed)
			c2 := ledgrid.NewEllipse(pos2, cSize, color.MediumSeaGreen)
			c3 := ledgrid.NewEllipse(pos3, cSize, color.SkyBlue)
			c4 := ledgrid.NewEllipse(pos4, cSize, color.Gold)
			c.Add(c1, c2, c3, c4)

			c1Path := ledgrid.NewPathAnimation(&c1.Pos, pathB, geom.Point{float64(width - 4), 6.0}, duration)
			c1Path.AutoReverse = true
			c3Path := ledgrid.NewPathAnimation(&c3.Pos, pathB, geom.Point{-float64(width - 4), -6.0}, duration)
			c3Path.AutoReverse = true

			c2Path := ledgrid.NewPathAnimation(&c2.Pos, pathA, geom.Point{float64(width) / 3.0, 6.0}, duration)
			c2Path.AutoReverse = true
			c4Path := ledgrid.NewPathAnimation(&c4.Pos, pathA, geom.Point{-float64(width) / 3.0, -6.0}, duration)
			c4Path.AutoReverse = true

			aGrp := ledgrid.NewGroup(c1Path, c2Path, c3Path, c4Path)
			aGrp.RepeatCount = ledgrid.AnimationRepeatForever
			aGrp.Start()
		})

	PolygonPathTest = NewLedGridProgram("Polygon path test",
		func(c *ledgrid.Canvas) {

			cPos := geom.Point{1, 1}
			cSize := geom.Point{1, 1}

			polyPath1 := ledgrid.NewPolygonPath(
				geom.Point{1, 1},
				geom.Point{float64(width) - 1, 1},
				geom.Point{float64(width) - 1, float64(height) - 1},
				geom.Point{1, float64(height) - 1},

				geom.Point{1, 2},
				geom.Point{float64(width) - 2, 2},
				geom.Point{float64(width) - 2, float64(height) - 2},
				geom.Point{2, float64(height) - 2},

				geom.Point{2, 3},
				geom.Point{float64(width) - 3, 3},
				geom.Point{float64(width) - 3, float64(height) - 3},
				geom.Point{3, float64(height) - 3},

				geom.Point{3, 4},
				geom.Point{float64(width) - 4, 4},
				geom.Point{float64(width) - 4, float64(height) - 4},
				geom.Point{4, float64(height) - 4},
			)

			polyPath2 := ledgrid.NewPolygonPath(
				geom.Point{1, 1},
				geom.Point{4, 9},
				geom.Point{7, 2},
				geom.Point{10, 8},
				geom.Point{13, 3},
				geom.Point{16, 7},
				geom.Point{19, 4},
				geom.Point{22, 6},
			)

			c1 := ledgrid.NewEllipse(cPos, cSize, color.GreenYellow)
			c.Add(c1)

			aPath1 := ledgrid.NewPolyPathAnimation(&c1.Pos, polyPath1, 7*time.Second)
			aPath1.AutoReverse = true

			aPath2 := ledgrid.NewPolyPathAnimation(&c1.Pos, polyPath2, 7*time.Second)
			aPath2.AutoReverse = true

			seq := ledgrid.NewSequence(aPath1, aPath2)
			seq.RepeatCount = ledgrid.AnimationRepeatForever

			seq.Start()
		})

	RandomWalk = NewLedGridProgram("Random walk",
		func(c *ledgrid.Canvas) {
			rect := geom.Rectangle{Min: geom.Point{1.5, 1.5}, Max: geom.Point{float64(width) - 0.5, float64(height) - 0.5}}
			pos1 := geom.Point{1.5, 1.5}
			pos2 := geom.Point{18.5, 8.5}
			size1 := geom.Point{2.0, 2.0}
			size2 := geom.Point{4.0, 4.0}

			c1 := ledgrid.NewEllipse(pos1, size1, color.SkyBlue)
			c2 := ledgrid.NewEllipse(pos2, size2, color.GreenYellow)
			c.Add(c1, c2)

			aPos1 := ledgrid.NewPositionAnimation(&c1.Pos, geom.Point{}, 1300*time.Millisecond)
			aPos1.Cont = true
			aPos1.ValFunc = ledgrid.RandPointTrunc(rect, 1.0)
			aPos1.RepeatCount = ledgrid.AnimationRepeatForever

			aPos2 := ledgrid.NewPositionAnimation(&c2.Pos, geom.Point{}, 901*time.Millisecond)
			aPos2.Cont = true
			aPos2.ValFunc = ledgrid.RandPoint(rect)
			aPos2.RepeatCount = ledgrid.AnimationRepeatForever

			aPos1.Start()
			aPos2.Start()
		})

	CirclingCircles = NewLedGridProgram("Circling circles",
		func(c *ledgrid.Canvas) {
			pos1 := geom.Point{1.5, 1.5}
			pos2 := geom.Point{10.5, 8.5}
			pos3 := geom.Point{19.5, 1.5}
			pos4 := geom.Point{28.5, 8.5}
			pos5 := geom.Point{37.5, 1.5}
			cSize := geom.Point{2.0, 2.0}

			c1 := ledgrid.NewEllipse(pos1, cSize, color.OrangeRed)
			c2 := ledgrid.NewEllipse(pos2, cSize, color.MediumSeaGreen)
			c3 := ledgrid.NewEllipse(pos3, cSize, color.SkyBlue)
			c4 := ledgrid.NewEllipse(pos4, cSize, color.Gold)
			c5 := ledgrid.NewEllipse(pos5, cSize, color.YellowGreen)

			stepRD := geom.Point{9.0, 7.0}
			stepLU := stepRD.Neg()
			stepRU := geom.Point{9.0, -7.0}
			stepLD := stepRU.Neg()

			c1Path1 := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.QuarterCirclePathA, stepRD, time.Second)
			c1Path1.Cont = true
			c1Path2 := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.QuarterCirclePathA, stepRU, time.Second)
			c1Path2.Cont = true
			c1Path3 := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.QuarterCirclePathA, stepLD, time.Second)
			c1Path3.Cont = true
			c1Path4 := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.QuarterCirclePathA, stepLU, time.Second)
			c1Path4.Cont = true

			c2Path1 := ledgrid.NewPathAnimation(&c2.Pos, ledgrid.QuarterCirclePathA, stepLU, time.Second)
			c2Path1.Cont = true
			c2Path2 := ledgrid.NewPathAnimation(&c2.Pos, ledgrid.QuarterCirclePathA, stepRD, time.Second)
			c2Path2.Cont = true

			c3Path1 := ledgrid.NewPathAnimation(&c3.Pos, ledgrid.QuarterCirclePathA, stepLD, time.Second)
			c3Path1.Cont = true
			c3Path2 := ledgrid.NewPathAnimation(&c3.Pos, ledgrid.QuarterCirclePathA, stepRU, time.Second)
			c3Path2.Cont = true

			c4Path1 := ledgrid.NewPathAnimation(&c4.Pos, ledgrid.QuarterCirclePathA, stepLU, time.Second)
			c4Path1.Cont = true
			c4Path2 := ledgrid.NewPathAnimation(&c4.Pos, ledgrid.QuarterCirclePathA, stepRD, time.Second)
			c4Path2.Cont = true

			c5Path1 := ledgrid.NewPathAnimation(&c5.Pos, ledgrid.QuarterCirclePathA, stepLD, time.Second)
			c5Path1.Cont = true
			c5Path2 := ledgrid.NewPathAnimation(&c5.Pos, ledgrid.QuarterCirclePathA, stepRU, time.Second)
			c5Path2.Cont = true

			aGrp1 := ledgrid.NewGroup(c1Path1, c2Path1)
			aGrp2 := ledgrid.NewGroup(c1Path2, c3Path1)
			aGrp3 := ledgrid.NewGroup(c1Path1, c4Path1)
			aGrp4 := ledgrid.NewGroup(c1Path2, c5Path1)
			aGrp5 := ledgrid.NewGroup(c1Path3, c5Path2)
			aGrp6 := ledgrid.NewGroup(c1Path4, c4Path2)
			aGrp7 := ledgrid.NewGroup(c1Path3, c3Path2)
			aGrp8 := ledgrid.NewGroup(c1Path4, c2Path2)
			aSeq := ledgrid.NewSequence(aGrp1, aGrp2, aGrp3, aGrp4, aGrp5, aGrp6, aGrp7, aGrp8)
			aSeq.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(c1, c2, c3, c4, c5)
			aSeq.Start()
		})

	ChasingCircles = NewLedGridProgram("Chasing circles",
		func(c *ledgrid.Canvas) {
			c1Pos1 := geom.Point{37.0, 5.0}
			c1Size1 := geom.Point{10.0, 10.0}
			c1Size2 := geom.Point{3.0, 3.0}
			c1PosSize := geom.Point{-34.0, -5.0}
			c2Pos := geom.Point{3.0, 5.0}
			c2Size1 := geom.Point{5.0, 5.0}
			c2Size2 := geom.Point{3.0, 3.0}
			c2PosSize := geom.Point{34.0, 7.0}

			aGrp := ledgrid.NewGroup()

			pal := ledgrid.NewGradientPaletteByList("Palette", true,
				color.DeepSkyBlue,
				color.Lime,
				color.Teal,
				color.SkyBlue,
			)

			c1 := ledgrid.NewEllipse(c1Pos1, c1Size1, color.Gold)

			c1pos := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.FullCirclePathB, c1PosSize, 4*time.Second)
			c1pos.RepeatCount = ledgrid.AnimationRepeatForever
			c1pos.Curve = ledgrid.AnimationLinear

			c1size := ledgrid.NewSizeAnimation(&c1.Size, c1Size2, time.Second)
			c1size.AutoReverse = true
			c1size.RepeatCount = ledgrid.AnimationRepeatForever

			c1bcolor := ledgrid.NewColorAnimation(&c1.BorderColor, color.OrangeRed, time.Second)
			c1bcolor.AutoReverse = true
			c1bcolor.RepeatCount = ledgrid.AnimationRepeatForever

			aGrp.Add(c1pos, c1size, c1bcolor)

			c2 := ledgrid.NewEllipse(c2Pos, c2Size1, color.Lime)

			c2pos := ledgrid.NewPathAnimation(&c2.Pos, ledgrid.FullCirclePathB, c2PosSize, 4*time.Second)
			c2pos.RepeatCount = ledgrid.AnimationRepeatForever
			c2pos.Curve = ledgrid.AnimationLinear

			c2size := ledgrid.NewSizeAnimation(&c2.Size, c2Size2, time.Second)
			c2size.AutoReverse = true
			c2size.RepeatCount = ledgrid.AnimationRepeatForever

			c2color := ledgrid.NewPaletteAnimation(&c2.BorderColor, pal, 2*time.Second)
			c2color.RepeatCount = ledgrid.AnimationRepeatForever
			c2color.Curve = ledgrid.AnimationLinear

			aGrp.Add(c2pos, c2size, c2color)

			c.Add(c2, c1)
			aGrp.Start()
		})

	CircleAnimation = NewLedGridProgram("Circle animation",
		func(c *ledgrid.Canvas) {
			c1Pos1 := geom.Point{2.0, 5.0}
			c1Pos3 := geom.Point{38.0, 5.0}

			c1Size1 := geom.Point{3.0, 3.0}
			c1Size2 := geom.Point{9.0, 9.0}

			c1 := ledgrid.NewEllipse(c1Pos1, c1Size1, color.OrangeRed)

			c1pos := ledgrid.NewPositionAnimation(&c1.Pos, c1Pos3, 2*time.Second)
			c1pos.AutoReverse = true
			c1pos.RepeatCount = ledgrid.AnimationRepeatForever

			c1radius := ledgrid.NewSizeAnimation(&c1.Size, c1Size2, time.Second)
			c1radius.AutoReverse = true
			c1radius.RepeatCount = ledgrid.AnimationRepeatForever

			c1color := ledgrid.NewColorAnimation(&c1.BorderColor, color.Gold, 4*time.Second)
			c1color.AutoReverse = true
			c1color.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(c1)

			c1pos.Start()
			c1radius.Start()
			c1color.Start()
		})

	PushingRectangles = NewLedGridProgram("Pushing rectangles",
		func(c *ledgrid.Canvas) {
			rSize1 := geom.Point{float64(width) - 3.0, 1.0}
			rSize2 := geom.Point{1.0, float64(height) - 1.0}

			r1Pos1 := geom.Point{1.0, float64(height) / 2.0}
			r1Pos2 := geom.Point{1.0 + float64(width-3)/2.0, float64(height) / 2.0}

			r2Pos1 := geom.Point{float64(width - 1), float64(height) / 2.0}
			r2Pos2 := geom.Point{float64(width-1) - float64(width-3)/2.0, float64(height) / 2.0}
			duration := 2 * time.Second

			r1 := ledgrid.NewRectangle(r1Pos1, rSize2, color.Crimson)

			a1Pos := ledgrid.NewPositionAnimation(&r1.Pos, r1Pos2, duration)
			a1Pos.AutoReverse = true
			a1Pos.RepeatCount = ledgrid.AnimationRepeatForever

			a1Size := ledgrid.NewSizeAnimation(&r1.Size, rSize1, duration)
			a1Size.AutoReverse = true
			a1Size.RepeatCount = ledgrid.AnimationRepeatForever

			a1Color := ledgrid.NewColorAnimation(&r1.BorderColor, color.GreenYellow, duration)
			a1Color.AutoReverse = true
			a1Color.RepeatCount = ledgrid.AnimationRepeatForever

			r2 := ledgrid.NewRectangle(r2Pos2, rSize1, color.SkyBlue)

			a2Pos := ledgrid.NewPositionAnimation(&r2.Pos, r2Pos1, duration)
			a2Pos.AutoReverse = true
			a2Pos.RepeatCount = ledgrid.AnimationRepeatForever

			a2Size := ledgrid.NewSizeAnimation(&r2.Size, rSize2, duration)
			a2Size.AutoReverse = true
			a2Size.RepeatCount = ledgrid.AnimationRepeatForever

			a2Color := ledgrid.NewColorAnimation(&r2.BorderColor, color.Crimson, duration)
			a2Color.AutoReverse = true
			a2Color.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(r1, r2)
			a1Pos.Start()
			a1Size.Start()
			a1Color.Start()
			a2Pos.Start()
			a2Size.Start()
			a2Color.Start()
		})

	RegularPolygonTest = NewLedGridProgram("Regular Polygon test",
		func(c *ledgrid.Canvas) {
			posList := []geom.Point{
				geom.Point{-5.5, 4.5},
				geom.Point{44.5, 4.5},
			}
			posCenter := geom.Point{19.5, 4.5}
			smallSize := geom.Point{9.0, 9.0}
			largeSize := geom.Point{80.0, 80.0}

			polyList := make([]*ledgrid.RegularPolygon, 9)

			aSeq := ledgrid.NewSequence()
			for n := 3; n <= 6; n++ {
				col := color.RandColor()
				polyList[n] = ledgrid.NewRegularPolygon(n, posList[n%2], smallSize, col)
				c.Add(polyList[n])
				dur := 2*time.Second + rand.N(time.Second)
				sign := []float64{+1.0, -1.0}[n%2]
				angle := sign * 2 * math.Pi
				animPos := ledgrid.NewPositionAnimation(&polyList[n].Pos, posCenter, dur)
				animAngle := ledgrid.NewFloatAnimation(&polyList[n].Angle, angle, dur)
				animSize := ledgrid.NewSizeAnimation(&polyList[n].Size, largeSize, 4*time.Second)
				animSize.Cont = true
				animFade := ledgrid.NewColorAnimation(&polyList[n].BorderColor, color.Black, 4*time.Second)
				animFade.Cont = true
				// animFade := NewAlphaAnimation(&polyList[n].FillColor.A, 0x00, time.Second)

				aGrpIn := ledgrid.NewGroup(animPos, animAngle)
				// aGrp.duration = dur
				aGrpOut := ledgrid.NewGroup(animSize, animFade)
				aObjSeq := ledgrid.NewSequence(aGrpIn, aGrpOut)
				aSeq.Add(aObjSeq)
			}
			aSeq.RepeatCount = ledgrid.AnimationRepeatForever
			aSeq.Start()
		})

	MovingText = NewLedGridProgram("Moving text",
		func(c *ledgrid.Canvas) {
			pts = []geom.Point{
				geom.Point{0, 0},
				geom.Point{0, float64(height)},
				geom.Point{float64(width), float64(height)},
				geom.Point{float64(width), 0},
			}

			t1 := ledgrid.NewText(randPoint(), "Mühlebach", color.LightSeaGreen)
			t2 := ledgrid.NewText(randPoint(), "Mathematik", color.YellowGreen)
			t3 := ledgrid.NewText(randPoint(), "Benedict", color.OrangeRed)
			c.Add(t1, t2, t3)

			aPos1 := ledgrid.NewPositionAnimation(&t1.Pos, geom.Point{}, 5*time.Second)
			aPos1.ValFunc = randPoint
			aPos1.RepeatCount = ledgrid.AnimationRepeatForever
			aPos1.Cont = true

			aPos2 := ledgrid.NewPositionAnimation(&t2.Pos, geom.Point{}, 3*time.Second)
			aPos2.ValFunc = randPoint
			aPos2.RepeatCount = ledgrid.AnimationRepeatForever
			aPos2.Cont = true

			aPos3 := ledgrid.NewPositionAnimation(&t3.Pos, geom.Point{}, 2*time.Second)
			aPos3.ValFunc = randPoint
			aPos3.RepeatCount = ledgrid.AnimationRepeatForever
			aPos3.Cont = true

			aAngle1 := ledgrid.NewFloatAnimation(&t1.Angle, 0.0, 3*time.Second)
			aAngle1.ValFunc = ledgrid.RandFloat(math.Pi/2.0, math.Pi)
			aAngle1.AutoReverse = true
			aAngle1.RepeatCount = ledgrid.AnimationRepeatForever

			aAngle2 := ledgrid.NewFloatAnimation(&t2.Angle, 0.0, 4*time.Second)
			aAngle2.ValFunc = ledgrid.RandFloat(math.Pi/6.0, math.Pi/2.0)
			aAngle2.AutoReverse = true
			aAngle2.RepeatCount = ledgrid.AnimationRepeatForever

			aAngle1.Start()
			aAngle2.Start()
			aPos1.Start()
			aPos2.Start()
			aPos3.Start()
		})

	BitmapText = NewLedGridProgram("Bitmap text",
		func(c *ledgrid.Canvas) {
			basePt := fixed.P(0, 5)
			baseColor1 := color.SkyBlue

			txt1 := ledgrid.NewFixedText(basePt, baseColor1.Alpha(0.0), "STEFAN")
			c.Add(txt1)

			aTxt1 := ledgrid.NewAlphaAnimation(&txt1.Color.A, 255, 2*time.Second)
			aTxt1.AutoReverse = true
			aTxt1.RepeatCount = ledgrid.AnimationRepeatForever

			aPos := ledgrid.NewFixedPosAnimation(&txt1.Pos, fixed.P(25, 5), 3*time.Second)
			aPos.AutoReverse = true
			aPos.RepeatCount = ledgrid.AnimationRepeatForever

			aTxt1.Start()
			aPos.Start()
		})

	FlyingImages = NewLedGridProgram("Flying images",
		func(c *ledgrid.Canvas) {
			pos1 := geom.Point{20, -5}
			pos2 := geom.Point{20, 5}
			pos3 := geom.Point{20, 15}

			size1 := geom.Point{float64(width) / 3.0, float64(height) / 3.0}
			size2 := geom.Point{3.0 * float64(width), 3.0 * float64(height)}

			size3 := geom.Point{1.0, 0.75}
			size4 := geom.Point{160.0, 120.0}

			img1 := ledgrid.NewImage(pos1, "images/ledgrid.png")
			img1.Size = size1
			// img1.Hide()

			img2 := ledgrid.NewImage(pos2, "images/testbild.png")
			img2.Size = size3
			img2.Hide()

			c.Add(img1, img2)

			aPos1 := ledgrid.NewPositionAnimation(&img1.Pos, pos2, 4*time.Second)
			aPos2 := ledgrid.NewPositionAnimation(&img1.Pos, pos3, 4*time.Second)
			aPos2.Cont = true

			aAngle1 := ledgrid.NewFloatAnimation(&img1.Angle, math.Pi, 4*time.Second)
			aAngle2 := ledgrid.NewFloatAnimation(&img1.Angle, 0.0, 4*time.Second)
			aAngle2.Cont = true

			aSize1 := ledgrid.NewSizeAnimation(&img1.Size, size2, 4*time.Second)
			aSize1.Cont = true
			aSize2 := ledgrid.NewSizeAnimation(&img1.Size, size1, 4*time.Second)
			aSize2.Cont = true

			aSize3 := ledgrid.NewSizeAnimation(&img2.Size, size4, 4*time.Second)
			aSize3.AutoReverse = true
			aSize3.RepeatCount = ledgrid.AnimationRepeatForever
			aAngle3 := ledgrid.NewFloatAnimation(&img2.Angle, 4*math.Pi, 4*time.Second)
			aAngle3.AutoReverse = true
			aAngle3.RepeatCount = ledgrid.AnimationRepeatForever

			// task := ledgrid.NewBackgroundTask(func() { img2.Show() })
			showHide := ledgrid.NewShowHideAnimation(img2)

			aSize3.Start()
			aAngle3.Start()
			aSeq := ledgrid.NewSequence(aPos1, aAngle1, showHide, aSize1, aPos2, showHide, aAngle2, aSize2, aSize3)
			aSeq.Start()
		})

	CameraTest = NewLedGridProgram("Camera test",
		func(c *ledgrid.Canvas) {
			pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			size := geom.Point{float64(width), float64(height)}

			pos1 := image.Point{1, 1}
			pos2 := image.Point{width - 2, 1}

			cam := ledgrid.NewCamera(pos, size)
			pix1 := ledgrid.NewPixel(pos1, color.YellowGreen)
			c.Add(cam, pix1)

			aPos1 := ledgrid.NewIntegerPosAnimation(&pix1.Pos, pos2, 5*time.Second)
			aPos1.AutoReverse = true
			aPos1.RepeatCount = ledgrid.AnimationRepeatForever

			aColor1 := ledgrid.NewColorAnimation(&pix1.Color, color.OrangeRed, 2000*time.Millisecond)
			aColor1.AutoReverse = true
			aColor1.RepeatCount = ledgrid.AnimationRepeatForever

			aPos1.Start()
			aColor1.Start()
			cam.Start()

			zoom := cam.Get(gocv.VideoCaptureZoom)
			gamma := cam.Get(gocv.VideoCaptureGamma)
            bright := cam.Get(gocv.VideoCaptureBrightness)
            satur := cam.Get(gocv.VideoCaptureSaturation)
			fmt.Printf("\nzoom: %f\ngamma: %f\nbright: %f\nsatur: %f\n", zoom, gamma, bright, satur)

			cam.Set(gocv.VideoCaptureZoom, 100.0)
            cam.Set(gocv.VideoCaptureGamma, 2.0)
            cam.Set(gocv.VideoCaptureBrightness, 128.0)
            cam.Set(gocv.VideoCaptureSaturation, 255.0)

			zoom = cam.Get(gocv.VideoCaptureZoom)
			gamma = cam.Get(gocv.VideoCaptureGamma)
            bright = cam.Get(gocv.VideoCaptureBrightness)
            satur = cam.Get(gocv.VideoCaptureSaturation)
			fmt.Printf("\nzoom: %f\ngamma: %f\nbright: %f\nsatur: %f\n", zoom, gamma, bright, satur)
		})

	BlinkenAnimation = NewLedGridProgram("Blinken animation",
		func(c *ledgrid.Canvas) {
			posA := geom.Point{20.0, 5.0}
			// posB := geom.Point{39.0, 5.0}

			bml := ledgrid.ReadBlinkenFile("blinken/mario.bml")

			imgList := ledgrid.NewImageList(posA)
			imgList.AddBlinkenLight(bml)
			imgList.Size = geom.Point{10.0, 10.0}

			// aPos := ledgrid.NewPositionAnimation(&imgList.Pos, posB, 3*time.Second)
			// aPos.Curve = AnimationLinear

			aImgList := ledgrid.NewImageAnimation(&imgList.ImgIdx)
			aImgList.RepeatCount = ledgrid.AnimationRepeatForever
			aImgList.AddBlinkenLight(bml)

			aGrp := ledgrid.NewGroup(aImgList)

			c.Add(imgList)
			aGrp.Start()
		})

	Piiiiixels = NewLedGridProgram("Piiiiixels",
		func(c *ledgrid.Canvas) {
			dPosX := image.Point{2, 0}
			dPosY := image.Point{0, 2}
			p1Pos1 := image.Point{1, 1}

			aGrp := ledgrid.NewGroup()
			aGrp.RepeatCount = ledgrid.AnimationRepeatForever
			for i := range 5 {
				for j := range 5 {
					palName := ledgrid.PaletteNames[5*i+j]
					pos := p1Pos1.Add(dPosX.Mul(j).Add(dPosY.Mul(i)))
					pix := ledgrid.NewPixel(pos, color.OrangeRed)
					c.Add(pix)

					pixColor := ledgrid.NewPaletteAnimation(&pix.Color, ledgrid.PaletteMap[palName], 4*time.Second)
					aGrp.Add(pixColor)
				}
			}
			aGrp.Start()
		})

	MovingPixels = NewLedGridProgram("Moving pixels",
		func(c *ledgrid.Canvas) {
			pos1 := image.Point{1, 1}
			pos2 := image.Point{width - 2, 1}

			pos3 := image.Point{width - 2, 3}
			pos4 := image.Point{1, 3}

			pix1 := ledgrid.NewPixel(pos1, color.YellowGreen)
			pix2 := ledgrid.NewPixel(pos3, color.LightSeaGreen)
			c.Add(pix1, pix2)

			aPos1 := ledgrid.NewIntegerPosAnimation(&pix1.Pos, pos2, 5*time.Second)
			aPos1.AutoReverse = true
			aPos1.RepeatCount = ledgrid.AnimationRepeatForever
			aPos2 := ledgrid.NewIntegerPosAnimation(&pix2.Pos, pos4, 4*time.Second)
			aPos2.AutoReverse = true
			aPos2.RepeatCount = ledgrid.AnimationRepeatForever
			aPos1.Start()
			aPos2.Start()

			aColor1 := ledgrid.NewColorAnimation(&pix1.Color, color.OrangeRed, 2000*time.Millisecond)
			aColor1.AutoReverse = true
			aColor1.RepeatCount = ledgrid.AnimationRepeatForever
			aColor1.Start()
			aColor2 := ledgrid.NewColorAnimation(&pix2.Color, color.Purple, 2123*time.Millisecond)
			aColor2.AutoReverse = true
			aColor2.RepeatCount = ledgrid.AnimationRepeatForever
			aColor2.Start()
		})

	GlowingPixels = NewLedGridProgram("Glowing pixels",
		func(c *ledgrid.Canvas) {
			aGrpPurple := ledgrid.NewGroup()
			aGrpYellow := ledgrid.NewGroup()
			aGrpGreen := ledgrid.NewGroup()
			aGrpGrey := ledgrid.NewGroup()

			for y := range c.Rect.Dy() {
				for x := range c.Rect.Dx() {
					pos := image.Point{x, y}
					t := rand.Float64()
					col := (color.DimGray.Dark(0.3)).Interpolate((color.DarkGrey.Dark(0.3)), t)
					pix := ledgrid.NewPixel(pos, col)
					c.Add(pix)

					dur := time.Second + rand.N(time.Second)
					aAlpha := ledgrid.NewAlphaAnimation(&pix.Color.A, 196, dur)
					aAlpha.AutoReverse = true
					aAlpha.RepeatCount = ledgrid.AnimationRepeatForever
					aAlpha.Start()

					aColor := ledgrid.NewColorAnimation(&pix.Color, col, 1*time.Second)
					aColor.Cont = true
					aGrpGrey.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.MediumPurple.Interpolate(color.Fuchsia, t), 5*time.Second)
					aColor.Cont = true
					aGrpPurple.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.Gold.Interpolate(color.Khaki, t), 5*time.Second)
					aColor.Cont = true
					aGrpYellow.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.GreenYellow.Interpolate(color.LightSeaGreen, t), 5*time.Second)
					aColor.Cont = true
					aGrpGreen.Add(aColor)
				}
			}

			txt1 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.GreenYellow.Alpha(0.0), "LORENZ")
			aTxt1 := ledgrid.NewAlphaAnimation(&txt1.Color.A, 255, 2*time.Second)
			aTxt1.AutoReverse = true
			txt2 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.DarkViolet.Alpha(0.0), "SIMON")
			aTxt2 := ledgrid.NewAlphaAnimation(&txt2.Color.A, 255, 2*time.Second)
			aTxt2.AutoReverse = true
			txt3 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.OrangeRed.Alpha(0.0), "REBEKKA")
			aTxt3 := ledgrid.NewAlphaAnimation(&txt3.Color.A, 255, 2*time.Second)
			aTxt3.AutoReverse = true
			c.Add(txt1, txt2, txt3)

			aTimel := ledgrid.NewTimeline(42 * time.Second)
			aTimel.Add(7*time.Second, aGrpPurple)
			aTimel.Add(12*time.Second, aTxt1)
			aTimel.Add(13*time.Second, aGrpGrey)

			aTimel.Add(22*time.Second, aGrpYellow)
			aTimel.Add(27*time.Second, aTxt2)
			aTimel.Add(28*time.Second, aGrpGrey)

			aTimel.Add(35*time.Second, aGrpGreen)
			aTimel.Add(40*time.Second, aTxt3)
			aTimel.Add(41*time.Second, aGrpGrey)
			aTimel.RepeatCount = ledgrid.AnimationRepeatForever

			aTimel.Start()
		})

	TestShaderFunc = func(x, y, t float64) float64 {
		t = t/4.0 + x
		_, f := math.Modf(math.Abs(t))
		return f
	}

	f1 = func(x, y, t, p1 float64) float64 {
		return math.Sin(x*p1 + t)
	}

	f2 = func(x, y, t, p1, p2, p3 float64) float64 {
		return math.Sin(p1*(x*math.Sin(t/p2)+y*math.Cos(t/p3)) + t)
	}

	f3 = func(x, y, t, p1, p2 float64) float64 {
		cx := 0.125*x + 0.5*math.Sin(t/p1)
		cy := 0.125*y + 0.5*math.Cos(t/p2)
		return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
	}

	PlasmaShaderFunc = func(x, y, t float64) float64 {
		v1 := f1(x, y, t, 1.2)
		v2 := f2(x, y, t, 1.6, 3.0, 1.5)
		v3 := f3(x, y, t, 5.0, 5.0)
		v := (v1+v2+v3)/6.0 + 0.5
		return v
	}

	ShowTheShader = NewLedGridProgram("Show the shader!",
		func(c *ledgrid.Canvas) {
			var xMin, yMax float64

			pal := ledgrid.PaletteMap["Hipster"]
			fader := ledgrid.NewPaletteFader(pal)
			aPal := ledgrid.NewPaletteFadeAnimation(fader, pal, 2*time.Second)
			aPal.ValFunc = ledgrid.SeqPalette()

			aPalTl := ledgrid.NewTimeline(10 * time.Second)
			aPalTl.Add(7*time.Second, aPal)
			aPalTl.RepeatCount = ledgrid.AnimationRepeatForever

			aGrp := ledgrid.NewGroup()
			dPix := 2.0 / float64(max(c.Rect.Dx(), c.Rect.Dy())-1)
			ratio := float64(c.Rect.Dx()) / float64(c.Rect.Dy())
			if ratio > 1.0 {
				xMin = -1.0
				yMax = ratio * 1.0
			} else if ratio < 1.0 {
				xMin = ratio * -1.0
				yMax = 1.0
			} else {
				xMin = -1.0
				yMax = 1.0
			}

			y := yMax
			for row := range c.Rect.Dy() {
				x := xMin
				for col := range c.Rect.Dx() {
					pix := ledgrid.NewPixel(image.Point{col, row}, color.Black)
					c.Add(pix)
					anim := ledgrid.NewShaderAnimation(&pix.Color, fader, x, y, TestShaderFunc)
					aGrp.Add(anim)
					x += dPix
				}
				y -= dPix
			}
			aPalTl.Start()
			aGrp.Start()
		})
)

//----------------------------------------------------------------------------

type BouncingEllipse struct {
	ledgrid.Ellipse
	Vel, Acc geom.Point
	Field    geom.Rectangle
}

func NewBouncingEllipse(pos, size geom.Point, col color.LedColor) *BouncingEllipse {
	b := &BouncingEllipse{}
	b.Ellipse = *ledgrid.NewEllipse(pos, size, col)
	b.Vel = geom.Point{}
	b.Acc = geom.Point{}
	return b
}

func (b *BouncingEllipse) Update(pit time.Time) bool {
	deltaVel := b.Acc.Mul(0.3)
	b.Vel = b.Vel.Add(deltaVel)
	b.Pos = b.Pos.Add(b.Vel)
	if b.Pos.X < b.Field.Min.X || b.Pos.X >= b.Field.Max.X {
		b.Vel.X = -b.Vel.X
	}
	if b.Pos.Y < b.Field.Min.Y || b.Pos.Y >= b.Field.Max.Y {
		b.Vel.Y = -b.Vel.Y
	}
	return true
}

func (b *BouncingEllipse) Duration() time.Duration {
	return time.Duration(0)
}
func (b *BouncingEllipse) SetDuration(dur time.Duration) {}
func (b *BouncingEllipse) Start()                        {}
func (b *BouncingEllipse) Stop()                         {}
func (b *BouncingEllipse) Continue()                     {}
func (b *BouncingEllipse) IsStopped() bool {
	return false
}

func BounceAround(c *ledgrid.Canvas) {
	pos1 := geom.Point{2.0, 2.0}
	pos2 := geom.Point{37.0, 7.0}
	size := geom.Point{4.0, 4.0}
	vel1 := geom.Point{0.15, 0.075}
	vel2 := geom.Point{-0.35, -0.25}

	obj1 := NewBouncingEllipse(pos1, size, color.GreenYellow)
	obj1.Vel = vel1
	obj1.Field = geom.NewRectangleIMG(c.Rect)
	obj2 := NewBouncingEllipse(pos2, size, color.LightSeaGreen)
	obj2.Vel = vel2
	obj2.Field = geom.NewRectangleIMG(c.Rect)

	c.Add(obj1, obj2)
	animCtrl.Add(obj1, obj2)
}

//----------------------------------------------------------------------------

// func pixIdx(x, y int) int {
// 	return y*width + x
// }

// func pixCoord(idx int) (x, y int) {
// 	return idx % width, idx / width
// }

var (
	pts     []geom.Point
	lastIdx = -1
)

func randPoint() geom.Point {
	idx0 := rand.IntN(len(pts))
	for idx0 == lastIdx {
		idx0 = rand.IntN(len(pts))
	}
	lastIdx = idx0
	idx1 := (idx0 + 1) % len(pts)

	return pts[idx0].Interpolate(pts[idx1], rand.Float64())
}

// func PlasmaShaderFunc(x, y, t float64) float64 {
// 	v1 := f1(x, y, t, 1.2)
// 	v2 := f2(x, y, t, 1.6, 3.0, 1.5)
// 	v3 := f3(x, y, t, 5.0, 5.0)
// 	v := (v1+v2+v3)/6.0 + 0.5
// 	return v
// }

// func f1(x, y, t, p1 float64) float64 {
// 	return math.Sin(x*p1 + t)
// }

// func f2(x, y, t, p1, p2, p3 float64) float64 {
// 	return math.Sin(p1*(x*math.Sin(t/p2)+y*math.Cos(t/p3)) + t)
// }

// func f3(x, y, t, p1, p2 float64) float64 {
// 	cx := 0.125*x + 0.5*math.Sin(t/p1)
// 	cy := 0.125*y + 0.5*math.Cos(t/p2)
// 	return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
// }

// func RandomGridPixels(g *Grid) {
// 	for y := range g.ledGrid.Rect.Dy() {
// 		for x := range g.ledGrid.Rect.Dx() {
// 			pos := image.Pt(x, y)
// 			colorGrp1 := color.ColorGroup(x/3) % color.NumColorGroups
// 			colorGrp2 := (colorGrp1 + 1) % color.NumColorGroups
// 			col := color.RandGroupColor(colorGrp1)
// 			pix := NewGridPixel(pos, col)
// 			g.Add(pix)
// 			dur := time.Second
// 			aColor := ledgrid.NewColorAnimation(&pix.Color, ledgrid.RandGroupColor(colorGrp2), dur)
// 			aColor.AutoReverse = true
// 			aColor.RepeatCount = ledgrid.AnimationRepeatForever
// 			aColor.Start()
// 		}
// 	}
// }

// func TextOnGrid(g *Grid) {
// 	basePt := fixed.P(0, 5)
// 	baseColor1 := color.SkyBlue

// 	txt1 := NewGridText(basePt, baseColor1, "STEFAN")
// 	g.Add(txt1)

// 	aPos := NewFixedPosAnimation(&txt1.Pos, fixed.P(25, 5), 3*time.Second)
// 	aPos.AutoReverse = true
// 	aPos.RepeatCount = ledgrid.AnimationRepeatForever
// 	aPos.Start()
// }

// func WalkingPixelOnGrid(g *Grid) {
// 	pos := image.Point{0, 0}
// 	col := color.GreenYellow
// 	pix := NewGridPixel(pos, col)
// 	g.Add(pix)

// 	go func() {
// 		idx := 0
// 		for {
// 			col := idx % width
// 			row := idx / width
// 			pix.Pos = image.Point{col, row}
// 			time.Sleep(time.Second / 5)
// 			idx++
// 		}
// 	}()
// }

// func ImagesOnGrid(g *Grid) {
// 	pos := image.Point{5, 2}
// 	size := image.Point{10, 10}

// 	img := NewGridImage(pos, size)
// 	for row := range size.Y {
// 		for col := range size.X {
// 			img.Img.SetRGBA(col, row, color.RGBA{0x8f, 0x8f, 0x8f, 0xff})
// 		}
// 	}

// 	g.Add(img)
// }

//----------------------------------------------------------------------------

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

//----------------------------------------------------------------------------

type canvasSceneRecord struct {
	name string
	fnc  func(canvas *ledgrid.Canvas)
}

var (
	canvasSceneList = []canvasSceneRecord{
		// {"Group test", GroupTest},
		// {"Sequence test", SequenceTest},
		// {"Timeline test", TimelineTest},
		// {"Path test", PathTest},
		// {"PolygonPath Test", PolygonPathTest},
		// {"Random walk", RandomWalk},
		{"Let's bounce around", BounceAround},
		// {"Piiiiixels", Piiiiixels},
		// {"Regular Polygon test", ledgrid.RegularPolygonTest},
		// {"Circling Circles", CirclingCircles},
		// {"Chasing Circles", ChasingCircles},
		// {"Circle Animation", CircleAnimation},
		// {"Pushing Rectangles", PushingRectangles},
		// {"Glowing Pixels", GlowingPixels},
		// {"Moving Text", MovingText},
		// {"Bitmap Text", BitmapText},
		// {"Flying images", FlyingImages},
		// {"Live Camera stream", CameraTest},
		// {"Animation from a BlinkenLight file", BlinkenAnimation},
		// {"Show the Shader", ShowTheShader},
	}

	programList = []LedGridProgram{
		GroupTest,
		SequenceTest,
		TimelineTest,
		StopContShowHideTest,
		PathTest,
		PolygonPathTest,
		RandomWalk,
		CirclingCircles,
		ChasingCircles,
		PushingRectangles,
		RegularPolygonTest,
		Piiiiixels,
		GlowingPixels,
		MovingPixels,
		CameraTest,
		BlinkenAnimation,
		MovingText,
		BitmapText,
		FlyingImages,
		ShowTheShader,
	}
)

func main() {
	var host string
	var port uint
	var input string
	var ch byte
	var progId int
	var runInteractive bool
	var progList string

	for i, prog := range programList {
		id := 'a' + i
		progList += ("\n" + string(id) + " - " + prog.Name())
	}

	flag.IntVar(&width, "width", defWidth, "Width of panel")
	flag.IntVar(&height, "height", defHeight, "Height of panel")
	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.StringVar(&input, "prog", input, "Play one single program"+progList)
	flag.Parse()

	gridSize = image.Point{width, height}

	if len(input) > 0 {
		runInteractive = false
		ch = input[0]
	} else {
		runInteractive = true
	}

	canvas := ledgrid.NewCanvas(gridSize)
	ledGrid := ledgrid.NewLedGrid(gridSize, nil)
	pixClient := ledgrid.NewNetPixelClient(host, port)

	animCtrl = ledgrid.NewAnimationController(canvas, ledGrid, pixClient)

	if runInteractive {
		progId = -1
		for {
			fmt.Printf("---------------------------------------\n")
			fmt.Printf("Program\n")
			fmt.Printf("---------------------------------------\n")
			for i, prog := range programList {
				if ch >= 'a' && ch <= 'z' && i == progId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%c] %s\n", 'a'+i, prog.Name())
			}
			fmt.Printf("---------------------------------------\n")
			fmt.Printf("  S - Stop animation\n")
			fmt.Printf("  C - Continue animation\n")
			fmt.Printf("---------------------------------------\n")

			fmt.Printf("Enter a character (or '0' for quit): ")

			fmt.Scanf("%s\n", &input)
			ch = input[0]
			if ch == '0' {
				break
			}

			if ch >= 'a' && ch <= 'z' {
				progId = int(ch - 'a')
				if progId < 0 || progId >= len(programList) {
					continue
				}
				// animCtrl.Stop()
				fmt.Printf("Program statistics:\n")
				fmt.Printf("  animation: %v\n", animCtrl.Watch())
				fmt.Printf("  painting : %v\n", canvas.Watch())
				fmt.Printf("  sending  : %v\n", pixClient.Watch())
				animCtrl.Purge()
				// animCtrl.Continue()
				canvas.Purge()
				animCtrl.Watch().Reset()
				canvas.Watch().Reset()
				pixClient.Watch().Reset()
				programList[progId].Run(canvas)
			}
			if ch == 'S' {
				animCtrl.Stop()
			}
			if ch == 'C' {
				animCtrl.Continue()
			}
		}
	} else {
		if ch >= 'a' && ch <= 'z' {
			progId = int(ch - 'a')
			if progId >= 0 && progId < len(programList) {
				programList[progId].Run(canvas)
			}
		}
		fmt.Printf("Quit by Ctrl-C\n")
		SignalHandler()
	}

	animCtrl.Stop()
	ledGrid.Clear(color.Black)
	pixClient.Send(ledGrid)
	pixClient.Close()

	fmt.Printf("Program statistics:\n")
	fmt.Printf("  animation: %v\n", animCtrl.Watch())
	fmt.Printf("  painting : %v\n", canvas.Watch())
	fmt.Printf("  sending  : %v\n", pixClient.Watch())
}
