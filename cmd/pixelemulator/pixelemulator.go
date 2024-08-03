package main

import (
	"errors"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/netip"
	"net/rpc"

	"github.com/stefan-muehlebach/ledgrid"

	// "github.com/stefan-muehlebach/gg/color"
	"fyne.io/fyne/v2/canvas"
	_ "github.com/stefan-muehlebach/ledgrid"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	_ "fyne.io/fyne/v2/driver/desktop"
)

const (
	Width      = 40
	Height     = 10
	PixelSize  = 50.0
	AppWidth   = Width * PixelSize
	AppHeight  = Height * PixelSize
	BufferSize = Width * Height * 3
	Port       = 5333
)

var (
	AppSize     = fyne.NewSize(AppWidth, AppHeight)
	FieldSize   = fyne.NewSize(PixelSize, PixelSize)
	App         fyne.App
	Win         fyne.Window
	Grid        *fyne.Container
	Buffer      []byte
	LEDField    [][]*canvas.Circle
	CoordMap    ledgrid.CoordMap
	udpAddr     *net.UDPAddr
	udpConn     *net.UDPConn
	tcpAddr     *net.TCPAddr
	tcpListener *net.TCPListener
	srv         PixelServer
	gammaValue  [3]float64
	maxValue    [3]uint8
	gamma       [3][256]byte
)

type PixelServer struct{}

// Retourniert die Gamma-Werte fuer die drei Farben.
func (p *PixelServer) Gamma() (r, g, b float64) {
	return gammaValue[0], gammaValue[1], gammaValue[2]
}

// Setzt die Gamma-Werte fuer die Farben und aktualisiert die Mapping-Tabelle.
func (p *PixelServer) SetGamma(r, g, b float64) {
	gammaValue[0], gammaValue[1], gammaValue[2] = r, g, b
	p.updateGammaTable()
}

func (p *PixelServer) MaxBright() (r, g, b uint8) {
	return maxValue[0], maxValue[1], maxValue[2]
}

func (p *PixelServer) SetMaxBright(r, g, b uint8) {
	maxValue[0], maxValue[1], maxValue[2] = r, g, b
	p.updateGammaTable()
}

func (p *PixelServer) updateGammaTable() {
	for color, val := range gammaValue {
		max := float64(maxValue[color])
		for i := range 256 {
			gamma[color][i] = byte(max * math.Pow(float64(i)/255.0, val))
		}
	}
}

func (p *PixelServer) RPCDraw(grid *ledgrid.LedGrid, reply *int) error {
	// var err error

	// for i := 0; i < len(grid.Pix); i++ {
	// 	grid.Pix[i] = p.gamma[i%3][grid.Pix[i]]
	// }
	// if p.onRaspi {
	// 	if err = p.spiConn.Tx(grid.Pix, nil); err != nil {
	// 		log.Printf("Error during communication via SPI: %v.", err)
	// 	}
	// } else {
	// 	log.Printf("Drawing grid.")
	// }
	// return err
	return nil
}

type GammaArg struct {
	RedVal, GreenVal, BlueVal float64
}

func (p *PixelServer) RPCSetGamma(arg GammaArg, reply *int) error {
	p.SetGamma(arg.RedVal, arg.GreenVal, arg.BlueVal)
	return nil
}

func (p *PixelServer) RPCGamma(arg int, reply *GammaArg) error {
	reply.RedVal, reply.GreenVal, reply.BlueVal = p.Gamma()
	return nil
}

type BrightArg struct {
	RedVal, GreenVal, BlueVal uint8
}

func (p *PixelServer) RPCSetMaxBright(arg BrightArg, reply *int) error {
	p.SetMaxBright(arg.RedVal, arg.GreenVal, arg.BlueVal)
	return nil
}

func (p *PixelServer) RPCMaxBright(arg int, reply *BrightArg) error {
	reply.RedVal, reply.GreenVal, reply.BlueVal = p.MaxBright()
	return nil
}

func Handle() {
	var bufferSize int
	var err error

	for {
		bufferSize, err = udpConn.Read(Buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			log.Fatal(err)
		}
		for i := 0; i < bufferSize; i += 3 {
			Buffer[i+0] = gamma[0][Buffer[i+0]]
			Buffer[i+1] = gamma[1][Buffer[i+1]]
			Buffer[i+2] = gamma[2][Buffer[i+2]]
		}
		UpdateField(Buffer[:bufferSize], LEDField)
	}

	// Vor dem Beenden des Programms werden alle LEDs Schwarz geschaltet
	// damit das Panel dunkel wird.
	for i := range Buffer {
		Buffer[i] = 0x00
	}
	UpdateField(Buffer[:bufferSize], LEDField)
}

func UpdateField(buffer []byte, field [][]*canvas.Circle) {
	var r, g, b uint8
	var idx int

	for i, val := range buffer {
		if i%3 == 0 {
			r = val
			idx = i / 3
		}
		if i%3 == 1 {
			g = val
		}
		if i%3 == 2 {
			b = val
			coord := CoordMap[idx]
			field[coord.X][coord.Y].FillColor = color.RGBA{R: r, G: g, B: b, A: 0xff}
		}
	}
	Grid.Refresh()
}

func main() {
	var addrPort netip.AddrPort
	var err error

	App = app.New()
	Win = App.NewWindow("LedGrid Emulator")

	Buffer = make([]byte, BufferSize)
	gridConf := ledgrid.DefaultModuleConfig(image.Point{Width, Height})
	CoordMap = gridConf.IndexMap().CoordMap()

    gammaValue = [3]float64{1.0, 1.0, 1.0}
	maxValue = [3]uint8{255, 255, 255}
	srv.updateGammaTable()


	// Jetzt wird der UDP-Port geoeffnet, resp. eine lesende Verbindung
	// dafuer erstellt.
	addrPort = netip.AddrPortFrom(netip.IPv4Unspecified(), uint16(Port))
	if !addrPort.IsValid() {
		log.Fatalf("Invalid address or port")
	}
	udpAddr = net.UDPAddrFromAddrPort(addrPort)
	udpConn, err = net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal("UDP listen error:", err)
	}

	// Anschliessend wird die RPC-Verbindung initiiert.
	rpc.Register(&srv)
	rpc.HandleHTTP()
	tcpAddr = net.TCPAddrFromAddrPort(addrPort)
	tcpListener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal("TCP listen error:", err)
	}
	go http.Serve(tcpListener, nil)

	LEDField = make([][]*canvas.Circle, Width)
	for i := range LEDField {
		LEDField[i] = make([]*canvas.Circle, Height)
	}

	Grid = container.NewGridWithRows(Height)
	for col := range Width {
		for row := range Height {
			ledColor := color.White
			led := canvas.NewCircle(ledColor)
			LEDField[col][row] = led
			Grid.Add(led)
		}
	}

	Win.Canvas().SetOnTypedKey(func(evt *fyne.KeyEvent) {
		switch evt.Name {
		case fyne.KeyEscape, fyne.KeyQ:
			App.Quit()
		case fyne.KeySpace:
			for i := range Buffer {
				Buffer[i] = byte(rand.Intn(256))
			}
			UpdateField(Buffer, LEDField)
		}
	})

	go Handle()

	Win.SetContent(Grid)
	Win.Resize(AppSize)
	Win.ShowAndRun()
}
