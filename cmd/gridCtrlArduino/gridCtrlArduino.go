//go:build tinygo

package main

import (
	"image"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/stefan-muehlebach/ledgrid"
	"github.com/stefan-muehlebach/ledgrid/conf"

	"tinygo.org/x/drivers/netlink"
	"tinygo.org/x/drivers/netlink/probe"
)

// Several settings which are controllable by command line arguments
// are here defined as constants.
const (
	Width          = 40
	Height         = 10
	CustomConfName = ""
	DataPort       = ledgrid.DefDataPort
	MissingIDs     = ""
	DefectIDs      = ""
	Ssid           = "Indernett"
	Passphrase     = "0000000000"
)

var (
	dataPort, rpcPort uint
)

func main() {
	var width, height int
	var customConfName string
	var gridSize image.Point
	var modConf conf.ModuleConfig

	var missingIDs, defectIDs string
	var ws2801 ledgrid.Displayer
	var gridServer *ledgrid.GridServer

	width = Width
	height = Height
	customConfName = CustomConfName
	dataPort = DataPort
	missingIDs = MissingIDs
	defectIDs = DefectIDs

	time.Sleep(2 * time.Second)

	// Establish the connection to the WiFi-Module and open it using Ssid and
	// Passphrase
    println("Here we go...")
    println("Probe network driver")
	link, _ := probe.Probe()

    println("Connect to WiFi")
	err := link.NetConnect(&netlink.ConnectParams{
		Ssid:       Ssid,
		Passphrase: Passphrase,
	})
	if err != nil {
		println("Couldn't connect to the network module:", err.Error())
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

    println("Setup connection to WS2891 by SPI bus")
	ws2801 = ledgrid.NewWS2801(modConf)

    println("Initialize Server")
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

    println("Start receiving data from the network")
	gridServer.HandleEvents()

    println("Server has been terminated")
}
