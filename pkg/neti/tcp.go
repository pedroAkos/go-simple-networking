package neti

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)


type tcpHostConn struct {
	conn net.Conn
	serviceId string
	sendLock *sync.Mutex
}

func (t tcpHostConn) ServiceId() string {
	return t.serviceId
}

func (t tcpHostConn) String() string {
	return t.Addr().String()
}

func (t tcpHostConn) Addr() net.Addr {
	return t.conn.RemoteAddr()
}

func (t tcpHostConn) Send(b []byte) error {
	t.sendLock.Lock()
	defer t.sendLock.Unlock()
	err := binary.Write(t.conn, binary.BigEndian, uint32(len(b)))
	if err != nil {
		return err
	}
	return writeFully(t.conn, b)
}


func (t tcpHostConn) Close() error {
	return t.conn.Close()
}

func (t tcpHostConn) Receive() ([]byte, error) {
	var size uint32
	err := binary.Read(t.conn, binary.BigEndian, &size)
	if err != nil {
		return nil, err
	}
	return readFully(t.conn, int(size))
}

func NewTcpNet(log *logrus.Logger) Net {
	return &tcp{
		listener:         nil,
		msgDeserializers: make(map[uint16]MessageDeserializer),
		log: log,
	}
}

type tcp struct {
	listener net.Listener
	msgDeserializers map[uint16]MessageDeserializer
	log *logrus.Logger
}

func (t tcp) RegisterMessage(message Message) {
	if _, ok := t.msgDeserializers[message.Code()]; !ok {
		t.msgDeserializers[message.Code()] = message.Deserialize
	} else {
		t.log.Warn(fmt.Sprintf("Message with Id: %v already registered, ignoring message: %v", message.Code(), message))
	}
}

func (t tcp) CloseListener() error {
	return t.listener.Close()
}

func (t tcp) Open(addr string) (HostConn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	hConn := tcpHostConn{conn: conn, sendLock: &sync.Mutex{}}

	return hConn, err
}

func (t tcp) OpenAsync(addr string, ch chan<- ReceivedConnection) {
	go func() {
		conn, err := t.Open(addr)
		ch <- ReceivedConnection{
			Addr: addr,
			Conn:  conn,
			Err: err,
		}
	}()
}


func (t tcp) recvAndDeserialize(conn HostConn) (Message, error) {
	b, err := conn.Receive()
	if err != nil {
		return nil, err
	}
	code := binary.BigEndian.Uint16(b)
	if d, ok := t.msgDeserializers[code]; ok {
		return d(b[binary.Size(code):])
	}
	return nil, errors.New(fmt.Sprintln("Unknown Msg code", code))
}

func (t tcp) RecvFrom(conn HostConn) (Message, error) {
	return t.recvAndDeserialize(conn)
}

func (t tcp) RecvFromAsync(conn HostConn, ch chan<- ReceivedMessage) {
	go func() {
		m, err := t.recvAndDeserialize(conn)
		ch <- ReceivedMessage{
			Conn:  conn,
			Msg:   m,
			Err: err,
		}
	}()
}

func (t tcp) SendTo(conn HostConn, message Message) error {
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
	t.log.WithFields(logrus.Fields{
		"Msg": message,
		"to": conn.Addr().String(),
		"size": len(payload),
	}).Debug("Sending")
	return conn.Send(buf.Bytes())
}

func (t tcp) SendToAsync(conn HostConn, message Message, ch chan<- SentMessage) {
	go func() {
		err := t.SendTo(conn, message)
		ch <- SentMessage{conn, message, err}
	}()
}

func (t *tcp) Listen(addr string) (<-chan HostConn, error) {

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	ch := make(chan HostConn)
	t.listener = listener
	go func() {
		for {
			conn, err := t.listener.Accept()
			if err != nil {
				fmt.Println("Error on accept", err.Error())
				close(ch)
				return
			} else {
				ch <- tcpHostConn{
					conn: conn,
					sendLock: &sync.Mutex{},
				}
			}
		}
	}()
	return ch, err
}
