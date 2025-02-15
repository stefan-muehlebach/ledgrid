package main

import (
	"bufio"
	"context"
	"image"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
	"golang.org/x/image/math/fixed"
)

func init() {
	programList.AddTitle("Pixel Animations")
	// programList.Add("Moving pixels", MovingPixels)
	programList.Add("Glowing pixels with changing text", GlowingPixels)
	programList.Add("Fireplace", Fireplace)
	programList.Add("Shader using palettes", PaletteShader)
	programList.Add("Shader using colors", ColorShader)
}

var (
	colorList = [][]color.LedColor{
		{color.NewLedColorHex(0xb9b90a), color.NewLedColorHex(0x0a5853)}, // Yellow to LightBlue
		{color.NewLedColorHex(0xa60c5f), color.NewLedColorHex(0x000080)}, // DeepPink to DarkBlue
		{color.NewLedColorHex(0x959500), color.NewLedColorHex(0x710e00)}, // Yellow to OrangeRed
		{color.NewLedColorHex(0x007400), color.NewLedColorHex(0x9a2222)}, // DarkGreen to DarkRed
		{color.NewLedColorHex(0x814124), color.NewLedColorHex(0x4c8baa)}, // Salmon to LightBlue
	}
)

func MovingPixels(ctx context.Context, c *ledgrid.Canvas) {
	mp := geom.Point{float64(width)/2 - 0.5, float64(height)/2 - 0.5}
	aSeq := ledgrid.NewSequence()
	for i := range 8 {
		grp := ledgrid.NewGroup()

		xMin, xMax := float64(3+i), float64(width-3-i)
		yMin, yMax := float64(3+i), float64(height-3-i)
		col1 := color.RandGroupColor(color.Blues).Dark(float64(5-i) * 0.1)
		col2 := color.RandGroupColor(color.Reds).Dark(float64(5-i) * 0.1)
		posList := []geom.Point{
			geom.Point{0.0, yMin},
			geom.Point{0.0, yMax - 1},
		}
		for x := xMin; x < xMax; x++ {
			for j := range 2 {
				pos := posList[j]
				pos.X = float64(x)
				dest := pos.Sub(mp).Normalize().Mul(20.0).Add(pos)
				pix := ledgrid.NewDot(pos, col1)
				c.Add(pix)
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
				pix := ledgrid.NewPixel(pos.Int(), col2)
				c.Add(pix)
				aPos := ledgrid.NewIntegerPosAnim(pix, dest.Int(), time.Second+rand.N(time.Second))
				aPos.AutoReverse = true
				grp.Add(aPos)
			}
		}
		aSeq.Add(grp)
	}
	aSeq.RepeatCount = ledgrid.AnimationRepeatForever
	aSeq.Start()
}

