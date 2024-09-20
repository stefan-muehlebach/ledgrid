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
	modConf   conf.ModuleConfig
	size      image.Point
	indexMap  conf.IndexMap
	coordMap  conf.CoordMap
	field     [][]*canvas.Circle
	numPixels int
}

// A new grid object must only know it's size in order to get the
// configuration of the emulated modules.
func NewGridWindowBySize(pixelSize float64, size image.Point) *GridWindow {
	modConf := conf.DefaultModuleConfig(size)
	return NewGridWindow(pixelSize, modConf)
}

func NewGridWindow(pixelSize float64, modConf conf.ModuleConfig) *GridWindow {
	e := &GridWindow{}
	e.Grid = container.NewWithoutLayout()
	e.modConf = modConf
	e.size = e.modConf.Size()
	e.coordMap = e.modConf.CoordMap()
	e.indexMap = e.modConf.IndexMap()
	e.field = make([][]*canvas.Circle, e.size.X)
	for i := range e.field {
		e.field[i] = make([]*canvas.Circle, e.size.Y)
	}
	ledSize := fyne.NewSize(float32(pixelSize-1), float32(pixelSize-1))
	for _, coord := range e.coordMap {
		col, row := coord.X, coord.Y
		ledPos := fyne.NewPos(float32(col)*float32(pixelSize),
			float32(row)*float32(pixelSize))
		led := canvas.NewCircle(color.RGBA{200, 200, 200, 255})
		led.Resize(ledSize)
		led.Move(ledPos)
		led.StrokeWidth = 0.0
		// led.StrokeColor = color.Black
		e.field[col][row] = led
		e.Grid.Add(led)
	}
	e.numPixels = e.size.X * e.size.Y
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
	return e.size.X * e.size.Y
}

// Takes the bytes in buffer, and uses them exactly as the real hardware
// would to recolor the individual LEDs (circle objects) of the emulation.
// Only if the colors really change, a fyne refresh is issued.
// Do not try to refresh the individual circles - this takes waaay to much
// time!
func (e *GridWindow) Send(buffer []byte) {
	var r, g, b uint8
	var i int
	var needsRefresh bool
	var newColor color.RGBA = color.RGBA{A: 0xff}
	var oldColor color.Color

	needsRefresh = false
	for i = 0; i < len(buffer); i += 3 {
		coord := e.coordMap[i/3]
		src := buffer[i : i+3 : i+3]
		r = src[0]
		g = src[1]
		b = src[2]
		newColor.R, newColor.G, newColor.B = r, g, b
		if !needsRefresh {
			oldColor = e.field[coord.X][coord.Y].FillColor
			if newColor != oldColor {
				needsRefresh = true
			}
		}
		e.field[coord.X][coord.Y].FillColor = newColor
	}
	if needsRefresh {
		e.Grid.Refresh()
	}
}

// func (e *GridWindow) Send(buffer []byte) {
// 	var r, g, b uint8
// 	var idx int
// 	var needsRefresh bool

// 	needsRefresh = false
// 	for i, val := range buffer {
// 		if i >= 3*e.numPixels {
// 			break
// 		}
// 		if i%3 == 0 {
// 			r = val
// 			idx = i / 3
// 		}
// 		if i%3 == 1 {
// 			g = val
// 		}
// 		if i%3 == 2 {
// 			b = val
// 			coord := e.coordMap[idx]
// 			newColor := color.RGBA{R: r, G: g, B: b, A: 0xff}
// 			if !needsRefresh {
// 				oldColor := e.field[coord.X][coord.Y].FillColor
// 				if newColor != oldColor {
// 					needsRefresh = true
// 				}
// 			}
// 			e.field[coord.X][coord.Y].FillColor = newColor
// 		}
// 	}
// 	if needsRefresh {
// 		e.Grid.Refresh()
// 	}
// }
