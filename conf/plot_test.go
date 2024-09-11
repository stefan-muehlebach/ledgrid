package conf

import (
	"fmt"
	"image"
	"testing"
)

func TestPlotConfig(t *testing.T) {
	dim := image.Point{40, 40}
	fileName := fmt.Sprintf("plot%dx%d.png", dim.X, dim.Y)
	modConf := DefaultModuleConfig(dim)
	modConf.Plot(fileName)
}
