package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand/v2"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
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

// ---------------------------------------------------------------------------

type ScanDir int
type ScanLine int
type ExitDir int

const (
	// Die Konstanten fuer ScanDir
	Top2Bottom ScanDir = iota
	Bottom2Top
	Left2Right
	Right2Left
	Out2Inside
	In2Outside
)

const (
	// Die Konstanten fuer ScanLine
	Forward, CW = iota, iota
	Backward, CCW
	Alternate = iota
)

const (
	// Die Konstanten fuer ExitDir
	ExitAway ExitDir = iota
	ExitOver
)

type PointPair struct {
	src, dst image.Point
}

var (
	spiralMap = [][]PointPair{
        []PointPair{},
        []PointPair{},
    }
)

func initSpiralMap(dir ScanLine) {
	var dx, dy, m12, m21, w, h int
	var pt0, pt1 image.Point
	var segList []int

    spiralMap[dir] = make([]PointPair, 0)

	switch dir {
	case CW:
		dx, dy = 1, 0
		m12, m21 = -1, 1
		w, h = width, height-1
		pt0 = image.Pt(-1, 0)
	case CCW:
		dx, dy = 0, 1
		m12, m21 = 1, -1
		w, h = width-1, height
		pt0 = image.Pt(0, -1)
	}
	n := 0
LOOP:
	for {
		switch dir {
		case CW:
			segList = []int{w, h, w - 1, h - 1}
		case CCW:
			segList = []int{h, w, h - 1, w - 1}
		}
		for _, seg := range segList {
			for range seg {
				pt0 = pt0.Add(image.Pt(dx, dy))
				pt1 = pt0.Add(image.Pt(-m12*dy*height, -m21*dx*height))
				spiralMap[dir] = append(spiralMap[dir], PointPair{pt0, pt1})
			}
			n += seg
			if n >= width*height {
				break LOOP
			}
			dx, dy = m12*dy, m21*dx
		}
		w, h = w-2, h-2
	}
}

func StartEndPoints(idx int, scanDir ScanDir, scanLine ScanLine,
	exitDir ExitDir) (pts PointPair) {

	switch scanDir {
	case Top2Bottom:
		pts.src.Y = idx / width
		if scanLine == Forward || (scanLine == Alternate && pts.src.Y%2 == 0) {
			pts.src.X = idx % width
		} else {
			pts.src.X = (width - 1) - idx%width
		}
		pts.dst = pts.src
		if exitDir == ExitAway {
			pts.dst = pts.dst.Add(image.Pt(0, -(height + 2)))
		} else {
			pts.dst = pts.dst.Add(image.Pt(0, (height + 2)))
		}
	case Bottom2Top:
		pts.src.Y = (height - 1) - idx/width
		if scanLine == Forward || (scanLine == Alternate && pts.src.Y%2 == 0) {
			pts.src.X = idx % width
		} else {
			pts.src.X = (width - 1) - idx%width
		}
		pts.dst = pts.src
		if exitDir == ExitAway {
			pts.dst = pts.dst.Add(image.Pt(0, (height + 2)))
		} else {
			pts.dst = pts.dst.Add(image.Pt(0, -(height + 2)))
		}
	case Left2Right:
		pts.src.X = idx / height
		if scanLine == Forward || (scanLine == Alternate && pts.src.X%2 == 0) {
			pts.src.Y = idx % height
		} else {
			pts.src.Y = (height - 1) - idx%height
		}
		pts.dst = pts.src
		if exitDir == ExitAway {
			pts.dst = pts.dst.Add(image.Pt(-(width + 2), 0))
		} else {
			pts.dst = pts.dst.Add(image.Pt((width + 2), 0))
		}
	case Right2Left:
		pts.src.X = (width - 1) - idx/height
		if scanLine == Forward || (scanLine == Alternate && pts.src.X%2 == 0) {
			pts.src.Y = idx % height
		} else {
			pts.src.Y = (height - 1) - idx%height
		}
		pts.dst = pts.src
		if exitDir == ExitAway {
			pts.dst = pts.dst.Add(image.Pt((width + 2), 0))
		} else {
			pts.dst = pts.dst.Add(image.Pt(-(width + 2), 0))
		}
	case Out2Inside:
		pts = spiralMap[scanLine][idx]
	case In2Outside:
		i := len(spiralMap[0]) - 1 - idx
		pts = spiralMap[(scanLine+1)%2][i]
	}

	return
}

// ---------------------------------------------------------------------------

type ShiftFilter struct {
	ledgrid.FilterBase
	Dx, Dy int
	W, H   int
}

