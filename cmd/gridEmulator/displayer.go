package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

// Die Emulation des LedGrids als fyne-Objekt. Implementiert die Methoden
// des Displayer-Interfaces und kann daher GridServer direkt als Anzeigegeraet
// uebergeben werden.
type GridEmulator struct {
	Grid      *fyne.Container
	gridConf  conf.ModuleConfig
	coordMap  conf.CoordMap
	field     [][]*canvas.Circle
	size      image.Point
	numPixels int
}

func NewGridEmulator(size image.Point) *GridEmulator {
	e := &GridEmulator{size: size}
	e.Grid = container.NewGridWithRows(size.Y)
	e.gridConf = conf.DefaultModuleConfig(size)
	e.coordMap = e.gridConf.IndexMap().CoordMap()
	e.field = make([][]*canvas.Circle, size.X)
	for i := range e.field {
		e.field[i] = make([]*canvas.Circle, size.Y)
	}
	for col := range size.X {
		for row := range size.Y {
			ledColor := color.Black
			led := canvas.NewCircle(ledColor)
			led.StrokeWidth = 0.0
			e.field[col][row] = led
			e.Grid.Add(led)
		}
	}
	e.numPixels = size.X * size.Y
	return e
}

func (e *GridEmulator) DefaultGamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (e *GridEmulator) Close() {}

func (e *GridEmulator) Size() image.Point {
	return e.size
}

func (e *GridEmulator) Send(buffer []byte) {
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
