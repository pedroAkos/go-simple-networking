package neti

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type basicTcpClient struct {
	self * string
	id uint16
	net Net
	rcv chan ReceivedMessage
	acpt chan HostConn
}

func (b *basicTcpClient) Type() TransportType {
	return TCP
}

func (b *basicTcpClient) Accept() <-chan HostConn {
	return b.acpt
}

func (b *basicTcpClient) RegisterMessage(message Message) {
	b.net.RegisterMessage(messageWrap{b.id, message})
}

func (b *basicTcpClient) RecvFrom(conn HostConn) (Message, error) {
	return b.net.RecvFrom(conn)
}

func (b *basicTcpClient) SendTo(conn HostConn, message Message) error {
	return b.net.SendTo(conn, messageWrap{id: b.id, msg: message})
}

func (b *basicTcpClient) Open(addr string) (conn HostConn, err error) {
	return b.net.Open(addr)
}

func (b *basicTcpClient) Self() string {
	return *b.self
}


func createTcpClient(self *string, id uint16, net Net) *basicTcpClient {
	return &basicTcpClient{
		self: self,
		id:   id,
		net:  net,
		rcv:  make(chan ReceivedMessage),
		acpt: make(chan HostConn),
	}
}

type basicTcpService struct {
	self string
	net Net
	listeners map[uint16]  *basicTcpClient

	logger *log.Logger
}

func (b *basicTcpService) GetConfiguration() Configuration {
	parts := strings.Split(b.self, ":")
	port, _ := strconv.Atoi(parts[1])
	return Configuration{
		ip: parts[0],
		port: port,
	}
}

func (b *basicTcpService) RegisterListener(id uint16) NetClient {
	client := createTcpClient(&b.self, id, b.net)
	b.listeners[id] = client
	return client
}


func (b *basicTcpService) accept(bid []byte, conn HostConn) {
	var id uint16;
	if err := DecodeNumberFromBuffer(&id, bytes.NewBuffer(bid)); err != nil {
		b.logger.Error(err)
	} else if c, ok := b.listeners[id]; ok {
		c.acpt <- conn
	} else {
		b.logger.Error("Received connection for ", id, ", but don't have listener registered")
	}
}

func InitBaseTcpService(listenAddr string, logger *log.Logger) NetService {
	net := NewTcpNet(logger)
	listen, err := net.Listen(listenAddr)
	if err != nil {
		panic(err)
	}

	service := &basicTcpService{
		self:      listenAddr,
		net:       net,
		listeners: make(map[uint16]*basicTcpClient),
	}
	go func() {
		for {
			select {
			case conn := <- listen:
				if bid, err := conn.Receive();  err != nil {
					logger.Error(err)
				} else {
					go service.accept(bid, conn)
				}
			}
		}
	}()
	return service
}