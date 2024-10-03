package conf

import (
	"fmt"
	"image"
	"testing"
)

func TestPlotDefaultConfig(t *testing.T) {
	dim := image.Point{40, 10}
	fileName := fmt.Sprintf("plots/default%dx%d.png", dim.X, dim.Y)
	modConf := DefaultModuleConfig(dim)
	modConf.Plot(fileName)

	dim = image.Point{40, 40}
	fileName = fmt.Sprintf("plots/default%dx%d.png", dim.X, dim.Y)
	modConf = DefaultModuleConfig(dim)
	modConf.Plot(fileName)

	dim = image.Point{10, 40}
	fileName = fmt.Sprintf("plots/default%dx%d.png", dim.X, dim.Y)
	modConf = DefaultModuleConfig(dim)
	modConf.Plot(fileName)
}

func TestPlotCustomConfig(t *testing.T) {
    var modConf ModuleConfig
    var fileName string

	// modConf = TetrisTile
	// fileName = fmt.Sprintf("plotTetrisTile.png")
	// modConf.Plot(path.Join(dirName, fileName))

	modConf.Load("data/tetris.json")
	fileName = "plots/tetris.png"
	modConf.Plot(fileName)

	modConf.Load("data/squareWithHole.json")
	fileName = "plots/squareWithHole.png"
	modConf.Plot(fileName)

	modConf.Load("data/lowerCurve.json")
	fileName = "plots/lowerCurve.png"
	modConf.Plot(fileName)

	modConf.Load("data/lowerCurve.json")
	fileName = "plots/lowerCurve.png"
	modConf.Plot(fileName)

	modConf.Load("data/chessBoardSmall.json")
	fileName = "plots/chessBoardSmall.png"
	modConf.Plot(fileName)

	modConf.Load("data/chessBoard.json")
	fileName = "plots/chessBoard.png"
	modConf.Plot(fileName)
}
