package neti

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type basicUpdClient struct {
	self *string
	id   string
	net  Net
	//rcv chan ReceivedMessage
	listenCh chan *ServiceHostConn

	buffered map[string][]ReceivedMessage

	msgs map[uint16]MessageDeserializer
}

func (b *basicUpdClient) Accept() <-chan *ServiceHostConn {
	return b.listenCh
}

func (b *basicUpdClient) Id() string {
	return b.id
}

func (b *basicUpdClient) Type() TransportType {
	return UDP
}

func (b *basicUpdClient) OpenTo(addr string, serviceId string) (*ServiceHostConn, error) {
	if conn, err := b.net.Open(addr); err != nil {
		return nil, err
	} else {
		return &ServiceHostConn{conn, serviceId, nil}, err
	}
}

func (b *basicUpdClient) RegisterMessage(message Message) {
	if _, ok := b.msgs[message.Code()]; ok {
		panic("Message already registered")
	} else {
		b.msgs[message.Code()] = message.Deserialize
	}
}

func (b *basicUpdClient) RecvFrom(conn *ServiceHostConn) (Message, error) {
	if conn.Msg != nil {
		defer func() { conn.Msg = nil }()
		return conn.Msg, nil
	}
	return nil, errors.New(fmt.Sprintf("Nothing to receive from connection"))
}

func (b *basicUpdClient) SendTo(conn *ServiceHostConn, message Message) error {
	return b.net.SendTo(conn, MessageWrap{Id: b.id, Msg: message})
}

func (b *basicUpdClient) Self() string {
	return *b.self
}

func (b *basicUpdClient) deliver(msg MessageWrap, conn *ServiceHostConn) {
	conn.ServiceId = msg.Id
	if d, ok := b.msgs[msg.code]; ok {
		conn.Msg, _ = d(msg.buff)
		b.listenCh <- conn
	} else {
		log.Warn("Unable to deserialize message with code ", msg.code, " for protocol ", b.id, ": Unknown serializer")
	}
}

func createUpdClient(self *string, id string, net Net) *basicUpdClient {
	return &basicUpdClient{
		self:     self,
		id:       id,
		net:      net,
		listenCh: make(chan *ServiceHostConn),
		buffered: make(map[string][]ReceivedMessage),
		msgs:     make(map[uint16]MessageDeserializer),
	}
}

type basicUdpService struct {
	self      string
	net       Net
	listeners map[string]*basicUpdClient
}

func (b *basicUdpService) GetConfiguration() Configuration {
	parts := strings.Split(b.self, ":")
	port, _ := strconv.Atoi(parts[1])
	return Configuration{
		ip:   parts[0],
		port: port,
	}
}

func (b *basicUdpService) RegisterListener(id string) NetClient {
	client := createUpdClient(&b.self, id, b.net)
	b.listeners[id] = client
	return client
}

func (b *basicUdpService) deliver(msg MessageWrap, conn *ServiceHostConn, err error) error {
	if err != nil {
		return err
	}
	if c, ok := b.listeners[conn.ServiceId]; ok {
		go c.deliver(msg, conn)
		return nil
	}
	return errors.New(fmt.Sprintf("Listener with Id %v is not registered, entries: %v", conn.ServiceId, b.listeners))
}

// InitBaseUdpService creates a new basicUdpService
func InitBaseUdpService(listenAddr string, buffsize int) NetService {
	net := NewUdpNet(buffsize)
	net.RegisterMessage(MessageWrap{})
	listen, err := net.Listen(listenAddr)
	if err != nil {
		panic(err)
	}
	service := &basicUdpService{
		self:      listenAddr,
		net:       net,
		listeners: make(map[string]*basicUpdClient),
	}
	go func(listen <-chan HostConn, net Net, service *basicUdpService) {
		for {
			select {
			case c := <-listen:
				conn := &ServiceHostConn{Conn: c}
				if msg, err := net.RecvFrom(conn); err != nil {
					panic(err)
				} else if err = service.deliver(msg.(MessageWrap), conn, err); err != nil {
					//panic(err)
					log.Warn(err)
				}
			}
		}
	}(listen, net, service)

	return service
}