func GlowingPixels(ctx context.Context, c *ledgrid.Canvas) {
	aGrpLedColor := ledgrid.NewGroup()
	dur := 3 * time.Second
	numReps := 1

	f, err := os.Open("Faust.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for y := range c.Rect.Dy() {
		for x := range c.Rect.Dx() {
			pt := image.Point{x, y}
			pix := ledgrid.NewPixel(pt, colorList[0][0])

			c.Add(pix)

			aColorCyc := ledgrid.NewColorAnim(pix, colorList[0][1], dur)
			aColorCyc.AutoReverse = true
			aColorCyc.RepeatCount = numReps
			aColorCyc.Curve = ledgrid.AnimationLinear
			aColorCyc.Pos = rand.Float64()

			aColorSeq := ledgrid.NewSequence(aColorCyc)

			for _, colPair := range colorList[1:] {

				aColorTrans := ledgrid.NewColorAnim(pix, colPair[0], dur)
				aColorTrans.Curve = ledgrid.AnimationLinear

				aColorCyc := ledgrid.NewColorAnim(pix, colPair[1], dur)
				aColorCyc.AutoReverse = true
				aColorCyc.RepeatCount = numReps
				aColorCyc.Curve = ledgrid.AnimationLinear

				aColorSeq.Add(aColorTrans, aColorCyc)
			}

			aColorTrans := ledgrid.NewColorAnim(pix, colorList[0][0], dur)
			aColorTrans.Curve = ledgrid.AnimationLinear

			aColorSeq.Add(aColorTrans)
			aColorSeq.RepeatCount = ledgrid.AnimationRepeatForever

			aGrpLedColor.Add(aColorSeq)
		}
	}

	txt := ledgrid.NewFixedText(fixed.P(width/2, height/2), "", color.White.Alpha(0))
	txt.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)

	txtFadeIn := ledgrid.NewFadeAnim(txt, ledgrid.FadeIn, 1000*time.Millisecond)
	txtColorOut := ledgrid.NewColorAnim(txt, color.Black, 1000*time.Millisecond)
	txtFadeOut := ledgrid.NewFadeAnim(txt, ledgrid.FadeOut, 3000*time.Millisecond)

	txtNextWord := ledgrid.NewTask(func() {
		if scanner.Scan() {
			txt.SetText(scanner.Text())
		}
		txt.Color = color.White.Alpha(0)
	})
	txtSeq := ledgrid.NewSequence(txtNextWord, txtFadeIn,
		ledgrid.NewDelay(time.Second), txtColorOut, txtFadeOut)
	txtSeq.RepeatCount = ledgrid.AnimationRepeatForever

	c.Add(txt)
	txtSeq.Start()
	aGrpLedColor.Start()
}


func Fireplace(ctx context.Context, c *ledgrid.Canvas) {
	fire := ledgrid.NewFire(image.Point{}, image.Point{width, height})
	c.Add(fire)
	fire.Start()
}

func PaletteShader(ctx context.Context, c *ledgrid.Canvas) {
	var xMin, yMax float64
	var txt *ledgrid.FixedText
	var palName string = "Hipster"
	var ptStart, pt, ptEnd fixed.Point26_6

	pt = fixed.P(1, height-1)
	ptStart = pt.Add(fixed.P(width, 0))
	ptEnd = pt.Sub(fixed.P(width, 0))

	txt = ledgrid.NewFixedText(pt, palName, color.Gold)

	pal := ledgrid.PaletteMap[palName]
	fader := ledgrid.NewPaletteFader(pal)
	aPal := ledgrid.NewPaletteFadeAnimation(fader, pal, 2*time.Second)
	aPal.ValFunc = ledgrid.SeqPalette()

	txtLeave := ledgrid.NewFixedPosAnim(txt, ptEnd, time.Second)
	txtLeave.Curve = ledgrid.AnimationEaseIn
	txtEnter := ledgrid.NewFixedPosAnim(txt, pt, time.Second)
	txtEnter.Curve = ledgrid.AnimationEaseOut
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
			c.Add(pix)
			anim := ledgrid.NewShaderAnim(pix, fader, x, y, PlasmaShaderFunc)
			aGrp.Add(anim)
			x += dPix
		}
		y -= dPix
	}
	c.Add(txt)

	aGrp.Start()
	aPalTl.Start()
}

// A color shader - the new thing people are talking about.
func ColorShader(ctx context.Context, c *ledgrid.Canvas) {
	var xMin, xMax, yMin, yMax, dx, dy float64

	randomList = make([]float64, width*height)
	for i := range width * height {
		randomList[i] = rand.Float64()
	}
	w, h := float64(width)/10.0, float64(height)/10.0
	xMin = -w / 2
	xMax = w / 2
	yMin = -h / 2
	yMax = h / 2
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
			c.Add(pix)
			anim := ledgrid.NewColorShaderAnim(pix, x, 0.0, y, idx, nPix, RandomShader)
			idx += 1
			aGrp.Add(anim)
			x += dx
		}
		y -= dy
	}
	aGrp.Start()
}

