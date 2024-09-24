package main

import (
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

			cPos := geom.Point{1, 1}

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

			c1 := ledgrid.NewDot(cPos, color.GreenYellow)
			c.Add(c1)

			aPath1 := ledgrid.NewPolyPathAnimation(&c1.Pos, polyPath1, 7*time.Second)
			aPath1.AutoReverse = true

			aPath2 := ledgrid.NewPolyPathAnimation(&c1.Pos, polyPath2, 7*time.Second)
			aPath2.AutoReverse = true

			seq := ledgrid.NewSequence(aPath1, aPath2)
			seq.RepeatCount = ledgrid.AnimationRepeatForever

			seq.Start()
		})
)
