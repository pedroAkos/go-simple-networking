package main

import (
	"github.com/pedroAkos/network/pkg/neti"
)

type Ping struct {

}

func (p Ping) String() string {
	return "Ping"
}

func (p Ping) Name() string {
	return "Ping"
}

func (p Ping) Code() uint16 {
	return 1
}

func (p Ping) Serialize() ([]byte, error) {
	b := []byte{}
	return b, nil
}

func (p Ping) Deserialize(bytes []byte) (neti.Message, error) {
	return Ping{}, nil
}


type Pong struct {
	Ping
}

func (p Pong) Code() uint16 {
	return 2
}

func (p Pong) String() string {
	return "Pong"
}

func (p Pong) Name() string {
	return "Pong"
}

func (p Pong) Deserialize(bytes []byte) (neti.Message, error) {
	return Pong{}, nil
}



