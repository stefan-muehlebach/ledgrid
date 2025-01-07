package main

import (
	"errors"
	"flag"
	"image"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/conf"
)

const (
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
			log.Printf("Server statistics:")
			log.Printf("   %v", gridServer.Watch())
			log.Printf("   %v bytes received by the controller", gridServer.RecvBytes)
			log.Printf("   %v bytes sent by the controller", gridServer.SentBytes)
			log.Printf("Current gamma values:")
			r, g, b := gridServer.Gamma()
			log.Printf("   R: %.1f, G: %.1f, B: %.1f", r, g, b)
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

func PlayFile(fileName string) {
	var client ledgrid.GridClient
	var buffer []byte

	client = ledgrid.NewNetGridClient("localhost", "udp", dataPort, rpcPort)
	buffer = make([]byte, 3*client.NumLeds())

	fh, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Couldn't open move file: %v", err)
	}
	ticker := time.NewTicker(30 * time.Millisecond)
	for range ticker.C {
		n, err := fh.Read(buffer)
		if n == 0 && errors.Is(err, io.EOF) {
			break
		}
		client.Send(buffer)
	}
	fh.Close()
    client.Close()
}

var (
	dataPort, rpcPort uint
)

func main() {
	var inFile string

	var width, height int
	var customConfName string
	var gridSize image.Point
	var modConf conf.ModuleConfig

	var baud int
	var missingIDs, defectIDs string
	var spiDevFile string = "/dev/spidev0.0"
	var ws2801 ledgrid.Displayer
	var gridServer *ledgrid.GridServer

	flag.StringVar(&inFile, "play", "", "Play the animation in this file instead of running as a daemon")

	// Verarbeite als erstes die Kommandozeilen-Optionen
	flag.IntVar(&width, "width", 0, "Width of panel")
	flag.IntVar(&height, "height", 0, "Height of panel")
	flag.StringVar(&customConfName, "custom", "", "Use a non standard module configuration")

	flag.UintVar(&dataPort, "data", ledgrid.DefDataPort, "Data port (UPD as well as TCP)")
	flag.UintVar(&rpcPort, "rpc", ledgrid.DefRPCPort, "RPC port")
	flag.IntVar(&baud, "baud", defBaud, "SPI baudrate in Hz")
	flag.StringVar(&missingIDs, "missing", defMissingIDs, "Comma separated list with IDs of missing LEDs (they will be skipped)")
	flag.StringVar(&defectIDs, "defect", defDefectIDs, "Comma separated list with IDs of defect LEDs (they will be blacked out)")
	flag.Parse()

	StartProfiling()
	defer StopProfiling()

	if inFile != "" {
		PlayFile(inFile)
		return
	}

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

	gridServer.HandleEvents()

	// Damit der Daemon kontrolliert beendet werden kann, installieren wir
	// einen Handler fuer das INT-Signal, welches bspw. durch Ctrl-C erzeugt
	// wird oder auch von systemd beim Stoppen eines Services verwendet wird.
	SignalHandler(gridServer)
}
