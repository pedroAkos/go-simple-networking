package neti

import (
	"fmt"
	"github.com/pkg/errors"
	"net"
)

type simAddr struct {
	id string
}

func (s simAddr) Network() string {
	panic("implement me")
}

func (s simAddr) String() string {
	return s.id
}

type simConn struct {
	id   string
	addr simAddr
}

func (s *simConn) String() string {
	return s.id
}

func (s *simConn) Addr() net.Addr {
	return s.addr
}

func (s *simConn) Send(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s *simConn) Receive() ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *simConn) Close() error {
	//noop
	return nil
}

type simClient struct {
	id       string
	service  *simService
	listenCh chan *ServiceHostConn
}

func (s *simClient) RegisterMessage(message Message) {
	//noop
}

func (s *simClient) RecvFrom(conn *ServiceHostConn) (Message, error) {
	if conn.Msg != nil {
		defer func() { conn.Msg = nil }()
		return conn.Msg, nil
	}
	return nil, errors.New(fmt.Sprintf("Nothing to receive from connection"))
}

func (s *simClient) SendTo(conn *ServiceHostConn, message Message) error {
	c := &ServiceHostConn{conn.Conn, s.id, message}
	go s.service.deliver(c, s.id)
	return nil
}

func (s *simClient) OpenTo(addr string, id string) (*ServiceHostConn, error) {
	conn := &ServiceHostConn{&simConn{addr, simAddr{addr}}, id, nil}
	return conn, nil
}

func (s *simClient) Accept() <-chan *ServiceHostConn {
	return s.listenCh
}

func (s *simClient) Self() string {
	return s.id
}

func (s *simClient) Type() TransportType {
	return UDP
}

func (s *simClient) Id() string {
	return s.id
}

func (s *simClient) deliver(conn *ServiceHostConn) {
	s.listenCh <- conn
}

type simService struct {
	protos     map[string]uint64
	listenners map[string]*simClient
}

// NewSimUDPService creates a new SimUDPService
// This is a service that can be used to simulate a network.
func NewSimUDPService() NetService {
	return &simService{
		protos:     make(map[string]uint64),
		listenners: make(map[string]*simClient),
	}

}

func (s *simService) RegisterListener(id string) NetClient {
	if _, ok := s.protos[id]; !ok {
		s.protos[id] = 0
	}
	seqnum := s.protos[id]
	s.protos[id] = seqnum + 1
	id = fmt.Sprintf("%v/%v", id, seqnum)
	client := s.newSimNetClient(id)
	s.listenners[id] = client
	return client
}

func (s *simService) newSimNetClient(id string) *simClient {
	return &simClient{id, s, make(chan *ServiceHostConn)}
}

func (s *simService) GetConfiguration() Configuration {
	//TODO implement me
	panic("implement me")
}

func (s *simService) deliver(conn *ServiceHostConn, sender_id string) {
	id := conn.Conn.String()
	conn.Conn = &simConn{
		id:   conn.ServiceId,
		addr: simAddr{conn.ServiceId},
	}
	s.listenners[id].deliver(conn)
}