func NewShiftFilter(img draw.Image) *ShiftFilter {
	f := &ShiftFilter{}
	f.FilterBase.Extend(img, f)
	f.W = img.Bounds().Dx()
	f.H = img.Bounds().Dy()
	return f
}

func (f *ShiftFilter) FF(x, y int) (int, int) {
	x = (x + f.Dx) % f.W
	y = (y + f.Dy) % f.H
	return x, y
}

// ---------------------------------------------------------------------------

type ShuffleFilter struct {
	ledgrid.FilterBase
	stride, lastPos int
	pt              []image.Point
}

func NewShuffleFilter(img draw.Image) *ShuffleFilter {
	f := &ShuffleFilter{}
	f.FilterBase.Extend(img, f)
	size := img.Bounds().Size()
	f.stride = size.X
	f.lastPos = 0
	f.pt = make([]image.Point, size.X*size.Y)
	for i := range f.pt {
		x, y := i%f.stride, i/f.stride
		f.pt[i] = image.Point{x, y}
	}
	return f
}

func (f *ShuffleFilter) FF(x, y int) (int, int) {
	pt := f.pt[y*f.stride+x]
	return pt.X, pt.Y
}

func (f *ShuffleFilter) Reset() {
	for i := range f.pt {
		x, y := i%f.stride, i/f.stride
		f.pt[i] = image.Point{x, y}
	}
}

func (f *ShuffleFilter) ShuffleOne() {
	i, j := rand.IntN(len(f.pt)), rand.IntN(len(f.pt))
	f.pt[i], f.pt[j] = f.pt[j], f.pt[i]
}

func (f *ShuffleFilter) Shuffle() {
	rand.Shuffle(len(f.pt), func(i, j int) {
		f.pt[i], f.pt[j] = f.pt[j], f.pt[i]
	})
}

func (f *ShuffleFilter) OrderOne() {
	var i, j int
	var pt image.Point
	var isOrdered bool

	isOrdered = true
	i = f.lastPos
	for range len(f.pt) {
		x, y := i%f.stride, i/f.stride
		pt = f.pt[i]
		if pt.X != x || pt.Y != y {
			isOrdered = false
			break
		}
		i++
		if i >= len(f.pt) {
			i %= len(f.pt)
		}
	}
	if isOrdered {
		return
	}
	j = pt.Y*f.stride + pt.X
	f.pt[i], f.pt[j] = f.pt[j], f.pt[i]
	f.lastPos = j
}

func (f *ShuffleFilter) SwapRows(i, j int) {
	i = i * f.stride
    j = j * f.stride

    for idx := range f.stride {
        f.pt[i+idx], f.pt[j+idx] = f.pt[j+idx], f.pt[i+idx]
    }
}

func (f *ShuffleFilter) SwapCols(i, j int) {
    for range len(f.pt)/f.stride {
        f.pt[i], f.pt[j] = f.pt[j], f.pt[i]
        i, j = i+f.stride, j+f.stride
    }
}




// ---------------------------------------------------------------------------

