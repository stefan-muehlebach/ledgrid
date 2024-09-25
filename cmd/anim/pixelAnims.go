package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
	"golang.org/x/image/math/fixed"
)

var (
	MovingPixels = NewLedGridProgram("Moving pixels",
		func(c *ledgrid.Canvas) {
			mp := geom.Point{float64(width)/2 - 0.5, float64(height)/2 - 0.5}
			// aKillGrp := ledgrid.NewGroup()
			aSeq := ledgrid.NewSequence()
			for i := range 8 {
				grp := ledgrid.NewGroup()

				xMin, xMax := float64(i), float64(width-i)
				yMin, yMax := float64(i), float64(height-i)
				col := color.RandGroupColor(color.Purples).Dark(float64(5-i) * 0.1)
				posList := []geom.Point{
					geom.Point{0.0, yMin},
					geom.Point{0.0, yMax - 1},
				}
				for x := xMin; x < xMax; x++ {
					for j := range 2 {
						pos := posList[j]
						pos.X = float64(x)
						dest := pos.Sub(mp).Normalize().Mul(20.0).Add(pos)
						pix := ledgrid.NewDot(pos, col)
						// aKillGrp.Add(ledgrid.NewTask(func() {
						// 	pix.Kill()
						// }))
						c.Add(pix)
						aPos := ledgrid.NewPositionAnimation(&pix.Pos, dest, time.Second+rand.N(time.Second))
						aPos.AutoReverse = true
						grp.Add(aPos)
					}
				}
				posList = []geom.Point{
					geom.Point{xMin, 0.0},
					geom.Point{xMax - 1, 0.0},
				}
				for y := yMin + 1; y < yMax-1; y++ {
					for j := range 2 {
						pos := posList[j]
						pos.Y = float64(y)
						dest := pos.Sub(mp).Normalize().Mul(20.0).Add(pos)
						pix := ledgrid.NewPixel(pos.Int(), col)
						// aKillGrp.Add(ledgrid.NewTask(func() {
						// 	pix.Kill()
						// }))
						c.Add(pix)
						aPos := ledgrid.NewIntegerPosAnimation(&pix.Pos, dest.Int(), time.Second+rand.N(time.Second))
						aPos.AutoReverse = true
						grp.Add(aPos)
					}
				}
				aSeq.Put(grp)
			}
			aSeq.RepeatCount = ledgrid.AnimationRepeatForever
			aSeq.Start()

			// time.Sleep(5 * time.Second)
			// aKillGrp.Start()
		})

	GlowingPixels = NewLedGridProgram("Glowing pixels",
		func(c *ledgrid.Canvas) {
			aGrpPurple := ledgrid.NewGroup()
			aGrpYellow := ledgrid.NewGroup()
			aGrpGreen := ledgrid.NewGroup()
			aGrpGrey := ledgrid.NewGroup()

			for y := range c.Rect.Dy() {
				for x := range c.Rect.Dx() {
					t := rand.Float64()
					col := (color.DimGray.Dark(0.3)).Interpolate((color.DarkGrey.Dark(0.3)), t)

					pt := image.Point{x, y}
					// pix := ledgrid.NewPixel(pt, col)

					pos := geom.NewPointIMG(pt)
					pix := ledgrid.NewDot(pos, col)

					c.Add(pix)

					dur := time.Second + time.Duration(x)*time.Millisecond
					aAlpha := ledgrid.NewFadeAnimation(&pix.Color.A, 196, dur)
					aAlpha.AutoReverse = true
					aAlpha.RepeatCount = ledgrid.AnimationRepeatForever
					aAlpha.Start()

					aColor := ledgrid.NewColorAnimation(&pix.Color, col, 1*time.Second)
					aColor.Cont = true
					aGrpGrey.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.MediumPurple.Interpolate(color.Fuchsia, t), 5*time.Second)
					aColor.Cont = true
					aGrpPurple.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.Gold.Interpolate(color.Khaki, t), 5*time.Second)
					aColor.Cont = true
					aGrpYellow.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.GreenYellow.Interpolate(color.LightSeaGreen, t), 5*time.Second)
					aColor.Cont = true
					aGrpGreen.Add(aColor)
				}
			}

			txt1 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.GreenYellow.Alpha(0.0), "LORENZ")
			txt1.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt1 := ledgrid.NewFadeAnimation(&txt1.Color.A, ledgrid.FadeIn, 2*time.Second)
			aTxt1.AutoReverse = true
			txt2 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.DarkViolet.Alpha(0.0), "SIMON")
			txt2.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt2 := ledgrid.NewFadeAnimation(&txt2.Color.A, ledgrid.FadeIn, 2*time.Second)
			aTxt2.AutoReverse = true
			txt3 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.OrangeRed.Alpha(0.0), "REBEKKA")
			txt3.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt3 := ledgrid.NewFadeAnimation(&txt3.Color.A, ledgrid.FadeIn, 2*time.Second)
			aTxt3.AutoReverse = true
			c.Add(txt1, txt2, txt3)

			aTimel := ledgrid.NewTimeline(42 * time.Second)
			aTimel.Add(7*time.Second, aGrpPurple)
			aTimel.Add(12*time.Second, aTxt1)
			aTimel.Add(13*time.Second, aGrpGrey)

			aTimel.Add(22*time.Second, aGrpYellow)
			aTimel.Add(27*time.Second, aTxt2)
			aTimel.Add(28*time.Second, aGrpGrey)

			aTimel.Add(35*time.Second, aGrpGreen)
			aTimel.Add(40*time.Second, aTxt3)
			aTimel.Add(41*time.Second, aGrpGrey)
			aTimel.RepeatCount = ledgrid.AnimationRepeatForever

			aTimel.Start()
		})

	TestShaderFunc = func(x, y, t float64) float64 {
		t = t/4.0 + x
		_, f := math.Modf(math.Abs(t))
		return f
	}

	f1 = func(x, y, t, p1 float64) float64 {
		return math.Sin(x*p1 + t)
	}

	f2 = func(x, y, t, p1, p2, p3 float64) float64 {
		return math.Sin(p1*(x*math.Sin(t/p2)+y*math.Cos(t/p3)) + t)
	}

	f3 = func(x, y, t, p1, p2 float64) float64 {
		cx := 0.125*x + 0.5*math.Sin(t/p1)
		cy := 0.125*y + 0.5*math.Cos(t/p2)
		return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
	}

	PlasmaShaderFunc = func(x, y, t float64) float64 {
		v1 := f1(x, y, t, 1.2)
		v2 := f2(x, y, t, 1.6, 3.0, 1.5)
		v3 := f3(x, y, t, 5.0, 5.0)
		v := (v1+v2+v3)/6.0 + 0.5
		return v
	}

	ShowTheShader = NewLedGridProgram("Show the shader!",
		func(c *ledgrid.Canvas) {
			var xMin, yMax float64
			var txt *ledgrid.FixedText
			var palId int

			pal := ledgrid.PaletteMap["Hipster"]
			fader := ledgrid.NewPaletteFader(pal)
			aPal := ledgrid.NewPaletteFadeAnimation(fader, pal, 2*time.Second)
			aPal.ValFunc = func() ledgrid.ColorSource {
				name := ledgrid.PaletteNames[palId]
				palId = (palId + 1) % len(ledgrid.PaletteNames)
				log.Printf(">>> Switch palette, new name: '%s'", name)
				txt.SetText(name)
				return ledgrid.PaletteMap[name]
			}

			aPalTl := ledgrid.NewTimeline(10 * time.Second)
			aPalTl.Add(7*time.Second, aPal)
			aPalTl.RepeatCount = ledgrid.AnimationRepeatForever

			aGrp := ledgrid.NewGroup()
			dPix := 2.0 / float64(max(c.Rect.Dx(), c.Rect.Dy())-1)
			ratio := float64(c.Rect.Dx()) / float64(c.Rect.Dy())
			if ratio > 1.0 {
				xMin = -1.0
				yMax = ratio * 1.0
			} else if ratio < 1.0 {
				xMin = ratio * -1.0
				yMax = 1.0
			} else {
				xMin = -1.0
				yMax = 1.0
			}

			y := yMax
			for row := range c.Rect.Dy() {
				x := xMin
				for col := range c.Rect.Dx() {
					pix := ledgrid.NewPixel(image.Point{col, row}, color.Black)
					c.Add(pix)
					anim := ledgrid.NewShaderAnimation(&pix.Color, fader, x, y, PlasmaShaderFunc)
					aGrp.Add(anim)
					x += dPix
				}
				y -= dPix
			}
			txt = ledgrid.NewFixedText(fixed.P(width/2, height/2), color.YellowGreen, "Hipster")
			txt.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			c.Add(txt)
			aPalTl.Start()
			aGrp.Start()
		})

	ColorFields = NewLedGridProgram("Fields of named colors",
		func(c *ledgrid.Canvas) {
			var input int
			var colGrp color.ColorGroup

			cs := NewColorSampler(color.Purples)
			c.Add(cs)

			for {
				fmt.Printf("Enter a number in 0..%d (or 99 to quit): ", color.NumColorGroups-1)
				fmt.Scanf("%d\n", &input)
				if input == 99 {
					return
				}
				colGrp = color.ColorGroup(input)
				if colGrp >= color.NumColorGroups {
					continue
				}
				fmt.Printf("Selected color group: %v\n", colGrp)
				cs.colGrp = colGrp
			}
		})
)

type ColorSampler struct {
	ledgrid.CanvasObjectEmbed
	colGrp color.ColorGroup
}

func NewColorSampler(colGrp color.ColorGroup) *ColorSampler {
	c := &ColorSampler{}
	c.CanvasObjectEmbed.Extend(c)
	c.colGrp = colGrp
	return c
}

func (c *ColorSampler) Draw(canv *ledgrid.Canvas) {
	for i, colorName := range color.Groups[c.colGrp] {
		col := color.Map[colorName]
		for j := range 2 {
			x := 2*i + j
			if x >= width {
				return
			}
			for y := range height {
				canv.GC.SetPixel(x, y, col)
			}
		}
	}
}
