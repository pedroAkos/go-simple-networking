package neti

type NetClient interface {
	RegisterMessage(message Message)
	Recv() <-chan ReceivedMessage
	SendTo(conn HostConn, message Message) error
	Open(addr string) (conn HostConn, err error)
	Self() string
}

type NetService interface {
	RegisterListener(id uint16) NetClient
}

