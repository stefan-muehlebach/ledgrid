package ledgrid

import (
	"errors"
	"log"
	"math"
	"net"
	"net/http"
	"net/netip"
	"net/rpc"
	"time"
)

// This type has been introduced in order to mark some LEDs on the chain as
// 'ok', defect or missing (see constants PixelOK, PixelDefect, PixelMissing)
type LedStatusType byte

const (
	// This is the default state of a LED
	LedOK LedStatusType = iota
	// LEDs with this status will be blacked out, this mean we send color data
	// for this LED, but we send (0,0,0). This status can be used if a NeoPixel
	// receives data but does not display them correctly.
	LedDefect
	// This status can be used, if a NeoPixel does not even transmit the data
	// to the NeoPixels further down the chain. Such a pixel needs to be cut
	// out of the chain and for the time till a replacement Pixel is organized
	// and soldered in, the pixel has status missing.
	LedMissing
)

// Der GridServer wird auf jenem Geraet gestartet, an dem das LedGrid via
// SPI angeschlossen ist oder allenfalls der Emulator laeuft.
type GridServer struct {
	Disp                 Displayer
	udpAddr              *net.UDPAddr
	udpConn              *net.UDPConn
	tcpAddr              *net.TCPAddr
	tcpListener          *net.TCPListener
	buffer               []byte
	statusList           []LedStatusType
	gammaValue           [3]float64
	maxValue             [3]uint8
	gamma                [3][256]byte
	drawTestPattern      bool
	sendWatch            *Stopwatch
	RecvBytes, SentBytes int
}

// Damit wird eine neue Instanz eines GridServers erzeugt. Mit port wird
// sowohl die UDP- als auch die TCP-Portnummer bezeichnet. spiDev enthaelt
// das Device-File des SPI-Anschlusses und mit baud wird die Geschwindigkeit
// des SPI-Interfaces in Baud bezeichnet.
func NewGridServer(port uint, disp Displayer) *GridServer {
	var err error
	var addrPort netip.AddrPort
	var bufferSize int

	p := &GridServer{Disp: disp}
	bufferSize = 3 * disp.Size()
	// Dann erstellen wir einen Buffer fuer die via Netzwerk eintreffenden
	// Daten und initialisieren, die Slices fuer die fehlenden (d.h. aus
	// der LED-Kette entfernten) und die fehlerhaften (d.h. die LEDs, welche
	// als Farbe immer Schwarz erhalten sollen).
	p.buffer = make([]byte, bufferSize)
	p.statusList = make([]LedStatusType, bufferSize/3)

	// Anschliessend werden die Tabellen fuer die Farbwertkorrektur und die
	// maximale Helligkeit erstellt.
	p.gammaValue[0], p.gammaValue[1], p.gammaValue[2] = p.Disp.DefaultGamma()
	p.maxValue = [3]uint8{255, 255, 255}
	p.updateGammaTable()

	p.sendWatch = NewStopwatch()

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
func (p *GridServer) Close() {
	p.udpConn.Close()
	p.tcpListener.Close()
}

// Dies ist die zentrale Verarbeitungs-Funktion des GridServers. In ihr
// wird laufend ein Datenpaket via UDP empfangen, die empfangenen Werte gem.
// Gamma-Korrektur umgeschrieben und auf ein Ausgabegeraet uebertragen
// (SPI-Bus, Emulation, etc.) Die genaue Konfiguration des LED-Grids
// (Anordnung der Lichterketten) ist dem GridServer nicht bekannt.
func (p *GridServer) Handle() {
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
		p.sendWatch.Start()
		numLEDs = bufferSize / 3
		for srcIdx, dstIdx := 0, 0; srcIdx < numLEDs; srcIdx++ {
			if p.statusList[srcIdx] == LedMissing {
				continue
			}
			dst = p.buffer[3*dstIdx : 3*dstIdx+3 : 3*dstIdx+3]
			if p.statusList[srcIdx] == LedDefect {
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
		p.sendWatch.Stop()
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

func (p *GridServer) Watch() *Stopwatch {
	return p.sendWatch
}

// Retourniert die Gamma-Werte fuer die drei Farben.
func (p *GridServer) Gamma() (r, g, b float64) {
	return p.gammaValue[0], p.gammaValue[1], p.gammaValue[2]
}

// Setzt die Gamma-Werte fuer die Farben und aktualisiert die Mapping-Tabelle.
func (p *GridServer) SetGamma(r, g, b float64) {
	p.gammaValue[0], p.gammaValue[1], p.gammaValue[2] = r, g, b
	p.updateGammaTable()
}

// Setzt pro Farbe den maximal erlaubten Farbwert als uint8-Wert
func (p *GridServer) MaxBright() (r, g, b uint8) {
	return p.maxValue[0], p.maxValue[1], p.maxValue[2]
}

func (p *GridServer) SetMaxBright(r, g, b uint8) {
	p.maxValue[0], p.maxValue[1], p.maxValue[2] = r, g, b
	p.updateGammaTable()
}

func (p *GridServer) SetPixelStatus(idx int, stat LedStatusType) {
	p.statusList[idx] = stat
}

func (p *GridServer) updateGammaTable() {
	for color, val := range p.gammaValue {
		max := float64(p.maxValue[color])
		for i := range 256 {
			p.gamma[color][i] = byte(max * math.Pow(float64(i)/255.0, val))
		}
	}
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

	var numTestLeds = p.Disp.Size()
	var testBufferSize = 3 * numTestLeds

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
					p.buffer[3*i+0] = colorValue
					p.buffer[3*i+1] = 0x00
					p.buffer[3*i+2] = 0x00
				}
			case TestGreen:
				for i := range numTestLeds {
					p.buffer[3*i+0] = 0x00
					p.buffer[3*i+1] = colorValue
					p.buffer[3*i+2] = 0x00
				}
			case TestBlue:
				for i := range numTestLeds {
					p.buffer[3*i+0] = 0x00
					p.buffer[3*i+1] = 0x00
					p.buffer[3*i+2] = colorValue
				}
			case TestYellow:
				for i := range numTestLeds {
					p.buffer[3*i+0] = colorValue
					p.buffer[3*i+1] = colorValue
					p.buffer[3*i+2] = 0x00
				}
			case TestMagenta:
				for i := range numTestLeds {
					p.buffer[3*i+0] = colorValue
					p.buffer[3*i+1] = 0x00
					p.buffer[3*i+2] = colorValue
				}
			case TestCyan:
				for i := range numTestLeds {
					p.buffer[3*i+0] = 0x00
					p.buffer[3*i+1] = colorValue
					p.buffer[3*i+2] = colorValue
				}
			case TestWhite:
				for i := range numTestLeds {
					p.buffer[3*i+0] = colorValue
					p.buffer[3*i+1] = colorValue
					p.buffer[3*i+2] = colorValue
				}
			}

			if colorValue < 0xff {
				colorValue = (colorValue << 1) | 1
			} else {
				colorValue = 0x00
				colorMode = (colorMode + 1) % NumColorModes
			}
			p.sendWatch.Start()
			p.Disp.Send(p.buffer[:testBufferSize])
			p.sendWatch.Stop()
			time.Sleep(300 * time.Millisecond)
		}
		for i := range testBufferSize {
			p.buffer[i] = 0x00
		}
		p.sendWatch.Start()
		p.Disp.Send(p.buffer)
		p.sendWatch.Stop()
	}()

	return true
}

