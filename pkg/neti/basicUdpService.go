package neti

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)


type basicUpdClient struct {
	self * string
	id uint16
	net Net
	//rcv chan ReceivedMessage
	listenCh chan HostConn

	buffered map[string][]ReceivedMessage
}

func (b *basicUpdClient) Type() TransportType {
	return UDP
}

func (b *basicUpdClient) Accept() <-chan HostConn {
	return b.listenCh
}

func (b *basicUpdClient) Open(addr string) (conn HostConn, err error) {
	return b.net.Open(addr)
}

func (b *basicUpdClient) RegisterMessage(message Message) {
	b.net.RegisterMessage(messageWrap{b.id, message})
}

func (b *basicUpdClient) RecvFrom(conn HostConn) (Message, error) {
	if _, ok := b.buffered[conn.Addr().String()]; !ok {
		return nil, errors.New(fmt.Sprintf("Nothing to receive from connection"))
	}
	m := b.buffered[conn.Addr().String()][0]
	b.buffered[conn.Addr().String()] = append(b.buffered[conn.Addr().String()][:0],
		b.buffered[conn.Addr().String()][0:]...)
	return m.Msg, m.Err
}

func (b *basicUpdClient) SendTo(conn HostConn, message Message) error {
	return b.net.SendTo(conn, messageWrap{id: b.id, msg: message})
}

func (b *basicUpdClient) Self() string {
	return *b.self
}

func (b *basicUpdClient) addmsg(msg Message, conn *HostConn, err error) {
	if _, ok := b.buffered[(*conn).Addr().String()]; !ok {
		b.buffered[(*conn).Addr().String()] = make([]ReceivedMessage, 1)
		b.buffered[(*conn).Addr().String()][0] = ReceivedMessage{
			Conn: conn,
			Msg:  msg,
			Err:  err,
		}
	} else {
		b.buffered[(*conn).Addr().String()] = append(b.buffered[(*conn).Addr().String()], ReceivedMessage{
			Conn: conn,
			Msg:  msg,
			Err:  err,
		})
	}
	b.listenCh <- *conn
}


func createUpdClient(self *string, id uint16, net Net) *basicUpdClient {
	return &basicUpdClient{
		self: self,
		id:   id,
		net:  net,
		listenCh:  make(chan HostConn),
	}
}

type basicUdpService struct {
	self string
	net Net
	listeners map[uint16] *basicUpdClient
}

func (b *basicUdpService) GetConfiguration() Configuration {
	parts := strings.Split(b.self, ":")
	port, _ := strconv.Atoi(parts[1])
	return Configuration{
		ip: parts[0],
		port: port,
	}
}

func (b *basicUdpService) RegisterListener(id uint16) NetClient {
	client := createUpdClient(&b.self, id, b.net)
	b.listeners[id] = client
	return client
}

func (b *basicUdpService) deliver(msg messageWrap, conn* HostConn, err error) error {
	if err != nil {
		return err
	}
	if c, ok := b.listeners[msg.id]; ok {
		c.addmsg(msg.msg, conn, err)
		return nil
	}
	return errors.New(fmt.Sprintf("Listener with id %d is not registered", msg.id))
}


func InitBaseUdpService(listenAddr string, buffsize int) NetService {
	net := NewUdpNet(buffsize)
	listen, err := net.Listen(listenAddr)
	if err != nil {
		panic(err)
	}
	service := &basicUdpService{
		self:      listenAddr,
		net:       net,
		listeners: make(map[uint16]*basicUpdClient),
	}
	go func() {
		for {
			select {
			case conn := <- listen:
				msg, err := net.RecvFrom(conn)
				if err = service.deliver(msg.(messageWrap), &conn, err); err != nil {
					panic(err)
				}
			}
		}
	}()

	return service
}

