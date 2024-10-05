package ledgrid

import (
	"log"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/sysfs"
)

// Um einerseits nicht nur von einer Library zum Ansteuern des SPI-Bus
// abhaengig zu sein, aber auch um verschiedene SPI-Libraries miteinander zu
// vergleichen, wird die Verbindung zu den LEDs via SPI mit periph.io und
// gobot.io realisiert.

type WS2801 struct {
    DisplayEmbed
	spiPort   spi.PortCloser
	spiConn   spi.Conn
	maxTxSize int
}

func NewWS2801(spiDev string, baud int, size int) *WS2801 {
	var err error
	p := &WS2801{}

    p.DisplayEmbed.Init(p, size)
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

func (p *WS2801) DefaultGamma() (r, g, b float64) {
	return 2.5, 2.5, 2.5
}

func (p *WS2801) Close() {
	p.spiPort.Close()
}

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
