package ledgrid

import (
	"fmt"
	"image"
	"math/rand"
	"testing"
)

const (
	Width, Height = 20, 10
	RandSeed      = 123_456
)

type Coord2Offset struct {
	coord image.Point
	idx   int
}

var (
	lg        *LedGrid
	idx       int
	coordList = []Coord2Offset{
		{image.Point{0, 0}, 0},
		{image.Point{9, 0}, 3 * (ModuleSize.X*ModuleSize.Y - 1)},
		{image.Point{10, 0}, 3 * (ModuleSize.X * ModuleSize.Y)},
		{image.Point{10, 9}, 3 * (2*ModuleSize.X*ModuleSize.Y - 1)},
		// {image.Point{10, 10}, 3 * (2 * ModuleSize.X * ModuleSize.Y)},
		// {image.Point{10, 19}, 3 * (3*ModuleSize.X*ModuleSize.Y - 1)},
		// {image.Point{9, 19}, 3 * (3 * ModuleSize.X * ModuleSize.Y)},
		// {image.Point{0, 19}, 3 * (4*ModuleSize.X*ModuleSize.Y - 1)},
	}
)

func init() {
	lg = NewLedGrid(image.Point{Width, Height})
}

func TestPixOffset(t *testing.T) {
	for _, rec := range coordList {
		pt := rec.coord
		refIdx := rec.idx
		idx = lg.PixOffset(pt.X, pt.Y)
		if refIdx == idx {
			t.Logf("(%d,%d) -> %d, OK", pt.X, pt.Y, idx)
		} else {
			t.Errorf("(%d,%d) -> %d, should be %d", pt.X, pt.Y, idx, refIdx)
		}
	}
}

func TestMarkDefect(t *testing.T) {
    lg.idxMap.MarkDefect(image.Point{6,3})
}

// func BenchmarkPixOffsetCalc(b *testing.B) {
// 	rand.Seed(RandSeed)
// 	for i := 0; i < b.N; i++ {
// 		x, y := rand.Intn(Width), rand.Intn(Height)
// 		idx = lg.pixOffset(x, y)
// 	}
// }

func BenchmarkPixOffset(b *testing.B) {
	rand.Seed(RandSeed)
	for i := 0; i < b.N; i++ {
		x, y := rand.Intn(Width), rand.Intn(Height)
		idx = lg.PixOffset(x, y)
	}
}

func TestModuleScanner(t *testing.T) {
	var mod Module
	input := "LR:90"

	n, err := fmt.Sscanf(input, "%v", &mod)
	t.Logf("n: %d, err: %v, mod: %+v", n, err, mod)
}
