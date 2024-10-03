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

var (
	CirclingCircles = NewLedGridProgram("Circling circles",
		func(c *ledgrid.Canvas) {
			pos1 := geom.Point{1.5, 1.5}
			pos2 := geom.Point{10.5, float64(height) - 1.5}
			pos3 := geom.Point{19.5, 1.5}
			pos4 := geom.Point{28.5, float64(height) - 1.5}
			pos5 := geom.Point{37.5, 1.5}
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
			c1Path1.Cont = true
			c1Path2 := ledgrid.NewPathAnim(c1, quartCirc, stepRU, time.Second)
			c1Path2.Cont = true
			c1Path3 := ledgrid.NewPathAnim(c1, quartCirc, stepLD, time.Second)
			c1Path3.Cont = true
			c1Path4 := ledgrid.NewPathAnim(c1, quartCirc, stepLU, time.Second)
			c1Path4.Cont = true

			c2Path1 := ledgrid.NewPathAnim(c2, quartCirc, stepLU, time.Second)
			c2Path1.Cont = true
			c2Path2 := ledgrid.NewPathAnim(c2, quartCirc, stepRD, time.Second)
			c2Path2.Cont = true

			c3Path1 := ledgrid.NewPathAnim(c3, quartCirc, stepLD, time.Second)
			c3Path1.Cont = true
			c3Path2 := ledgrid.NewPathAnim(c3, quartCirc, stepRU, time.Second)
			c3Path2.Cont = true

			c4Path1 := ledgrid.NewPathAnim(c4, quartCirc, stepLU, time.Second)
			c4Path1.Cont = true
			c4Path2 := ledgrid.NewPathAnim(c4, quartCirc, stepRD, time.Second)
			c4Path2.Cont = true

			c5Path1 := ledgrid.NewPathAnim(c5, quartCirc, stepLD, time.Second)
			c5Path1.Cont = true
			c5Path2 := ledgrid.NewPathAnim(c5, quartCirc, stepRU, time.Second)
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

			c2color := ledgrid.NewPaletteAnimation(&c2.Color, pal, 2*time.Second)
			c2color.Curve = ledgrid.AnimationLinear

			aGrp.Add(c1pos, c1size, c1bcolor, c2pos, c2size, c2color)
			aGrp.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(c1, c2)
			aGrp.Start()
		})

	CircleAnimation = NewLedGridProgram("Circle animation",
		func(c *ledgrid.Canvas) {
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
		})

	PushingRectangles = NewLedGridProgram("Pushing rectangles",
		func(c *ledgrid.Canvas) {
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
		})

	RegularPolygonTest = NewLedGridProgram("Regular Polygon test",
		func(c *ledgrid.Canvas) {
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
				animSize.Cont = true
				animFade := ledgrid.NewColorAnim(polyList[n], color.Black, 4*time.Second)
				animFade.Cont = true

				aGrpIn := ledgrid.NewGroup(animPos, animAngle)
				aGrpOut := ledgrid.NewGroup(animSize, animFade)
				aObjSeq := ledgrid.NewSequence(aGrpIn, aGrpOut)
				aSeq.Add(aObjSeq)
			}
			aSeq.RepeatCount = ledgrid.AnimationRepeatForever
			aSeq.Start()
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
func (b *BouncingEllipse) Suspend()                      {}
func (b *BouncingEllipse) Continue()                     {}
func (b *BouncingEllipse) IsRunning() bool {
	return true
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
	ledgrid.AnimCtrl.Add(obj1, obj2)
}
