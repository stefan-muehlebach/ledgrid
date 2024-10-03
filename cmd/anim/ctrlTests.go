package main

import (
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

var (
	GroupTest = NewLedGridProgram("Group test",
		func(c *ledgrid.Canvas) {
			rPos1 := geom.Point{5.0, float64(height) / 2.0}
			rPos2 := geom.Point{float64(width) - 5.0, float64(height) / 2.0}
			rSize1 := geom.Point{7.0, 7.0}
			rSize2 := geom.Point{1.0, 1.0}
			rColor1 := color.SkyBlue
			rColor2 := color.GreenYellow

			r := ledgrid.NewRectangle(rPos1, rSize1, rColor1)
			c.Add(r)

			aPos := ledgrid.NewPositionAnim(r, rPos2, time.Second)
			aPos.AutoReverse = true
			aSize := ledgrid.NewSizeAnim(r, rSize2, 2*time.Second)
			aSize.AutoReverse = true
			aColor := ledgrid.NewColorAnim(r, rColor2, 2*time.Second)
			aColor.AutoReverse = true
			aAngle := ledgrid.NewAngleAnim(r, math.Pi, 2*time.Second)
			aAngle.AutoReverse = true

			aGroup := ledgrid.NewGroup(aPos, aSize, aColor, aAngle)
			aGroup.RepeatCount = ledgrid.AnimationRepeatForever

			aGroup.Start()
		})

	SequenceTest = NewLedGridProgram("Sequence test",
		func(c *ledgrid.Canvas) {
			rPos := geom.NewPointIMG(gridSize).Mul(0.5)
			rSize1 := geom.NewPointIMG(gridSize).SubXY(1, 1)
			rSize4 := geom.Point{5.0, 3.0}
			// rSize2 := rSize1
			// rSize2.X = rSize4.X
			rSize3 := rSize1
			rSize3.Y = rSize4.Y

			r := ledgrid.NewRectangle(rPos, rSize1, color.SkyBlue)
			c.Add(r)

			// aSize2 := ledgrid.NewSizeAnim(r, rSize2, time.Second)
			aColor1 := ledgrid.NewColorAnim(r, color.OrangeRed, time.Second)
			aColor1.Cont = true
			aSize3 := ledgrid.NewSizeAnim(r, rSize3, time.Second)
			aSize3.Cont = true
			aColor2 := ledgrid.NewColorAnim(r, color.YellowGreen, time.Second)
			aColor2.Cont = true
			aSize4 := ledgrid.NewSizeAnim(r, rSize4, time.Second)
			aSize4.Cont = true
			aColor3 := ledgrid.NewColorAnim(r, color.Gold, time.Second)
			aColor3.Cont = true
			aColor4 := ledgrid.NewColorAnim(r, color.MediumOrchid, time.Second)
			aColor4.Cont = true
			aSize1 := ledgrid.NewSizeAnim(r, rSize1, time.Second)
			aSize1.Cont = true
			aColor5 := ledgrid.NewColorAnim(r, color.SkyBlue, time.Second)
			aColor5.Cont = true

			aSeq := ledgrid.NewSequence(aColor1, aSize3, aColor2, aSize4, aColor3, aColor4, aSize1, aColor5)
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

			aAngle1 := ledgrid.NewAngleAnim(r1, math.Pi, time.Second)
			aAngle2 := ledgrid.NewAngleAnim(r1, 0.0, time.Second)
			aAngle2.Cont = true

			aColor1 := ledgrid.NewColorAnim(r1, color.OrangeRed, 200*time.Millisecond)
			aColor1.AutoReverse = true
			aColor1.RepeatCount = 3
			aColor2 := ledgrid.NewColorAnim(r1, color.Purple, 500*time.Millisecond)
			aColor3 := ledgrid.NewColorAnim(r1, color.GreenYellow, 500*time.Millisecond)
			aColor3.Cont = true

			aPos1 := ledgrid.NewPositionAnim(r1, r2.Pos.SubXY(r2Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos1.AutoReverse = true

			aAngle3 := ledgrid.NewAngleAnim(r3, -math.Pi, time.Second)
			aAngle4 := ledgrid.NewAngleAnim(r3, 0.0, time.Second)
			aAngle4.Cont = true

			aColor4 := ledgrid.NewColorAnim(r3, color.DarkOrange, 200*time.Millisecond)
			aColor4.AutoReverse = true
			aColor4.RepeatCount = 3
			aColor5 := ledgrid.NewColorAnim(r3, color.Purple, 500*time.Millisecond)
			aColor6 := ledgrid.NewColorAnim(r3, color.SkyBlue, 500*time.Millisecond)
			aColor6.Cont = true

			aPos2 := ledgrid.NewPositionAnim(r3, r4.Pos.AddXY(r4Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos2.AutoReverse = true

			aColor7 := ledgrid.NewColorAnim(r2, color.Cornsilk, 500*time.Millisecond)
			aColor7.AutoReverse = true
			aBorder1 := ledgrid.NewStrokeWidthAnim(r2, 2.0, 500*time.Millisecond)
			aBorder1.AutoReverse = true

			aColor8 := ledgrid.NewColorAnim(r4, color.Cornsilk, 500*time.Millisecond)
			aColor8.AutoReverse = true
			aBorder2 := ledgrid.NewStrokeWidthAnim(r4, 2.0, 500*time.Millisecond)
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
)
