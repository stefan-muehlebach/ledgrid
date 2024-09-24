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

			aPos := ledgrid.NewPositionAnimation(&r.Pos, rPos2, time.Second)
			aPos.AutoReverse = true
			aSize := ledgrid.NewSizeAnimation(&r.Size, rSize2, 2*time.Second)
			aSize.AutoReverse = true
			aColor := ledgrid.NewColorAnimation(&r.StrokeColor, rColor2, 2*time.Second)
			aColor.AutoReverse = true
			aAngle := ledgrid.NewFloatAnimation(&r.Angle, math.Pi, 2*time.Second)
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
			aColor1 := ledgrid.NewColorAnimation(&r.StrokeColor, color.OrangeRed, time.Second/2)
			aColor1.AutoReverse = true
			aColor2 := ledgrid.NewColorAnimation(&r.StrokeColor, color.Crimson, time.Second/2)
			aColor2.AutoReverse = true
			aColor3 := ledgrid.NewColorAnimation(&r.StrokeColor, color.Coral, time.Second/2)
			aColor3.AutoReverse = true
			aColor4 := ledgrid.NewColorAnimation(&r.StrokeColor, color.FireBrick, time.Second/2)
			aSize2 := ledgrid.NewSizeAnimation(&r.Size, rSize1, time.Second)
			aSize2.Cont = true
			aColor5 := ledgrid.NewColorAnimation(&r.StrokeColor, color.SkyBlue, time.Second)
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

			aColor1 := ledgrid.NewColorAnimation(&r1.StrokeColor, color.OrangeRed, 200*time.Millisecond)
			aColor1.AutoReverse = true
			aColor1.RepeatCount = 3
			aColor2 := ledgrid.NewColorAnimation(&r1.StrokeColor, color.Purple, 500*time.Millisecond)
			aColor3 := ledgrid.NewColorAnimation(&r1.StrokeColor, color.GreenYellow, 500*time.Millisecond)
			aColor3.Cont = true

			aPos1 := ledgrid.NewPositionAnimation(&r1.Pos, r2.Pos.SubXY(r2Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos1.AutoReverse = true

			aAngle3 := ledgrid.NewFloatAnimation(&r3.Angle, -math.Pi, time.Second)
			aAngle4 := ledgrid.NewFloatAnimation(&r3.Angle, 0.0, time.Second)
			aAngle4.Cont = true

			aColor4 := ledgrid.NewColorAnimation(&r3.StrokeColor, color.DarkOrange, 200*time.Millisecond)
			aColor4.AutoReverse = true
			aColor4.RepeatCount = 3
			aColor5 := ledgrid.NewColorAnimation(&r3.StrokeColor, color.Purple, 500*time.Millisecond)
			aColor6 := ledgrid.NewColorAnimation(&r3.StrokeColor, color.SkyBlue, 500*time.Millisecond)
			aColor6.Cont = true

			aPos2 := ledgrid.NewPositionAnimation(&r3.Pos, r4.Pos.AddXY(r4Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos2.AutoReverse = true

			aColor7 := ledgrid.NewColorAnimation(&r2.StrokeColor, color.Cornsilk, 500*time.Millisecond)
			aColor7.AutoReverse = true
			aBorder1 := ledgrid.NewFloatAnimation(&r2.StrokeWidth, 2.0, 500*time.Millisecond)
			aBorder1.AutoReverse = true

			aColor8 := ledgrid.NewColorAnimation(&r4.StrokeColor, color.Cornsilk, 500*time.Millisecond)
			aColor8.AutoReverse = true
			aBorder2 := ledgrid.NewFloatAnimation(&r4.StrokeWidth, 2.0, 500*time.Millisecond)
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
