package neti

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

type basicTcpClient struct {
	self *string
	id   string
	net  Net
	rcv  chan ReceivedMessage
	acpt chan *ServiceHostConn

	msgs map[uint16]MessageDeserializer
}

func (b *basicTcpClient) Id() string {
	return b.id
}

func (b *basicTcpClient) Type() TransportType {
	return TCP
}

func (b *basicTcpClient) Accept() <-chan *ServiceHostConn {
	return b.acpt
}

func (b *basicTcpClient) RegisterMessage(message Message) {
	if _, ok := b.msgs[message.Code()]; ok {
		panic("Message already registered")
	} else {
		b.msgs[message.Code()] = message.Deserialize
	}
}

func (b *basicTcpClient) RecvFrom(conn *ServiceHostConn) (Message, error) {
	m, err := b.net.RecvFrom(conn)
	if err == nil {
		msg := m.(MessageWrap)
		if d, ok := b.msgs[msg.code]; ok {
			return d(msg.buff)
		} else {
			log.Warn("Unable to deserialize message with code ", msg.code, " for protocol ", b.id, ": Unknown serializer")
		}
		return nil, errors.New(fmt.Sprintf("Unable to deserialize message with code %v: Unknown serializer", msg.Code()))
	} else {
		return nil, err
	}
}

func (b *basicTcpClient) SendTo(conn *ServiceHostConn, message Message) error {
	return b.net.SendTo(conn, MessageWrap{Id: b.id, Msg: message})
}

func (b *basicTcpClient) OpenTo(addr string, id string) (*ServiceHostConn, error) {
	if conn, err := b.net.Open(addr); err == nil {

		buff := new(bytes.Buffer)
		if err = EncodeStringToBuffer(id, buff); err != nil {
			return nil, err
		}
		if err = EncodeStringToBuffer(b.id, buff); err != nil {
			return nil, err
		}
		if err = conn.Send(buff.Bytes()); err != nil {
			return nil, err
		}
		if _, err = conn.Receive(); err != nil {
			return nil, err
		}
		return &ServiceHostConn{ServiceId: id, Conn: conn}, err
	} else {
		return nil, err
	}
}

func (b *basicTcpClient) Self() string {
	return *b.self
}

func createTcpClient(self *string, id string, net Net) *basicTcpClient {
	return &basicTcpClient{
		self: self,
		id:   id,
		net:  net,
		rcv:  make(chan ReceivedMessage),
		acpt: make(chan *ServiceHostConn),
		msgs: make(map[uint16]MessageDeserializer),
	}
}

type basicTcpService struct {
	self      string
	net       Net
	listeners map[string]*basicTcpClient

	logger *log.Logger
}

func (b *basicTcpService) GetConfiguration() Configuration {
	parts := strings.Split(b.self, ":")
	port, _ := strconv.Atoi(parts[1])
	return Configuration{
		ip:   parts[0],
		port: port,
	}
}

func (b *basicTcpService) RegisterListener(id string) NetClient {
	client := createTcpClient(&b.self, id, b.net)
	b.listeners[id] = client
	return client
}

func (b *basicTcpService) accept(bid []byte, conn HostConn) {
	b.logger.Debug("Accepting: ", conn)
	buff := bytes.NewBuffer(bid)
	if id, err := DecodeStringFromBuffer(buff); err != nil {
		b.logger.Error(err)
	} else if c, ok := b.listeners[id]; ok {
		if id, err = DecodeStringFromBuffer(buff); err != nil {
			b.logger.Error(err)
		} else {
			if err := conn.Send([]byte{}); err != nil {
				b.logger.Error(err)
				_ = conn.Close()
			} else {
				c.acpt <- &ServiceHostConn{conn, id, nil}
			}
		}
	} else {
		_ = conn.Close()
		b.logger.Error("Received connection for ", id, ", but don't have listener registered")
	}
}

// InitBaseTcpService creates a new basic tcp service
func InitBaseTcpService(listenAddr string, logger *log.Logger) NetService {
	net := NewTcpNet(logger)
	net.RegisterMessage(MessageWrap{})
	listen, err := net.Listen(listenAddr)
	if err != nil {
		panic(err)
	}

	service := &basicTcpService{
		self:      listenAddr,
		net:       net,
		listeners: make(map[string]*basicTcpClient),
		logger:    logger,
	}
	go func() {
		for {
			select {
			case conn := <-listen:
				if bid, err := conn.Receive(); err != nil {
					logger.Error(err)
				} else {
					go service.accept(bid, conn)
				}
			}
		}
	}()
	return service
}
