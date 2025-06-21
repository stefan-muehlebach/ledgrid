package main

import (
	"context"
	"image"
	"image/color"
	"image/draw"
	"iter"
	"math/rand/v2"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colors"
)

// ---------------------------------------------------------------------------

// With ScanDir, the direction of a moving effect can be specified.
type ScanDir int

const (
	Top2Bottom ScanDir = iota
	Bottom2Top
	Left2Right
	Right2Left
	// Out2Inside
	// In2Outside
)

// Within a given line (Top2Bottom or Bottom2Top) or a given column
// (Left2Right or Right2Left) the effect can move 'forward' or 'backward' or
// alternate after each line/column. Forward means 'increasing X or Y value'
// while backward means 'decreasing X or Y value'.
type LineDir int

const (
	Forward LineDir = iota
	Backward
	Alternate
)

// ExitDir specifies the direction, in which the pixels will exit the screen.
// ExitAway guides them away for the image over already masked area,
// while ExitOver guides the pixel over the image.
type ExitDir int

const (
	ExitAway ExitDir = iota
	ExitOver
)

type EffectType struct {
	scanDir ScanDir
	lineDir LineDir
	exitDir ExitDir
}

type PointPair struct {
	Src, Dst image.Point
}

func EffectFader(typ EffectType, size image.Point) iter.Seq2[int, []PointPair] {
	var l1Count, l1Step int
	var l2Count, l2Step int
	var dDst image.Point
	var lastDir LineDir
	w, h := size.X, size.Y

	if typ.lineDir == Alternate {
		lastDir = Forward
	} else {
		lastDir = typ.lineDir
	}

	switch typ.scanDir {
	case Top2Bottom:
		l1Count, l1Step = h, 1
		l2Count, l2Step = w, 1
		dDst = image.Pt(0, -(h + 2))
	case Bottom2Top:
		l1Count, l1Step = h, -1
		l2Count, l2Step = w, 1
		dDst = image.Pt(0, h+2)
	case Left2Right:
		l1Count, l1Step = w, 1
		l2Count, l2Step = h, 1
		dDst = image.Pt(-(w + 2), 0)
	case Right2Left:
		l1Count, l1Step = w, -1
		l2Count, l2Step = h, 1
		dDst = image.Pt(w+2, 0)
	}

	if typ.exitDir == ExitOver {
		dDst = dDst.Mul(-1)
	}

	return func(yield func(id int, pts []PointPair) bool) {
		pts := make([]PointPair, 0)
		for l1 := range l1Count {
			if lastDir == Forward {
				l2Step = +1
			} else {
				l2Step = -1
			}

			pts = pts[:0]
			if l1Step < 0 {
				l1 = (l1Count - 1) - l1
			}
			for l2 := range l2Count {
				if l2Step < 0 {
					l2 = (l2Count - 1) - l2
				}
				src := image.Pt(l1, l2)
				if typ.scanDir == Top2Bottom || typ.scanDir == Bottom2Top {
					src.X, src.Y = src.Y, src.X
				}
				dst := src
				if typ.scanDir == Top2Bottom || typ.scanDir == Bottom2Top {
					dst.Y = dDst.Y
				} else {
					dst.X = dDst.X
				}
				pts = append(pts, PointPair{src, dst})
			}
			if !yield(l1, pts) {
				break
			}
			if typ.lineDir == Alternate {
				if lastDir == Forward {
					lastDir = Backward
				} else {
					lastDir = Forward
				}
			}
		}
	}
	return nil
}

func EffectFaderTest(ctx context.Context, canv1 *ledgrid.Canvas) {
	pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	size := geom.Point{float64(width), float64(height)}

	opaque := image.NewUniform(color.Alpha{0xff})

	canv2, _ := ledGrid.NewCanvas()
	mask := image.NewAlpha(canv2.Rect)
	draw.Draw(mask, canv2.Rect, opaque, image.Point{}, draw.Src)
	canv2.Mask = mask

	cam := NewCamera(pos, size)
	canv2.Add(cam)
	cam.Start()

	backgroundTask := func(ctx0 context.Context) {
		var src, dst geom.Point

		effectList := []EffectType{
			// {In2Outside, Backward, ExitOver},
			// {In2Outside, Forward, ExitOver},
			// {In2Outside, Backward, ExitAway},
			// {In2Outside, Forward, ExitAway},
			{Top2Bottom, Alternate, ExitOver},
			{Top2Bottom, Forward, ExitOver},
			{Bottom2Top, Backward, ExitAway},
			{Left2Right, Forward, ExitAway},
			{Right2Left, Backward, ExitOver},
		}
		colorList := []colors.LedColor{
			colors.OrangeRed,
			colors.Crimson,
			colors.FireBrick,
			colors.Gold,
			colors.BlueViolet,
			colors.SkyBlue,
			colors.Lime,
			colors.YellowGreen,
			colors.Teal,
		}

		for {
			time.Sleep(1 * time.Second)
			for i, effect := range effectList {
				ledgrid.AnimCtrl.Purge()
				canv1.Purge()
				for _, pts := range EffectFader(effect, canv1.Rect.Size()) {
					for _, pp := range pts {
						p0, p1 := pp.Src, pp.Dst
						src = geom.NewPointIMG(p0)
						dst = geom.NewPointIMG(p1)
						pixAway := ledgrid.NewDot(src, colorList[i].Alpha(0.0))
						canv1.Add(pixAway)

						aDur := 200*time.Millisecond + rand.N(300*time.Millisecond)
						aFadeIn := ledgrid.NewFadeAnim(pixAway, ledgrid.FadeIn, aDur)
						aFadeIn.Curve = ledgrid.AnimationCubicIn

						aDur = time.Second + rand.N(time.Second)
						aPos := ledgrid.NewPositionAnim(pixAway, dst, aDur)
						aPos.Curve = ledgrid.AnimationCubicIn
						aSeq := ledgrid.NewSequence(
							aFadeIn,
							ledgrid.NewGroup(ledgrid.NewTask(func() {
								mask.Set(p0.X, p0.Y, color.Alpha{0x00})
							}), aPos),
						)

						aSeq.Start()
						time.Sleep(20 * time.Millisecond)
					}

					select {
					case <-ctx0.Done():
						return
					default:
					}

					time.Sleep(250 * time.Millisecond)
				}
				time.Sleep(3 * time.Second)
				draw.Draw(mask, canv2.Rect, opaque, image.Point{}, draw.Src)
			}
		}
	}

	go backgroundTask(ctx)
}
