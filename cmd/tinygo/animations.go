package main

import (
	"bufio"
	"image"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/gg/colors"
	"golang.org/x/image/math/fixed"
)

var (
	colorList = [][]colors.RGBA{
		{colors.RGBA{R: 0xb9, G: 0xb9, B: 0x0a}, colors.RGBA{R: 0x0a, G: 0x58, B: 0x53}}, // Yellow to LightBlue
		{colors.RGBA{R: 0xa6, G: 0x0c, B: 0x5f}, colors.RGBA{R: 0x00, G: 0x00, B: 0x80}}, // DeepPink to DarkBlue
		{colors.RGBA{R: 0x95, G: 0x95, B: 0x00}, colors.RGBA{R: 0x71, G: 0x0e, B: 0x00}}, // Yellow to OrangeRed
		{colors.RGBA{R: 0x00, G: 0x74, B: 0x00}, colors.RGBA{R: 0x9a, G: 0x22, B: 0x22}}, // DarkGreen to DarkRed
		{colors.RGBA{R: 0x81, G: 0x41, B: 0x24}, colors.RGBA{R: 0x4c, G: 0x8b, B: 0xaa}}, // Salmon to LightBlue
	}
)

func GlowingPixels(c *ledgrid.Canvas) {
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

	txt := ledgrid.NewFixedText(fixed.P(width/2, height/2), "", colors.Black.Alpha(0.0))
	txt.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
	txtNextWord := ledgrid.NewTask(func() {
		if scanner.Scan() {
			txt.SetText(scanner.Text())
		}
	})
	txtFadeIn := ledgrid.NewFadeAnim(txt, ledgrid.FadeIn, 3*time.Second)
	txtColor := ledgrid.NewColorAnim(txt, colors.White, 1*time.Second)
	txtColor.AutoReverse = true
	txtColor.Cont = false
	txtFadeOut := ledgrid.NewFadeAnim(txt, ledgrid.FadeOut, 1*time.Second)
	txtSeq := ledgrid.NewSequence(txtNextWord, txtFadeIn, txtColor, txtFadeOut)
	txtSeq.RepeatCount = ledgrid.AnimationRepeatForever

	c.Add(txt)
	txtSeq.Start()

	aGrpLedColor.Start()
}