var (
	randomList []float64

	// Einer der Shader (oder dort: color functions) aus den Testprogrammen
	// von OpenPixelController.

	NyanCatShader = func(t, x, y, z float64, idx, nPix int) color.LedColor {
		y += myCos(x+0.2*z, 0, 1, 0, 0.6)
		z += myCos(x, 0, 1, 0, 0.3)
		x += myCos(y+z, 0, 1.5, 0, 0.2)

		x, y, z = y, z, x

		if idx%7 == 0 {
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

		r = clamp(r*256.0, 0.0, 255.0)
		g = clamp(g*256.0, 0.0, 255.0)
		b = clamp(b*256.0, 0.0, 255.0)

		return color.LedColor{uint8(r), uint8(g), uint8(b), 0xff}
	}

	// Einer der Shader (oder dort: color functions) aus den Testprogrammen
	// von OpenPixelController.

	LavaLampShader = func(t, x, y, z float64, idx, nPix int) color.LedColor {
		y += myCos(x+0.2*z, 0, 1, 0, 0.6)
		z += myCos(x, 0, 1, 0, 0.3)
		x += myCos(y+z, 0, 1.5, 0, 0.2)

		x, y, z = y, z, x

		r := myCos(x, t/4.0, 2, 0, 1)
		g := myCos(y, t/4.0, 2, 0, 1)
		b := myCos(z, t/4.0, 2, 0, 1)
		r, g, b = contrast(r, g, b, 0.5, 1.5)

		r2 := myCos(x, t/10.0+12.345, 3, 0, 1)
		g2 := myCos(y, t/10.0+24.536, 3, 0, 1)
		b2 := myCos(z, t/10.0+34.675, 3, 0, 1)
		clampDown := (r2 + g2 + b2) / 2.0
		clampDown = remap(clampDown, 0.8, 0.9, 0, 1)
		clampDown = clamp(clampDown, 0, 1)
		r *= clampDown
		g *= clampDown
		b *= clampDown

		g = g*0.6 + ((r+b)/2.0)*0.4

		r = clamp(r*256.0, 0.0, 255.0)
		g = clamp(g*256.0, 0.0, 255.0)
		b = clamp(b*256.0, 0.0, 255.0)

		return color.LedColor{uint8(r), uint8(g), uint8(b), 0xff}
	}

	BlinkPeriod = 11.5

	RandomShader = func(t, x, y, z float64, idx, nPix int) color.LedColor {
		var col color.LedColor

		blinkTime := BlinkPeriod * randomList[idx]
		relTime := math.Mod(t, BlinkPeriod)
		if abs(blinkTime-relTime) <= 0.1 {
			col = color.OrangeRed
		} else {
			if relTime < blinkTime {
				relTime += BlinkPeriod
			}
			t := (relTime - blinkTime) / BlinkPeriod
			col = color.OrangeRed.Interpolate(color.Black, t)
		}
		return col
	}

	// Eine Sammlung von Farben-Hilfsfunktionen (ebenfalls aus dem Umfeld von
	// OpenPixelController).

	remap = func(x, minIn, maxIn, minOut, maxOut float64) float64 {
		t := (x - minIn) / (maxIn - minIn)
		return (1.0-t)*minOut + t*maxOut
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
)

// ---------------------------------------------------------------------------

// type ColorSampler struct {
// 	ledgrid.CanvasObjectEmbed
// 	colGrp color.ColorGroup
// }

// func NewColorSampler(colGrp color.ColorGroup) *ColorSampler {
// 	c := &ColorSampler{}
// 	c.CanvasObjectEmbed.Extend(c)
// 	c.colGrp = colGrp
// 	return c
// }

// func (c *ColorSampler) Draw(canv *ledgrid.Canvas) {
// 	for i, colorName := range color.Groups[c.colGrp] {
// 		col := color.Map[colorName]
// 		for j := range 2 {
// 			x := 2*i + j
// 			if x >= width {
// 				return
// 			}
// 			for y := range height {
// 				canv.GC.SetPixel(x, y, col)
// 			}
// 		}
// 	}
// }
