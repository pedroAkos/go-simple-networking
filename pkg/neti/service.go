package neti

type NetClient interface {
	RegisterMessage(message Message)
	Recv() <-chan ReceivedMessage
	Send(conn HostConn, message Message)
	Self() string
}

type NetService interface {
	RegisterListener(id uint16) NetClient
}

