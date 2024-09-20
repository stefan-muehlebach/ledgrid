package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"math"
	"math/rand/v2"
	"os"
	"os/signal"
	"slices"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
	"github.com/stefan-muehlebach/ledgrid/conf"
	"golang.org/x/image/math/fixed"
)

const (
	defHost = "raspi-3"
	defPort = 5333
)

var (
	width, height int
	gridSize      image.Point
	backAlpha     = 1.0
	ledGrid       *ledgrid.LedGrid
	canvas        *ledgrid.Canvas
)

//----------------------------------------------------------------------------

type LedGridProgram interface {
	Name() string
	Run(c *ledgrid.Canvas)
}

func NewLedGridProgram(name string, runFunc func(c *ledgrid.Canvas)) LedGridProgram {
	return &simpleProgram{name, runFunc}
}

type simpleProgram struct {
	name    string
	runFunc func(c *ledgrid.Canvas)
}

func (p *simpleProgram) Name() string {
	return p.name
}

func (p *simpleProgram) Run(c *ledgrid.Canvas) {
	p.runFunc(c)
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

type TraceFunc func(yield func(id int, pt image.Point) bool)

func Trace(p1, p2 image.Point) TraceFunc {
	if (p1.X != p2.X && p1.Y != p2.Y) || (p1.X == p2.X && p1.Y == p2.Y) {
		return nil
	}
	dp := p2.Sub(p1)
	if dp.X != 0 {
		dp.X /= abs(dp.X)
	} else {
		dp.Y /= abs(dp.Y)
	}
	return func(yield func(id int, pt image.Point) bool) {
		id := 0
		for pt := p1; pt != p2; pt = pt.Add(dp) {
			if !yield(id, pt) {
				break
			}
			id += 1
		}
	}
}

// ---------------------------------------------------------------------------

type ScanDir int
type LineDir int
type ExitDir int

const (
	Top2Bottom ScanDir = iota
	Bottom2Top
	Left2Right
	Right2Left
	Out2Inside
	In2Outside

	Forward LineDir = iota
	CW
	Backward
	CCW
	Alternate

	ExitAway ExitDir = iota
	ExitOver
)

type DissolveType struct {
	dir     ScanDir
	lineDir LineDir
	exitDir ExitDir
}

type PointPair struct {
	src, dst image.Point
}

type DissolverFunc func(yield func(row int, pts []PointPair) bool)

func Dissolver(typ DissolveType) DissolverFunc {
	var l1Count, l1Step int
	var l2Count, l2Step int
	var swapValues bool
	var dDst image.Point
	var minSize, numRounds int
	var ptMin, ptMax, dp image.Point

	minSize = min(width, height)
	numRounds = minSize / 2

	switch typ.dir {
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

	switch typ.dir {
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
					for _, src := range Trace(p0, p1) {
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

			// for x, y := ptMin.X, ptMin.Y; x < ptMax.X; x++ {
			// 	src := image.Pt(x, y)
			// 	dst := src.Add(image.Pt(0, -height))
			// 	pts = append(pts, PointPair{src, dst})
			// }
			// for x, y := ptMax.X, ptMin.Y; y < ptMax.Y; y++ {
			// 	src := image.Pt(x, y)
			// 	dst := src.Add(image.Pt(width, 0))
			// 	pts = append(pts, PointPair{src, dst})
			// }
			// for x, y := ptMax.X, ptMax.Y; x > ptMin.X; x-- {
			// 	src := image.Pt(x, y)
			// 	dst := src.Add(image.Pt(0, height))
			// 	pts = append(pts, PointPair{src, dst})
			// }
			// for x, y := ptMin.X, ptMax.Y; y > ptMin.Y; y-- {
			// 	src := image.Pt(x, y)
			// 	dst := src.Add(image.Pt(-width, 0))
			// 	pts = append(pts, PointPair{src, dst})
			// }
			// if !yield(rnd, pts) {
			// 	break
			// }
			// ptMin = ptMin.Add(dp)
			// ptMax = ptMax.Sub(dp)
			// }
		}
	}
	return nil
}

// func SpiralDissolverO2I(yield func(row int, pts []PointPair) bool) {
// 	minSize := min(width, height)
// 	numRounds := minSize / 2

// 	dp := image.Point{1, 1}
// 	ptMin := image.Point{0, 0}
// 	ptMax := image.Point{width - 1, height - 1}

// 	dp := image.Point{-1, -1}
// 	ptMin := image.Point{numRounds - 1, numRounds - 1}
// 	ptMax := image.Point{width - numRounds, height - numRounds}

// 	return func(yield func(id int, pts []PointPair) bool) {
// 		pts := make([]PointPair, 0)
// 		for rnd := range numRounds {
// 			pts = pts[:0]
// 			for x, y := ptMin.X, ptMin.Y; x < ptMax.X; x++ {
// 				src := image.Pt(x, y)
// 				dst := src.Add(image.Pt(0, -height))
// 				pts = append(pts, PointPair{src, dst})
// 			}
// 			for x, y := ptMax.X, ptMin.Y; y < ptMax.Y; y++ {
// 				src := image.Pt(x, y)
// 				dst := src.Add(image.Pt(width, 0))
// 				pts = append(pts, PointPair{src, dst})
// 			}
// 			for x, y := ptMax.X, ptMax.Y; x > ptMin.X; x-- {
// 				src := image.Pt(x, y)
// 				dst := src.Add(image.Pt(0, height))
// 				pts = append(pts, PointPair{src, dst})
// 			}
// 			for x, y := ptMin.X, ptMax.Y; y > ptMin.Y; y-- {
// 				src := image.Pt(x, y)
// 				dst := src.Add(image.Pt(-width, 0))
// 				pts = append(pts, PointPair{src, dst})
// 			}
// 			if !yield(rnd, pts) {
// 				break LOOP
// 			}
// 			ptMin = ptMin.Add(dp)
// 			ptMax = ptMax.Sub(dp)
// 		}
// 	}
// }

// func SpiralDissolverI2O(yield func(row int, pts []PointPair) bool) {
// 	pts := make([]PointPair, 0)
// 	minSize := min(width, height)
// 	numRounds := minSize / 2
// 	dp := image.Point{-1, -1}
// 	ptMin := image.Point{numRounds - 1, numRounds - 1}
// 	ptMax := image.Point{width - numRounds, height - numRounds}
// LOOP:
// 	for rnd := range numRounds {
// 		// fmt.Printf("spiral round %d\n", rnd)
// 		pts = pts[:0]
// 		// fmt.Printf("  side 1\n")
// 		for x, y := ptMin.X, ptMin.Y; x < ptMax.X; x++ {
// 			src := image.Pt(x, y)
// 			dst := src.Add(image.Pt(0, -height))
// 			pts = append(pts, PointPair{src, dst})
// 		}
// 		// fmt.Printf("  side 2\n")
// 		for x, y := ptMax.X, ptMin.Y; y < ptMax.Y; y++ {
// 			src := image.Pt(x, y)
// 			dst := src.Add(image.Pt(width, 0))
// 			pts = append(pts, PointPair{src, dst})
// 		}
// 		// fmt.Printf("  side 3\n")
// 		for x, y := ptMax.X, ptMax.Y; x > ptMin.X; x-- {
// 			src := image.Pt(x, y)
// 			dst := src.Add(image.Pt(0, height))
// 			pts = append(pts, PointPair{src, dst})
// 		}
// 		// fmt.Printf("  side 4\n")
// 		for x, y := ptMin.X, ptMax.Y; y > ptMin.Y; y-- {
// 			src := image.Pt(x, y)
// 			dst := src.Add(image.Pt(-width, 0))
// 			pts = append(pts, PointPair{src, dst})
// 		}
// 		// fmt.Printf("spiral yielding back\n")
// 		if !yield(rnd, pts) {
// 			break LOOP
// 		}
// 		ptMin = ptMin.Add(dp)
// 		ptMax = ptMax.Sub(dp)
// 	}
// }

// ---------------------------------------------------------------------------

var (
	FarewellGery = NewLedGridProgram("Farewell Gery!",
		func(c *ledgrid.Canvas) {
			aGrpFadeIn := ledgrid.NewGroup()
			aGrpPurple := ledgrid.NewGroup()
			aGrpYellow := ledgrid.NewGroup()
			aGrpGreen := ledgrid.NewGroup()
			aGrpGrey := ledgrid.NewGroup()
			aGrpRed := ledgrid.NewGroup()
			aGrpBlack := ledgrid.NewGroup()
			aSeqColor := ledgrid.NewSequence(aGrpRed, aGrpGreen)

			for y := range c.Rect.Dy() {
				for x := range c.Rect.Dx() {
					pt := image.Point{x, y}
					pos := geom.NewPointIMG(pt)
					t := rand.Float64()
					col := color.Black
					pix := ledgrid.NewDot(pos, col)
					c.Add(pix)

					dur := time.Second + time.Duration(10*x+20*y)*time.Millisecond
					aAlpha := ledgrid.NewFadeAnimation(&pix.Color.A, 196, dur)
					aAlpha.AutoReverse = true
					aAlpha.RepeatCount = ledgrid.AnimationRepeatForever
					aAlpha.Start()

					aColor := ledgrid.NewColorAnimation(&pix.Color, (color.DimGray.Dark(0.5)).Interpolate((color.DarkGrey.Dark(0.5)), t), 9*time.Second)
					aColor.Cont = true
					aGrpFadeIn.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, (color.DimGray.Dark(0.5)).Interpolate((color.DarkGrey.Dark(0.5)), t), 1*time.Second)
					aColor.Cont = true
					aGrpGrey.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.MediumPurple.Interpolate(color.Fuchsia, t), 4*time.Second)
					aColor.Cont = true
					aGrpPurple.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.Gold.Interpolate(color.LemonChiffon, t), 4*time.Second)
					aColor.Cont = true
					aGrpYellow.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.Crimson.Interpolate(color.Orange, t), 4*time.Second)
					aColor.Cont = true
					aGrpRed.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.LightSeaGreen.Interpolate(color.GreenYellow, t), 500*time.Millisecond)
					aColor.Cont = true
					aGrpGreen.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.Black, 2*time.Second)
					aColor.Cont = true
					aGrpBlack.Add(aColor)
				}
			}

			txt1 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.GreenYellow.Alpha(0.0), "LIEBER")
			aTxt1 := ledgrid.NewFadeAnimation(&txt1.Color.A, ledgrid.FadeIn, 1*time.Second)
			aTxt1.AutoReverse = true
			txt2 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.DarkViolet.Alpha(0.0), "GERY")
			aTxt2 := ledgrid.NewFadeAnimation(&txt2.Color.A, ledgrid.FadeIn, 2*time.Second)
			aTxt2.AutoReverse = true
			txt3 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.RoyalBlue.Alpha(0.0), "MACH'S GUET")
			aTxt3 := ledgrid.NewFadeAnimation(&txt3.Color.A, ledgrid.FadeIn, 5*time.Second)
			aTxt3.AutoReverse = true
			c.Add(txt1, txt2, txt3)

			aTimel := ledgrid.NewTimeline(40 * time.Second)
			aTimel.Add(0, aGrpFadeIn)
			aTimel.Add(10*time.Second, aGrpPurple)
			aTimel.Add(13*time.Second, aTxt1)
			aTimel.Add(14*time.Second, aGrpGrey)

			aTimel.Add(20*time.Second, aGrpYellow)
			aTimel.Add(23*time.Second, aTxt2)
			aTimel.Add(24*time.Second, aGrpGrey)

			aTimel.Add(30*time.Second, aSeqColor)
			aTimel.Add(33*time.Second, aTxt3)
			aTimel.Add(35*time.Second, aGrpBlack)
			aTimel.RepeatCount = ledgrid.AnimationRepeatForever

			aTimel.Start()
		})

	GroupTest = NewLedGridProgram("Group test",
		func(c *ledgrid.Canvas) {
			rPos1 := geom.Point{5.0, float64(height) / 2.0}
			rPos2 := geom.Point{float64(width) - 5.0, float64(height) / 2.0}
			rSize1 := geom.Point{7.0, 7.0}
			rSize2 := geom.Point{1.0, 1.0}
			rColor1 := color.SkyBlue
			rColor2 := color.GreenYellow

			r := ledgrid.NewRectangle(rPos1, rSize1, rColor1)
			c.Add(r)

			aPos := ledgrid.NewPositionAnimation(&r.Pos, rPos2, time.Second)
			aPos.AutoReverse = true
			aSize := ledgrid.NewSizeAnimation(&r.Size, rSize2, 2*time.Second)
			aSize.AutoReverse = true
			aColor := ledgrid.NewColorAnimation(&r.StrokeColor, rColor2, 2*time.Second)
			aColor.AutoReverse = true
			aAngle := ledgrid.NewFloatAnimation(&r.Angle, math.Pi, 2*time.Second)
			aAngle.AutoReverse = true

			aGroup := ledgrid.NewGroup(aPos, aSize, aColor, aAngle)
			aGroup.RepeatCount = ledgrid.AnimationRepeatForever

			aGroup.Start()
		})

	SequenceTest = NewLedGridProgram("Sequence test",
		func(c *ledgrid.Canvas) {
			rPos := geom.NewPointIMG(gridSize).Mul(0.5)
			rSize1 := geom.NewPointIMG(gridSize).SubXY(1, 1)
			rSize2 := geom.Point{5.0, 3.0}

			r := ledgrid.NewRectangle(rPos, rSize1, color.SkyBlue)
			c.Add(r)

			aSize1 := ledgrid.NewSizeAnimation(&r.Size, rSize2, time.Second)
			aColor1 := ledgrid.NewColorAnimation(&r.StrokeColor, color.OrangeRed, time.Second/2)
			aColor1.AutoReverse = true
			aColor2 := ledgrid.NewColorAnimation(&r.StrokeColor, color.Crimson, time.Second/2)
			aColor2.AutoReverse = true
			aColor3 := ledgrid.NewColorAnimation(&r.StrokeColor, color.Coral, time.Second/2)
			aColor3.AutoReverse = true
			aColor4 := ledgrid.NewColorAnimation(&r.StrokeColor, color.FireBrick, time.Second/2)
			aSize2 := ledgrid.NewSizeAnimation(&r.Size, rSize1, time.Second)
			aSize2.Cont = true
			aColor5 := ledgrid.NewColorAnimation(&r.StrokeColor, color.SkyBlue, time.Second)
			aColor5.Cont = true

			aSeq := ledgrid.NewSequence(aSize1, aColor1, aColor2, aColor3, aColor4, aSize2, aColor5)
			aSeq.RepeatCount = ledgrid.AnimationRepeatForever
			aSeq.Start()
		})

	TimelineTest = NewLedGridProgram("Timeline test",
		func(c *ledgrid.Canvas) {
			r1Pos := geom.Point{6, float64(height) / 2.0}
			r1Size := geom.Point{9.0, 3.0}
			r2Pos := geom.Point{float64(width)/2.0 - 3.0, float64(height) / 2.0}
			r2Size := geom.Point{3.0, 9.0}
			r3Pos := geom.Point{float64(width) - 6, float64(height) / 2.0}
			r3Size := geom.Point{9.0, 3.0}
			r4Pos := geom.Point{float64(width)/2.0 + 3.0, float64(height) / 2.0}
			r4Size := geom.Point{3.0, 9.0}

			r1 := ledgrid.NewRectangle(r1Pos, r1Size, color.GreenYellow)
			r2 := ledgrid.NewRectangle(r2Pos, r2Size, color.Gold)
			r3 := ledgrid.NewRectangle(r3Pos, r3Size, color.SkyBlue)
			r4 := ledgrid.NewRectangle(r4Pos, r4Size, color.Gold)
			c.Add(r1, r3, r2, r4)

			aAngle1 := ledgrid.NewFloatAnimation(&r1.Angle, math.Pi, time.Second)
			aAngle2 := ledgrid.NewFloatAnimation(&r1.Angle, 0.0, time.Second)
			aAngle2.Cont = true

			aColor1 := ledgrid.NewColorAnimation(&r1.StrokeColor, color.OrangeRed, 200*time.Millisecond)
			aColor1.AutoReverse = true
			aColor1.RepeatCount = 3
			aColor2 := ledgrid.NewColorAnimation(&r1.StrokeColor, color.Purple, 500*time.Millisecond)
			aColor3 := ledgrid.NewColorAnimation(&r1.StrokeColor, color.GreenYellow, 500*time.Millisecond)
			aColor3.Cont = true

			aPos1 := ledgrid.NewPositionAnimation(&r1.Pos, r2.Pos.SubXY(r2Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos1.AutoReverse = true

			aAngle3 := ledgrid.NewFloatAnimation(&r3.Angle, -math.Pi, time.Second)
			aAngle4 := ledgrid.NewFloatAnimation(&r3.Angle, 0.0, time.Second)
			aAngle4.Cont = true

			aColor4 := ledgrid.NewColorAnimation(&r3.StrokeColor, color.DarkOrange, 200*time.Millisecond)
			aColor4.AutoReverse = true
			aColor4.RepeatCount = 3
			aColor5 := ledgrid.NewColorAnimation(&r3.StrokeColor, color.Purple, 500*time.Millisecond)
			aColor6 := ledgrid.NewColorAnimation(&r3.StrokeColor, color.SkyBlue, 500*time.Millisecond)
			aColor6.Cont = true

			aPos2 := ledgrid.NewPositionAnimation(&r3.Pos, r4.Pos.AddXY(r4Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos2.AutoReverse = true

			aColor7 := ledgrid.NewColorAnimation(&r2.StrokeColor, color.Cornsilk, 500*time.Millisecond)
			aColor7.AutoReverse = true
			aBorder1 := ledgrid.NewFloatAnimation(&r2.StrokeWidth, 2.0, 500*time.Millisecond)
			aBorder1.AutoReverse = true

			aColor8 := ledgrid.NewColorAnimation(&r4.StrokeColor, color.Cornsilk, 500*time.Millisecond)
			aColor8.AutoReverse = true
			aBorder2 := ledgrid.NewFloatAnimation(&r4.StrokeWidth, 2.0, 500*time.Millisecond)
			aBorder2.AutoReverse = true

			tl := ledgrid.NewTimeline(5 * time.Second)
			tl.RepeatCount = ledgrid.AnimationRepeatForever

			// Timeline positions for the first rectangle
			tl.Add(300*time.Millisecond, aColor1)
			tl.Add(1800*time.Millisecond, aAngle1)
			tl.Add(2300*time.Millisecond, aColor2, aPos1)
			tl.Add(2400*time.Millisecond, aColor7, aBorder1)
			tl.Add(2900*time.Millisecond, aAngle2)
			tl.Add(3400*time.Millisecond, aColor3)

			// Timeline positions for the second rectangle
			tl.Add(500*time.Millisecond, aColor4)
			tl.Add(2000*time.Millisecond, aAngle3)
			tl.Add(2500*time.Millisecond, aColor5, aPos2)
			tl.Add(2600*time.Millisecond, aColor8, aBorder2)
			tl.Add(3100*time.Millisecond, aAngle4)
			tl.Add(3600*time.Millisecond, aColor6)

			tl.Start()
		})

	StopContShowHideTest = NewLedGridProgram("Hide/Show vs. Suspend/Continue by Tasks",
		func(c *ledgrid.Canvas) {
			rPos1 := geom.Point{5.0, float64(height) / 2.0}
			rPos2 := geom.Point{float64(width) - 5.0, float64(height) / 2.0}
			rSize1 := geom.Point{7.0, 7.0}
			rSize2 := geom.Point{6.0, 6.0}
			rColor1 := color.SkyBlue
			rColor2 := color.GreenYellow

			r := ledgrid.NewRectangle(rPos1, rSize1, rColor1)
			c.Add(r)

			aPos := ledgrid.NewPositionAnimation(&r.Pos, rPos2, 4*time.Second)
			aPos.AutoReverse = true
			aPos.RepeatCount = ledgrid.AnimationRepeatForever
			aSize := ledgrid.NewSizeAnimation(&r.Size, rSize2, 4*time.Second)
			aSize.AutoReverse = true
			aSize.RepeatCount = ledgrid.AnimationRepeatForever
			aColor := ledgrid.NewColorAnimation(&r.StrokeColor, rColor2, 4*time.Second)
			aColor.AutoReverse = true
			aColor.RepeatCount = ledgrid.AnimationRepeatForever
			aAngle := ledgrid.NewFloatAnimation(&r.Angle, math.Pi, 4*time.Second)
			aAngle.AutoReverse = true
			aAngle.RepeatCount = ledgrid.AnimationRepeatForever

			aGroup := ledgrid.NewGroup(aPos, aSize, aColor, aAngle)

			aTimeline := ledgrid.NewTimeline(4 * time.Second)
			aTimeline.RepeatCount = ledgrid.AnimationRepeatForever
			aTimeline.Add(1000*time.Millisecond, ledgrid.NewSuspContAnimation(aColor))
			aTimeline.Add(1500*time.Millisecond, ledgrid.NewSuspContAnimation(aColor))
			aTimeline.Add(2500*time.Millisecond, ledgrid.NewSuspContAnimation(aAngle))
			aTimeline.Add(3000*time.Millisecond, ledgrid.NewSuspContAnimation(aAngle))
			aTimeline.Add(1900*time.Millisecond, ledgrid.NewHideShowAnimation(r))
			aTimeline.Add(2100*time.Millisecond, ledgrid.NewHideShowAnimation(r))

			aGroup.Start()
			aTimeline.Start()
		})

	PathTest = NewLedGridProgram("Path test",
		func(c *ledgrid.Canvas) {
			duration := 4 * time.Second
			pathA := ledgrid.CirclePath
			pathB := ledgrid.CirclePath.NewStart(0.25)

			pos1 := geom.Point{float64(width) / 2.0, 2}
			pos2 := geom.Point{float64(width) - 2, float64(height) / 2.0}
			pos3 := geom.Point{float64(width) / 2.0, float64(height) - 2}
			pos4 := geom.Point{2, float64(height) / 2.0}
			cSize := geom.Point{3.0, 3.0}

			c1 := ledgrid.NewEllipse(pos1, cSize, color.OrangeRed)
			c2 := ledgrid.NewEllipse(pos2, cSize, color.MediumSeaGreen)
			c3 := ledgrid.NewEllipse(pos3, cSize, color.SkyBlue)
			c4 := ledgrid.NewEllipse(pos4, cSize, color.Gold)
			c.Add(c1, c2, c3, c4)

			c1Path := ledgrid.NewPathAnimation(&c1.Pos, pathA, geom.Point{float64(width) / 3.0, float64(height - 4)}, duration)
			c1Path.AutoReverse = true

			c2Path := ledgrid.NewPathAnimation(&c2.Pos, pathB, geom.Point{float64(width - 4), float64(height - 4)}, duration)
			c2Path.AutoReverse = true

			c3Path := ledgrid.NewPathAnimation(&c3.Pos, pathA, geom.Point{-float64(width) / 3.0, -float64(height - 4)}, duration)
			c3Path.AutoReverse = true

			c4Path := ledgrid.NewPathAnimation(&c4.Pos, pathB, geom.Point{-float64(width - 4), -float64(height - 4)}, duration)
			c4Path.AutoReverse = true

			aGrp := ledgrid.NewGroup(c1Path, c3Path, c2Path, c4Path)
			aGrp.RepeatCount = ledgrid.AnimationRepeatForever
			aGrp.Start()
		})

	PolygonPathTest = NewLedGridProgram("Polygon path test",
		func(c *ledgrid.Canvas) {

			cPos := geom.Point{1.5, 1.5}
			cSize := geom.Point{math.Sqrt2, math.Sqrt2}

			polyPath1 := ledgrid.NewPolygonPath(
				geom.Point{1.5, 1.5},
				geom.Point{float64(width) - 1.5, 1.5},
				geom.Point{float64(width) - 1.5, float64(height) - 1.5},
				geom.Point{1.5, float64(height) - 1.5},

				geom.Point{1.5, 2.5},
				geom.Point{float64(width) - 2.5, 2.5},
				geom.Point{float64(width) - 2.5, float64(height) - 2.5},
				geom.Point{2.5, float64(height) - 2.5},

				geom.Point{2.5, 3.5},
				geom.Point{float64(width) - 3.5, 3.5},
				geom.Point{float64(width) - 3.5, float64(height) - 3.5},
				geom.Point{3.5, float64(height) - 3.5},

				geom.Point{3.5, 4.5},
				geom.Point{float64(width) - 4.5, 4.5},
				geom.Point{float64(width) - 4.5, float64(height) - 4.5},
				geom.Point{4.5, float64(height) - 4.5},
			)

			polyPath2 := ledgrid.NewPolygonPath(
				geom.Point{1.5, 1.5},
				geom.Point{4.5, 9.5},
				geom.Point{7.5, 2.5},
				geom.Point{10.5, 8.5},
				geom.Point{13.5, 3.5},
				geom.Point{16.5, 7.5},
				geom.Point{19.5, 4.5},
				geom.Point{22.5, 6.5},
			)

			c1 := ledgrid.NewEllipse(cPos, cSize, color.GreenYellow)
			c1.StrokeWidth = 0.0
			c1.FillColor = color.GreenYellow
			c.Add(c1)

			aPath1 := ledgrid.NewPolyPathAnimation(&c1.Pos, polyPath1, 7*time.Second)
			aPath1.AutoReverse = true

			aPath2 := ledgrid.NewPolyPathAnimation(&c1.Pos, polyPath2, 7*time.Second)
			aPath2.AutoReverse = true

			seq := ledgrid.NewSequence(aPath1, aPath2)
			seq.RepeatCount = ledgrid.AnimationRepeatForever

			seq.Start()
		})

	RandomWalk = NewLedGridProgram("Random walk",
		func(c *ledgrid.Canvas) {
			rect := geom.Rectangle{Min: geom.Point{1.5, 1.5}, Max: geom.Point{float64(width) - 0.5, float64(height) - 0.5}}
			pos1 := geom.Point{1.5, 1.5}
			pos2 := geom.Point{float64(width) - 1.5, float64(height) - 1.5}
			size1 := geom.Point{2.0, 2.0}
			size2 := geom.Point{4.0, 4.0}

			c1 := ledgrid.NewEllipse(pos1, size1, color.SkyBlue)
			c2 := ledgrid.NewEllipse(pos2, size2, color.GreenYellow)
			c.Add(c1, c2)

			aPos1 := ledgrid.NewPositionAnimation(&c1.Pos, geom.Point{}, 1300*time.Millisecond)
			aPos1.Cont = true
			aPos1.Val2 = ledgrid.RandPointTrunc(rect, 1.0)
			aPos1.RepeatCount = ledgrid.AnimationRepeatForever

			aPos2 := ledgrid.NewPositionAnimation(&c2.Pos, geom.Point{}, 901*time.Millisecond)
			aPos2.Cont = true
			aPos2.Val2 = ledgrid.RandPoint(rect)
			aPos2.RepeatCount = ledgrid.AnimationRepeatForever

			aPos1.Start()
			aPos2.Start()
		})

	CirclingCircles = NewLedGridProgram("Circling circles",
		func(c *ledgrid.Canvas) {
			pos1 := geom.Point{1.5, 1.5}
			pos2 := geom.Point{10.5, float64(height) - 1.5}
			pos3 := geom.Point{19.5, 1.5}
			pos4 := geom.Point{28.5, float64(height) - 1.5}
			pos5 := geom.Point{37.5, 1.5}
			cSize := geom.Point{2.0, 2.0}

			c1 := ledgrid.NewEllipse(pos1, cSize, color.OrangeRed)
			c2 := ledgrid.NewEllipse(pos2, cSize, color.MediumSeaGreen)
			c3 := ledgrid.NewEllipse(pos3, cSize, color.SkyBlue)
			c4 := ledgrid.NewEllipse(pos4, cSize, color.Gold)
			c5 := ledgrid.NewEllipse(pos5, cSize, color.YellowGreen)

			stepRD := geom.Point{18.0, 2.0 * (float64(height) - 3.0)}
			stepLU := stepRD.Neg()
			stepRU := geom.Point{18.0, -2.0 * (float64(height) - 3.0)}
			stepLD := stepRU.Neg()

			quartCirc := ledgrid.CirclePath.NewStartLen(0, 0.25)

			c1Path1 := ledgrid.NewPathAnimation(&c1.Pos, quartCirc, stepRD, time.Second)
			c1Path1.Cont = true
			c1Path2 := ledgrid.NewPathAnimation(&c1.Pos, quartCirc, stepRU, time.Second)
			c1Path2.Cont = true
			c1Path3 := ledgrid.NewPathAnimation(&c1.Pos, quartCirc, stepLD, time.Second)
			c1Path3.Cont = true
			c1Path4 := ledgrid.NewPathAnimation(&c1.Pos, quartCirc, stepLU, time.Second)
			c1Path4.Cont = true

			c2Path1 := ledgrid.NewPathAnimation(&c2.Pos, quartCirc, stepLU, time.Second)
			c2Path1.Cont = true
			c2Path2 := ledgrid.NewPathAnimation(&c2.Pos, quartCirc, stepRD, time.Second)
			c2Path2.Cont = true

			c3Path1 := ledgrid.NewPathAnimation(&c3.Pos, quartCirc, stepLD, time.Second)
			c3Path1.Cont = true
			c3Path2 := ledgrid.NewPathAnimation(&c3.Pos, quartCirc, stepRU, time.Second)
			c3Path2.Cont = true

			c4Path1 := ledgrid.NewPathAnimation(&c4.Pos, quartCirc, stepLU, time.Second)
			c4Path1.Cont = true
			c4Path2 := ledgrid.NewPathAnimation(&c4.Pos, quartCirc, stepRD, time.Second)
			c4Path2.Cont = true

			c5Path1 := ledgrid.NewPathAnimation(&c5.Pos, quartCirc, stepLD, time.Second)
			c5Path1.Cont = true
			c5Path2 := ledgrid.NewPathAnimation(&c5.Pos, quartCirc, stepRU, time.Second)
			c5Path2.Cont = true

			aGrp1 := ledgrid.NewGroup(c1Path1, c2Path1)
			aGrp2 := ledgrid.NewGroup(c1Path2, c3Path1)
			aGrp3 := ledgrid.NewGroup(c1Path1, c4Path1)
			aGrp4 := ledgrid.NewGroup(c1Path2, c5Path1)
			aGrp5 := ledgrid.NewGroup(c1Path3, c5Path2)
			aGrp6 := ledgrid.NewGroup(c1Path4, c4Path2)
			aGrp7 := ledgrid.NewGroup(c1Path3, c3Path2)
			aGrp8 := ledgrid.NewGroup(c1Path4, c2Path2)
			aSeq := ledgrid.NewSequence(aGrp1, aGrp2, aGrp3, aGrp4, aGrp5, aGrp6, aGrp7, aGrp8)
			aSeq.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(c1, c2, c3, c4, c5)
			aSeq.Start()
		})

	ChasingCircles = NewLedGridProgram("Chasing circles",
		func(c *ledgrid.Canvas) {
			c1Pos1 := geom.Point{float64(width) - 5.0, float64(height) / 2.0}
			c1Size1 := geom.Point{9.0, 9.0}
			c1Size2 := geom.Point{3.0, 3.0}
			c1PosSize := geom.Point{float64(width - 10), -float64(height) / 2.0}
			c2Pos := geom.Point{5.0, float64(height) / 2.0}
			c2Size1 := geom.Point{5.0, 5.0}
			c2Size2 := geom.Point{3.0, 3.0}
			c2PosSize := geom.Point{-float64(width - 10), float64(height)/2.0 + 2.0}

			aGrp := ledgrid.NewGroup()

			pal := ledgrid.NewGradientPaletteByList("Palette", true,
				color.DeepSkyBlue,
				color.Lime,
				color.Teal,
				color.SkyBlue,
			)

			c1 := ledgrid.NewEllipse(c1Pos1, c1Size1, color.Gold)

			path := ledgrid.CirclePath.NewStart(0.25)

			c1pos := ledgrid.NewPathAnimation(&c1.Pos, path, c1PosSize, 4*time.Second)
			c1pos.Curve = ledgrid.AnimationLinear

			c1size := ledgrid.NewSizeAnimation(&c1.Size, c1Size2, time.Second)
			c1size.AutoReverse = true

			c1bcolor := ledgrid.NewColorAnimation(&c1.StrokeColor, color.OrangeRed, time.Second)
			c1bcolor.AutoReverse = true

			c2 := ledgrid.NewEllipse(c2Pos, c2Size1, color.Lime)

			c2pos := ledgrid.NewPathAnimation(&c2.Pos, path, c2PosSize, 4*time.Second)
			c2pos.Curve = ledgrid.AnimationLinear

			c2size := ledgrid.NewSizeAnimation(&c2.Size, c2Size2, time.Second)
			c2size.AutoReverse = true

			c2color := ledgrid.NewPaletteAnimation(&c2.StrokeColor, pal, 2*time.Second)
			c2color.Curve = ledgrid.AnimationLinear

			aGrp.Add(c1pos, c1size, c1bcolor, c2pos, c2size, c2color)
			aGrp.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(c1, c2)
			aGrp.Start()
		})

	CircleAnimation = NewLedGridProgram("Circle animation",
		func(c *ledgrid.Canvas) {
			c1Pos1 := geom.Point{2.0, float64(height) / 2.0}
			c1Pos3 := geom.Point{float64(width) - 2.0, float64(height) / 2.0}

			c1Size1 := geom.Point{3.0, 3.0}
			c1Size2 := geom.Point{9.0, 9.0}

			c1 := ledgrid.NewEllipse(c1Pos1, c1Size1, color.OrangeRed)

			c1pos := ledgrid.NewPositionAnimation(&c1.Pos, c1Pos3, 2*time.Second)
			c1pos.AutoReverse = true
			c1pos.RepeatCount = ledgrid.AnimationRepeatForever

			c1radius := ledgrid.NewSizeAnimation(&c1.Size, c1Size2, time.Second)
			c1radius.AutoReverse = true
			c1radius.RepeatCount = ledgrid.AnimationRepeatForever

			c1color := ledgrid.NewColorAnimation(&c1.StrokeColor, color.Gold, 4*time.Second)
			c1color.AutoReverse = true
			c1color.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(c1)

			c1pos.Start()
			c1radius.Start()
			c1color.Start()
		})

	PushingRectangles = NewLedGridProgram("Pushing rectangles",
		func(c *ledgrid.Canvas) {
			r1Pos1 := geom.Point{1.0, float64(height) / 2.0}
			r1Pos2 := geom.Point{0.5 + float64(width-3)/2.0, float64(height) / 2.0}

			r2Pos1 := geom.Point{float64(width - 1), float64(height) / 2.0}
			r2Pos2 := geom.Point{float64(width) - 0.5 - float64(width-3)/2.0, float64(height) / 2.0}

			rSize1 := geom.Point{float64(width - 3), 1.0}
			rSize2 := geom.Point{1.0, float64(height - 1)}

			duration := 2 * time.Second

			r1 := ledgrid.NewRectangle(r1Pos1, rSize2, color.Crimson)

			a1Pos := ledgrid.NewPositionAnimation(&r1.Pos, r1Pos2, duration)
			a1Pos.AutoReverse = true
			a1Pos.RepeatCount = ledgrid.AnimationRepeatForever

			a1Size := ledgrid.NewSizeAnimation(&r1.Size, rSize1, duration)
			a1Size.AutoReverse = true
			a1Size.RepeatCount = ledgrid.AnimationRepeatForever

			a1Color := ledgrid.NewColorAnimation(&r1.StrokeColor, color.GreenYellow, duration)
			a1Color.AutoReverse = true
			a1Color.RepeatCount = ledgrid.AnimationRepeatForever

			r2 := ledgrid.NewRectangle(r2Pos2, rSize1, color.SkyBlue)

			a2Pos := ledgrid.NewPositionAnimation(&r2.Pos, r2Pos1, duration)
			a2Pos.AutoReverse = true
			a2Pos.RepeatCount = ledgrid.AnimationRepeatForever

			a2Size := ledgrid.NewSizeAnimation(&r2.Size, rSize2, duration)
			a2Size.AutoReverse = true
			a2Size.RepeatCount = ledgrid.AnimationRepeatForever

			a2Color := ledgrid.NewColorAnimation(&r2.StrokeColor, color.Crimson, duration)
			a2Color.AutoReverse = true
			a2Color.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(r1, r2)
			a1Pos.Start()
			a1Size.Start()
			a1Color.Start()
			a2Pos.Start()
			a2Size.Start()
			a2Color.Start()
		})

	RegularPolygonTest = NewLedGridProgram("Regular Polygon test",
		func(c *ledgrid.Canvas) {
			posList := []geom.Point{
				geom.Point{-6.0, float64(height) / 2.0},
				geom.Point{float64(width) + 5.0, float64(height) / 2.0},
			}
			posCenter := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			smallSize := geom.Point{9.0, 9.0}
			largeSize := geom.Point{80.0, 80.0}

			polyList := make([]*ledgrid.RegularPolygon, 9)

			aSeq := ledgrid.NewSequence()
			for n := 3; n <= 6; n++ {
				col := color.RandColor()
				polyList[n] = ledgrid.NewRegularPolygon(n, posList[n%2], smallSize, col)
				c.Add(polyList[n])
				dur := 2*time.Second + rand.N(time.Second)
				sign := []float64{+1.0, -1.0}[n%2]
				angle := sign * 2 * math.Pi
				animPos := ledgrid.NewPositionAnimation(&polyList[n].Pos, posCenter, dur)
				animAngle := ledgrid.NewFloatAnimation(&polyList[n].Angle, angle, dur)
				animSize := ledgrid.NewSizeAnimation(&polyList[n].Size, largeSize, 4*time.Second)
				animSize.Cont = true
				animFade := ledgrid.NewColorAnimation(&polyList[n].StrokeColor, color.Black, 4*time.Second)
				animFade.Cont = true

				aGrpIn := ledgrid.NewGroup(animPos, animAngle)
				aGrpOut := ledgrid.NewGroup(animSize, animFade)
				aObjSeq := ledgrid.NewSequence(aGrpIn, aGrpOut)
				aSeq.Add(aObjSeq)
			}
			aSeq.RepeatCount = ledgrid.AnimationRepeatForever
			aSeq.Start()
		})

	MovingText = NewLedGridProgram("Moving text",
		func(c *ledgrid.Canvas) {
			t1 := ledgrid.NewText(geom.Point{0, float64(height) / 2.0}, "Stefan", color.LightSeaGreen)
			t1.SetAlign(ledgrid.AlignLeft | ledgrid.AlignMiddle)
			t2 := ledgrid.NewText(geom.Point{float64(width), float64(height) / 2.0}, "Beni", color.YellowGreen)
			t2.SetAlign(ledgrid.AlignRight | ledgrid.AlignMiddle)
			t3 := ledgrid.NewText(geom.Point{float64(width) + 60.0, float64(height) / 2.0}, "wohnen im Lochbach", color.OrangeRed)
			t3.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)

			c.Add(t1, t2, t3)

			aAngle1 := ledgrid.NewFloatAnimation(&t1.Angle, -2*math.Pi, 7*time.Second)
			aAngle1.Curve = ledgrid.AnimationLinear
			aAngle1.RepeatCount = ledgrid.AnimationRepeatForever

			aAngle2 := ledgrid.NewFloatAnimation(&t2.Angle, -2*math.Pi, 8*time.Second)
			aAngle2.Curve = ledgrid.AnimationLinear
			aAngle2.RepeatCount = ledgrid.AnimationRepeatForever

			aPos := ledgrid.NewPositionAnimation(&t3.Pos, geom.Point{-100, float64(height) / 2.0}, 6*time.Second)
			aPos.Curve = ledgrid.AnimationEaseInOut

			aTimeline := ledgrid.NewTimeline(15 * time.Second)
			aTimeline.Add(10*time.Second, aPos)
			aTimeline.RepeatCount = ledgrid.AnimationRepeatForever

			aAngle1.Start()
			aAngle2.Start()
			aTimeline.Start()
		})

	BitmapText = NewLedGridProgram("Bitmap text",
		func(c *ledgrid.Canvas) {
			pos1 := fixed.P(55, height/2)
			pos2 := fixed.P(-15, height/2)
			color1 := color.Aquamarine
			color2 := color.BlueViolet

			txt1 := ledgrid.NewFixedText(pos1, color1, "STEFAN")
			c.Add(txt1)

			aPos := ledgrid.NewFixedPosAnimation(&txt1.Pos, pos2, 10*time.Second)
			aPos.AutoReverse = true
			aPos.RepeatCount = ledgrid.AnimationRepeatForever

			aColor := ledgrid.NewColorAnimation(&txt1.Color, color2, 2*time.Second)
			aColor.AutoReverse = true
			aColor.RepeatCount = ledgrid.AnimationRepeatForever

			aPos.Start()
			aColor.Start()
		})

	SingleImageAlign = NewLedGridProgram("Align this lonely image!",
		func(c *ledgrid.Canvas) {
			imgPos := geom.Point{float64(width / 2), float64(height / 2)}
			imgSize := geom.Point{30, 30}
			img := ledgrid.NewImage(imgPos, "images/raster.png")
			img.Size = imgSize
			img.SetAlign(ledgrid.AlignRight | ledgrid.AlignBottom)
			c.Add(img)

			aAngle := ledgrid.NewFloatAnimation(&img.Angle, 2*math.Pi, 4*time.Second)
			aAngle.Curve = ledgrid.AnimationLinear
			aAngle.AutoReverse = true
			aAngle.RepeatCount = ledgrid.AnimationRepeatForever
			aAngle.Start()
		})

	SlideTheShow = NewLedGridProgram("Slide-the-Show",
		func(c *ledgrid.Canvas) {
			pos2 := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			files := []string{
				"images/raster.png",
				"images/square1.png",
				"images/square2.png",
			}
			aTimeline := ledgrid.NewTimeline(time.Duration(2*len(files)) * time.Second)
			dstSize := geom.NewPointIMG(c.Bounds().Size())
			dstRatio := dstSize.X / dstSize.Y
			for i, fileName := range files {
				img := ledgrid.NewImage(pos2, fileName)
				img.Hide()
				srcRatio := img.Size.X / img.Size.Y
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
			aTimeline.Start()
		})

	CameraTest = NewLedGridProgram("Camera test",
		func(c *ledgrid.Canvas) {
			pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			size := geom.Point{float64(width), float64(height)}

			cam := NewCamera(pos, size)
			c.Add(cam)
			cam.Start()
			mask := cam.Mask

			go func() {
				var src, dst geom.Point
				effectList := []DissolveType{
					{Top2Bottom, Forward, ExitOver},
					{Top2Bottom, Backward, ExitOver},
					{Top2Bottom, Forward, ExitAway},
					{Top2Bottom, Backward, ExitAway},
					// {Out2Inside, Forward, ExitAway},
					// {Out2Inside, Backward, ExitAway},
					// {Left2Right, Backward, ExitAway},
					// {Left2Right, Forward, ExitAway},
					// {Bottom2Top, Forward, ExitAway},
					// {Right2Left, Forward, ExitAway},
					// {In2Outside, Forward, ExitAway},
				}

				time.Sleep(1 * time.Second)
				for _, effect := range effectList {
					for _, pts := range Dissolver(effect) {
						for _, pp := range pts {
							p0, p1 := pp.src, pp.dst
							src = geom.NewPointIMG(p0)
							dst = geom.NewPointIMG(p1)
							pixAway := ledgrid.NewDot(src, color.FireBrick.Alpha(0.0))
							c.Add(pixAway)

							aMask := ledgrid.NewTask(func() {
								idx := mask.PixOffset(p0.X, p0.Y)
								mask.Pix[idx] = 0x00
							})

							aDur := rand.N(time.Second)
							aFadeIn := ledgrid.NewFadeAnimation(&pixAway.Color.A, ledgrid.FadeIn, aDur)
							aFadeIn.Curve = ledgrid.AnimationLazeIn
							aDur = time.Second + rand.N(time.Second)
							aFadeOut := ledgrid.NewFadeAnimation(&pixAway.Color.A, ledgrid.FadeOut, aDur)
							aFadeOut.Curve = ledgrid.AnimationEaseIn
							aFadeOut.Cont = true
							aColor2 := ledgrid.NewColorAnimation(&pixAway.Color, color.DarkRed, aDur)
							aColor2.Curve = ledgrid.AnimationEaseIn
							aColor2.Cont = true
							aPos := ledgrid.NewPositionAnimation(&pixAway.Pos, dst, aDur)
							aPos.Curve = ledgrid.AnimationEaseIn
							aGrp2 := ledgrid.NewGroup(aMask, aColor2, aFadeOut, aPos)
							aSeq := ledgrid.NewSequence(aFadeIn, aGrp2)
							aSeq.Start()
							// time.Sleep(30 * time.Millisecond)
						}
						time.Sleep(500 * time.Millisecond)
					}
					time.Sleep(3 * time.Second)
					for i := range mask.Pix {
						mask.Pix[i] = 0xff
					}
				}
			}()
		})

	BlinkenAnimation = NewLedGridProgram("Blinken animation",
		func(c *ledgrid.Canvas) {
			posFlame1 := geom.Point{4.5, float64(height - 1)}
			posFlame2 := geom.Point{float64(width) - 4.5, float64(height - 1)}
			pos1Mario := geom.Point{10.0, float64(height - 1)}
			pos2Mario := geom.Point{float64(width - 11), float64(height - 1)}

			bmlFlame := ledgrid.ReadBlinkenFile("blinken/flameNew.bml")
			bmlFlame.SetAllDuration(32)

			flame1 := ledgrid.NewImageList(posFlame1)
			flame1.SetAlign(ledgrid.AlignCenter | ledgrid.AlignBottom)
			flame1.AddBlinkenLight(bmlFlame)
			flame1.RepeatCount = ledgrid.AnimationRepeatForever

			bmlFlame.SetAllDuration(43)

			flame2 := ledgrid.NewImageList(posFlame2)
			flame2.SetAlign(ledgrid.AlignCenter | ledgrid.AlignBottom)
			flame2.AddBlinkenLight(bmlFlame)
			flame2.RepeatCount = ledgrid.AnimationRepeatForever

			bmlMario := ledgrid.ReadBlinkenFile("blinken/marioWalkRight.bml")

			mario := ledgrid.NewImageList(pos1Mario)
			mario.SetAlign(ledgrid.AlignCenter | ledgrid.AlignBottom)
			mario.AddBlinkenLight(bmlMario)
			mario.RepeatCount = ledgrid.AnimationRepeatForever
			// mario.Size = geom.Point{20.0, 20.0}

			aPos := ledgrid.NewPositionAnimation(&mario.Pos, pos2Mario, 3*time.Second)
			aPos.Curve = ledgrid.AnimationLinear
			aPos.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(flame1, flame2, mario)

			aGrp := ledgrid.NewGroup(flame1, flame2, mario, aPos)
			aGrp.Start()
		})

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
						c.Add(pix)
						aPos := ledgrid.NewPositionAnimation(&pix.Pos, dest, time.Second+rand.N(time.Second))
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
						c.Add(pix)
						aPos := ledgrid.NewIntegerPosAnimation(&pix.Pos, dest.Int(), time.Second+rand.N(time.Second))
						aPos.AutoReverse = true
						grp.Add(aPos)
					}
				}
				aSeq.Put(grp)
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

					pt := image.Point{x, y}
					// pix := ledgrid.NewPixel(pt, col)

					pos := geom.NewPointIMG(pt)
					pix := ledgrid.NewDot(pos, col)

					c.Add(pix)

					dur := time.Second + time.Duration(x)*time.Millisecond
					aAlpha := ledgrid.NewFadeAnimation(&pix.Color.A, 196, dur)
					aAlpha.AutoReverse = true
					aAlpha.RepeatCount = ledgrid.AnimationRepeatForever
					aAlpha.Start()

					aColor := ledgrid.NewColorAnimation(&pix.Color, col, 1*time.Second)
					aColor.Cont = true
					aGrpGrey.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.MediumPurple.Interpolate(color.Fuchsia, t), 5*time.Second)
					aColor.Cont = true
					aGrpPurple.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.Gold.Interpolate(color.Khaki, t), 5*time.Second)
					aColor.Cont = true
					aGrpYellow.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.GreenYellow.Interpolate(color.LightSeaGreen, t), 5*time.Second)
					aColor.Cont = true
					aGrpGreen.Add(aColor)
				}
			}

			txt1 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.GreenYellow.Alpha(0.0), "LORENZ")
			txt1.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt1 := ledgrid.NewFadeAnimation(&txt1.Color.A, ledgrid.FadeIn, 2*time.Second)
			aTxt1.AutoReverse = true
			txt2 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.DarkViolet.Alpha(0.0), "SIMON")
			txt2.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt2 := ledgrid.NewFadeAnimation(&txt2.Color.A, ledgrid.FadeIn, 2*time.Second)
			aTxt2.AutoReverse = true
			txt3 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.OrangeRed.Alpha(0.0), "REBEKKA")
			txt3.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			aTxt3 := ledgrid.NewFadeAnimation(&txt3.Color.A, ledgrid.FadeIn, 2*time.Second)
			aTxt3.AutoReverse = true
			c.Add(txt1, txt2, txt3)

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

	TestShaderFunc = func(x, y, t float64) float64 {
		t = t/4.0 + x
		_, f := math.Modf(math.Abs(t))
		return f
	}

	f1 = func(x, y, t, p1 float64) float64 {
		return math.Sin(x*p1 + t)
	}

	f2 = func(x, y, t, p1, p2, p3 float64) float64 {
		return math.Sin(p1*(x*math.Sin(t/p2)+y*math.Cos(t/p3)) + t)
	}

	f3 = func(x, y, t, p1, p2 float64) float64 {
		cx := 0.125*x + 0.5*math.Sin(t/p1)
		cy := 0.125*y + 0.5*math.Cos(t/p2)
		return math.Sin(math.Sqrt(100.0*(cx*cx+cy*cy)+1.0) + t)
	}

	PlasmaShaderFunc = func(x, y, t float64) float64 {
		v1 := f1(x, y, t, 1.2)
		v2 := f2(x, y, t, 1.6, 3.0, 1.5)
		v3 := f3(x, y, t, 5.0, 5.0)
		v := (v1+v2+v3)/6.0 + 0.5
		return v
	}

	ShowTheShader = NewLedGridProgram("Show the shader!",
		func(c *ledgrid.Canvas) {
			var xMin, yMax float64
			var txt *ledgrid.FixedText
			var palId int

			pal := ledgrid.PaletteMap["Hipster"]
			fader := ledgrid.NewPaletteFader(pal)
			aPal := ledgrid.NewPaletteFadeAnimation(fader, pal, 2*time.Second)
			aPal.ValFunc = func() ledgrid.ColorSource {
				name := ledgrid.PaletteNames[palId]
				palId = (palId + 1) % len(ledgrid.PaletteNames)
				log.Printf(">>> Switch palette, new name: '%s'", name)
				txt.SetText(name)
				return ledgrid.PaletteMap[name]
			}

			aPalTl := ledgrid.NewTimeline(10 * time.Second)
			aPalTl.Add(7*time.Second, aPal)
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
					anim := ledgrid.NewShaderAnimation(&pix.Color, fader, x, y, PlasmaShaderFunc)
					aGrp.Add(anim)
					x += dPix
				}
				y -= dPix
			}
			txt = ledgrid.NewFixedText(fixed.P(width/2, height/2), color.YellowGreen, "Hipster")
			txt.SetAlign(ledgrid.AlignCenter | ledgrid.AlignMiddle)
			c.Add(txt)
			aPalTl.Start()
			aGrp.Start()
		})

	ColorFields = NewLedGridProgram("Fields of named colors",
		func(c *ledgrid.Canvas) {
			var input int
			var colGrp color.ColorGroup

			cs := NewColorSampler(color.Purples)
			c.Add(cs)

			for {
				fmt.Printf("Enter a number in 0..%d (or 99 to quit): ", color.NumColorGroups-1)
				fmt.Scanf("%d\n", &input)
				if input == 99 {
					return
				}
				colGrp = color.ColorGroup(input)
				if colGrp >= color.NumColorGroups {
					continue
				}
				fmt.Printf("Selected color group: %v\n", colGrp)
				cs.colGrp = colGrp
			}
		})
)

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

