package main

import (
	"math/rand"
	"time"

	"github.com/stefan-muehlebach/gg/geom"
	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
)

var (
	PathTest = NewLedGridProgram("Path test",
		func(c *ledgrid.Canvas) {
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
			c.Add(0, c1, c2, c3, c4)

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
		})

	PolygonPathTest = NewLedGridProgram("Polygon path test",
		func(c *ledgrid.Canvas) {

			cPos := geom.Point{1, 1}

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
				geom.Point{1, 1},
			}
			for range 20 {
				ptList = append(ptList, geom.Point{float64(rand.Intn(width - 1)), float64(rand.Intn(height - 1))})
			}
			polyPath3 := ledgrid.NewPolygonPath(ptList...)

			c1 := ledgrid.NewDot(cPos, color.GreenYellow)
			c.Add(0, c1)

			aPath1 := ledgrid.NewPolyPathAnim(c1, polyPath1, 14*time.Second)
			aPath1.AutoReverse = true

			aPath2 := ledgrid.NewPolyPathAnim(c1, polyPath2, 5*time.Second)
			aPath2.AutoReverse = true

			aPath3 := ledgrid.NewPolyPathAnim(c1, polyPath3, 5*time.Second)
			aPath3.AutoReverse = true

			seq := ledgrid.NewSequence(aPath1, aPath3, aPath2)
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
			c.Add(0, c1, c2)

			aPos1 := ledgrid.NewPositionAnim(c1, geom.Point{}, 1300*time.Millisecond)
			aPos1.Cont = true
			aPos1.Val2 = ledgrid.RandPointTrunc(rect, 1.0)
			aPos1.RepeatCount = ledgrid.AnimationRepeatForever

			aPos2 := ledgrid.NewPositionAnim(c2, geom.Point{}, 901*time.Millisecond)
			aPos2.Cont = true
			aPos2.Val2 = func() geom.Point { return c1.Pos }
			aPos2.Curve = ledgrid.AnimationLinear
			aPos2.RepeatCount = ledgrid.AnimationRepeatForever

			aPos1.Start()
			aPos2.Start()
		})
)
