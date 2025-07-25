package main

import (
	"context"
	"fmt"
	"image"
	"math"
	"math/rand/v2"
	"time"

	"golang.org/x/image/math/fixed"

	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colors"
)

func init() {
	programList.Add("Rotating, Floating Words", "Text", MovingText)
	programList.Add("All Named Colors", "Text", NamedColors)
	programList.Add("Clock Animation", "Text", ClockAnimation)
	programList.Add("Countdown (Rectangle)", "Text", RectangleCountdown)
	programList.Add("Countdown (Glowing)", "Text", GlowingCountdown)
	programList.Add("Punktestand", "Text", BoredScores)
}

func f2f(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}
func p2p(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{f2f(x), f2f(y)}
}

func ClockAnimation(ctx context.Context, c *ledgrid.Canvas) {
	var clockText *ledgrid.FixedText
	var binDigits [32]*ledgrid.Pixel

	pos1 := p2p(5.0, 6.0)
	clockText = ledgrid.NewFixedText(pos1, "00:00:00", colors.Blue.Dark(0.3))
	c.Add(clockText)

	pos2 := image.Point{4, 8}
	for i := range binDigits {
		binDigits[31-i] = ledgrid.NewPixel(pos2.Add(image.Point{i, 0}), colors.Red.Dark(0.3))
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

	timeLine1.Start()
}

func RectangleCountdown(ctx context.Context, c *ledgrid.Canvas) {
	dur1 := 10 * time.Second
	dur2 := 7 * time.Second
	pal1 := ledgrid.PaletteMap["BackPinkBlue"]
	pal2 := ledgrid.PaletteMap["BackYellowBlue"]
	rectPos := geom.Point{float64(width) / 2, float64(height)/2 - 0.5}
	rectSize := geom.Point{float64(width) - 1, float64(height) - 2}

	textPos := p2p(3.0, 7.0)
	pit, _ := time.Parse("02.01.2006", "29.06.2025")

	rect := ledgrid.NewRectangle(rectPos, rectSize, colors.Indigo)
	clockText := ledgrid.NewFixedText(textPos, "0", colors.Indigo)
	c.Add(rect, clockText)

	pal1Anim := ledgrid.NewPaletteAnim(rect, pal1, dur1)
	pal1Anim.AutoReverse = true
	pal1Anim.Curve = ledgrid.AnimationLinear
	pal1Seq := ledgrid.NewSequence(pal1Anim)
	pal1Seq.RepeatCount = ledgrid.AnimationRepeatForever

	pal2Anim := ledgrid.NewPaletteAnim(clockText, pal2, dur2)
	pal2Anim.AutoReverse = true
	pal2Anim.Curve = ledgrid.AnimationLinear
	pal2Seq := ledgrid.NewSequence(pal2Anim)
	pal2Seq.RepeatCount = ledgrid.AnimationRepeatForever

	timeLine1 := ledgrid.NewTimeline(10 * time.Millisecond)
	timeLine1.RepeatCount = ledgrid.AnimationRepeatForever
	timeLine1.Add(0, ledgrid.NewTask(func() {
		duration := time.Until(pit)
		txt := fmt.Sprintf("%9.1f", duration.Seconds())
		clockText.SetText(txt)
	}))

	pal1Seq.Start()
	pal2Seq.Start()
	timeLine1.Start()
}

func GlowingCountdown(ctx context.Context, c *ledgrid.Canvas) {
	aGrpLedColor := ledgrid.NewGroup()
	fadeOutPixels := ledgrid.NewGroup()
	pulseDur := 2 * time.Second
    fadeOutDur := 4 * time.Second
	textPos := p2p(3.0, 8.0)
	clockText := ledgrid.NewFixedText(textPos, "0", colors.FireBrick)

	backPal := ledgrid.PaletteMap["BackYellowBlue"]
	palFader := ledgrid.NewPaletteFader(backPal)
	aPalFade := ledgrid.NewPaletteFadeAnim(palFader, ledgrid.PaletteMap["BackPinkBlue"], 5*time.Second)
	aTxtFade := ledgrid.NewColorAnim(clockText, colors.Yellow, 5*time.Second)
	fadeOutText := ledgrid.NewFadeAnim(clockText, ledgrid.FadeOut, 5*time.Second)

	pit, _ := time.Parse("02.01.2006 15:04:05", "28.06.2025 10:30:00")
	finalDurA := 6 * time.Second
	finalDurB := 6 * time.Second
	finalDur := finalDurA + finalDurB
	// pit := time.Now().Add(finalDur + 5*time.Second)

	for y := range c.Rect.Dy() {
		for x := range c.Rect.Dx() {
			pt := image.Point{x, y}
			pix := ledgrid.NewPixel(pt, colorList[0][0])

			c.Add(pix)

			aPal := ledgrid.NewPaletteAnim(pix, palFader, pulseDur)
			aPal.Pos = rand.Float64()

			aFade := ledgrid.NewFadeAnim(pix, ledgrid.FadeOut, fadeOutDur)
			fadeOutPixels.Add(aFade)

			aColorSeq := ledgrid.NewSequence(aPal)
			aColorSeq.RepeatCount = ledgrid.AnimationRepeatForever
			aGrpLedColor.Add(aColorSeq)
		}
	}

	c.Add(clockText /*, alertText*/)

	updateCounter := ledgrid.NewTimeline(10 * time.Millisecond)
	updateCounter.RepeatCount = ledgrid.AnimationRepeatForever
	updateCounter.Add(0, ledgrid.NewTask(func() {
		duration := time.Until(pit)
		if duration < 0 {
			duration = 0
		}
		txt := fmt.Sprintf("%9.1f", duration.Seconds())
		clockText.SetText(txt)
	}))

	finalizeCounter := ledgrid.NewTimeline(finalDur + 10*time.Second)
	finalizeCounter.Add(0, ledgrid.NewTask(func() {
		aPalFade.Start()
		aTxtFade.Start()
	}))
	finalizeCounter.Add(finalDur-finalDurA, ledgrid.NewTask(func() {
		aPalFade.Val2 = ledgrid.PaletteMap["BackYellowOrange"]
		aPalFade.Start()
		aTxtFade.Val2 = ledgrid.Const(colors.WhiteSmoke)
		aTxtFade.Start()
	}))
	finalizeCounter.Add(finalDur, ledgrid.NewTask(func() {
		updateCounter.Suspend()
		aGrpLedColor.Suspend()
		fadeOutPixels.Start()
	}))
	finalizeCounter.Add(finalDur+4*time.Second, ledgrid.NewTask(func() {
		fadeOutText.Start()
	}))

	checkFinale := ledgrid.NewSequence()
	checkFinale.Add(ledgrid.NewTask(func() {
		if time.Until(pit) <= finalDur {
			finalizeCounter.Start()
			checkFinale.Suspend()
		}
	}))
	checkFinale.SetDuration(1 * time.Second)
	checkFinale.RepeatCount = ledgrid.AnimationRepeatForever

	updateCounter.Start()
	checkFinale.Start()
	aGrpLedColor.Start()
	// fadeSeq.Start()

}

func MovingText(ctx context.Context, c *ledgrid.Canvas) {
	t1 := ledgrid.NewText(geom.Point{0, float64(height) / 2.0}, "Stefan", colors.LightSeaGreen)
	t1.SetAlign(ledgrid.AlignLeft)
	t2 := ledgrid.NewText(geom.Point{float64(width), float64(height) / 2.0}, "Beni", colors.YellowGreen)
	t2.SetAlign(ledgrid.AlignRight)

	t4 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "werden", colors.Gold)
	t5 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "immer", colors.Gold)
	t6 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "im", colors.Gold)
	t7 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "Lochbach", colors.Gold)
	t8 := ledgrid.NewText(geom.Point{float64(width) / 2.0, float64(height) * 1.5}, "wohnen", colors.Gold)

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
	aPosSeq.RepeatCount = ledgrid.AnimationRepeatForever

	aAngle1.Start()
	aAngle2.Start()
	aPosSeq.Start()
}