//----------------------------------------------------------------------------

type BouncingEllipse struct {
	ledgrid.Ellipse
	Vel, Acc geom.Point
	Field    geom.Rectangle
}

func NewBouncingEllipse(pos, size geom.Point, col color.LedColor) *BouncingEllipse {
	b := &BouncingEllipse{}
	b.Ellipse = *ledgrid.NewEllipse(pos, size, col)
	b.Vel = geom.Point{}
	b.Acc = geom.Point{}
	return b
}

func (b *BouncingEllipse) Update(pit time.Time) bool {
	deltaVel := b.Acc.Mul(0.3)
	b.Vel = b.Vel.Add(deltaVel)
	b.Pos = b.Pos.Add(b.Vel)
	if b.Pos.X < b.Field.Min.X || b.Pos.X >= b.Field.Max.X {
		b.Vel.X = -b.Vel.X
	}
	if b.Pos.Y < b.Field.Min.Y || b.Pos.Y >= b.Field.Max.Y {
		b.Vel.Y = -b.Vel.Y
	}
	return true
}

func (b *BouncingEllipse) Duration() time.Duration {
	return time.Duration(0)
}
func (b *BouncingEllipse) SetDuration(dur time.Duration) {}
func (b *BouncingEllipse) Start()                        {}
func (b *BouncingEllipse) Suspend()                      {}
func (b *BouncingEllipse) Continue()                     {}
func (b *BouncingEllipse) IsRunning() bool {
	return true
}

