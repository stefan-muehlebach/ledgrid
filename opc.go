//go:build ignore

package ledgrid

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

var (
	DispList []Displayer
)

func init() {
	DispList = make([]Displayer, 64)
}

func RegisterDisplayer(ch int, disp Displayer) {
	DispList[ch] = disp
}

type Message struct {
	Channel, Command byte
	Length           uint16
	Data             []byte
}

func NewMessage(b []byte) *Message {
	m := &Message{}
	m.Channel = b[0]
	m.Command = b[1]
	m.Length = (uint16(b[2]) << 8) | uint16(b[3])
	m.Data = make([]byte, m.Length)
	copy(m.Data, b[4:])
	return m
}

func HandleOPC(port uint) {
	hostPort := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", hostPort)
	if err != nil {
		log.Fatalf("Failed Listen(): %v", err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Failed Accept(): %v", err)
		}
		go HandleMessage(conn)
	}
}

func HandleMessage(conn net.Conn) {
	var buffer []byte

	defer conn.Close()
	buffer = make([]byte, 65539)
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				return
			}
			log.Fatalf("Failed Read(): %v", err)
		}
		m := NewMessage(buffer)
		disp := DispList[m.Channel]
		if disp == nil {
			log.Printf("No displayer on channel %d\n", m.Channel)
			continue
		}
		disp.Display(m.Data)
	}
}
