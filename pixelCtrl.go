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

// Der PixelServer wird auf jenem Geraet gestartet, an dem das LedGrid via
// SPI angeschlossen ist.
type PixelServer struct {
	onRaspi         bool
	udpAddr         *net.UDPAddr
	udpConn         *net.UDPConn
	tcpAddr         *net.TCPAddr
	tcpListener     *net.TCPListener
	spiPort         spi.PortCloser
	spiConn         spi.Conn
	buffer          []byte
	maxTxSize       int
	gammaValue      [3]float64
	maxValue        [3]uint8
	gamma           [3][256]byte
	drawTestPattern bool
}

// Damit wird eine neue Instanz eines PixelServers erzeugt. Mit port wird
// sowohl die UDP- als auch die TCP-Portnummer bezeichnet. spiDev enthaelt
// das Device-File des SPI-Anschlusses und mit baud wird die Geschwindigkeit
// des SPI-Interfaces in Baud bezeichnet.
func NewPixelServer(port uint, spiDev string, baud int) *PixelServer {
	var err error
	var addrPort netip.AddrPort

	p := &PixelServer{}
	_, err = host.Init()
	if err != nil {
		log.Fatal(err)
	}
	if rpi.Present() {
		p.onRaspi = true
	}

	// Dann erstellen wir einen Buffer fuer die via Netzwerk eintreffenden
	// Daten und oeffnen die Verbindung zum LED-Grid via SPI.
	p.buffer = make([]byte, bufferSize)
	spiFs, _ := sysfs.NewSPI(0, 0)
	p.maxTxSize = spiFs.MaxTxSize()
	spiFs.Close()

	// Anschliessend werden die Tabellen fuer die Farbwertkorrektur und die
	// maximale Helligkeit erstellt.
	p.gammaValue = [3]float64{1.0, 1.0, 1.0}
	p.maxValue = [3]uint8{255, 255, 255}
	p.updateGammaTable()

	// Dann wird der SPI-Bus initialisiert.
	if p.onRaspi {
		p.spiPort, err = spireg.Open(spiDev)
		if err != nil {
			log.Fatal(err)
		}
		p.spiConn, err = p.spiPort.Connect(physic.Frequency(baud)*physic.Hertz,
			spi.Mode0, 8)
		if err != nil {
			log.Fatal(err)
		}
	}

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

func (p *PixelServer) MaxBright() (r, g, b uint8) {
	return p.maxValue[0], p.maxValue[1], p.maxValue[2]
}

func (p *PixelServer) SetMaxBright(r, g, b uint8) {
	p.maxValue[0], p.maxValue[1], p.maxValue[2] = r, g, b
	p.updateGammaTable()
}

func (p *PixelServer) updateGammaTable() {
	for color, val := range p.gammaValue {
		max := float64(p.maxValue[color])
		for i := range 256 {
			p.gamma[color][i] = byte(max * math.Pow(float64(i)/255.0, val))
		}
	}
}

func (p *PixelServer) Draw(lg *LedGrid) {
	var bufferSize int
	var err error

	bufferSize = len(lg.Pix)
	for i := 0; i < bufferSize; i += 3 {
		p.buffer[i+0] = p.gamma[0][lg.Pix[i+0]]
		p.buffer[i+1] = p.gamma[1][lg.Pix[i+1]]
		p.buffer[i+2] = p.gamma[2][lg.Pix[i+2]]
	}
	if p.onRaspi {
		for idx := 0; idx < bufferSize; idx += p.maxTxSize {
			txSize := min(p.maxTxSize, bufferSize-idx)
			if err = p.spiConn.Tx(p.buffer[idx:idx+txSize], nil); err != nil {
				log.Fatalf("Couldn't send data: %v", err)
			}
		}
		time.Sleep(20 * time.Microsecond)
	} else {
		log.Printf("Received %d bytes", bufferSize)
	}
}

func (p *PixelServer) ToggleTestPattern() {
	if p.drawTestPattern {
		p.drawTestPattern = false
		return
	} else {
		p.drawTestPattern = true
	}

	bufferSize := 3 * 20 * 20
	go func() {
		idx := 0
		for p.drawTestPattern {
			if idx == 0 {
				for i := range bufferSize {
					p.buffer[i] = 0x00
				}
			} else if idx%100 == 0 {
				p.buffer[3*(idx-1)+0] = 0xff
				p.buffer[3*(idx-1)+1] = 0x3f
				p.buffer[3*(idx-1)+2] = 0x00
			} else if idx%10 == 0 {
				p.buffer[3*(idx-1)+0] = 0x00
				p.buffer[3*(idx-1)+1] = 0x8f
				p.buffer[3*(idx-1)+2] = 0x8f
			} else if idx%5 == 0 {
				p.buffer[3*(idx-1)+0] = 0x00
				p.buffer[3*(idx-1)+1] = 0x63
				p.buffer[3*(idx-1)+2] = 0x00
			} else {
				p.buffer[3*(idx-1)+0] = 0xbf
				p.buffer[3*(idx-1)+1] = 0xbf
				p.buffer[3*(idx-1)+2] = 0xbf
			}
			if p.onRaspi {
				for i := 0; i < bufferSize; i += p.maxTxSize {
					txSize := min(p.maxTxSize, bufferSize-i)
					if err := p.spiConn.Tx(p.buffer[i:i+txSize], nil); err != nil {
						log.Fatalf("Couldn't send data: %v", err)
					}
				}
				time.Sleep(150 * time.Millisecond)
			} else {
				log.Printf("Sending %d bytes", bufferSize)
			}
			idx = (idx + 1) % (400)
		}
		for i := range bufferSize {
			p.buffer[i] = 0x00
		}
		if p.onRaspi {
			for i := 0; i < bufferSize; i += p.maxTxSize {
				txSize := min(p.maxTxSize, bufferSize-i)
				if err := p.spiConn.Tx(p.buffer[i:i+txSize], nil); err != nil {
					log.Fatalf("Couldn't send data: %v", err)
				}
			}
			time.Sleep(20 * time.Microsecond)
		} else {
			log.Printf("Sending %d bytes", bufferSize)
		}

	}()
}

// Dies ist die zentrale Verarbeitungs-Funktion des Pixel-Controllers. In ihr
// wird laufend ein Datenpaket via UDP empfangen, die empfangenen Werte gem.
// Gamma-Korrektur umgeschrieben und via SPI-Bus auf das LED-Grid uebertragen.
// Die genaue Konfiguration des LED-Grids (Anordnung der Lichterketten) ist
// dem Pixel-Controller nicht bekannt.
func (p *PixelServer) Handle() {
	var bufferSize int
	var err error

	for {
		bufferSize, err = p.udpConn.Read(p.buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			log.Fatal(err)
		}
		for i := 0; i < bufferSize; i += 3 {
			p.buffer[i+0] = p.gamma[0][p.buffer[i+0]]
			p.buffer[i+1] = p.gamma[1][p.buffer[i+1]]
			p.buffer[i+2] = p.gamma[2][p.buffer[i+2]]
		}
		if p.onRaspi {
			for idx := 0; idx < bufferSize; idx += p.maxTxSize {
				txSize := min(p.maxTxSize, bufferSize-idx)
				if err = p.spiConn.Tx(p.buffer[idx:idx+txSize], nil); err != nil {
					log.Fatalf("Couldn't send data: %v", err)
				}
			}
			time.Sleep(20 * time.Microsecond)
		} else {
			log.Printf("Received %d bytes", bufferSize)
		}
	}

	// Vor dem Beenden des Programms werden alle LEDs Schwarz geschaltet
	// damit das Panel dunkel wird.
	//
	for i := range p.buffer {
		p.buffer[i] = 0x00
	}
	if p.onRaspi {
		if err = p.spiConn.Tx(p.buffer, nil); err != nil {
			log.Printf("Error during communication via SPI: %v\n", err)
		}
	} else {
		log.Printf("Turning all LEDs off.")
	}

	if p.onRaspi {
		p.spiPort.Close()
	}
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
	if p.onRaspi {
		if err = p.spiConn.Tx(grid.Pix, nil); err != nil {
			log.Printf("Error during communication via SPI: %v.", err)
		}
	} else {
		log.Printf("Drawing grid.")
	}
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
	Draw(lg *LedGrid)
	Gamma() (r, g, b float64)
	SetGamma(r, g, b float64)
	MaxBright() (r, g, b uint8)
	SetMaxBright(r, g, b uint8)
}

// Falls die Software zur Erzeugung der Bilder auf dem gleichen Node laeuft
// an dem auch das LED-Grid angeschlossen ist, dient der PixelServer auch
// gleich als Client.
type LocalPixelClient PixelServer

func NewLocalPixelClient(port uint, spiDev string, baud int) PixelClient {
	p := NewPixelServer(port, spiDev, baud)
	return p
}

// Mit diesem Typ wird die klassische Verwendung auf zwei Nodes realisiert.
type NetPixelClient struct {
	addr      *net.UDPAddr
	conn      *net.UDPConn
	rpcClient *rpc.Client
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

	return p
}

// Schliesst die Verbindung zum Controller.
func (p *NetPixelClient) Close() {
	p.conn.Close()
}

// Sendet die Daten im Buffer b zum Controller.
func (p *NetPixelClient) Draw(lg *LedGrid) {
	var err error

	_, err = p.conn.Write(lg.Pix)
	if err != nil {
		log.Fatal(err)
	}
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

func (p *DummyPixelClient) Draw(lg *LedGrid) {

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
