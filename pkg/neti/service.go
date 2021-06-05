package neti

type TransportType uint8

const (
	UDP TransportType = 1
	TCP TransportType = 2
)

type NetClient interface {
	RegisterMessage(message Message)
	RecvFrom(conn HostConn) (Message, error)
	SendTo(conn HostConn, message Message) error
	Open(addr string) (conn HostConn, err error)
	Accept() <-chan HostConn
	Self() string
	Type() TransportType
}

type NetService interface {
	RegisterListener(id uint16) NetClient
	GetConfiguration() Configuration
}

