package ledgrid

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"net/netip"
	"net/rpc"
	"time"

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
	"periph.io/x/host/v3/sysfs"
)

const (
	// Dies ist die Groesse des Buffers, welcher fuer den Empfang der
	// LED-Daten zur Verfuegung steht. Er ist bewusst extrem grosszuegig
	// dimensioniert... ;-)
	bufferSize = 320 * 240 * 3
)

type Displayer interface {
	Close()
	Send(bufffer []byte)
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

func OpenSPIBus(spiDev string, baud int) *SPIBus {
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

type PixelStatusType byte

const (
	PixelFine PixelStatusType = iota
	PixelMissing
	PixelDefect
)

// Der PixelServer wird auf jenem Geraet gestartet, an dem das LedGrid via
// SPI angeschlossen ist oder allenfalls der Emulator laeuft.
type PixelServer struct {
	Disp                 Displayer
	onRaspi              bool
	udpAddr              *net.UDPAddr
	udpConn              *net.UDPConn
	tcpAddr              *net.TCPAddr
	tcpListener          *net.TCPListener
	buffer               []byte
	statusList           []PixelStatusType
	gammaValue           [3]float64
	maxValue             [3]uint8
	gamma                [3][256]byte
	drawTestPattern      bool
	SendWatch            *Stopwatch
	RecvBytes, SentBytes int
}

// Damit wird eine neue Instanz eines PixelServers erzeugt. Mit port wird
// sowohl die UDP- als auch die TCP-Portnummer bezeichnet. spiDev enthaelt
// das Device-File des SPI-Anschlusses und mit baud wird die Geschwindigkeit
// des SPI-Interfaces in Baud bezeichnet.
func NewPixelServer(port uint /*, spiDev string, baud int*/) *PixelServer {
	var err error
	var addrPort netip.AddrPort

	p := &PixelServer{}
	// _, err = host.Init()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	if rpi.Present() {
		p.onRaspi = true
	}

	// Dann erstellen wir einen Buffer fuer die via Netzwerk eintreffenden
	// Daten und initialisieren, die Slices fuer die fehlenden (d.h. aus
	// der LED-Kette entfernten) und die fehlerhaften (d.h. die LEDs, welche
	// als Farbe immer Schwarz erhalten sollen).
	p.buffer = make([]byte, bufferSize)
	p.statusList = make([]PixelStatusType, bufferSize/3)

	// p.missingList = make([]int, 0)
	// p.defectList = make([]int, 0)

	// spiFs, _ := sysfs.NewSPI(0, 0)
	// p.maxTxSize = spiFs.MaxTxSize()
	// spiFs.Close()

	// Anschliessend werden die Tabellen fuer die Farbwertkorrektur und die
	// maximale Helligkeit erstellt.
	p.gammaValue = [3]float64{1.0, 1.0, 1.0}
	p.maxValue = [3]uint8{255, 255, 255}
	p.updateGammaTable()

	p.SendWatch = NewStopwatch()

	// Dann wird der SPI-Bus initialisiert.
	// if p.onRaspi {
	// 	p.spiPort, err = spireg.Open(spiDev)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	p.spiConn, err = p.spiPort.Connect(physic.Frequency(baud)*physic.Hertz,
	// 		spi.Mode0, 8)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	// Jetzt wird der UDP-Port geoeffnet, resp. eine lesende Verbindung
	// dafuer erstellt.
	addrPort = netip.AddrPortFrom(netip.IPv4Unspecified(), uint16(port))
	if !addrPort.IsValid() {
		log.Fatalf("Invalid address or port")
	}
	p.udpAddr = net.UDPAddrFromAddrPort(addrPort)
	p.udpConn, err = net.ListenUDP("udp", p.udpAddr)
	if err != nil {
		log.Fatal("UDP listen error:", err)
	}

	// Anschliessend wird die RPC-Verbindung initiiert.
	rpc.Register(p)
	rpc.HandleHTTP()
	p.tcpAddr = net.TCPAddrFromAddrPort(addrPort)
	p.tcpListener, err = net.ListenTCP("tcp", p.tcpAddr)
	if err != nil {
		log.Fatal("TCP listen error:", err)
	}
	go http.Serve(p.tcpListener, nil)

	return p
}

// Schliesst die diversen Verbindungen.
func (p *PixelServer) Close() {
	p.udpConn.Close()
    p.tcpListener.Close()
}

// Retourniert die Gamma-Werte fuer die drei Farben.
func (p *PixelServer) Gamma() (r, g, b float64) {
	return p.gammaValue[0], p.gammaValue[1], p.gammaValue[2]
}

// Setzt die Gamma-Werte fuer die Farben und aktualisiert die Mapping-Tabelle.
func (p *PixelServer) SetGamma(r, g, b float64) {
	p.gammaValue[0], p.gammaValue[1], p.gammaValue[2] = r, g, b
	p.updateGammaTable()
}

// Setzt pro Farbe den maximal erlaubten Farbwert als uint8-Wert
func (p *PixelServer) MaxBright() (r, g, b uint8) {
	return p.maxValue[0], p.maxValue[1], p.maxValue[2]
}

func (p *PixelServer) SetMaxBright(r, g, b uint8) {
	p.maxValue[0], p.maxValue[1], p.maxValue[2] = r, g, b
	p.updateGammaTable()
}

func (p *PixelServer) SetPixelStatus(idx int, stat PixelStatusType) {
	p.statusList[idx] = stat
}

func (p *PixelServer) updateGammaTable() {
	for color, val := range p.gammaValue {
		max := float64(p.maxValue[color])
		for i := range 256 {
			p.gamma[color][i] = byte(max * math.Pow(float64(i)/255.0, val))
		}
	}
}

// func (p *PixelServer) SPISendBuffer(buffer []byte) {
// 	var err error
// 	var bufferSize int

// 	bufferSize = len(buffer)
// 	if p.onRaspi {
// 		for idx := 0; idx < bufferSize; idx += p.maxTxSize {
// 			txSize := min(p.maxTxSize, bufferSize-idx)
// 			if err = p.spiConn.Tx(buffer[idx:idx+txSize], nil); err != nil {
// 				log.Fatalf("Couldn't send data: %v", err)
// 			}
// 		}
// 		time.Sleep(20 * time.Microsecond)
// 	} else {
// 		log.Printf("Sent %d bytes to the SPI bus", bufferSize)
// 	}
// }

const (
	TestRed = iota
	TestGreen
	TestBlue
	TestYellow
	TestMagenta
	TestCyan
	NumColorModes
)

const (
	TestDimmed = iota
	TestFull
	NumBrightModes
)

const (
	NumTestLeds    = 425
	TestBufferSize = 3 * NumTestLeds
)

func (p *PixelServer) ToggleTestPattern() bool {
	var colorMode, brightMode int
	var colorValue byte

	if p.drawTestPattern {
		p.drawTestPattern = false
		return false
	} else {
		p.drawTestPattern = true
		colorMode = TestRed
		brightMode = TestDimmed
	}

	go func() {
		for p.drawTestPattern {
			switch brightMode {
			case TestDimmed:
				colorValue = 0x0f
			case TestFull:
				colorValue = 0xff
			}
			switch colorMode {
			case TestRed:
				for i := range NumTestLeds {
					p.buffer[3*i+0] = colorValue
					p.buffer[3*i+1] = 0x00
					p.buffer[3*i+2] = 0x00
				}
			case TestGreen:
				for i := range NumTestLeds {
					p.buffer[3*i+0] = 0x00
					p.buffer[3*i+1] = colorValue
					p.buffer[3*i+2] = 0x00
				}
			case TestBlue:
				for i := range NumTestLeds {
					p.buffer[3*i+0] = 0x00
					p.buffer[3*i+1] = 0x00
					p.buffer[3*i+2] = colorValue
				}
			case TestYellow:
				for i := range NumTestLeds {
					p.buffer[3*i+0] = colorValue
					p.buffer[3*i+1] = colorValue
					p.buffer[3*i+2] = 0x00
				}
			case TestMagenta:
				for i := range NumTestLeds {
					p.buffer[3*i+0] = colorValue
					p.buffer[3*i+1] = 0x00
					p.buffer[3*i+2] = colorValue
				}
			case TestCyan:
				for i := range NumTestLeds {
					p.buffer[3*i+0] = 0x00
					p.buffer[3*i+1] = colorValue
					p.buffer[3*i+2] = colorValue
				}
			}

			brightMode = (brightMode + 1) % NumBrightModes
			if brightMode == 0 {
				colorMode = (colorMode + 1) % NumColorModes
			}

			p.Disp.Send(p.buffer[:TestBufferSize])
			time.Sleep(time.Second)
		}
		for i := range TestBufferSize {
			p.buffer[i] = 0x00
		}
		p.Disp.Send(p.buffer)
	}()

	return true
}

// Dies ist die zentrale Verarbeitungs-Funktion des Pixel-Controllers. In ihr
// wird laufend ein Datenpaket via UDP empfangen, die empfangenen Werte gem.
// Gamma-Korrektur umgeschrieben und via SPI-Bus auf das LED-Grid uebertragen.
// Die genaue Konfiguration des LED-Grids (Anordnung der Lichterketten) ist
// dem Pixel-Controller nicht bekannt.
func (p *PixelServer) Handle() {
	var bufferSize, numLEDs int
	var src, dst []byte
	var err error

	for {
		bufferSize, err = p.udpConn.Read(p.buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			log.Fatal(err)
		}
		p.RecvBytes += bufferSize
		p.SendWatch.Start()
		numLEDs = bufferSize / 3
		for srcIdx, dstIdx := 0, 0; srcIdx < numLEDs; srcIdx++ {
			if p.statusList[srcIdx] == PixelMissing {
				continue
			}
			dst = p.buffer[3*dstIdx : 3*dstIdx+3 : 3*dstIdx+3]
			if p.statusList[srcIdx] == PixelDefect {
				dst[0] = 0x00
				dst[1] = 0x00
				dst[2] = 0x00
			} else {
				src = p.buffer[3*srcIdx : 3*srcIdx+3 : 3*srcIdx+3]
				dst[0] = p.gamma[0][src[0]]
				dst[1] = p.gamma[1][src[1]]
				dst[2] = p.gamma[2][src[2]]
			}
			dstIdx++
		}
		p.Disp.Send(p.buffer[:bufferSize])
		p.SentBytes += bufferSize
		p.SendWatch.Stop()
	}

	// Vor dem Beenden des Programms werden alle LEDs Schwarz geschaltet
	// damit das Panel dunkel wird.
	for i := range p.buffer {
		p.buffer[i] = 0x00
	}
	p.Disp.Send(p.buffer)
	p.SentBytes += len(p.buffer)
	p.Disp.Close()
}

// Die folgenden Methoden koennen via RPC vom Client aufgerufen werden.
// Die Methode RPCDraw ist nur der Vollstaendigkeit halber vorhanden. In
// der Praxis hat sich das Senden der Bilddaten via RPC als zu langsam
// erwiesen und wurde auf UDP umgestellt.
func (p *PixelServer) RPCDraw(grid *LedGrid, reply *int) error {
	var err error

	for i := 0; i < len(grid.Pix); i++ {
		grid.Pix[i] = p.gamma[i%3][grid.Pix[i]]
	}
	p.Disp.Send(grid.Pix)
	return err
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

// Um den clientseitigen Code so generisch wie moeglich zu halten, ist der
// PixelClient als Interface definiert. Drei konkrete Implementationen
// stehen zur Verfuegung:
// - LocalPixelClient
// - NetPixelClient
// - DummyPixelClient
type PixelClient interface {
	Close()
	Send(lg *LedGrid)
	Gamma() (r, g, b float64)
	SetGamma(r, g, b float64)
	MaxBright() (r, g, b uint8)
	SetMaxBright(r, g, b uint8)
    Watch() (*Stopwatch)
}

// Falls die Software zur Erzeugung der Bilder auf dem gleichen Node laeuft
// an dem auch das LED-Grid angeschlossen ist, dient der PixelServer auch
// gleich als Client.
// type LocalPixelClient PixelServer

// func NewLocalPixelClient(port uint, spiDev string, baud int) PixelClient {
// 	p := NewPixelServer(port, spiDev, baud)
// 	return p
// }

// Mit diesem Typ wird die klassische Verwendung auf zwei Nodes realisiert.
type NetPixelClient struct {
	addr      *net.UDPAddr
	conn      *net.UDPConn
	rpcClient *rpc.Client
    sendWatch *Stopwatch
}

func NewNetPixelClient(host string, port uint) PixelClient {
	var hostPort string
	var err error

	p := &NetPixelClient{}
	hostPort = fmt.Sprintf("%s:%d", host, port)
	p.addr, err = net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		log.Fatal(err)
	}
	p.conn, err = net.DialUDP("udp", nil, p.addr)
	if err != nil {
		log.Fatal(err)
	}

	p.rpcClient, err = rpc.DialHTTP("tcp", hostPort)
	if err != nil {
		log.Fatal("Dialing:", err)
	}
    p.sendWatch = NewStopwatch()

	return p
}

// Schliesst die Verbindung zum Controller.
func (p *NetPixelClient) Close() {
	p.conn.Close()
}

// Sendet die Bilddaten in der LedGrid-Struktur zum Controller.
func (p *NetPixelClient) Send(lg *LedGrid) {
	var err error

    p.sendWatch.Start()
	_, err = p.conn.Write(lg.Pix)
	if err != nil {
		log.Fatal(err)
	}
    p.sendWatch.Stop()
}

func (p *NetPixelClient) Gamma() (r, g, b float64) {
	var reply GammaArg
	var err error

	err = p.rpcClient.Call("PixelServer.RPCGamma", 0, &reply)
	if err != nil {
		log.Fatal("Gamma error:", err)
	}
	return reply.RedVal, reply.GreenVal, reply.BlueVal
}

func (p *NetPixelClient) SetGamma(r, g, b float64) {
	var reply int
	var err error

	err = p.rpcClient.Call("PixelServer.RPCSetGamma", GammaArg{r, g, b}, &reply)
	if err != nil {
		log.Fatal("SetGamma error:", err)
	}
}

func (p *NetPixelClient) MaxBright() (r, g, b uint8) {
	var reply BrightArg
	var err error

	err = p.rpcClient.Call("PixelServer.RPCMaxBright", 0, &reply)
	if err != nil {
		log.Fatal("MaxBright error:", err)
	}
	return reply.RedVal, reply.GreenVal, reply.BlueVal
}

func (p *NetPixelClient) SetMaxBright(r, g, b uint8) {
	var reply int
	var err error

	err = p.rpcClient.Call("PixelServer.RPCSetMaxBright", BrightArg{r, g, b}, &reply)
	if err != nil {
		log.Fatal("SetMaxBright error:", err)
	}
}

func (p *NetPixelClient) Watch() *Stopwatch {
    return p.sendWatch
}

// Mit dieser Implementation des PixelClient-Interfaces kann man ohne Zugriff
// auf ein reales LED-Grid Software testen.
type DummyPixelClient struct {
}

func NewDummyPixelClient() PixelClient {
	p := &DummyPixelClient{}
	return p
}

func (p *DummyPixelClient) Close() {

}

func (p *DummyPixelClient) Send(lg *LedGrid) {

}

func (p *DummyPixelClient) Gamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (p *DummyPixelClient) SetGamma(r, g, b float64) {

}

func (p *DummyPixelClient) MaxBright() (r, g, b uint8) {
	return 0xff, 0xff, 0xff
}

func (p *DummyPixelClient) SetMaxBright(r, g, b uint8) {

}

func (p *DummyPixelClient) Watch() *Stopwatch {
    return nil
}
