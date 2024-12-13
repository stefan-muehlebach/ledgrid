// geomShapeAnims.go
package main

import (
	"math"
	"math/rand/v2"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

func init() {
	programList.Add("Circling circles", CirclingCircles)
	programList.Add("Chasing circles", ChasingCircles)
	programList.Add("Circle animation", CircleAnimation)
	programList.Add("Pushing rectangles", PushingRectangles)
	programList.Add("Regular polygons", RegularPolygon)
	programList.Add("Rectangles journey", RectanglesJourney)
	programList.Add("Async multiple color fade", AsyncColorFade)
}

func AsyncColorFade(c *ledgrid.Canvas) {
	var posList []geom.Point
	var objList []*ledgrid.Rectangle
	// var fadeList []*ledgrid.ColorAnimation
	// var seqList []*ledgrid.Sequence
	var numObjs int = 10

	color1 := color.NewLedColorRGB(0x0e, 0x4f, 0x92)
	color2 := color.NewLedColorRGB(0xff, 0x00, 0xff)
	color3 := color.NewLedColorRGB(0xab, 0xff, 0x5e)

	posList = make([]geom.Point, numObjs)
	posList[0] = geom.Point{1.5, float64(height) / 2.0}
	for i := range posList[1:] {
		posList[i+1] = posList[i].AddXY(4.0, 0.0)
	}
	size := geom.Point{2.0, 9.0}

	objList = make([]*ledgrid.Rectangle, numObjs)
	for i := range objList {
		obj := ledgrid.NewRectangle(posList[i], size, color1)
		objList[i] = obj
		c.Add(obj)
	}

	animGrp := ledgrid.NewGroup()
	// fadeList = make([]*ledgrid.ColorAnimation, numObjs)
	// seqList = make([]*ledgrid.Sequence, numObjs)
	for i := range objList {
		fade1 := ledgrid.NewColorAnim(objList[i], color2, 1*time.Second)
		fade1.AutoReverse = true
		fade1.RepeatCount = 6
		fade1.Curve = ledgrid.AnimationLinear
		fade1.Pos = float64(i) * 0.5 / float64(numObjs-1)
		// fadeList[i] = fade1

		fade2 := ledgrid.NewColorAnim(objList[i], color3, 1*time.Second)
		fade2.AutoReverse = true
		fade2.RepeatCount = 6
		fade2.Curve = ledgrid.AnimationLinear

		seq := ledgrid.NewSequence(fade1, fade2)
		seq.RepeatCount = ledgrid.AnimationRepeatForever

		// seqList[i] = seq
		animGrp.Add(seq)
	}
	animGrp.Start()

	// ticker := time.NewTicker(30 * time.Second)
	// go func() {
	// 	for range ticker.C {
	// 		fmt.Printf("\nList of sequences:\n")
	// 		for i, seq := range seqList {
	// 			sStart, sEnd := seq.TimeInfo()
	// 			fStart, fEnd, fTotal := fadeList[i].TimeInfo()
	// 			fmt.Printf("[%d]\n", i)
	// 			fmt.Printf("  Duration: %v\n", seq.Duration())
	// 			fmt.Printf("  Start: %v; End: %v\n", sStart, sEnd)
	// 			fmt.Printf("  Start: %v; End: %v\n", fStart, fEnd)
	// 			fmt.Printf("  Total: %f\n", fTotal)
	// 		}
	// 	}
	// }()
}

func CirclingCircles(c *ledgrid.Canvas) {
	pos1 := geom.Point{1.5, 1.5}
	pos2 := geom.Point{10.5, float64(height) - 1.5}
	pos3 := geom.Point{19.5, 1.5}
	pos4 := geom.Point{28.5, float64(height) - 1.5}
	pos5 := geom.Point{float64(width) - 1.5, 1.5}
	cSize := geom.Point{2.0, 2.0}

	c1 := ledgrid.NewEllipse(pos1, cSize, color.OrangeRed)
	c2 := ledgrid.NewEllipse(pos2, cSize, color.MediumSeaGreen)
	c3 := ledgrid.NewEllipse(pos3, cSize, color.SkyBlue)
	c4 := ledgrid.NewEllipse(pos4, cSize, color.Gold)
	c5 := ledgrid.NewEllipse(pos5, cSize, color.YellowGreen)

	stepRD := geom.Point{18.0, 2.0 * (float64(height) - 3.0)}
	stepLU := stepRD.Neg()
	stepRU := geom.Point{18.0, -2.0 * (float64(height) - 3.0)}
	stepLD := stepRU.Neg()

	quartCirc := ledgrid.CirclePath.NewStartLen(0, 0.25)

	c1Path1 := ledgrid.NewPathAnim(c1, quartCirc, stepRD, time.Second)
	c1Path2 := ledgrid.NewPathAnim(c1, quartCirc, stepRU, time.Second)
	c1Path3 := ledgrid.NewPathAnim(c1, quartCirc, stepLD, time.Second)
	c1Path4 := ledgrid.NewPathAnim(c1, quartCirc, stepLU, time.Second)

	c2Path1 := ledgrid.NewPathAnim(c2, quartCirc, stepLU, time.Second)
	c2Path2 := ledgrid.NewPathAnim(c2, quartCirc, stepRD, time.Second)

	c3Path1 := ledgrid.NewPathAnim(c3, quartCirc, stepLD, time.Second)
	c3Path2 := ledgrid.NewPathAnim(c3, quartCirc, stepRU, time.Second)

	c4Path1 := ledgrid.NewPathAnim(c4, quartCirc, stepLU, time.Second)
	c4Path2 := ledgrid.NewPathAnim(c4, quartCirc, stepRD, time.Second)

	c5Path1 := ledgrid.NewPathAnim(c5, quartCirc, stepLD, time.Second)
	c5Path2 := ledgrid.NewPathAnim(c5, quartCirc, stepRU, time.Second)

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
}

func ChasingCircles(c *ledgrid.Canvas) {
	c1Pos1 := geom.Point{float64(width) - 5.0, float64(height) / 2.0}
	c1Size1 := geom.Point{9.0, 9.0}
	c1Size2 := geom.Point{3.0, 3.0}
	c1PosSize := geom.Point{float64(width - 10), -float64(height) / 2.0}
	c2Pos := geom.Point{5.0, float64(height) / 2.0}
	c2Size1 := geom.Point{5.0, 5.0}
	c2Size2 := geom.Point{3.0, 3.0}
	c2PosSize := geom.Point{-float64(width - 10), float64(height)/2.0 + 2.0}

	aGrp := ledgrid.NewGroup()

	pal := ledgrid.NewGradientPaletteByList("Palette", true,
		color.DeepSkyBlue,
		color.Lime,
		color.Teal,
		color.SkyBlue,
	)

	c1 := ledgrid.NewEllipse(c1Pos1, c1Size1, color.Gold)

	path := ledgrid.CirclePath.NewStart(0.25)

	c1pos := ledgrid.NewPathAnim(c1, path, c1PosSize, 4*time.Second)
	c1pos.Curve = ledgrid.AnimationLinear

	c1size := ledgrid.NewSizeAnim(c1, c1Size2, time.Second)
	c1size.AutoReverse = true

	c1bcolor := ledgrid.NewColorAnim(c1, color.OrangeRed, time.Second)
	c1bcolor.AutoReverse = true

	c2 := ledgrid.NewEllipse(c2Pos, c2Size1, color.Lime)

	c2pos := ledgrid.NewPathAnim(c2, path, c2PosSize, 4*time.Second)
	c2pos.Curve = ledgrid.AnimationLinear

	c2size := ledgrid.NewSizeAnim(c2, c2Size2, time.Second)
	c2size.AutoReverse = true

	c2color := ledgrid.NewPaletteAnim(c2, pal, 2*time.Second)

	aGrp.Add(c1pos, c1size, c1bcolor, c2pos, c2size, c2color)
	aGrp.RepeatCount = ledgrid.AnimationRepeatForever

	c.Add(c1, c2)
	aGrp.Start()
}

func CircleAnimation(c *ledgrid.Canvas) {
	c1Pos1 := geom.Point{2.0, float64(height) / 2.0}
	c1Pos3 := geom.Point{float64(width) - 2.0, float64(height) / 2.0}

	c1Size1 := geom.Point{3.0, 3.0}
	c1Size2 := geom.Point{9.0, 9.0}

	c1 := ledgrid.NewEllipse(c1Pos1, c1Size1, color.OrangeRed)

	c1pos := ledgrid.NewPositionAnim(c1, c1Pos3, 2*time.Second)
	c1pos.AutoReverse = true
	c1pos.RepeatCount = ledgrid.AnimationRepeatForever

	c1radius := ledgrid.NewSizeAnim(c1, c1Size2, time.Second)
	c1radius.AutoReverse = true
	c1radius.RepeatCount = ledgrid.AnimationRepeatForever

	c1color := ledgrid.NewColorAnim(c1, color.Gold, 4*time.Second)
	c1color.AutoReverse = true
	c1color.RepeatCount = ledgrid.AnimationRepeatForever

	c.Add(c1)

	c1pos.Start()
	c1radius.Start()
	c1color.Start()
}

func PushingRectangles(c *ledgrid.Canvas) {
	r1Pos1 := geom.Point{1.0, float64(height) / 2.0}
	r1Pos2 := geom.Point{0.5 + float64(width-3)/2.0, float64(height) / 2.0}

	r2Pos1 := geom.Point{float64(width - 1), float64(height) / 2.0}
	r2Pos2 := geom.Point{float64(width) - 0.5 - float64(width-3)/2.0, float64(height) / 2.0}

	rSize1 := geom.Point{float64(width - 3), 1.0}
	rSize2 := geom.Point{1.0, float64(height - 1)}

	duration := 2 * time.Second

	r1 := ledgrid.NewRectangle(r1Pos1, rSize2, color.Crimson)

	a1Pos := ledgrid.NewPositionAnim(r1, r1Pos2, duration)
	a1Pos.AutoReverse = true
	a1Pos.RepeatCount = ledgrid.AnimationRepeatForever

	a1Size := ledgrid.NewSizeAnim(r1, rSize1, duration)
	a1Size.AutoReverse = true
	a1Size.RepeatCount = ledgrid.AnimationRepeatForever

	a1Color := ledgrid.NewColorAnim(r1, color.GreenYellow, duration)
	a1Color.AutoReverse = true
	a1Color.RepeatCount = ledgrid.AnimationRepeatForever

	r2 := ledgrid.NewRectangle(r2Pos2, rSize1, color.SkyBlue)

	a2Pos := ledgrid.NewPositionAnim(r2, r2Pos1, duration)
	a2Pos.AutoReverse = true
	a2Pos.RepeatCount = ledgrid.AnimationRepeatForever

	a2Size := ledgrid.NewSizeAnim(r2, rSize2, duration)
	a2Size.AutoReverse = true
	a2Size.RepeatCount = ledgrid.AnimationRepeatForever

	a2Color := ledgrid.NewColorAnim(r2, color.Crimson, duration)
	a2Color.AutoReverse = true
	a2Color.RepeatCount = ledgrid.AnimationRepeatForever

	c.Add(r1, r2)
	a1Pos.Start()
	a1Size.Start()
	a1Color.Start()
	a2Pos.Start()
	a2Size.Start()
	a2Color.Start()
}

func RegularPolygon(c *ledgrid.Canvas) {
	posList := []geom.Point{
		geom.Point{-6.0, float64(height) / 2.0},
		geom.Point{float64(width) + 5.0, float64(height) / 2.0},
	}
	posCenter := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
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
		animPos := ledgrid.NewPositionAnim(polyList[n], posCenter, dur)
		animAngle := ledgrid.NewAngleAnim(polyList[n], angle, dur)
		animSize := ledgrid.NewSizeAnim(polyList[n], largeSize, 4*time.Second)
		animFade := ledgrid.NewColorAnim(polyList[n], color.Black, 4*time.Second)

		aGrpIn := ledgrid.NewGroup(animPos, animAngle)
		aGrpOut := ledgrid.NewGroup(animSize, animFade)
		aObjSeq := ledgrid.NewSequence(aGrpIn, aGrpOut)
		aSeq.Add(aObjSeq)
	}
	aSeq.RepeatCount = ledgrid.AnimationRepeatForever
	aSeq.Start()
}

