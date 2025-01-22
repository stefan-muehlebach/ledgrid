package main

import (
	"context"
	"image"
	"math"
	"time"

	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

func init() {
    programList.AddTitle("Text Animations")
	programList.Add("Clock animation", ClockAnimation)
	programList.Add("Moving text", MovingText)
	programList.Add("Named colors", NamedColors)
}

func f2f(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}
func p2p(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{f2f(x), f2f(y)}
}

func ClockAnimation(ctx context.Context, c *ledgrid.Canvas) {
	var clockText *ledgrid.FixedText
	var colorFade *ledgrid.ColorAnimation
	var binDigits [32]*ledgrid.Pixel

	pos1 := p2p(5.0, 6.0)
	clockText = ledgrid.NewFixedText(pos1, color.Blue, "00:00:00")
	c.Add(clockText)

	pos2 := image.Point{4, 8}
	for i := range binDigits {
		binDigits[31-i] = ledgrid.NewPixel(pos2.Add(image.Point{i, 0}), color.Red)
		c.Add(binDigits[31-i])
	}

	timeLine1 := ledgrid.NewTimeline(time.Second)
	timeLine1.RepeatCount = ledgrid.AnimationRepeatForever
	timeLine1.Add(0, ledgrid.NewTask(func() {
		txt := time.Now().Format("15:04:05")
		clockText.SetText(txt)
		secSinceEpoc := time.Now().Unix()
		for i := range 32 {
			if secSinceEpoc&(1<<i) != 0 {
				binDigits[i].Show()
			} else {
				binDigits[i].Hide()
			}
		}
	}))

	digitColor := color.Blue
	colorFade = ledgrid.NewColorAnim(clockText, digitColor, 2*time.Second)
	colorFade.Val2 = ledgrid.RandColor(true)

	seq2 := ledgrid.NewSequence(colorFade, ledgrid.NewDelay(5*time.Second))
	seq2.RepeatCount = ledgrid.AnimationRepeatForever

	timeLine1.Start()
	seq2.Start()
}

func MovingText(ctx context.Context, c *ledgrid.Canvas) {
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

	aPosSeq := ledgrid.NewSequence(
        ledgrid.NewDelay(4*time.Second),
        aPos4,
        aPos5,
        aPos6,
        aPos7,
        aPos8,
    )
	// aPosSeq.SetDuration(aPosSeq.Duration() + 4*time.Second)
	aPosSeq.RepeatCount = ledgrid.AnimationRepeatForever

	aAngle1.Start()
	aAngle2.Start()
	aPosSeq.Start()
}