func BounceAround(c *ledgrid.Canvas) {
	pos1 := geom.Point{2.0, 2.0}
	pos2 := geom.Point{37.0, 7.0}
	size := geom.Point{4.0, 4.0}
	vel1 := geom.Point{0.15, 0.075}
	vel2 := geom.Point{-0.35, -0.25}

	obj1 := NewBouncingEllipse(pos1, size, color.GreenYellow)
	obj1.Vel = vel1
	obj1.Field = geom.NewRectangleIMG(c.Rect)
	obj2 := NewBouncingEllipse(pos2, size, color.LightSeaGreen)
	obj2.Vel = vel2
	obj2.Field = geom.NewRectangleIMG(c.Rect)

	c.Add(obj1, obj2)
	ledgrid.AnimCtrl.Add(obj1, obj2)
}

//----------------------------------------------------------------------------

func SignalHandler() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

//----------------------------------------------------------------------------

var (
	programList = []LedGridProgram{
		FarewellGery,
		GroupTest,
		SequenceTest,
		TimelineTest,
		StopContShowHideTest,
		PathTest,
		PolygonPathTest,
		RandomWalk,
		CirclingCircles,
		ChasingCircles,
		PushingRectangles,
		RegularPolygonTest,
		GlowingPixels,
		MovingPixels,
		// ImageFilterTest,
		CameraTest,
		BlinkenAnimation,
		MovingText,
		BitmapText,
		SlideTheShow,
		ShowTheShader,
		ColorFields,
		SingleImageAlign,
	}
)

