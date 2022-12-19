package pingpong

import (
	"bytes"
	"github.com/pedroAkos/go-simple-networking/pkg/neti"
)

type Ping struct {
}

func (p Ping) String() string {
	return p.Name()
}

func (p Ping) Name() string {
	return "Ping"
}

func (p Ping) Code() uint16 {
	return 1
}

func (p Ping) Serialize(buff *bytes.Buffer) error {
	return nil
}

func (p Ping) Deserialize(buff *bytes.Buffer) (neti.Message, error) {
	return Ping{}, nil
}

type Pong struct {
	Ping
}

func (p Pong) Code() uint16 {
	return 2
}

func (p Pong) Name() string {
	return "Pong"
}
