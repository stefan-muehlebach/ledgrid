package main

import (
	"image"
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
						c.Add(0, pix)
						aPos := ledgrid.NewPositionAnim(pix, dest, time.Second+rand.N(time.Second))
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
						c.Add(0, pix)
						aPos := ledgrid.NewIntegerPosAnimation(&pix.Pos, dest.Int(), time.Second+rand.N(time.Second))
						aPos.AutoReverse = true
						grp.Add(aPos)
					}
				}
				aSeq.Add(grp)
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

					// pt := image.Point{x, y}
					// pix := ledgrid.NewPixel(pt, col)

					pos := geom.NewPoint(float64(x), float64(y))
					// pos := geom.NewPointIMG(pt)
					pix := ledgrid.NewDot(pos, col)

					c.Add(0, pix)

					dur := time.Second + time.Duration(x)*time.Millisecond
					aAlpha := ledgrid.NewFadeAnim(pix, 196, dur)
					aAlpha.AutoReverse = true
					aAlpha.RepeatCount = ledgrid.AnimationRepeatForever
					aAlpha.Start()

					aColor := ledgrid.NewColorAnim(pix, col, 1*time.Second)
					aColor.Cont = true
					aGrpGrey.Add(aColor)

					aColor = ledgrid.NewColorAnim(pix, color.MediumPurple.Interpolate(color.Fuchsia, t), 5*time.Second)
					aColor.Cont = true
					aGrpPurple.Add(aColor)

					aColor = ledgrid.NewColorAnim(pix, color.Gold.Interpolate(color.Khaki, t), 5*time.Second)
					aColor.Cont = true
					aGrpYellow.Add(aColor)

					aColor = ledgrid.NewColorAnim(pix, color.GreenYellow.Interpolate(color.LightSeaGreen, t), 5*time.Second)
					aColor.Cont = true
					aGrpGreen.Add(aColor)
				}
			}

			txt1 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.GreenYellow.Alpha(0.0), "LORENZ")
			txt1.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt1 := ledgrid.NewFadeAnim(txt1, ledgrid.FadeIn, 2*time.Second)
			aTxt1.AutoReverse = true
			txt2 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.DarkViolet.Alpha(0.0), "SIMON")
			txt2.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt2 := ledgrid.NewFadeAnim(txt2, ledgrid.FadeIn, 2*time.Second)
			aTxt2.AutoReverse = true
			txt3 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.OrangeRed.Alpha(0.0), "REBEKKA")
			txt3.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt3 := ledgrid.NewFadeAnim(txt3, ledgrid.FadeIn, 2*time.Second)
			aTxt3.AutoReverse = true
			c.Add(0, txt1, txt2, txt3)

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

	ShowThePaletteShader = NewLedGridProgram("Show a palette shader!",
		func(c *ledgrid.Canvas) {
			var xMin, yMax float64
			var txt *ledgrid.FixedText
			var palName string = "Hipster"
			var ptStart, pt, ptEnd fixed.Point26_6

			pt = fixed.P(1, height-1)
			ptStart = pt.Add(fixed.P(width, 0))
			ptEnd = pt.Sub(fixed.P(width, 0))

			txt = ledgrid.NewFixedText(pt, color.Gold, palName)

			pal := ledgrid.PaletteMap[palName]
			fader := ledgrid.NewPaletteFader(pal)
			aPal := ledgrid.NewPaletteFadeAnimation(fader, pal, 2*time.Second)
			aPal.ValFunc = ledgrid.SeqPalette()

			txtLeave := ledgrid.NewFixedPosAnim(txt, ptEnd, time.Second)
			txtLeave.Curve = ledgrid.AnimationEaseIn
			txtLeave.Cont = true
			txtEnter := ledgrid.NewFixedPosAnim(txt, pt, time.Second)
			txtEnter.Curve = ledgrid.AnimationEaseOut
			txtEnter.Cont = true
			txtNewText := ledgrid.NewTask(func() {
				txt.SetText(fader.Name())
				txt.Pos = ptStart
			})
			txtChange := ledgrid.NewSequence(txtLeave, txtNewText, txtEnter)

			aPalTl := ledgrid.NewTimeline(10 * time.Second)
			aPalTl.Add(7*time.Second, aPal, txtChange)
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
					c.Add(0, pix)
					anim := ledgrid.NewShaderAnim(pix, fader, x, y, PlasmaShaderFunc)
					aGrp.Add(anim)
					x += dPix
				}
				y -= dPix
			}
			c.Add(0, txt)

			aGrp.Start()
			aPalTl.Start()
		})

	ShowTheColorShader = NewLedGridProgram("Show a color shader!",
		func(c *ledgrid.Canvas) {
			var xMin, xMax, yMin, yMax, dx, dy float64

			randomList = make([]float64, width*height)
			for i := range width * height {
				randomList[i] = rand.Float64()
			}

            xMin =  0.0
            xMax =  float64(width)/10.0
            yMin =  0.0
            yMax =  float64(height)/10.0
            dx = (xMax - xMin) / float64(width)
            dy = (yMax - yMin) / float64(height)

			aGrp := ledgrid.NewGroup()
			nPix := c.Rect.Dx() * c.Rect.Dy()
			idx := 0
			y := yMax
			for row := range c.Rect.Dy() {
				x := xMin
				for col := range c.Rect.Dx() {
					pix := ledgrid.NewPixel(image.Point{col, row}, color.Black)
					c.Add(0, pix)
					anim := ledgrid.NewColorShaderAnim(pix, x, 0.0, y, idx, nPix, LavaLampShader)
					idx += 1
					aGrp.Add(anim)
					x += dx
				}
				y -= dy
			}
			aGrp.Start()
		})

	NyanCatShader = func(t, x, y, z float64, idx, nPix int) color.LedColor {
		y += myCos(x+0.2*z, 0, 1, 0, 0.6)
		z += myCos(x, 0, 1, 0, 0.3)
		x += myCos(y+z, 0, 1.5, 0, 0.2)

		x, y, z = y, z, x

		if idx % 7 == 0 {
			x += float64((idx*123)%5) / float64(nPix) * 32.12
			y += float64((idx*137)%5) / float64(nPix) * 22.23
			z += float64((idx*147)%7) / float64(nPix) * 44.34
		}

		r := myCos(x, t/4.0, 2.0, 0.0, 1.0)
		g := myCos(y, t/4.0, 2.0, 0.0, 1.0)
		b := myCos(z, t/4.0, 2.0, 0.0, 1.0)
		r, g, b = contrast(r, g, b, 0.5, 1.5)

		fade := math.Pow(myCos(t-float64(idx)/float64(nPix), 0, 7, 0, 1), 20.0)
		r *= fade
		g *= fade
		b *= fade

		twinkleSpeed := 0.07
		twinkleDensity := 0.1
		_, twinkle := math.Modf(randomList[idx]*7 + t*twinkleSpeed)
		twinkle = abs(twinkle*2 - 1)
		twinkle = remap(twinkle, 0, 1, -1/twinkleDensity, 1.1)
		twinkle = clamp(twinkle, -0.5, 1.1)
		twinkle = math.Pow(twinkle, 5.0)
		twinkle *= fade
		twinkle = clamp(twinkle, -0.3, 1)
		r += twinkle
		g += twinkle
		b += twinkle

        r = clamp(r * 256.0, 0.0, 255.0)
        g = clamp(g * 256.0, 0.0, 255.0)
        b = clamp(b * 256.0, 0.0, 255.0)

		return color.LedColor{uint8(r), uint8(g), uint8(b), 0xff}
	}

	LavaLampShader = func(t, x, y, z float64, idx, nPix int) color.LedColor {
		y += myCos(x+0.2*z, 0, 1, 0, 0.6)
		z += myCos(x, 0, 1, 0, 0.3)
		x += myCos(y+z, 0, 1.5, 0, 0.2)

		x, y, z = y, z, x

        r := myCos(x, t/4.0, 2, 0, 1)
        g := myCos(y, t/4.0, 2, 0, 1)
        b := myCos(z, t/4.0, 2, 0, 1)
		r, g, b = contrast(r, g, b, 0.5, 1.5)

        r2 := myCos(x, t/10.0 + 12.345, 3, 0, 1)
        g2 := myCos(y, t/10.0 + 24.536, 3, 0, 1)
        b2 := myCos(z, t/10.0 + 34.675, 3, 0, 1)
        clampDown := (r2 + g2 + b2)/2.0
        clampDown = remap(clampDown, 0.8, 0.9, 0, 1)
        clampDown = clamp(clampDown, 0, 1)
        r *= clampDown
        g *= clampDown
        b *= clampDown

        g = g * 0.6 + ((r + b) / 2.0) * 0.4

        r = clamp(r * 256.0, 0.0, 255.0)
        g = clamp(g * 256.0, 0.0, 255.0)
        b = clamp(b * 256.0, 0.0, 255.0)

		return color.LedColor{uint8(r), uint8(g), uint8(b), 0xff}
	}

	remap = func(x, minIn, maxIn, minOut, maxOut float64) float64 {
		t := (x - minIn) / (maxIn - minIn)
		return minOut + t*(maxOut-minOut)
	}

	clamp = func(x, lb, ub float64) float64 {
		return max(lb, min(ub, x))
	}

	myCos = func(x, off, period, lb, ub float64) float64 {
		val := math.Cos((x/period-off)*2.0*math.Pi)/2.0 + 0.5
		return val*(ub-lb) + lb
	}

	contrast = func(r0, g0, b0, center, scale float64) (r1, g1, b1 float64) {
		r1 = (r0-center)*scale + center
		g1 = (g0-center)*scale + center
		b1 = (b0-center)*scale + center
		return
	}

	randomList []float64

	// ---------------------------------------------------------------------------

	f1 = func(t, x, y, p1 float64) float64 {
		return math.Sin(x*p1 + t)
	}

	f2 = func(t, x, y, p1, p2, p3 float64) float64 {
		return math.Sin(p1*(x*math.Sin(t/p2)+y*math.Cos(t/p3)) + t)
	}

	f3 = func(t, x, y, p1, p2 float64) float64 {
		cx := 0.125*x + 0.5*math.Sin(t/p1)
		cy := 0.125*y + 0.5*math.Cos(t/p2)
		return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
	}

	PlasmaShaderFunc = func(t, x, y float64) float64 {
		v1 := f1(t, x, y, 10)       // old param: 1.2
		v2 := f2(t, x, y, 10, 2, 3) // old param: 1.6, 3.0, 1.5
		v3 := f3(t, x, y, 5, 3)     // old param: 5.0, 5.0
		v := (v1+v2+v3)/6.0 + 0.5
		return v
	}

	ColorFields = NewLedGridProgram("All the named colors",
		func(c *ledgrid.Canvas) {
			// var colGrpIdx color.ColorGroup = 0
			// var colIdx int = 0
			var colName string
			var nameList []string
			var nameIdx int = 0

			nameList = make([]string, len(color.Names))
			copy(nameList, color.Names)
			rand.Shuffle(len(nameList), func(i, j int) {
				nameList[i], nameList[j] = nameList[j], nameList[i]
			})
			colName = nameList[0]

			rectPos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			rectSize := geom.Point{float64(width), float64(height)}

			rect := ledgrid.NewRectangle(rectPos, rectSize, color.Black)
			rect.StrokeWidth = 0.0
			rect.FillColor = color.Black

			txtPos1 := fixed.P(width+1, height-1)
			txtPos2 := fixed.P(1, height-1)
			txtPos3 := fixed.P(-2*width, height-1)
			txt := ledgrid.NewFixedText(txtPos1, color.Black, "")
			c.Add(0, rect, txt)

			posAnim1 := ledgrid.NewFixedPosAnim(txt, txtPos2, 1*time.Second)
			posAnim1.Curve = ledgrid.AnimationEaseOut
			posAnim2 := ledgrid.NewFixedPosAnim(txt, txtPos3, 1*time.Second)
			posAnim2.Curve = ledgrid.AnimationEaseIn
			posAnim2.Cont = true

			fadeIn := ledgrid.NewFillColorAnim(rect, color.Map[colName], 1*time.Second)
			fadeIn.Curve = ledgrid.AnimationEaseOut
			fadeIn.Cont = true
			fadeOut := ledgrid.NewFillColorAnim(rect, color.Black, 1*time.Second)
			fadeOut.Curve = ledgrid.AnimationEaseIn
			fadeOut.Cont = true
			txtTask := ledgrid.NewTask(func() {
				var txtColor color.LedColor

				col := color.Map[colName]
				h, s, l := col.HSL()
				switch {
				case s == 0:
					txtColor = color.Gray
				case h >= 60:
					txtColor = color.Red
				case h >= -60:
					txtColor = color.Blue
				default:
					txtColor = color.Green
				}

				if l > 0.4 {
					txtColor = color.Gray.Dark(0.7)
				} else {
					txtColor = color.Gray.Bright(0.5)
				}
				txt.SetText(colName)
				txt.Color = txtColor
			})
			colTask := ledgrid.NewTask(func() {
				oldColor := color.Map[colName]
				nameIdx = (nameIdx + 1) % len(nameList)

				// if colIdx >= len(color.Groups[colGrpIdx]) {
				// 	colGrpIdx = (colGrpIdx + 1) % color.NumColorGroups
				// 	colIdx = 0
				// }
				// colName = color.Groups[colGrpIdx][colIdx]

				colName = nameList[nameIdx]
				newColor := color.Map[colName]
				fadeOut.Val2 = ledgrid.Const(oldColor.Interpolate(newColor, 0.5))
				fadeIn.Val2 = ledgrid.Const(newColor)
			})

			timeLine := ledgrid.NewTimeline(4 * time.Second)
			timeLine.Add(0*time.Second, txtTask, posAnim1, fadeIn)
			timeLine.Add(2*time.Second, colTask)
			timeLine.Add(3*time.Second, fadeOut, posAnim2)
			timeLine.RepeatCount = ledgrid.AnimationRepeatForever

			timeLine.Start()
		})

	FirePlace = NewLedGridProgram("Fireplace",
		func(c *ledgrid.Canvas) {
			fire := ledgrid.NewFire(image.Point{}, image.Point{width, height})
			c.Add(0, fire)
			fire.Start()
		})

	SpecialCamera = NewLedGridProgram("Camera in differential mode",
		func(c *ledgrid.Canvas) {
			pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			size := geom.Point{float64(width), float64(height)}

			cam := NewHistCamera(pos, size, 100, color.SkyBlue)
			c.Add(0, cam)
			cam.Start()
		})
)

// ---------------------------------------------------------------------------

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