func main() {
	var host string
	var port uint
	var input string
	var ch byte
	var progId int
	var runInteractive bool
	var useCustomLayout bool
	var progList string
	var gR, gG, gB float64
	var customConf conf.ModuleConfig = conf.ChessBoard

	for i, prog := range programList {
		id := 'a' + i
		progList += fmt.Sprintf("\n%c - %s", id, prog.Name())
	}

	flag.IntVar(&width, "width", 40, "Width of LedGrid")
	flag.IntVar(&height, "height", 10, "Height of LedGrid")
	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.StringVar(&input, "prog", input, "Play one single program"+progList)
	flag.BoolVar(&useCustomLayout, "custom", false, "Use custom module configuration")
	flag.Parse()

	if len(input) > 0 {
		runInteractive = false
		ch = input[0]
	} else {
		runInteractive = true
	}

	if useCustomLayout {
		ledGrid = ledgrid.NewLedGrid(host, port, customConf)
	} else {
		ledGrid = ledgrid.NewLedGridBySize(host, port, image.Pt(width, height))
	}
	gR, gG, gB = ledGrid.Client.Gamma()

	gridSize = ledGrid.Rect.Size()
	width = gridSize.X
	height = gridSize.Y

	canvas = ledgrid.NewCanvas(gridSize)
	ledgrid.NewAnimationController(canvas, ledGrid)

	// initSpiralMap(CW)
	// initSpiralMap(CCW)

	if runInteractive {
		progId = -1
		for {
			fmt.Printf("---------------------------------------\n")
			fmt.Printf("  Program\n")
			fmt.Printf("---------------------------------------\n")
			for i, prog := range programList {
				if ch >= 'a' && ch <= 'z' && i == progId {
					fmt.Printf("> ")
				} else {
					fmt.Printf("  ")
				}
				fmt.Printf("[%c] %s\n", 'a'+i, prog.Name())
			}
			fmt.Printf("---------------------------------------\n")
			fmt.Printf("  Gamma values: %.1f, %.1f, %.1f\n", gR, gG, gB)
			fmt.Printf("   +/-: increase/decreases by 0.1\n")
			fmt.Printf("---------------------------------------\n")

			fmt.Printf("Enter a character (or '0' for quit): ")

			fmt.Scanf("%s\n", &input)
			ch = input[0]
			if ch == '0' {
				break
			}

			if ch >= 'a' && ch <= 'z' {
				progId = int(ch - 'a')
				if progId < 0 || progId >= len(programList) {
					continue
				}
				// ledgrid.AnimCtrl.Stop()
				fmt.Printf("Program statistics:\n")
				fmt.Printf("  animation: %v\n", ledgrid.AnimCtrl.Watch())
				fmt.Printf("  painting : %v\n", canvas.Watch())
				fmt.Printf("  sending  : %v\n", ledGrid.Client.Watch())
				ledgrid.AnimCtrl.Purge()
				// ledgrid.AnimCtrl.Continue()
				canvas.Purge()
				ledgrid.AnimCtrl.Watch().Reset()
				canvas.Watch().Reset()
				ledGrid.Client.Watch().Reset()
				programList[progId].Run(canvas)
			}
			if ch == 'S' {
				ledgrid.AnimCtrl.Save("gobs/program01.gob")
			}
			if ch == 'L' {
				ledgrid.AnimCtrl.Suspend()
				ledgrid.AnimCtrl.Purge()
				ledgrid.AnimCtrl.Watch().Reset()
				canvas.Purge()
				canvas.Watch().Reset()
				time.Sleep(60 * time.Millisecond)
				ledgrid.AnimCtrl.Load("gobs/program01.gob")
				ledgrid.AnimCtrl.Continue()
				// fmt.Printf("canvas  >>> %+v\n", canvas)
				// for i, obj := range canvas.ObjList {
				i := 0
				for ele := canvas.ObjList.Front(); ele != nil; ele = ele.Next() {
					obj := ele.Value.(ledgrid.CanvasObject)
					if obj == nil {
						continue
					}
					fmt.Printf(">>> obj[%d] : %[2]T %+[2]v\n", i, obj)
					i++
				}
				// fmt.Printf("animCtrl>>> %+v\n", ledgrid.AnimCtrl)
				for i, anim := range ledgrid.AnimCtrl.AnimList {
					if anim == nil {
						continue
					}
					fmt.Printf(">>> anim[%d]: %[2]T %+[2]v\n", i, anim)
				}
			}
			if ch == '+' {
				gR += 0.1
				gG += 0.1
				gB += 0.1
				ledGrid.Client.SetGamma(gR, gG, gB)
			}
			if ch == '-' {
				if gR > 0.1 {
					gR -= 0.1
					gG -= 0.1
					gB -= 0.1
					ledGrid.Client.SetGamma(gR, gG, gB)
				}
			}
		}
	} else {
		if ch >= 'a' && ch <= 'z' {
			progId = int(ch - 'a')
			if progId >= 0 && progId < len(programList) {
				programList[progId].Run(canvas)
			}
		}
		fmt.Printf("Quit by Ctrl-C\n")
		SignalHandler()
	}

	ledgrid.AnimCtrl.Suspend()
	ledGrid.Clear(color.Black)
	ledGrid.Close()

	fmt.Printf("Program statistics:\n")
	fmt.Printf("  animation: %v\n", ledgrid.AnimCtrl.Watch())
	fmt.Printf("  painting : %v\n", canvas.Watch())
	fmt.Printf("  sending  : %v\n", ledGrid.Client.Watch())
}
