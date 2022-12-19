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

// NetClient is an interface for a network client for a NetService.
// Register the messages that the client can receive and send.
// It can be used to send and receive messages.
// It can be used to connect to other hosts.
type NetClient interface {
	RegisterMessage(message Message)                         //Register Message in the NetClient (Known how to deserialize)
	RecvFrom(conn *ServiceHostConn) (Message, error)         //Receive Message from ServiceHostConn
	SendTo(conn *ServiceHostConn, message Message) error     //Send Message to ServiceHostConn
	OpenTo(addr string, id string) (*ServiceHostConn, error) //Open ServiceHostConn to peer with addr to NetClient NodeID
	Accept() <-chan *ServiceHostConn                         //Accept Connection
	Self() string                                            //Self Address
	Type() TransportType                                     //Transport Type
	Id() string                                              //NodeID of NetClient
}

// NetService is an interface for a network service for a NetClient.
// Register in the NetService the client id to get a NetClient.
// NetService multiplexes the connections to the NetClient.
type NetService interface {
	RegisterListener(id string) NetClient
	GetConfiguration() Configuration
}

// ServiceHostConn is a HostConn that multiplexes the connections to the NetClient.
type ServiceHostConn struct {
	Conn      HostConn
	ServiceId string
	Msg       Message
}

// String returns the string representation of the ServiceHostConn.
func (s *ServiceHostConn) String() string {
	return fmt.Sprintf("%v %v", s.ServiceId, s.Conn.String())
}

// Addr returns the address of the ServiceHostConn.
func (s *ServiceHostConn) Addr() net.Addr {
	return s.Conn.Addr()
}

// Send sends the  bytes to the Host on the other end of the ServiceHostConn.
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

// Receive receives the bytes from the Host on the other end of the ServiceHostConn.
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

// Close closes the ServiceHostConn.
func (s *ServiceHostConn) Close() error {
	return s.Conn.Close()
}
