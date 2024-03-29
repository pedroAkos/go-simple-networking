package main

import (
	"bufio"
	"fmt"
	"github.com/pedroAkos/go-simple-networking/cmd/pingpong"
	"github.com/pedroAkos/go-simple-networking/pkg/neti"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
)

func main() {

	host := os.Args[1]
	target := os.Args[2]

	tcp := neti.NewTcpNet(logrus.StandardLogger())

	listener, err := tcp.Listen(host)
	if err != nil {
		log.Panicln(err.Error())
	}

	tcp.RegisterMessage(pingpong.Ping{})
	tcp.RegisterMessage(pingpong.Pong{})

	go func() {
		for {
			select {
			case conn := <-listener:
				go func() {
					for {
						msg, err := tcp.RecvFrom(conn)
						if err != nil {
							log.Panicln(err.Error())
						}
						log.Println(fmt.Sprintf("Received msg %s from conn %s", msg.String(), conn.Addr().String()))
						err = tcp.SendTo(conn, pingpong.Pong{})
						if err != nil {
							log.Panicln(err.Error())
						}
					}
				}()
			}
		}
	}()

	fmt.Println("Press to continue...")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	conn, err := tcp.Open(target)
	if err != nil {
		log.Panicln(err.Error())
	}

	for {
		err = tcp.SendTo(conn, pingpong.Ping{})
		if err != nil {
			log.Panicln(err.Error())
		}
		msg, err := tcp.RecvFrom(conn)
		if err != nil {
			log.Panicln(err.Error())
		}
		log.Println(fmt.Sprintf("Received msg %s from conn %s", msg.String(), conn.Addr().String()))

		time.Sleep(time.Second * 2)
	}

}
