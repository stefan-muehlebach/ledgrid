package main

import (
	"image"
	"math/rand/v2"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
	"golang.org/x/image/math/fixed"
)

var (
	FarewellGery = NewLedGridProgram("Farewell Gery!",
		func(c *ledgrid.Canvas) {
			aGrpFadeIn := ledgrid.NewGroup()
			aGrpPurple := ledgrid.NewGroup()
			aGrpYellow := ledgrid.NewGroup()
			aGrpGreen := ledgrid.NewGroup()
			aGrpGrey := ledgrid.NewGroup()
			aGrpRed := ledgrid.NewGroup()
			aGrpBlack := ledgrid.NewGroup()
			aSeqColor := ledgrid.NewSequence(aGrpRed, aGrpGreen)

			for y := range c.Rect.Dy() {
				for x := range c.Rect.Dx() {
					pt := image.Point{x, y}
					pos := geom.NewPointIMG(pt)
					t := rand.Float64()
					col := color.Black
					pix := ledgrid.NewDot(pos, col)
					c.Add(pix)

					dur := time.Second + time.Duration(10*x+20*y)*time.Millisecond
					aAlpha := ledgrid.NewFadeAnim(pix, 196, dur)
					aAlpha.AutoReverse = true
					aAlpha.RepeatCount = ledgrid.AnimationRepeatForever
					aAlpha.Start()

					aColor := ledgrid.NewColorAnim(pix, (color.DimGray.Dark(0.5)).Interpolate((color.DarkGrey.Dark(0.5)), t), 9*time.Second)
					aColor.Cont = true
					aGrpFadeIn.Add(aColor)

					aColor = ledgrid.NewColorAnim(pix, (color.DimGray.Dark(0.5)).Interpolate((color.DarkGrey.Dark(0.5)), t), 1*time.Second)
					aColor.Cont = true
					aGrpGrey.Add(aColor)

					aColor = ledgrid.NewColorAnim(pix, color.MediumPurple.Interpolate(color.Fuchsia, t), 4*time.Second)
					aColor.Cont = true
					aGrpPurple.Add(aColor)

					aColor = ledgrid.NewColorAnim(pix, color.Gold.Interpolate(color.LemonChiffon, t), 4*time.Second)
					aColor.Cont = true
					aGrpYellow.Add(aColor)

					aColor = ledgrid.NewColorAnim(pix, color.Crimson.Interpolate(color.Orange, t), 4*time.Second)
					aColor.Cont = true
					aGrpRed.Add(aColor)

					aColor = ledgrid.NewColorAnim(pix, color.LightSeaGreen.Interpolate(color.GreenYellow, t), 500*time.Millisecond)
					aColor.Cont = true
					aGrpGreen.Add(aColor)

					aColor = ledgrid.NewColorAnim(pix, color.Black, 2*time.Second)
					aColor.Cont = true
					aGrpBlack.Add(aColor)
				}
			}

			txt1 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.GreenYellow.Alpha(0.0), "LIEBER")
			txt1.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt1 := ledgrid.NewFadeAnim(txt1, ledgrid.FadeIn, 1*time.Second)
			aTxt1.AutoReverse = true
			txt2 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.DarkViolet.Alpha(0.0), "GERY")
			txt2.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt2 := ledgrid.NewFadeAnim(txt2, ledgrid.FadeIn, 2*time.Second)
			aTxt2.AutoReverse = true
			txt3 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.RoyalBlue.Alpha(0.0), "FAREWELL")
			txt3.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt3 := ledgrid.NewFadeAnim(txt3, ledgrid.FadeIn, 5*time.Second)
			aTxt3.AutoReverse = true
			c.Add(txt1, txt2, txt3)

			aTimel := ledgrid.NewTimeline(40 * time.Second)
			aTimel.Add(0, aGrpFadeIn)
			aTimel.Add(10*time.Second, aGrpPurple)
			aTimel.Add(13*time.Second, aTxt1)
			aTimel.Add(14*time.Second, aGrpGrey)

			aTimel.Add(20*time.Second, aGrpYellow)
			aTimel.Add(23*time.Second, aTxt2)
			aTimel.Add(24*time.Second, aGrpGrey)

			aTimel.Add(30*time.Second, aSeqColor)
			aTimel.Add(33*time.Second, aTxt3)
			aTimel.Add(35*time.Second, aGrpBlack)
			aTimel.RepeatCount = ledgrid.AnimationRepeatForever

			aTimel.Start()
		})

	FirePlace = NewLedGridProgram("Fireplace",
		func(c *ledgrid.Canvas) {
			fire := ledgrid.NewFire(image.Point{}, image.Point{width, height})
			c.Add(fire)
			fire.Start()
		})
)
