package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

// Die Emulation des LedGrids als fyne-Applikation
type PixelEmulator struct {
	Grid      *fyne.Container
	gridConf  conf.ModuleConfig
	coordMap  conf.CoordMap
	field     [][]*canvas.Circle
	size      image.Point
	numPixels int
}

func NewPixelEmulator(width, height int) *PixelEmulator {
	e := &PixelEmulator{}
	e.Grid = container.NewGridWithRows(height)
	e.gridConf = conf.DefaultModuleConfig(image.Point{width, height})
	e.coordMap = e.gridConf.IndexMap().CoordMap()
	e.field = make([][]*canvas.Circle, width)
	for i := range e.field {
		e.field[i] = make([]*canvas.Circle, height)
	}
	for col := range width {
		for row := range height {
			ledColor := color.Black
			led := canvas.NewCircle(ledColor)
			led.StrokeWidth = 0.0
			e.field[col][row] = led
			e.Grid.Add(led)
		}
	}
	e.size = image.Point{width, height}
	e.numPixels = width * height
	return e
}

func (e *PixelEmulator) DefaultGamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (e *PixelEmulator) Close() {}

func (e *PixelEmulator) Size() image.Point {
	return e.size
}

func (e *PixelEmulator) Send(buffer []byte) {
	var r, g, b uint8
	var idx int
	var src []byte
	var needsRefresh bool

	src = buffer
	needsRefresh = false
	for i, val := range src {
		if i >= 3*e.numPixels {
			break
		}
		if i%3 == 0 {
			r = val
			idx = i / 3
		}
		if i%3 == 1 {
			g = val
		}
		if i%3 == 2 {
			b = val
			coord := e.coordMap[idx]
			newColor := color.RGBA{R: r, G: g, B: b, A: 0xff}
			if !needsRefresh {
				oldColor := e.field[coord.X][coord.Y].FillColor
				if newColor != oldColor {
					needsRefresh = true
				}
			}
			e.field[coord.X][coord.Y].FillColor = newColor
		}
	}
	if needsRefresh {
		e.Grid.Refresh()
	}
}
