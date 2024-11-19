package main

import (
	"fmt"
	"golang.org/x/image/math/fixed"
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

func f2f(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}
func p2p(x, y float64) fixed.Point26_6 {
    return fixed.Point26_6{f2f(x), f2f(y)}
}

var (
    ClockAnim = NewLedGridProgram("Clock Animation",
        func(c *ledgrid.Canvas) {
            var digitList[4] *ledgrid.FixedText
            var colFade[4] *ledgrid.ColorAnimation
            var d[4] int

            for i := range 4 {
                pos := p2p(float64(3-i)*float64(width)/4.0+2.0, float64(height)-1.0)
                digit := ledgrid.NewFixedText(pos, color.Blue, "0")
                digit.SetFont(ledgrid.Fixed5x7)
                digitList[i] = digit

                c.Add(digit)
            }

            timeLine1 := ledgrid.NewTimeline(time.Second)
            timeLine1.RepeatCount = ledgrid.AnimationRepeatForever
            timeLine1.Add(0, ledgrid.NewTask(func() {
                d[0] = (d[0]+1) % 10
                if d[0] == 0 {
                    d[1] = (d[1]+1) % 6
                    if d[1] == 0 {
                        d[2] = (d[2]+1) % 10
                        if d[2] == 0 {
                            d[3] = (d[3]+1) % 6
                        }
                    }
                }
                for i := range 4 {
                    ch := fmt.Sprintf("%d", d[i])
                    digitList[i].SetText(ch)
                }
            }))

            digitColor := color.Blue
            for i := range 4 {
                colFade[i] = ledgrid.NewColorAnim(digitList[i], digitColor, 2*time.Second)
                colFade[i].Val2 = func() color.LedColor {
                    return digitColor
                }
                colFade[i].Cont = true
                colFade[i].Curve = ledgrid.AnimationLinear
            }

            timeLine2 := ledgrid.NewTimeline(3 * time.Second)
            timeLine2.RepeatCount = ledgrid.AnimationRepeatForever
            timeLine2.Add(0, ledgrid.NewTask(func() {
                digitColor = color.RandColor()
            }))
            timeLine2.Add(time.Second, colFade[0], colFade[1], colFade[2], colFade[3])

            timeLine1.Start()
            timeLine2.Start()
        })

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
