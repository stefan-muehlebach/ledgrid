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
			c.Add(0, r)

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
			sizeList := [4]geom.Point{
				geom.NewPointIMG(gridSize).SubXY(1, 1),
				geom.NewPoint(5.0, float64(gridSize.Y-1)),
				geom.NewPoint(5.0, 3.0),
				geom.NewPoint(float64(gridSize.X-1), 3.0),
			}
			sizeAnims := [4]*ledgrid.PathAnimation{}
			colorList := [4]color.LedColor{
				color.SkyBlue,
				color.OrangeRed,
				color.Gold,
				color.MediumOrchid,
			}
			colorAnims := [4]*ledgrid.ColorAnimation{}

			rSize1 := geom.NewPointIMG(gridSize).SubXY(1, 1)
			rSize4 := geom.Point{5.0, 3.0}
			rSize2 := rSize1
			rSize2.X = rSize4.X
			rSize3 := rSize1
			rSize3.Y = rSize4.Y

			r := ledgrid.NewRectangle(rPos, sizeList[0], colorList[0])
			c.Add(0, r)

			for i, size := range sizeList {
				sizeAnims[i] = ledgrid.NewSizeAnim(r, size, time.Second)
				sizeAnims[i].Cont = true
			}
			for i, color := range colorList {
				colorAnims[i] = ledgrid.NewColorAnim(r, color, time.Second)
				colorAnims[i].Cont = true
			}

			aSeq := ledgrid.NewSequence()
			for i := range colorAnims {
				j := (i + 1) % len(colorAnims)
				aSeq.Add(colorAnims[j], sizeAnims[j])
			}
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
			c.Add(0, r1, r3, r2, r4)

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
