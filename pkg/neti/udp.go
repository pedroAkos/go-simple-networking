package neti

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

type udpHostConn struct {
	b []byte
	addr net.Addr
	conn   net.PacketConn
}

func (u udpHostConn) String() string {
	return u.addr.String()
}

func (u udpHostConn) Addr() net.Addr {
	return u.addr
}

func (u udpHostConn) Send(b []byte) error {
	n, err := u.conn.WriteTo(b, u.addr)
	if n != len(b) && err == nil {
		return errors.New(fmt.Sprint("Expected to send ", len(b), " bytes, sent ", n))
	}
	if err !=  nil {
		return err
	}
	return nil
}

func (u udpHostConn) Receive() ([]byte, error) {
	return u.b, nil
}


func (u udpHostConn) Close() error {
	return nil
}

func NewUdpNet(buffsize int) Net {
	return &udp{
		conn: nil,
		msgDeserializers: make(map[uint16]MessageDeserializer),
		buffsize: buffsize,
	}
}

type udp struct {
	conn             net.PacketConn
	msgDeserializers map[uint16]MessageDeserializer
	buffsize         int
}

func (u udp) RegisterMessage(message Message) {
	if _, ok := u.msgDeserializers[message.Code()]; !ok {
		u.msgDeserializers[message.Code()] = message.Deserialize
	}
}

func (u*udp) Listen(addr string) (<-chan HostConn, error) {
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return nil, err
	}
	u.conn = conn
	ch := make(chan HostConn)
	go func() {
		for {
			p := make([]byte, u.buffsize)
			_, addr, err := u.conn.ReadFrom(p)
			if err != nil {
				fmt.Println("Error on receive", err.Error())
				close(ch)
				return
			} else {
				ch <- udpHostConn{
					conn: u.conn,
					addr: addr,
					b: p,
				}
			}
		}
	}()

	return ch, err
}

func (u udp) CloseListener() error {
	return u.conn.Close()
}

func (u udp) Open(addr string) (HostConn, error) {
	if u.conn == nil {
		return nil, errors.New("no socket ready for UDP, call Listen first")
	}

	_addr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}
	return udpHostConn{
		addr: _addr,
		conn:  u.conn,
		b: nil,
	}, nil

}

func (u udp) OpenAsync(addr string, ch chan<- ReceivedConnection) {
	go func() {
		h, err := u.Open(addr)
		ch <- ReceivedConnection{addr, h, err}
	}()
}

func (u udp) recvAndDeserialize(conn HostConn) (Message, error) {
	b, err := conn.Receive()
	code := binary.BigEndian.Uint16(b)
	if err != nil {
		return nil, err
	}
	if d, ok := u.msgDeserializers[code]; ok {
		return d(b[binary.Size(code):])
	}
	return nil, errors.New(fmt.Sprintln("Unknown msg code", code))
}

func (u udp) RecvFromAsync(conn HostConn, ch chan<- ReceivedMessage) {
	go func() {
		m, err := u.recvAndDeserialize(conn)
		ch <- ReceivedMessage{
			Conn: &conn,
			Msg:  m,
			Err:  err,
		}
	}()
}

func (u udp) RecvFrom(conn HostConn) (Message, error) {
	return u.recvAndDeserialize(conn)
}


func (u udp) SendTo(conn HostConn, message Message) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, message.Code())
	if err != nil {
		return err
	}
	payload, err := message.Serialize()
	if err != nil {
		return err
	}
	err = writeFully(buf, payload)
	if err != nil {
		return err
	}
	return conn.Send(buf.Bytes())
}


func (u udp) SendToAsync(conn HostConn, message Message, ch chan<- SentMessage) {
	go func() {
		err := u.SendTo(conn, message)
		ch <- SentMessage{&conn, message, err}
	}()
}


