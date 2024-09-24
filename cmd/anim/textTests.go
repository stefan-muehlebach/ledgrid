package main

import (
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
	"golang.org/x/image/math/fixed"
)

var (
	MovingText = NewLedGridProgram("Moving text",
		func(c *ledgrid.Canvas) {
			t1 := ledgrid.NewText(geom.Point{0, float64(height) / 2.0}, "Stefan", color.LightSeaGreen)
			t1.SetAlign(ledgrid.AlignLeft | ledgrid.AlignMiddle)
			t2 := ledgrid.NewText(geom.Point{float64(width), float64(height) / 2.0}, "Beni", color.YellowGreen)
			t2.SetAlign(ledgrid.AlignRight | ledgrid.AlignMiddle)
			t3 := ledgrid.NewText(geom.Point{float64(width) + 60.0, float64(height) / 2.0}, "wohnen im Lochbach", color.OrangeRed)
			t3.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)

			c.Add(t1, t2, t3)

			aAngle1 := ledgrid.NewFloatAnimation(&t1.Angle, -2*math.Pi, 7*time.Second)
			aAngle1.Curve = ledgrid.AnimationLinear
			aAngle1.RepeatCount = ledgrid.AnimationRepeatForever

			aAngle2 := ledgrid.NewFloatAnimation(&t2.Angle, -2*math.Pi, 8*time.Second)
			aAngle2.Curve = ledgrid.AnimationLinear
			aAngle2.RepeatCount = ledgrid.AnimationRepeatForever

			aPos := ledgrid.NewPositionAnimation(&t3.Pos, geom.Point{-100, float64(height) / 2.0}, 6*time.Second)
			aPos.Curve = ledgrid.AnimationEaseInOut

			aTimeline := ledgrid.NewTimeline(15 * time.Second)
			aTimeline.Add(10*time.Second, aPos)
			aTimeline.RepeatCount = ledgrid.AnimationRepeatForever

			aAngle1.Start()
			aAngle2.Start()
			aTimeline.Start()
		})

	FixedFontTest = NewLedGridProgram("Fixed font",
		func(c *ledgrid.Canvas) {
			pos1 := fixed.P(55, height/2)
			pos2 := fixed.P(-15, height/2)
			color1 := color.Aquamarine
			color2 := color.BlueViolet

			txt1 := ledgrid.NewFixedText(pos1, color1, "STEFAN")
			c.Add(txt1)

			aPos := ledgrid.NewFixedPosAnimation(&txt1.Pos, pos2, 10*time.Second)
			aPos.AutoReverse = true
			aPos.RepeatCount = ledgrid.AnimationRepeatForever

			aColor := ledgrid.NewColorAnimation(&txt1.Color, color2, 2*time.Second)
			aColor.AutoReverse = true
			aColor.RepeatCount = ledgrid.AnimationRepeatForever

			aPos.Start()
			aColor.Start()
		})
)
