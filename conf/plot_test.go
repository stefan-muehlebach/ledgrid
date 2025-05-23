package conf

import (
	"fmt"
	"image"
	"testing"
)

func TestPlotDefaultConfig(t *testing.T) {
    dimList := []image.Point{
        {10, 10},
        {40, 10},
        {10, 40},
        {40, 40},
    }

    for _, dim := range dimList {
        	fileName := fmt.Sprintf("plots/default%dx%d.png", dim.X, dim.Y)
        	modConf := DefaultModuleConfig(dim)
        	modConf.Plot(fileName)
    }
}

var (
    customList = []string{
        "LR-000",
        "LR-090",
        "LR-180",
        "LR-270",
        "RL-000",
        "RL-090",
        "RL-180",
        "RL-270",
        "sample01",
        "sample02",
        "sample03",
        "sample04",
        "cap",
        "cup",
        "customA",
        "tetris",
        "squareWithHole",
        "lowerCurve",
        "chessBoard",
        "chessBoardSmall",
    }
)

func TestPlotCustomConfig(t *testing.T) {
    var modConf ModuleConfig
    var confFileName, pngFileName string

    for _, customName := range customList {
        confFileName = "data/" + customName + ".json"
        pngFileName = "plots/" + customName + ".png"
	    modConf = Load(confFileName)
        modConf.Plot(pngFileName)
    }
}
