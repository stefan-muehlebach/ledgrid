package conf

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
	idx       int
	coordList = []Coord2Offset{
		{image.Point{0, 0}, 0},
		{image.Point{9, 0}, 3 * (ModuleDim.X*ModuleDim.Y - 1)},
		{image.Point{10, 0}, 3 * (ModuleDim.X * ModuleDim.Y)},
		{image.Point{10, 9}, 3 * (2*ModuleDim.X*ModuleDim.Y - 1)},
	}

    ptList = []image.Point{
        image.Point{0, 0},
        image.Point{0, 9},
        image.Point{9, 9},
        image.Point{9, 0},
    }
    idxList = []int{0, 9, 90, 99}
    modTypeList = []ModuleType{ModLR, ModRL}
    modList = []Module{
        Module{ModLR, Rot180},
        Module{ModRL, Rot090},
    }
)

func TestModuleTypeIndex(t *testing.T) {
    for _, modType := range modTypeList {
        t.Logf("ModuleType: %v", modType)
        for _, pt := range ptList {
            t.Logf("  %v -> %2d", pt, modType.Index(pt))
        }
    }
}

func TestModuleTypeCoord(t *testing.T) {
    for _, modType := range modTypeList {
        t.Logf("ModuleType: %v", modType)
        for _, idx := range idxList {
            t.Logf("  %2d -> %v", idx, modType.Coord(idx))
        }
    }
}

func TestModuleIndex(t *testing.T) {
    for _, mod := range modList {
        t.Logf("Module: %v", mod)
        for _, pt := range ptList {
            t.Logf("  %v -> %2d", pt, mod.Index(pt))
        }
    }
}

func TestModuleCoord(t *testing.T) {
    for _, mod := range modList {
        t.Logf("Module: %v", mod)
        for _, idx := range idxList {
            t.Logf("  %2d -> %v", idx, mod.Coord(idx))
        }
    }
}

func TestDefaultModuleConfig(t *testing.T) {
    ptList = []image.Point{
        image.Point{0, 0},
        image.Point{0, 19},
        image.Point{19, 19},
        image.Point{19, 0},
    }
    idxList = []int{0, 109, 290, 399}

	modConf := DefaultModuleConfig(image.Point{20, 20})
	t.Logf("Module config: %v", modConf)
    t.Logf("Testing index map")
	idxMap := modConf.IndexMap()
    for _, pt := range ptList {
	    t.Logf("  %v -> %2d", pt, idxMap[pt.X][pt.Y])
	}
    t.Logf("Testing coordinate map")
    coordMap := modConf.CoordMap()
    for _, idx := range idxList {
        t.Logf("  %2d -> %v", idx, coordMap[idx])
    }
}

// func TestCoordMap(t *testing.T) {
// 	modConf := DefaultModuleConfig(image.Point{20, 20})
// 	t.Logf("module config: %v", modConf)
// 	idxMap := modConf.IndexMap()
//     coordMap := modConf.CoordMap()
// 	t.Logf("  idxMap: %v", idxMap)
// 	t.Logf("  coordMap: %v", coordMap)
// }

// func TestModuleConfig(t *testing.T) {
// 	var modConf ModuleConfig

//     modConf = DefaultModuleConfig(image.Point{30, 30})
// 	t.Logf("module configuration: %v", modConf)
// }

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