func FlyingRectangle(c *ledgrid.Canvas) {
	r1Pos1 := geom.Point{4, float64(height) / 2.0}
	r1Pos2 := geom.Point{float64(width) + 4.0, float64(height) / 2.0}
	r1Size := geom.Point{3.0, 7.0}

	r2Pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	r2Size1 := geom.Point{3.0, 7.0}
	r2Size2 := geom.Point{7.0, 9.0}
	r4Pos := geom.Point{float64(width) - 4.0, float64(height) / 2.0}
	r4Size1 := geom.Point{3.0, 7.0}
	r4Size2 := geom.Point{7.0, 9.0}

	r1 := ledgrid.NewRectangle(r1Pos1, r1Size, color.GreenYellow)

	r2 := ledgrid.NewEllipse(r2Pos, r2Size1, color.Gold)
	r4 := ledgrid.NewEllipse(r4Pos, r4Size1, color.Gold)
	c.Add(r1, r2, r4)

	aAngle1 := ledgrid.NewAngleAnim(r1, 0.5*math.Pi, time.Second)
	aAngle2 := ledgrid.NewAngleAnim(r1, 0.0, time.Second)

	aColor1 := ledgrid.NewColorAnim(r1, color.OrangeRed, time.Second)
	aColor2 := ledgrid.NewColorAnim(r1, color.Purple, 500*time.Millisecond)
	aColor3 := ledgrid.NewColorAnim(r1, color.GreenYellow, 500*time.Millisecond)

	aPos1 := ledgrid.NewPositionAnim(r1, r1Pos2, 1000*time.Millisecond)
	aPos1.AutoReverse = true

	aBorder1 := ledgrid.NewStrokeWidthAnim(r2, 2.0, 300*time.Millisecond)
	aBorder1.AutoReverse = true
	aSize2 := ledgrid.NewSizeAnim(r2, r2Size2, 300*time.Millisecond)
	// aSize3 := ledgrid.NewSizeAnim(r2, r2Size1, 300*time.Millisecond)
	// aSize3.Cont = true
	aSize2.AutoReverse = true

	// aColor8 := ledgrid.NewColorAnim(r4, color.Cornsilk, 500*time.Millisecond)
	// aColor8.AutoReverse = true
	aBorder2 := ledgrid.NewStrokeWidthAnim(r4, 2.0, 300*time.Millisecond)
	aBorder2.AutoReverse = true
	aSize4 := ledgrid.NewSizeAnim(r4, r4Size2, 300*time.Millisecond)
	// aSize5 := ledgrid.NewSizeAnim(r4, r4Size1, 300*time.Millisecond)
	// aSize5.Cont = true
	aSize4.AutoReverse = true

	tl := ledgrid.NewTimeline(6 * time.Second)
	tl.RepeatCount = ledgrid.AnimationRepeatForever

	tl.Add(300*time.Millisecond, aColor1)
	tl.Add(500*time.Millisecond, aAngle1)
	tl.Add(2300*time.Millisecond, aPos1)
	tl.Add(2500*time.Millisecond, aColor2, aBorder1, aSize2)
	tl.Add(2700*time.Millisecond, aBorder2, aSize4)
	tl.Add(3500*time.Millisecond, aBorder2, aSize4)
	tl.Add(3700*time.Millisecond, aBorder1, aSize2)
	tl.Add(3900*time.Millisecond, aColor3)
	tl.Add(4400*time.Millisecond, aAngle2)

	tl.Start()
}

