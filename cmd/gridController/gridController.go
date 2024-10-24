package main

import (
	"image"
	"github.com/stefan-muehlebach/ledgrid/conf"
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/stefan-muehlebach/ledgrid"
)

const (
	defWidth           = 40
	defHeight          = 10

	defMissingIDs = ""
	defDefectIDs  = ""
	defBaud       = 2_000_000
)

func SignalHandler(gridServer *ledgrid.GridServer) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGUSR1, syscall.SIGUSR2)
	for sig := range sigChan {
		switch sig {
		case os.Interrupt:
			gridServer.Close()
			return

		case syscall.SIGUSR1:
			log.Printf("Some server statistics:")
			log.Printf("   %v", gridServer.Watch())
			log.Printf("   %d bytes received by the controller", gridServer.RecvBytes)
			log.Printf("   %d bytes sent by the controller", gridServer.SentBytes)
			log.Printf("Current gamma values:")
			r, g, b := gridServer.Gamma()
			log.Printf("   R: %.1f, G: %.1f, B: %.1f", r, g, b)
			log.Printf("Current settings for max values (brightness):")
			br, bg, bb := gridServer.MaxBright()
			log.Printf("   R: %3d, G: %3d, B: %3d", br, bg, bb)
			gridServer.RecvBytes = 0
			gridServer.SentBytes = 0
			gridServer.Watch().Reset()

		case syscall.SIGUSR2:
			if gridServer.ToggleTestPattern() {
				log.Printf("Drawing test pattern is ON now.")
			} else {
				log.Printf("Drawing test pattern is OFF now.")
			}
		}
	}
}

func main() {
	var width, height int
	var customConfName string
	var gridSize image.Point
	var modConf conf.ModuleConfig

	var dataPort, rpcPort uint
	var baud int
	var missingIDs, defectIDs string
	var spiDevFile string = "/dev/spidev0.0"
	var ws2801 ledgrid.Displayer
	var gridServer *ledgrid.GridServer

	// Verarbeite als erstes die Kommandozeilen-Optionen
    	flag.IntVar(&width, "width", defWidth, "Width of panel")
	flag.IntVar(&height, "height", defHeight, "Height of panel")
	flag.StringVar(&customConfName, "custom", "", "Use a non standard module configuration")

	// flag.IntVar(&numPix, "numpix", defNumPix, "Number of pixels (for buffers and such)")
	flag.UintVar(&dataPort, "data", ledgrid.DefDataPort, "Data port (UPD as well as TCP)")
	flag.UintVar(&rpcPort, "rpc", ledgrid.DefRPCPort, "RPC port")
	flag.IntVar(&baud, "baud", defBaud, "SPI baudrate in Hz")
	flag.StringVar(&missingIDs, "missing", defMissingIDs, "Comma separated list with IDs of missing LEDs (they will be skipped)")
	flag.StringVar(&defectIDs, "defect", defDefectIDs, "Comma separated list with IDs of defect LEDs (they will be blacked out)")
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

	ws2801 = ledgrid.NewWS2801(spiDevFile, baud, modConf)
	gridServer = ledgrid.NewGridServer(dataPort, rpcPort, ws2801)

	if len(missingIDs) > 0 {
		for _, str := range strings.Split(missingIDs, ",") {
			val, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				log.Fatalf("Failed to parse 'missing': wrong format: %s", str)
			}
			gridServer.SetPixelStatus(int(val), ledgrid.LedMissing)
		}
	}

	if len(defectIDs) > 0 {
		for _, str := range strings.Split(defectIDs, ",") {
			val, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				log.Fatalf("Failed to parse 'defect': wrong format: %s", str)
			}
			gridServer.SetPixelStatus(int(val), ledgrid.LedDefect)
		}
	}

	// Damit der Daemon kontrolliert beendet werden kann, installieren wir
	// einen Handler fuer das INT-Signal, welches bspw. durch Ctrl-C erzeugt
	// wird oder auch von systemd beim Stoppen eines Services verwendet wird.
	gridServer.HandleEvents()
	SignalHandler(gridServer)
}
