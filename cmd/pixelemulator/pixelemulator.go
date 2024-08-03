package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"os"
	"os/signal"
	"syscall"

	"fyne.io/fyne/v2/container"

	"github.com/stefan-muehlebach/ledgrid"

	"fyne.io/fyne/v2/canvas"
	_ "github.com/stefan-muehlebach/ledgrid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "fyne.io/fyne/v2/driver/desktop"
)

const (
	defPort        = 5333
	defGammaValues = "1.0,1.0,1.0"
	defWidth       = 10
	defHeight      = 10
	defPixelSize   = 50.0
)

var (
	App fyne.App
	Win fyne.Window
)

type PixelEmulator struct {
	grid     *fyne.Container
	gridConf ledgrid.ModuleConfig
	coordMap ledgrid.CoordMap
	field    [][]*canvas.Circle
}

func NewPixelEmulator(width, height int) *PixelEmulator {
	e := &PixelEmulator{}
	e.grid = container.NewGridWithRows(height)
	e.gridConf = ledgrid.DefaultModuleConfig(image.Point{width, height})
	e.coordMap = e.gridConf.IndexMap().CoordMap()
	e.field = make([][]*canvas.Circle, width)
	for i := range e.field {
		e.field[i] = make([]*canvas.Circle, height)
	}
	for col := range width {
		for row := range height {
			ledColor := color.White
			led := canvas.NewCircle(ledColor)
			e.field[col][row] = led
			e.grid.Add(led)
		}
	}
	return e
}

func (e *PixelEmulator) Close() {}

func (e *PixelEmulator) Send(buffer []byte) {
	var r, g, b uint8
	var idx int
	var src []byte
    var needsRefresh bool

	src = buffer
    needsRefresh = false
	for i, val := range src {
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
	    e.grid.Refresh()
    }
}

func SignalHandler(pixelServer *ledgrid.PixelServer) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP, syscall.SIGUSR1)
	for sig := range sigChan {
		switch sig {
		case os.Interrupt:
			pixelServer.Close()
			return
		case syscall.SIGHUP:
			log.Printf("Server Statistics:")
			num, total, avg := pixelServer.SendWatch.Stats()
			log.Printf("   %d sends to fyne.io took %v (%v per send)", num, total, avg)
			log.Printf("   %d bytes received by the controller", pixelServer.RecvBytes)
			log.Printf("   %d bytes sent by the controller", pixelServer.SentBytes)
		case syscall.SIGUSR1:
			if pixelServer.ToggleTestPattern() {
				log.Printf("Drawing test pattern is ON now.")
			} else {
				log.Printf("Drawing test pattern is OFF now.")
			}
		}
	}
}

func main() {
	var width, height int
	var port uint
	var appWidth, appHeight float32
    var pixelSize float64
	var appSize fyne.Size
	var pixelServer *ledgrid.PixelServer
	var pixelEmulator *PixelEmulator

	flag.IntVar(&width, "width", defWidth, "Width of panel")
	flag.IntVar(&height, "height", defHeight, "Height of panel")
	flag.UintVar(&port, "port", defPort, "UDP port")
    flag.Float64Var(&pixelSize, "size", defPixelSize, "Size of one LED")
	flag.Parse()

	appWidth = float32(width) * float32(pixelSize)
	appHeight = float32(height) * float32(pixelSize)
	appSize = fyne.NewSize(appWidth, appHeight)

	App = app.New()
	Win = App.NewWindow("LedGrid Emulator")

	pixelServer = ledgrid.NewPixelServer(port)
	pixelEmulator = NewPixelEmulator(width, height)
	pixelServer.Disp = pixelEmulator
    pixelServer.SetGamma(1.0, 1.0, 1.0)

	Win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyEscape, fyne.KeyQ:
			App.Quit()
		}
	})

	go SignalHandler(pixelServer)
	go pixelServer.Handle()

	Win.SetContent(pixelEmulator.grid)
	Win.Resize(appSize)
	Win.ShowAndRun()
}
