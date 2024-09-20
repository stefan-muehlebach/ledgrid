package conf

import (
	"fmt"
	"image"
	"testing"
)


func TestPlotDefaultConfig(t *testing.T) {
	dim := image.Point{40, 10}
	fileName := fmt.Sprintf("plot%dx%d.png", dim.X, dim.Y)
	modConf := DefaultModuleConfig(dim)
	modConf.Plot(fileName)

	dim = image.Point{40, 40}
	fileName = fmt.Sprintf("plot%dx%d.png", dim.X, dim.Y)
	modConf = DefaultModuleConfig(dim)
	modConf.Plot(fileName)
}

func TestPlotCustomConfig(t *testing.T) {
    var modConf ModuleConfig
    var fileName string

	modConf = TetrisTile
	fileName = fmt.Sprintf("plotTetrisTile.png")
	modConf.Plot(fileName)

	modConf = SquareWithHole
	fileName = fmt.Sprintf("plotSquareWithHole.png")
	modConf.Plot(fileName)

	modConf = LowerCurve
	fileName = fmt.Sprintf("plotLowerCurve.png")
	modConf.Plot(fileName)

	modConf = ChessBoard
	fileName = fmt.Sprintf("plotChessBoard.png")
	modConf.Plot(fileName)

	modConf = SmallChessBoard
	fileName = fmt.Sprintf("plotSmallChessBoard.png")
	modConf.Plot(fileName)
}
