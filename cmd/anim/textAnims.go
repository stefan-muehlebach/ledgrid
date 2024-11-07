package main

import (
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

var (
	MovingText = NewLedGridProgram("Moving text",
		func(c *ledgrid.Canvas) {
			t1 := ledgrid.NewText(geom.Point{0, float64(height) / 2.0}, "Stefan", color.LightSeaGreen)
			t1.SetAlign(ledgrid.AlignLeft)
			t2 := ledgrid.NewText(geom.Point{float64(width), float64(height) / 2.0}, "Beni", color.YellowGreen)
			t2.SetAlign(ledgrid.AlignRight)

			t4 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "werden", color.Gold)
			t5 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "immer", color.Gold)
			t6 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "im", color.Gold)
			t7 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "Lochbach", color.Gold)
			t8 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "wohnen", color.Gold)

			c.Add(t1, t2, t4, t5, t6, t7, t8)

			aAngle1 := ledgrid.NewAngleAnim(t1, -2*math.Pi, 7*time.Second)
			aAngle1.Curve = ledgrid.AnimationLinear
			aAngle1.RepeatCount = ledgrid.AnimationRepeatForever

			aAngle2 := ledgrid.NewAngleAnim(t2, -2*math.Pi, 8*time.Second)
			aAngle2.Curve = ledgrid.AnimationLinear
			aAngle2.RepeatCount = ledgrid.AnimationRepeatForever

			aPos4 := ledgrid.NewPositionAnim(t4, geom.Point{float64(width) / 2.0, -float64(height) / 2.0}, 4*time.Second)
			aPos4.Curve = ledgrid.AnimationLinear
			aPos5 := ledgrid.NewPositionAnim(t5, geom.Point{float64(width) / 2.0, -float64(height) / 2.0}, 4*time.Second)
			aPos5.Curve = ledgrid.AnimationLinear
			aPos6 := ledgrid.NewPositionAnim(t6, geom.Point{float64(width) / 2.0, -float64(height) / 2.0}, 4*time.Second)
			aPos6.Curve = ledgrid.AnimationLinear
			aPos7 := ledgrid.NewPositionAnim(t7, geom.Point{float64(width) / 2.0, -float64(height) / 2.0}, 4*time.Second)
			aPos7.Curve = ledgrid.AnimationLinear
			aPos8 := ledgrid.NewPositionAnim(t8, geom.Point{float64(width) / 2.0, -float64(height) / 2.0}, 4*time.Second)
			aPos8.Curve = ledgrid.AnimationLinear

			aPosSeq := ledgrid.NewSequence(ledgrid.NewDelay(4*time.Second), aPos4, aPos5, aPos6, aPos7, aPos8)
			// aPosSeq.SetDuration(aPosSeq.Duration() + 4*time.Second)
			aPosSeq.RepeatCount = ledgrid.AnimationRepeatForever

			aAngle1.Start()
			aAngle2.Start()
			aPosSeq.Start()
		})
)
