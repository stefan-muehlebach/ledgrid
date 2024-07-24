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
	defGammaValues = "3.0,3.0,3.0"
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
			num, total, avg := pixelServer.SendWatch.Stats()
			log.Printf("   %d sends to SPI took %v (%v per send)", num, total, avg)
			log.Printf("   %d bytes received by the controller", pixelServer.RecvBytes)
			log.Printf("   %d bytes sent by the controller", pixelServer.SentBytes)
		case syscall.SIGUSR1:
			log.Printf("Toggle drawing of a test pattern.")
			pixelServer.ToggleTestPattern()
		}
	}
}

func main() {
	var port uint
	var baud int
	var gammaValues string
	var gammaValue [3]float64
	var spiDevFile string = "/dev/spidev0.0"
	var pixelServer *ledgrid.PixelServer

	// Verarbeite als erstes die Kommandozeilen-Optionen
	flag.UintVar(&port, "port", defPort, "UDP port")
	flag.IntVar(&baud, "baud", defBaud, "SPI baudrate in Hz")
	flag.StringVar(&gammaValues, "gamma", defGammaValues, "Gamma values")
	flag.Parse()

	for i, str := range strings.Split(gammaValues, ",") {
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			log.Fatalf("Wrong format: %s", str)
		}
		gammaValue[i] = val
	}

	pixelServer = ledgrid.NewPixelServer(port, spiDevFile, baud)
	pixelServer.SetGamma(gammaValue[red], gammaValue[green], gammaValue[blue])

	// Damit der Daemon kontrolliert beendet werden kann, installieren wir
	// einen Handler fuer das INT-Signal, welches bspw. durch Ctrl-C erzeugt
	// wird oder auch von systemd beim Stoppen eines Services verwendet wird.
	go SignalHandler(pixelServer)
	//
	// go func() {
	// 	sigChan := make(chan os.Signal)
	// 	signal.Notify(sigChan, os.Interrupt)
	// 	<-sigChan
	// 	pixelServer.Close()
	// }()

	pixelServer.Handle()
}
