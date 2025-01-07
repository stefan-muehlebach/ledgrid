package main

import (
	"bufio"
	"image"
	"log"
	"math/rand/v2"
	"os"
	"time"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
	"golang.org/x/image/math/fixed"
)

var (
	colorList = [][]color.LedColor{
		{color.NewLedColorHex(0xb9b90a), color.NewLedColorHex(0x0a5853)}, // Yellow to LightBlue
		{color.NewLedColorHex(0xa60c5f), color.NewLedColorHex(0x000080)}, // DeepPink to DarkBlue
		{color.NewLedColorHex(0x959500), color.NewLedColorHex(0x710e00)}, // Yellow to OrangeRed
		{color.NewLedColorHex(0x007400), color.NewLedColorHex(0x9a2222)}, // DarkGreen to DarkRed
		{color.NewLedColorHex(0x814124), color.NewLedColorHex(0x4c8baa)}, // Salmon to LightBlue
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

	txt := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.Black.Alpha(0.0), "")
	txt.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
	txtNextWord := ledgrid.NewTask(func() {
		if scanner.Scan() {
			txt.SetText(scanner.Text())
		}
	})
	txtFadeIn := ledgrid.NewFadeAnim(txt, ledgrid.FadeIn, 3*time.Second)
	txtColor := ledgrid.NewColorAnim(txt, color.White, 1*time.Second)
	txtColor.AutoReverse = true
	txtColor.Cont = false
	txtFadeOut := ledgrid.NewFadeAnim(txt, ledgrid.FadeOut, 1*time.Second)
	txtSeq := ledgrid.NewSequence(txtNextWord, txtFadeIn, txtColor, txtFadeOut)
	txtSeq.RepeatCount = ledgrid.AnimationRepeatForever

	c.Add(txt)
	txtSeq.Start()

	aGrpLedColor.Start()
}
