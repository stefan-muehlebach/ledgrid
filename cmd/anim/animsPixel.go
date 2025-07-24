package main

import (
	"bufio"
	"context"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"github.com/stefan-muehlebach/gg/fonts"
	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colors"
	"golang.org/x/image/math/fixed"
)

func init() {
	// programList.AddTitle("Pixel Animations")
	programList.Add("Zollweg Biel", "Pixel", ZollwegBiel)
	programList.Add("Moving pixels", "Pixel", MovingPixels)
	programList.Add("Pixel im Stau", "Pixel", CrowdedPixels)
	programList.Add("Glowing pixels with changing text", "Pixel", GlowingPixels)
	programList.Add("Waves of colors", "Pixel", ColorWaves)
	programList.Add("Fireplace", "Pixel", Fireplace)
	programList.Add("Shader using palettes", "Pixel", PaletteShader)
	programList.Add("Shader using colors", "Pixel", ColorShader)
}

var (
	message = "Stefan und Benedict haben sich entschieden"
)

// Hilfsfunktioenchen (sogar generisch!)
func abs[T ~int | ~float64](i T) T {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

func ZollwegBiel(ctx context.Context, c *ledgrid.Canvas) {
	aGrpLedColor := ledgrid.NewGroup()
    	dur := 3 * time.Second
	pal := ledgrid.PaletteMap["BackYellowBlue"]

	wordList := strings.Split(message, " ")
	wordIndex := 0

	for y := range c.Rect.Dy() {
		for x := range c.Rect.Dx() {
			pt := image.Point{x, y}
			pix := ledgrid.NewPixel(pt, colorList[0][0])

			c.Add(pix)

			aColorPal := ledgrid.NewPaletteAnim(pix, pal, dur)
			aColorPal.AutoReverse = true
			aColorPal.Curve = ledgrid.AnimationLinear
			aColorPal.Pos = rand.Float64() / 2.0

			aColorSeq := ledgrid.NewSequence(aColorPal)
			aColorSeq.RepeatCount = ledgrid.AnimationRepeatForever
			aGrpLedColor.Add(aColorSeq)
		}
	}

	pt := geom.Point{float64(width)/2, float64(height)-1.0}
	ptStart := pt.Add(geom.Point{float64(width), 0})
	ptEnd := pt.Sub(geom.Point{float64(width), 0})

	txt := ledgrid.NewText(pt, "", colors.FireBrick)
    txt.SetFont(fonts.SeafordBold, 10.5)
    txt.SetAlign(ledgrid.AlignCenter | ledgrid.AlignBottom)
	c.Add(txt)

	txtLeave := ledgrid.NewPositionAnim(txt, ptEnd, 2*time.Second)
	txtLeave.Curve = ledgrid.AnimationEaseIn
	txtEnter := ledgrid.NewPositionAnim(txt, pt, 2*time.Second)
	txtEnter.Curve = ledgrid.AnimationEaseOut
	txtNewText := ledgrid.NewTask(func() {
		txt.Text = wordList[wordIndex]
    		wordIndex = (wordIndex + 1) % len(wordList)
		txt.Pos = ptStart
	})
	txtChange := ledgrid.NewSequence(/*ledgrid.NewDelay(50 * time.Millisecond),*/
        txtLeave, txtNewText, txtEnter)
    txtChange.RepeatCount = ledgrid.AnimationRepeatForever

	txtChange.Start()
	aGrpLedColor.Start()
}

func CrowdedPixels(ctx context.Context, c *ledgrid.Canvas) {
	mainSeq := ledgrid.NewSequence()
	grp := ledgrid.NewGroup()

	posList := make([]geom.Point, 0)
	for y := 0; y < height; y++ {
		for x := 0; x < width/2; x++ {
			posList = append(posList, geom.Point{float64(x), float64(y)})
		}
	}
	permList := rand.Perm(len(posList))

	min := geom.Point{0, 0}
	max := geom.Point{float64(width / 2), float64(height)}
    sz := max.Sub(min)
    t := sz.X
    if sz.Y > t {
        t = sz.Y
    }
    t += 2.5
	dp := geom.Point{float64(width / 2), 0}
	for i, pos := range posList {
		pixSeq := ledgrid.NewSequence()
		t1 := pos.Distance(min) / t
		t2 := pos.Distance(max) / t
		dest := posList[permList[i]].Add(dp)
		dot := ledgrid.NewDot(pos, colors.LedColor{uint8(255 * t1 * t2), uint8(255 * (1 - t1)), uint8(255 * (1 - t2)), 0xFF})
		c.Add(dot)
		posAnim1 := ledgrid.NewPositionAnim(dot, dest, time.Second+rand.N(time.Second))
		posAnim2 := ledgrid.NewPositionAnim(dot, pos, time.Second+rand.N(time.Second))
		pixSeq.Add(posAnim1, ledgrid.NewDelay(1*time.Second), posAnim2)
		grp.Add(pixSeq)
	}
	mainSeq.Add(ledgrid.NewDelay(time.Second), grp)
	mainSeq.RepeatCount = ledgrid.AnimationRepeatForever
	mainSeq.Start()
}

func MovingPixels(ctx context.Context, c *ledgrid.Canvas) {
	mainSeq := ledgrid.NewSequence()
	grp := ledgrid.NewGroup()

	xMin, xMax := float64(0), float64(width/2-1)
	yMin, yMax := float64(0), float64(height-1)

	// Je zwei gegenueberliegende Seiten des Rechtecks werden in einer Schleife
	// erstellt. Als erstes die horizontalen, d.h. obere und untere Seite.
	colorList := [][]colors.LedColor{
		{colors.RandGroupColor(colors.Blues).Dark(0.1),
			colors.RandGroupColor(colors.Greens).Dark(0.1)},
		{colors.RandGroupColor(colors.Reds).Dark(0.1),
			colors.RandGroupColor(colors.Pinks).Dark(0.1)},
	}
	edgeList := [][]geom.Point{
		{geom.Point{0.0, yMin},
			geom.Point{0.0, yMax}},
		{geom.Point{xMin, 0.0},
			geom.Point{xMax, 0.0}},
	}

	newPos := func(id, side int, v float64) geom.Point {
		pos := edgeList[id][side]
		if id == 0 {
			pos.X = v
		} else {
			pos.Y = v
		}
		return pos
	}

	dp := geom.Point{20, 0}

	posList := make([]geom.Point, 0)

	// Zuerst werden die horizontalen, d.h. die obere und untere Seite
	// erstellt.
	for j := range 2 {
		posList = posList[:0]
		for x := xMin + 1; x <= xMax-1; x++ {
			posList = append(posList, newPos(0, j, float64(x)))
		}
		posPerm := rand.Perm(len(posList))

		for i, pos := range posList {
			pixSeq := ledgrid.NewSequence()
			dest := newPos(0, (j+1)%2, posList[posPerm[i]].X).Add(dp)
			pix := ledgrid.NewDot(pos, colorList[0][j])
			c.Add(pix)
			posAnim1 := ledgrid.NewPositionAnim(pix, dest, time.Second+rand.N(time.Second))
			posAnim2 := ledgrid.NewPositionAnim(pix, pos, time.Second+rand.N(time.Second))
			pixSeq.Add(posAnim1, ledgrid.NewDelay(2*time.Second), posAnim2)
			grp.Add(pixSeq)
		}
	}

	// Anschliessend werden die vertikalen, d.h. linke und rechte Seite
	// erstellt.
	for j := range 2 {
		posList = posList[:0]
		for y := yMin + 1; y <= yMax-1; y++ {
			posList = append(posList, newPos(1, j, float64(y)))
		}
		posPerm := rand.Perm(len(posList))

		for i, pos := range posList {
			pixSeq := ledgrid.NewSequence()
			dest := posList[posPerm[i]].Add(dp)
			pix := ledgrid.NewDot(pos, colorList[1][j])
			c.Add(pix)
			posAnim1 := ledgrid.NewPositionAnim(pix, dest, time.Second+rand.N(time.Second))
			posAnim2 := ledgrid.NewPositionAnim(pix, pos, time.Second+rand.N(time.Second))
			pixSeq.Add(posAnim1, ledgrid.NewDelay(2*time.Second), posAnim2)
			grp.Add(pixSeq)
		}
	}

	mainSeq.Add(grp, ledgrid.NewDelay(time.Second))
	mainSeq.RepeatCount = ledgrid.AnimationRepeatForever
	mainSeq.Start()
}

var (
	colorList = [][]colors.LedColor{
		{colors.LedColor{0xb9, 0xb9, 0x0a, 0xff}, colors.LedColor{0x0a, 0x58, 0x53, 0xff}}, // Yellow to LightBlue
		{colors.LedColor{0xa6, 0x0c, 0x5f, 0xff}, colors.LedColor{0x00, 0x00, 0x80, 0xff}}, // DeepPink to DarkBlue
		{colors.LedColor{0x95, 0x95, 0x00, 0xff}, colors.LedColor{0x71, 0x0e, 0x00, 0xff}}, // Yellow to OrangeRed
		{colors.LedColor{0x00, 0x74, 0x00, 0xff}, colors.LedColor{0x9a, 0x22, 0x22, 0xff}}, // DarkGreen to DarkRed
		{colors.LedColor{0x81, 0x41, 0x24, 0xff}, colors.LedColor{0x4c, 0x8b, 0xaa, 0xff}}, // Salmon to LightBlue
	}
)

func GlowingPixels(ctx context.Context, c *ledgrid.Canvas) {
	aGrpLedColor := ledgrid.NewGroup()
	dur := 3 * time.Second
	numReps := 3

	f, err := os.Open("Faust.txt")
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for y := range c.Rect.Dy() {
		for x := range c.Rect.Dx() {
			// tx := float64(x) / float64(c.Rect.Dx()-1)
			pt := image.Point{x, y}
			pix := ledgrid.NewPixel(pt, colorList[0][0])

			c.Add(pix)

			aColorCyc := ledgrid.NewColorAnim(pix, colorList[0][1], dur)
			aColorCyc.AutoReverse = true
			aColorCyc.RepeatCount = numReps
			aColorCyc.Curve = ledgrid.AnimationLinear
			// aColorCyc.Pos = tx/2.0
			aColorCyc.Pos = rand.Float64() / 2.0

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

	txt := ledgrid.NewFixedText(fixed.P(width/2, height/2), "", colors.White.Alpha(0))
	txt.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)

	txtFadeIn := ledgrid.NewFadeAnim(txt, ledgrid.FadeIn, 500*time.Millisecond)
	txtColorOut := ledgrid.NewColorAnim(txt, colors.Black, 1000*time.Millisecond)
	txtFadeOut := ledgrid.NewFadeAnim(txt, ledgrid.FadeOut, 2000*time.Millisecond)

	txtNextWord := ledgrid.NewTask(func() {
		if scanner.Scan() {
			txt.SetText(scanner.Text())
		}
		txt.Color = colors.White.Alpha(0)
	})
	txtSeq := ledgrid.NewSequence(txtNextWord, txtFadeIn,
		ledgrid.NewDelay(time.Second), txtColorOut, txtFadeOut)
	txtSeq.RepeatCount = ledgrid.AnimationRepeatForever

	c.Add(txt)
	txtSeq.Start()
	aGrpLedColor.Start()
}

func ColorWaves(ctx context.Context, c *ledgrid.Canvas) {
	aGrpLedColor := ledgrid.NewGroup()
	dur := 3 * time.Second
	pal := ledgrid.PaletteMap["Viridis"]

	for y := range c.Rect.Dy() {
		for x := range c.Rect.Dx() {
			pt := image.Point{x, y}
			pix := ledgrid.NewPixel(pt, colorList[0][0])

			c.Add(pix)

			aColorPal := ledgrid.NewPaletteAnim(pix, pal, dur)
			aColorPal.AutoReverse = true
			aColorPal.Curve = ledgrid.AnimationLinear
			aColorPal.Pos = rand.Float64() / 4.0
            if y == 0 && x == 0 {
				fmt.Printf("Pos: %f\n", aColorPal.Pos)
			}

			aColorSeq := ledgrid.NewSequence(aColorPal)
			aColorSeq.RepeatCount = ledgrid.AnimationRepeatForever
			aGrpLedColor.Add(aColorSeq)
		}
	}
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
	var palName string = ledgrid.PaletteNames[0]
	var ptStart, pt, ptEnd fixed.Point26_6

	pt = fixed.P(1, height-1)
	ptStart = pt.Add(fixed.P(width, 0))
	ptEnd = pt.Sub(fixed.P(width, 0))

	txt = ledgrid.NewFixedText(pt, palName, colors.Gold)

	pal := ledgrid.PaletteMap[palName]
	fader := ledgrid.NewPaletteFader(pal)
	aPal := ledgrid.NewPaletteFadeAnim(fader, pal, 2*time.Second)
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
			pix := ledgrid.NewPixel(image.Point{col, row}, colors.Black)
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
	i := 3 // rand.IntN(len(ColorShaderList))

	aGrp := ledgrid.NewGroup()
	nPix := c.Rect.Dx() * c.Rect.Dy()
	idx := 0
	y := yMax
	for row := range c.Rect.Dy() {
		x := xMin
		for col := range c.Rect.Dx() {
			pix := ledgrid.NewPixel(image.Point{col, row}, colors.Black)
			c.Add(pix)
			anim := ledgrid.NewColorShaderAnim(pix, x, y, y, idx, nPix, ColorShaderList[i])
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

	ColorShaderList = []ledgrid.ColorShaderFunc{
		NyanCatShader,
		LavaLampShader,
		RandomShader,
		QuasicrystalShader,
	}

	NyanCatShader = func(t, x, y, z float64, idx, nPix int) colors.LedColor {
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

		return colors.LedColor{uint8(r), uint8(g), uint8(b), 0xff}
	}

	// Einer der Shader (oder dort: color functions) aus den Testprogrammen
	// von OpenPixelController.

	LavaLampShader = func(t, x, y, z float64, idx, nPix int) colors.LedColor {
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

		return colors.LedColor{uint8(r), uint8(g), uint8(b), 0xff}
	}

	BlinkPeriod = 11.5

	RandomShader = func(t, x, y, z float64, idx, nPix int) colors.LedColor {
		var col colors.LedColor

		blinkTime := BlinkPeriod * randomList[idx]
		relTime := math.Mod(t, BlinkPeriod)
		if abs(blinkTime-relTime) <= 0.1 {
			col = colors.OrangeRed
		} else {
			if relTime < blinkTime {
				relTime += BlinkPeriod
			}
			t := (relTime - blinkTime) / BlinkPeriod
			col = colors.OrangeRed.Interpolate(colors.Black, t)
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
		v1 := f1(t, x, y, 5.0)           // old param: 1.2
		v2 := f2(t, x, y, 6.0, 2.0, 3.0) // old param: 1.6, 3.0, 1.5
		v3 := f3(t, x, y, 5.0, 5.0)      // old param: 5.0, 5.0
		v := (v1+v2+v3)/6.0 + 0.5
		return v
	}

	wrap = func(x float64) float64 {
		return math.Abs(math.Mod(x, 2.0) - 1.0)
	}

	wave = func(p geom.Point, angle float64) float64 {
		dir := geom.Point{math.Cos(angle), math.Sin(angle)}
		return math.Cos(p.X*dir.X + p.Y*dir.Y)
	}

	QuasicrystalShader = func(t, x, y, z float64, idx, nPix int) colors.LedColor {
		p := geom.Point{x, y}.Mul(2.0)
		bright := 0.0
		for i := 1.0; i <= 11.0; i += 1.0 {
			bright += wave(p, t/i)
		}
		bright = wrap(bright)
		return colors.LedColor{uint8(255 * bright), uint8(255 * bright), uint8(255 * bright), 0xFF}
	}
)

/*
// Quasicristal

precision mediump float;

varying vec2 position;
uniform float time;

float wave(vec2 p, float angle) {
  vec2 direction = vec2(cos(angle), sin(angle));
  return cos(dot(p, direction));
}

float wrap(float x) {
  return abs(mod(x, 2.)-1.);
}

void main() {
  vec2 p = (position - 0.5) * 50.;

  float brightness = 0.;
  for (float i = 1.; i <= 11.; i++) {
    brightness += wave(p, time / i);
  }

  brightness = wrap(brightness);

  gl_FragColor.rgb = vec3(brightness);
  gl_FragColor.a = 1.;
}

*/

/*
// Noise

precision mediump float;

varying vec2 position;
uniform float time;

float random(float p) {
  return fract(sin(p)*10000.);
}

float noise(vec2 p) {
  return random(p.x + p.y*10000.);
}

vec2 sw(vec2 p) {return vec2( floor(p.x) , floor(p.y) );}
vec2 se(vec2 p) {return vec2( ceil(p.x)  , floor(p.y) );}
vec2 nw(vec2 p) {return vec2( floor(p.x) , ceil(p.y)  );}
vec2 ne(vec2 p) {return vec2( ceil(p.x)  , ceil(p.y)  );}

float smoothNoise(vec2 p) {
  vec2 inter = smoothstep(0., 1., fract(p));
  float s = mix(noise(sw(p)), noise(se(p)), inter.x);
  float n = mix(noise(nw(p)), noise(ne(p)), inter.x);
  return mix(s, n, inter.y);
  return noise(nw(p));
}

float movingNoise(vec2 p) {
  float total = 0.0;
  total += smoothNoise(p     - time);
  total += smoothNoise(p*2.  + time) / 2.;
  total += smoothNoise(p*4.  - time) / 4.;
  total += smoothNoise(p*8.  + time) / 8.;
  total += smoothNoise(p*16. - time) / 16.;
  total /= 1. + 1./2. + 1./4. + 1./8. + 1./16.;
  return total;
}

float nestedNoise(vec2 p) {
  float x = movingNoise(p);
  float y = movingNoise(p + 100.);
  return movingNoise(p + vec2(x, y));
}

void main() {
  vec2 p = position * 6.;
  float brightness = nestedNoise(p);
  gl_FragColor.rgb = vec3(brightness);
  gl_FragColor.a = 1.;
}

*/
