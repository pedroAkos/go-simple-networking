package neti

import (
	"fmt"
	"github.com/pkg/errors"
)

type messageWrap struct {
	id uint16
	msg Message
}

func (m messageWrap) String() string {
	return m.msg.String()
}

func (m messageWrap) Name() string {
	return m.msg.Name()
}

func (m messageWrap) Code() uint16 {
	return m.msg.Code()
}

func (m messageWrap) Serialize() ([]byte, error) {
	return m.msg.Serialize()
}

func (m messageWrap) Deserialize(bytes []byte) (Message, error) {
	msg, err := m.msg.Deserialize(bytes)
	if err != nil {
		return nil, err
	}
	m.msg = msg
	return m, nil
}

type basicUpdClient struct {
	self * string
	id uint16
	net Net
	rcv chan ReceivedMessage
}

func (b basicUpdClient) RegisterMessage(message Message) {
	b.net.RegisterMessage(messageWrap{b.id, message})
}

func (b basicUpdClient) Recv() <-chan ReceivedMessage {
	return b.rcv
}

func (b basicUpdClient) Send(conn HostConn, message Message) {
	go b.net.SendTo(conn, message)
}

func (b basicUpdClient) Self() string {
	return *b.self
}


func createUpdClient(self *string, id uint16, net Net) *basicUpdClient {
	return &basicUpdClient{
		self: self,
		id:   id,
		net:  net,
		rcv:  make(chan ReceivedMessage),
	}
}

type basicUdpService struct {
	self string
	net Net
	listeners map[uint16] *basicUpdClient
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
		c.rcv <- ReceivedMessage{
			Conn: conn,
			Msg:  msg.msg,
			Err:  err,
		}
		return nil
	}
	return errors.New(fmt.Sprintf("Listener with id %d is not registered", msg.id))
}


func Init(listenAddr string, buffsize int) NetService {
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

