//go:build !tinygo

package ledgrid

import (
	"github.com/stefan-muehlebach/ledgrid/conf"
	"log"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/sysfs"
)

// Dies ist die Implementation eines Displayers, welcher eine Lichterkette von
// NeoPixeln mit WS2801 via SPI-Bus auf einem RaspberryPi ansteuert.
type WS2801 struct {
	DisplayEmbed
	spiPort   spi.PortCloser
	spiConn   spi.Conn
	maxTxSize int
}

// Erstellt eine neue Instanz. spiDev ist das Device-File des SPI-Buses, baud
// die Taktrate (in Bit pro Sekunde) und numLeds die Anzahl NeoPixel auf der
// Lichterkette - ohne die entfernten NeoPixel zu beruecksichtigen.
func NewWS2801(spiDev string, baud int, modConf conf.ModuleConfig) *WS2801 {
	var err error
	p := &WS2801{}

	p.DisplayEmbed.Init(p, len(modConf)*conf.ModuleDim.X*conf.ModuleDim.Y)
    p.SetModuleConfig(modConf)
	_, err = host.Init()
	if err != nil {
		log.Fatal(err)
	}

	spiFs, _ := sysfs.NewSPI(0, 0)
	p.maxTxSize = spiFs.MaxTxSize()
	spiFs.Close()

	p.spiPort, err = spireg.Open(spiDev)
	if err != nil {
		log.Fatal(err)
	}
	p.spiConn, err = p.spiPort.Connect(physic.Frequency(baud)*physic.Hertz,
		spi.Mode0, 8)
	if err != nil {
		log.Fatal(err)
	}

	return p
}

// Diese Methode gehoert zum Displayer-Interface und retourniert die
// empfohlenen Gamma-Werte fuer die drei Farbkanaele Rot, Gruen und Blau.
func (p *WS2801) DefaultGamma() (r, g, b float64) {
	return 2.5, 2.5, 2.5
}

// Schliesst den Displayer, in diesem Fall den SPI-Port.
func (p *WS2801) Close() {
	p.spiPort.Close()
}

// Sendet die Farbwerte in buffer via SPI-Bus zur NeoPixel Lichterkette. Die
// Reihenfolge der Pixel muss bereits vorgaengig der effektiven Verkabelung
// angepasst worden sein, ebenso die Farbwertkorrektur. Diese Methode wird
// ueblicherweise vom DisplayEmbed und nicht von Benutzercode aufgerufen.
func (p *WS2801) Send(buffer []byte) {
	var err error
	var bufferSize int

	bufferSize = len(buffer)
	for idx := 0; idx < bufferSize; idx += p.maxTxSize {
		txSize := min(p.maxTxSize, bufferSize-idx)
		if err = p.spiConn.Tx(buffer[idx:idx+txSize:idx+txSize], nil); err != nil {
			log.Fatalf("Couldn't send data: %v", err)
		}
	}
	time.Sleep(20 * time.Microsecond)
}
