//go:build !tinygo

package ledgrid

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/netip"
	"net/rpc"
	"time"

	"github.com/stefan-muehlebach/ledgrid/conf"
)

const (
	DefDataPort = 5333
	DefRPCPort  = 5332
)

// Der Datentyp ByteCount kann zum Zaehlen von Bytes verwendet werden (bspw.
// bei Netzwerk- oder Datei-IO). Die Ausgabe als Text wird formatiert. So
// werden 1024 als '1.0 kB' oder 4096 als '4.0 kB' dargestellt.
type ByteCount int64

func (b ByteCount) String() string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

// Der GridServer wird auf jenem Geraet gestartet, an dem das LedGrid via
// SPI angeschlossen ist oder allenfalls der Emulator laeuft.
type GridServer struct {
	Disp                 Displayer
	RecvBytes, SentBytes ByteCount
	udpAddr              *net.UDPAddr
	udpConn              *net.UDPConn
	tcpAddr              *net.TCPAddr
	tcpListener          *net.TCPListener
	rpcAddr              *net.TCPAddr
	rpcListener          *net.TCPListener
	bufferSize           int
	maxValue             [3]uint8
	drawTestPattern      bool
	stopwatch            *Stopwatch
}

// Damit wird eine neue Instanz eines GridServers erzeugt. Mit dataPort wird
// der Port sowohl fuer die UDP-, als auch fuer die TCP-Verbindung angegeben
// mit mit rpcPort der Port fuer die RPC-Calls. Mit disp wird dem Server
// ein konkretes, anzeigefaehiges Geraet (sog. Displayer) mitgegeben.
//
func NewGridServer(dataPort, rpcPort uint, disp Displayer) *GridServer {
	var err error
	var addrPort netip.AddrPort

	p := &GridServer{Disp: disp}
	p.bufferSize = 3 * disp.NumLeds()
	p.maxValue = [3]uint8{255, 255, 255}

	p.stopwatch = NewStopwatch()

	// Jetzt wird der UDP-Port geoeffnet, resp. eine lesende Verbindung
	// dafuer erstellt und der entsprechende Handler dafuer gestartet.
	addrPort = netip.AddrPortFrom(netip.IPv4Unspecified(), uint16(dataPort))
	if !addrPort.IsValid() {
		log.Fatalf("Invalid address or port: %v", addrPort)
	}
	p.udpAddr = net.UDPAddrFromAddrPort(addrPort)
	p.udpConn, err = net.ListenUDP("udp", p.udpAddr)
	if err != nil {
		log.Fatal("UDP listen error:", err)
	}

	// Jetzt wird der TCP-Port geoeffnet, resp. eine lesende Verbindung
	// dafuer erstellt und der entsprechende Handler dafuer gestartet.
	addrPort = netip.AddrPortFrom(netip.IPv4Unspecified(), uint16(dataPort))
	if !addrPort.IsValid() {
		log.Fatalf("Invalid address or port: %v", addrPort)
	}
	p.tcpAddr = net.TCPAddrFromAddrPort(addrPort)
	p.tcpListener, err = net.ListenTCP("tcp", p.tcpAddr)
	if err != nil {
		log.Fatal("TCP listen error:", err)
	}

	// Anschliessend wird die RPC-Verbindung initiiert.
	rpc.Register(p)
	rpc.HandleHTTP()
	addrPort = netip.AddrPortFrom(netip.IPv4Unspecified(), uint16(rpcPort))
	p.rpcAddr = net.TCPAddrFromAddrPort(addrPort)
	p.rpcListener, err = net.ListenTCP("tcp", p.rpcAddr)
	if err != nil {
		log.Fatal("RPC listen error:", err)
	}

	return p
}

func (p *GridServer) HandleEvents() {
	go p.HandleMessage(p.udpConn)
	go p.HandleTCP(p.tcpListener)
	go http.Serve(p.rpcListener, nil)
}

// Schliesst die diversen Verbindungen.
func (p *GridServer) Close() {
	p.udpConn.Close()
	p.tcpListener.Close()
	p.rpcListener.Close()
	p.Disp.Close()
}

// Dies ist die zentrale Verarbeitungs-Funktion des GridServers. In ihr
// wird laufend ein Datenpaket via UDP empfangen und die empfangenen Werte auf
// ein Ausgabegeraet uebertragen (SPI-Bus, Emulation, etc.) Die genaue
// Konfiguration des LED-Grids (Anordnung der Lichterketten) ist dem
// GridServer nicht bekannt.
func (p *GridServer) HandleMessage(conn net.Conn) {
	var bufferSize int
	var err error
	var buffer []byte

	buffer = make([]byte, p.bufferSize)
	for {
		bufferSize, err = conn.Read(buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				break
			}
			log.Fatalf("Failed Read(): %v", err)
		}
		p.RecvBytes += ByteCount(bufferSize)
		p.stopwatch.Start()
		p.Disp.Display(buffer)
		p.SentBytes += ByteCount(bufferSize)
		p.stopwatch.Stop()
	}

	// Vor dem Beenden des Programms werden alle LEDs Schwarz geschaltet
	// damit das Panel dunkel wird.
	for i := range buffer {
		buffer[i] = 0x00
	}
	p.Disp.Send(buffer)
	p.SentBytes += ByteCount(len(buffer))
}

