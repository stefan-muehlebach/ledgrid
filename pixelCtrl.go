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
	bufferSize = 320 * 240 * 3
)

// Der PixelServer wird auf jenem Geraet gestartet, an dem das LedGrid via
// SPI angeschlossen ist.
type PixelServer struct {
	onRaspi     bool
	udpAddr     *net.UDPAddr
	udpConn     *net.UDPConn
	tcpAddr     *net.TCPAddr
	tcpListener *net.TCPListener
	spiPort     spi.PortCloser
	spiConn     spi.Conn
	buffer      []byte
	maxTxSize   int
	gammaValue  [3]float64
	gamma       [3][256]byte
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
	// Daten.
	//
	p.buffer = make([]byte, bufferSize)
	spiFs, _ := sysfs.NewSPI(0, 0)
	p.maxTxSize = spiFs.MaxTxSize()
	spiFs.Close()

	// Anschliessend werden die Tabellen fuer die Farbwertkorrektur erstellt.
	//
	p.SetGamma(1.0, 1.0, 1.0)

	// Dann wird der SPI-Bus initialisiert.
	//
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
	//
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
	//
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

func (p *PixelServer) Close() {
	p.udpConn.Close()
}

func (p *PixelServer) Gamma() (r, g, b float64) {
	return p.gammaValue[0], p.gammaValue[1], p.gammaValue[2]
}

func (p *PixelServer) SetGamma(r, g, b float64) {
	p.gammaValue[0], p.gammaValue[1], p.gammaValue[2] = r, g, b
	for color, val := range p.gammaValue {
		for i := range 256 {
			p.gamma[color][i] = byte(255.0 * math.Pow(float64(i)/255.0, val))
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

// Um den clientseitigen Code so generisch wie moeglich zu halten, ist der
// PixelClient als Interface definiert. Zwei konkrete Implementationen
// stehen zur Verfuegung:
// - LocalPixelClient
// - NetPixelClient
type PixelClient interface {
    Close()
    Draw(lg *LedGrid)
    Gamma() (r, g, b float64)
    SetGamma(r, g, b float64)
}

// Dieser Typ wird client-seitig fuer die Ansteuerung des LedGrid verwendet.
// Im Wesentlichen ist dies eine Abstraktion der Ansteuerung via UDP.
type NetPixelClient struct {
	addr      *net.UDPAddr
	conn      *net.UDPConn
	rpcClient *rpc.Client
}

// Erzeugt ein neues Controller-Objekt, welches das LedGrid ueber die Adresse
// in Host und den UDP-Port in Port anspricht.
func NewNetPixelClient(host string, port uint) *NetPixelClient {
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
