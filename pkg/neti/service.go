package neti

import (
	"bytes"
	"fmt"
	"net"
)

type TransportType uint8

const (
	UDP TransportType = 1
	TCP TransportType = 2
)

type NetClient interface {
	RegisterMessage(message Message)                         //Register Message in the NetClient (Known how to deserialize)
	RecvFrom(conn *ServiceHostConn) (Message, error)         //Receive Message from ServiceHostConn
	SendTo(conn *ServiceHostConn, message Message) error     //Send Message to ServiceHostConn
	OpenTo(addr string, id string) (*ServiceHostConn, error) //Open ServiceHostConn to peer with addr to NetClient ID
	Accept() <-chan *ServiceHostConn                         //Accept Connection
	Self() string                                            //Self Address
	Type() TransportType                                     //Transport Type
	Id() string                                              //ID of NetClient
}

type NetService interface {
	RegisterListener(id string) NetClient
	GetConfiguration() Configuration
}

type ServiceHostConn struct {
	Conn      HostConn
	ServiceId string
	Msg       Message
}

func (s *ServiceHostConn) String() string {
	return fmt.Sprintf("%v %v", s.ServiceId, s.Conn.String())
}

func (s *ServiceHostConn) Addr() net.Addr {
	return s.Conn.Addr()
}

func (s *ServiceHostConn) Send(b []byte) error {
	buff := new(bytes.Buffer)
	if err := EncodeStringToBuffer(s.ServiceId, buff); err != nil {
		return err
	}
	if err := EncodeBytesToBuffer(b, buff); err != nil {
		return err
	}
	return s.Conn.Send(buff.Bytes())
}

func (s *ServiceHostConn) Receive() ([]byte, error) {
	b, err := s.Conn.Receive()
	buff := bytes.NewBuffer(b)
	if s.ServiceId, err = DecodeStringFromBuffer(buff); err != nil {
		return nil, err
	}
	if s.ServiceId == "" {
		panic("Service Id is empty")
	}
	return DecodeBytesFromBuffer(buff)
}

func (s *ServiceHostConn) Close() error {
	return s.Conn.Close()
}
