package neti

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type MessageDeserializer func(*bytes.Buffer) (Message, error)

// Message is the interface that wraps the basic Message methods.
type Message interface {
	String() string
	Name() string
	Code() uint16
	Serialize(buff *bytes.Buffer) error
	Deserialize(buff *bytes.Buffer) (Message, error)
}

// HostConn is the interface that wraps the basic HostConn methods.
// This represents a connection to a host.
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
	Msg  Message
	Err  error
}

// Net is the interface that wraps the basic Net methods.
// This represents a network. It can be a TCP, UDP, etc.
// It can be used to send and receive messages.
// It can be used to listen for new connections.
// It can be used to connect to other hosts.
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
		//log.Warn("Expected to read ", toRead, " read ", n)
		log.Panicln("Expected to read ", toRead, " read ", n)
	}
	return b, err
}

// GetInterfaceIpv4Addr returns the IPv4 address of the interface with the given name.
// taken from https://gist.github.com/schwarzeni/f25031a3123f895ff3785970921e962c
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
