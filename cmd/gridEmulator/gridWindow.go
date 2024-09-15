package main

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

// This is the essential part of the gridEmulator: GridWindow, which
// implements the ledgrid.Displayer interface and can be used by GridServer
// as an output object. Sadly, I wasn't able to implement it as a fyne object
// (widget or container) as it should be...
// It does not only show a grid of colorful LEDs, but exactly emulates the
// configuration of the original LedGrid using 10x10 modules as well as the
// snake cabeling.
type GridWindow struct {
    // This is the fyne object, which must be added to a fyne application in
    // order to experience the whole glory of the emulation.
	Grid      *fyne.Container
	modConf  conf.ModuleConfig
	coordMap  conf.CoordMap
	field     [][]*canvas.Circle
	size      image.Point
	numPixels int
}

// A new grid object must only know it's size in order to get the
// configuration of the emulated modules.
func NewGridWindowBySize(size image.Point) *GridWindow {
    modConf := conf.DefaultModuleConfig(size)
    return NewGridWindow(size, modConf)
}

func NewGridWindow(size image.Point, modConf conf.ModuleConfig) *GridWindow {
	e := &GridWindow{size: size}
	e.Grid = container.NewGridWithRows(size.Y)
    	e.modConf = modConf
	e.coordMap = e.modConf.CoordMap()
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

// Since we implement ledgrid.Displayer, we must provide a default gamma
// setting. In contrast to the real hardware, the emulation must not correct
// any colors, so the gamma values are 1.0 for all colors.
func (e *GridWindow) DefaultGamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (e *GridWindow) Close() {}

func (e *GridWindow) Size() int {
	return e.size.X*e.size.Y
}

// Takes the bytes in buffer, and uses them exactly as the real hardware
// would to recolor the individual LEDs (circle objects) of the emulation.
// Only if the colors really change, a fyne refresh is issued.
// Do not try to refresh the individual circles - this takes waaay to much
// time!
func (e *GridWindow) Send(buffer []byte) {
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
