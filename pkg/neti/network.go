package neti

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type MessageDeserializer func([]byte) (Message, error)

type Message interface {
	String() string
	Name() string
	Code() uint16
	Serialize() ([]byte, error)
	Deserialize([]byte) (Message, error)
}

type HostConn interface {
	fmt.Stringer
	Addr() net.Addr
	Send([]byte) error
	Receive() ([]byte, error)
	Close() error
}

type ReceivedMessage struct {
	Conn HostConn
	Msg  Message
	Err  error
}

type ReceivedConnection struct {
	Addr string
	Conn HostConn
	Err  error
}

type SentMessage struct {
	Conn HostConn
	Msg Message
	Err error
}

type Net interface {
	RegisterMessage(message Message)
	Listen(addr string) (<-chan HostConn, error)
	CloseListener() error
	Open(addr string) (HostConn, error)
	OpenAsync(addr string, ch chan<- ReceivedConnection)
	RecvFromAsync(conn HostConn, ch chan<- ReceivedMessage)
	RecvFrom(conn HostConn) (Message, error)
	SendTo(conn HostConn, m Message) error
	SendToAsync(conn HostConn, m Message, ch chan<- SentMessage)
}


func writeFully(writer io.Writer, b []byte) error {
	total := len(b)
	n, err := writer.Write(b)
	if n != total {
		log.Warn("Expected to write ", total, " wrote ", n)
	}
	return err
}

func readFully(reader io.Reader, toRead int) ([]byte, error) {
	b := make([]byte, toRead)
	n, err := reader.Read(b)
	if n != toRead {
		log.Warn("Expected to read ", toRead, " read ", n)
	}
	return b, err
}


//taken from https://gist.github.com/schwarzeni/f25031a3123f895ff3785970921e962c
func GetInterfaceIpv4Addr(interfaceName string) (addr string, err error) {
	var (
		ief      *net.Interface
		addrs    []net.Addr
		ipv4Addr net.IP
	)
	if ief, err = net.InterfaceByName(interfaceName); err != nil { // get interface
		return
	}
	if addrs, err = ief.Addrs(); err != nil { // get addresses
		return
	}
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		return "", errors.New(fmt.Sprintf("interface %s don't have an ipv4 address\n", interfaceName))
	}
	return ipv4Addr.String(), nil
}
