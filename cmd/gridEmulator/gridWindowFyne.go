//go:build guiFyne

package main

import (
	"fmt"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

// Since the default padding between adjacent elements in a GridContainer is
// waaaay to large, we had to define a custom theme with only one divergent
// property: theme.SizeNamePadding
type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}
func (t myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}
func (t myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}
func (t myTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNamePadding {
		return 1.0
	}
	return theme.DefaultTheme().Size(name)
}

// var (
// 	App fyne.App
// 	Win fyne.Window
// )

// This is the essential part of the gridEmulator: GridWindow, which
// implements the ledgrid.Displayer interface and can be used by GridServer
// as an output object. Sadly, I wasn't able to implement it as a fyne object
// (widget or container) as it should be...
// It does not only show a grid of colorful LEDs, but exactly emulates the
// configuration of the original LedGrid using 10x10 modules as well as the
// snake cabeling.
type Window struct {
	ledgrid.DisplayEmbed
	// This is the fyne object, which must be added to a fyne application in
	// order to experience the whole glory of the emulation.
	App       fyne.App
	Win       fyne.Window
	Grid      *fyne.Container
	size      image.Point
	indexMap  conf.IndexMap
	coordMap  conf.CoordMap
	field     [][]*canvas.Circle
	numPixels int
}

// A new grid object must only know it's size in order to get the
// configuration of the emulated modules.
func NewWindowBySize(title string, pixelSize float64, size image.Point) *Window {
	modConf := conf.DefaultModuleConfig(size)
	return NewWindow(title, pixelSize, modConf)
}

func NewWindow(title string, pixelSize float64, modConf conf.ModuleConfig) *Window {
	e := &Window{}
	e.DisplayEmbed.Init(e, len(modConf)*conf.ModuleDim.X*conf.ModuleDim.Y)

	e.App = app.New()
	e.App.SetIcon(resourceIconIco)
	e.App.Settings().SetTheme(&myTheme{})
	e.Win = e.App.NewWindow(title)

	e.Grid = container.NewWithoutLayout()
	e.ModConf = modConf
	e.size = e.ModConf.Size()
	e.coordMap = e.ModConf.CoordMap()
	e.indexMap = e.ModConf.IndexMap()
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

	winSize := fyne.NewSize(float32(e.size.X)*float32(pixelSize), float32(e.size.Y)*float32(pixelSize))
	e.Win.SetContent(e.Grid)
	e.Win.Resize(winSize)
	e.Win.SetFixedSize(true)

	return e
}

func (e *Window) HandleEvents() {
	e.Win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyH:
			fmt.Printf("Use the following keys to control the software:\n")
			fmt.Printf("  h   Show this help (again)\n")
			fmt.Printf("  s   Show some statistics\n")
			fmt.Printf("  t   Start test pattern, press 't' again to stop\n")
			fmt.Printf("  q   Quit the program\n")
			fmt.Printf(" ESC  Same as 'q'\n")
		// case fyne.KeyS:
		// 	PrintStatistics(gridServer)
		// 	ResetStatistics(gridServer)
		// case fyne.KeyT:
		// 	ToggleTests(gridServer)
		case fyne.KeyEscape, fyne.KeyQ:
			e.App.Quit()
		}
	})
	e.Win.ShowAndRun()
}

// Since we implement ledgrid.Displayer, we must provide a default gamma
// setting. In contrast to the real hardware, the emulation must not correct
// any colors, so the gamma values are 1.0 for all colors.
func (e *Window) DefaultGamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (e *Window) Close() {}

// Takes the bytes in buffer, and uses them exactly as the real hardware
// would to recolor the individual LEDs (circle objects) of the emulation.
// Only if the colors really change, a fyne refresh is issued.
// Do not try to refresh the individual circles - this takes waaay to much
// time!
func (e *Window) Send(buffer []byte) {
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
