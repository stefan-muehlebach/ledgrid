package ledgrid

import (
	"fmt"
	"image"
	"testing"
)

const (
	Width, Height = 20, 20
	RandSeed      = 123_456
)

type Coord2Offset struct {
	coord image.Point
	idx   int
}

var (
	lg        *LedGrid
	idx       int
	modLayout = ModuleLayout{
		{
			{ModLR, Rot000}, {ModRL, Rot090},
		},
	}
	coordList = []Coord2Offset{
		{image.Point{0, 0}, 0},
		{image.Point{Width - 1, 0}, 3 * (Height*Width - 1)},
		{image.Point{0, Height - 1}, 3 * (Width - 1)},
		{image.Point{Width - 1, Height - 1}, 3 * Width * (Height - 1)},
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

// func BenchmarkPixOffsetCalc(b *testing.B) {
// 	rand.Seed(RandSeed)
// 	for i := 0; i < b.N; i++ {
// 		x, y := rand.Intn(Width), rand.Intn(Height)
// 		idx = lg.pixOffset(x, y)
// 	}
// }

// func BenchmarkPixOffsetMap(b *testing.B) {
// 	rand.Seed(RandSeed)
// 	for i := 0; i < b.N; i++ {
// 		x, y := rand.Intn(Width), rand.Intn(Height)
// 		idx = lg.PixOffset(x, y)
// 	}
// }

func TestModuleScanner(t *testing.T) {
	var mod Module
	input := "LR:90"

	n, err := fmt.Sscanf(input, "%v", &mod)
	t.Logf("n: %d, err: %v, mod: %+v", n, err, mod)
}
