package ledgrid

import (
	"fmt"
	"errors"
	"log"
	"net"
)

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

func HandleOPC() {
	l, err := net.Listen("tcp", ":5333")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleMessage(conn)
	}
}

func HandleMessage(conn net.Conn) {
	var buffer []byte

	buffer = make([]byte, 65539)
	_, err := conn.Read(buffer)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			return
		}
		log.Fatal(err)
	}
    m := NewMessage(buffer)
    fmt.Printf("m: %+v\n", m)
    conn.Close()
}
