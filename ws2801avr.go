//go:build tinygo

package ledgrid

import (
	"time"

	"github.com/stefan-muehlebach/ledgrid/conf"

	"machine"
)

// Dies ist die Implementation eines Displayers, welcher eine Lichterkette von
// NeoPixeln mit WS2801 via SPI-Bus auf einem RaspberryPi ansteuert.
type WS2801avr struct {
	DisplayEmbed
	spi machine.SPI
}

// Erstellt eine neue Instanz. spiDev ist das Device-File des SPI-Buses, baud
// die Taktrate (in Bit pro Sekunde) und numLeds die Anzahl NeoPixel auf der
// Lichterkette - ohne die entfernten NeoPixel zu beruecksichtigen.
func NewWS2801avr(modConf conf.ModuleConfig) *WS2801avr {
	var err error
	p := &WS2801avr{}

	p.DisplayEmbed.Init(p, len(modConf)*conf.ModuleDim.X*conf.ModuleDim.Y)
	p.SetModuleConfig(modConf)

	p.spi = machine.SPI0
	if err = p.spi.Configure(machine.SPIConfig{
		Frequency: 3 * machine.MHz,
	}); err != nil {
		println("Error configuring SPI:")
		println(err.Error())
		return nil
	}

	return p
}

// Diese Methode gehoert zum Displayer-Interface und retourniert die
// empfohlenen Gamma-Werte fuer die drei Farbkanaele Rot, Gruen und Blau.
func (p *WS2801avr) DefaultGamma() (r, g, b float64) {
	return 2.5, 2.5, 2.5
}

// Schliesst den Displayer, in diesem Fall den SPI-Port.
func (p *WS2801avr) Close() {

}

// Sendet die Farbwerte in buffer via SPI-Bus zur NeoPixel Lichterkette. Die
// Reihenfolge der Pixel muss bereits vorgaengig der effektiven Verkabelung
// angepasst worden sein, ebenso die Farbwertkorrektur. Diese Methode wird
// ueblicherweise vom DisplayEmbed und nicht von Benutzercode aufgerufen.
func (p *WS2801avr) Send(buffer []byte) {
	var err error

	if err = p.spi.Tx(buffer, nil); err != nil {
		println("Couldn't send data:", err.Error())
		return
	}
	time.Sleep(20 * time.Microsecond)
}
