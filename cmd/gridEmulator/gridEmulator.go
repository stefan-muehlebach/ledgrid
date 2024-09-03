//go:generate fyne bundle -o icon.go Icon.ico

package main

import (
	"image"
	"flag"
	"fmt"
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

// func SignalHandler(gridServer *ledgrid.GridServer) {
// 	sigChan := make(chan os.Signal)
// 	signal.Notify(sigChan, os.Interrupt)
// 	for sig := range sigChan {
// 		switch sig {
// 		case os.Interrupt:
// 			gridServer.Close()
// 			return
// 		}
// 	}
// }

func SignalHandler(gridServer *ledgrid.GridServer) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP, syscall.SIGUSR1)
	for sig := range sigChan {
		switch sig {
		case os.Interrupt:
			gridServer.Close()
			return
		case syscall.SIGHUP:
			PrintStatistics(gridServer)
		case syscall.SIGUSR1:
			ToggleTests(gridServer)
		}
	}
}

func PrintStatistics(gridServer *ledgrid.GridServer) {
	log.Printf("Server Statistics:")
	log.Printf("   %v", gridServer.Watch())
	log.Printf("   %d bytes received by the controller", gridServer.RecvBytes)
	log.Printf("   %d bytes sent by the controller", gridServer.SentBytes)
	log.Printf("Current gamma values:")
	r, g, b := gridServer.Gamma()
	log.Printf("   R: %.1f, G: %.1f, B: %.1f", r, g, b)
	log.Printf("Current settings for max values (brightness):")
	br, bg, bb := gridServer.MaxBright()
	log.Printf("   R: %3d, G: %3d, B: %3d", br, bg, bb)
}

func ToggleTests(gridServer *ledgrid.GridServer) {
	if gridServer.ToggleTestPattern() {
		log.Printf("Drawing test pattern is ON now.")
	} else {
		log.Printf("Drawing test pattern is OFF now.")
	}
}

func main() {
	var width, height int
	var port uint
	var appWidth, appHeight float32
	var pixelSize float64
	var appSize fyne.Size
	var gridServer *ledgrid.GridServer
	var gridEmulator *GridEmulator

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
	winTitle := fmt.Sprintf("LEDGrid Emulator (Size: %d x %d; Port: %d)", width, height, port)
	Win = App.NewWindow(winTitle)

	gridEmulator = NewGridEmulator(image.Pt(width, height))
	gridServer = ledgrid.NewGridServer(port, gridEmulator)

	Win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyEscape, fyne.KeyQ:
			App.Quit()
		case fyne.KeyH:
			fmt.Printf("Use the following keys to control the software:\n")
			fmt.Printf("  h   Show this help (again)\n")
			fmt.Printf("  s   Show some statistics\n")
			fmt.Printf("  t   Start test pattern, press 't' again to stop\n")
			fmt.Printf("  q   Quit the program\n")
			fmt.Printf(" ESC  Same as 'q'\n")
		case fyne.KeyS:
			PrintStatistics(gridServer)
		case fyne.KeyT:
			ToggleTests(gridServer)
		}
	})

	go SignalHandler(gridServer)
	go gridServer.Handle()

	Win.SetContent(gridEmulator.Grid)
	Win.Resize(appSize)
	Win.ShowAndRun()
}