package main

import (
	"context"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
)

func init() {
    // programList.AddTitle("Transitions")
	programList.Add("Fade between three canvases", "Transitions", FadeCanvases)
	programList.Add("Show the different wipe transitions", "Transitions", WipeTrans)
    programList.Add("Like at a theatre...", "Transitions", TheaterKulissen)
	programList.Add("Camera images with some nice fading effects", "Transitions", EffectFaderTest)
}

func FadeCanvases(ctx context.Context, c1 *ledgrid.Canvas) {
	c2, _ := ledGrid.NewCanvas()
	c3, _ := ledGrid.NewCanvas()

	fader1 := ledgrid.NewFadeTransition(c1, 0x00, 0xff, 2*time.Second)
	fader1.AutoReverse = true
	fader1.RepeatCount = ledgrid.AnimationRepeatForever

	fader2 := ledgrid.NewFadeTransition(c2, 0x00, 0xff, 3*time.Second)
	fader2.AutoReverse = true
	fader2.RepeatCount = ledgrid.AnimationRepeatForever

	fader3 := ledgrid.NewFadeTransition(c3, 0x00, 0xff, 5*time.Second)
	fader3.AutoReverse = true
	fader3.RepeatCount = ledgrid.AnimationRepeatForever

	pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}

	img1 := ledgrid.NewImage(pos, "images/img01.png")
	img2 := ledgrid.NewImage(pos, "images/img02.png")
	img3 := ledgrid.NewImage(pos, "images/img03.png")

	c1.Add(img1)
	c2.Add(img2)
	c3.Add(img3)

	fader1.Start()
	fader2.Start()
	fader3.Start()
}

func WipeTrans(ctx context.Context, c1 *ledgrid.Canvas) {
	// c2, _ := ledGrid.NewCanvas()

	pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	size := geom.Point{float64(width), float64(height)}

	cam := NewCamera(pos, size)
	c1.Add(cam)
	cam.Start()

	// img1 := ledgrid.NewImage(pos, "images/img02.png")
	// c1.Add(img1)


	wiper1 := ledgrid.NewWipeTransition(c1, ledgrid.WipeR2L, ledgrid.WipeIn, 3*time.Second)
	wiper1.AutoReverse = true
	wiper2 := ledgrid.NewWipeTransition(c1, ledgrid.WipeL2R, ledgrid.WipeIn, 3*time.Second)
	wiper2.AutoReverse = true
	// wiper5 := ledgrid.NewWipeTransition(c1, ledgrid.WipeT2B, ledgrid.WipeIn, 3*time.Second)
	// wiper5.AutoReverse = true
	// wiper7 := ledgrid.NewWipeTransition(c1, ledgrid.WipeB2T, ledgrid.WipeIn, 3*time.Second)
	// wiper7.AutoReverse = true

    seq := ledgrid.NewSequence(wiper1, wiper2) //, wiper3, wiper5, wiper7)
    seq.RepeatCount = ledgrid.AnimationRepeatForever


	seq.Start()
}

func TheaterKulissen(ctx context.Context, c1 *ledgrid.Canvas) {
	c2, _ := ledGrid.NewCanvas()
	c3, _ := ledGrid.NewCanvas()

	pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	size := geom.Point{float64(width), float64(height)}

    imgCurtain := ledgrid.NewImage(pos, "images/curtain.png")
    c1.Add(imgCurtain)

    imgRocks := ledgrid.NewImage(pos, "images/floor.png")
    c2.Add(imgRocks)

	cam := NewCamera(pos, size)
	c3.Add(cam)
	cam.Start()

}

