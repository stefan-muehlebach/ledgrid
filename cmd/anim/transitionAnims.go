package main

import (
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
)

var (
	FadeCanvases = NewLedGridProgram("Fade between three canvases",
		func(c1 *ledgrid.Canvas) {
			c2, _ := ledGrid.NewCanvas()
			c3, _ := ledGrid.NewCanvas()

			fader1 := ledgrid.NewFadeTransition(0x00, 0xff, 2*time.Second)
			fader1.AutoReverse = true
			fader1.RepeatCount = ledgrid.AnimationRepeatForever
			c1.Mask = fader1

			fader2 := ledgrid.NewFadeTransition(0x00, 0xff, 3*time.Second)
			fader2.AutoReverse = true
			fader2.RepeatCount = ledgrid.AnimationRepeatForever
			c2.Mask = fader2

			fader3 := ledgrid.NewFadeTransition(0x00, 0xff, 5*time.Second)
			fader3.AutoReverse = true
			fader3.RepeatCount = ledgrid.AnimationRepeatForever
			c3.Mask = fader3

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
		})

	WipeTrans = NewLedGridProgram("Show the different wipe transitions",
		func(c1 *ledgrid.Canvas) {
			c2, _ := ledGrid.NewCanvas()

            wiper1 := ledgrid.NewWipeTransition(c1.Bounds(), 3*time.Second)
            wiper1.AutoReverse = true
            wiper1.RepeatCount = ledgrid.AnimationRepeatForever
            c1.Mask = wiper1

			pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}

			img1 := ledgrid.NewImage(pos, "images/img03.png")
			img2 := ledgrid.NewImage(pos, "images/img02.png")

			c1.Add(img1)
			c2.Add(img2)

            wiper1.Start()
		})
)
