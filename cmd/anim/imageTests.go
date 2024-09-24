package main

import (
	"math"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
)

var (
	BlinkenAnimation = NewLedGridProgram("Blinken animation",
		func(c *ledgrid.Canvas) {
			posFlame1 := geom.Point{4.5, float64(height)}
			posFlame2 := geom.Point{float64(width) - 4.5, float64(height)}
			pos1Mario := geom.Point{0.0, float64(height)}
			pos2Mario := geom.Point{float64(width), float64(height)}

			bmlFlame := ledgrid.ReadBlinkenFile("blinken/flameNew.bml")
			bmlFlame.SetAllDuration(32)

			flame1 := ledgrid.NewSprite(posFlame1)
			flame1.SetAlign(ledgrid.AlignCenter | ledgrid.AlignBottom)
			flame1.AddBlinkenLight(bmlFlame)
			flame1.RepeatCount = ledgrid.AnimationRepeatForever

			bmlFlame.SetAllDuration(43)

			flame2 := ledgrid.NewSprite(posFlame2)
			flame2.SetAlign(ledgrid.AlignCenter | ledgrid.AlignBottom)
			flame2.AddBlinkenLight(bmlFlame)
			flame2.RepeatCount = ledgrid.AnimationRepeatForever

			bmlMario := ledgrid.ReadBlinkenFile("blinken/marioWalkRight.bml")

			mario := ledgrid.NewSprite(pos1Mario)
			mario.SetAlign(ledgrid.AlignCenter | ledgrid.AlignBottom)
			mario.AddBlinkenLight(bmlMario)
			mario.RepeatCount = ledgrid.AnimationRepeatForever
			mario.Size = geom.Point{20.0, 20.0}

			aPos := ledgrid.NewPositionAnimation(&mario.Pos, pos2Mario, 4*time.Second)
			aPos.Curve = ledgrid.AnimationLinear
			aPos.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(flame1, flame2, mario)

			aGrp := ledgrid.NewGroup(flame1, flame2, mario, aPos)
			aGrp.Start()
		})

	SlideTheShow = NewLedGridProgram("Slide-the-Show",
		func(c *ledgrid.Canvas) {
			pos := geom.Point{float64(width / 2), float64(height / 2)}
			files := []string{
				"images/raster.png",
				"images/square1.png",
				"images/square2.png",
			}
			aTimeline := ledgrid.NewTimeline(time.Duration(len(files)) * 4 * time.Second)
			dstSize := geom.NewPointIMG(c.Bounds().Size())
			dstRatio := dstSize.X / dstSize.Y
			for i, fileName := range files {
				img := ledgrid.NewImage(pos, fileName)
				img.Hide()
				srcRatio := float64(img.Img.Bounds().Dx()) / float64(img.Img.Bounds().Dy())
				if dstRatio > srcRatio {
					h := dstSize.Y
					w := h * srcRatio
					img.Size = geom.Point{w, h}
				} else {
					w := dstSize.X
					h := w / srcRatio
					img.Size = geom.Point{w, h}
				}
				t0 := time.Duration(4*i+1) * time.Second
				t1 := t0 + 300*time.Millisecond
				t2 := t1 + 3300*time.Millisecond
				aTimeline.Add(t0, ledgrid.NewHideShowAnimation(img))
				aTimeline.Add(t1, ledgrid.NewFloatAnimation(&img.Angle, 6*math.Pi, 3*time.Second))
				aTimeline.Add(t2, ledgrid.NewHideShowAnimation(img))
				c.Add(img)
			}
			aTimeline.RepeatCount = ledgrid.AnimationRepeatForever
			aTimeline.Start()
		})

	SingleImageAlign = NewLedGridProgram("Align this lonely image!",
		func(c *ledgrid.Canvas) {
			imgPos := geom.Point{float64(width / 2), float64(height / 2)}
			img := ledgrid.NewImage(imgPos, "images/raster.png")
			img.Size = geom.Point{20, 20}
			img.SetAlign(ledgrid.AlignBottom)
			c.Add(img)

			aAlignRight := ledgrid.NewTask(func() {
				img.SetAlign(ledgrid.AlignRight)
			})
			aAlignCenter := ledgrid.NewTask(func() {
				img.SetAlign(ledgrid.AlignCenter)
			})
			aAlignLeft := ledgrid.NewTask(func() {
				img.SetAlign(ledgrid.AlignLeft)
			})
			aAlignBottom := ledgrid.NewTask(func() {
				img.SetAlign(ledgrid.AlignBottom)
			})
			aAlignMiddle := ledgrid.NewTask(func() {
				img.SetAlign(ledgrid.AlignMiddle)
			})
			aAlignTop := ledgrid.NewTask(func() {
				img.SetAlign(ledgrid.AlignTop)
			})

			aAngle := ledgrid.NewFloatAnimation(&img.Angle, 2*math.Pi, 4*time.Second)
            aAngle.Curve = ledgrid.AnimationLazeInOut
			aHoriSeq := ledgrid.NewSequence(
				aAlignRight, aAngle,
				aAlignCenter, aAngle,
				aAlignLeft, aAngle,
			)
			aVertSeq := ledgrid.NewSequence(
				aAlignBottom, aHoriSeq,
				aAlignMiddle, aHoriSeq,
				aAlignTop, aHoriSeq,
			)
			aVertSeq.Start()
		})
)
