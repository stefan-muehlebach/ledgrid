package ledgrid

import (
	"fmt"
	"log"
	"net"
)

// Dieser Typ wird client-seitig fuer die Ansteuerung des LedGrid verwendet.
// Im Wesentlichen ist dies eine Abstraktion der Ansteuerung via UDP.
type PixelCtrl struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

// Erzeugt ein neues Controller-Objekt, welches das LedGrid ueber die Adresse
// in Host und den UDP-Port in Port anspricht.
func NewPixelCtrl(host string, port uint) *PixelCtrl {
	var hostPort string
	var err error

	p := &PixelCtrl{}
	hostPort = fmt.Sprintf("%s:%d", host, port)
	p.addr, err = net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		log.Fatal(err)
	}
	p.conn, err = net.DialUDP("udp", nil, p.addr)
	if err != nil {
		log.Fatal(err)
	}
	return p
}

// Schliesst die Verbindung zum Controller.
func (p *PixelCtrl) Close() {
    p.conn.Close()
}

// Sendet die Daten im Buffer b zum Controller.
func (p *PixelCtrl) Send(b []byte) {
    var err error

    	_, err = p.conn.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}