func RectanglesJourney(c *ledgrid.Canvas) {
	var posList [3]geom.Point
	var animList [3]*ledgrid.PathAnimation
	var dotList [3]*ledgrid.Dot

	r1Pos1 := geom.Point{4.0, float64(height) / 2.0}
	r1Pos2 := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	r1Pos3 := geom.Point{float64(width) - 4.0, float64(height) / 2.0}
	r1Size := geom.Point{7.0, 3.0}

	// c1Pos1 := geom.Point{2.0 * float64(width) / 3.0, float64(height) / 4.0}
	// c1Pos2 := c1Pos1.SubXY(float64(width)+4, 0.0)
	// c2Pos1 := geom.Point{float64(width) / 3.0, 3.0 * float64(height) / 4.0}
	// c2Pos2 := c2Pos1.SubXY(float64(width)+4, 0.0)
	// c3Pos1 := geom.Point{3.0 * float64(width) / 4.0, float64(height) / 2.0}
	// c3Pos2 := c3Pos1.SubXY(float64(width)+4, 0.0)

	// c4Pos1 := geom.Point{float64(width) + 2.0, float64(height) / 4.0}
	// c4Pos2 := c4Pos1.SubXY(float64(width)+4, 0.0)

	r1 := ledgrid.NewRectangle(r1Pos1, r1Size, color.GreenYellow)

	// c1 := ledgrid.NewDot(c1Pos1, color.LightBlue)
	// c2 := ledgrid.NewDot(c2Pos1, color.LightBlue)
	// c3 := ledgrid.NewDot(c3Pos1, color.LightBlue)
	dotList[0] = ledgrid.NewDot(geom.Point{float64(width + 2), 0.0}, color.LightBlue)
	dotList[1] = ledgrid.NewDot(geom.Point{float64(width + 2), 0.0}, color.LightBlue.Dark(0.2))
	dotList[2] = ledgrid.NewDot(geom.Point{float64(width + 2), 0.0}, color.LightBlue.Dark(0.4))

	// randTask := ledgrid.NewTask(func() {
	// 	pos1 := geom.Point{float64(width) + 2.0, float64(height) * rand.Float64()}
	// 	pos2 := pos1.SubXY(float64(width)+4.0, 0.0)
	// 	dot := ledgrid.NewDot(pos1, color.LightBlue)
	// 	c.Add(dot)
	// 	posAnim := ledgrid.NewPositionAnim(dot, pos2, time.Second)
	// 	posAnim.Curve = ledgrid.AnimationLinear
	// 	posAnim.Start()
	// })

	c.Add(dotList[2], dotList[1], dotList[0], r1)

	aPos2 := ledgrid.NewPositionAnim(r1, r1Pos2, 1000*time.Millisecond)
	aPos2.Curve = ledgrid.AnimationEaseIn
	aPos2.Cont = false
	aPos3 := ledgrid.NewPositionAnim(r1, r1Pos3, 1000*time.Millisecond)
	aPos3.Curve = ledgrid.AnimationEaseOut

	for i := range 3 {
		animList[i] = ledgrid.NewPositionAnim(dotList[i], geom.Point{}, time.Duration(i+1)*time.Second)
		animList[i].Curve = ledgrid.AnimationLinear
		animList[i].Val1 = func() geom.Point {
			posList[i] = geom.Point{float64(width + 2), float64(rand.IntN(10))}
			return posList[i]
		}
		animList[i].Val2 = func() geom.Point {
			return posList[i].SubXY(float64(width+5), 0.0)
		}
		animList[i].Cont = false
	}

	tl := ledgrid.NewTimeline(5 * time.Second)
	tl.RepeatCount = 1

	// tl.Add(1000*time.Millisecond, aPos2)
	// tl.Add(2000*time.Millisecond, posAnim, aPos5, aPos6)
	tl.Add(0000*time.Millisecond, animList[0])
	tl.Add(1100*time.Millisecond, animList[0])
	tl.Add(2500*time.Millisecond, animList[0])
	tl.Add(3900*time.Millisecond, animList[0])

	tl.Add(900*time.Millisecond, animList[1])
	tl.Add(3000*time.Millisecond, animList[1])

	tl.Add(200*time.Millisecond, animList[2])

	// tl.Add(3500*time.Millisecond, randTask)
	// tl.Add(4200*time.Millisecond, randTask)
	// tl.Add(9000*time.Millisecond, aPos3)

	seq := ledgrid.NewSequence(aPos2, tl, aPos3)

	seq.Start()
}