// Damit werden Meldungen via TCP empfangen und verarbeitet.
func (p *GridServer) HandleTCP(lsnr *net.TCPListener) {
	var conn net.Conn
	var err error

	for {
		conn, err = lsnr.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			log.Fatalf("Failed TCP Accept(): %v", err)
		}
		go p.HandleMessage(conn)
	}
}

func (p *GridServer) Stopwatch() *Stopwatch {
	return p.stopwatch
}

// Retourniert die Gamma-Werte fuer die drei Farben.
func (p *GridServer) Gamma() (r, g, b float64) {
	return p.Disp.Gamma()
}

// Setzt die Gamma-Werte fuer die Farben und aktualisiert die Mapping-Tabelle.
func (p *GridServer) SetGamma(r, g, b float64) {
	p.Disp.SetGamma(r, g, b)
}

func (p *GridServer) ModuleConfig() conf.ModuleConfig {
	return p.Disp.ModuleConfig()
}

// Setzt pro Farbe den maximal erlaubten Farbwert als uint8-Wert
/*
func (p *GridServer) MaxBright() (r, g, b uint8) {
	return p.maxValue[0], p.maxValue[1], p.maxValue[2]
}

func (p *GridServer) SetMaxBright(r, g, b uint8) {
	p.maxValue[0], p.maxValue[1], p.maxValue[2] = r, g, b
}
*/

func (p *GridServer) SetPixelStatus(idx int, stat LedStatusType) {
	p.Disp.SetPixelStatus(idx, stat)
}

const (
	TestRed = iota
	TestGreen
	TestBlue
	TestYellow
	TestMagenta
	TestCyan
	TestWhite
	NumColorModes
)

func (p *GridServer) ToggleTestPattern() bool {
	var colorMode int
	var colorValue byte
	var numTestLeds = p.Disp.NumLeds()
	var testBufferSize = 3 * numTestLeds
	var buffer []byte

	buffer = make([]byte, p.bufferSize)
	if p.drawTestPattern {
		p.drawTestPattern = false
		return false
	} else {
		p.drawTestPattern = true
		colorMode = TestRed
	}

	go func() {
		colorValue = 0x00
		for p.drawTestPattern {
			switch colorMode {
			case TestRed:
				for i := range numTestLeds {
					buffer[3*i+0] = colorValue
					buffer[3*i+1] = 0x00
					buffer[3*i+2] = 0x00
				}
			case TestGreen:
				for i := range numTestLeds {
					buffer[3*i+0] = 0x00
					buffer[3*i+1] = colorValue
					buffer[3*i+2] = 0x00
				}
			case TestBlue:
				for i := range numTestLeds {
					buffer[3*i+0] = 0x00
					buffer[3*i+1] = 0x00
					buffer[3*i+2] = colorValue
				}
			case TestYellow:
				for i := range numTestLeds {
					buffer[3*i+0] = colorValue
					buffer[3*i+1] = colorValue
					buffer[3*i+2] = 0x00
				}
			case TestMagenta:
				for i := range numTestLeds {
					buffer[3*i+0] = colorValue
					buffer[3*i+1] = 0x00
					buffer[3*i+2] = colorValue
				}
			case TestCyan:
				for i := range numTestLeds {
					buffer[3*i+0] = 0x00
					buffer[3*i+1] = colorValue
					buffer[3*i+2] = colorValue
				}
			case TestWhite:
				for i := range numTestLeds {
					buffer[3*i+0] = colorValue
					buffer[3*i+1] = colorValue
					buffer[3*i+2] = colorValue
				}
			}

			if colorValue < 0xff {
				colorValue = (colorValue << 1) | 1
			} else {
				colorValue = 0x00
				colorMode = (colorMode + 1) % NumColorModes
			}
			p.stopwatch.Start()
			p.Disp.Send(buffer)
			p.stopwatch.Stop()
			time.Sleep(300 * time.Millisecond)
		}
		for i := range testBufferSize {
			buffer[i] = 0x00
		}
		p.stopwatch.Start()
		p.Disp.Send(buffer)
		p.stopwatch.Stop()
	}()

	return true
}

// Die folgenden Methoden koennen via RPC vom Client aufgerufen werden.
type NumLedsArg int

func (p *GridServer) RPCNumLeds(arg int, reply *NumLedsArg) error {
	*reply = NumLedsArg(p.Disp.NumLeds())
	return nil
}

type GammaArg struct {
	RedVal, GreenVal, BlueVal float64
}

func (p *GridServer) RPCGamma(arg int, reply *GammaArg) error {
	reply.RedVal, reply.GreenVal, reply.BlueVal = p.Gamma()
	return nil
}

func (p *GridServer) RPCSetGamma(arg GammaArg, reply *int) error {
	p.SetGamma(arg.RedVal, arg.GreenVal, arg.BlueVal)
	return nil
}

/*
type BrightArg struct {
	RedVal, GreenVal, BlueVal uint8
}

func (p *GridServer) RPCMaxBright(arg int, reply *BrightArg) error {
	reply.RedVal, reply.GreenVal, reply.BlueVal = p.MaxBright()
	return nil
}

func (p *GridServer) RPCSetMaxBright(arg BrightArg, reply *int) error {
	p.SetMaxBright(arg.RedVal, arg.GreenVal, arg.BlueVal)
	return nil
}
*/

type ModuleConfigArg struct {
	ModuleConfig conf.ModuleConfig
}

func (p *GridServer) RPCModuleConfig(arg int, reply *ModuleConfigArg) error {
	reply.ModuleConfig = p.Disp.ModuleConfig()
	return nil
}