// Die folgenden Methoden koennen via RPC vom Client aufgerufen werden.
// Die Methode RPCDraw ist nur der Vollstaendigkeit halber vorhanden. In
// der Praxis hat sich das Senden der Bilddaten via RPC als zu langsam
// erwiesen und wurde auf UDP umgestellt.
func (p *GridServer) RPCDraw(grid *LedGrid, reply *int) error {
	var err error

	for i := 0; i < len(grid.Pix); i++ {
		grid.Pix[i] = p.gamma[i%3][grid.Pix[i]]
	}
	p.Disp.Send(grid.Pix)
	return err
}

type SizeArg int

func (p *GridServer) RPCSize(arg int, reply *SizeArg) error {
	*reply = SizeArg(p.Disp.Size())
	return nil
}

type GammaArg struct {
	RedVal, GreenVal, BlueVal float64
}

func (p *GridServer) RPCSetGamma(arg GammaArg, reply *int) error {
	p.SetGamma(arg.RedVal, arg.GreenVal, arg.BlueVal)
	return nil
}

func (p *GridServer) RPCGamma(arg int, reply *GammaArg) error {
	reply.RedVal, reply.GreenVal, reply.BlueVal = p.Gamma()
	return nil
}

type BrightArg struct {
	RedVal, GreenVal, BlueVal uint8
}

func (p *GridServer) RPCSetMaxBright(arg BrightArg, reply *int) error {
	p.SetMaxBright(arg.RedVal, arg.GreenVal, arg.BlueVal)
	return nil
}

func (p *GridServer) RPCMaxBright(arg int, reply *BrightArg) error {
	reply.RedVal, reply.GreenVal, reply.BlueVal = p.MaxBright()
	return nil
}
