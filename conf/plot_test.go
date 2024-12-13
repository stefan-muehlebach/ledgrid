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

func TestPlotCustomConfig(t *testing.T) {
    var modConf ModuleConfig
    var fileName string

	modConf = Load("data/cap.json")
	fileName = "plots/cap.png"
	modConf.Plot(fileName)

	modConf = Load("data/cup.json")
	fileName = "plots/cup.png"
	modConf.Plot(fileName)

	modConf = Load("data/customConf.json")
	fileName = "plots/customConf.png"
	modConf.Plot(fileName)

	modConf = Load("data/tetris.json")
	fileName = "plots/tetris.png"
	modConf.Plot(fileName)

	modConf = Load("data/squareWithHole.json")
	fileName = "plots/squareWithHole.png"
	modConf.Plot(fileName)

	modConf = Load("data/lowerCurve.json")
	fileName = "plots/lowerCurve.png"
	modConf.Plot(fileName)

	modConf = Load("data/chessBoardSmall.json")
	fileName = "plots/chessBoardSmall.png"
	modConf.Plot(fileName)

	modConf = Load("data/chessBoard.json")
	fileName = "plots/chessBoard.png"
	modConf.Plot(fileName)
}