//----------------------------------------------------------------------------

// type BouncingEllipse struct {
// 	ledgrid.Ellipse
// 	Vel, Acc geom.Point
// 	Field    geom.Rectangle
// }

// func NewBouncingEllipse(pos, size geom.Point, col color.LedColor) *BouncingEllipse {
// 	b := &BouncingEllipse{}
// 	b.Ellipse = *ledgrid.NewEllipse(pos, size, col)
// 	b.Vel = geom.Point{}
// 	b.Acc = geom.Point{}
// 	return b
// }

// func (b *BouncingEllipse) Update(pit time.Time) bool {
// 	deltaVel := b.Acc.Mul(0.3)
// 	b.Vel = b.Vel.Add(deltaVel)
// 	b.Pos = b.Pos.Add(b.Vel)
// 	if b.Pos.X < b.Field.Min.X || b.Pos.X >= b.Field.Max.X {
// 		b.Vel.X = -b.Vel.X
// 	}
// 	if b.Pos.Y < b.Field.Min.Y || b.Pos.Y >= b.Field.Max.Y {
// 		b.Vel.Y = -b.Vel.Y
// 	}
// 	return true
// }

// func (b *BouncingEllipse) Duration() time.Duration {
// 	return time.Duration(0)
// }
// func (b *BouncingEllipse) SetDuration(dur time.Duration) {}
// func (b *BouncingEllipse) Start()                        {}
// func (b *BouncingEllipse) Suspend()                      {}
// func (b *BouncingEllipse) Continue()                     {}
// func (b *BouncingEllipse) IsRunning() bool {
// 	return true
// }

// func BounceAround(c *ledgrid.Canvas) {
// 	pos1 := geom.Point{2.0, 2.0}
// 	pos2 := geom.Point{37.0, 7.0}
// 	size := geom.Point{4.0, 4.0}
// 	vel1 := geom.Point{0.15, 0.075}
// 	vel2 := geom.Point{-0.35, -0.25}

// 	obj1 := NewBouncingEllipse(pos1, size, color.GreenYellow)
// 	obj1.Vel = vel1
// 	obj1.Field = geom.NewRectangleIMG(c.Rect)
// 	obj2 := NewBouncingEllipse(pos2, size, color.LightSeaGreen)
// 	obj2.Vel = vel2
// 	obj2.Field = geom.NewRectangleIMG(c.Rect)

// 	c.Add(obj1, obj2)
// 	ledgrid.AnimCtrl.Add(0, obj1, obj2)
// }
