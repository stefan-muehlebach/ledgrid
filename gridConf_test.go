package ledgrid

import (
	"image"
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

// func init() {
// 	lg = NewLedGrid(image.Point{Width, Height})
// }

func TestDefaultModuleConfig(t *testing.T) {
	conf := DefaultModuleConfig(image.Point{20, 20})
	t.Logf("module config: %v", conf)
	idxMap := conf.IndexMap()
	t.Logf("  idx( 0, 0): %d", idxMap[0][0]/3)
	t.Logf("  idx(19, 0): %d", idxMap[19][0]/3)
	t.Logf("  idx(19,19): %d", idxMap[19][19]/3)
	t.Logf("  idx( 0,19): %d", idxMap[0][19]/3)
}

func TestCoordMap(t *testing.T) {
	conf := DefaultModuleConfig(image.Point{40, 10})
	t.Logf("module config: %v", conf)
	idxMap := conf.IndexMap()
    coordMap := idxMap.CoordMap()
	t.Logf("  idxMap: %v", idxMap)
	t.Logf("  coordMap: %v", coordMap)
}

func TestModuleConfig(t *testing.T) {
	var modConf ModuleConfig
	var err error

	modConf, err = modConf.Append(0, 0, Module{ModLR, Rot000})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(1, 0, Module{ModLR, Rot000})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(2, 0, Module{ModLR, Rot000})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(3, 0, Module{ModRL, Rot090})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(3, 1, Module{ModRL, Rot090})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(2, 1, Module{ModLR, Rot180})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(1, 1, Module{ModLR, Rot180})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(0, 1, Module{ModLR, Rot180})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(0, 2, Module{ModLR, Rot000})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(1, 2, Module{ModLR, Rot000})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(2, 2, Module{ModLR, Rot000})
	if err != nil {
		t.Errorf("%v", err)
	}
	modConf, err = modConf.Append(3, 2, Module{ModRL, Rot090})
	if err != nil {
		t.Errorf("%v", err)
	}

	t.Logf("module configuration: %v", modConf)
}

func TestPlotModuleConfig(t *testing.T) {
	conf := DefaultModuleConfig(image.Point{20, 20})
	conf.Plot("moduleConfig.png")
}

// func TestPixOffset(t *testing.T) {
// 	for _, rec := range coordList {
// 		pt := rec.coord
// 		refIdx := rec.idx
// 		idx = lg.PixOffset(pt.X, pt.Y)
// 		if refIdx == idx {
// 			t.Logf("(%d,%d) -> %d, OK", pt.X, pt.Y, idx)
// 		} else {
// 			t.Errorf("(%d,%d) -> %d, should be %d", pt.X, pt.Y, idx, refIdx)
// 		}
// 	}
// }

// func TestMarkDefect(t *testing.T) {
//     lg.idxMap.MarkDefect(image.Point{6,3})
// }

// func BenchmarkPixOffsetCalc(b *testing.B) {
// 	rand.Seed(RandSeed)
// 	for i := 0; i < b.N; i++ {
// 		x, y := rand.Intn(Width), rand.Intn(Height)
// 		idx = lg.pixOffset(x, y)
// 	}
// }

// func BenchmarkPixOffset(b *testing.B) {
// 	rand.Seed(RandSeed)
// 	for i := 0; i < b.N; i++ {
// 		x, y := rand.Intn(Width), rand.Intn(Height)
// 		idx = lg.PixOffset(x, y)
// 	}
// }

// func TestModuleScanner(t *testing.T) {
// 	var mod Module
// 	input := "LR:90"

// 	n, err := fmt.Sscanf(input, "%v", &mod)
// 	t.Logf("n: %d, err: %v, mod: %+v", n, err, mod)
// }
