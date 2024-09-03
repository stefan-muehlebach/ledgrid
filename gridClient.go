package ledgrid

import (
	"fmt"
	"image"
	"log"
	"net"
	"net/rpc"
)

// Um den clientseitigen Code so generisch wie moeglich zu halten, ist der
// GridClient als Interface definiert. Zwei konkrete Implementationen
// stehen zur Verfuegung:
// - NetGridClient
// - DummyGridClient
type GridClient interface {
	Send(lg *LedGrid)
	Size() image.Point
	Gamma() (r, g, b float64)
	SetGamma(r, g, b float64)
	MaxBright() (r, g, b uint8)
	SetMaxBright(r, g, b uint8)
	Watch() *Stopwatch
	Close()
}

// Mit diesem Typ wird die klassische Verwendung auf zwei Nodes realisiert.
type NetGridClient struct {
	addr      *net.UDPAddr
	conn      *net.UDPConn
	rpcClient *rpc.Client
	sendWatch *Stopwatch
}

func NewNetGridClient(host string, port uint) GridClient {
	var hostPort string
	var err error

	p := &NetGridClient{}
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

// Sendet die Bilddaten in der LedGrid-Struktur zum Controller.
func (p *NetGridClient) Send(lg *LedGrid) {
	var err error

	p.sendWatch.Start()
	_, err = p.conn.Write(lg.Pix)
	if err != nil {
		log.Fatal(err)
	}
	p.sendWatch.Stop()
}

func (p *NetGridClient) Size() (size image.Point) {
	var reply SizeArg
	var err error

	err = p.rpcClient.Call("GridServer.RPCSize", 0, &reply)
	if err != nil {
		log.Fatal("Size error:", err)
	}
	return image.Point(reply)
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

func (p *NetGridClient) Watch() *Stopwatch {
	return p.sendWatch
}

// Schliesst die Verbindung zum Controller.
func (p *NetGridClient) Close() {
	p.conn.Close()
}

// Mit dieser Implementation des GridClient-Interfaces kann man ohne Zugriff
// auf ein reales LED-Grid Software testen.
type DummyGridClient struct {
	size image.Point
}

func NewDummyGridClient(size image.Point) GridClient {
	p := &DummyGridClient{size: size}
	return p
}

func (p *DummyGridClient) Send(lg *LedGrid) { }

func (p *DummyGridClient) Size() image.Point {
	return p.size
}

func (p *DummyGridClient) Gamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (p *DummyGridClient) SetGamma(r, g, b float64) {}

func (p *DummyGridClient) MaxBright() (r, g, b uint8) {
	return 0xff, 0xff, 0xff
}

func (p *DummyGridClient) SetMaxBright(r, g, b uint8) {}

func (p *DummyGridClient) Watch() *Stopwatch {
    // TO DO: even a dummy implementation of the client should return a
    // usable Stopwatch. Otherwise, the calling function may crash if we just
    // return nil...
	return nil
}

func (p *DummyGridClient) Close() {}


