package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/stefan-muehlebach/ledgrid"
)

type colorType int

const (
	red colorType = iota
	green
	blue
)

const (
	defPort       = 5333
	defMissingIDs = ""
	defDefectIDs  = ""
	defBaud       = 2_000_000
	defUseTCP     = false
	defNumPix     = 400
)

func SignalHandler(gridServer *ledgrid.GridServer) {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGHUP, syscall.SIGUSR1)
	for sig := range sigChan {
		switch sig {
		case os.Interrupt:
			gridServer.Close()
			return

		case syscall.SIGHUP:
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

		case syscall.SIGUSR1:
			if gridServer.ToggleTestPattern() {
				log.Printf("Drawing test pattern is ON now.")
			} else {
				log.Printf("Drawing test pattern is OFF now.")
			}
		}
	}
}

func main() {
	var numPix int
	var port uint
	var baud int
	var missingIDs, defectIDs string
	var spiDevFile string = "/dev/spidev0.0"
	var spiBus ledgrid.Displayer
	var gridServer *ledgrid.GridServer

	// Verarbeite als erstes die Kommandozeilen-Optionen
	flag.IntVar(&numPix, "numpix", defNumPix, "Number of pixels (for fancy module configurations)")
	flag.UintVar(&port, "port", defPort, "UDP port")
	flag.IntVar(&baud, "baud", defBaud, "SPI baudrate in Hz")
	flag.StringVar(&missingIDs, "missing", defMissingIDs, "Comma separated list with IDs of missing LEDs (they will be skipped)")
	flag.StringVar(&defectIDs, "defect", defDefectIDs, "Comma separated list with IDs of defect LEDs (they will be blacked out)")
	flag.Parse()

	spiBus = ledgrid.NewSPIBus(spiDevFile, baud, numPix)
	gridServer = ledgrid.NewGridServer(port, spiBus)

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
	go SignalHandler(gridServer)

	gridServer.Handle()
}
