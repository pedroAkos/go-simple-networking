package neti

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)


type basicUpdClient struct {
	self * string
	id string
	net Net
	//rcv chan ReceivedMessage
	listenCh chan *ServiceHostConn

	buffered map[string][]ReceivedMessage
}

func (b *basicUpdClient) Accept() <-chan *ServiceHostConn {
	return b.listenCh
}

func (b *basicUpdClient) ServiceId() string {
	return b.id
}

func (b *basicUpdClient) Type() TransportType {
	return UDP
}

func (b *basicUpdClient) OpenTo(addr string, serviceId string) (*ServiceHostConn, error) {
	if conn, err :=  b.net.Open(addr); err != nil {
		return nil, err
	} else {
		return &ServiceHostConn{conn, serviceId, nil}, err
	}
}

func (b *basicUpdClient) RegisterMessage(message Message) {
	b.net.RegisterMessage(messageWrap{b.id, message})
}

func (b *basicUpdClient) RecvFrom(conn *ServiceHostConn) (Message, error) {
	if conn.Msg != nil {
		defer func() {conn.Msg = nil}()
		return conn.Msg, nil
	}
	return nil, errors.New(fmt.Sprintf("Nothing to receive from connection"))
}

func (b *basicUpdClient) SendTo(conn *ServiceHostConn, message Message) error {
	return b.net.SendTo(conn, messageWrap{id: b.id, msg: message})
}

func (b *basicUpdClient) Self() string {
	return *b.self
}

func createUpdClient(self *string, id string, net Net) *basicUpdClient {
	return &basicUpdClient{
		self: self,
		id:   id,
		net:  net,
		listenCh:  make(chan *ServiceHostConn),
		buffered: make(map[string][]ReceivedMessage),
	}
}

type basicUdpService struct {
	self string
	net Net
	listeners map[string] *basicUpdClient
}

func (b *basicUdpService) GetConfiguration() Configuration {
	parts := strings.Split(b.self, ":")
	port, _ := strconv.Atoi(parts[1])
	return Configuration{
		ip: parts[0],
		port: port,
	}
}

func (b *basicUdpService) RegisterListener(id string) NetClient {
	client := createUpdClient(&b.self, id, b.net)
	b.listeners[id] = client
	return client
}

func (b *basicUdpService) deliver(msg messageWrap, conn* ServiceHostConn, err error) error {
	if err != nil {
		return err
	}
	if c, ok := b.listeners[conn.ServiceId]; ok {
		conn.ServiceId = msg.id
		conn.Msg = msg.msg
		c.listenCh <- conn
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
		listeners: make(map[string]*basicUpdClient),
	}
	go func() {
		for {
			select {
			case c := <- listen:
				conn := &ServiceHostConn{Conn: c}
				msg, err := net.RecvFrom(conn)
				if err = service.deliver(msg.(messageWrap), conn, err); err != nil {
					panic(err)
				}
			}
		}
	}()

	return service
}

