package main

import (
	"context"
	"image"
	"image/color"
	"image/draw"
	"iter"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/colors"
)

// The function Trace may be used in a range clause at the for statements.
// It allows to iterate over all points between p1 and p2. Since p1 and p2
// are of type image.Point, their coordinates are integers and the points
// generated between p1 and p2 are of type image.Point as well. In order to
// get a continous series of points, I use the line algorithm of Bresenham.
func TraceIntPoint(p1, p2 image.Point) iter.Seq2[int, image.Point] {
	var err, e2 int
	var dp, sp image.Point

	// If p1 and p2 are the same point, there is nothing to iterate over.
	if p1.Eq(p2) {
		return nil
	}
	dp = p2.Sub(p1)
	dp.X = abs(dp.X)
	dp.Y = -abs(dp.Y)

	sp = image.Point{1, 1}
	if p1.X >= p2.X {
		sp.X = -1
	}
	if p1.Y >= p2.Y {
		sp.Y = -1
	}
	err = dp.X + dp.Y

	return func(yield func(id int, pt image.Point) bool) {
		id := 0
		for {
			if !yield(id, p1) || p1.Eq(p2) {
				break
			}
			e2 = 2 * err
			if e2 > dp.Y {
				err += dp.Y
				p1.X += sp.X
			}
			if e2 < dp.X {
				err += dp.X
				p1.Y += sp.Y
			}
			id += 1
		}
	}
}

// This function can be used as a iterator at the range clause of the for
// statement. It's very similar to Trace() (see above) but generates n
// equidistant points between p1 and p2.
func TraceFloatPoint(p1, p2 geom.Point, n int) iter.Seq2[int, geom.Point] {
	return func(yield func(id int, pt geom.Point) bool) {
		id := 0
		for i := 0; i < n; i++ {
			t := float64(i) / float64(n-1)
			pt := p1.Interpolate(p2, t)
			if !yield(id, pt) {
				break
			}
			id += 1
		}
	}
}

