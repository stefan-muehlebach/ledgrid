package ledgrid

import (
	"image"
	"log"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/sysfs"
)

type Displayer interface {
    Size() image.Point
	DefaultGamma() (r, g, b float64)
	Close()
	Send(buffer []byte)
}

// Um einerseits nicht nur von einer Library zum Ansteuern des SPI-Bus
// abhaengig zu sein, aber auch um verschiedene SPI-Libraries miteinander zu
// vergleichen, wird die Verbindung zu den LEDs via SPI mit periph.io und
// gobot.io realisiert.

type SPIBus struct {
	spiPort   spi.PortCloser
	spiConn   spi.Conn
	maxTxSize int
}

func NewSPIBus(spiDev string, baud int) *SPIBus {
	var err error
	p := &SPIBus{}

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

func (p *SPIBus) DefaultGamma() (r, g, b float64) {
	return 3.0, 3.0, 3.0
}

func (p *SPIBus) Size() image.Point {
    return image.Point{40, 40}
}

func (p *SPIBus) Close() {
	p.spiPort.Close()
}

func (p *SPIBus) Send(buffer []byte) {
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
