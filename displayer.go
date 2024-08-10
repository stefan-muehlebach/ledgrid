package ledgrid

import (
	"image"
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/sysfs"
)

type Displayer interface {
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

// Die Emulation des LedGrids als fyne-Applikation
type PixelEmulator struct {
	Grid     *fyne.Container
	gridConf ModuleConfig
	coordMap CoordMap
	field    [][]*canvas.Circle
}

func NewPixelEmulator(width, height int) *PixelEmulator {
	e := &PixelEmulator{}
	e.Grid = container.NewGridWithRows(height)
	e.gridConf = DefaultModuleConfig(image.Point{width, height})
	e.coordMap = e.gridConf.IndexMap().CoordMap()
	e.field = make([][]*canvas.Circle, width)
	for i := range e.field {
		e.field[i] = make([]*canvas.Circle, height)
	}
	for col := range width {
		for row := range height {
			ledColor := color.White
			led := canvas.NewCircle(ledColor)
			e.field[col][row] = led
			e.Grid.Add(led)
		}
	}
	return e
}

func (e *PixelEmulator) DefaultGamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (e *PixelEmulator) Close() {}

func (e *PixelEmulator) Send(buffer []byte) {
	var r, g, b uint8
	var idx int
	var src []byte
	var needsRefresh bool

	src = buffer
	needsRefresh = false
	for i, val := range src {
		if i%3 == 0 {
			r = val
			idx = i / 3
		}
		if i%3 == 1 {
			g = val
		}
		if i%3 == 2 {
			b = val
			coord := e.coordMap[idx]
			newColor := color.RGBA{R: r, G: g, B: b, A: 0xff}
			if !needsRefresh {
				oldColor := e.field[coord.X][coord.Y].FillColor
				if newColor != oldColor {
					needsRefresh = true
				}
			}
			e.field[coord.X][coord.Y].FillColor = newColor
		}
	}
	if needsRefresh {
		e.Grid.Refresh()
	}
}
