//go:generate fyne bundle -o icon.go Icon.ico

package main

import (
	"fmt"
	"flag"
	"image/color"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stefan-muehlebach/ledgrid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

const (
	defPort      = 5333
	defWidth     = 40
	defHeight    = 10
	defPixelSize = 50.0
)

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
			log.Printf("   %v", pixelServer.Watch())
			log.Printf("   %d bytes received by the controller", pixelServer.RecvBytes)
			log.Printf("   %d bytes sent by the controller", pixelServer.SentBytes)
			log.Printf("Current gamma values:")
			r, g, b := pixelServer.Gamma()
			log.Printf("   R: %.1f, G: %.1f, B: %.1f", r, g, b)
			log.Printf("Current settings for max values (brightness):")
			br, bg, bb := pixelServer.MaxBright()
			log.Printf("   R: %3d, G: %3d, B: %3d", br, bg, bb)
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
	App.SetIcon(resourceIconIco)
	App.Settings().SetTheme(&myTheme{})
    winTitle := fmt.Sprintf("LEDGrid Emulator (%d x %d)", width, height)
	Win = App.NewWindow(winTitle)

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
