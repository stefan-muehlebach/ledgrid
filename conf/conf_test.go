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

	goodConf01 = ModuleConfig{
		ModulePosition{Col: 0, Row: 1, Idx: 0,   Mod: ModLR000},
		ModulePosition{Col: 0, Row: 0, Idx: 100, Mod: ModLR180},
	}
	goodConf02 = ModuleConfig{
		ModulePosition{Col: 0, Row: 1, Idx: 0,   Mod: ModLR000},
		ModulePosition{Col: 1, Row: 0, Idx: 100, Mod: ModRL180},
	}
	goodConf03 = ModuleConfig{
		ModulePosition{Col: 0, Row: 0, Idx: 0,   Mod: ModRL180},
		ModulePosition{Col: 1, Row: 1, Idx: 100, Mod: ModLR000},
		ModulePosition{Col: 2, Row: 0, Idx: 200, Mod: ModRL180},
		ModulePosition{Col: 3, Row: 1, Idx: 300, Mod: ModLR000},
	}

	badConf01 = ModuleConfig{
		ModulePosition{Col: 0, Row: 0, Idx: 0,   Mod: ModLR000},
		ModulePosition{Col: 1, Row: 0, Idx: 100, Mod: ModLR180},
	}
	badConf02 = ModuleConfig{
		ModulePosition{Col: 0, Row: 0, Idx: 0,   Mod: ModLR000},
		ModulePosition{Col: 1, Row: 0, Idx: 100, Mod: ModRL000},
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

func TestVerify(t *testing.T) {
    t.Logf("Verify Default Configuration")
    modConf := DefaultModuleConfig(image.Point{30, 30})
    err := modConf.Verify()
    if err != nil {
        t.Error(err)
    }

    t.Logf("Verify Custom Configuration (Tetris)")
    modConf = TetrisTile
    err = modConf.Verify()
    if err != nil {
        t.Error(err)
    }

    t.Logf("Verify Custom Configuration (LowerCurve)")
    modConf = LowerCurve
    err = modConf.Verify()
    if err != nil {
        t.Error(err)
    }

    t.Logf("Verify Custom Configuration (SquareWithHole)")
    modConf = SquareWithHole
    err = modConf.Verify()
    if err != nil {
        t.Error(err)
    }

    t.Logf("Verify Good Config 01")
    modConf = goodConf01
    err = modConf.Verify()
    if err != nil {
        t.Error(err)
    }

    t.Logf("Verify Good Config 02")
    modConf = goodConf02
    err = modConf.Verify()
    if err != nil {
        t.Error(err)
    }

    t.Logf("Verify Good Config 03")
    modConf = goodConf03
    err = modConf.Verify()
    if err != nil {
        t.Error(err)
    }

    t.Logf("Verify Bad Config 01")
    modConf = badConf01
    err = modConf.Verify()
    if err != nil {
        t.Error(err)
    }

    t.Logf("Verify Bad Config 02")
    modConf = badConf02
    err = modConf.Verify()
    if err != nil {
        t.Error(err)
    }
}

func TestSaveCustomConf(t *testing.T) {
    t.Logf("Save custom configuration")
//     TetrisTile.Save("tetris.json")
//     LowerCurve.Save("lowerCurve.json")
//     SquareWithHole.Save("squareWithHole.json")
//     SmallChessBoard.Save("chessBoardSmall.json")
//     ChessBoard.Save("chessBoard.json")
    CustomConf.Save("data/customConf.json")
}


func TestLoadCustomConf(t *testing.T) {
    var conf ModuleConfig

    t.Logf("Load custom configuration")
    conf.Load("data/squareWithHole.json")

    t.Logf("%v", conf)
}


