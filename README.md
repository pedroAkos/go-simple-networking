# GO simple networking

A simple socket library in GO with multiplexers.

Contains functions for TCP and UDP sockets.

## Basic Usage

```go

import "github.com/pedroAkos/network/pkg/neti"


func main() {
    tcp := neti.newTcpNet(logrus.StandardLogger())

    listener, err := tcp.Listen("0.0.0.0:10000")
    if err != nil { panic(err) }

    tcp.RegisterMessage(ping{})
    tcp.RegisterMessage(pong{})

    go func() {
      for {
        select {
          case conn := <- listener:
            go func() {
              msg, err := tcp.RecvFrom(conn)
              if err != nil {
                 panic(err)
              }
              err := tcp.SendTo(conn, pong{})
              if err != nil {
                panic(err)
              }
            }()
        }
      }
    }()

    conn, err := tcp.Open("127.0.0.1:10000")
    if err != nil { panic(err) }

    err := tcp.SendTo(conn, ping{})
    if err != nil { panic(err) }

    msg, err := tcp.RecvFrom(conn)
    if err != nil { panic(err) }
}
```

## Messages

Messages are structs that implement the Message interface.

```go
type Message interface {
	String() string
	Name() string
	Code() uint16
	Serialize(buff *bytes.Buffer) error
	Deserialize(buff *bytes.Buffer) (Message, error)
}
```

The code is used to identify the message type. This is latter used in the message serialization.

### Example

```go
type ping struct {
    seqnum uint64
}

func (p ping) String() string {
    return "ping"
}

func (p ping) Name() string {
    return "ping"
}

func (p ping) Code() uint16 {
    return 0
}

func (p ping) Serialize(buff *bytes.Buffer) error {
    err := neti.EncodeNumberToBuffer(p.seqnum, buff)
    return err
}

func (p ping) Deserialize(buff *bytes.Buffer) (Message, error) {
    err := neti.DecodeNumberFromBuffer(&p.seqnum, buff)
    return p, err
}
```


## Multiplexers

Multiplexers are used to multiplex messages over a single connection.

```go

func main() {
    netServ := neti.InitBaseTcpService("0.0.0.0:10000", logrus.StandardLogger())

    client1ID := "client1"
    client2ID := "client2"

    client1 := netserv.RegisterListener(client1ID)
    client2 := netserv.RegisterListener(client2ID)

    client1.RegisterMessage(ping{})
    client2.RegisterMessage(ping{})

    client1.RegisterMessage(pong{})
    client2.RegisterMessage(pong{})

    go clientRcvLoop(client1)
    go clientRcvLoop(client2)

   conn, _ := client1.OpenTo("127.0.0.1:10000", client1ID)
   _ := client1.SendTo(conn, ping{})
   conn.Close()

   conn, _ := client1.OpenTo("127.0.0.1:10000", client2ID)
   _ := client1.SendTo(conn, ping{})
   conn.Close()
}

func clientRcvLoop(client *neti.TcpClient) {
    for {
        select {
        case conn := <- client.Accept():
            go func() {
                msg, err := client.RecvFrom(conn)
                if err != nil {
                    panic(err)
                }
                err := client.SendTo(conn, pong{})
                if err != nil {
                    panic(err)
                }
            }()
        }
    }
}
```
