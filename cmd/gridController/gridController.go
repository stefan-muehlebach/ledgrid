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
	defPort        = 5333
	defMissingIDs  = ""
	defDefectIDs   = ""
	defBaud        = 2_000_000
	defUseTCP      = false
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
	var port uint
	var baud int
	var missingIDs, defectIDs string
	var spiDevFile string = "/dev/spidev0.0"
    var spiBus ledgrid.Displayer
	var pixelServer *ledgrid.PixelServer

	// Verarbeite als erstes die Kommandozeilen-Optionen
	flag.UintVar(&port, "port", defPort, "UDP port")
	flag.IntVar(&baud, "baud", defBaud, "SPI baudrate in Hz")
	flag.StringVar(&missingIDs, "missing", defMissingIDs, "Comma separated list with IDs of missing LEDs")
	flag.StringVar(&defectIDs, "defect", defDefectIDs, "Comma separated list with IDs of LEDs to black out")
	flag.Parse()

    spiBus = ledgrid.NewSPIBus(spiDevFile, baud)
	pixelServer = ledgrid.NewPixelServer(port, spiBus)

	if len(missingIDs) > 0 {
		for _, str := range strings.Split(missingIDs, ",") {
			val, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				log.Fatalf("Failed to parse 'missing': wrong format: %s", str)
			}
			pixelServer.SetPixelStatus(int(val), ledgrid.PixelMissing)
		}
	}

	if len(defectIDs) > 0 {
		for _, str := range strings.Split(defectIDs, ",") {
			val, err := strconv.ParseInt(str, 10, 32)
			if err != nil {
				log.Fatalf("Failed to parse 'defect': wrong format: %s", str)
			}
			pixelServer.SetPixelStatus(int(val), ledgrid.PixelDefect)
		}
	}

	// Damit der Daemon kontrolliert beendet werden kann, installieren wir
	// einen Handler fuer das INT-Signal, welches bspw. durch Ctrl-C erzeugt
	// wird oder auch von systemd beim Stoppen eines Services verwendet wird.
	go SignalHandler(pixelServer)

	pixelServer.Handle()
}
