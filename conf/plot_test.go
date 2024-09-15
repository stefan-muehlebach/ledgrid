package conf

import (
	"fmt"
	"image"
	"testing"
)

func TestPlotDefaultConfig(t *testing.T) {
	dim := image.Point{40, 40}
	fileName := fmt.Sprintf("plot%dx%d.png", dim.X, dim.Y)
	modConf := DefaultModuleConfig(dim)
	modConf.Plot(fileName)
}

func TestPlotCustomConfig(t *testing.T) {
	dim := image.Point{30, 20}
	fileName := fmt.Sprintf("plot%dx%d.png", dim.X, dim.Y)
	modConf := ModuleConfig{
        ModulePosition{Col: 0, Row: 0, Idx: 0, Mod: Module{ModLR, Rot000}},
        ModulePosition{Col: 1, Row: 0, Idx: 100, Mod: Module{ModRL, Rot090}},
        ModulePosition{Col: 1, Row: 1, Idx: 200, Mod: Module{ModLR, Rot000}},
        ModulePosition{Col: 2, Row: 1, Idx: 300, Mod: Module{ModLR, Rot000}},
    }
	modConf.Plot(fileName)
}
