package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

func init() {
    // programList.AddTitle("Path Animations")
	programList.Add("Path test", "Path", PathTest)
	programList.Add("Polygon path test", "Path", PolygonPathTest)
}

func PathTest(ctx context.Context, c *ledgrid.Canvas) {
	duration := 4 * time.Second
	pathA := ledgrid.CirclePath
	pathB := ledgrid.CirclePath.NewStart(0.25)

	pos1 := geom.Point{float64(width) / 2.0, 2.0}
	pos2 := geom.Point{float64(width) - 2.0, float64(height) / 2.0}
	pos3 := geom.Point{float64(width) / 2.0, float64(height) - 2.0}
	pos4 := geom.Point{2.0, float64(height) / 2.0}
	cSize := geom.Point{3.0, 3.0}

	c1 := ledgrid.NewEllipse(pos1, cSize, color.OrangeRed)
	c2 := ledgrid.NewEllipse(pos2, cSize, color.MediumSeaGreen)
	c3 := ledgrid.NewEllipse(pos3, cSize, color.SkyBlue)
	c4 := ledgrid.NewEllipse(pos4, cSize, color.Gold)
	c.Add(c1, c2, c3, c4)

	c1Path := ledgrid.NewPathAnim(c1, pathA, geom.Point{float64(width) / 3.0, float64(height - 4)}, duration)
	c1Path.AutoReverse = true

	c2Path := ledgrid.NewPathAnim(c2, pathB, geom.Point{float64(width - 4), float64(height - 4)}, duration)
	c2Path.AutoReverse = true

	c3Path := ledgrid.NewPathAnim(c3, pathA, geom.Point{-float64(width) / 3.0, -float64(height - 4)}, duration)
	c3Path.AutoReverse = true

	c4Path := ledgrid.NewPathAnim(c4, pathB, geom.Point{-float64(width - 4), -float64(height - 4)}, duration)
	c4Path.AutoReverse = true

	aGrp := ledgrid.NewGroup(c1Path, c3Path, c2Path, c4Path)
	aGrp.RepeatCount = ledgrid.AnimationRepeatForever
	aGrp.Start()
}

func PolygonPathTest(ctx context.Context, c *ledgrid.Canvas) {

	cPos := geom.Point{0, 0}

	polyPath1 := ledgrid.NewPolygonPath(
		geom.Point{1, 1},
		geom.Point{float64(width) - 2, 1},
		geom.Point{float64(width) - 2, float64(height) - 2},
		geom.Point{1, float64(height) - 2},
		geom.Point{1, 2},
		geom.Point{float64(width) - 3, 2},
		geom.Point{float64(width) - 3, float64(height) - 3},
		geom.Point{2, float64(height) - 3},
		geom.Point{2, 3},
		geom.Point{float64(width) - 4, 3},
		geom.Point{float64(width) - 4, float64(height) - 4},
		geom.Point{3, float64(height) - 4},
		geom.Point{3, 4},
		geom.Point{float64(width) - 5, 4},
		geom.Point{float64(width) - 5, float64(height) - 5},
		geom.Point{4, float64(height) - 5},
	)

	polyPath2 := ledgrid.NewPolygonPath(
		geom.Point{1, 1},
		geom.Point{4, 9},
		geom.Point{7, 2},
		geom.Point{10, 8},
		geom.Point{13, 3},
		geom.Point{16, 7},
		geom.Point{19, 4},
		geom.Point{22, 6},
	)

	ptList := []geom.Point{
		geom.Point{0, 0},
	}
	lastSide := 0
	for range 20 {
		pt := geom.Point{}
		side := rand.Intn(3)
		if side >= lastSide {
			side += 1
		}
		switch side {
		case 0:
			pt = geom.Point{0.0, float64(rand.Intn(height - 1))}
		case 1:
			pt = geom.Point{float64(rand.Intn(width - 1)), 0.0}
		case 2:
			pt = geom.Point{float64(width - 1), float64(rand.Intn(height - 1))}
		case 3:
			pt = geom.Point{float64(rand.Intn(width - 1)), float64(height - 1)}
		}
		ptList = append(ptList, pt)
		lastSide = side
	}
	polyPath3 := ledgrid.NewPolygonPath(ptList...)

	c1 := ledgrid.NewDot(cPos, color.GreenYellow)
	c.Add(c1)

	aPath1 := ledgrid.NewPolyPathAnim(c1, polyPath1, 14*time.Second)
	aPath1.AutoReverse = true

	aPath2 := ledgrid.NewPolyPathAnim(c1, polyPath2, 5*time.Second)
	aPath2.AutoReverse = true

	aPath3 := ledgrid.NewPolyPathAnim(c1, polyPath3, 10*time.Second)
	aPath3.AutoReverse = true

	seq := ledgrid.NewSequence(aPath1, aPath2, aPath3)
	seq.RepeatCount = ledgrid.AnimationRepeatForever

	seq.Start()
}