// Hilfsfunktioenchen (sogar generisch!)
func abs[T ~int | ~float64](i T) T {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

// ---------------------------------------------------------------------------

// With ScanDir, the direction of a moving effect can be specified.
type ScanDir int

const (
	Top2Bottom ScanDir = iota
	Bottom2Top
	Left2Right
	Right2Left
	Out2Inside
	In2Outside
)

// With LineDir, the direction of the effect on a given row or column is
// specified.
type LineDir int

const (
	Forward LineDir = iota
	Backward
)

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

func EffectFader(typ EffectType) iter.Seq2[int, []PointPair] {
	var l1Count, l1Step int
	var l2Count, l2Step int
	var swapValues bool
	var dDst image.Point
	var minSize, numRounds int
	var ptMin, ptMax, dp image.Point

	minSize = min(width, height)
	numRounds = minSize / 2

	switch typ.scanDir {
	case Top2Bottom:
		l1Count, l1Step = height, 1
		l2Count, l2Step = width, 1
		swapValues = true
		dDst = image.Pt(0, -(height + 2))
	case Bottom2Top:
		l1Count, l1Step = height, -1
		l2Count, l2Step = width, 1
		swapValues = true
		dDst = image.Pt(0, height+2)
	case Left2Right:
		l1Count, l1Step = width, 1
		l2Count, l2Step = height, 1
		swapValues = false
		dDst = image.Pt(-(width + 2), 0)
	case Right2Left:
		l1Count, l1Step = width, -1
		l2Count, l2Step = height, 1
		swapValues = false
		dDst = image.Pt(width+2, 0)
	case Out2Inside:
		dp = image.Point{1, 1}
		ptMin = image.Point{0, 0}
		ptMax = image.Point{width - 1, height - 1}
	case In2Outside:
		dp = image.Point{-1, -1}
		ptMin = image.Point{numRounds - 1, numRounds - 1}
		ptMax = image.Point{width - numRounds, height - numRounds}
	}

	if typ.lineDir == Backward {
		l2Step = -1
	}
	if typ.exitDir == ExitOver {
		dDst = dDst.Mul(-1)
	}

	switch typ.scanDir {
	case Top2Bottom, Bottom2Top, Left2Right, Right2Left:
		return func(yield func(id int, pts []PointPair) bool) {
			pts := make([]PointPair, 0)
			for l1 := range l1Count {
				pts = pts[:0]
				if l1Step < 0 {
					l1 = (l1Count - 1) - l1
				}
				for l2 := range l2Count {
					if l2Step < 0 {
						l2 = (l2Count - 1) - l2
					}
					src := image.Pt(l1, l2)
					if swapValues {
						src.X, src.Y = src.Y, src.X
					}
					dst := src.Add(dDst)
					pts = append(pts, PointPair{src, dst})
				}
				if !yield(l1, pts) {
					break
				}
			}
		}

	case Out2Inside, In2Outside:
		return func(yield func(id int, pts []PointPair) bool) {
			pts := make([]PointPair, 0)
			for rnd := range numRounds {
				pts = pts[:0]
				corners := []image.Point{
					image.Point{ptMin.X, ptMin.Y},
					image.Point{ptMax.X, ptMin.Y},
					image.Point{ptMax.X, ptMax.Y},
					image.Point{ptMin.X, ptMax.Y},
					image.Point{ptMin.X, ptMin.Y},
				}
				if typ.lineDir == Backward {
					slices.Reverse(corners)
				}
				p0 := corners[0]
				for _, p1 := range corners[1:] {
					dDst := p1.Sub(p0)
					dDst.X, dDst.Y = dDst.Y, -dDst.X
					if typ.lineDir == Backward {
						dDst = dDst.Mul(-1)
					}
					if typ.exitDir == ExitAway {
						dDst = dDst.Mul(-1)
					}
					for _, src := range TraceIntPoint(p0, p1) {
						dst := src.Add(dDst)
						pts = append(pts, PointPair{src, dst})
					}
					p0 = p1
				}
				if !yield(rnd, pts) {
					break
				}
				ptMin = ptMin.Add(dp)
				ptMax = ptMax.Sub(dp)
			}
		}
	}
	return nil
}

func EffectFaderTest(ctx context.Context, canv1 *ledgrid.Canvas) {
	pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
	size := geom.Point{float64(width), float64(height)}

	canv2, _ := ledGrid.NewCanvas()
	mask := image.NewAlpha(canv2.Rect)
	draw.Draw(mask, canv2.Rect, image.NewUniform(color.Alpha{0xff}), image.Point{}, draw.Src)
	canv2.Mask = mask

	cam := NewCamera(pos, size)
	canv2.Add(cam)
	cam.Start()

	backgroundTask := func(ctx0 context.Context) {
		var src, dst geom.Point

		effectList := []EffectType{
			{In2Outside, Backward, ExitOver},
			{In2Outside, Forward, ExitOver},
			// {In2Outside, Backward, ExitAway},
			// {In2Outside, Forward, ExitAway},
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
				for _, pts := range EffectFader(effect) {
					for _, pp := range pts {
						p0, p1 := pp.Src, pp.Dst
						src = geom.NewPointIMG(p0)
						dst = geom.NewPointIMG(p1)
						pixAway := ledgrid.NewDot(src, colorList[i].Alpha(0.0))
						canv1.Add(pixAway)

						aMask := ledgrid.NewTask(func() {
							mask.Set(p0.X, p0.Y, color.Alpha{0x00})
						})

						aDur := 200*time.Millisecond + rand.N(300*time.Millisecond)
						aFadeIn := ledgrid.NewFadeAnim(pixAway, ledgrid.FadeIn, aDur)
						aFadeIn.Curve = ledgrid.AnimationCubicIn

						aDur = time.Second + rand.N(time.Second)
						aFadeOut := ledgrid.NewFadeAnim(pixAway, ledgrid.FadeOut, 3*aDur/2)
						aFadeOut.Curve = ledgrid.AnimationEaseIn
						aPos := ledgrid.NewPositionAnim(pixAway, dst, aDur)
						aPos.Curve = ledgrid.AnimationEaseIn
						aSeq := ledgrid.NewSequence(aFadeIn,
							ledgrid.NewGroup(aMask, aFadeOut, aPos),
						)
						aSeq.Start()
						time.Sleep(20 * time.Millisecond)
					}
					select {
					case <-ctx0.Done():
						return
					default:
					}
					// if i%2 == 0 {
					// 	if effect.scanDir == Left2Right || effect.scanDir == Right2Left {
					time.Sleep(250 * time.Millisecond)
					// } else {
					// 	time.Sleep(900 * time.Millisecond)
					// }
					// }
				}
				time.Sleep(3 * time.Second)
				draw.Draw(mask, canv2.Rect, image.NewUniform(color.Alpha{0xff}), image.Point{}, draw.Src)
			}
		}
	}

	go backgroundTask(ctx)
}
