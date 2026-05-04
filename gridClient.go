package ledgrid

import (
	"fmt"
	"image"
	"log"
	"net"
	"net/rpc"
	"os"

	"github.com/stefan-muehlebach/ledgrid/conf"
)

// Um den clientseitigen Code so generisch wie moeglich zu halten, ist der
// GridClient als Interface definiert. Aktuell stehen zwei Implementationen
// am Start:
//
//		NetGridClient  - Verbindet sich via UDP und RPC mit einem externen
//		                 gridController.
//	 FileSaveClient - Schreibt die Bilddaten in ein File, welches dann auf das
//	                  System mit dem Grid-Controller kopiert und dort direkt
//	                  abgespielt werden kann.
type GridClient interface {
	Send(buffer []byte)
	NumLeds() int
	Gamma() (r, g, b float64)
	SetGamma(r, g, b float64)
	ModuleConfig() conf.ModuleConfig
	Stopwatch() *Stopwatch
	Close()
}

// Mit diesem Typ wird die klassische Verwendung auf zwei Nodes realisiert.
type NetGridClient struct {
	conn        net.Conn
	rpcDisabled bool
	rpcClient   *rpc.Client
	stopwatch   *Stopwatch
}

func NewNetGridClient(host string, port, rpcPort uint) GridClient {
	var hostPortData, hostPortRPC string
	var err error

	p := &NetGridClient{}

	hostPortData = fmt.Sprintf("%s:%d", host, port)
	p.conn, err = net.Dial("tcp", hostPortData)
	if err != nil {
		log.Fatalf("Error in Dial(dataPort): %v", err)
	}

	hostPortRPC = fmt.Sprintf("%s:%d", host, rpcPort)
	if rpcPort != 0 {
		p.rpcClient, err = rpc.DialHTTP("tcp", hostPortRPC)
		if err != nil {
			log.Fatalf("Error in Dial(rpcPort): %v", err)
		}
	}

	p.stopwatch = NewStopwatch()

	return p
}

// Sendet die Bilddaten in der LedGrid-Struktur zum Controller.
func (p *NetGridClient) Send(buffer []byte) {
	var err error

	p.stopwatch.Start()
	_, err = p.conn.Write(buffer)
	if err != nil {
		log.Fatal(err)
	}
	p.stopwatch.Stop()
}

// Die folgenden Methoden verpacken die entsprechenden RPC-Calls zum
// Grid-Server.
func (p *NetGridClient) NumLeds() int {
	var reply NumLedsArg
	var err error

	if p.rpcClient == nil {
		return 400
	}
	err = p.rpcClient.Call("GridServer.RPCNumLeds", 0, &reply)
	if err != nil {
		log.Fatal("NumLeds error:", err)
	}
	return int(reply)
}

func (p *NetGridClient) Gamma() (r, g, b float64) {
	var reply GammaArg
	var err error

	if p.rpcClient == nil {
		return 2.5, 2.5, 2.5
	}
	err = p.rpcClient.Call("GridServer.RPCGamma", 0, &reply)
	if err != nil {
		log.Fatal("Gamma error:", err)
	}
	return reply.RedVal, reply.GreenVal, reply.BlueVal
}

func (p *NetGridClient) SetGamma(r, g, b float64) {
	var reply int
	var err error

	if p.rpcClient == nil {
		return
	}
	err = p.rpcClient.Call("GridServer.RPCSetGamma", GammaArg{r, g, b}, &reply)
	if err != nil {
		log.Fatal("SetGamma error:", err)
	}
}

func (p *NetGridClient) ModuleConfig() conf.ModuleConfig {
	var reply ModuleConfigArg
	var err error

	if p.rpcClient == nil {
		return conf.DefaultModuleConfig(image.Point{40, 10})
	}
	err = p.rpcClient.Call("GridServer.RPCModuleConfig", 0, &reply)
	if err != nil {
		log.Fatal("ModuleConfig error:", err)
	}
	return reply.ModuleConfig
}

func (p *NetGridClient) Stopwatch() *Stopwatch {
	return p.stopwatch
}

// Schliesst die Verbindung zum Controller.
func (p *NetGridClient) Close() {
	p.conn.Close()
}

func (p *NetGridClient) Address() string {
	return p.conn.RemoteAddr().String()
}

// Diese Client-Variante kommt dann zum Einsatz, wenn die LED-Hardware am
// selben Rechner angeschlossen ist, auf dem auch die Animationen gerechnet
// werden.
type DirectGridClient struct {
	Disp      Displayer
	stopwatch *Stopwatch
}

// Erstellt wird dieser Client-Typ mit einem Displayer, dem Interface zur
// LED-Hardware.
func NewDirectGridClient(Disp Displayer) GridClient {
	c := &DirectGridClient{}

	c.Disp = Disp
	c.stopwatch = NewStopwatch()
	return c
}

func (c *DirectGridClient) Send(buffer []byte) {
	c.Disp.Display(buffer)
}

func (c *DirectGridClient) NumLeds() int {
	return c.Disp.NumLeds()
}

func (c *DirectGridClient) Gamma() (r, g, b float64) {
	return c.Disp.Gamma()
}

func (c *DirectGridClient) SetGamma(r, g, b float64) {
	c.Disp.SetGamma(r, g, b)
}

func (c *DirectGridClient) ModuleConfig() conf.ModuleConfig {
	return c.Disp.ModuleConfig()
}

func (c *DirectGridClient) Stopwatch() *Stopwatch {
	return c.stopwatch
}

func (c *DirectGridClient) Close() {
	c.Disp.Close()
}


// Dieser Client-Typ schreibt alle Bilddaten in eine Datei, welche im
// Anschluss auf ein System mit echter Hardware transferiert und dort
// wie ein Film abgespielt wird.
type FileSaveClient struct {
	fh        *os.File
	modConf   conf.ModuleConfig
	stopwatch *Stopwatch
}

func NewFileSaveClient(fileName string, modConf conf.ModuleConfig) GridClient {
	var err error

	p := &FileSaveClient{}

	p.fh, err = os.Create(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	p.modConf = modConf
	p.stopwatch = NewStopwatch()

	return p
}

func (p *FileSaveClient) Send(buffer []byte) {
	if _, err := p.fh.Write(buffer); err != nil {
		log.Fatalf("Couldnt' write data to file: %v", err)
	}
}

func (p *FileSaveClient) NumLeds() int {
	return len(p.modConf) * (conf.ModuleDim.X * conf.ModuleDim.Y)
}

func (p *FileSaveClient) Gamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (p *FileSaveClient) SetGamma(r, g, b float64) {}

func (p *FileSaveClient) ModuleConfig() conf.ModuleConfig {
	return p.modConf
}

func (p *FileSaveClient) Stopwatch() *Stopwatch {
	return p.stopwatch
}

func (p *FileSaveClient) Close() {
	p.fh.Close()
}