var (
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
			aColor := ledgrid.NewColorAnimation(&r.BorderColor, rColor2, 2*time.Second)
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
			aColor1 := ledgrid.NewColorAnimation(&r.BorderColor, color.OrangeRed, time.Second/2)
			aColor1.AutoReverse = true
			aColor2 := ledgrid.NewColorAnimation(&r.BorderColor, color.Crimson, time.Second/2)
			aColor2.AutoReverse = true
			aColor3 := ledgrid.NewColorAnimation(&r.BorderColor, color.Coral, time.Second/2)
			aColor3.AutoReverse = true
			aColor4 := ledgrid.NewColorAnimation(&r.BorderColor, color.FireBrick, time.Second/2)
			aSize2 := ledgrid.NewSizeAnimation(&r.Size, rSize1, time.Second)
			aSize2.Cont = true
			aColor5 := ledgrid.NewColorAnimation(&r.BorderColor, color.SkyBlue, time.Second)
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

			aColor1 := ledgrid.NewColorAnimation(&r1.BorderColor, color.OrangeRed, 200*time.Millisecond)
			aColor1.AutoReverse = true
			aColor1.RepeatCount = 3
			aColor2 := ledgrid.NewColorAnimation(&r1.BorderColor, color.Purple, 500*time.Millisecond)
			aColor3 := ledgrid.NewColorAnimation(&r1.BorderColor, color.GreenYellow, 500*time.Millisecond)
			aColor3.Cont = true

			aPos1 := ledgrid.NewPositionAnimation(&r1.Pos, r2.Pos.SubXY(r2Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos1.AutoReverse = true

			aAngle3 := ledgrid.NewFloatAnimation(&r3.Angle, -math.Pi, time.Second)
			aAngle4 := ledgrid.NewFloatAnimation(&r3.Angle, 0.0, time.Second)
			aAngle4.Cont = true

			aColor4 := ledgrid.NewColorAnimation(&r3.BorderColor, color.DarkOrange, 200*time.Millisecond)
			aColor4.AutoReverse = true
			aColor4.RepeatCount = 3
			aColor5 := ledgrid.NewColorAnimation(&r3.BorderColor, color.Purple, 500*time.Millisecond)
			aColor6 := ledgrid.NewColorAnimation(&r3.BorderColor, color.SkyBlue, 500*time.Millisecond)
			aColor6.Cont = true

			aPos2 := ledgrid.NewPositionAnimation(&r3.Pos, r4.Pos.AddXY(r4Size.X/2.0, 0.0), 500*time.Millisecond)
			aPos2.AutoReverse = true

			aColor7 := ledgrid.NewColorAnimation(&r2.BorderColor, color.Cornsilk, 500*time.Millisecond)
			aColor7.AutoReverse = true
			aBorder1 := ledgrid.NewFloatAnimation(&r2.BorderWidth, 2.0, 500*time.Millisecond)
			aBorder1.AutoReverse = true

			aColor8 := ledgrid.NewColorAnimation(&r4.BorderColor, color.Cornsilk, 500*time.Millisecond)
			aColor8.AutoReverse = true
			aBorder2 := ledgrid.NewFloatAnimation(&r4.BorderWidth, 2.0, 500*time.Millisecond)
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
			aColor := ledgrid.NewColorAnimation(&r.BorderColor, rColor2, 4*time.Second)
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
			pathA := ledgrid.FullCirclePathA
			pathB := ledgrid.FullCirclePathB

			pos1 := geom.Point{2, float64(height) / 2.0}
			pos2 := geom.Point{float64(width) / 2.0, 2}
			pos3 := geom.Point{float64(width) - 2, float64(height) / 2.0}
			pos4 := geom.Point{float64(width) / 2.0, float64(height) - 2}
			cSize := geom.Point{3.0, 3.0}

			c1 := ledgrid.NewEllipse(pos1, cSize, color.OrangeRed)
			c2 := ledgrid.NewEllipse(pos2, cSize, color.MediumSeaGreen)
			c3 := ledgrid.NewEllipse(pos3, cSize, color.SkyBlue)
			c4 := ledgrid.NewEllipse(pos4, cSize, color.Gold)
			c.Add(c1, c2, c3, c4)

			c1Path := ledgrid.NewPathAnimation(&c1.Pos, pathB, geom.Point{float64(width - 4), float64(height) - 4.0}, duration)
			c1Path.AutoReverse = true
			c3Path := ledgrid.NewPathAnimation(&c3.Pos, pathB, geom.Point{-float64(width - 4), -(float64(height) - 4.0)}, duration)
			c3Path.AutoReverse = true

			c2Path := ledgrid.NewPathAnimation(&c2.Pos, pathA, geom.Point{float64(width) / 3.0, float64(height) - 4.0}, duration)
			c2Path.AutoReverse = true
			c4Path := ledgrid.NewPathAnimation(&c4.Pos, pathA, geom.Point{-float64(width) / 3.0, -(float64(height) - 4.0)}, duration)
			c4Path.AutoReverse = true

			aGrp := ledgrid.NewGroup(c1Path, c2Path, c3Path, c4Path)
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
			c1.BorderWidth = 0.0
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

			stepRD := geom.Point{9.0, float64(height) - 3.0}
			stepLU := stepRD.Neg()
			stepRU := geom.Point{9.0, -(float64(height) - 3.0)}
			stepLD := stepRU.Neg()

			c1Path1 := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.QuarterCirclePathA, stepRD, time.Second)
			c1Path1.Cont = true
			c1Path2 := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.QuarterCirclePathA, stepRU, time.Second)
			c1Path2.Cont = true
			c1Path3 := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.QuarterCirclePathA, stepLD, time.Second)
			c1Path3.Cont = true
			c1Path4 := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.QuarterCirclePathA, stepLU, time.Second)
			c1Path4.Cont = true

			c2Path1 := ledgrid.NewPathAnimation(&c2.Pos, ledgrid.QuarterCirclePathA, stepLU, time.Second)
			c2Path1.Cont = true
			c2Path2 := ledgrid.NewPathAnimation(&c2.Pos, ledgrid.QuarterCirclePathA, stepRD, time.Second)
			c2Path2.Cont = true

			c3Path1 := ledgrid.NewPathAnimation(&c3.Pos, ledgrid.QuarterCirclePathA, stepLD, time.Second)
			c3Path1.Cont = true
			c3Path2 := ledgrid.NewPathAnimation(&c3.Pos, ledgrid.QuarterCirclePathA, stepRU, time.Second)
			c3Path2.Cont = true

			c4Path1 := ledgrid.NewPathAnimation(&c4.Pos, ledgrid.QuarterCirclePathA, stepLU, time.Second)
			c4Path1.Cont = true
			c4Path2 := ledgrid.NewPathAnimation(&c4.Pos, ledgrid.QuarterCirclePathA, stepRD, time.Second)
			c4Path2.Cont = true

			c5Path1 := ledgrid.NewPathAnimation(&c5.Pos, ledgrid.QuarterCirclePathA, stepLD, time.Second)
			c5Path1.Cont = true
			c5Path2 := ledgrid.NewPathAnimation(&c5.Pos, ledgrid.QuarterCirclePathA, stepRU, time.Second)
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
			c1PosSize := geom.Point{-(float64(width) - 10.0), -float64(height) / 2.0}
			c2Pos := geom.Point{5.0, float64(height) / 2.0}
			c2Size1 := geom.Point{5.0, 5.0}
			c2Size2 := geom.Point{3.0, 3.0}
			c2PosSize := geom.Point{float64(width) - 10.0, float64(height)/2.0 + 2.0}

			aGrp := ledgrid.NewGroup()

			pal := ledgrid.NewGradientPaletteByList("Palette", true,
				color.DeepSkyBlue,
				color.Lime,
				color.Teal,
				color.SkyBlue,
			)

			c1 := ledgrid.NewEllipse(c1Pos1, c1Size1, color.Gold)

			c1pos := ledgrid.NewPathAnimation(&c1.Pos, ledgrid.FullCirclePathB, c1PosSize, 4*time.Second)
			c1pos.Curve = ledgrid.AnimationLinear

			c1size := ledgrid.NewSizeAnimation(&c1.Size, c1Size2, time.Second)
			c1size.AutoReverse = true

			c1bcolor := ledgrid.NewColorAnimation(&c1.BorderColor, color.OrangeRed, time.Second)
			c1bcolor.AutoReverse = true

			c2 := ledgrid.NewEllipse(c2Pos, c2Size1, color.Lime)

			c2pos := ledgrid.NewPathAnimation(&c2.Pos, ledgrid.FullCirclePathB, c2PosSize, 4*time.Second)
			c2pos.Curve = ledgrid.AnimationLinear

			c2size := ledgrid.NewSizeAnimation(&c2.Size, c2Size2, time.Second)
			c2size.AutoReverse = true

			c2color := ledgrid.NewPaletteAnimation(&c2.BorderColor, pal, 2*time.Second)
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

			c1color := ledgrid.NewColorAnimation(&c1.BorderColor, color.Gold, 4*time.Second)
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

			a1Color := ledgrid.NewColorAnimation(&r1.BorderColor, color.GreenYellow, duration)
			a1Color.AutoReverse = true
			a1Color.RepeatCount = ledgrid.AnimationRepeatForever

			r2 := ledgrid.NewRectangle(r2Pos2, rSize1, color.SkyBlue)

			a2Pos := ledgrid.NewPositionAnimation(&r2.Pos, r2Pos1, duration)
			a2Pos.AutoReverse = true
			a2Pos.RepeatCount = ledgrid.AnimationRepeatForever

			a2Size := ledgrid.NewSizeAnimation(&r2.Size, rSize2, duration)
			a2Size.AutoReverse = true
			a2Size.RepeatCount = ledgrid.AnimationRepeatForever

			a2Color := ledgrid.NewColorAnimation(&r2.BorderColor, color.Crimson, duration)
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
				animFade := ledgrid.NewColorAnimation(&polyList[n].BorderColor, color.Black, 4*time.Second)
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
			t1.AX = 0.0
			t2 := ledgrid.NewText(geom.Point{float64(width), float64(height) / 2.0}, "Beni", color.YellowGreen)
			t2.AX = 1.0
			t3 := ledgrid.NewText(geom.Point{float64(width) + 60.0, float64(height) / 2.0}, "wohnen im Lochbach", color.OrangeRed)

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

			/*
			   			pos1 := geom.Point{float64(width) / 2.0, -float64(height) / 2.0}
			   			pos2 := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			   			pos3 := geom.Point{float64(width) / 2.0, 3 * float64(height) / 2.0}

			   			size1 := geom.Point{float64(width) / 3.0, float64(height) / 3.0}
			   			size2 := geom.Point{3.0 * float64(width), 3.0 * float64(height)}

			   			size3 := geom.Point{1.0, 0.75}
			   			size4 := geom.Point{160.0, 120.0}

			               img1 := ledgrid.NewImage(pos2, "images/colormap01.png")

			   			img1 := ledgrid.NewImage(pos1, "images/ledgrid.png")
			   			img1.Size = size1
			   			// img1.Hide()

			   			img2 := ledgrid.NewImage(pos2, "images/testbild.png")
			   			img2.Size = size3
			   			img2.Hide()

			   			c.Add(img1, img2)

			   			aPos1 := ledgrid.NewPositionAnimation(&img1.Pos, pos2, 4*time.Second)
			   			aPos2 := ledgrid.NewPositionAnimation(&img1.Pos, pos3, 4*time.Second)
			   			aPos2.Cont = true

			   			aAngle1 := ledgrid.NewFloatAnimation(&img1.Angle, math.Pi, 4*time.Second)
			   			aAngle2 := ledgrid.NewFloatAnimation(&img1.Angle, 0.0, 4*time.Second)
			   			aAngle2.Cont = true

			   			aSize1 := ledgrid.NewSizeAnimation(&img1.Size, size2, 4*time.Second)
			   			aSize1.Cont = true
			   			aSize2 := ledgrid.NewSizeAnimation(&img1.Size, size1, 4*time.Second)
			   			aSize2.Cont = true

			   			aSize3 := ledgrid.NewSizeAnimation(&img2.Size, size4, 4*time.Second)
			   			aSize3.AutoReverse = true
			   			aSize3.RepeatCount = ledgrid.AnimationRepeatForever
			   			aAngle3 := ledgrid.NewFloatAnimation(&img2.Angle, 4*math.Pi, 4*time.Second)
			   			aAngle3.AutoReverse = true
			   			aAngle3.RepeatCount = ledgrid.AnimationRepeatForever

			   			// task := ledgrid.NewBackgroundTask(func() { img2.Show() })
			   			showHide := ledgrid.NewShowHideAnimation(img2)

			   			aSize3.Start()
			   			aAngle3.Start()
			   			aSeq := ledgrid.NewSequence(aPos1, aAngle1, showHide, aSize1, aPos2, showHide, aAngle2, aSize2, aSize3)
			   			aSeq.Start()
			*/
		})

	ImageFilterTest = NewLedGridProgram("Image filter test",
		func(c *ledgrid.Canvas) {
			filter := NewShuffleFilter(c)
			ledgrid.AnimCtrl.Filter = filter
			// shuffleTask := ledgrid.NewTask(filter.ShuffleOne)
			// filter.Shuffle()
			// orderTask := ledgrid.NewTask(filter.OrderOne)
			// aGrpShuffle := ledgrid.NewGroup(shuffleTask)
			// aGrpShuffle.SetDuration(200 * time.Millisecond)
			// aGrpShuffle.RepeatCount = 100
			// aGrpOrder := ledgrid.NewGroup(orderTask)
			// aGrpOrder.SetDuration(50 * time.Millisecond)
			// aGrpOrder.RepeatCount = 400
			bisectTask := ledgrid.NewTask(func() {
                i, j := rand.IntN(width), rand.IntN(width)
				filter.SwapCols(i, j)
			})
            resetTask := ledgrid.NewTask(filter.Reset)
			aGrpBisect := ledgrid.NewGroup(bisectTask)
			aGrpBisect.SetDuration(1000 * time.Millisecond)
			aGrpBisect.RepeatCount = 10
			infoTask := ledgrid.NewTask(func() {
				fmt.Printf("Done...\n")
			})
			aSeq := ledgrid.NewSequence(aGrpBisect, infoTask, resetTask)

			pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			size := geom.Point{float64(width), float64(height)}
			cam := NewCamera(pos, size)
			c.Add(cam)
			cam.Start()
			aSeq.Start()
		})

	CameraTest = NewLedGridProgram("Camera test",
		func(c *ledgrid.Canvas) {
			// c.BackColor = color.DarkGray
			pos := geom.Point{float64(width) / 2.0, float64(height) / 2.0}
			size := geom.Point{float64(width), float64(height)}

			cam := NewCamera(pos, size)
			c.Add(cam)
			cam.Start()
			mask := cam.Mask

			go func() {
				var pt, ptDest geom.Point
				// scanDirList := []ScanDir{Top2Bottom, Border2Center, Center2Border, Bottom2Top, Left2Right, Right2Left}
				// scanLineList := []ScanLine{Forward, Forward, Forward, Alternate, Backward, Alternate}
				// exitDirList := []ExitDir{ExitOver, ExitAway, ExitAway, ExitAway, ExitOver, ExitAway}

				scanDirList := []ScanDir{Out2Inside, In2Outside}
				scanLineList := []ScanLine{CW, CW}
				exitDirList := []ExitDir{ExitAway, ExitAway}

				time.Sleep(3 * time.Second)
				for i, scanDir := range scanDirList {
					scanLine := scanLineList[i]
					exitDir := exitDirList[i]
					for i := range width * height {
						pp := StartEndPoints(i, scanDir, scanLine, exitDir)
						p0, p1 := pp.src, pp.dst
						pt = geom.NewPointIMG(p0)
						ptDest = geom.NewPointIMG(p1)
						pixAway := ledgrid.NewDot(pt, color.OrangeRed)
						c.Add(pixAway)
						aMask := ledgrid.NewTask(func() {
							idx := mask.PixOffset(p0.X, p0.Y)
							mask.Pix[idx] = 0x00
						})

						aDur := time.Second + rand.N(time.Second)
						aColor2 := ledgrid.NewUint8Animation(&pixAway.Color.A, 0, aDur)
						aColor2.Curve = ledgrid.AnimationEaseIn
						aPos := ledgrid.NewPositionAnimation(&pixAway.Pos, ptDest, aDur)
						aPos.Curve = ledgrid.AnimationEaseIn
						aGrp2 := ledgrid.NewGroup(aMask, aColor2, aPos)
						aGrp2.Start()
						// if i%width == 0 {
						time.Sleep(30 * time.Millisecond)
						// }
					}
					time.Sleep(2 * time.Second)
					for i := range mask.Pix {
						mask.Pix[i] = 0xff
					}
				}
			}()
		})

	BlinkenAnimation = NewLedGridProgram("Blinken animation",
		func(c *ledgrid.Canvas) {
			posFlame1 := geom.Point{4.5, float64(height) - 4.5}
			posFlame2 := geom.Point{float64(width) - 4.5, float64(height) - 4.5}
			pos1Mario := geom.Point{10.0, float64(height) / 2.0}
			pos2Mario := geom.Point{float64(width) - 10.0, float64(height) / 2.0}

			bmlFlame := ledgrid.ReadBlinkenFile("blinken/flameNew.bml")
			bmlFlame.SetAllDuration(32)

			flame1 := ledgrid.NewImageList(posFlame1)
			flame1.AddBlinkenLight(bmlFlame)
			flame1.RepeatCount = ledgrid.AnimationRepeatForever

			bmlFlame.SetAllDuration(43)

			flame2 := ledgrid.NewImageList(posFlame2)
			flame2.AddBlinkenLight(bmlFlame)
			flame2.RepeatCount = ledgrid.AnimationRepeatForever

			bmlMario := ledgrid.ReadBlinkenFile("blinken/marioWalkRight.bml")

			mario := ledgrid.NewImageList(pos1Mario)
			mario.AddBlinkenLight(bmlMario)
			mario.RepeatCount = ledgrid.AnimationRepeatForever
			mario.Size = geom.Point{10.0, 10.0}

			aPos := ledgrid.NewPositionAnimation(&mario.Pos, pos2Mario, 3*time.Second)
			aPos.Curve = ledgrid.AnimationLinear
			aPos.RepeatCount = ledgrid.AnimationRepeatForever

			c.Add(flame1, flame2, mario)

			aGrp := ledgrid.NewGroup(flame1, flame2, mario, aPos)
			aGrp.Start()
		})

	MovingPixels = NewLedGridProgram("Moving pixels",
		func(c *ledgrid.Canvas) {
			pos1 := image.Point{0, 0}
			pos2 := image.Point{width - 1, 0}

			pos5 := geom.Point{0, 2}
			pos6 := geom.Point{float64(width) - 1, 2}

			pos3 := image.Point{width - 1, 0}
			pos4 := image.Point{0, height - 3}

			pos7 := geom.Point{float64(width) - 1, 2}
			pos8 := geom.Point{0, float64(height) - 1}

			pix1 := ledgrid.NewPixel(pos1, color.YellowGreen)
			pix2 := ledgrid.NewPixel(pos3, color.LightSeaGreen)
			pix3 := ledgrid.NewDot(pos5, color.YellowGreen)
			pix4 := ledgrid.NewDot(pos7, color.LightSeaGreen)
			c.Add(pix1, pix2, pix3, pix4)

			aPos1 := ledgrid.NewIntegerPosAnimation(&pix1.Pos, pos2, 5*time.Second)
			aPos1.AutoReverse = true
			aPos1.RepeatCount = ledgrid.AnimationRepeatForever
			aPos2 := ledgrid.NewIntegerPosAnimation(&pix2.Pos, pos4, 5*time.Second)
			aPos2.AutoReverse = true
			aPos2.RepeatCount = ledgrid.AnimationRepeatForever
			aPos3 := ledgrid.NewPositionAnimation(&pix3.Pos, pos6, 5*time.Second)
			aPos3.AutoReverse = true
			aPos3.RepeatCount = ledgrid.AnimationRepeatForever
			aPos4 := ledgrid.NewPositionAnimation(&pix4.Pos, pos8, 5*time.Second)
			aPos4.AutoReverse = true
			aPos4.RepeatCount = ledgrid.AnimationRepeatForever
			aPos1.Start()
			aPos2.Start()
			aPos3.Start()
			aPos4.Start()

			aColor1 := ledgrid.NewColorAnimation(&pix1.Color, color.OrangeRed, 2000*time.Millisecond)
			aColor1.AutoReverse = true
			aColor1.RepeatCount = ledgrid.AnimationRepeatForever
			aColor1.Start()
			aColor2 := ledgrid.NewColorAnimation(&pix2.Color, color.Purple, 2123*time.Millisecond)
			aColor2.AutoReverse = true
			aColor2.RepeatCount = ledgrid.AnimationRepeatForever
			aColor2.Start()
			aColor3 := ledgrid.NewColorAnimation(&pix3.Color, color.OrangeRed, 2000*time.Millisecond)
			aColor3.AutoReverse = true
			aColor3.RepeatCount = ledgrid.AnimationRepeatForever
			aColor3.Start()
			aColor4 := ledgrid.NewColorAnimation(&pix4.Color, color.Purple, 2123*time.Millisecond)
			aColor4.AutoReverse = true
			aColor4.RepeatCount = ledgrid.AnimationRepeatForever
			aColor4.Start()
		})

	GlowingPixels = NewLedGridProgram("Glowing pixels",
		func(c *ledgrid.Canvas) {
			aGrpPurple := ledgrid.NewGroup()
			aGrpYellow := ledgrid.NewGroup()
			aGrpGreen := ledgrid.NewGroup()
			aGrpGrey := ledgrid.NewGroup()

			for y := range c.Rect.Dy() {
				for x := range c.Rect.Dx() {
					pt := image.Point{x, y}
					// pos := pt
					pos := geom.NewPointIMG(pt)
					t := rand.Float64()
					col := (color.DimGray.Dark(0.3)).Interpolate((color.DarkGrey.Dark(0.3)), t)
					// pix := ledgrid.NewPixel(pos, col)
					pix := ledgrid.NewDot(pos, col)
					c.Add(pix)

					dur := time.Second + time.Duration(x)*time.Millisecond
					aAlpha := ledgrid.NewUint8Animation(&pix.Color.A, 196, dur)
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

			txt1 := ledgrid.NewFixedText(fixed.P(1, height-1), color.GreenYellow.Alpha(0.0), "LORENZ")
			aTxt1 := ledgrid.NewUint8Animation(&txt1.Color.A, 255, 2*time.Second)
			aTxt1.AutoReverse = true
			txt2 := ledgrid.NewFixedText(fixed.P(1, height-1), color.DarkViolet.Alpha(0.0), "SIMON")
			aTxt2 := ledgrid.NewUint8Animation(&txt2.Color.A, 255, 2*time.Second)
			aTxt2.AutoReverse = true
			txt3 := ledgrid.NewFixedText(fixed.P(1, height-1), color.OrangeRed.Alpha(0.0), "REBEKKA")
			aTxt3 := ledgrid.NewUint8Animation(&txt3.Color.A, 255, 2*time.Second)
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

	FarewellGery = NewLedGridProgram("Farewell Gery!",
		func(c *ledgrid.Canvas) {
			aGrpPurple := ledgrid.NewGroup()
			aGrpYellow := ledgrid.NewGroup()
			aGrpGreen := ledgrid.NewGroup()
			aGrpGrey := ledgrid.NewGroup()
			aGrpRed := ledgrid.NewGroup()
			aGrpBlack := ledgrid.NewGroup()

			for y := range c.Rect.Dy() {
				for x := range c.Rect.Dx() {
					pt := image.Point{x, y}
					pos := geom.NewPointIMG(pt)
					t := rand.Float64()
					col := (color.DimGray.Dark(0.5)).Interpolate((color.DarkGrey.Dark(0.5)), t)
					pix := ledgrid.NewDot(pos, col)
					c.Add(pix)

					dur := time.Second + time.Duration(10*x)*time.Millisecond
					aAlpha := ledgrid.NewUint8Animation(&pix.Color.A, 196, dur)
					aAlpha.AutoReverse = true
					aAlpha.RepeatCount = ledgrid.AnimationRepeatForever
					aAlpha.Start()

					aColor := ledgrid.NewColorAnimation(&pix.Color, col, 1*time.Second)
					aColor.Cont = true
					aGrpGrey.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.MediumPurple.Interpolate(color.Fuchsia, t), 4*time.Second)
					aColor.Cont = true
					aGrpPurple.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.Gold.Interpolate(color.Khaki, t), 4*time.Second)
					aColor.Cont = true
					aGrpYellow.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.GreenYellow.Interpolate(color.LightSeaGreen, t), 4*time.Second)
					aColor.Cont = true
					aGrpGreen.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.OrangeRed.Interpolate(color.DarkOrange, t), 4*time.Second)
					aColor.Cont = true
					aGrpRed.Add(aColor)

					aColor = ledgrid.NewColorAnimation(&pix.Color, color.Black, 2*time.Second)
					aColor.Cont = true
					aGrpBlack.Add(aColor)
				}
			}

			txt1 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.GreenYellow.Alpha(0.0), "LIEBER")
			aTxt1 := ledgrid.NewUint8Animation(&txt1.Color.A, 255, 1*time.Second)
			aTxt1.AutoReverse = true
			txt2 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.DarkViolet.Alpha(0.0), "GERY")
			aTxt2 := ledgrid.NewUint8Animation(&txt2.Color.A, 255, 2*time.Second)
			aTxt2.AutoReverse = true
			txt3 := ledgrid.NewFixedText(fixed.P(width/2, height/2), color.OrangeRed.Alpha(0.0), "FAREWELL")
			aTxt3 := ledgrid.NewUint8Animation(&txt3.Color.A, 255, 5*time.Second)
			aTxt3.AutoReverse = true
			c.Add(txt1, txt2, txt3)

			aTimel := ledgrid.NewTimeline(40 * time.Second)
			aTimel.Add(5*time.Second, aGrpPurple)
			aTimel.Add(8*time.Second, aTxt1)
			aTimel.Add(9*time.Second, aGrpGrey)

			aTimel.Add(15*time.Second, aGrpYellow)
			aTimel.Add(18*time.Second, aTxt2)
			aTimel.Add(19*time.Second, aGrpGrey)

			aTimel.Add(25*time.Second, aGrpGreen)
			aTimel.Add(28*time.Second, aTxt3)
			aTimel.Add(32*time.Second, aGrpBlack)
			aTimel.Add(38*time.Second, aGrpGrey)
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
				txt.SetText(name)
				// log.Printf(">>> %s", name)
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
			c.Add(txt)
			aPalTl.Start()
			aGrp.Start()
		})

	WaterDrops = NewLedGridProgram("Water drops",
		func(c *ledgrid.Canvas) {

		})
)

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
		ImageFilterTest,
		CameraTest,
		BlinkenAnimation,
		MovingText,
		BitmapText,
		SlideTheShow,
		ShowTheShader,
	}
)

func main() {
	var host string
	var port uint
	var input string
	var ch byte
	var progId int
	var runInteractive bool
	var progList string

	for i, prog := range programList {
		id := 'a' + i
		progList += fmt.Sprintf("\n%c - %s", id, prog.Name())
	}

	flag.StringVar(&host, "host", defHost, "Controller hostname")
	flag.UintVar(&port, "port", defPort, "Controller port")
	flag.StringVar(&input, "prog", input, "Play one single program"+progList)
	flag.Parse()

	if len(input) > 0 {
		runInteractive = false
		ch = input[0]
	} else {
		runInteractive = true
	}

	ledGrid = ledgrid.NewLedGridV2(host, port)

	gridSize = ledGrid.Rect.Size()
	width = gridSize.X
	height = gridSize.Y

	canvas = ledgrid.NewCanvas(gridSize)
	ledgrid.NewAnimationController(canvas, ledGrid)

    	initSpiralMap(CW)
	initSpiralMap(CCW)

	if runInteractive {
		progId = -1
		for {
			fmt.Printf("---------------------------------------\n")
			fmt.Printf("Program\n")
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
			fmt.Printf("  S - Save Canvas\n")
			fmt.Printf("  L - Load Canvas from File\n")
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
