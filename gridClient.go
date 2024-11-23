package ledgrid

import (
	"fmt"
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
//	NetGridClient  - Verbindet sich via UDP und RPC mit einem externen
//	                 gridController.
//  FileSaveClient - Schreibt die Bilddaten in ein File, welches dann auf das
//                   System mit dem Grid-Controller kopiert und dort direkt
//                   abgespielt werden kann.
type GridClient interface {
	Send(buffer []byte)
	NumLeds() int
	Gamma() (r, g, b float64)
	SetGamma(r, g, b float64)
	MaxBright() (r, g, b uint8)
	SetMaxBright(r, g, b uint8)
	ModuleConfig() conf.ModuleConfig
	Watch() *Stopwatch
	Close()
}

// Mit diesem Typ wird die klassische Verwendung auf zwei Nodes realisiert.
type NetGridClient struct {
	conn      net.Conn
	rpcClient *rpc.Client
	sendWatch *Stopwatch
}

func NewNetGridClient(host string, network string, port, rpcPort uint) GridClient {
	var hostPortData, hostPortRPC string
	var err error

	p := &NetGridClient{}
	hostPortData = fmt.Sprintf("%s:%d", host, port)
	hostPortRPC = fmt.Sprintf("%s:%d", host, rpcPort)

	p.conn, err = net.Dial(network, hostPortData)
	if err != nil {
		log.Fatalf("Error in Dial(dataPort): %v", err)
	}

	p.rpcClient, err = rpc.DialHTTP("tcp", hostPortRPC)
	if err != nil {
		log.Fatalf("Error in Dial(rpcPort): %v", err)
	}
	p.sendWatch = NewStopwatch()

	return p
}

// Sendet die Bilddaten in der LedGrid-Struktur zum Controller.
func (p *NetGridClient) Send(buffer []byte) {
	var err error

	p.sendWatch.Start()
	_, err = p.conn.Write(buffer)
	if err != nil {
		log.Fatal(err)
	}
	p.sendWatch.Stop()
}

// Die folgenden Methoden verpacken die entsprechenden RPC-Calls zum
// Grid-Server.
func (p *NetGridClient) NumLeds() int {
	var reply NumLedsArg
	var err error

	err = p.rpcClient.Call("GridServer.RPCNumLeds", 0, &reply)
	if err != nil {
		log.Fatal("NumLeds error:", err)
	}
	return int(reply)
}

func (p *NetGridClient) Gamma() (r, g, b float64) {
	var reply GammaArg
	var err error

	err = p.rpcClient.Call("GridServer.RPCGamma", 0, &reply)
	if err != nil {
		log.Fatal("Gamma error:", err)
	}
	return reply.RedVal, reply.GreenVal, reply.BlueVal
}

func (p *NetGridClient) SetGamma(r, g, b float64) {
	var reply int
	var err error

	err = p.rpcClient.Call("GridServer.RPCSetGamma", GammaArg{r, g, b}, &reply)
	if err != nil {
		log.Fatal("SetGamma error:", err)
	}
}

func (p *NetGridClient) MaxBright() (r, g, b uint8) {
	var reply BrightArg
	var err error

	err = p.rpcClient.Call("GridServer.RPCMaxBright", 0, &reply)
	if err != nil {
		log.Fatal("MaxBright error:", err)
	}
	return reply.RedVal, reply.GreenVal, reply.BlueVal
}

func (p *NetGridClient) SetMaxBright(r, g, b uint8) {
	var reply int
	var err error

	err = p.rpcClient.Call("GridServer.RPCSetMaxBright", BrightArg{r, g, b}, &reply)
	if err != nil {
		log.Fatal("SetMaxBright error:", err)
	}
}

func (p *NetGridClient) ModuleConfig() conf.ModuleConfig {
	var reply ModuleConfigArg
	var err error

	err = p.rpcClient.Call("GridServer.RPCModuleConfig", 0, &reply)
	if err != nil {
		log.Fatal("ModuleConfig error:", err)
	}
	return reply.ModuleConfig
}

func (p *NetGridClient) Watch() *Stopwatch {
	return p.sendWatch
}

// Schliesst die Verbindung zum Controller.
func (p *NetGridClient) Close() {
	p.conn.Close()
}

// Dieser Client-Typ schreibt alle Bilddaten in eine Datei, welche im
// Anschluss auf ein System mit echter Hardware transferiert und dort
// wie ein Film abgespielt wird.
type FileSaveClient struct {
	fh        *os.File
	modConf   conf.ModuleConfig
	sendWatch *Stopwatch
}

func NewFileSaveClient(fileName string, modConf conf.ModuleConfig) GridClient {
	var err error

	p := &FileSaveClient{}

	p.fh, err = os.Create(fileName)
	if err != nil {
		log.Fatalf("Couldn't create file: %v", err)
	}
	p.modConf = modConf
	p.sendWatch = NewStopwatch()

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

func (p *FileSaveClient) MaxBright() (r, g, b uint8) {
	return 0xff, 0xff, 0xff
}

func (p *FileSaveClient) SetMaxBright(r, g, b uint8) {}

func (p *FileSaveClient) ModuleConfig() conf.ModuleConfig {
	return p.modConf
}

func (p *FileSaveClient) Watch() *Stopwatch {
	return p.sendWatch
}

func (p *FileSaveClient) Close() {
	p.fh.Close()
}