func NamedColors(ctx context.Context, c *ledgrid.Canvas) {
	var colName string
	var nameList []string
	var nameIdx int = 0

	nameList = make([]string, len(colors.Names))
	copy(nameList, colors.Names)
	colName = nameList[0]

	rectPos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	rectSize := geom.Point{float64(width), float64(height)}

	rect := ledgrid.NewRectangle(rectPos, rectSize, colors.Black)
	rect.StrokeWidth = 0.0
	rect.FillColor = colors.Black

	txtPos1 := geom.Point{float64(width + 1), float64(height - 1)}
	txtPos2 := geom.Point{1.5, float64(height - 1)}
	txtPos3 := geom.Point{1.5, -3}
	txt := ledgrid.NewText(txtPos1, "", colors.Black)
	txt.SetAlign(ledgrid.AlignLeft | ledgrid.AlignBottom)
	txt.SetFont(fonts.GoBold, 10.0)

	c.Add(rect, txt)

	posAnim1 := ledgrid.NewPositionAnim(txt, txtPos2, 3*time.Second/2)
	posAnim1.Cont = false
	posAnim1.Curve = ledgrid.AnimationEaseOut
	posAnim2 := ledgrid.NewPositionAnim(txt, txtPos3, time.Second/2)
	posAnim2.Curve = ledgrid.AnimationEaseIn

	fadeIn := ledgrid.NewFillColorAnim(rect, colors.Map[colName], 1*time.Second)
	fadeIn.Curve = ledgrid.AnimationEaseOut
	fadeOut := ledgrid.NewFillColorAnim(rect, colors.Black, 1*time.Second)
	fadeOut.Curve = ledgrid.AnimationEaseIn
	txtTask := ledgrid.NewTask(func() {
		var txtColor colors.LedColor

		col := colors.Map[colName]
		h, s, l := col.HSL()
		switch {
		case s == 0:
			txtColor = colors.Gray
		case h >= 60:
			txtColor = colors.Red
		case h >= -60:
			txtColor = colors.Blue
		default:
			txtColor = colors.Green
		}

		if l > 0.4 {
			txtColor = colors.Gray.Dark(0.7)
		} else {
			txtColor = colors.Gray.Bright(0.5)
		}
		txt.Text = colName
		txt.Color = txtColor
	})
	colTask := ledgrid.NewTask(func() {
		oldColor := colors.Map[colName]
		nameIdx = (nameIdx + 1) % len(nameList)
		colName = nameList[nameIdx]
		newColor := colors.Map[colName]
		fadeOut.Val2 = ledgrid.Const(oldColor.Interpolate(newColor, 0.5))
		fadeIn.Val2 = ledgrid.Const(newColor)
	})

	timeLine := ledgrid.NewTimeline(3 * time.Second)
	timeLine.Add(0*time.Second, txtTask, posAnim1, fadeIn)
	timeLine.Add(1500*time.Millisecond, colTask)
	timeLine.Add(2*time.Second, fadeOut, posAnim2)
	timeLine.RepeatCount = ledgrid.AnimationRepeatForever

	timeLine.Start()
}

func BoredScores(ctx context.Context, c *ledgrid.Canvas) {
    var a, b int
	var scoreA, scoreB *ledgrid.FixedText

	posA := p2p(1.0, 9.0)
    posB := p2p(31.0, 9.0)
	scoreA = ledgrid.NewFixedText(posA, "0", colors.MediumSeaGreen)
    scoreA.SetFont(ledgrid.Lightdot8x8)
	scoreB = ledgrid.NewFixedText(posB, "0", colors.MediumSeaGreen)
    scoreB.SetFont(ledgrid.Lightdot6x8)
	c.Add(scoreA, scoreB)

    a, b = 0, 0
    	timeLine1 := ledgrid.NewTimeline(time.Second)
	timeLine1.RepeatCount = ledgrid.AnimationRepeatForever
	timeLine1.Add(0, ledgrid.NewTask(func() {
		txtA := fmt.Sprintf("%d", a)
        txtB := fmt.Sprintf("%d", b)
		scoreA.SetText(txtA)
		scoreB.SetText(txtB)
        a, b = (a+1)%10, (b+1)%10
	}))

	timeLine1.Start()


}
