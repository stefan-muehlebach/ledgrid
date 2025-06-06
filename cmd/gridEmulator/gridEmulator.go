package main

import (
	"flag"
	"fmt"
	"image"
	"log"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

const (
	defWidth           = 40
	defHeight          = 10
	defPixelSize       = 40.0
)

func ResetStatistics(gridServer *ledgrid.GridServer) {
	gridServer.Stopwatch().Reset()
	gridServer.RecvBytes = 0
	gridServer.SentBytes = 0
}

func PrintStatistics(gridServer *ledgrid.GridServer) {
	log.Printf("Emulator statistics:")
	log.Printf("   %v", gridServer.Stopwatch())
	log.Printf("   %d bytes received by the controller", gridServer.RecvBytes)
	log.Printf("   %d bytes sent by the controller", gridServer.SentBytes)
	log.Printf("Current gamma values:")
	r, g, b := gridServer.Gamma()
	log.Printf("   R: %.1f, G: %.1f, B: %.1f", r, g, b)
	log.Printf("Current settings for max values (brightness):")
}

func ToggleTests(gridServer *ledgrid.GridServer) {
	if gridServer.ToggleTestPattern() {
		log.Printf("Drawing test pattern is ON now.")
	} else {
		log.Printf("Drawing test pattern is OFF now.")
	}
}

var (
    gridServer *ledgrid.GridServer
)

func main() {
	var width, height int
	var dataPort, rpcPort uint
	var pixelSize float64
	var gridWindow *Window
	var customConfName string
	var gridSize image.Point
	var modConf conf.ModuleConfig

	flag.IntVar(&width, "width", defWidth, "Width of panel")
	flag.IntVar(&height, "height", defHeight, "Height of panel")
	flag.UintVar(&dataPort, "data", ledgrid.DefDataPort, "Data port")
	flag.UintVar(&rpcPort, "rpc", ledgrid.DefRPCPort, "RPC port")
	flag.Float64Var(&pixelSize, "size", defPixelSize, "Diameter of one LED in pixels")
	flag.StringVar(&customConfName, "custom", "", "Use a non standard module configuration")
	flag.Parse()

    StartProfiling()
    defer StopProfiling()

	if customConfName != "" {
		modConf = conf.Load("data/" + customConfName + ".json")
		gridSize = modConf.Size()
		width, height = gridSize.X, gridSize.Y
	} else {
		gridSize = image.Point{width, height}
		modConf = conf.DefaultModuleConfig(gridSize)
	}

    title := fmt.Sprintf("LEDGrid Emulator (Size: %d x %d; Port: %d)", gridSize.X, gridSize.Y, dataPort)

	gridWindow = NewWindow(title, pixelSize, modConf)
	gridServer = ledgrid.NewGridServer(dataPort, rpcPort, gridWindow)

    gridServer.HandleEvents()
    gridWindow.HandleEvents()
}
