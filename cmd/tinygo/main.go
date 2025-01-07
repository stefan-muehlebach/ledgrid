package main

import (
	"fmt"
	"image"
	"math"
	"os"
	"os/signal"
	"time"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/color"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

const (
	defWidth  = 10
	defHeight = 10
)

var (
	width      int = defWidth
	height     int = defHeight
	gridSize   image.Point
	gridClient ledgrid.GridClient
	ledGrid    *ledgrid.LedGrid
	animCtrl   *ledgrid.AnimationController
	canvas     *ledgrid.Canvas
)

// ---------------------------------------------------------------------------

func SignalHandler(timeout time.Duration) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	if timeout == 0 {
		timeout = time.Duration(math.MaxInt64)
	}
	timer := time.NewTimer(timeout)
	select {
	case <-sigChan:
	case <-timer.C:
	}
}

//----------------------------------------------------------------------------

func main() {
	var modConf conf.ModuleConfig
	var timeout time.Duration
	// var gR, gG, gB float64

	gridClient = ledgrid.NewDirectGridClient(ledgrid.NewWS2801avr(conf.DefaultModuleConfig(image.Point{width, height})))
	modConf = gridClient.ModuleConfig()
	ledGrid = ledgrid.NewLedGrid(gridClient, modConf)
	// gR, gG, gB = ledGrid.Client.Gamma()

	gridSize = ledGrid.Rect.Size()
	width = gridSize.X
	height = gridSize.Y

	canvas = ledGrid.Canvas(0)
	animCtrl = ledGrid.AnimCtrl

	ledGrid.StartRefresh()

	GlowingPixels(canvas)
	SignalHandler(timeout)

	ledgrid.AnimCtrl.Suspend()
	ledGrid.Clear(color.Black)
	ledGrid.Close()

	fmt.Printf("Program statistics:\n")
	fmt.Printf("  animation: %v\n", ledgrid.AnimCtrl.Watch())
	fmt.Printf("  painting : %v\n", canvas.Watch())
	fmt.Printf("  sending  : %v\n", ledGrid.Client.Watch())
}
