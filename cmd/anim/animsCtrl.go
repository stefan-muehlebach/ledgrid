package main

import (
	"context"
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

func init() {
	// programList.AddTitle("Animation Controllers")
	programList.Add("Group test", "Controllers", GroupTest)
	programList.Add("Sequence test", "Controllers", SequenceTest)
	programList.Add("Timeline test", "Controllers", TimelineTest)
}

func GroupTest(ctx context.Context, c *ledgrid.Canvas) {
	rPos1 := geom.Point{5.0, float64(height) / 2.0}
	rPos2 := geom.Point{float64(width) - 5.0, float64(height) / 2.0}
	rSize1 := geom.Point{7.0, 7.0}
	rSize2 := geom.Point{1.0, 1.0}
	rColor1 := color.SkyBlue
	rColor2 := color.GreenYellow

	r := ledgrid.NewRectangle(rPos1, rSize1, rColor1)
	c.Add(r)

	aPos := ledgrid.NewPositionAnim(r, rPos2, 3000*time.Millisecond)
	aPos.AutoReverse = true
	aSize := ledgrid.NewSizeAnim(r, rSize2, 3000*time.Millisecond)
	aSize.AutoReverse = true
	aColor := ledgrid.NewColorAnim(r, rColor2, 2200*time.Millisecond)
	aColor.AutoReverse = true
	aAngle := ledgrid.NewAngleAnim(r, math.Pi, 2500*time.Millisecond)
	aAngle.AutoReverse = true

	aGroup := ledgrid.NewGroup(aPos, aSize, aColor, aAngle)
	aGroup.RepeatCount = ledgrid.AnimationRepeatForever

	aGroup.Start()
}

func SequenceTest(ctx context.Context, c *ledgrid.Canvas) {
	rPos := geom.NewPointIMG(gridSize).Mul(0.5)
	sizeList := [4]geom.Point{
		geom.NewPointIMG(gridSize).SubXY(1, 1),
		geom.NewPoint(5.0, float64(gridSize.Y-1)),
		geom.NewPoint(5.0, 3.0),
		geom.NewPoint(float64(gridSize.X-1), 3.0),
	}
	colorList := [4]color.LedColor{
		color.SkyBlue,
		color.OrangeRed,
		color.Gold,
		color.MediumOrchid,
	}
	sizeAnims := [4]*ledgrid.SizeAnimation{}
	colorAnims := [4]*ledgrid.ColorAnimation{}

	// rSize1 := geom.NewPointIMG(gridSize).SubXY(1, 1)
	// rSize4 := geom.Point{5.0, 3.0}
	// rSize2 := rSize1
	// rSize2.X = rSize4.X
	// rSize3 := rSize1
	// rSize3.Y = rSize4.Y

	r := ledgrid.NewRectangle(rPos, sizeList[0], colorList[0])
	c.Add(r)

	for i, size := range sizeList {
		sizeAnims[i] = ledgrid.NewSizeAnim(r, size, time.Second)
	}
	for i, color := range colorList {
		colorAnims[i] = ledgrid.NewColorAnim(r, color, time.Second)
	}
	angleAnim := ledgrid.NewAngleAnim(r, math.Pi, time.Second)

	aSeq := ledgrid.NewSequence()
	for i := range colorAnims {
		j := (i + 1) % len(colorAnims)
		aSeq.Add(colorAnims[j], sizeAnims[j])
		if i == 1 {
			aSeq.Add(angleAnim)
		}
	}
	aSeq.RepeatCount = ledgrid.AnimationRepeatForever
	aSeq.Start()
}

func TimelineTest(ctx context.Context, c *ledgrid.Canvas) {
	r1Pos := geom.Point{4, float64(height) / 2.0}
	r1Size := geom.Point{3.0, 7.0}
	r3Pos := geom.Point{float64(width) - 4, float64(height) / 2.0}
	r3Size := geom.Point{3.0, 7.0}

	r2Pos := geom.Point{float64(width)/2.0 - 5.0, float64(height) / 2.0}
	r2Size1 := geom.Point{3.0, 7.0}
	r2Size2 := geom.Point{7.0, 9.0}
	r4Pos := geom.Point{float64(width)/2.0 + 5.0, float64(height) / 2.0}
	r4Size1 := geom.Point{3.0, 7.0}
	r4Size2 := geom.Point{7.0, 9.0}

	r1 := ledgrid.NewRectangle(r1Pos, r1Size, color.GreenYellow)
	r3 := ledgrid.NewRectangle(r3Pos, r3Size, color.SkyBlue)

	r2 := ledgrid.NewEllipse(r2Pos, r2Size1, color.Gold)
	r4 := ledgrid.NewEllipse(r4Pos, r4Size1, color.Gold)
	c.Add(r1, r3, r2, r4)

	aAngle1 := ledgrid.NewAngleAnim(r1, 0.5*math.Pi, time.Second)
	aAngle2 := ledgrid.NewAngleAnim(r1, 0.0, time.Second)

	aColor1 := ledgrid.NewColorAnim(r1, color.OrangeRed, 200*time.Millisecond)
	aColor1.AutoReverse = true
	aColor1.RepeatCount = 3
	aColor2 := ledgrid.NewColorAnim(r1, color.Purple, 500*time.Millisecond)
	aColor3 := ledgrid.NewColorAnim(r1, color.GreenYellow, 500*time.Millisecond)

	aPos1 := ledgrid.NewPositionAnim(r1, r2.Pos.SubXY(r2Size1.X/2.0, 0.0), 500*time.Millisecond)
	aPos1.AutoReverse = true

	aAngle3 := ledgrid.NewAngleAnim(r3, -0.5*math.Pi, time.Second)
	aAngle4 := ledgrid.NewAngleAnim(r3, 0.0, time.Second)

	aColor4 := ledgrid.NewColorAnim(r3, color.DarkOrange, 200*time.Millisecond)
	aColor4.AutoReverse = true
	aColor4.RepeatCount = 3
	aColor5 := ledgrid.NewColorAnim(r3, color.Purple, 500*time.Millisecond)
	aColor6 := ledgrid.NewColorAnim(r3, color.SkyBlue, 500*time.Millisecond)

	aPos2 := ledgrid.NewPositionAnim(r3, r4.Pos.AddXY(r4Size1.X/2.0, 0.0), 500*time.Millisecond)
	aPos2.AutoReverse = true

	aColor7 := ledgrid.NewColorAnim(r2, color.Cornsilk, 500*time.Millisecond)
	aColor7.AutoReverse = true
	aBorder1 := ledgrid.NewStrokeWidthAnim(r2, 2.0, 500*time.Millisecond)
	aBorder1.AutoReverse = true
	aSize2 := ledgrid.NewSizeAnim(r2, r2Size2, 500*time.Millisecond)
	aSize2.AutoReverse = true

	aColor8 := ledgrid.NewColorAnim(r4, color.Cornsilk, 500*time.Millisecond)
	aColor8.AutoReverse = true
	aBorder2 := ledgrid.NewStrokeWidthAnim(r4, 2.0, 500*time.Millisecond)
	aBorder2.AutoReverse = true
	aSize4 := ledgrid.NewSizeAnim(r4, r4Size2, 500*time.Millisecond)
	aSize4.AutoReverse = true

	tl := ledgrid.NewTimeline(6 * time.Second)
	tl.RepeatCount = ledgrid.AnimationRepeatForever

	// Timeline positions for the first rectangle
	tl.Add(300*time.Millisecond, aColor1)
	tl.Add(500*time.Millisecond, aAngle1)
	tl.Add(2300*time.Millisecond, aPos1)
	tl.Add(2400*time.Millisecond, aColor2, aColor7, aBorder1, aSize2)
	tl.Add(3900*time.Millisecond, aColor3)
	tl.Add(4400*time.Millisecond, aAngle2)

	// Timeline positions for the second rectangle
	tl.Add(500*time.Millisecond, aColor4)
	tl.Add(700*time.Millisecond, aAngle3)
	tl.Add(2500*time.Millisecond, aPos2)
	tl.Add(2600*time.Millisecond, aColor5, aColor8, aBorder2, aSize4)
	tl.Add(4100*time.Millisecond, aColor6)
	tl.Add(4600*time.Millisecond, aAngle4)

	tl.Start()
}
