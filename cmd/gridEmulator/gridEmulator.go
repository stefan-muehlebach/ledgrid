//go:generate fyne bundle -o icon.go Icon.ico

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/conf"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
)

const (
	defWidth           = 40
	defHeight          = 10
	defPixelSize       = 50.0
	defUseCustomLayout = false
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

var (
	App fyne.App
	Win fyne.Window
)

func ResetStatistics(gridServer *ledgrid.GridServer) {
	gridServer.Watch().Reset()
	gridServer.RecvBytes = 0
	gridServer.SentBytes = 0
}

func PrintStatistics(gridServer *ledgrid.GridServer) {
	log.Printf("Emulator statistics:")
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
	var dataPort, rpcPort uint
	var appWidth, appHeight float32
	var pixelSize float64
	var appSize fyne.Size
	var gridServer *ledgrid.GridServer
	var gridEmulator *GridWindow
	var customConfName string
	var gridSize image.Point
	var modConf conf.ModuleConfig

	flag.IntVar(&width, "width", defWidth, "Width of panel")
	flag.IntVar(&height, "height", defHeight, "Height of panel")
	flag.UintVar(&dataPort, "port", ledgrid.DefDataPort, "Data port")
	flag.UintVar(&rpcPort, "rpc", ledgrid.DefRPCPort, "RPC port")
	flag.Float64Var(&pixelSize, "size", defPixelSize, "Size of one LED")
	flag.StringVar(&customConfName, "custom", "", "Use a non standard module configuration")
	flag.Parse()

    StartProfiling()
    defer StopProfiling()

	if customConfName != "" {
		modConf = conf.Load("data/" + customConfName + ".json")
		gridSize = modConf.Size()
		width, height = gridSize.X, gridSize.Y
	} else {
		gridSize = image.Pt(width, height)
		modConf = conf.DefaultModuleConfig(gridSize)
	}

	appWidth = float32(width) * float32(pixelSize)
	appHeight = float32(height) * float32(pixelSize)
	appSize = fyne.NewSize(appWidth, appHeight)

	App = app.New()
	App.SetIcon(resourceIconIco)
	App.Settings().SetTheme(&myTheme{})
	winTitle := fmt.Sprintf("LEDGrid Emulator (Size: %d x %d; Port: %d)", width, height, dataPort)
	Win = App.NewWindow(winTitle)

	gridEmulator = NewGridWindow(pixelSize, modConf)
	gridServer = ledgrid.NewGridServer(dataPort, rpcPort, gridEmulator)

	Win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyH:
			fmt.Printf("Use the following keys to control the software:\n")
			fmt.Printf("  h   Show this help (again)\n")
			fmt.Printf("  s   Show some statistics\n")
			fmt.Printf("  t   Start test pattern, press 't' again to stop\n")
			fmt.Printf("  q   Quit the program\n")
			fmt.Printf(" ESC  Same as 'q'\n")
		case fyne.KeyS:
			PrintStatistics(gridServer)
			ResetStatistics(gridServer)
		case fyne.KeyT:
			ToggleTests(gridServer)
		case fyne.KeyEscape, fyne.KeyQ:
			App.Quit()
		}
	})

	Win.SetContent(gridEmulator.Grid)
	Win.Resize(appSize)
    Win.SetFixedSize(true)
	Win.ShowAndRun()
}
