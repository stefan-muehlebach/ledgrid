package ledgrid

import (
	"github.com/stefan-muehlebach/ledgrid/conf"
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
	Send(buffer []byte)
	NumLeds() int
	Gamma() (r, g, b float64)
	SetGamma(r, g, b float64)
	MaxBright() (r, g, b uint8)
	SetMaxBright(r, g, b uint8)
    ModuleConfig() (conf.ModuleConfig)
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
		log.Fatal("Error in Dial(dataPort): %v", err)
	}

	p.rpcClient, err = rpc.DialHTTP("tcp", hostPortRPC)
	if err != nil {
		log.Fatal("Error in Dial(rpcPort): %v", err)
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

func (p *NetGridClient) NumLeds() int {
	var reply NumLedsArg
	var err error

	err = p.rpcClient.Call("GridServer.RPCNumLeds", 0, &reply)
	if err != nil {
		log.Fatal("Size error:", err)
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

func (p *NetGridClient) ModuleConfig() (conf.ModuleConfig) {
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

// // Mit diesem Typ wird die klassische Verwendung auf zwei Nodes realisiert.
// type OPCGridClient struct {
// 	conn      net.Conn
// 	rpcClient *rpc.Client
// 	sendWatch *Stopwatch
// 	buffer    []byte
// }

// func NewOPCGridClient(host string, dataPort, rpcPort uint) GridClient {
// 	var hostPortData, hostPortRPC string
// 	var err error

// 	p := &OPCGridClient{}
// 	hostPortData = fmt.Sprintf("%s:%d", host, dataPort)
// 	hostPortRPC = fmt.Sprintf("%s:%d", host, rpcPort)
// 	p.conn, err = net.Dial("tcp", hostPortData)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	p.rpcClient, err = rpc.DialHTTP("tcp", hostPortRPC)
// 	if err != nil {
// 		log.Fatal("Dialing:", err)
// 	}
// 	p.sendWatch = NewStopwatch()
// 	p.buffer = make([]byte, 65539)

// 	return p
// }

// // Sendet die Bilddaten in der LedGrid-Struktur zum Controller.
// func (p *OPCGridClient) Send(buffer []byte) {
// 	var err error

// 	p.sendWatch.Start()
// 	length := len(buffer)
// 	p.buffer[0] = 0
// 	p.buffer[1] = 0
// 	p.buffer[2] = byte((length >> 8) & 0xff)
// 	p.buffer[3] = byte(length & 0xff)
// 	copy(p.buffer[4:], buffer)
// 	_, err = p.conn.Write(p.buffer[:length+4])
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	p.sendWatch.Stop()
// }

// func (p *OPCGridClient) NumLeds() int {
// 	var reply NumLedsArg
// 	var err error

// 	err = p.rpcClient.Call("GridServer.RPCNumLeds", 0, &reply)
// 	if err != nil {
// 		log.Fatal("Size error:", err)
// 	}
// 	return int(reply)
// }

// func (p *OPCGridClient) Gamma() (r, g, b float64) {
// 	var reply GammaArg
// 	var err error

// 	err = p.rpcClient.Call("GridServer.RPCGamma", 0, &reply)
// 	if err != nil {
// 		log.Fatal("Gamma error:", err)
// 	}
// 	return reply.RedVal, reply.GreenVal, reply.BlueVal
// }

// func (p *OPCGridClient) SetGamma(r, g, b float64) {
// 	var reply int
// 	var err error

// 	err = p.rpcClient.Call("GridServer.RPCSetGamma", GammaArg{r, g, b}, &reply)
// 	if err != nil {
// 		log.Fatal("SetGamma error:", err)
// 	}
// }

// func (p *OPCGridClient) MaxBright() (r, g, b uint8) {
// 	var reply BrightArg
// 	var err error

// 	err = p.rpcClient.Call("GridServer.RPCMaxBright", 0, &reply)
// 	if err != nil {
// 		log.Fatal("MaxBright error:", err)
// 	}
// 	return reply.RedVal, reply.GreenVal, reply.BlueVal
// }

// func (p *OPCGridClient) SetMaxBright(r, g, b uint8) {
// 	var reply int
// 	var err error

// 	err = p.rpcClient.Call("GridServer.RPCSetMaxBright", BrightArg{r, g, b}, &reply)
// 	if err != nil {
// 		log.Fatal("SetMaxBright error:", err)
// 	}
// }

// func (p *OPCGridClient) Watch() *Stopwatch {
// 	return p.sendWatch
// }

// // Schliesst die Verbindung zum Controller.
// func (p *OPCGridClient) Close() {
// 	p.conn.Close()
// }

// Mit dieser Implementation des GridClient-Interfaces kann man ohne Zugriff
// auf ein reales LED-Grid Software testen.
type DummyGridClient struct {
	size image.Point
}

func NewDummyGridClient(size image.Point) GridClient {
	p := &DummyGridClient{size: size}
	return p
}

func (p *DummyGridClient) Send(buffer []byte) {}

func (p *DummyGridClient) NumLeds() int {
	return p.size.X * p.size.Y
}

func (p *DummyGridClient) Gamma() (r, g, b float64) {
	return 1.0, 1.0, 1.0
}

func (p *DummyGridClient) SetGamma(r, g, b float64) {}

func (p *DummyGridClient) MaxBright() (r, g, b uint8) {
	return 0xff, 0xff, 0xff
}

func (p *DummyGridClient) SetMaxBright(r, g, b uint8) {}

func (p *DummyGridClient) ModuleConfig() (conf.ModuleConfig) {
    return conf.DefaultModuleConfig(p.size)
}

func (p *DummyGridClient) Watch() *Stopwatch {
	// TO DO: even a dummy implementation of the client should return a
	// usable Stopwatch. Otherwise, the calling function may crash if we just
	// return nil...
	return nil
}

func (p *DummyGridClient) Close() {}
