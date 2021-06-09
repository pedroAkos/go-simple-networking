package main

import (
	"bytes"
	"fmt"
	"github.com/pedroAkos/network/pkg/neti"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type msg struct {
	code uint16
	mark string
	seqnum uint32
}

func (m msg) String() string {
	return fmt.Sprintf("%v{%v-%v%v}", m.Name(),m.code, m.mark, m.seqnum)
}

func (m msg) Name() string {
	return "Msg"
}

func (m msg) Code() uint16 {
	return m.code
}

func (m msg) Serialize() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := neti.EncodeStringToBuffer(m.mark, buf); err != nil {
		return nil, err
	}
	if err := neti.EncodeNumberToBuffer(m.seqnum, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m msg) Deserialize(b []byte) (neti.Message, error) {
	buf := bytes.NewBuffer(b)
	var err error
	if m.mark, err = neti.DecodeStringFromBuffer(buf); err != nil {
		return nil, err
	}
	if err = neti.DecodeNumberFromBuffer(&m.seqnum, buf); err != nil {
		return nil, err
	}
	fmt.Println(m)
	return m, nil
}

func main() {

	listenAddr := pflag.String("listen", "127.0.0.1:10000", "Listen address")
	dstAddr := pflag.String("dst", "", "Dest address")

	pflag.Parse()

	log.SetLevel(log.DebugLevel)

	netserv := neti.InitBaseUdpService(*listenAddr, 1024)

	client1ID := "client1"
	client2ID := "client2"

	client1 := netserv.RegisterListener(client1ID)
	client2 := netserv.RegisterListener(client2ID)

	client1.RegisterMessage(msg{code: 1})
	client1.RegisterMessage(msg{code: 1})
	client1.RegisterMessage(msg{code: 2})
	client1.RegisterMessage(msg{code: 2})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	go clientRcvLoop(client1)
	go clientRcvLoop(client2)

	if *dstAddr != "" {
		go clientSendLoop(client1, *dstAddr, client1ID, "A")
		go clientSendLoop(client1, *dstAddr, client2ID, "B")

		go clientSendLoop(client2, *dstAddr, client1ID, "C")
		go clientSendLoop(client2, *dstAddr, client2ID, "D")
	}


	<-stop

}

func clientRcvLoop(client neti.NetClient) {
	for {
		select {
		case conn := <- client.Accept():
			if m, err := client.RecvFrom(conn); err != nil {
				log.Error(client.ServiceId(), ": ", err)
			} else {
				log.Info(client.ServiceId(),": Received: ", m.String(), " from: ", conn.ServiceId)
				if m.Code() == 1 {
					m2 := m.(msg)
					m2.code = 2
					err := client.SendTo(conn, m2)
					if err != nil {
						log.Error(client.ServiceId(), ": ", err)
					}
				}
			}
		}
	}
}

func clientSendLoop(client neti.NetClient, dst string, dstId string, mark string) {
	i := 0
	for {
		select {
		case <- time.After(1*time.Second):
			if conn, err := client.OpenTo(dst, dstId); err == nil {
				m := msg{1, mark, uint32(i)}
				if err = client.SendTo(conn, m); err != nil {
					log.Error(client.ServiceId(), ": ", err)
				}
				log.Info(client.ServiceId(), ": Sent: ", m, " to: ", dstId)
				i++
			} else {
				log.Error(client.ServiceId(), ": ", err)
			}
		}
	}
}
