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

	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

const (
	bufferSize = 1024
)

type PixelServer struct {
	onRaspi     bool
	udpAddr     *net.UDPAddr
	udpConn     *net.UDPConn
	tcpAddr     *net.TCPAddr
	tcpListener *net.TCPListener
	spiPort     spi.PortCloser
	spiConn     spi.Conn
	buffer      []byte
	gammaValue  [3]float64
	gamma       [3][256]byte
}

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
	// Daten. 1kB sollten aktuell reichen (entspricht rund 340 RGB-Werten).
	//
	p.buffer = make([]byte, bufferSize)

	// Anschliessend wird die Tabelle fuer die Farbwertkorrektur erstellt.
	//
	p.SetGamma(0, 1.0)
	p.SetGamma(1, 1.0)
	p.SetGamma(2, 1.0)

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

func (p *PixelServer) Gamma(color int) (float64) {
    return p.gammaValue[color]
}

func (p *PixelServer) SetGamma(color int, value float64) {
	p.gammaValue[color] = value
	for i := 0; i < 256; i++ {
		p.gamma[color][i] = byte(255.0 * math.Pow(float64(i)/255.0,
			p.gammaValue[color]))
	}
}

func (p *PixelServer) Handle() {
	var len int
	var err error

	for {
		len, err = p.udpConn.Read(p.buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			log.Fatal(err)
		}
		if len != 300 {
			log.Printf("Only %d bytes received instead of 300.\n", len)
		}
		for i := 0; i < len; i += 3 {
			p.buffer[i+0] = p.gamma[0][p.buffer[i+0]]
			p.buffer[i+1] = p.gamma[1][p.buffer[i+1]]
			p.buffer[i+2] = p.gamma[2][p.buffer[i+2]]
		}
		if p.onRaspi {
			if err = p.spiConn.Tx(p.buffer[:len], nil); err != nil {
				log.Printf("Error during communication via SPI: %v\n", err)
			}
		} else {
			log.Printf("Received %d bytes", len)
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

// Die folgenden Methoden werden via RPC vom Client aufgerufen.
func (p *PixelServer) DrawRPC(grid *LedGrid, reply *int) error {
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
	Color int
	Value float64
}

func (p *PixelServer) SetGammaRPC(arg GammaArg, reply *int) error {
	p.SetGamma(arg.Color, arg.Value)
	return nil
}

// Dieser Typ wird client-seitig fuer die Ansteuerung des LedGrid verwendet.
// Im Wesentlichen ist dies eine Abstraktion der Ansteuerung via UDP.
type PixelClient struct {
	addr      *net.UDPAddr
	conn      *net.UDPConn
	rpcClient *rpc.Client
}

// Erzeugt ein neues Controller-Objekt, welches das LedGrid ueber die Adresse
// in Host und den UDP-Port in Port anspricht.
func NewPixelClient(host string, port uint) *PixelClient {
	var hostPort string
	var err error

	p := &PixelClient{}
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
func (p *PixelClient) Close() {
	p.conn.Close()
}

// Sendet die Daten im Buffer b zum Controller.
func (p *PixelClient) Draw(ledGrid *LedGrid) {
	var err error

	_, err = p.conn.Write(ledGrid.Pix)
	if err != nil {
		log.Fatal(err)
	}
}

// func (p *PixelClient) Draw(grid *LedGrid) {
// 	var reply int
// 	var err error

// 	err = p.rpcClient.Call("PixelServer.DrawRPC", grid, &reply)
// 	if err != nil {
// 		log.Fatal("Draw error:", err)
// 	}
// }

func (p *PixelClient) SetGamma(color int, value float64) {
	var reply int
	var err error

	err = p.rpcClient.Call("PixelServer.SetGammaRPC", GammaArg{color, value}, &reply)
	if err != nil {
		log.Fatal("SetGamma error:", err)
	}
}
