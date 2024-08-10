package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stefan-muehlebach/ledgrid"

	_ "github.com/stefan-muehlebach/ledgrid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "fyne.io/fyne/v2/driver/desktop"
)

const (
	defPort        = 5333
	defWidth       = 40
	defHeight      = 10
	defPixelSize   = 50.0
)

var (
	App fyne.App
	Win fyne.Window
)

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

	pixelEmulator = NewPixelEmulator(width, height)
	pixelServer = ledgrid.NewPixelServer(port, pixelEmulator)

	Win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyEscape, fyne.KeyQ:
			App.Quit()
		}
	})

	go SignalHandler(pixelServer)
	go pixelServer.Handle()

	Win.SetContent(pixelEmulator.Grid)
	Win.Resize(appSize)
	Win.ShowAndRun()
}
